package klog

import (
	"bytes"
	"fmt"
)

type NoopLogSink struct {
}

func (s NoopLogSink) WriteMsg(logName string,
	level int,
	message string,
	functionName string,
	fileName string,
	lineNo int) {
}

type NoopLogger struct {
	Name string
	Skip int
}

func (l NoopLogger) Trace(s string) {}
func (l NoopLogger) Debug(s string) {}
func (l NoopLogger) Info(s string)  {}
func (l NoopLogger) Warn(s string)  {}
func (l NoopLogger) Error(s string) {}

func formatMessage(logName string, level int, useColor bool,
	message string,
	functionName string,
	fileName string,
	lineNo int) string {
	var color string
	var levelName string
	switch level {
	case LogTrace:
		levelName = "TRACE"
		color = "35"
	case LogDebug:
		levelName = "DEBUG"
		color = "36"
	case LogInfo:
		levelName = "INFO "
		color = "32"
	case LogWarn:
		levelName = "WARN "
		color = "33"
	case LogError:
		levelName = "ERROR"
		color = "31"
	}
	var buf bytes.Buffer
	_, err := fmt.Fprintf(&buf,
		"%s:%s %s[%s:%d]: %s",
		levelName, logName, functionName, fileName, lineNo, message)
	if err != nil {
		return err.Error()
	}
	if useColor {
		return "\u001b[0;" + color + "m" + buf.String() + "\u001b[0m"
	} else {
		return buf.String()
	}
}

type StdoutLogSink struct {
	UseColor bool
}

func (s StdoutLogSink) WriteMsg(
	logName string,
	level int,
	message string,
	functionName string,
	fileName string,
	lineNo int,
) {
	println(formatMessage(logName, level, s.UseColor, message, functionName, fileName, lineNo))
}

var KLog = Logger{Name: "KIPFS", Skip: 2, LogSink: StdoutLogSink{UseColor: true}}
