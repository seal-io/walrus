package auths

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types/crypto"
	"github.com/seal-io/seal/pkg/dao/types/oid"
	"github.com/seal-io/seal/utils/strs"
)

func decodeToken(r *http.Request) (subjectID, tokenID oid.ID, tokenValue string) {
	v := getAccessToken(r)
	if v == "" {
		return
	}

	subjectID, tokenID, tokenValue, _ = unwrapAccessToken(v)

	return
}

func getAccessToken(r *http.Request) string {
	// Get basic auth.
	v := strings.TrimSpace(r.Header.Get("Authorization"))
	if v == "" {
		return ""
	}

	switch {
	case strings.Contains(v, "Basic "):
		// Basic auth.
		bb64 := strings.TrimPrefix(v, "Basic ")

		// Decode basic auth.
		b, err := strs.DecodeBase64(bb64)
		if err != nil {
			return ""
		}

		// Get token.
		ss := strings.SplitN(b, ":", 2)
		if len(ss) == 2 {
			return ss[1]
		}
	case strings.Contains(v, "Bearer "):
		// Bearer token.
		return strings.TrimPrefix(v, "Bearer ")
	}

	return ""
}

func unwrapAccessToken(accessToken string) (subjectID, tokenID oid.ID, tokenValue string, err error) {
	ct, err := strs.DecodeBase64(accessToken)
	if err != nil {
		return
	}
	ctbs := strs.ToBytes(&ct)

	ptbs, err := EncryptorConfig.Get().Decrypt(ctbs, nil)
	if err != nil {
		return
	}
	pt := strs.FromBytes(&ptbs)

	if ss := strings.SplitN(pt, ":", 3); len(ss) == 3 {
		return oid.ID(ss[0]), oid.ID(ss[1]), ss[2], nil
	}

	return "", "", pt, nil
}

func wrapAccessToken(subjectID, tokenID oid.ID, tokenValue string) (accessToken string, err error) {
	pt := strs.Join(":", string(subjectID), string(tokenID), tokenValue)
	ptbs := strs.ToBytes(&pt)

	ctbs, err := EncryptorConfig.Get().Encrypt(ptbs, nil)
	if err != nil {
		return
	}
	ct := strs.FromBytes(&ctbs)

	return strs.EncodeBase64(ct), nil
}

type AccessToken struct {
	Raw   *model.Token
	Value string
}

// CreateAccessToken creates a token with the given kind, name and expiration in seconds.
func CreateAccessToken(
	ctx context.Context,
	mc model.ClientSet,
	subjectID oid.ID,
	kind, name string,
	expirationSeconds *int,
) (*AccessToken, error) {
	entity := &model.Token{
		SubjectID: subjectID,
		Kind:      kind,
		Name:      name,
		Value:     crypto.String(strs.String(32)),
	}

	if expirationSeconds != nil {
		e := time.Now().Add(time.Duration(*expirationSeconds) * time.Second)
		entity.Expiration = &e
	}

	creates, err := dao.TokenCreates(mc, entity)
	if err != nil {
		return nil, err
	}

	entity, err = creates[0].Save(ctx)
	if err != nil {
		return nil, err
	}

	at, err := wrapAccessToken(entity.SubjectID, entity.ID, string(entity.Value))
	if err != nil {
		return nil, err
	}

	return &AccessToken{
		Raw:   entity,
		Value: at,
	}, nil
}
