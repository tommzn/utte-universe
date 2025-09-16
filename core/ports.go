package core

type Random interface {
	Seek() float64
	Of(n int) int
	OfRange(min, max int) int
	OfIntRange(r intRange) int
}

type Log interface {
	Error(message string, v ...interface{})
	Info(message string, v ...interface{})
	Debug(message string, v ...interface{})
}
