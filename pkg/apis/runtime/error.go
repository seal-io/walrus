package runtime

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/seal-io/walrus/utils/json"
)

// Errorc wraps an HTTP status code as error.
func Errorc(c int) error {
	return beGinErr(c, nil, gin.ErrorTypePublic)
}

// Error wraps an HTTP status code and string message or error as error.
func Error(c int, v any) error {
	if v != nil {
		switch t := v.(type) {
		case error:
			return beGinErr(c, t, gin.ErrorTypePublic)
		case string:
			if t != "" {
				return beGinErr(c, errors.New(t), gin.ErrorTypePublic)
			}
		}
	}

	return Errorc(c)
}

// Errorf wraps an HTTP status code, format, and arguments as error.
func Errorf(c int, format string, a ...any) error {
	if format == "" {
		return Errorc(c)
	}

	if len(a) == 0 {
		return beGinErr(c, errors.New(format), gin.ErrorTypePublic)
	}

	return beGinErr(c, fmt.Errorf(format, a...), gin.ErrorTypePublic)
}

// Errorp wraps an HTTP status code and string message as error,
// but logs the given error internally.
func Errorp(c int, err error, msg string) error {
	if msg == "" {
		return Errorc(c)
	}

	return beGinErr(c, wrapError{internal: err, external: errors.New(msg)}, gin.ErrorTypePrivate)
}

// Errorpf wraps an HTTP status code, format, and arguments as error,
// but logs the given error internally.
func Errorpf(c int, err error, format string, a ...any) error {
	if format == "" {
		return Errorc(c)
	}

	if len(a) == 0 {
		return Errorp(c, err, format)
	}

	return beGinErr(c, wrapError{internal: err, external: fmt.Errorf(format, a...)}, gin.ErrorTypePrivate)
}

// Errorw is similar to Errorp,
// it gains the HTTP status code from the given error,
// and wraps the HTTP status code and string message as error,
// but logs the given error internally.
func Errorw(err error, msg string) error {
	if err == nil {
		return nil
	}

	return &gin.Error{
		Err: wrapError{
			internal: err,
			external: errors.New(msg),
		},
		Type: gin.ErrorTypePrivate,
	}
}

// Errorwf is similar to Errorpf,
// it gains the HTTP status code from the given error,
// and wraps the HTTP status code and format, and arguments as error,
// but logs the given error internally.
func Errorwf(err error, format string, a ...any) error {
	if err == nil {
		return nil
	}

	if len(a) == 0 {
		return Errorw(err, format)
	}

	return &gin.Error{
		Err: wrapError{
			internal: err,
			external: fmt.Errorf(format, a...),
		},
		Type: gin.ErrorTypePrivate,
	}
}

func beGinErr(c int, err error, typ gin.ErrorType) error {
	if c == http.StatusOK {
		return nil
	}

	return &gin.Error{
		Err: httpError{
			code:  c,
			cause: err,
		},
		Type: typ,
	}
}

type httpError struct {
	code  int
	cause error
	brief string
}

func (e httpError) Error() string {
	var sb strings.Builder

	sb.WriteString(strconv.Itoa(e.code))
	sb.WriteString(" ")
	sb.WriteString(http.StatusText(e.code))

	if e.cause != nil {
		ev := reflect.ValueOf(e.cause)
		switch ev.Kind() {
		case reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
			if ev.IsNil() {
				return sb.String()
			}
		}

		sb.WriteString(": ")
		sb.WriteString(e.cause.Error())
	} else if e.brief != "" {
		sb.WriteString(": ")
		sb.WriteString(e.brief)
	}

	return sb.String()
}

func (e httpError) JSON() any {
	jsonData := gin.H{}
	jsonData["status"] = e.code
	jsonData["statusText"] = http.StatusText(e.code)

	if e.cause != nil {
		ev := reflect.ValueOf(e.cause)
		switch ev.Kind() {
		case reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
			if ev.IsNil() {
				return jsonData
			}
		}
		jsonData["message"] = e.cause.Error()
	} else if e.brief != "" {
		jsonData["message"] = e.brief
	}

	return jsonData
}

func (e httpError) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.JSON())
}

type wrapError struct {
	internal error
	external error
}

func (e wrapError) Error() string {
	return fmt.Sprintf("%v: %v", e.external, e.internal)
}

func (e wrapError) Unwrap() error {
	return e.internal
}
