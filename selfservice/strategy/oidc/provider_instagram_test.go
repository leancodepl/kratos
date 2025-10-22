// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oidc_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/kratos/internal"
	"github.com/ory/kratos/selfservice/strategy/oidc"
)

func TestProviderInstagram_OAuth2(t *testing.T) {
	_, reg := internal.NewVeryFastRegistryWithoutDB(t)

	p := oidc.NewProviderInstagram(&oidc.Configuration{
		Provider:     "instagram",
		ID:           "valid",
		ClientID:     "client",
		ClientSecret: "secret",
		Mapper:       "file://./stub/hydra.schema.json",
		Scope:        []string{"instagram_business_basic"},
	}, reg)

	c, err := p.(oidc.OAuth2Provider).OAuth2(context.Background())
	require.NoError(t, err)

	assert.Equal(t, "https://www.instagram.com/oauth/authorize", c.Endpoint.AuthURL)
	assert.Equal(t, "https://api.instagram.com/oauth/access_token", c.Endpoint.TokenURL)
	assert.Equal(t, "client", c.ClientID)
	assert.Equal(t, "secret", c.ClientSecret)
	assert.Contains(t, c.Scopes, "instagram_business_basic")
	// Critical: Instagram is OAuth2, not OIDC - "openid" scope must NOT be present
	assert.NotContains(t, c.Scopes, "openid")
}

func TestProviderInstagram_Config(t *testing.T) {
	_, reg := internal.NewVeryFastRegistryWithoutDB(t)

	p := oidc.NewProviderInstagram(&oidc.Configuration{
		Provider:     "instagram",
		ID:           "instagram-test",
		ClientID:     "test-client",
		ClientSecret: "test-secret",
		Scope:        []string{"instagram_business_basic", "instagram_business_manage_messages"},
	}, reg)

	config := p.Config()
	assert.Equal(t, "instagram-test", config.ID)
	assert.Equal(t, "test-client", config.ClientID)
	assert.Equal(t, "test-secret", config.ClientSecret)
	assert.Equal(t, "https://www.instagram.com", config.IssuerURL)
	assert.Contains(t, config.Scope, "instagram_business_basic")
	assert.Contains(t, config.Scope, "instagram_business_manage_messages")
}
