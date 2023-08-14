package server

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func privateInitiationName(context.Context, initOptions) error {
	panic("test only")
}

func PublicInitiationName(context.Context, initOptions) error {
	panic("test only")
}

type _X struct{}

func (_X) StructInitiationName(context.Context, initOptions) error {
	panic("test only")
}

func Test_loadInitiationName(t *testing.T) {
	anonymityInitiationName := func(context.Context, initOptions) error {
		panic("test only")
	}

	testCases := []struct {
		given    initiation
		expected string
	}{
		{
			given:    privateInitiationName,
			expected: "private initiation name",
		},
		{
			given:    PublicInitiationName,
			expected: "public initiation name",
		},
		{
			given:    _X{}.StructInitiationName,
			expected: "struct initiation name",
		},
		{
			given:    anonymityInitiationName,
			expected: "func1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			actual := loadInitiationName(tc.given)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
