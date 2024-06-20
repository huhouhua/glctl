package util

import (
	"github.com/huhouhua/gitlab-repo-operator/cmd/types"
	"github.com/xanzy/go-gitlab"
)

type factoryImpl struct {
	clientGetter RESTClientGetter
}

func NewFactory(clientGetter RESTClientGetter) Factory {
	if clientGetter == nil {
		panic("attempt to instantiate client_access_factory with nil clientGetter")
	}

	f := &factoryImpl{
		clientGetter: clientGetter,
	}
	return f
}

func (f *factoryImpl) ToRESTConfig() (*types.Config, error) {
	return f.clientGetter.ToRESTConfig()
}

func (f *factoryImpl) ToRawGrepoConfigLoader() ClientConfig {
	return f.clientGetter.ToRawGrepoConfigLoader()
}

func (f *factoryImpl) GitlabClient() (*gitlab.Client, error) {
	clientConfig, err := f.ToRESTConfig()
	if err != nil {
		return nil, err
	}
	return NewForConfig(clientConfig)
}
