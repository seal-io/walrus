package manifest

import (
	"errors"
	"fmt"
	"time"

	utilwait "k8s.io/apimachinery/pkg/util/wait"

	"github.com/seal-io/walrus/pkg/cli/common"
	"github.com/seal-io/walrus/pkg/cli/config"
	"github.com/seal-io/walrus/utils/strs"
)

// Operator is an interface for applying and deleting objects.
type Operator interface {
	Operate(set ObjectSet) (r OperateResult, err error)
	PrintResult(r OperateResult)
}

// OperateResult is a type that represents the result of an operation.
type OperateResult struct {
	Success   ObjectSet
	Failed    ObjectSet
	NotFound  ObjectSet
	UnChanged ObjectSet
}

// DefaultApplyOperator returns a apply Operator.
func DefaultApplyOperator(sc *config.Config, wait bool) Operator {
	return &ApplyOperator{
		operatorConfig: newOperatorConfig(sc, wait),
	}
}

// DefaultDeleteOperator returns a delete Operator.
func DefaultDeleteOperator(sc *config.Config, wait bool) Operator {
	return &DeleteOperator{
		operatorConfig: newOperatorConfig(sc, wait),
	}
}

// DefaultPreviewOperator returns preview Operator.
func DefaultPreviewOperator(sc *config.Config, wait bool) Operator {
	return &PreviewOperator{
		operatorConfig: newOperatorConfig(sc, wait),
		extraBodyParams: map[string]any{
			"preview": true,
		},
	}
}

type operatorConfig struct {
	serverContext *config.Config
	groupSequence []string
	backoff       *utilwait.Backoff
	wait          bool
}

func newOperatorConfig(sc *config.Config, wait bool) operatorConfig {
	return operatorConfig{
		serverContext: sc,
		groupSequence: GroupSequence,
		backoff: &utilwait.Backoff{
			Duration: 100 * time.Millisecond,
			Factor:   2,
			Steps:    3,
		},
		wait: wait,
	}
}

// ApplyOperator is a type that represents an apply operator.
type ApplyOperator struct {
	operatorConfig
}

// Operate applies the provided ObjectSet.
func (o *ApplyOperator) Operate(set ObjectSet) (OperateResult, error) {
	if o.backoff == nil {
		return o.apply(set)
	}

	var (
		err         error
		retryErr    *common.RetryableError
		finalResult = OperateResult{
			Success:   ObjectSet{},
			UnChanged: ObjectSet{},
		}
	)
	err = utilwait.ExponentialBackoff(*o.backoff, func() (bool, error) {
		if set.Len() == 0 {
			return true, nil
		}

		r, err := o.apply(set)
		if err != nil {
			if errors.As(err, &retryErr) {
				err = nil
			}
			set = r.Failed

			return false, err
		}

		finalResult.UnChanged.Add(r.UnChanged.All()...)
		finalResult.Success.Add(r.Success.All()...)

		set = ObjectSet{}

		return true, nil
	})

	finalResult.Failed = set

	return finalResult, err
}

func (o *ApplyOperator) apply(set ObjectSet) (r OperateResult, err error) {
	if set.Len() == 0 {
		return
	}

	r = OperateResult{
		Success:   ObjectSet{},
		Failed:    set,
		UnChanged: ObjectSet{},
	}

	for _, group := range o.groupSequence {
		objByScope := set.ByGroup(group)
		if len(objByScope) == 0 {
			continue
		}

		unchanged, notFound, changed, err := GetObjects(o.serverContext, group, objByScope, true)
		if err != nil {
			return r, err
		}

		// Unchanged.
		r.UnChanged.Add(unchanged.All()...)
		r.Failed.Remove(unchanged.All()...)

		// Patch.
		successPatched, _, err := PatchObjects(o.serverContext, group, changed.ByGroup(group), nil)
		if err != nil {
			return r, err
		}

		r.Success.Add(successPatched.All()...)
		r.Failed.Remove(successPatched.All()...)

		// Batch create.
		successCreated, _, err := BatchCreateObjects(o.serverContext, group, notFound.ByGroup(group), nil)
		if err != nil {
			return r, err
		}

		r.Success.Add(successCreated.All()...)
		r.Failed.Remove(successCreated.All()...)
	}

	return r, nil
}

// PrintResult prints the result of an operation.
func (o *ApplyOperator) PrintResult(r OperateResult) {
	var (
		failed    = "apply failed"
		notFound  = "not found"
		unchanged = "unchanged"

		success = map[ObjectStatus]string{
			statusNotFound: "created",
			statusChanged:  "patched",
		}

		wait = map[ObjectStatus]string{
			statusNotFound: "creating",
			statusChanged:  "patching",
		}
	)

	for _, os := range r.NotFound.All() {
		fmt.Printf("%s %s %s\n", strs.Singularize(os.Group), os.Key(), notFound)
	}

	for _, os := range r.UnChanged.All() {
		fmt.Printf("%s %s %s\n", strs.Singularize(os.Group), os.Key(), unchanged)
	}

	for _, os := range r.Success.All() {
		msg := success[os.Status]
		if o.wait {
			msg = wait[os.Status]
		}

		fmt.Printf("%s %s %s\n", strs.Singularize(os.Group), os.Key(), msg)
	}

	for _, os := range r.Failed.All() {
		fmt.Printf("%s %s %s\n", strs.Singularize(os.Group), os.Key(), failed)
	}
}

// DeleteOperator is a type that represents an delete operator.
type DeleteOperator struct {
	operatorConfig
}

// Operate deletes the provided ObjectSet.
func (o *DeleteOperator) Operate(set ObjectSet) (OperateResult, error) {
	if o.backoff == nil {
		return o.delete(set)
	}

	var (
		retryErr *common.RetryableError
		result   = OperateResult{
			Success:  ObjectSet{},
			Failed:   ObjectSet{},
			NotFound: ObjectSet{},
		}
	)

	err := utilwait.ExponentialBackoff(*o.backoff, func() (bool, error) {
		if set.Len() == 0 {
			return true, nil
		}

		r, err := o.delete(set)
		if err != nil {
			if errors.As(err, &retryErr) {
				err = nil
			}

			set = r.Failed
			result.NotFound.Add(r.NotFound.All()...)
			result.Success.Add(r.Success.All()...)

			return false, err
		}

		set = ObjectSet{}

		result.NotFound.Add(r.NotFound.All()...)
		result.Success.Add(r.Success.All()...)

		return true, nil
	})

	result.Failed.Add(set.All()...)

	return result, err
}

func (o *DeleteOperator) delete(set ObjectSet) (r OperateResult, err error) {
	if set.Len() == 0 {
		return
	}

	r = OperateResult{
		Success:  ObjectSet{},
		Failed:   set,
		NotFound: ObjectSet{},
	}
	// Delete in reverse order.
	for i := len(o.groupSequence) - 1; i >= 0; i-- {
		group := o.groupSequence[i]

		objByScope := set.ByGroup(group)
		if len(objByScope) == 0 {
			continue
		}

		unchanged, notFound, _, err := GetObjects(o.serverContext, group, objByScope, false)
		if err != nil {
			return r, err
		}

		r.NotFound.Add(notFound.All()...)

		if unchanged.Len() == 0 {
			continue
		}

		successDeleted, _, err := DeleteObjects(o.serverContext, group, unchanged.ByGroup(group))
		if err != nil {
			return r, err
		}

		r.Failed.Remove(successDeleted.All()...)
		r.Success.Add(successDeleted.All()...)
	}

	return r, nil
}

// PrintResult prints the result of an operation.
func (o *DeleteOperator) PrintResult(r OperateResult) {
	var (
		waiting  = "deleting"
		failed   = "delete failed"
		notFound = "not found"
	)

	for _, os := range r.Success.All() {
		success := "deleted"
		if o.wait {
			success = waiting
		}

		fmt.Printf("%s %s %s\n", strs.Singularize(os.Group), os.Key(), success)
	}

	for _, os := range r.Failed.All() {
		fmt.Printf("%s %s %s\n", strs.Singularize(os.Group), os.Key(), failed)
	}

	for _, os := range r.NotFound.All() {
		fmt.Printf("%s %s %s\n", strs.Singularize(os.Group), os.Key(), notFound)
	}
}

// PreviewOperator is a type that represents preview operator.
type PreviewOperator struct {
	operatorConfig
	extraBodyParams map[string]any
}

// Operate generate preview changes of the provided ObjectSet.
func (o *PreviewOperator) Operate(set ObjectSet) (OperateResult, error) {
	if o.backoff == nil {
		return o.preview(set)
	}

	var (
		err         error
		retryErr    *common.RetryableError
		finalResult = OperateResult{
			Success:   ObjectSet{},
			UnChanged: ObjectSet{},
		}
	)
	err = utilwait.ExponentialBackoff(*o.backoff, func() (bool, error) {
		if set.Len() == 0 {
			return true, nil
		}

		r, err := o.preview(set)
		if err != nil {
			if errors.As(err, &retryErr) {
				err = nil
			}
			set = r.Failed

			return false, err
		}

		finalResult.UnChanged.Add(r.UnChanged.All()...)
		finalResult.Success.Add(r.Success.All()...)

		set = ObjectSet{}

		return true, nil
	})

	finalResult.Failed = set

	return finalResult, err
}

func (o *PreviewOperator) PrintResult(r OperateResult) {
	var (
		failed = "preview generation failed"

		success = map[ObjectStatus]string{
			statusNotFound:  "preview generated",
			statusUnchanged: "preview generated",
		}

		wait = map[ObjectStatus]string{
			statusNotFound:  "preview generating",
			statusUnchanged: "preview generating",
		}
	)

	for _, os := range r.Success.All() {
		msg := success[os.Status]
		if o.wait {
			msg = wait[os.Status]
		}

		fmt.Printf("%s %s %s\n", strs.Singularize(os.Group), os.Key(), msg)
	}

	for _, os := range r.Failed.All() {
		fmt.Printf("%s %s %s\n", strs.Singularize(os.Group), os.Key(), failed)
	}
}

func (o *PreviewOperator) preview(set ObjectSet) (r OperateResult, err error) {
	if set.Len() == 0 {
		return
	}

	r = OperateResult{
		Success:   ObjectSet{},
		Failed:    set,
		UnChanged: ObjectSet{},
	}

	for _, group := range o.groupSequence {
		objByScope := set.ByGroup(group)
		if len(objByScope) == 0 {
			continue
		}

		unchanged, notFound, _, err := GetObjects(o.serverContext, group, objByScope, false)
		if err != nil {
			return r, err
		}

		// Unchanged.
		r.UnChanged.Add(unchanged.All()...)

		// Patch.
		successPatched, _, err := PatchObjects(o.serverContext, group, unchanged.ByGroup(group), o.extraBodyParams)
		if err != nil {
			return r, err
		}

		r.Success.Add(successPatched.All()...)
		r.Failed.Remove(successPatched.All()...)

		// Batch create.
		successCreated, _, err := BatchCreateObjects(o.serverContext, group, notFound.ByGroup(group), o.extraBodyParams)
		if err != nil {
			return r, err
		}

		r.Success.Add(successCreated.All()...)
		r.Failed.Remove(successCreated.All()...)
	}

	return r, nil
}

type extraBodyParams map[string]any

func (e extraBodyParams) applyToBody(body map[string]any) map[string]any {
	if e == nil {
		return body
	}

	if body == nil {
		body = make(map[string]any)
	}

	for k, v := range e {
		body[k] = v
	}
	return body
}
