// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package model

import (
	"time"

	"github.com/seal-io/seal/pkg/dao/types/oid"
)

// TokenQueryInput is the input for the Token query.
type TokenQueryInput struct {
	// ID holds the value of the "id" field.
	ID oid.ID `uri:"id,omitempty" json:"id,omitempty"`
}

// Model converts the TokenQueryInput to Token.
func (in TokenQueryInput) Model() *Token {
	return &Token{
		ID: in.ID,
	}
}

// TokenCreateInput is the input for the Token creation.
type TokenCreateInput struct {
	// The token name of casdoor.
	CasdoorTokenName string `json:"casdoorTokenName,omitempty"`
	// The token owner of casdoor.
	CasdoorTokenOwner string `json:"casdoorTokenOwner,omitempty"`
	// The name of token.
	Name string `json:"name"`
	// Expiration in seconds.
	Expiration *int `json:"expiration,omitempty"`
}

// Model converts the TokenCreateInput to Token.
func (in TokenCreateInput) Model() *Token {
	var entity = &Token{
		CasdoorTokenName:  in.CasdoorTokenName,
		CasdoorTokenOwner: in.CasdoorTokenOwner,
		Name:              in.Name,
		Expiration:        in.Expiration,
	}
	return entity
}

// TokenUpdateInput is the input for the Token modification.
type TokenUpdateInput struct {
	// ID holds the value of the "id" field.
	ID oid.ID `uri:"id" json:"-"`
	// The token name of casdoor.
	CasdoorTokenName string `json:"casdoorTokenName,omitempty"`
	// The token owner of casdoor.
	CasdoorTokenOwner string `json:"casdoorTokenOwner,omitempty"`
	// The name of token.
	Name string `json:"name,omitempty"`
	// Expiration in seconds.
	Expiration *int `json:"expiration,omitempty"`
}

// Model converts the TokenUpdateInput to Token.
func (in TokenUpdateInput) Model() *Token {
	var entity = &Token{
		ID:                in.ID,
		CasdoorTokenName:  in.CasdoorTokenName,
		CasdoorTokenOwner: in.CasdoorTokenOwner,
		Name:              in.Name,
		Expiration:        in.Expiration,
	}
	return entity
}

// TokenOutput is the output for the Token.
type TokenOutput struct {
	// ID holds the value of the "id" field.
	ID oid.ID `json:"id,omitempty"`
	// Describe creation time.
	CreateTime *time.Time `json:"createTime,omitempty"`
	// Describe modification time.
	UpdateTime *time.Time `json:"updateTime,omitempty"`
	// The name of token.
	Name string `json:"name,omitempty"`
	// Expiration in seconds.
	Expiration *int `json:"expiration,omitempty"`
}

// ExposeToken converts the Token to TokenOutput.
func ExposeToken(in *Token) *TokenOutput {
	if in == nil {
		return nil
	}
	var entity = &TokenOutput{
		ID:         in.ID,
		CreateTime: in.CreateTime,
		UpdateTime: in.UpdateTime,
		Name:       in.Name,
		Expiration: in.Expiration,
	}
	return entity
}

// ExposeTokens converts the Token slice to TokenOutput pointer slice.
func ExposeTokens(in []*Token) []*TokenOutput {
	var out = make([]*TokenOutput, 0, len(in))
	for i := 0; i < len(in); i++ {
		var o = ExposeToken(in[i])
		if o == nil {
			continue
		}
		out = append(out, o)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}
