package fail_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/codeus-node/fail"
	"github.com/codeus-node/icons"
)

func TestWrap_Nil(t *testing.T) {
	err := fail.Wrap(nil)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestWrap_CreateCustomError(t *testing.T) {
	err := fail.Wrap(nil, "this is my new special custom error")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	expectedMessage := icons.Error + "[AP_ERROR]: this is my new special custom error"
	if !strings.HasPrefix(err.Error(), expectedMessage) {
		t.Fatalf("expected error message, %s == \n\n%s", expectedMessage, err.Error())
	}
	if !strings.Contains(err.Error(), "api_test.go") {
		t.Fatalf("expected stack trace on error message but get %s", err.Error())
	}
	println(err.Error())
}

func TestWrap_WrapAnExistingError(t *testing.T) {
	err := fail.Wrap(fmt.Errorf("this is my new special custom error"), "this is my additional error")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	expectedMessage := icons.Error + "[AP_ERROR]: this is my new special custom error\n   > this is my additional error"
	if !strings.HasPrefix(err.Error(), expectedMessage) {
		t.Fatalf("expected error message, %s == \n\n%s", expectedMessage, err.Error())
	}
	println(err.Error())
}

func TestWrap_ChainedErrors(t *testing.T) {
	err := a()
	if err == nil {
		t.Fatalf("expected error, got nil")
		return
	}
	expectedMessage := icons.Error + "[AP_ERROR]: this is my external error\n   > called in c()\n   > called in b()\n   > called in a()"
	if !strings.HasPrefix(err.Error(), expectedMessage) {
		t.Fatalf("expected error message, %s == \n\n%s", expectedMessage, err.Error())
	}
	println(err.Error())
}

func TestWrap_Warning(t *testing.T) {
	err := fail.Wrap(nil, "this is my warning").AsWarning()
	if err == nil {
		t.Fatalf("expected warning, got nil")
	}
	expectedMessage := icons.Warning + "[AP_WARNING]: this is my warning"
	if !strings.HasPrefix(err.Error(), expectedMessage) {
		t.Fatalf("expected warning message, %s == \n\n%s", expectedMessage, err.Error())
	}
	println(err.Error())
}

func TestWrap_IsWarning_And_IsError(t *testing.T) {
	err := fail.Wrap(nil, "this is my warning").AsWarning()
	if !fail.IsWarning(err) {
		t.Fatalf("expected warning, got %s", err.Error())
	}
	err = fail.Wrap(nil, "this is my error")
	if fail.IsWarning(err) {
		t.Fatalf("expected error, got %s", err.Error())
	}
	if !fail.IsError(err) {
		t.Fatalf("expected error, got %s", err.Error())
	}
	if fail.IsError(nil) {
		t.Fatalf("expected nil was no Error")
	}
	if fail.IsWarning(nil) {
		t.Fatalf("expected nil was no Warning")
	}
}

func TestFormatJsonError(t *testing.T) {
	err := fail.Wrap(nil, "this is my error")
	expected := "{\"exitCode\":-1,\"message\":\"[AP_ERROR]: this is my error"
	tmp := fail.FormatJsonError(err)
	if !strings.HasPrefix(tmp, expected) {
		t.Fatalf("expected prefix %s, got %s", expected, fail.FormatJsonError(err))
	}
	println(tmp)

	err = fail.Wrap(nil, "this is my warning").AsWarning()
	expected = "{\"exitCode\":-2,\"message\":\"[AP_WARNING]: this is my warning"
	tmp = fail.FormatJsonError(err)
	if !strings.HasPrefix(tmp, expected) {
		t.Fatalf("expected prefix %s, got %s", expected, fail.FormatJsonError(err))
	}
	println(tmp)
}

func TestCombine(t *testing.T) {
	err := fail.Combine(
		fmt.Errorf("this is error 1"),
		fmt.Errorf("this is error 2"),
		fmt.Errorf("this is error 3"))
	if err == nil {
		t.Fatalf("expected combined error, got %v", err)
	}
	println(err.Error())
}

func TestCombine_Nil(t *testing.T) {
	err := fail.Combine(
		nil,
		fmt.Errorf("this is error 1"),
		fmt.Errorf("this is error 2"),
		fmt.Errorf("this is error 3"))
	if err == nil {
		t.Fatalf("expected combined error, got %v", err)
	}
	println(err.Error())
}

func TestCombine_CustomError(t *testing.T) {
	err := fail.Combine(
		fmt.Errorf("this is error 1"),
		fail.Wrap(fmt.Errorf("this is error 2")),
		fail.Wrap(fmt.Errorf("this is error 3")).AsWarning(),
	)
	if err == nil {
		t.Fatalf("expected combined error, got %v", err)
	}
	println(err.Error())
}

func TestHasErrors_SliceOfStandardErrors(t *testing.T) {
	errs := []error{
		fmt.Errorf("this is error 1"),
		fmt.Errorf("this is error 2"),
		fmt.Errorf("this is error 3"),
	}
	hasErrors := fail.HasErrors(errs...)
	if !hasErrors {
		t.Fatalf("expected hasErrors to be true, got %v", hasErrors)
	}
}

func TestHasWarnings_SliceOfStandardErrors(t *testing.T) {
	errs := []error{
		fmt.Errorf("this is error 1"),
		fmt.Errorf("this is error 2"),
		fmt.Errorf("this is error 3"),
	}
	hasWarnings := fail.HasWarnings(errs...)
	if hasWarnings {
		t.Fatalf("expected hasWarnings to be false, got %v", hasWarnings)
	}
}

func TestHasErrors_SliceOfCustomErrors(t *testing.T) {
	errs := []error{
		fmt.Errorf("this is error 1"),
		fail.Wrap(fmt.Errorf("this is error 2")),
		fail.Wrap(fmt.Errorf("this is error 3")).AsWarning(),
	}
	hasErrors := fail.HasErrors(errs...)
	if !hasErrors {
		t.Fatalf("expected hasErrors to be true, got %v", hasErrors)
	}
}

func TestHasWarnings_SliceOfCustomErrors(t *testing.T) {
	errs := []error{
		fmt.Errorf("this is error 1"),
		fail.Wrap(fmt.Errorf("this is error 2")),
		fail.Wrap(fmt.Errorf("this is error 3")).AsWarning(),
	}
	hasWarnings := fail.HasWarnings(errs...)
	if !hasWarnings {
		t.Fatalf("expected hasWarnings to be true, got %v", hasWarnings)
	}
}

func TestHasErrors_Nil(t *testing.T) {
	hasErrors := fail.HasErrors()
	if hasErrors {
		t.Fatalf("expected hasErrors to be false, got %v", hasErrors)
	}
	errs := []error{nil}
	hasErrors = fail.HasErrors(errs...)
	if hasErrors {
		t.Fatalf("expected hasErrors to be false, got %v", hasErrors)
	}
}

func TestHasWarnings_Nil(t *testing.T) {
	hasWarnings := fail.HasWarnings()
	if hasWarnings {
		t.Fatalf("expected hasWarnings to be false, got %v", hasWarnings)
	}
	errs := []error{nil}
	hasWarnings = fail.HasWarnings(errs...)
	if hasWarnings {
		t.Fatalf("expected hasWarnings to be false, got %v", hasWarnings)
	}
}

func TestHasErrors_OnlyWarnings(t *testing.T) {
	errs := []error{
		fail.Wrap(fmt.Errorf("this is Warning 1")).AsWarning(),
		fail.Wrap(fmt.Errorf("this is Warning 2")).AsWarning(),
		fail.Wrap(fmt.Errorf("this is Warning 3")).AsWarning(),
	}
	hasErrors := fail.HasErrors(errs...)
	if hasErrors {
		t.Fatalf("expected hasErrors to be false, got %v", hasErrors)
	}
}

func TestHasWarnings_OnlyWarnings(t *testing.T) {
	errs := []error{
		fail.Wrap(fmt.Errorf("this is Warning 1")).AsWarning(),
		fail.Wrap(fmt.Errorf("this is Warning 2")).AsWarning(),
		fail.Wrap(fmt.Errorf("this is Warning 3")).AsWarning(),
	}
	hasWarnings := fail.HasWarnings(errs...)
	if !hasWarnings {
		t.Fatalf("expected hasWarnings to be true, got %v", hasWarnings)
	}
}

func TestGetErrors_OnlyErrors(t *testing.T) {
	errs := []error{
		fail.Wrap(fmt.Errorf("this is Error 1")),
		fail.Wrap(fmt.Errorf("this is Error 2")),
		fail.Wrap(fmt.Errorf("this is Error 3")),
	}
	apErrors := fail.GetErrors(errs...)
	if apErrors.Count() != 3 {
		t.Fatalf("expected 3 errors, got %v", apErrors.Count())
	}
}

func TestGetErrors_Nil(t *testing.T) {
	apErrors := fail.GetErrors()
	if apErrors.Count() != 0 {
		t.Fatalf("expected 0 errors, got %v", apErrors.Count())
	}
}

func TestGetErrors_MultipleErrors(t *testing.T) {
	errs := []error{
		fmt.Errorf("this is Error 1"),
		fail.Wrap(fmt.Errorf("this is Error 2")),
		fail.Wrap(fmt.Errorf("this is Warning 1")).AsWarning(),
	}
	apErrors := fail.GetErrors(errs...)
	if apErrors.Count() != 1 {
		t.Fatalf("expected 1 error, got %v", apErrors.Count())
	}
}

func TestGetWarnings_OnlyWarnings(t *testing.T) {
	errs := []error{
		fail.Wrap(fmt.Errorf("this is Warning 1")).AsWarning(),
		fail.Wrap(fmt.Errorf("this is Warning 2")).AsWarning(),
		fail.Wrap(fmt.Errorf("this is Warning 3")).AsWarning(),
	}
	warnings := fail.GetWarnings(errs...)
	if warnings.Count() != 3 {
		t.Fatalf("expected 3 warnings, got %v", warnings.Count())
	}
}

func TestGetWarnings_Nil(t *testing.T) {
	warnings := fail.GetWarnings()
	if warnings.Count() != 0 {
		t.Fatalf("expected 0 warnings, got %v", warnings.Count())
	}
}

func TestGetWarnings_MultipleErrors(t *testing.T) {
	errs := []error{
		fmt.Errorf("this is Warning 1"),
		fail.Wrap(fmt.Errorf("this is Warning 2")),
		fail.Wrap(fmt.Errorf("this is Warning 3")).AsWarning(),
	}
	warnings := fail.GetWarnings(errs...)
	if warnings.Count() != 1 {
		t.Fatalf("expected 1 warnings, got %v", warnings.Count())
	}
}

func Test_HasDefaultErrorIcon(t *testing.T) {
	err := fail.Wrap(nil, "error with icon")
	if !strings.HasPrefix(err.Error(), icons.Error) {
		t.Fatalf("expect icon %s at the begin of error message: %s", icons.Error, err.Error())
	}
	println(err.Error())
}

func Test_HasDefaultWarningIcon(t *testing.T) {
	err := fail.Wrap(nil, "error with icon").AsWarning()
	if !strings.HasPrefix(err.Error(), icons.Warning) {
		t.Fatalf("expect icon %s at the begin of error message: %s", icons.Warning, err.Error())
	}
	println(err.Error())
}

func Test_PrintWithCustomIcon(t *testing.T) {
	err := fail.Wrap(nil, "error with icon").WithIcon(icons.Knowledge)
	if !strings.HasPrefix(err.Error(), icons.Knowledge) {
		t.Fatalf("expect icon %s at the begin of error message: %s", icons.Knowledge, err.Error())
	}
	println(err.Error())
}

func Test_IsCustomError(t *testing.T) {
	err := fail.Wrap(nil, "error with icon").WithIcon(icons.Knowledge)
	if !fail.IsCustomError(err) {
		t.Fatalf("expect custom Error: %s", err.Error())
	}
	normalErr := fmt.Errorf("normal error")
	if fail.IsCustomError(normalErr) {
		t.Fatalf("expect no custom Error: %s", err.Error())
	}
}

func Test_AsCustomError(t *testing.T) {
	var custom fail.CustomError
	var convertedCustom fail.CustomError
	var convertErr fail.CustomError

	custom = fail.Wrap(nil, "error with icon").WithIcon(icons.Knowledge)
	convertedCustom, convertErr = fail.AsCustomError(custom)
	if convertErr != nil {
		t.Fatalf("expect no converter Error: %s", convertErr.Error())
	}
	if convertedCustom.GetExitCode() != -1 {
		t.Fatalf("expect custom Error: %s", convertedCustom.Error())
	}

	_, convertErr = fail.AsCustomError(fmt.Errorf("normal error"))
	if convertErr == nil {
		t.Fatalf("expect converter Error: nil")
	}
}

func Test_Equals(t *testing.T) {
	tests := []struct {
		name string
		a    error
		b    error
		want bool
	}{
		{
			name: "nil equals nil",
			a:    nil,
			b:    nil,
			want: true,
		},
		{
			name: "nil does not equal custom error",
			a:    nil,
			b:    fail.Wrap(nil, "error with icon"),
			want: false,
		},
		{
			name: "custom error does not equal nil",
			a:    fail.Wrap(nil, "error with icon"),
			b:    nil,
			want: false,
		},
		{
			name: "custom errors with same message and exit code are equal",
			a:    fail.Wrap(nil, "error with icon").WithIcon(icons.Knowledge),
			b:    fail.Wrap(nil, "error with icon").WithIcon(icons.Knowledge),
			want: true,
		},
		{
			name: "custom errors with different messages are not equal",
			a:    fail.Wrap(nil, "error1 with icon").WithIcon(icons.Knowledge),
			b:    fail.Wrap(nil, "error2 with icon").WithIcon(icons.Knowledge),
			want: false,
		},
		{
			name: "custom errors with same message but different exit code are not equal",
			a:    fail.Wrap(nil, "error with icon").AsError(),
			b:    fail.Wrap(nil, "error with icon").AsWarning(),
			want: false,
		},
		{
			name: "normal error equals custom error with same message",
			a:    errors.New("error with icon"),
			b:    fail.Wrap(nil, "error with icon"),
			want: true,
		},
		{
			name: "custom error equals normal error with same message",
			a:    fail.Wrap(nil, "error with icon"),
			b:    errors.New("error with icon"),
			want: true,
		},
		{
			name: "normal error does not equal custom error with different message",
			a:    errors.New("normal error"),
			b:    fail.Wrap(nil, "custom error"),
			want: false,
		},
		{
			name: "custom error does not equal normal error with different message",
			a:    fail.Wrap(nil, "custom error"),
			b:    errors.New("normal error"),
			want: false,
		},
		{
			name: "normal errors with same message are equal",
			a:    errors.New("1"),
			b:    errors.New("1"),
			want: true,
		},
		{
			name: "normal errors with different messages are not equal",
			a:    errors.New("1"),
			b:    errors.New("2"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fail.Equals(tt.a, tt.b)
			if got != tt.want {
				t.Fatalf("Equals(%v, %v) = %v; want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func a() error {
	err := b()
	return fail.Wrap(err, "called in a()")
}

func b() error {
	err := c()
	return fail.Wrap(err, "called in b()")
}

func c() error {
	return fail.Wrap(external(), "called in c()")
}

func external() error {
	return fmt.Errorf("this is my external error")
}
