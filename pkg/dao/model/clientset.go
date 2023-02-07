// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package model

import "context"

// ClientSet is an interface that allows getting all clients.
type ClientSet interface {
	// Roles returns the client for interacting with the Role builders.
	Roles() *RoleClient

	// Settings returns the client for interacting with the Setting builders.
	Settings() *SettingClient

	// Subjects returns the client for interacting with the Subject builders.
	Subjects() *SubjectClient

	// Tokens returns the client for interacting with the Token builders.
	Tokens() *TokenClient

	// WithTx gives a new transactional client in the callback function,
	// if already in a transaction, this will keep in the same transaction.
	WithTx(context.Context, func(tx *Tx) error) error

	// Dialect returns the dialect name of the driver.
	Dialect() string
}

// RoleClientGetter is an interface that allows getting RoleClient.
type RoleClientGetter interface {
	// Roles returns the client for interacting with the Role builders.
	Roles() *RoleClient
}

// SettingClientGetter is an interface that allows getting SettingClient.
type SettingClientGetter interface {
	// Settings returns the client for interacting with the Setting builders.
	Settings() *SettingClient
}

// SubjectClientGetter is an interface that allows getting SubjectClient.
type SubjectClientGetter interface {
	// Subjects returns the client for interacting with the Subject builders.
	Subjects() *SubjectClient
}

// TokenClientGetter is an interface that allows getting TokenClient.
type TokenClientGetter interface {
	// Tokens returns the client for interacting with the Token builders.
	Tokens() *TokenClient
}
