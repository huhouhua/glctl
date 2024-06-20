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

package util

import (
	"errors"
	"github.com/AlekSi/pointer"
	"github.com/huhouhua/gitlab-repo-operator/cmd/types"
	"github.com/spf13/pflag"
	"sync"
)

type RESTClientGetter interface {
	// ToRESTConfig returns restconfig
	ToRESTConfig() (*types.Config, error)

	// ToRawGrepoConfigLoader return grepoconfig loader as-is
	ToRawGrepoConfigLoader() ClientConfig
}

var _ RESTClientGetter = &ConfigFlags{}

type ConfigFlags struct {
	Env          *types.GitLabOathFormEnv
	Oath         *types.GitLabOauthInfo
	clientConfig ClientConfig
	lock         sync.Mutex
	// If set to true, will use persistent client config and
	// propagate the config to the places that need it, rather than
	// loading the config multiple times
	usePersistentConfig bool
}

func (f *ConfigFlags) ToRESTConfig() (*types.Config, error) {
	return f.ToRawGrepoConfigLoader().ClientConfig()
}

// ToRawGrepoConfigLoader binds config flag values to config overrides
// Returns an interactive clientConfig if the password flag is enabled,
// or a non-interactive clientConfig otherwise.
func (f *ConfigFlags) ToRawGrepoConfigLoader() ClientConfig {
	if f.usePersistentConfig {
		return f.toRawGrepoPersistentConfigLoader()
	}

	return f.toRawGrepoConfigLoader()
}

// toRawGrepoPersistentConfigLoader binds config flag values to config overrides
// Returns a persistent clientConfig for propagation.
func (f *ConfigFlags) toRawGrepoPersistentConfigLoader() ClientConfig {
	f.lock.Lock()
	defer f.lock.Unlock()

	if f.clientConfig == nil {
		f.clientConfig = f.toRawGrepoConfigLoader()
	}

	return f.clientConfig
}

func (f *ConfigFlags) toRawGrepoConfigLoader() ClientConfig {
	oathInfoCfg, infoErr := LoadOathWithInfoConfig()
	oathInfoEnvCfg, envErr := LoadOathWithEnvConfig()
	if infoErr == nil && envErr == nil {
		panic(errors.Join(infoErr, envErr))
	}
	return NewClientConfigFromConfig(oathInfoCfg, oathInfoEnvCfg)
}

// NewConfigFlags returns ConfigFlags with default values set.
func NewConfigFlags(usePersistentConfig bool) *ConfigFlags {
	return &ConfigFlags{
		Env: &types.GitLabOathFormEnv{
			URL:          pointer.ToString(""),
			UserName:     pointer.ToString(""),
			Password:     pointer.ToString(""),
			PrivateToken: pointer.ToString(""),
			OauthToken:   pointer.ToString(""),
		},
		Oath: &types.GitLabOauthInfo{
			AccessToken:  pointer.ToString(""),
			CreatedAt:    pointer.ToString(""),
			HostUrl:      pointer.ToString(""),
			RefreshToken: pointer.ToString(""),
			Scope:        pointer.ToString(""),
			TokenType:    pointer.ToString(""),
			UserName:     pointer.ToString(""),
		},
		usePersistentConfig: usePersistentConfig,
	}
}

// AddConfig binds client configuration flags to a given flagset.
func (f *ConfigFlags) AddFlags(flags *pflag.FlagSet) {

}
