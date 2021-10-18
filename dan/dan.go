package dan

import (
  "bytes"
  "fmt"
  "runtime"
)

var DanHacks = true

type Log interface {
  Trace(s string, a ...interface{})
  Debug(s string, a ...interface{})
  Info(s string, a ...interface{})
  Warn(s string, a ...interface{})
  Error(s string, a ...interface{})
}

const (
  LogTrace = 0
  LogDebug = 1
  LogInfo  = 2
  LogWarn  = 3
  LogError = 4
)

type LogSink interface {
  WriteMsg(
      logName string,
      level int,
      message string,
      functionName string,
      fileName string,
      lineNo int,
  )
}

type Logger struct {
  LogSink
  Name string
  Skip int
}

func (l Logger) log(level int, format string, a ...interface{}) {
  pc, fn, line, _ := runtime.Caller(l.Skip)
  funcName := runtime.FuncForPC(pc).Name()
  var buf bytes.Buffer
  _, err := fmt.Fprintf(&buf, format, a...)
  var msg string
  if err != nil {
    msg = err.Error()
  } else {
    msg = buf.String()
  }
  l.LogSink.WriteMsg(l.Name, level, msg, funcName, fn, line)
  //l.LogSink.WriteMsg(createLogMsg(l.Name, LogDebug, l.Skip, s))
}

func (l Logger) Trace(s string, a ...interface{}) {
  l.log(LogTrace, s, a...)
}

func (l Logger) Debug(s string, a ...interface{}) {
  l.log(LogDebug, s, a...)
}

func (l Logger) Info(s string, a ...interface{}) {
  l.log(LogInfo, s, a...)
}
func (l Logger) Warn(s string, a ...interface{}) {
  l.log(LogWarn, s, a...)
}
func (l Logger) Error(s string, a ...interface{}) {
  l.log(LogError, s, a...)
}

type NoopLogSink struct {
}

func (s NoopLogSink) WriteMsg(logName string,
    level int,
    message string,
    functionName string,
    fileName string,
    lineNo int, ) {
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
    levelName = " INFO"
    color = "32"
  case LogWarn:
    levelName = " WARN"
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

var DanLog = Logger{Name: "GOMOBILE:", Skip: 2, LogSink: StdoutLogSink{UseColor: true}}
