package types

type Config struct {
	OathInfo *GitLabOauthInfo
	OathEnv  *GitLabOathFormEnv
}

func NewConfig() *Config {
	return &Config{
		OathEnv:  &GitLabOathFormEnv{},
		OathInfo: &GitLabOauthInfo{},
	}
}
