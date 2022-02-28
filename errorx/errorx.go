package errorx

import (
	"bytes"
	"fmt"
)

type ErrorFormatter	interface {
	GetArgs()			[]interface{}
	GetMessage()		string
	FormattedMessage()	string
}

type ErrorX struct {
	messageFormat	string
	cause 			error
	args			[]interface{}
}

func NewErrorX(messageFormat string, args ...interface{}) *ErrorX {
	err := &ErrorX{
		messageFormat: 	messageFormat,
		args: 			args,
	}
	return err
}

func (e *ErrorX) Message() string {
	return e.messageFormat
}

func (e *ErrorX) GetArgs() []interface{} {
	return e.args
}

func (e *ErrorX) GetMessage() string {
	return e.messageFormat
}

func (e *ErrorX) Wrap(err error) error {
	e.cause = err
	return e
}

func (e *ErrorX) Error() string {
	return fmt.Sprintf(e.messageFormat, e.args...)
}

func (w *ErrorX) Cause() error { return w.cause }

func GetErrorMessages(e error) string {
	return extractFullErrorMessage(e, false)
}

func GetErrorMessagesWithStack(e error) string {
	return extractFullErrorMessage(e, true)
}

func extractFullErrorMessage(err error, includeStack bool) string {
	type causer interface {
		Causer()	error
	}
	var ok 			bool
	var lastClErr	error
	errMsg := bytes.NewBuffer(make([]byte, 0, 1024))
	dbxErr := err
	for {
		_, ok := dbxErr.(StackTracer)
		if ok  {
			lastClErr = dbxErr
		}
		errorWithFormat, ok := dbxErr.(ErrorFormatter)
		if ok {
			errMsg.WriteString(errorWithFormat.FormattedMessage())
		}
		errorCauser, ok := dbxErr.(causer)
		if ok {
			innerErr := errorCauser.Causer()
			if innerErr == nil {
				break
			}
			dbxErr = innerErr
		} else {
			// We have reached the end and traveresed all inner errors.
			// Add last message and exit loop.
			errMsg.WriteString(dbxErr.Error())
			break
		}
		errMsg.WriteString(", ")
	}
	stackError, ok := lastClErr.(StackTracer)
	if includeStack && ok {
		errMsg.WriteString("\nSTACK TRACE:\n")
		errMsg.WriteString(stackError.GetStack())
	}
	return errMsg.String()
}

func Causer(err error) error {
	type causer interface {
		Causer() error
	}

	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Causer()
	}
	return err
}