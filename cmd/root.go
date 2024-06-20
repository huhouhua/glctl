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

package cmd

import (
	"fmt"
	"github.com/huhouhua/gitlab-repo-operator/cmd/get"
	"github.com/huhouhua/gitlab-repo-operator/cmd/login"
	cmdutil "github.com/huhouhua/gitlab-repo-operator/cmd/util"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"strings"
)

var AuthDoc = `
There are two options to authenticate the command-line client to Gitlab interface:

1.) Using the 'login' command by passing the host url, username and password.

$ grepo login

The login token will be saved in $HOME/.grepo.yaml file.

2.) Using Environment variables.

* Basic Authentication (if using a username and password)
    - GITLAB_USERNAME
    - GITLAB_PASSWORD
    - GITLAB_HTTP_URL

* Private Token (if using a private token)
    - GITLAB_PRIVATE_TOKEN
    - GITLAB_API_HTTP_URL

* OAuth Token (if using an oauth token)
    - GITLAB_OAUTH_TOKEN
    - GITLAB_API_HTTP_URL`

var cfgFile string
var globalUsage = `The gitlab repository operator for the command-line.

This client helps you view, update, create, and delete Gitlab resources from the
command-line interface.
`

func NewRootCmd(out io.Writer) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:           "grepo",
		Short:         "the gitlab repository operator",
		Long:          fmt.Sprintf("%s\n%s", globalUsage, AuthDoc),
		SilenceErrors: true,
		Run:           runHelp,
	}
	flags := cmd.PersistentFlags()
	flags.StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.grepo.yaml)")

	configFlags := cmdutil.NewConfigFlags(false)
	configFlags.AddFlags(flags)
	f := cmdutil.NewFactory(configFlags)

	cmd.AddCommand(
		login.NewLoginCmd(),
		get.NewGetCmd(f))
	return cmd, nil
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			cmdutil.Error(err)
		}

		// Search config in home directory with name ".grepo" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".grepo")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		// NOTE: the config file is not required to exists
		// raise an error if error is other than config file not found
		if !strings.Contains(err.Error(), `Config File ".grepo" Not Found`) {
			cmdutil.Error(err)
		}
	}
}
func runHelp(cmd *cobra.Command, args []string) {
	_ = cmd.Help()
}
