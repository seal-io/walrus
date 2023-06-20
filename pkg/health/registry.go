package health

import (
	"context"
	"errors"
	"sync"
)

type (
	// Checker defines the operations of a health checker.
	Checker interface {
		Name() string
		Check(context.Context) error
	}

	// Checkers holds the list of Checker.
	Checkers []Checker
)

var (
	checkers Checkers
	o        sync.Once
)

// Register registers all health checkers.
func Register(ctx context.Context, cs Checkers) (err error) {
	err = errors.New("not allowed duplicated registering")

	o.Do(func() {
		checkers = cs
		err = nil
	})

	return
}

// Check defines the stereotype for health checking.
type Check func(context.Context) error

// CheckerFunc wraps the given Check as a Checker.
func CheckerFunc(name string, fn Check) Checker {
	return checker{n: name, f: fn}
}

type checker struct {
	n string
	f Check
}

func (c checker) Name() string {
	return c.n
}

func (c checker) Check(ctx context.Context) error {
	if c.f == nil {
		return nil
	}

	return c.f(ctx)
}
