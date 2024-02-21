package types

import (
	"errors"

	"github.com/seal-io/walrus/pkg/dao/types/crypto"
	"github.com/seal-io/walrus/pkg/dao/types/property"
)

const (
	SubscriptionID = "subscription_id"
	TenantID       = "tenant_id"
	ClientID       = "client_id"
	ClientSecret   = "client_secret"
)

type Credential struct {
	SubscriptionID string
	TenantID       string
	ClientID       string
	ClientSecret   string
}

func GetCredential(configData crypto.Properties) (*Credential, error) {
	var (
		cred = &Credential{}
		ok   bool
		err  error
	)

	cred.SubscriptionID, ok, err = property.GetString(configData[SubscriptionID].Value)
	if !ok || err != nil || cred.SubscriptionID == "" {
		return nil, errors.New("subscriptionID is empty")
	}

	cred.TenantID, ok, err = property.GetString(configData[TenantID].Value)
	if !ok || err != nil || cred.TenantID == "" {
		return nil, errors.New("tenantID is empty")
	}

	cred.ClientID, ok, err = property.GetString(configData[ClientID].Value)
	if !ok || err != nil || cred.ClientID == "" {
		return nil, errors.New("clientID is empty")
	}

	cred.ClientSecret, ok, err = property.GetString(configData[ClientSecret].Value)
	if !ok || err != nil || cred.ClientSecret == "" {
		return nil, errors.New("clientSecret is empty")
	}

	return cred, nil
}
