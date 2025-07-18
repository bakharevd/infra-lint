package compose

import (
	"strings"

	"gopkg.in/yaml.v3"

	"infra-lint/pkg/printer"
)

func IsDockerCompose(path string) bool {
	return strings.HasSuffix(path, "docker-compose.yml") || strings.HasSuffix(path, "docker-compose.yaml")
}

type ComposeFile struct {
	Services map[string]Service `yaml:"services"`
}

type Service struct {
	Ports   []string               `yaml:"ports"`
	Restart string                 `yaml:"restart"`
	Deploy  map[string]interface{} `yaml:"deploy"`
}

func Lint(data []byte) {
	var compose ComposeFile
	err := yaml.Unmarshal(data, &compose)
	if err != nil {
		printer.Fatal("Failed to parse docker-compose file: %v", err)
	}

	for name, service := range compose.Services {
		for _, port := range service.Ports {
			printer.CheckError(
				!(strings.HasPrefix(port, "0.0.0.0") || strings.HasPrefix(port, "80:") || strings.HasPrefix(port, "443:")),
				"[%s] Port %s is publicly exposed. Consider using internal networking.", name, port)
		}

		printer.CheckWarn(service.Restart != "always", "[%s] Using 'restart: always'. Ensure this is intended.", name)

		printer.CheckError(service.Deploy != nil && service.Deploy["resources"] != nil,
			"[%s] Missing resource limits under 'deploy.resources'.", name)
	}
}
