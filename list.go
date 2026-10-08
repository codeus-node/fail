package fail

import "fmt"

// ErrorList is a list of errors
type ErrorList struct {
	error
	errors []error
}

// ==========================================
// Error Interface Implementation
// ==========================================

// Error returns a string representation of the error list
func (e *ErrorList) Error() string {
	result := ""
	for _, err := range e.errors {
		result = fmt.Sprintf("%s%s\n----------------\n", result, err.Error())
	}
	return result
}

// ==========================================
// Stringer Interface Implementation
// ==========================================

func (e *ErrorList) String() string {
	return e.Error()
}

// ToErrorList returns the underlying error list
func (e *ErrorList) ToErrorList() []error {
	return e.errors
}

// Count gets the number of Elements in the List
func (e *ErrorList) Count() int {
	return len(e.errors)
}

func NewErrorList(errors ...error) *ErrorList {
	takenErrors := make([]error, 0)
	for _, err := range errors {
		if err == nil {
			continue
		}
		takenErrors = append(takenErrors, err)
	}
	return &ErrorList{errors: takenErrors}
}
