package fail

import (
	"encoding/json"
	"errors"
	"fmt"
)

// CustomError is a custom error type that can be an Error or a Warning based on the ExitCode
type CustomError interface {
	error
	GetExitCode() int
	GetMessage() string
	GetStack() string
	AsWarning() CustomError
	AsError() CustomError
	WithIcon(icon string) CustomError
}

// Combine combines multiple CustomError into one CustomError
func Combine(errs ...error) CustomError {
	if len(errs) < 1 {
		return nil
	}

	takenErrors := make([]error, 0)
	for _, err := range errs {
		if err == nil {
			continue
		}
		takenErrors = append(takenErrors, err)
	}
	if len(takenErrors) < 1 {
		return nil
	}

	errorList := NewErrorList(takenErrors...)
	return newError(errorList, false)
}

// HasErrors checks if the given errors contain any errors
func HasErrors(errs ...error) bool {
	for _, err := range errs {
		if err == nil {
			continue
		}
		if IsWarning(err) {
			continue
		}
		if IsError(err) || err != nil {
			return true
		}
	}
	return false
}

// HasWarnings checks if the given errors contain any warnings
func HasWarnings(errs ...error) bool {
	for _, err := range errs {
		if err == nil {
			continue
		}
		if IsError(err) {
			continue
		}
		if IsWarning(err) {
			return true
		}
	}
	return false
}

// GetWarnings returns a list of warnings from the given errors
func GetWarnings(errs ...error) ErrorList {
	res := ErrorList{
		errors: make([]error, 0),
	}
	for _, err := range errs {
		if !IsWarning(err) {
			continue
		}
		var warning CustomError
		errors.As(err, &warning)
		res.errors = append(res.errors, warning)
	}
	return res
}

// GetErrors returns a list of Errors from the given errors
func GetErrors(errs ...error) ErrorList {
	res := ErrorList{
		errors: make([]error, 0),
	}
	for _, err := range errs {
		if !IsError(err) {
			continue
		}
		var e CustomError
		errors.As(err, &e)
		res.errors = append(res.errors, e)
	}
	return res
}

// Wrap wraps the given error with additional information
func Wrap(err error, additional ...string) CustomError {
	if err == nil && len(additional) == 0 {
		return nil
	}

	if err == nil && len(additional) > 0 {
		err = errors.New(additional[0])
		additional = make([]string, 0)
	}
	if e, ok := errors.AsType[*Error](err); ok {
		for _, addi := range additional {
			e.Message += fmt.Sprintf("\n   > %s", addi)
		}
		return e
	}
	return newError(err, true, additional...)
}

// IsError checks if the given error is an Error
func IsError(err error) bool {
	var e CustomError
	if !errors.As(err, &e) {
		return false
	}
	return e.GetExitCode() == -1
}

// IsWarning checks if the given error is an Warning
func IsWarning(err error) bool {
	var e CustomError
	if !errors.As(err, &e) {
		return false
	}
	return e.GetExitCode() == -2
}

// IsCustomError checks if the given error is an CustomError
func IsCustomError(err error) bool {
	var e CustomError
	return errors.As(err, &e)
}

// AsCustomError casts the given error to an CustomError
func AsCustomError(err error) (CustomError, CustomError) {
	var e CustomError
	if !errors.As(err, &e) {
		return e, Wrap(nil, "error is not an CustomError")
	}
	return e, nil
}

// Equals checks if the given errors are equal
func Equals(a, b error) bool {
	if a == nil || b == nil {
		return a == b
	}

	aErr, aIsCustom := errors.AsType[CustomError](a)
	bErr, bIsCustom := errors.AsType[CustomError](b)

	if aIsCustom && bIsCustom {
		return aErr.GetExitCode() == bErr.GetExitCode() &&
			aErr.GetMessage() == bErr.GetMessage()
	}

	if aIsCustom {
		return aErr.GetMessage() == b.Error()
	}

	if bIsCustom {
		return a.Error() == bErr.GetMessage()
	}

	return a.Error() == b.Error()
}

type JsonErrorResponse struct {
	ExitCode int    `json:"exitCode"`
	Message  string `json:"message"`
}

// FormatJsonError returns a JSON formatted string representation of the given error
func FormatJsonError(err error) string {
	resp := JsonErrorResponse{
		ExitCode: -1,
		Message:  err.Error(),
	}

	if e, ok := errors.AsType[CustomError](err); ok {
		label := errorLabel
		if e.GetExitCode() == -2 {
			label = warningLabel
		}
		resp.ExitCode = e.GetExitCode()
		resp.Message = fmt.Sprintf("[%s]: %s\n%s", label, e.GetMessage(), e.GetStack())
	}

	jsonBytes, marshalErr := json.Marshal(resp)
	if marshalErr != nil {
		return `{"exitCode":-1,"message":"Failed to marshal error response"}`
	}

	return string(jsonBytes)
}
