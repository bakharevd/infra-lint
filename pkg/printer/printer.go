package printer

import (
	"fmt"
	"os"
)

var NoColor = false

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
)

type Level int

const (
	LevelOK Level = iota
	LevelWarn
	LevelError
)

func (l Level) String() string {
	switch l {
	case LevelOK:
		return "[OK]"
	case LevelWarn:
		return "[WARN]"
	case LevelError:
		return "[ERROR]"
	default:
		return "[INFO]"
	}
}

func (l Level) Color() string {
	switch l {
	case LevelOK:
		return ColorGreen
	case LevelWarn:
		return ColorYellow
	case LevelError:
		return ColorRed
	default:
		return ColorBlue
	}
}

var (
	errorCount int
	warnCount  int
	checkCount int
)

func Check(condition bool, level Level, format string, args ...interface{}) bool {
	checkCount++

	if condition {
		if level == LevelOK {
			print(level.String(), level.Color(), format, args...)
		}
		return true
	}
	switch level {
	case LevelWarn:
		warnCount++
	case LevelError:
		errorCount++
	}

	print(level.String(), level.Color(), format, args...)
	return false
}

func CheckOK(cond bool, msg string, args ...interface{}) bool {
	return Check(cond, LevelOK, msg, args...)
}

func CheckWarn(cond bool, msg string, args ...interface{}) bool {
	return Check(cond, LevelWarn, msg, args...)
}

func CheckError(cond bool, msg string, args ...interface{}) bool {
	return Check(cond, LevelError, msg, args...)
}

func Error(format string, args ...interface{}) {
	errorCount++
	print("[ERROR]", ColorRed, format, args...)
}

func Warn(format string, args ...interface{}) {
	warnCount++
	print("[WARN]", ColorYellow, format, args...)
}

func Info(format string, args ...interface{}) {
	print("[INFO]", ColorBlue, format, args...)
}

func OK(format string, args ...interface{}) {
	print("[OK]", ColorGreen, format, args...)
}

func Fatal(format string, args ...interface{}) {
	print("[FATAL]", ColorRed, format, args...)
	os.Exit(1)
}

func Stats() (errors int, warnings int, total int, health int) {
	okChecks := checkCount - errorCount - warnCount
	if checkCount == 0 {
		return errorCount, warnCount, 0, 100
	}
	health = (okChecks * 100) / checkCount
	return errorCount, warnCount, checkCount, health
}

func PrintStats() {
	errs, warns, total, health := Stats()

	fmt.Println()
	if NoColor {
		fmt.Printf("Checks: %d | Errors: %d | Warnings: %d | Health: %d%%\n", total, errs, warns, health)
	} else {
		color := ColorGreen
		switch {
		case health < 60:
			color = ColorRed
		case health < 85:
			color = ColorYellow
		}
		fmt.Printf("%s%sChecks: %d | Errors: %d | Warnings: %d | Health: %d%%%s\n", color, "[SUMMARY] ", total, errs, warns, health, ColorReset)
	}
}

func HasErrors() bool {
	return errorCount > 0
}

func print(label string, color string, format string, args ...interface{}) {
	if NoColor {
		fmt.Printf("%s %s\n", label, fmt.Sprintf(format, args...))
	} else {
		fmt.Printf("%s %s\n", colorize(label, color), fmt.Sprintf(format, args...))
	}
}

func colorize(label, color string) string {
	return fmt.Sprintf("%s%s%s", color, label, ColorReset)
}
