package datamessage

import (
	"github.com/seal-io/seal/pkg/dao/model"
)

type EventType uint8

const (
	EventCreate = iota + 1
	EventUpdate
	EventDelete
)

func EventTypeFor(op model.Op) EventType {
	switch {
	case op.Is(model.OpCreate):
		return EventCreate
	case op.Is(model.OpUpdate) || op.Is(model.OpUpdateOne):
		return EventUpdate
	case op.Is(model.OpDelete) || op.Is(model.OpDeleteOne):
		return EventDelete
	}

	panic("unknown event type")
}
