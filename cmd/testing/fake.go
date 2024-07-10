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

package testing

import (
	"github.com/AlekSi/pointer"
	"github.com/huhouhua/glctl/cmd/types"
	cmdutil "github.com/huhouhua/glctl/cmd/util"
)

var _ cmdutil.RESTClientGetter = &FakeRESTClientGetter{}

type FakeRESTClientGetter struct {
	cfg       types.Config
	clientCfg cmdutil.ClientConfig
}

func (f FakeRESTClientGetter) ToRESTConfig() (*types.Config, error) {
	return f.ToRawGLConfigLoader().ClientConfig()
}

func (f FakeRESTClientGetter) ToRawGLConfigLoader() cmdutil.ClientConfig {
	return f.clientCfg
}

func NewFakeRESTClientGetter() *FakeRESTClientGetter {
	cfg := types.Config{
		OathEnv: &types.GitLabOathFormEnv{
			Url:          pointer.ToString(""),
			UserName:     pointer.ToString(""),
			Password:     pointer.ToString(""),
			PrivateToken: pointer.ToString(""),
			OauthToken:   pointer.ToString(""),
		},
		OathInfo: &types.GitLabOauthInfo{
			AccessToken:  pointer.ToString("86e2f1f41672758ebcdb1c5ffe17ee463809b6c84aa8ddfa050eee7f0fa4756f"),
			CreatedAt:    pointer.ToFloat64(1.693392912e+09),
			HostUrl:      pointer.ToString("http://172.17.162.204"),
			RefreshToken: pointer.ToString("b1b1f718927007332bcc19b120a8a7e268aa91a5892a411a890682cc3e5692fc"),
			Scope:        pointer.ToString("api"),
			TokenType:    pointer.ToString("Bearer"),
			UserName:     pointer.ToString("v-huhouhua@ruijie.com.cn"),
		},
	}
	return &FakeRESTClientGetter{
		cfg:       cfg,
		clientCfg: cmdutil.NewClientConfigFromConfig(cfg.OathInfo, cfg.OathEnv),
	}
}
