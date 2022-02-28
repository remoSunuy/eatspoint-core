package errorx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path"
	"runtime"
	"sync"
)

type StackTracer interface {
	Error()				string
	StackAddrs()		string
	StackFrames()		[]StackFrame
	GetStack()			string
	GetStackAsJson()	interface{}
}

type StackFrame struct {
	PC			uintptr
	Func 		*runtime.Func
	FuncName	string
	File		string
	LineNumber	int
}

type ErrorStackTrace struct {
	stack		[]uintptr
	framesOnce	sync.Once
	stackFrames	[]StackFrame
}

func NewErrorWithStackTrace(stackLen int, callerLen int) *ErrorStackTrace {
	stack := make([]uintptr, stackLen)
	stackLength := runtime.Callers(callerLen, stack)
	return &ErrorStackTrace{
		stack: stack[:stackLength],
	}
}

func (err *ErrorStackTrace) StackAddrs() string {
	buf := bytes.NewBuffer(make([]byte, 0, len(err.stack) * 8))
	for _, pc := range err.stack {
		fmt.Fprintf(buf, "0x%x ", pc)
	}

	bufBytes := buf.Bytes()
	return string(bufBytes[:len(bufBytes) -1 ])
}

func (err *ErrorStackTrace) StackFrames() []StackFrame {
	err.framesOnce.Do(func() {
		err.stackFrames = make([]StackFrame, len(err.stack))
		for i, pc := range err.stack {
			frame := &err.stackFrames[i]
			frame.PC = pc
			frame.Func = runtime.FuncForPC(pc)
			if frame.Func != nil {
				frame.FuncName = frame.Func.Name()
				frame.File, frame.LineNumber = frame.Func.FileLine(frame.PC -1)
			}
		}
	})
	return err.stackFrames
}

func (err *ErrorStackTrace) GetStack() string {
	stackFrames := err.StackFrames()
	buf := bytes.NewBuffer(make([]byte, 0, 256))
	for _, frame := range stackFrames {
		_, _ = buf.WriteString(frame.FuncName)
		_, _ = buf.WriteString("\n")
		fmt.Fprintf(buf, "\t%s:%d +0x%x\n", frame.File, frame.LineNumber, frame.PC)
	}
	return buf.String()
}

func (err *ErrorStackTrace) GetStackAsJson() interface{} {
	stackFrames := err.stackFrames
	buf := bytes.NewBuffer(make([]byte, 0, 256))
	var (
		data	[]byte
		i 		interface {}
	)

	data = append(data, '[')
	for i, frame := range stackFrames {
		if i != 0 {
			data = append(data, ',')
		}
		name := path.Base(frame.FuncName)
		frameBytes := []byte(
			fmt.Sprintf(`{"filepath": "%s", "name": "%s", "line": %d}`, frame.File, name, frame.LineNumber),
		)
		data = append(data, frameBytes...)
	}
	data = append(data, ']')
	buf.Write(data)
	_ = json.Unmarshal(data, &i)
	return i
}