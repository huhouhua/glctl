// Copyright 2024 The Kevin Berger <huhouhuam@outlook.com> Authors
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

package util

import (
	"fmt"
	"strings"

	"github.com/xanzy/go-gitlab"

	"github.com/huhouhua/glctl/cmd/types"
)

func NewForConfig(config *types.Config) (*gitlab.Client, error) {
	authorization := newGitLabAuthorization(config.OathInfo, config.OathEnv)
	switch {
	case authorization.HasAuth():
		return gitlab.NewOAuthClient(
			*authorization.OathInfo.AccessToken,
			gitlab.WithBaseURL(withApiUrl(*authorization.OathInfo.HostUrl)),
		)
	case authorization.HasPasswordAuth():
		return gitlab.NewBasicAuthClient(
			*authorization.OathEnv.UserName,
			*authorization.OathEnv.Password,
			gitlab.WithBaseURL(withApiUrl(*authorization.OathEnv.Url)),
		)
	case authorization.HasBasicAuth():
		return gitlab.NewClient(*authorization.OathEnv.PrivateToken, gitlab.WithBaseURL(*authorization.OathEnv.Url))
	case authorization.HasOathAuth():
		return gitlab.NewOAuthClient(*authorization.OathEnv.OauthToken, gitlab.WithBaseURL(*authorization.OathEnv.Url))
	default:
		return nil, fmt.Errorf("no client was created. "+
			"gitlab configuration was not set properly. \n %s", "")
	}
}

func withApiUrl(url string) string {
	if strings.HasSuffix(url, "/api") {
		return fmt.Sprintf("%s/v4", url)
	}
	return url
}
