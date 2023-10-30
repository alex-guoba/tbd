package logger

import (
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/viper"
)

type Level int8

type Fields map[string]interface{}

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

var colorsOutput = map[string]func(format string, a ...interface{}) string{
	"white":        color.WhiteString,
	"green":        color.GreenString,
	"yellow":       color.YellowString,
	"red":          color.RedString,
	"blue":         color.BlueString,
	"magenta":      color.MagentaString,
	"cyan":         color.CyanString,
	"white_bold":   color.New(color.FgWhite, color.Bold).SprintfFunc(),
	"green_bold":   color.New(color.FgGreen, color.Bold).SprintfFunc(),
	"yellow_bold":  color.New(color.FgYellow, color.Bold).SprintfFunc(),
	"red_bold":     color.New(color.FgRed, color.Bold).SprintfFunc(),
	"blue_bold":    color.New(color.FgBlue, color.Bold).SprintfFunc(),
	"magenta_bold": color.New(color.FgMagenta, color.Bold).SprintfFunc(),
	"cyan_bold":    color.New(color.FgCyan, color.Bold).SprintfFunc(),
}

var gl *Logger

func init() {
	// viper not loaded, so this would not work, you should set level mannually in main
	// level := viper.GetInt("log.level")
	gl = NewLogger("", LevelDebug)
}

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	case LevelFatal:
		return "fatal"
	}
	return ""
}

type Logger struct {
	level  Level
	logger *log.Logger
}

func NewLogger(prefix string, level Level) *Logger {
	l := log.New(os.Stdout, prefix, 0)
	return &Logger{logger: l, level: level}
}

func (l *Logger) outWithColor(level Level, content string) {
	switch level {
	case LevelDebug:
		if l.level <= LevelDebug {
			l.logger.Print(color.WhiteString(content))
		}

	case LevelInfo:
		if l.level <= LevelInfo {
			l.logger.Print(color.GreenString(content))
		}

	case LevelWarn:
		if l.level <= LevelWarn {
			l.logger.Print(color.YellowString(content))
		}

	case LevelError:
		if l.level <= LevelError {
			l.logger.Print(color.RedString(content))
		}

	case LevelFatal:
		if l.level <= LevelFatal {
			l.logger.Fatal(color.HiRedString(content))
		}
	}
}

func SetLevel(level int) {
	gl.SetLevel(level)
}

func (l *Logger) SetLevel(level int) {
	l.level = Level(level)
}

func Debug(v ...interface{}) { gl.Debug(v...) }

func (l *Logger) Debug(v ...interface{}) {
	l.outWithColor(LevelDebug, fmt.Sprint(v...))
}

func Debugf(format string, v ...interface{}) { gl.Debugf(format, v...) }

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.outWithColor(LevelDebug, fmt.Sprintf(format, v...))
}

func Info(v ...interface{}) { gl.Info(v...) }

func (l *Logger) Info(v ...interface{}) {
	l.outWithColor(LevelInfo, fmt.Sprint(v...))
}

func Infof(format string, v ...interface{}) { gl.Infof(format, v...) }

func (l *Logger) Infof(format string, v ...interface{}) {
	l.outWithColor(LevelInfo, fmt.Sprintf(format, v...))
}

func Warn(v ...interface{}) { gl.Warn(v...) }

func (l *Logger) Warn(v ...interface{}) {
	l.outWithColor(LevelWarn, fmt.Sprint(v...))
}

func Warnf(format string, v ...interface{}) { gl.Warnf(format, v...) }

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.outWithColor(LevelWarn, fmt.Sprintf(format, v...))
}

func Error(v ...interface{}) { gl.Error(v...) }

func (l *Logger) Error(v ...interface{}) {
	l.outWithColor(LevelError, fmt.Sprint(v...))
}

func Errorf(format string, v ...interface{}) { gl.Errorf(format, v...) }

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.outWithColor(LevelError, fmt.Sprintf(format, v...))
}

func Fatal(v ...interface{}) { gl.Fatal(v...) }

func (l *Logger) Fatal(v ...interface{}) {
	l.outWithColor(LevelFatal, fmt.Sprint(v...))
}

func Fatalf(format string, v ...interface{}) { gl.Fatalf(format, v...) }

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.outWithColor(LevelFatal, fmt.Sprintf(format, v...))
}

// replay chat message with specified formats
func ChatReplay(format string, v ...interface{}) { gl.ChatReplay(format, v...) }

func (l *Logger) ChatReplay(format string, v ...interface{}) {
	clr := viper.GetString("chat.color")

	// Note: sys logger will add a newline Automatically, that's not friendly in streaming dialogs. so I replace logger with fmt lib.
	if output, ok := colorsOutput[clr]; ok {
		fmt.Print(output(fmt.Sprintf(format, v...)))
		// l.logger.Print(prompt, msg)
	} else {
		// default
		// l.logger.Print(prompt, color.MagentaString(fmt.Sprintf(format, v...)))
		fmt.Print(color.MagentaString(fmt.Sprintf(format, v...)))
	}
}

func ChatReplayNewline() { gl.ChatReplayNewline() }
func (l *Logger) ChatReplayNewline() {
	fmt.Println("")
}

func ChatReplayPrompt() { gl.ChatReplayPrompt() }
func (l *Logger) ChatReplayPrompt() {
	prompt := viper.GetString("chat.prompt")

	if len(prompt) > 0 {
		fmt.Print(color.CyanString(prompt))
	}

}
