// Copyright 2024 The Kevin Berger <huhouhuam@gmail.com> Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http:www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-cleanhttp"
	gitlab "gitlab.com/gitlab-org/api/client-go"
	"golang.org/x/oauth2"
)

const apiVersionPath = "api/v4/"

// PasswordCredentialsAuthSource implements the AuthSource interface for the OAuth 2.0
// resource owner password credentials flow.
type PasswordCredentialsAuthSource struct {
	username   string
	password   string
	httpClient *http.Client
	gitlab.AuthSource
}

func NewPasswordCredentialsAuthSource(username, password string) *PasswordCredentialsAuthSource {
	return &PasswordCredentialsAuthSource{
		username:   username,
		password:   password,
		httpClient: cleanhttp.DefaultPooledClient(),
	}
}

func (as *PasswordCredentialsAuthSource) Init(ctx context.Context, client *gitlab.Client) error {
	ctx = context.WithValue(ctx, oauth2.HTTPClient, as.httpClient)

	baseURL := strings.TrimSuffix(client.BaseURL().String(), apiVersionPath)
	config := &oauth2.Config{
		Endpoint: oauth2.Endpoint{
			AuthURL:       baseURL + "oauth/authorize",
			TokenURL:      baseURL + "oauth/token",
			DeviceAuthURL: baseURL + "oauth/authorize_device",
		},
	}

	pct, err := config.PasswordCredentialsToken(ctx, as.username, as.password)
	if err != nil {
		return fmt.Errorf("PasswordCredentialsToken(%q, ******): %w", as.username, err)
	}

	as.AuthSource = gitlab.OAuthTokenSource{
		TokenSource: config.TokenSource(ctx, pct),
	}
	return nil
}
