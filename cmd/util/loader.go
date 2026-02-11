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

package util

import (
	"strings"

	"github.com/AlekSi/pointer"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/huhouhua/glctl/cmd/types"
)

// Load takes a byte slice and deserializes the contents into Config object.
// Encapsulates deserialization without assuming the source is a file.
func Load(data []byte) (*types.Config, error) {
	config := types.NewConfig()
	// if there's no data in a file, return the default object instead of failing (DecodeInto reject empty input)
	if len(data) == 0 {
		return config, nil
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}

func LoadOathWithInfoConfig() (*types.GitLabOauthInfo, error) {
	var oathUserInfo types.GitLabOauthInfo
	err := viper.Unmarshal(&oathUserInfo)
	if err != nil {
		return nil, err
	}
	return &oathUserInfo, nil
}

func LoadOathWithEnvConfig() (*types.GitLabOathFormEnv, error) {
	return &types.GitLabOathFormEnv{
		Url:          pointer.ToString(viper.GetString("GITLAB_URL")),
		UserName:     pointer.ToString(viper.GetString("GITLAB_USERNAME")),
		Password:     pointer.ToString(viper.GetString("GITLAB_PASSWORD")),
		PrivateToken: pointer.ToString(viper.GetString("GITLAB_PRIVATE_TOKEN")),
		OauthToken:   pointer.ToString(viper.GetString("GITLAB_OAUTH_TOKEN")),
	}, nil
}

type GitLabAuthorization struct {
	*types.Config
}

func newGitLabAuthorization(info *types.GitLabOauthInfo, env *types.GitLabOathFormEnv) *GitLabAuthorization {
	return &GitLabAuthorization{
		Config: &types.Config{
			OathInfo: info,
			OathEnv:  env,
		},
	}
}

func (g *GitLabAuthorization) HasAuth() bool {
	info := g.OathInfo
	if info == nil {
		return false
	}
	if info.AccessToken == nil || strings.TrimSpace(*info.AccessToken) == "" {
		return false
	}
	if info.HostUrl == nil || strings.TrimSpace(*info.HostUrl) == "" {
		return false
	}
	return true
}

func (g *GitLabAuthorization) HasPasswordAuth() bool {
	env := g.OathEnv
	if env == nil {
		return false
	}
	if env.Url == nil || strings.TrimSpace(*env.Url) == "" {
		return false
	}
	if env.UserName == nil || strings.TrimSpace(*env.UserName) == "" {
		return false
	}
	if env.Password == nil || strings.TrimSpace(*env.Password) == "" {
		return false
	}
	return true
}
func (g *GitLabAuthorization) HasBasicAuth() bool {
	env := g.OathEnv
	if env == nil {
		return false
	}
	if env.Url == nil || strings.TrimSpace(*env.Url) == "" {
		return false
	}
	if env.PrivateToken == nil || strings.TrimSpace(*env.PrivateToken) == "" {
		return false
	}
	return true
}
func (g *GitLabAuthorization) HasOathAuth() bool {
	env := g.OathEnv
	if env == nil {
		return false
	}
	if env.Url == nil || strings.TrimSpace(*env.Url) == "" {
		return false
	}
	if env.OauthToken == nil || strings.TrimSpace(*env.OauthToken) == "" {
		return false
	}
	return true
}
