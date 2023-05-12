package bus

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

// Bus defines the interface to decouple the logic.
type Bus interface {
	// Subscribe subscribes a Handler to process.
	Subscribe(string, Handler) error
	// Publish publishes the Message and callback all subscribed Handler.
	Publish(context.Context, Message) error
}

// Message defines the message to be handled.
type Message interface{}

// Handler defines the handler to process any publishing message,
// for example, func(context.Context, T) error.
// The incoming context.Context comes from the Publish function,
// if the inside processing is background,
// do not rely on that context.Context.
type Handler interface{}

// bus implements the Bus interface.
type bus map[string][]namedHandler

// handler wraps the Handler with name.
type namedHandler struct {
	n string
	h Handler
}

var globalBus = New()

// Subscribe implements the Bus interface.
func (b bus) Subscribe(n string, h Handler) error {
	if b == nil {
		return errors.New("nil bus")
	}

	ht := reflect.TypeOf(h)
	if ht.NumIn() != 2 {
		return errors.New("handler must has two parameters")
	}

	if ht.In(0).String() != "context.Context" {
		return errors.New("handler must uses 'context.Context' as the first parameter")
	}

	if ht.NumOut() != 1 {
		return errors.New("handler must has one result")
	}

	if ht.Out(0).String() != "error" {
		return errors.New("handler must uses 'error' as the only result")
	}

	switch ht.In(1).Kind() {
	case reflect.Chan, reflect.Interface, reflect.UnsafePointer, reflect.Uintptr, reflect.Func, reflect.Invalid:
		return errors.New("the second parameter of handler is invalid")
	default:
	}

	mt := getTypeSymbol(ht.In(1))
	b[mt] = append(b[mt], namedHandler{n: n, h: h})

	return nil
}

// Publish implements the Bus interface.
func (b bus) Publish(ctx context.Context, m Message) error {
	if b == nil {
		return errors.New("nil bus")
	}

	mt := getTypeSymbol(reflect.TypeOf(m))

	hs, exist := b[mt]
	if !exist {
		return nil
	}

	in := []reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(m),
	}
	for i := range hs {
		r := reflect.ValueOf(hs[i].h).Call(in)

		err := r[0].Interface()
		if err != nil {
			return fmt.Errorf("error calling %q handler: %v", hs[i].n, err)
		}
	}

	return nil
}

func getTypeSymbol(t reflect.Type) string {
	var (
		p string
		s = t.String()
	)

	switch t.Kind() {
	default:
		p = t.PkgPath()
	case reflect.Interface, reflect.Pointer:
		p = t.Elem().PkgPath()
	}

	if p != "" {
		s = p + "/" + s
	}

	return s
}

// New returns a new Bus.
func New() Bus {
	return bus{}
}

// Subscribe subscribes the handler to global Bus.
func Subscribe(n string, h Handler) error {
	return globalBus.Subscribe(n, h)
}

// MustSubscribe likes Subscribe, but panic if error found.
func MustSubscribe(n string, h Handler) {
	err := Subscribe(n, h)
	if err != nil {
		panic(err)
	}
}

// Publish publishes the message to all subscribed Handler.
func Publish(ctx context.Context, m Message) error {
	return globalBus.Publish(ctx, m)
}

// MustPublish likes Publish, but panic if error found.
func MustPublish(ctx context.Context, m Message) {
	err := Publish(ctx, m)
	if err != nil {
		panic(err)
	}
}
