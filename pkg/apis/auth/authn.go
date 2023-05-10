package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/apis/auth/cache"
	"github.com/seal-io/seal/pkg/apis/auth/session"
	"github.com/seal-io/seal/utils/log"
	"github.com/seal-io/seal/utils/req"

	"github.com/seal-io/seal/pkg/casdoor"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/subject"
	"github.com/seal-io/seal/pkg/settings"
)

func authn(c *gin.Context, modelClient model.ClientSet) error {
	token := casdoor.GetToken(c.Request.Header)
	if token != "" {
		return authnWithToken(c, modelClient, token)
	}

	internalSession := casdoor.GetSession(c.Request.Cookies())
	if internalSession == nil {
		// Anonymous.
		return nil
	}
	return authnWithSession(c, modelClient, internalSession)
}

func authnWithToken(c *gin.Context, modelClient model.ClientSet, token string) error {
	logger := log.WithName("api").WithName("auth")

	if sj, active := cache.LoadTokenSubject(token); sj != nil {
		if !active {
			// Anonymous.
			return nil
		}
		g, n, err := session.ParseSubjectKey(*sj)
		if err != nil {
			return fmt.Errorf("failed to parse subject key: %w", err)
		}
		groups, err := getGroups(c, modelClient, g, n)
		if err != nil {
			return err
		}
		session.StoreSubjectAuthnInfo(c, groups, n)
		return nil
	}

	var cred casdoor.ApplicationCredential
	if err := settings.CasdoorCred.ValueJSONUnmarshal(c, modelClient, &cred); err != nil {
		return fmt.Errorf("failed to unmarshal casdoor secret: %w", err)
	}
	r, err := casdoor.IntrospectToken(c, cred.ClientID, cred.ClientSecret, token)
	if err != nil {
		// Avoid d-dos.
		logger.Errorf("error verifying user token: %v", err)
		cache.StoreTokenSubject(token, "", false)
		return nil
	}
	if !r.Active || r.Exp < time.Now().Unix() {
		// Expired.
		cache.StoreTokenSubject(token, "", false)
		return nil
	}
	// Cache.
	groups, err := getGroups(c, modelClient, "", r.UserName)
	if err != nil {
		return err
	}
	cache.StoreTokenSubject(token, session.ToSubjectKey(groups[len(groups)-1], r.UserName), true)
	session.StoreSubjectAuthnInfo(c, groups, r.UserName)
	return nil
}

func authnWithSession(c *gin.Context, modelClient model.ClientSet, internalSession *req.HttpCookie) error {
	logger := log.WithName("api").WithName("auth")

	if sj, active := cache.LoadSessionSubject(string(internalSession.Value())); sj != nil {
		if !active {
			// Anonymous.
			return nil
		}
		g, n, err := session.ParseSubjectKey(*sj)
		if err != nil {
			return fmt.Errorf("failed to parse subject key: %w", err)
		}
		groups, err := getGroups(c, modelClient, g, n)
		if err != nil {
			return err
		}
		session.StoreSubjectAuthnInfo(c, groups, n)
		return nil
	}

	r, err := casdoor.GetUserInfo(c, []*req.HttpCookie{internalSession})
	if err != nil {
		// Avoid d-dos.
		logger.Errorf("error getting user account: %v", err)
		cache.StoreSessionSubject(string(internalSession.Value()), "", false)
		return casdoor.InterruptSession(c.Writer, []*req.HttpCookie{internalSession})
	}
	// Cache.
	groups, err := getGroups(c, modelClient, "", r.Name)
	if err != nil {
		return err
	}
	cache.StoreSessionSubject(string(internalSession.Value()),
		session.ToSubjectKey(groups[len(groups)-1], r.Name), true)
	session.StoreSubjectAuthnInfo(c, groups, r.Name)
	return casdoor.HoldSession(c.Writer, []*req.HttpCookie{internalSession})
}

// getGroups returns the groups with the given user,
// if not group is blank, retries the proper groups.
func getGroups(ctx context.Context, modelClient model.ClientSet, group string, user string) ([]string, error) {
	query := modelClient.Subjects().Query().
		Where(subject.Kind("user"), subject.Name(user))
	if group == "" {
		// Get specified login group(loginTo=true) or default login group(mountTo=false).
		query.Where(subject.Or(subject.LoginTo(true), subject.MountTo(false)))
	} else {
		// Specified group.
		query.Where(subject.Group(group))
	}

	users, err := query.
		Select(
			subject.FieldLoginTo,
			subject.FieldMountTo,
			subject.FieldPaths).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	var groups []string
	for i := 0; i < len(users); i++ {
		u := users[i]
		if *u.LoginTo {
			return u.Paths[:len(u.Paths)-1], nil
		}
		if !*u.MountTo {
			groups = u.Paths[:len(u.Paths)-1]
		}
	}
	return groups, nil
}
