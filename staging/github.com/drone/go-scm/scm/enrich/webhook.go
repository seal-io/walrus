// Copyright 2022 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package enrich

import (
	"context"

	"github.com/drone/go-scm/scm"
)

// Webhook enriches the webhook payload with missing
// information not included in the webhook payload.
func Webhook(ctx context.Context, client *scm.Client, webhook *scm.Webhook) error {
	return nil // TODO
}
