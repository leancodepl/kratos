package oidc

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"

	"github.com/ory/herodot"
	"github.com/ory/x/httpx"
)

type ProviderInstagram struct {
	*ProviderGenericOIDC
}

func NewProviderInstagram(
	config *Configuration,
	reg Dependencies,
) Provider {
	config.IssuerURL = "https://www.instagram.com"
	return &ProviderInstagram{
		ProviderGenericOIDC: &ProviderGenericOIDC{
			config: config,
			reg:    reg,
		},
	}
}

func (g *ProviderInstagram) OAuth2(ctx context.Context) (*oauth2.Config, error) {
	endpoint := oauth2.Endpoint{
		AuthURL:  "https://www.instagram.com/oauth/authorize",
		TokenURL: "https://api.instagram.com/oauth/access_token",
	}
	// Instagram is pure OAuth2, not OIDC - do NOT add "openid" scope
	// which would cause "Invalid platform app" error
	return &oauth2.Config{
		ClientID:     g.config.ClientID,
		ClientSecret: g.config.ClientSecret,
		Endpoint:     endpoint,
		Scopes:       g.config.Scope, // Use only configured scopes, no automatic "openid"
		RedirectURL:  g.config.Redir(g.reg.Config().OIDCRedirectURIBase(ctx)),
	}, nil
}

func (g *ProviderInstagram) Claims(ctx context.Context, exchange *oauth2.Token, query url.Values) (*Claims, error) {
	o, err := g.OAuth2(ctx)
	if err != nil {
		return nil, err
	}

	// Instagram Graph API endpoint for user profile
	// Note: Instagram Business Login typically provides comprehensive profile data
	// Email is generally not available through Instagram Business Login
	fields := "id,username,biography,profile_picture_url"
	u, err := url.Parse("https://graph.instagram.com/me?fields=" + fields)
	if err != nil {
		return nil, errors.WithStack(herodot.ErrInternalServerError.WithReasonf("%s", err))
	}

	ctx, client := httpx.SetOAuth2(ctx, g.reg.HTTPClient(ctx), o, exchange)
	req, err := retryablehttp.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, errors.WithStack(herodot.ErrInternalServerError.WithReasonf("%s", err))
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.WithStack(herodot.ErrInternalServerError.WithReasonf("%s", err))
	}
	defer resp.Body.Close()

	if err := logUpstreamError(g.reg.Logger(), resp); err != nil {
		return nil, err
	}

	var user struct {
		Id             string `json:"id,omitempty"`
		Username       string `json:"username,omitempty"`
		Biography      string `json:"biography,omitempty"`
		ProfilePicture string `json:"profile_picture_url,omitempty"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, errors.WithStack(herodot.ErrInternalServerError.WithReasonf("%s", err))
	}

	rawClaims := map[string]interface{}{
		"id":                  user.Id,
		"username":            user.Username,
		"biography":           user.Biography,
		"profile_picture_url": user.ProfilePicture,
	}

	return &Claims{
		Issuer:            u.String(),
		Subject:           user.Id,
		PreferredUsername: user.Username,
		Nickname:          user.Username,
		Picture:           user.ProfilePicture,
		RawClaims:         rawClaims,
	}, nil
}
