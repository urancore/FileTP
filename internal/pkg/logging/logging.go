// autor: https://github.com/urancore
package logging

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time"
)

type Logger struct {
	enableColor bool
}

const (
	colorRed    = "\x1b[31m"
	colorGreen  = "\x1b[32m"
	colorBlue   = "\x1b[34m"
	colorYellow = "\x1b[33m"
	colorPurple = "\x1b[35m"
	colorWhite  = "\x1b[37m"
	colorGray   = "\x1b[90m"
	italicGray  = "\x1b[3;90m" // Курсив + серый цвет
	colorReset  = "\x1b[0m"
)

var levelColors = map[string]string{
	"INFO":  colorGreen,
	"ERROR": colorRed,
	"DEBUG": colorBlue,
}

func NewLogger(enableColor bool) *Logger {
	return &Logger{enableColor: enableColor}
}

func (l *Logger) SetColorEnabled(enabled bool) {
	l.enableColor = enabled
}

func (l *Logger) Info(msg string) {
	l.log("INFO", msg)
}

func (l *Logger) Error(msg string) {
	l.log("ERROR", msg)
}

func (l *Logger) Debug(msg string) {
	l.log("DEBUG", msg)
}

func (l *Logger) log(level string, msg string) {
	fmt.Print(l.formatLog(level, msg))
}

func (l *Logger) formatLog(level string, msg string) string {
	now := time.Now().Format(time.RFC3339)
	file, line := l.getCallerInfo()
	file = formatFilepath(file)

	// Форматирование уровня логирования
	levelPart := level
	if color, ok := levelColors[level]; ok && l.enableColor {
		levelPart = fmt.Sprintf("%s%s%s", color, level, colorReset)
	}

	// Форматирование времени (курсив + серый)
	nowPart := now
	if l.enableColor {
		nowPart = fmt.Sprintf("%s%s%s", italicGray, now, colorReset)
	}

	// Форматирование файла и номера строки (фиолетовый)
	filePart := file
	lineStr := fmt.Sprintf("%d", line)
	if l.enableColor {
		filePart = fmt.Sprintf("%s%s%s", colorPurple, file, colorReset)
		lineStr = fmt.Sprintf("%s:%s%s", colorPurple, lineStr, colorReset)
	}

	// Форматирование сообщения (белый)
	msgPart := msg
	if l.enableColor {
		msgPart = fmt.Sprintf("%s%s%s", colorWhite, msg, colorReset)
	}

	return fmt.Sprintf("%s [%s] %s%s %s\n",
		nowPart,
		levelPart,
		filePart,
		lineStr,
		msgPart,
	)
}

func formatFilepath(path string) string {
	abs_path, err := filepath.Abs(".")
	if err != nil {
		fmt.Printf("%v", err)
	}

	path, err = filepath.Rel(abs_path, path)
	if err != nil {
		fmt.Printf("%v", err)
	}

	path = filepath.ToSlash(filepath.Clean(path))

	return path
}

func (l *Logger) getCallerInfo() (string, int) {
	// 4 - глубина вызова для получения места вызова метода логирования
	_, file, line, ok := runtime.Caller(4)
	if !ok {
		return "", 0
	}
	return file, line
}
