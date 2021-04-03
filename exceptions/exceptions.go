package exceptions

import "fmt"

type NotFound struct {
	msg string
}

func (a NotFound) Error() string {
	return a.msg
}

func NewNotFound(format string, v ...interface{}) *NotFound {
	return &NotFound{msg: fmt.Sprintf(format, v...)}
}

// ===================

type Conflict struct {
	msg string
}

func (a Conflict) Error() string {
	return a.msg
}

func NewConflict(format string, v ...interface{}) *Conflict {
	return &Conflict{msg: fmt.Sprintf(format, v...)}
}
