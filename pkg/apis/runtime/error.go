package runtime

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/utils/json"
)

// Errorc wraps an HTTP status code as error.
func Errorc(c int) error {
	return asGinErr(c, nil, gin.ErrorTypePublic)
}

// Error wraps an HTTP status code and string message as error,
// and only outputs the string message to frontend.
func Error(c int, msg string) error {
	if msg == "" {
		return Errorc(c)
	}
	return asGinErr(c, errors.New(msg), gin.ErrorTypePublic)
}

// ErrorP wraps an HTTP status code and string message as error,
// but only logs the string message to backend.
func ErrorP(c int, msg string) error {
	if msg == "" {
		return Errorc(c)
	}
	return asGinErr(c, errors.New(msg), gin.ErrorTypePrivate)
}

// Errorf wraps an HTTP status code, format, and arguments as error,
// and only outputs the formatted message to frontend.
func Errorf(c int, format string, a ...any) error {
	if len(a) == 0 {
		return asGinErr(c, errors.New(format), gin.ErrorTypePublic)
	}
	return asGinErr(c, fmt.Errorf(format, a...), gin.ErrorTypePublic)
}

// ErrorfP wraps an HTTP status code, format, and arguments as error,
// but only logs the formatted message to backend.
func ErrorfP(c int, format string, a ...any) error {
	if len(a) == 0 {
		return asGinErr(c, errors.New(format), gin.ErrorTypePrivate)
	}
	return asGinErr(c, fmt.Errorf(format, a...), gin.ErrorTypePrivate)
}

// Errorw wraps an HTTP status code and the given error as error,
// and only outputs the formatted message to frontend.
func Errorw(c int, err error) error {
	return asGinErr(c, err, gin.ErrorTypePublic)
}

// ErrorwP wraps an HTTP status code and the given error as error,
// but only logs the formatted message to backend.
func ErrorwP(c int, err error) error {
	return asGinErr(c, err, gin.ErrorTypePrivate)
}

func asGinErr(c int, err error, typ gin.ErrorType) error {
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

func isGinError(err error) bool {
	if err == nil {
		return false
	}
	var ge *gin.Error
	return errors.As(err, &ge)
}

type httpError struct {
	code  int
	cause error
}

func (e httpError) Error() string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(e.code))
	sb.WriteString(" ")
	sb.WriteString(http.StatusText(e.code))
	if e.cause != nil {
		var ev = reflect.ValueOf(e.cause)
		switch ev.Kind() {
		case reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
			if ev.IsNil() {
				return sb.String()
			}
		}
		sb.WriteString(": ")
		sb.WriteString(e.cause.Error())
	}
	return sb.String()
}

func (e httpError) JSON() any {
	var jsonData = gin.H{}
	jsonData["status"] = e.code
	jsonData["statusText"] = http.StatusText(e.code)
	if e.cause != nil {
		var ev = reflect.ValueOf(e.cause)
		switch ev.Kind() {
		case reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
			if ev.IsNil() {
				return jsonData
			}
		}
		jsonData["message"] = e.cause.Error()
	}
	return jsonData
}

func (e httpError) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.JSON())
}
