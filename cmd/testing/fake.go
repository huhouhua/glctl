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
	"github.com/spf13/viper"

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
	viper.AutomaticEnv()
	cfg, err := cmdutil.NewConfigFlags(false).ToRESTConfig()
	if err != nil {
		panic(err)
	}
	return &FakeRESTClientGetter{
		cfg:       *cfg,
		clientCfg: cmdutil.NewClientConfigFromConfig(cfg.OathInfo, cfg.OathEnv),
	}
}
