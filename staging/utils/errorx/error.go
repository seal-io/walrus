package errorx

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

// Public returns the public error message and status code.
func Public(errs []error) (int, string) {
	var (
		status int
		pes    []PublicError
	)

	for i := range errs {
		unwraps := UnwrapErrors(errs[i])

		var pe PublicError

		for ui := range unwraps {
			e := unwraps[ui]

			// Get last public error, it should include all errors under itself.
			if pe == nil && errors.As(e, &pe) {
				pes = append(pes, pe)
			}

			// Get status while it isn't set.
			var he HttpError
			if status == 0 && errors.As(e, &he) {
				status = he.Status
			}
		}
	}

	return status, PublicFormat(pes)
}

// UnwrapErrors returns all errors in the chain.
func UnwrapErrors(err error) []error {
	errs := make([]error, 0)

	for err != nil {
		errs = append(errs, err)
		err = errors.Unwrap(err)
	}

	return errs
}

// Format returns the formatted error message.
func Format(errs []error) string {
	if len(errs) == 0 {
		return ""
	}

	if len(errs) == 1 {
		return fmt.Sprintf("1 error occurred:\n\t* %s\n\n", errs[0].Error())
	}

	msg := make([]string, len(errs))
	for i, err := range errs {
		msg[i] = fmt.Sprintf("* %s", err.Error())
	}

	return fmt.Sprintf(
		"%d errors occurred:\n\t%s\n\n",
		len(msg), strings.Join(msg, "\n\t"))
}

// PublicFormat returns the formatted public error message.
func PublicFormat(errs []PublicError) string {
	if len(errs) == 0 {
		return ""
	}

	if len(errs) == 1 {
		return errs[0].Public()
	}

	msg := make([]string, len(errs))
	for i, err := range errs {
		msg[i] = fmt.Sprintf("* %s", err.Public())
	}

	return fmt.Sprintf(
		"%d errors occurred: %s",
		len(msg), strings.Join(msg, ";"))
}

// PublicError is the public error message interface.
type PublicError interface {
	Public() string
}

// New returns an error with the supplied message.
func New(message string) ErrorX {
	return ErrorX{
		Message: message,
	}
}

// Errorf returns error with the supplied format and args.
func Errorf(format string, args ...any) ErrorX {
	return ErrorX{
		Message: fmt.Sprintf(format, args...),
	}
}

// Wrap returns an error with the supplied message.
func Wrap(err error, message string) ErrorX {
	return ErrorX{
		Cause:   err,
		Message: message,
	}
}

// Wrapf returns error with the supplied format and args.
func Wrapf(err error, format string, args ...any) ErrorX {
	return ErrorX{
		Cause:   err,
		Message: fmt.Sprintf(format, args...),
	}
}

// ErrorX is an implementation of PublicError interface.
//
//nolint:errname
type ErrorX struct {
	Cause   error
	Message string
}

// Error returns the error message.
func (e ErrorX) Error() string {
	var sb strings.Builder

	if e.Message != "" {
		sb.WriteString(e.Message)

		if e.Cause != nil {
			sb.WriteString(": ")
		}
	}

	if e.Cause != nil {
		sb.WriteString(e.Cause.Error())
	}

	return sb.String()
}

// Unwrap returns the cause error.
func (e ErrorX) Unwrap() error {
	return e.Cause
}

// Public returns the public error message.
func (e ErrorX) Public() string {
	return e.Error()
}

// NewHttpError returns an error with the http status and supplied message.
func NewHttpError(status int, message string) HttpError {
	return HttpError{
		ErrorX: ErrorX{
			Message: message,
		},
		Status: status,
	}
}

// HttpErrorf returns error with the http status, format and args.
func HttpErrorf(status int, format string, args ...any) HttpError {
	return HttpError{
		ErrorX: ErrorX{
			Message: fmt.Sprintf(format, args...),
		},
		Status: status,
	}
}

// WrapHttpError returns an error with the http status and supplied message.
func WrapHttpError(status int, err error, message string) HttpError {
	return HttpError{
		ErrorX: ErrorX{
			Cause:   err,
			Message: message,
		},
		Status: status,
	}
}

// WrapfHttpError returns error with the http status, format and args.
func WrapfHttpError(status int, err error, format string, args ...any) HttpError {
	return HttpError{
		ErrorX: ErrorX{
			Cause:   err,
			Message: fmt.Sprintf(format, args...),
		},
		Status: status,
	}
}

// HttpError is an implementation of PublicError interface.
type HttpError struct {
	Status int
	ErrorX
}

// Error returns the error message.
func (e HttpError) Error() string {
	var sb strings.Builder

	sb.WriteString(strconv.Itoa(e.Status))
	sb.WriteString(" ")
	sb.WriteString(http.StatusText(e.Status))

	if e.Cause != nil {
		ev := reflect.ValueOf(e.Cause)
		switch ev.Kind() {
		case reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
			if ev.IsNil() {
				return sb.String()
			}
		}

		sb.WriteString(": ")
		sb.WriteString(e.Cause.Error())
	} else if e.Message != "" {
		sb.WriteString(": ")
		sb.WriteString(e.Message)
	}

	return sb.String()
}

// Unwrap returns the cause error.
func (e HttpError) Unwrap() error {
	return e.Cause
}

// Public returns the public error message.
func (e HttpError) Public() string {
	return e.Error()
}
