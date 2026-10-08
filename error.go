package fail

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/codeus-node/icons"
)

const (
	errorLabel   = "AP_ERROR"
	warningLabel = "AP_WARNING"
)

// Error is a custom error type that can be an Error or a Warning based on the ExitCode
type Error struct {
	ExitCode int
	Message  string
	Stack    string
	Icon     string
}

// ==========================================
// Error Interface Implementation
// ==========================================

// Error returns a string representation of the Error
func (e *Error) Error() string {
	label := errorLabel
	if e.ExitCode == -2 {
		label = warningLabel
	}
	icon := ""
	if len(e.Icon) > 0 {
		icon = fmt.Sprintf("%s", e.Icon)
	} else if e.ExitCode == -2 {
		icon = icons.Warning
	} else {
		icon = icons.Error
	}
	result := fmt.Sprintf("%s[%s]:", icon, label)
	if len(e.Message) > 0 {
		result = fmt.Sprintf("%s %s", result, e.Message)
	}
	if len(e.Stack) > 0 {
		result = fmt.Sprintf("%s\n%s", result, e.Stack)
	}
	return result
}

// ==========================================
// Stringer Interface Implementation
// ==========================================

func (e *Error) String() string {
	return e.Error()
}

// ==========================================
// CustomError Interface Implementation
// ==========================================

// GetExitCode returns the exit code of the Error
func (e *Error) GetExitCode() int {
	return e.ExitCode
}

// GetMessage returns the message of the Error
func (e *Error) GetMessage() string {
	return e.Message
}

// GetStack returns the stacktrace of the Error
func (e *Error) GetStack() string {
	return e.Stack
}

// AsWarning changes the exit code of the Error to -2
func (e *Error) AsWarning() CustomError {
	e.ExitCode = -2
	return e
}

func (e *Error) WithIcon(icon string) CustomError {
	e.Icon = icon
	return e
}

// AsError changes the exit code of the Error to -1
func (e *Error) AsError() CustomError {
	e.ExitCode = -1
	return e
}

func newError(err error, createStack bool, additional ...string) *Error {
	message := ""
	if err != nil {
		message = fmt.Sprintf("%s", err.Error())
	}
	for _, addi := range additional {
		message += fmt.Sprintf("\n   > %s", addi)
	}

	stack := ""
	if createStack {
		stack = takeStacktrace(3)
	}
	return &Error{
		ExitCode: -1,
		Message:  message,
		Stack:    stack,
	}
}

// takeStacktrace returns the stacktrace and removes the first entries based on the skip parameter
func takeStacktrace(skip int) string {
	pcs := make([]uintptr, 128)

	n := runtime.Callers(skip+1, pcs)
	if n == 0 {
		return ""
	}

	pcs = pcs[:n]
	var sb strings.Builder
	frames := runtime.CallersFrames(pcs)

	for {
		frame, more := frames.Next()
		_, err := fmt.Fprintf(&sb, "%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line)
		if err != nil {
			break
		}
		if !more {
			break
		}
	}

	return sb.String()
}
