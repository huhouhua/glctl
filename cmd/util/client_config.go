package util

import "github.com/huhouhua/gitlab-repo-operator/cmd/types"

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

// NewClientConfigFromBytes takes your grepo config and gives you back a ClientConfig.
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
