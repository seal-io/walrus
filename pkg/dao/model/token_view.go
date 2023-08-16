// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "seal". DO NOT EDIT.

package model

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/seal-io/walrus/pkg/dao/model/token"
	"github.com/seal-io/walrus/pkg/dao/types/crypto"
	"github.com/seal-io/walrus/pkg/dao/types/object"
)

// TokenCreateInput holds the creation input of the Token entity,
// please tags with `path:",inline" json:",inline"` if embedding.
type TokenCreateInput struct {
	inputConfig `path:"-" query:"-" json:"-"`

	// The value of token, store in string.
	Value crypto.String `path:"-" query:"-" json:"value"`
	// The name of token.
	Name string `path:"-" query:"-" json:"name"`
	// The kind of token.
	Kind string `path:"-" query:"-" json:"kind,omitempty"`
	// The time of expiration, empty means forever.
	Expiration *time.Time `path:"-" query:"-" json:"expiration,omitempty"`
}

// Model returns the Token entity for creating,
// after validating.
func (tci *TokenCreateInput) Model() *Token {
	if tci == nil {
		return nil
	}

	_t := &Token{
		Value:      tci.Value,
		Name:       tci.Name,
		Kind:       tci.Kind,
		Expiration: tci.Expiration,
	}

	return _t
}

// Validate checks the TokenCreateInput entity.
func (tci *TokenCreateInput) Validate() error {
	if tci == nil {
		return errors.New("nil receiver")
	}

	return tci.ValidateWith(tci.inputConfig.Context, tci.inputConfig.Client, nil)
}

// ValidateWith checks the TokenCreateInput entity with the given context and client set.
func (tci *TokenCreateInput) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if tci == nil {
		return errors.New("nil receiver")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	return nil
}

// TokenCreateInputs holds the creation input item of the Token entities.
type TokenCreateInputsItem struct {
	// The value of token, store in string.
	Value crypto.String `path:"-" query:"-" json:"value"`
	// The name of token.
	Name string `path:"-" query:"-" json:"name"`
	// The kind of token.
	Kind string `path:"-" query:"-" json:"kind,omitempty"`
	// The time of expiration, empty means forever.
	Expiration *time.Time `path:"-" query:"-" json:"expiration,omitempty"`
}

// ValidateWith checks the TokenCreateInputsItem entity with the given context and client set.
func (tci *TokenCreateInputsItem) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if tci == nil {
		return errors.New("nil receiver")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	return nil
}

// TokenCreateInputs holds the creation input of the Token entities,
// please tags with `path:",inline" json:",inline"` if embedding.
type TokenCreateInputs struct {
	inputConfig `path:"-" query:"-" json:"-"`

	// Items holds the entities to create, which MUST not be empty.
	Items []*TokenCreateInputsItem `path:"-" query:"-" json:"items"`
}

// Model returns the Token entities for creating,
// after validating.
func (tci *TokenCreateInputs) Model() []*Token {
	if tci == nil || len(tci.Items) == 0 {
		return nil
	}

	_ts := make([]*Token, len(tci.Items))

	for i := range tci.Items {
		_t := &Token{
			Value:      tci.Items[i].Value,
			Name:       tci.Items[i].Name,
			Kind:       tci.Items[i].Kind,
			Expiration: tci.Items[i].Expiration,
		}

		_ts[i] = _t
	}

	return _ts
}

// Validate checks the TokenCreateInputs entity .
func (tci *TokenCreateInputs) Validate() error {
	if tci == nil {
		return errors.New("nil receiver")
	}

	return tci.ValidateWith(tci.inputConfig.Context, tci.inputConfig.Client, nil)
}

// ValidateWith checks the TokenCreateInputs entity with the given context and client set.
func (tci *TokenCreateInputs) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if tci == nil {
		return errors.New("nil receiver")
	}

	if len(tci.Items) == 0 {
		return errors.New("empty items")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	for i := range tci.Items {
		if tci.Items[i] == nil {
			continue
		}

		if err := tci.Items[i].ValidateWith(ctx, cs, cache); err != nil {
			return err
		}
	}

	return nil
}

// TokenDeleteInput holds the deletion input of the Token entity,
// please tags with `path:",inline"` if embedding.
type TokenDeleteInput struct {
	TokenQueryInput `path:",inline"`
}

// TokenDeleteInputs holds the deletion input item of the Token entities.
type TokenDeleteInputsItem struct {
	// ID of the Token entity.
	ID object.ID `path:"-" query:"-" json:"id"`
}

// TokenDeleteInputs holds the deletion input of the Token entities,
// please tags with `path:",inline" json:",inline"` if embedding.
type TokenDeleteInputs struct {
	inputConfig `path:"-" query:"-" json:"-"`

	// Items holds the entities to create, which MUST not be empty.
	Items []*TokenDeleteInputsItem `path:"-" query:"-" json:"items"`
}

// Model returns the Token entities for deleting,
// after validating.
func (tdi *TokenDeleteInputs) Model() []*Token {
	if tdi == nil || len(tdi.Items) == 0 {
		return nil
	}

	_ts := make([]*Token, len(tdi.Items))
	for i := range tdi.Items {
		_ts[i] = &Token{
			ID: tdi.Items[i].ID,
		}
	}
	return _ts
}

// IDs returns the ID list of the Token entities for deleting,
// after validating.
func (tdi *TokenDeleteInputs) IDs() []object.ID {
	if tdi == nil || len(tdi.Items) == 0 {
		return nil
	}

	ids := make([]object.ID, len(tdi.Items))
	for i := range tdi.Items {
		ids[i] = tdi.Items[i].ID
	}
	return ids
}

// Validate checks the TokenDeleteInputs entity.
func (tdi *TokenDeleteInputs) Validate() error {
	if tdi == nil {
		return errors.New("nil receiver")
	}

	return tdi.ValidateWith(tdi.inputConfig.Context, tdi.inputConfig.Client, nil)
}

// ValidateWith checks the TokenDeleteInputs entity with the given context and client set.
func (tdi *TokenDeleteInputs) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if tdi == nil {
		return errors.New("nil receiver")
	}

	if len(tdi.Items) == 0 {
		return errors.New("empty items")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	q := cs.Tokens().Query()

	ids := make([]object.ID, 0, len(tdi.Items))

	for i := range tdi.Items {
		if tdi.Items[i] == nil {
			return errors.New("nil item")
		}

		if tdi.Items[i].ID != "" {
			ids = append(ids, tdi.Items[i].ID)
		} else {
			return errors.New("found item hasn't identify")
		}
	}

	if len(ids) != cap(ids) {
		return errors.New("found unrecognized item")
	}

	idsCnt, err := q.Where(token.IDIn(ids...)).
		Count(ctx)
	if err != nil {
		return err
	}

	if idsCnt != cap(ids) {
		return errors.New("found unrecognized item")
	}

	return nil
}

// TokenQueryInput holds the query input of the Token entity,
// please tags with `path:",inline"` if embedding.
type TokenQueryInput struct {
	inputConfig `path:"-" query:"-" json:"-"`

	// Refer holds the route path reference of the Token entity.
	Refer *object.Refer `path:"token,default=" query:"-" json:"-"`
	// ID of the Token entity.
	ID object.ID `path:"-" query:"-" json:"id"`
}

// Model returns the Token entity for querying,
// after validating.
func (tqi *TokenQueryInput) Model() *Token {
	if tqi == nil {
		return nil
	}

	return &Token{
		ID: tqi.ID,
	}
}

// Validate checks the TokenQueryInput entity.
func (tqi *TokenQueryInput) Validate() error {
	if tqi == nil {
		return errors.New("nil receiver")
	}

	return tqi.ValidateWith(tqi.inputConfig.Context, tqi.inputConfig.Client, nil)
}

// ValidateWith checks the TokenQueryInput entity with the given context and client set.
func (tqi *TokenQueryInput) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if tqi == nil {
		return errors.New("nil receiver")
	}

	if tqi.Refer != nil && *tqi.Refer == "" {
		return fmt.Errorf("model: %s : %w", token.Label, ErrBlankResourceRefer)
	}

	if cache == nil {
		cache = map[string]any{}
	}

	q := cs.Tokens().Query()

	if tqi.Refer != nil {
		if tqi.Refer.IsID() {
			q.Where(
				token.ID(tqi.Refer.ID()))
		} else {
			return errors.New("invalid identify refer of token")
		}
	} else if tqi.ID != "" {
		q.Where(
			token.ID(tqi.ID))
	} else {
		return errors.New("invalid identify of token")
	}

	q.Select(
		token.FieldID,
	)

	var e *Token
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
			e = cv.(*Token)
		}
	}

	tqi.ID = e.ID
	return nil
}

// TokenQueryInputs holds the query input of the Token entities,
// please tags with `path:",inline" query:",inline"` if embedding.
type TokenQueryInputs struct {
	inputConfig `path:"-" query:"-" json:"-"`
}

// Validate checks the TokenQueryInputs entity.
func (tqi *TokenQueryInputs) Validate() error {
	if tqi == nil {
		return errors.New("nil receiver")
	}

	return tqi.ValidateWith(tqi.inputConfig.Context, tqi.inputConfig.Client, nil)
}

// ValidateWith checks the TokenQueryInputs entity with the given context and client set.
func (tqi *TokenQueryInputs) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if tqi == nil {
		return errors.New("nil receiver")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	return nil
}

// TokenUpdateInput holds the modification input of the Token entity,
// please tags with `path:",inline" json:",inline"` if embedding.
type TokenUpdateInput struct {
	TokenQueryInput `path:",inline" query:"-" json:"-"`
}

// Model returns the Token entity for modifying,
// after validating.
func (tui *TokenUpdateInput) Model() *Token {
	if tui == nil {
		return nil
	}

	_t := &Token{
		ID: tui.ID,
	}

	return _t
}

// Validate checks the TokenUpdateInput entity.
func (tui *TokenUpdateInput) Validate() error {
	if tui == nil {
		return errors.New("nil receiver")
	}

	return tui.ValidateWith(tui.inputConfig.Context, tui.inputConfig.Client, nil)
}

// ValidateWith checks the TokenUpdateInput entity with the given context and client set.
func (tui *TokenUpdateInput) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if cache == nil {
		cache = map[string]any{}
	}

	if err := tui.TokenQueryInput.ValidateWith(ctx, cs, cache); err != nil {
		return err
	}

	return nil
}

// TokenUpdateInputs holds the modification input item of the Token entities.
type TokenUpdateInputsItem struct {
	// ID of the Token entity.
	ID object.ID `path:"-" query:"-" json:"id"`
}

// ValidateWith checks the TokenUpdateInputsItem entity with the given context and client set.
func (tui *TokenUpdateInputsItem) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if tui == nil {
		return errors.New("nil receiver")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	return nil
}

// TokenUpdateInputs holds the modification input of the Token entities,
// please tags with `path:",inline" json:",inline"` if embedding.
type TokenUpdateInputs struct {
	inputConfig `path:"-" query:"-" json:"-"`

	// Items holds the entities to create, which MUST not be empty.
	Items []*TokenUpdateInputsItem `path:"-" query:"-" json:"items"`
}

// Model returns the Token entities for modifying,
// after validating.
func (tui *TokenUpdateInputs) Model() []*Token {
	if tui == nil || len(tui.Items) == 0 {
		return nil
	}

	_ts := make([]*Token, len(tui.Items))

	for i := range tui.Items {
		_t := &Token{
			ID: tui.Items[i].ID,
		}

		_ts[i] = _t
	}

	return _ts
}

// IDs returns the ID list of the Token entities for modifying,
// after validating.
func (tui *TokenUpdateInputs) IDs() []object.ID {
	if tui == nil || len(tui.Items) == 0 {
		return nil
	}

	ids := make([]object.ID, len(tui.Items))
	for i := range tui.Items {
		ids[i] = tui.Items[i].ID
	}
	return ids
}

// Validate checks the TokenUpdateInputs entity.
func (tui *TokenUpdateInputs) Validate() error {
	if tui == nil {
		return errors.New("nil receiver")
	}

	return tui.ValidateWith(tui.inputConfig.Context, tui.inputConfig.Client, nil)
}

// ValidateWith checks the TokenUpdateInputs entity with the given context and client set.
func (tui *TokenUpdateInputs) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if tui == nil {
		return errors.New("nil receiver")
	}

	if len(tui.Items) == 0 {
		return errors.New("empty items")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	q := cs.Tokens().Query()

	ids := make([]object.ID, 0, len(tui.Items))

	for i := range tui.Items {
		if tui.Items[i] == nil {
			return errors.New("nil item")
		}

		if tui.Items[i].ID != "" {
			ids = append(ids, tui.Items[i].ID)
		} else {
			return errors.New("found item hasn't identify")
		}
	}

	if len(ids) != cap(ids) {
		return errors.New("found unrecognized item")
	}

	idsCnt, err := q.Where(token.IDIn(ids...)).
		Count(ctx)
	if err != nil {
		return err
	}

	if idsCnt != cap(ids) {
		return errors.New("found unrecognized item")
	}

	for i := range tui.Items {
		if tui.Items[i] == nil {
			continue
		}

		if err := tui.Items[i].ValidateWith(ctx, cs, cache); err != nil {
			return err
		}
	}

	return nil
}

// TokenOutput holds the output of the Token entity.
type TokenOutput struct {
	ID          object.ID  `json:"id,omitempty"`
	CreateTime  *time.Time `json:"createTime,omitempty"`
	Kind        string     `json:"kind,omitempty"`
	Name        string     `json:"name,omitempty"`
	Expiration  *time.Time `json:"expiration,omitempty"`
	AccessToken string     `json:"accessToken,omitempty"`

	Subject *SubjectOutput `json:"subject,omitempty"`
}

// View returns the output of Token entity.
func (_t *Token) View() *TokenOutput {
	return ExposeToken(_t)
}

// View returns the output of Token entities.
func (_ts Tokens) View() []*TokenOutput {
	return ExposeTokens(_ts)
}

// ExposeToken converts the Token to TokenOutput.
func ExposeToken(_t *Token) *TokenOutput {
	if _t == nil {
		return nil
	}

	to := &TokenOutput{
		ID:          _t.ID,
		CreateTime:  _t.CreateTime,
		Kind:        _t.Kind,
		Name:        _t.Name,
		Expiration:  _t.Expiration,
		AccessToken: _t.AccessToken,
	}

	if _t.Edges.Subject != nil {
		to.Subject = ExposeSubject(_t.Edges.Subject)
	} else if _t.SubjectID != "" {
		to.Subject = &SubjectOutput{
			ID: _t.SubjectID,
		}
	}
	return to
}

// ExposeTokens converts the Token slice to TokenOutput pointer slice.
func ExposeTokens(_ts []*Token) []*TokenOutput {
	if len(_ts) == 0 {
		return nil
	}

	tos := make([]*TokenOutput, len(_ts))
	for i := range _ts {
		tos[i] = ExposeToken(_ts[i])
	}
	return tos
}
