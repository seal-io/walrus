// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "seal". DO NOT EDIT.

package model

import (
	"context"
	"errors"
	"fmt"

	"github.com/seal-io/walrus/pkg/dao/model/distributelock"
	"github.com/seal-io/walrus/pkg/dao/types/object"
)

// DistributeLockCreateInput holds the creation input of the DistributeLock entity,
// please tags with `path:",inline" json:",inline"` if embedding.
type DistributeLockCreateInput struct {
	inputConfig `path:"-" query:"-" json:"-"`

	// Expiration timestamp to prevent the lock be occupied for long time.
	ExpireAt int64 `path:"-" query:"-" json:"expireAt"`
}

// Model returns the DistributeLock entity for creating,
// after validating.
func (dlci *DistributeLockCreateInput) Model() *DistributeLock {
	if dlci == nil {
		return nil
	}

	_dl := &DistributeLock{
		ExpireAt: dlci.ExpireAt,
	}

	return _dl
}

// Validate checks the DistributeLockCreateInput entity.
func (dlci *DistributeLockCreateInput) Validate() error {
	if dlci == nil {
		return errors.New("nil receiver")
	}

	return dlci.ValidateWith(dlci.inputConfig.Context, dlci.inputConfig.Client, nil)
}

// ValidateWith checks the DistributeLockCreateInput entity with the given context and client set.
func (dlci *DistributeLockCreateInput) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if dlci == nil {
		return errors.New("nil receiver")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	return nil
}

// DistributeLockCreateInputs holds the creation input item of the DistributeLock entities.
type DistributeLockCreateInputsItem struct {
	// Expiration timestamp to prevent the lock be occupied for long time.
	ExpireAt int64 `path:"-" query:"-" json:"expireAt"`
}

// ValidateWith checks the DistributeLockCreateInputsItem entity with the given context and client set.
func (dlci *DistributeLockCreateInputsItem) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if dlci == nil {
		return errors.New("nil receiver")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	return nil
}

// DistributeLockCreateInputs holds the creation input of the DistributeLock entities,
// please tags with `path:",inline" json:",inline"` if embedding.
type DistributeLockCreateInputs struct {
	inputConfig `path:"-" query:"-" json:"-"`

	// Items holds the entities to create, which MUST not be empty.
	Items []*DistributeLockCreateInputsItem `path:"-" query:"-" json:"items"`
}

// Model returns the DistributeLock entities for creating,
// after validating.
func (dlci *DistributeLockCreateInputs) Model() []*DistributeLock {
	if dlci == nil || len(dlci.Items) == 0 {
		return nil
	}

	_dls := make([]*DistributeLock, len(dlci.Items))

	for i := range dlci.Items {
		_dl := &DistributeLock{
			ExpireAt: dlci.Items[i].ExpireAt,
		}

		_dls[i] = _dl
	}

	return _dls
}

// Validate checks the DistributeLockCreateInputs entity .
func (dlci *DistributeLockCreateInputs) Validate() error {
	if dlci == nil {
		return errors.New("nil receiver")
	}

	return dlci.ValidateWith(dlci.inputConfig.Context, dlci.inputConfig.Client, nil)
}

// ValidateWith checks the DistributeLockCreateInputs entity with the given context and client set.
func (dlci *DistributeLockCreateInputs) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if dlci == nil {
		return errors.New("nil receiver")
	}

	if len(dlci.Items) == 0 {
		return errors.New("empty items")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	for i := range dlci.Items {
		if dlci.Items[i] == nil {
			continue
		}

		if err := dlci.Items[i].ValidateWith(ctx, cs, cache); err != nil {
			return err
		}
	}

	return nil
}

// DistributeLockDeleteInput holds the deletion input of the DistributeLock entity,
// please tags with `path:",inline"` if embedding.
type DistributeLockDeleteInput struct {
	DistributeLockQueryInput `path:",inline"`
}

// DistributeLockDeleteInputs holds the deletion input item of the DistributeLock entities.
type DistributeLockDeleteInputsItem struct {
	// ID of the DistributeLock entity.
	ID string `path:"-" query:"-" json:"id"`
}

// DistributeLockDeleteInputs holds the deletion input of the DistributeLock entities,
// please tags with `path:",inline" json:",inline"` if embedding.
type DistributeLockDeleteInputs struct {
	inputConfig `path:"-" query:"-" json:"-"`

	// Items holds the entities to create, which MUST not be empty.
	Items []*DistributeLockDeleteInputsItem `path:"-" query:"-" json:"items"`
}

// Model returns the DistributeLock entities for deleting,
// after validating.
func (dldi *DistributeLockDeleteInputs) Model() []*DistributeLock {
	if dldi == nil || len(dldi.Items) == 0 {
		return nil
	}

	_dls := make([]*DistributeLock, len(dldi.Items))
	for i := range dldi.Items {
		_dls[i] = &DistributeLock{
			ID: dldi.Items[i].ID,
		}
	}
	return _dls
}

// IDs returns the ID list of the DistributeLock entities for deleting,
// after validating.
func (dldi *DistributeLockDeleteInputs) IDs() []string {
	if dldi == nil || len(dldi.Items) == 0 {
		return nil
	}

	ids := make([]string, len(dldi.Items))
	for i := range dldi.Items {
		ids[i] = dldi.Items[i].ID
	}
	return ids
}

// Validate checks the DistributeLockDeleteInputs entity.
func (dldi *DistributeLockDeleteInputs) Validate() error {
	if dldi == nil {
		return errors.New("nil receiver")
	}

	return dldi.ValidateWith(dldi.inputConfig.Context, dldi.inputConfig.Client, nil)
}

// ValidateWith checks the DistributeLockDeleteInputs entity with the given context and client set.
func (dldi *DistributeLockDeleteInputs) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if dldi == nil {
		return errors.New("nil receiver")
	}

	if len(dldi.Items) == 0 {
		return errors.New("empty items")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	q := cs.DistributeLocks().Query()

	ids := make([]string, 0, len(dldi.Items))

	for i := range dldi.Items {
		if dldi.Items[i] == nil {
			return errors.New("nil item")
		}

		if dldi.Items[i].ID != "" {
			ids = append(ids, dldi.Items[i].ID)
		} else {
			return errors.New("found item hasn't identify")
		}
	}

	if len(ids) != cap(ids) {
		return errors.New("found unrecognized item")
	}

	idsCnt, err := q.Where(distributelock.IDIn(ids...)).
		Count(ctx)
	if err != nil {
		return err
	}

	if idsCnt != cap(ids) {
		return errors.New("found unrecognized item")
	}

	return nil
}

// DistributeLockQueryInput holds the query input of the DistributeLock entity,
// please tags with `path:",inline"` if embedding.
type DistributeLockQueryInput struct {
	inputConfig `path:"-" query:"-" json:"-"`

	// Refer holds the route path reference of the DistributeLock entity.
	Refer *object.Refer `path:"distributelock,default=" query:"-" json:"-"`
	// ID of the DistributeLock entity.
	ID string `path:"-" query:"-" json:"id"`
}

// Model returns the DistributeLock entity for querying,
// after validating.
func (dlqi *DistributeLockQueryInput) Model() *DistributeLock {
	if dlqi == nil {
		return nil
	}

	return &DistributeLock{
		ID: dlqi.ID,
	}
}

// Validate checks the DistributeLockQueryInput entity.
func (dlqi *DistributeLockQueryInput) Validate() error {
	if dlqi == nil {
		return errors.New("nil receiver")
	}

	return dlqi.ValidateWith(dlqi.inputConfig.Context, dlqi.inputConfig.Client, nil)
}

// ValidateWith checks the DistributeLockQueryInput entity with the given context and client set.
func (dlqi *DistributeLockQueryInput) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if dlqi == nil {
		return errors.New("nil receiver")
	}

	if dlqi.Refer != nil && *dlqi.Refer == "" {
		return fmt.Errorf("model: %s : %w", distributelock.Label, ErrBlankResourceRefer)
	}

	if cache == nil {
		cache = map[string]any{}
	}

	q := cs.DistributeLocks().Query()

	if dlqi.Refer != nil {
		if dlqi.Refer.IsString() {
			q.Where(
				distributelock.ID(dlqi.Refer.String()))
		} else {
			return errors.New("invalid identify refer of distributelock")
		}
	} else if dlqi.ID != "" {
		q.Where(
			distributelock.ID(dlqi.ID))
	} else {
		return errors.New("invalid identify of distributelock")
	}

	q.Select(
		distributelock.FieldID,
	)

	var e *DistributeLock
	{
		// Get cache from previous validation.
		queryStmt, queryArgs := q.sqlQuery(setContextOp(ctx, q.ctx, "cache")).Query()
		ck := fmt.Sprintf("stmt=%v, args=%v", queryStmt, queryArgs)
		if cv, existed := cache[ck]; !existed {
			var err error
			e, err = q.Only(ctx)
			if err != nil {
				return err
			}

			// Set cache for other validation.
			cache[ck] = e
		} else {
			e = cv.(*DistributeLock)
		}
	}

	dlqi.ID = e.ID
	return nil
}

// DistributeLockQueryInputs holds the query input of the DistributeLock entities,
// please tags with `path:",inline" query:",inline"` if embedding.
type DistributeLockQueryInputs struct {
	inputConfig `path:"-" query:"-" json:"-"`
}

// Validate checks the DistributeLockQueryInputs entity.
func (dlqi *DistributeLockQueryInputs) Validate() error {
	if dlqi == nil {
		return errors.New("nil receiver")
	}

	return dlqi.ValidateWith(dlqi.inputConfig.Context, dlqi.inputConfig.Client, nil)
}

// ValidateWith checks the DistributeLockQueryInputs entity with the given context and client set.
func (dlqi *DistributeLockQueryInputs) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if dlqi == nil {
		return errors.New("nil receiver")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	return nil
}

// DistributeLockUpdateInput holds the modification input of the DistributeLock entity,
// please tags with `path:",inline" json:",inline"` if embedding.
type DistributeLockUpdateInput struct {
	DistributeLockQueryInput `path:",inline" query:"-" json:"-"`

	// Expiration timestamp to prevent the lock be occupied for long time.
	ExpireAt int64 `path:"-" query:"-" json:"expireAt,omitempty"`
}

// Model returns the DistributeLock entity for modifying,
// after validating.
func (dlui *DistributeLockUpdateInput) Model() *DistributeLock {
	if dlui == nil {
		return nil
	}

	_dl := &DistributeLock{
		ID:       dlui.ID,
		ExpireAt: dlui.ExpireAt,
	}

	return _dl
}

// Validate checks the DistributeLockUpdateInput entity.
func (dlui *DistributeLockUpdateInput) Validate() error {
	if dlui == nil {
		return errors.New("nil receiver")
	}

	return dlui.ValidateWith(dlui.inputConfig.Context, dlui.inputConfig.Client, nil)
}

// ValidateWith checks the DistributeLockUpdateInput entity with the given context and client set.
func (dlui *DistributeLockUpdateInput) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if cache == nil {
		cache = map[string]any{}
	}

	if err := dlui.DistributeLockQueryInput.ValidateWith(ctx, cs, cache); err != nil {
		return err
	}

	return nil
}

// DistributeLockUpdateInputs holds the modification input item of the DistributeLock entities.
type DistributeLockUpdateInputsItem struct {
	// ID of the DistributeLock entity.
	ID string `path:"-" query:"-" json:"id"`

	// Expiration timestamp to prevent the lock be occupied for long time.
	ExpireAt int64 `path:"-" query:"-" json:"expireAt"`
}

// ValidateWith checks the DistributeLockUpdateInputsItem entity with the given context and client set.
func (dlui *DistributeLockUpdateInputsItem) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if dlui == nil {
		return errors.New("nil receiver")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	return nil
}

// DistributeLockUpdateInputs holds the modification input of the DistributeLock entities,
// please tags with `path:",inline" json:",inline"` if embedding.
type DistributeLockUpdateInputs struct {
	inputConfig `path:"-" query:"-" json:"-"`

	// Items holds the entities to create, which MUST not be empty.
	Items []*DistributeLockUpdateInputsItem `path:"-" query:"-" json:"items"`
}

// Model returns the DistributeLock entities for modifying,
// after validating.
func (dlui *DistributeLockUpdateInputs) Model() []*DistributeLock {
	if dlui == nil || len(dlui.Items) == 0 {
		return nil
	}

	_dls := make([]*DistributeLock, len(dlui.Items))

	for i := range dlui.Items {
		_dl := &DistributeLock{
			ID:       dlui.Items[i].ID,
			ExpireAt: dlui.Items[i].ExpireAt,
		}

		_dls[i] = _dl
	}

	return _dls
}

// IDs returns the ID list of the DistributeLock entities for modifying,
// after validating.
func (dlui *DistributeLockUpdateInputs) IDs() []string {
	if dlui == nil || len(dlui.Items) == 0 {
		return nil
	}

	ids := make([]string, len(dlui.Items))
	for i := range dlui.Items {
		ids[i] = dlui.Items[i].ID
	}
	return ids
}

// Validate checks the DistributeLockUpdateInputs entity.
func (dlui *DistributeLockUpdateInputs) Validate() error {
	if dlui == nil {
		return errors.New("nil receiver")
	}

	return dlui.ValidateWith(dlui.inputConfig.Context, dlui.inputConfig.Client, nil)
}

// ValidateWith checks the DistributeLockUpdateInputs entity with the given context and client set.
func (dlui *DistributeLockUpdateInputs) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if dlui == nil {
		return errors.New("nil receiver")
	}

	if len(dlui.Items) == 0 {
		return errors.New("empty items")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	q := cs.DistributeLocks().Query()

	ids := make([]string, 0, len(dlui.Items))

	for i := range dlui.Items {
		if dlui.Items[i] == nil {
			return errors.New("nil item")
		}

		if dlui.Items[i].ID != "" {
			ids = append(ids, dlui.Items[i].ID)
		} else {
			return errors.New("found item hasn't identify")
		}
	}

	if len(ids) != cap(ids) {
		return errors.New("found unrecognized item")
	}

	idsCnt, err := q.Where(distributelock.IDIn(ids...)).
		Count(ctx)
	if err != nil {
		return err
	}

	if idsCnt != cap(ids) {
		return errors.New("found unrecognized item")
	}

	for i := range dlui.Items {
		if dlui.Items[i] == nil {
			continue
		}

		if err := dlui.Items[i].ValidateWith(ctx, cs, cache); err != nil {
			return err
		}
	}

	return nil
}

// DistributeLockOutput holds the output of the DistributeLock entity.
type DistributeLockOutput struct {
	ID       string `json:"id,omitempty"`
	ExpireAt int64  `json:"expireAt,omitempty"`
}

// View returns the output of DistributeLock entity.
func (_dl *DistributeLock) View() *DistributeLockOutput {
	return ExposeDistributeLock(_dl)
}

// View returns the output of DistributeLock entities.
func (_dls DistributeLocks) View() []*DistributeLockOutput {
	return ExposeDistributeLocks(_dls)
}

// ExposeDistributeLock converts the DistributeLock to DistributeLockOutput.
func ExposeDistributeLock(_dl *DistributeLock) *DistributeLockOutput {
	if _dl == nil {
		return nil
	}

	dlo := &DistributeLockOutput{
		ID:       _dl.ID,
		ExpireAt: _dl.ExpireAt,
	}

	return dlo
}

// ExposeDistributeLocks converts the DistributeLock slice to DistributeLockOutput pointer slice.
func ExposeDistributeLocks(_dls []*DistributeLock) []*DistributeLockOutput {
	if len(_dls) == 0 {
		return nil
	}

	dlos := make([]*DistributeLockOutput, len(_dls))
	for i := range _dls {
		dlos[i] = ExposeDistributeLock(_dls[i])
	}
	return dlos
}
