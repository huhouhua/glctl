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

package types

type GitLabOauthInfo struct {
	AccessToken  *string  `json:"access_token"  yaml:"access_token"  mapstructure:"access_token"`
	CreatedAt    *float64 `json:"created_at"    yaml:"created_at"    mapstructure:"created_at"`
	HostUrl      *string  `json:"host_url"      yaml:"host_url"      mapstructure:"host_url"`
	RefreshToken *string  `json:"refresh_token" yaml:"refresh_token" mapstructure:"refresh_token"`
	Scope        *string  `json:"scope"         yaml:"scope"         mapstructure:"scope"`
	TokenType    *string  `json:"token_type"    yaml:"token_type"    mapstructure:"token_type"`
	UserName     *string  `json:"user_name"     yaml:"user_name"     mapstructure:"user_name"`
}
