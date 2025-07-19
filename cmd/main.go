package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"infra-lint/internal/compose"
	"infra-lint/internal/dockerfile"
	"infra-lint/internal/jenkinsfile"
	"infra-lint/internal/nginx"
	"infra-lint/pkg/printer"
)

var (
	includedTypes []string
	excludedTypes []string
)

type linterInfo struct {
	name        string
	description string
	checkFunc   func(string) bool
}

var availableLinters = []linterInfo{
	{"docker", "Dockerfile", dockerfile.IsDockerfile},
	{"compose", "docker-compose.yml/yaml", compose.IsDockerCompose},
	{"jenkins", "Jenkinsfile", jenkinsfile.IsJenkinsfile},
	{"nginx", "nginx.conf and .vhost files", nginx.IsNginxConfig},
}

func printAvailableTypes() {
	fmt.Println("Available linter types:")
	fmt.Println()
	for _, linter := range availableLinters {
		fmt.Printf("  %-10s - %s\n", linter.name, linter.description)
	}
}

func main() {
	filePath := flag.String("file", "", "Path to the configuration file")
	dirPath := flag.String("dir", "", "Path to directory with configuration files")
	repoURL := flag.String("repo", "", "Git repository URL to clone and scan")
	includeTypes := flag.String("include-types", "", "Comma-separated list of linters to use")
	excludeTypes := flag.String("exclude-types", "", "Comma-separated list of linters to exclude")
	listTypes := flag.Bool("list-types", false, "List all available linter types")
	noColor := flag.Bool("no-color", false, "Disable color output")
	formatter := flag.Bool("formatter", false, "Format files instead of linting")
	flag.Parse()

	if *listTypes {
		printAvailableTypes()
		return
	}

	if *filePath == "" && *dirPath == "" && *repoURL == "" {
		log.Fatal("Provide path to file with -file, directory with -dir, or repository with -repo")
	}

	printer.NoColor = *noColor

	if *includeTypes != "" {
		includedTypes = strings.Split(strings.ToLower(*includeTypes), ",")
		for i := range includedTypes {
			includedTypes[i] = strings.TrimSpace(includedTypes[i])
		}
	}

	if *excludeTypes != "" {
		excludedTypes = strings.Split(strings.ToLower(*excludeTypes), ",")
		for i := range excludedTypes {
			excludedTypes[i] = strings.TrimSpace(excludedTypes[i])
		}
	}

	if *repoURL != "" {
		tmpDir, err := os.MkdirTemp("", "infra-lint-repo-*")
		if err != nil {
			log.Fatalf("Failed to create temp directory: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		log.Printf("Cloning repository %s...", *repoURL)
		cmd := exec.Command("git", "clone", *repoURL, tmpDir)
		if err := cmd.Run(); err != nil {
			log.Fatalf("Failed to clone repository: %v", err)
		}

		log.Printf("Scanning repository...")
		if *formatter {
			if err := scanDirAndFormat(tmpDir); err != nil {
				log.Fatalf("Failed to scan repository: %v", err)
			}
		} else {
			if err := scanDirAndLint(tmpDir); err != nil {
				log.Fatalf("Failed to scan repository: %v", err)
			}
			printer.PrintStats()
		}
		return
	}

	if *dirPath != "" {
		if *formatter {
			err := scanDirAndFormat(*dirPath)
			if err != nil {
				log.Fatalf("Failed to scan directory: %v", err)
			}
		} else {
			err := scanDirAndLint(*dirPath)
			if err != nil {
				log.Fatalf("Failed to scan directory: %v", err)
			}
			printer.PrintStats()
		}
		return
	}

	data, err := os.ReadFile(*filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	if *formatter {
		if err := formatFile(*filePath, data); err != nil {
			log.Fatalf("Failed to format file: %v", err)
		}
	} else {
		lintFile(*filePath, data)
		printer.PrintStats()
	}
}

func shouldRunLinter(linterType string) bool {
	if len(excludedTypes) > 0 {
		for _, excluded := range excludedTypes {
			if excluded == linterType {
				return false
			}
		}
	}

	if len(includedTypes) > 0 {
		for _, included := range includedTypes {
			if included == linterType {
				return true
			}
		}
		return false
	}

	return true
}

func lintFile(filePath string, data []byte) {
	switch {
	case dockerfile.IsDockerfile(filePath) && shouldRunLinter("docker"):
		log.Printf("Linting Dockerfile: %s", filePath)
		dockerfile.Lint(data)
	case compose.IsDockerCompose(filePath) && shouldRunLinter("compose"):
		log.Printf("Linting Docker Compose: %s", filePath)
		compose.Lint(data)
	case jenkinsfile.IsJenkinsfile(filePath) && shouldRunLinter("jenkins"):
		log.Printf("Linting Jenkinsfile: %s", filePath)
		jenkinsfile.Lint(data)
	case nginx.IsNginxConfig(filePath) && shouldRunLinter("nginx"):
		log.Printf("Linting Nginx config: %s", filePath)
		nginx.Lint(data)
	}
}

func scanDirAndLint(dir string) error {
	return walkDir(dir)
}

func formatFile(filePath string, data []byte) error {
	var formatted []byte
	var err error

	switch {
	case dockerfile.IsDockerfile(filePath):
		log.Printf("Formatting Dockerfile: %s", filePath)
		formatted, err = dockerfile.Format(data)
	case compose.IsDockerCompose(filePath):
		log.Printf("Formatting Docker Compose: %s", filePath)
		formatted, err = compose.Format(data)
	case jenkinsfile.IsJenkinsfile(filePath):
		log.Printf("Formatting Jenkinsfile: %s", filePath)
		formatted, err = jenkinsfile.Format(data)
	case nginx.IsNginxConfig(filePath):
		log.Printf("Formatting Nginx config: %s", filePath)
		formatted, err = nginx.Format(data)
	default:
		log.Printf("Unknown file type: %s", filePath)
		return nil
	}

	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, formatted, 0644)
	if err != nil {
		return err
	}

	printer.OK("Formatted: %s", filePath)
	return nil
}

func scanDirAndFormat(dir string) error {
	return walkDirAndFormat(dir)
}

func walkDirAndFormat(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		path := dir + string(os.PathSeparator) + entry.Name()
		if entry.IsDir() {
			if err := walkDirAndFormat(path); err != nil {
				return err
			}
			continue
		}

		shouldProcess := false
		if dockerfile.IsDockerfile(path) && shouldRunLinter("docker") {
			shouldProcess = true
		} else if compose.IsDockerCompose(path) && shouldRunLinter("compose") {
			shouldProcess = true
		} else if jenkinsfile.IsJenkinsfile(path) && shouldRunLinter("jenkins") {
			shouldProcess = true
		} else if nginx.IsNginxConfig(path) && shouldRunLinter("nginx") {
			shouldProcess = true
		}

		if shouldProcess {
			data, err := os.ReadFile(path)
			if err != nil {
				log.Printf("Failed to read file %s: %v", path, err)
				continue
			}
			if err := formatFile(path, data); err != nil {
				log.Printf("Failed to format file %s: %v", path, err)
			}
		}
	}
	return nil
}

func walkDir(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		path := dir + string(os.PathSeparator) + entry.Name()
		if entry.IsDir() {
			if err := walkDir(path); err != nil {
				return err
			}
			continue
		}

		shouldProcess := false
		if dockerfile.IsDockerfile(path) && shouldRunLinter("docker") {
			shouldProcess = true
		} else if compose.IsDockerCompose(path) && shouldRunLinter("compose") {
			shouldProcess = true
		} else if jenkinsfile.IsJenkinsfile(path) && shouldRunLinter("jenkins") {
			shouldProcess = true
		} else if nginx.IsNginxConfig(path) && shouldRunLinter("nginx") {
			shouldProcess = true
		}

		if shouldProcess {
			data, err := os.ReadFile(path)
			if err != nil {
				log.Printf("Failed to read file %s: %v", path, err)
				continue
			}
			lintFile(path, data)
		}
	}
	return nil
}
