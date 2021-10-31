package klog

import (
  "bytes"
  "fmt"
  "runtime"
)

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

/*func createLogMsg(logName string, level int, skip int, format string, a ...interface{}) LogMsg {
  pc, fn, line, _ := runtime.Caller(skip)
  funcName := runtime.FuncForPC(pc).Name()
  var buf bytes.Buffer
  _, err := fmt.Fprintf(&buf, format, a...)
  var msg string
  if err != nil {
    msg = err.Error()
  } else {
    msg = buf.String()
  }
  return LogMsg{LogName: logName, Level: level, Message: msg, FunctionName: funcName, FileName: fn, LineNo: line}
}

func log(name string, level string, color int, skip int, msg string) {
  pc, fn, line, _ := runtime.Caller(skip)
  fmt.Printf("\u001b[0;%dm%s:%s %s[%s:%d]: %s\u001b[0m\n", color, level, name, runtime.FuncForPC(pc).Name(), fn, line, msg)
}
*/
