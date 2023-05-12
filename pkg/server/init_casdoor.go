package server

import (
	"context"
	"fmt"

	"github.com/seal-io/seal/pkg/casdoor"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/settings"
	"github.com/seal-io/seal/utils/log"
	"github.com/seal-io/seal/utils/strs"
)

func (r *Server) initCasdoor(ctx context.Context, opts initOptions) error {
	// Short circuit for none first-login.
	var cred casdoor.ApplicationCredential
	if err := settings.CasdoorCred.ValueJSONUnmarshal(ctx, opts.ModelClient, &cred); err != nil {
		return err
	}

	if cred.ClientID != "" && cred.ClientSecret != "" {
		return nil
	}

	// Login the builtin admin with initialized password.
	adminSessions, err := casdoor.SignInUser(ctx, casdoor.BuiltinApp, casdoor.BuiltinOrg,
		casdoor.BuiltinAdmin, casdoor.BuiltinAdminInitPwd)
	if err != nil {
		// Nothing to do if failed login the builtin admin at bootstrap phase.
		return fmt.Errorf("cannot login the builtin admin with init password: %w", err)
	}

	// Get the credential of the builtin application,
	// so that boot the system token creation at below.
	appCred, err := casdoor.GetApplicationCredential(ctx, adminSessions,
		casdoor.BuiltinApp)
	if err != nil {
		return err
	}
	cred.ClientID, cred.ClientSecret = appCred.ClientID, appCred.ClientSecret

	// Create a "never expires" user demand token as system token,
	// the system token is used for internal interaction and password reset.
	token, err := casdoor.CreateToken(ctx, cred.ClientID, cred.ClientSecret,
		casdoor.BuiltinAdmin, nil)
	if err != nil {
		return err
	}

	defer func() {
		// NB(thxCode): revert the token if occurs error,
		// make the `initCasdoor` idempotent.
		if err != nil {
			_ = casdoor.DeleteToken(context.Background(), cred.ClientID, cred.ClientSecret,
				token.Owner, token.Name)
		}
	}()

	// Set bootstrap password.
	adminPassword := r.BootstrapPassword
	if adminPassword == "" {
		adminPassword = strs.Hex(16)
	}

	err = casdoor.UpdateUserPassword(ctx, cred.ClientID, cred.ClientSecret, casdoor.BuiltinOrg, casdoor.BuiltinAdmin,
		"", adminPassword)
	if err != nil {
		return err
	}

	defer func() {
		// NB(thxCode): revert the password if occurs error,
		// make the `initCasdoor` idempotent.
		if err != nil {
			_ = casdoor.UpdateUserPassword(ctx, cred.ClientID, cred.ClientSecret,
				casdoor.BuiltinOrg, casdoor.BuiltinAdmin,
				"", casdoor.BuiltinAdminInitPwd)
		}
	}()

	// Record the application credential.
	err = opts.ModelClient.WithTx(ctx, func(tx *model.Tx) error {
		if _, err = settings.CasdoorCred.Set(ctx, tx, cred); err != nil {
			return err
		}

		if _, err = settings.PrivilegeApiToken.Set(ctx, tx, token.AccessToken); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	if r.BootstrapPassword == "" {
		log.Infof("!!! Bootstrap Admin Password: %s !!!", adminPassword)
	}

	return nil
}
