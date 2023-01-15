package logger

type Logger interface {
	Fatal(string, ...interface{})
	Warn(string, ...interface{})
	Error(string, ...interface{})
	Info(string, ...interface{})
	Debug(string, ...interface{})
	Trace(string, ...interface{})
}