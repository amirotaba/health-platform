package domain

type Logger interface {
	LogError(string, ...interface{})
	LogAccess(string, ...interface{})
}
