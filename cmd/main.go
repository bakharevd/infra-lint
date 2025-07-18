package main

import (
	"flag"
	"log"
	"os"

	"infra-lint/internal/compose"
	"infra-lint/internal/dockerfile"
	"infra-lint/internal/jenkinsfile"
	"infra-lint/internal/nginx"
	"infra-lint/pkg/printer"
)

func main() {
	filePath := flag.String("file", "", "Path to the configuration file")
	noColor := flag.Bool("no-color", false, "Disable color output")
	flag.Parse()

	if *filePath == "" {
		log.Fatal("Provide path to file with -file")
	}

	printer.NoColor = *noColor

	data, err := os.ReadFile(*filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}
	switch {
	case dockerfile.IsDockerfile(*filePath):
		dockerfile.Lint(data)
	case compose.IsDockerCompose(*filePath):
		compose.Lint(data)
	case jenkinsfile.IsJenkinsfile(*filePath):
		jenkinsfile.Lint(data)
	case nginx.IsNginxConfig(*filePath):
		nginx.Lint(data)
	default:
		log.Fatalf("Unknown file type for: %s", *filePath)
	}
	printer.PrintStats()
}
