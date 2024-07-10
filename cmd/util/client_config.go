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

import "github.com/huhouhua/glctl/cmd/types"

// ClientConfig is used to make it easy to get an api server client.
type ClientConfig interface {
	// ClientConfig returns a complete client config
	ClientConfig() (*types.Config, error)
}

// DirectClientConfig wrap for Config.
type DirectClientConfig struct {
	oathInfo *types.GitLabOauthInfo
	oathEnv  *types.GitLabOathFormEnv
}

// NewClientConfigFromConfig takes your Config and gives you back a ClientConfig.
func NewClientConfigFromConfig(oathInfo *types.GitLabOauthInfo, oathEnv *types.GitLabOathFormEnv) ClientConfig {
	return &DirectClientConfig{
		oathInfo: oathInfo,
		oathEnv:  oathEnv,
	}
}

// NewClientConfigFromBytes takes your glctl config and gives you back a ClientConfig.
func NewClientConfigFromBytes(configBytes []byte) (ClientConfig, error) {
	config, err := Load(configBytes)
	if err != nil {
		return nil, err
	}
	return &DirectClientConfig{
		oathInfo: config.OathInfo,
		oathEnv:  config.OathEnv,
	}, nil
}

func (config *DirectClientConfig) ClientConfig() (*types.Config, error) {
	clientConfig := &types.Config{
		OathInfo: config.oathInfo,
		OathEnv:  config.oathEnv,
	}
	return clientConfig, nil
}
