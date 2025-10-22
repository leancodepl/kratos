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

func (g *ProviderInstagram) OAuth2(ctx context.Context) *oauth2.Config {
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
	}
}

func (g *ProviderInstagram) Claims(ctx context.Context, exchange *oauth2.Token, query url.Values) (*Claims, error) {
	o := g.OAuth2(ctx)

	// Instagram Graph API endpoint for user profile
	// Note: Instagram Business Login typically provides comprehensive profile data
	// Email is generally not available through Instagram Business Login
	fields := "id,username,account_type,name,biography,followers_count,follows_count," +
		"media_count,profile_picture_url,website"
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
		AccountType    string `json:"account_type,omitempty"`
		Name           string `json:"name,omitempty"`
		Biography      string `json:"biography,omitempty"`
		FollowersCount int64  `json:"followers_count,omitempty"`
		FollowsCount   int64  `json:"follows_count,omitempty"`
		MediaCount     int64  `json:"media_count,omitempty"`
		ProfilePicture string `json:"profile_picture_url,omitempty"`
		Website        string `json:"website,omitempty"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, errors.WithStack(herodot.ErrInternalServerError.WithReasonf("%s", err))
	}

	// Populate RawClaims with all Instagram-specific fields for jsonnet mapper access
	rawClaims := map[string]interface{}{
		"id":                  user.Id,
		"username":            user.Username,
		"account_type":        user.AccountType,
		"name":                user.Name,
		"biography":           user.Biography,
		"followers_count":     user.FollowersCount,
		"follows_count":       user.FollowsCount,
		"media_count":         user.MediaCount,
		"profile_picture_url": user.ProfilePicture,
		"website":             user.Website,
	}

	return &Claims{
		Issuer:            u.String(),
		Subject:           user.Id,
		PreferredUsername: user.Username,
		Nickname:          user.Username,
		Name:              user.Name,
		Picture:           user.ProfilePicture,
		Website:           user.Website,
		RawClaims:         rawClaims,
		// Note: Email is typically not provided by Instagram Business Login
		// All Instagram-specific fields are stored in RawClaims and accessible via jsonnet mapper:
		// account_type, biography, followers_count, follows_count, media_count
	}, nil
}
