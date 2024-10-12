package logger

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Logger struct {
	debugEnabled bool
	infoEnabled  bool
	warnEnabled  bool
	errorEnabled bool
	colorEnabled bool
}

//var instance *Logger
//var once sync.Once

// ANSI color codes
const (
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorReset  = "\033[0m"
)

var (
	once sync.Once
	// Global instance accessible from outside
	Console *Logger
)

func InitConsoleLogger(debug, info, warn, error, color bool) {
	once.Do(func() {
		Console = &Logger{
			debugEnabled: debug,
			infoEnabled:  info,
			warnEnabled:  warn,
			errorEnabled: error,
			colorEnabled: color,
		}
	})
}

// SetDebug 设置 debug 日志状态
func (l *Logger) SetDebug(enabled bool) {
	l.debugEnabled = enabled
}

// SetInfo 设置 info 日志状态
func (l *Logger) SetInfo(enabled bool) {
	l.infoEnabled = enabled
}

// SetWarn 设置 warn 日志状态
func (l *Logger) SetWarn(enabled bool) {
	l.warnEnabled = enabled
}

// SetError 设置 error 日志状态
func (l *Logger) SetError(enabled bool) {
	l.errorEnabled = enabled
}

// SetColor 设置彩色输出状态
func (l *Logger) SetColor(enabled bool) {
	l.colorEnabled = enabled
}

func (l *Logger) log(level string, color string, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	if l.colorEnabled && color != "" {
		fmt.Printf("%s%s| %s%s\n", color, level, message, colorReset)
	} else {
		fmt.Printf("%s| %s\n", level, message)
	}
}

func (l *Logger) Debug(format string, args ...interface{}) {
	if l.debugEnabled {
		l.logWithDetails("DEBUG", colorBlue, format, args...)
	}
}

func (l *Logger) Info(format string, args ...interface{}) {
	if l.infoEnabled {
		l.logWithDetails("INFO", "", format, args...)
	}
}

func (l *Logger) Warn(format string, args ...interface{}) {
	if l.warnEnabled {
		l.logWithDetails("WARN", colorYellow, format, args...)
	}
}

func (l *Logger) Error(format string, args ...interface{}) {
	if l.errorEnabled {
		l.logWithDetails("ERROR", colorRed, format, args...)
	}
}

// logWithDetails 输出日志信息，包括调用者的详细信息
func (l *Logger) logWithDetails(level string, color string, format string, args ...interface{}) {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	function := runtime.FuncForPC(pc).Name()
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	// 打印实例如下
	/**
	[DEBUG]|2024-05-11 11:04:29|/Users/zqzess/Code/iwhalecloud/whale-term/backend/sftp/sftp.go:636|whale-term/backend/sftp.copyFileOrDirectory.func1 |
	*/
	// 切割路径，只取最后2位
	parts := strings.Split(file, "/")
	lastTwoParts := parts[len(parts)-2:]
	twoPartString := lastTwoParts[0] + "/" + lastTwoParts[1]

	// 切割取最后一位
	parts = strings.Split(function, "/")
	lastParts := parts[len(parts)-1]
	//details := fmt.Sprintf("[%s]|%s|%s:%d|%s", level, timestamp, lastTwoParts, line, lastParts)
	details := fmt.Sprintf("[%s]|%s|%s:%d|%s ", level, timestamp, twoPartString, line, lastParts)
	l.log(details, color, format, args...)
}
