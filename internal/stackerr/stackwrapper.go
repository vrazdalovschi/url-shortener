package stackerr

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func Wrap(err error) HandledError {
	return internalWrapWithMessage(err, "")
}

func WrapWithMessage(err error, msg string) HandledError {
	return internalWrapWithMessage(err, msg)
}

func WrapWithMessagef(err error, format string, a ...interface{}) HandledError {
	return internalWrapWithMessage(err, fmt.Sprintf(format, a...))
}

func internalWrapWithMessage(err error, msg string) HandledError {
	if err == nil {
		return nil
	}
	pc, file, line, _ := runtime.Caller(2)

	function := runtime.FuncForPC(pc).Name()
	packageName := filepath.Dir(function)
	packageRefactor := strings.ReplaceAll(packageName, "\\", "/")
	methodName := filepath.Base(function)
	packagePart := strings.Split(methodName, ".")[0]
	stackFrame := StackFrameNoRef{
		methodName:  filepath.Ext(methodName),
		fileName:    filepath.Base(file) + ":" + strconv.Itoa(line),
		packageName: packageRefactor + "/" + packagePart,
	}
	handledErr, handled := err.(handledError)
	if handled {
		stackFrame.msg = msg
		handledErr = handledError{
			msg:   handledErr.msg,
			err:   handledErr.err,
			stack: append(handledErr.stack, stackFrame),
		}
	} else {
		handledErr = handledError{
			msg:   msg,
			err:   err,
			stack: []StackFrame{stackFrame},
		}
	}
	return handledErr
}
