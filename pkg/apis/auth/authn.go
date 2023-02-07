package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/apis/auth/cache"
	"github.com/seal-io/seal/pkg/apis/auth/session"
	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/utils/log"
	"github.com/seal-io/seal/utils/req"

	"github.com/seal-io/seal/pkg/casdoor"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/subject"
	"github.com/seal-io/seal/pkg/settings"
)

func authn(c *gin.Context, modelClient model.ClientSet) error {
	var token = casdoor.GetInternalToken(c.Request.Header)
	if token != "" {
		return authnWithToken(c, modelClient, token)
	}

	var internalSession = casdoor.GetInternalSession(c.Request.Cookies())
	if internalSession == nil {
		// anonymous
		return nil
	}
	return authnWithSession(c, modelClient, internalSession)
}

func authnWithToken(c *gin.Context, modelClient model.ClientSet, token string) error {
	if sj, active := cache.LoadTokenSubject(token); sj != nil {
		if !active {
			// anonymous
			return nil
		}
		var g, n, err = session.ParseSubjectKey(*sj)
		if err != nil {
			return runtime.ErrorfP(http.StatusInternalServerError, "failed to parse subject key: %w", err)
		}
		session.StoreSubjectAuthnInfo(c, g, n)
		return nil
	}

	var cred casdoor.ApplicationCredential
	if err := settings.CasdoorCred.ValueJSONUnmarshal(c, modelClient, &cred); err != nil {
		return runtime.ErrorfP(http.StatusInternalServerError, "failed to unmarshal casdoor secret: %w", err)
	}
	var r, err = casdoor.IntrospectToken(c, cred.ClientID, cred.ClientSecret, token)
	if err != nil {
		// avoid d-dos
		log.WithName("auth").Errorf("error verifying user token: %v", err)
		cache.StoreTokenSubject(token, "", false)
		return nil
	}
	if !r.Active || r.Exp < time.Now().Unix() {
		// expired
		cache.StoreTokenSubject(token, "", false)
		return nil
	}
	// cache
	loginGroup, err := getLoginGroup(c, modelClient, r.UserName)
	if err != nil {
		return err
	}
	cache.StoreTokenSubject(token, session.ToSubjectKey(loginGroup, r.UserName), true)
	session.StoreSubjectAuthnInfo(c, loginGroup, r.UserName)
	return nil
}

func authnWithSession(c *gin.Context, modelClient model.ClientSet, internalSession *req.HttpCookie) error {
	if sj, active := cache.LoadSessionSubject(string(internalSession.Value())); sj != nil {
		if !active {
			// anonymous
			return nil
		}
		var g, n, err = session.ParseSubjectKey(*sj)
		if err != nil {
			return runtime.ErrorfP(http.StatusInternalServerError, "failed to parse subject key: %w", err)
		}
		session.StoreSubjectAuthnInfo(c, g, n)
		return nil
	}

	var r, err = casdoor.GetUserInfo(c, []*req.HttpCookie{internalSession})
	if err != nil {
		// avoid d-dos
		log.WithName("auth").Errorf("error getting user account: %v", err)
		cache.StoreSessionSubject(string(internalSession.Value()), "", false)
		return nil
	}
	// cache
	loginGroup, err := getLoginGroup(c, modelClient, r.Name)
	if err != nil {
		return err
	}
	cache.StoreSessionSubject(string(internalSession.Value()), session.ToSubjectKey(loginGroup, r.Name), true)
	session.StoreSubjectAuthnInfo(c, loginGroup, r.Name)
	return nil
}

func getLoginGroup(ctx context.Context, modelClient model.ClientSet, name string) (string, error) {
	var users, err = modelClient.Subjects().Query().
		Where(
			subject.Kind("user"),
			subject.Name(name),
			subject.Or(
				subject.LoginTo(true),
				subject.MountTo(false),
			),
		).
		Select(subject.FieldGroup, subject.FieldLoginTo, subject.FieldMountTo).
		All(ctx)
	if err != nil {
		return "", runtime.ErrorfP(http.StatusInternalServerError, "failed to get user: %w", err)
	}
	var loginGroup string
	for i := 0; i < len(users); i++ {
		var u = users[i]
		if *u.LoginTo {
			loginGroup = u.Group
			break
		} else if !*u.MountTo {
			loginGroup = u.Group
		}
	}
	return loginGroup, nil
}
