package util

import "github.com/xanzy/go-gitlab"

type Factory interface {

	// GitlabClient gives you back an external gitlabClient
	GitlabClient() (*gitlab.Client, error)
}
