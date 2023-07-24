package auths

import (
	"net/http"
	"strings"
	"time"

	"github.com/seal-io/seal/pkg/dao/types/object"
	"github.com/seal-io/seal/utils/strs"
)

const SessionCookieName = "seal_session"

func decodeSession(r *http.Request) (subjectID object.ID, domain, value string) {
	v := getAccessSession(r)
	if v == "" {
		return
	}

	subjectID, domain, value, _ = unwrapAccessSession(v)

	return
}

func flushSession(r *http.Request, rw http.ResponseWriter, elements ...string) error {
	o, err := r.Cookie(SessionCookieName)
	if err != nil {
		if len(elements) != 3 {
			return nil
		}
	}

	var v string

	switch {
	case len(elements) == 3:
		v, err = wrapAccessSession(object.ID(elements[0]), elements[1], elements[2])
		if err != nil {
			return err
		}
	case o != nil:
		v = o.Value
	}

	n := &http.Cookie{
		Name:     SessionCookieName,
		Value:    v,
		Path:     "/",
		Domain:   "",
		Secure:   SecureConfig.Get(),
		HttpOnly: true,
		MaxAge:   int(MaxIdleDurationConfig.Get().Round(time.Second) / time.Second),
	}
	if n.MaxAge > 0 {
		n.Expires = time.Now().Add(time.Duration(n.MaxAge) * time.Second)
	}

	http.SetCookie(rw, n)

	return nil
}

func revertSession(r *http.Request, rw http.ResponseWriter) {
	_, err := r.Cookie(SessionCookieName)
	if err != nil {
		return
	}

	n := &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		Domain:   "",
		Secure:   SecureConfig.Get(),
		HttpOnly: true,
		MaxAge:   -1,
		Expires:  time.Now().Add(-time.Second),
	}

	http.SetCookie(rw, n)
}

func getAccessSession(r *http.Request) string {
	c, err := r.Cookie(SessionCookieName)
	if err != nil {
		return ""
	}

	return c.Value
}

func unwrapAccessSession(accessSession string) (subjectID object.ID, domain, value string, err error) {
	ct, err := strs.DecodeBase64(accessSession)
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
		return object.ID(ss[0]), ss[1], ss[2], nil
	}

	return "", "", pt, nil
}

func wrapAccessSession(subjectID object.ID, domain, value string) (accessSession string, err error) {
	pt := strs.Join(":", string(subjectID), domain, value)
	ptbs := strs.ToBytes(&pt)

	ctbs, err := EncryptorConfig.Get().Encrypt(ptbs, nil)
	if err != nil {
		return
	}
	ct := strs.FromBytes(&ctbs)

	return strs.EncodeBase64(ct), nil
}
