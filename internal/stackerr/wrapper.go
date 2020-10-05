package stackerr

import (
	"fmt"
	"strings"
)

type StackFrame interface {
	String() string
	InnerMessage() string
}

type StackFrameNoRef struct {
	msg         string
	methodName  string
	fileName    string
	packageName string
}

var _ = StackFrame(&StackFrameNoRef{})

type stackFrameWithRef struct {
	StackFrameNoRef
	reference string
}

func (s stackFrameWithRef) String() string {
	return fmt.Sprintf("%v/%v:%v (%v)", s.packageName, s.fileName, s.methodName, s.reference)
}

func (s stackFrameWithRef) InnerMessage() string {
	return s.msg
}

var _ = StackFrame(&stackFrameWithRef{})

func (s StackFrameNoRef) WithRef(ref string) StackFrame {
	return stackFrameWithRef{
		StackFrameNoRef: s,
		reference:       ref,
	}
}

func (s StackFrameNoRef) InnerMessage() string {
	return s.msg
}

func (s StackFrameNoRef) String() string {
	return fmt.Sprintf("%v/%v:%v", s.packageName, s.fileName, s.methodName)
}

type StackTrace []StackFrame

type HandledError interface {
	Inner() error
	Error() string
	StackTrace() StackTrace
	StackTraceString() string
}

type handledError struct {
	msg   string
	err   error
	stack StackTrace
}

// check it implements required interfaces
var _ = HandledError(&handledError{})
var _ = error(&handledError{})

func (h handledError) Inner() error {
	return h.err
}

func (h handledError) Error() string {
	msg := ""
	if h.msg != "" {
		msg = h.msg + ": "
	}
	msg += h.err.Error()
	msg += fmt.Sprintf("\n  at\n%v", h.StackTraceString())
	return msg
}

func (h handledError) StackTraceString() string {
	var b strings.Builder
	for i, s := range h.stack {
		app := s.InnerMessage()
		if app != "" {
			app += ":\n  "
		}
		_, _ = b.WriteString(fmt.Sprintf("%s[%02d] %v\n", app, len(h.stack)-i-1, s))
	}
	return b.String()
}

func (h handledError) StackTrace() StackTrace {
	return h.stack
}
