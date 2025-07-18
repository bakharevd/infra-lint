package jenkinsfile

import (
	"bufio"
	"bytes"
	"infra-lint/pkg/printer"
	"strings"
)

func IsJenkinsfile(path string) bool {
	return strings.HasSuffix(path, "Jenkinsfile")
}

func Lint(data []byte) {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	hasPostBlock := false
	hasTimeout := false
	stageCount := 0

	lineNo := 0

	for scanner.Scan() {
		lineNo++
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "stage(") || strings.HasPrefix(line, "stage ") {
			stageCount++
		}
		if strings.Contains(line, "timeout(") {
			hasTimeout = true
		}
		if strings.Contains(line, "post {") {
			hasPostBlock = true
		}
	}

	printer.CheckWarn(hasTimeout, "Jenkinsfile: No timeout() specified in pipeline. Risk of hanging builds.")
	printer.CheckError(hasPostBlock, "Jenkinsfile: Missing post block to handle failures/cleanup.")
	printer.CheckWarn(stageCount <= 10, "Jenkinsfile: Too many stages (>10). Consider modularizing.")
}
