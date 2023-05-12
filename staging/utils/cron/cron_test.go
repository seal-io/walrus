package cron

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testTask struct {
	sync.RWMutex

	outputs []interface{}
}

func (in *testTask) Name() string {
	return "test-task"
}

func (in *testTask) Process(ctx context.Context, args ...interface{}) error {
	in.Lock()
	defer in.Unlock()

	in.outputs = args
	if len(args) == 0 {
		in.outputs = []interface{}{"testing"}
	}

	return nil
}

func (in *testTask) Outputs() []interface{} {
	in.RLock()
	defer in.RUnlock()

	return in.outputs
}

func TestScheduler_Schedule(t *testing.T) {
	var err error
	err = Start(context.Background())
	assert.Nil(t, err, "error starting")

	var actual testTask
	err = Schedule("test", AwaitedExpr("0/1 * * * * * *"), &actual)
	assert.Equal(
		t,
		"invalid cron expression: expected exactly 6 fields, found 7: [0/1 * * * * * *]",
		err.Error(),
		"error scheduling",
	)

	actual = testTask{}
	err = Schedule("test", ImmediateExpr("* * * ? * *"), &actual)
	assert.Nil(t, err, "error none variables scheduling")
	time.Sleep(3 * time.Second) // Give an enough range to execute scheduling.
	assert.Equal(t, []interface{}{"testing"}, actual.Outputs(),
		"invalid output of none variables scheduling")

	actual = testTask{}
	err = Schedule("test", AwaitedExpr("* * * ? * *"), &actual, "test", "with", "variables")
	assert.Nil(t, err, "error variables scheduling")
	time.Sleep(5 * time.Second) // Give an enough range to execute scheduling.
	assert.Equal(t, []interface{}{"test", "with", "variables"}, actual.Outputs(),
		"invalid output of variables scheduling")

	err = Stop()
	assert.Nil(t, err, "error stopping")
}
