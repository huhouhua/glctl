// Copyright 2024 The huhouhua Authors
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

package cmd

import (
	"github.com/huhouhua/gitlab-repo-operator/cmd/types"
	"github.com/spf13/viper"
	"strings"
)

func withReadOathInfoConfig() (*types.GitLabOauthInfo, error) {
	var oathUserInfo types.GitLabOauthInfo
	err := viper.Unmarshal(&oathUserInfo)
	if err != nil {
		return nil, err
	}
	return &oathUserInfo, nil
}

func withReadOathEnvConfig() (*types.GitLabOathFormEnv, error) {
	var oathEnvInfo types.GitLabOathFormEnv
	viper.SetEnvPrefix("GITLAB")
	err := viper.Unmarshal(&oathEnvInfo)
	if err != nil {
		return nil, err
	}
	return &oathEnvInfo, nil
}

type GitLabAuthorization struct {
	oathInfo *types.GitLabOauthInfo
	oathEnv  *types.GitLabOathFormEnv
}

func newGitLabAuthorization(info *types.GitLabOauthInfo, env *types.GitLabOathFormEnv) *GitLabAuthorization {
	return &GitLabAuthorization{
		oathInfo: info,
		oathEnv:  env,
	}
}

func (g *GitLabAuthorization) HasAuth() bool {
	info := g.oathInfo
	if info == nil {
		return false
	}
	if strings.TrimSpace(info.AccessToken) == "" {
		return false
	}
	if strings.TrimSpace(info.HostUrl) == "" {
		return false
	}
	return true
}

func (g *GitLabAuthorization) HasPasswordAuth() bool {
	env := g.oathEnv
	if env == nil {
		return false
	}
	if strings.TrimSpace(env.URL) == "" {
		return false
	}
	if strings.TrimSpace(env.UserName) == "" {
		return false
	}
	if strings.TrimSpace(env.Password) == "" {
		return false
	}
	return true
}
func (g *GitLabAuthorization) HasBasicAuth() bool {
	env := g.oathEnv
	if env == nil {
		return false
	}
	if strings.TrimSpace(env.URL) == "" {
		return false
	}
	if strings.TrimSpace(env.PrivateToken) == "" {
		return false
	}
	return true
}
func (g *GitLabAuthorization) HasOathAuth() bool {
	env := g.oathEnv
	if env == nil {
		return false
	}
	if strings.TrimSpace(env.URL) == "" {
		return false
	}
	if strings.TrimSpace(env.OauthToken) == "" {
		return false
	}
	return true
}
