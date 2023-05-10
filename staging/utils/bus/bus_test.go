package bus

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testMessage struct {
	namespace string
	name      string
}

func (in testMessage) String() string {
	return in.namespace + "/" + in.name
}

func TestBus_PublishStruct(t *testing.T) {
	var actual testMessage
	err := Subscribe("struct", func(ctx context.Context, m testMessage) error {
		actual = m
		return nil
	})
	assert.Nil(t, err, "error subscribing")

	err = Publish(context.Background(), testMessage{namespace: "ns-abc", name: "n-efg"})
	assert.Nil(t, err, "error publishing")

	assert.Equal(t, actual, testMessage{namespace: "ns-abc", name: "n-efg"}, "unexpected")
}

func TestBus_PublishPointer(t *testing.T) {
	var actual *testMessage
	err := Subscribe("pointer", func(ctx context.Context, m *testMessage) error {
		actual = m
		return nil
	})
	assert.Nil(t, err, "error subscribing")

	err = Publish(context.Background(), &testMessage{namespace: "ns-abc", name: "n-efg"})
	assert.Nil(t, err, "error publishing")

	assert.Equal(t, actual, &testMessage{namespace: "ns-abc", name: "n-efg"}, "unexpected")
}

func TestBus_PublishInterfaceValue(t *testing.T) {
	var actual *testMessage
	err := Subscribe("interface value", func(ctx context.Context, m *testMessage) error {
		actual = m
		return nil
	})
	assert.Nil(t, err, "error subscribing")

	var v interface{} = &testMessage{namespace: "ns-abc", name: "n-efg"}
	err = Publish(context.Background(), v)
	assert.Nil(t, err, "error publishing")

	assert.Equal(t, actual, &testMessage{namespace: "ns-abc", name: "n-efg"}, "unexpected")
}

func TestBus_PublishMap(t *testing.T) {
	var actual map[string]int
	err := Subscribe("map", func(ctx context.Context, m map[string]int) error {
		actual = m
		return nil
	})
	assert.Nil(t, err, "error subscribing")

	err = Publish(context.Background(), map[string]int{"a": 0, "b": 1})
	assert.Nil(t, err, "error publishing")

	assert.Equal(t, actual, map[string]int{"a": 0, "b": 1}, "unexpected")
}

func TestBus_PublishSlice(t *testing.T) {
	var actual []string
	err := Subscribe("slice", func(ctx context.Context, m []string) error {
		actual = m
		return nil
	})
	assert.Nil(t, err, "error subscribing")

	err = Publish(context.Background(), []string{"a", "b"})
	assert.Nil(t, err, "error publishing")

	assert.Equal(t, actual, []string{"a", "b"}, "unexpected")
}

func TestBus_PublishArray(t *testing.T) {
	var actual [2]int
	err := Subscribe("array", func(ctx context.Context, m [2]int) error {
		actual = m
		return nil
	})
	assert.Nil(t, err, "error subscribing")

	err = Publish(context.Background(), [2]int{0, 1})
	assert.Nil(t, err, "error publishing")

	assert.Equal(t, actual, [2]int{0, 1}, "unexpected")
}
