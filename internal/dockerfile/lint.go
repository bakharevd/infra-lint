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

func Format(data []byte) ([]byte, error) {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	var formatted bytes.Buffer

	for scanner.Scan() {
		line := scanner.Text()

		if strings.TrimSpace(line) == "" {
			formatted.WriteString("\n")
			continue
		}

		trimmed := strings.TrimSpace(line)

		for _, instruction := range []string{"FROM", "RUN", "CMD", "LABEL", "EXPOSE", "ENV", "ADD", "COPY", "ENTRYPOINT", "VOLUME", "USER", "WORKDIR", "ARG", "ONBUILD", "STOPSIGNAL", "HEALTHCHECK", "SHELL"} {
			if strings.HasPrefix(strings.ToUpper(trimmed), instruction+" ") || strings.ToUpper(trimmed) == instruction {
				trimmed = instruction + trimmed[len(instruction):]
				break
			}
		}

		parts := strings.SplitN(trimmed, " ", 2)
		if len(parts) == 2 {
			formatted.WriteString(parts[0] + " " + strings.TrimSpace(parts[1]) + "\n")
		} else {
			formatted.WriteString(trimmed + "\n")
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return formatted.Bytes(), nil
}
