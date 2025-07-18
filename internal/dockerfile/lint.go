package dockerfile

import (
	"bufio"
	"bytes"
	"infra-lint/pkg/printer"
	"strings"
)

func IsDockerfile(path string) bool {
	return strings.HasSuffix(path, "Dockerfile") || strings.Contains(path, "Dockerfile")
}

func Lint(data []byte) {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	hasHealthCheck := false
	hasUser := false
	lineNo := 0

	for scanner.Scan() {
		lineNo++
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "FROM") {
			printer.CheckWarn(strings.Contains(line, ":latest"), "Line %d: Avoid using 'latest' tag in FROM", lineNo)
		}

		if strings.HasPrefix(line, "RUN") {
			printer.CheckWarn(strings.Contains(line, "&&"), "Line %d: Combined RUN commands with '&&' to reduce layers", lineNo)
		}

		if strings.HasPrefix(line, "HEALTHCHECK") {
			hasHealthCheck = true
		}

		if strings.HasPrefix(line, "USER") {
			hasUser = true
		}
	}

	printer.CheckWarn(hasHealthCheck, "Missing HEALTHCHECK instruction")
	printer.CheckError(hasUser, "No USER specified. Container runs as root")
}
