package server

import (
	"context"
	"fmt"
	"os"

	"github.com/seal-io/walrus/pkg/casdoor"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/settings"
	"github.com/seal-io/walrus/utils/log"
	"github.com/seal-io/walrus/utils/strs"
)

// configureCasdoor initializes the builtin Casdoor application and builtin admin user.
func (r *Server) configureCasdoor(ctx context.Context, opts initOptions) error {
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
		// make the `configureCasdoor` idempotent.
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
		// make the `configureCasdoor` idempotent.
		if err != nil {
			_ = casdoor.UpdateUserPassword(ctx, cred.ClientID, cred.ClientSecret,
				casdoor.BuiltinOrg, casdoor.BuiltinAdmin,
				"", casdoor.BuiltinAdminInitPwd)
		}
	}()

	// Record the application credential.
	err = opts.ModelClient.WithTx(ctx, func(tx *model.Tx) error {
		if err := settings.CasdoorCred.Set(ctx, tx, cred); err != nil {
			return err
		}

		return settings.CasdoorApiToken.Set(ctx, tx, token.AccessToken)
	})
	if err != nil {
		return err
	}

	if r.BootstrapPassword == "" {
		var fl string

		switch {
		case os.Getenv("KUBERNETES_SERVICE_HOST") != "":
			// Running inside a Kubernetes Pod.
			fl = "Kubernetes"
		case os.Getenv("_RUNNING_INSIDE_CONTAINER_") != "":
			// Running inside a container.
			fl = "Docker"
		default:
			// Running as a process.
			fl = "Process"
		}

		err = settings.BootPwdGainSource.Set(ctx, opts.ModelClient, fl)
		if err != nil {
			return err
		}

		log.Infof("!!! Bootstrap Admin Password: %s !!!", adminPassword)
	}

	return nil
}
