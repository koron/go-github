package github

type Logger interface {
	Printf(format string, v ...interface{})
}
