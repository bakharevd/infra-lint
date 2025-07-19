package nginx

import (
	"bufio"
	"bytes"
	"infra-lint/pkg/printer"
	"strings"
)

func IsNginxConfig(path string) bool {
	return strings.HasSuffix(path, ".conf") || strings.HasSuffix(path, ".vhost")
}

func Lint(data []byte) {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	foundSecurityHeaders := false
	foundServerTokensOff := false
	foundGzip := false

	lineNo := 0

	for scanner.Scan() {
		lineNo++
		line := strings.TrimSpace(scanner.Text())

		if strings.Contains(line, "server_tokens off") {
			foundServerTokensOff = true
		}
		if strings.Contains(line, "add_header") && strings.Contains(line, "X-Content-Type-Options") {
			foundSecurityHeaders = true
		}
		if strings.HasPrefix(line, "gzip") {
			foundGzip = true
		}
		if strings.HasPrefix(line, "access_log") && strings.Contains(line, "off") {
			printer.CheckWarn(false, "Line %d: Logging disabled. Useful for performance, but be careful in prod.", lineNo)
		}
	}

	printer.CheckError(foundServerTokensOff, "nginx.conf: Missing 'server_tokens off' for security.")
	printer.CheckWarn(foundSecurityHeaders, "nginx.conf: Missing security headers like 'X-Content-Type-Options'.")
	printer.CheckWarn(foundGzip, "nginx.conf: Gzip not enabled.")
}

func Format(data []byte) ([]byte, error) {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	var formatted bytes.Buffer
	indentLevel := 0

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if trimmed == "" {
			formatted.WriteString("\n")
			continue
		}

		if strings.HasPrefix(trimmed, "}") {
			indentLevel--
			if indentLevel < 0 {
				indentLevel = 0
			}
		}

		indent := strings.Repeat("    ", indentLevel)
		formatted.WriteString(indent + trimmed + "\n")

		if strings.HasSuffix(trimmed, "{") {
			indentLevel++
		}

		if !strings.HasSuffix(trimmed, "{") && !strings.HasSuffix(trimmed, "}") && !strings.HasSuffix(trimmed, ";") && !strings.HasPrefix(trimmed, "#") {
			//
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return formatted.Bytes(), nil
}
