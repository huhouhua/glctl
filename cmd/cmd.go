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
	"flag"
	"fmt"
	"github.com/huhouhua/gl/cmd/completion"
	delete "github.com/huhouhua/gl/cmd/delete"
	"github.com/huhouhua/gl/cmd/edit"
	"github.com/huhouhua/gl/cmd/get"
	"github.com/huhouhua/gl/cmd/login"
	cmdutil "github.com/huhouhua/gl/cmd/util"
	"github.com/huhouhua/gl/cmd/version"
	"github.com/huhouhua/gl/util/templates"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"os"
	"strings"
)

var AuthDoc = `
There are two options to authenticate the command-line client to Gitlab interface:

1.) Using the 'login' command by passing the host url, username and password.

$ gl login

The login token will be saved in $HOME/.gl.yaml file.

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

// NeDefaultGlCommand creates the `iamctl` command with default arguments.
func NeDefaultGlCommand() *cobra.Command {
	return NeGlCommand(os.Stdin, os.Stdout, os.Stderr)
}

func NeGlCommand(in io.Reader, out, err io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "gl",
		Short:         "the gitlab repository operator",
		Long:          fmt.Sprintf("%s\n%s", globalUsage, AuthDoc),
		SilenceErrors: true,
		Run:           runHelp,
		// Hook before and after Run initialize and write profiles to disk,
		// respectively.
		PersistentPreRunE: func(*cobra.Command, []string) error {
			return initProfiling()
		},
		PersistentPostRunE: func(*cobra.Command, []string) error {
			return flushProfiling()
		},
	}
	flags := cmd.PersistentFlags()
	flags.StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.gl.yaml)")
	flags.SetNormalizeFunc(cmdutil.WarnWordSepNormalizeFunc) // Warn for "_" flags

	// Normalize all flags that are coming from other packages or pre-configurations
	// a.k.a. change all "_" to "-". e.g. glog package
	flags.SetNormalizeFunc(cmdutil.WordSepNormalizeFunc)

	addProfilingFlags(flags)

	cobra.OnInitialize(initConfig)
	flags.AddGoFlagSet(flag.CommandLine)

	configFlags := cmdutil.NewConfigFlags(false)
	configFlags.AddFlags(flags)
	f := cmdutil.NewFactory(configFlags)
	// From this point and forward we get warnings on flags that contain "_" separators
	cmd.SetGlobalNormalizationFunc(cmdutil.WarnWordSepNormalizeFunc)
	ioStreams := cmdutil.IOStreams{In: in, Out: out, ErrOut: err}
	groups := templates.CommandGroups{
		{
			Message: "Basic Commands:",
			Commands: []*cobra.Command{
				get.NewGetCmd(f, ioStreams),
				edit.NewEditCmd(f, ioStreams),
				delete.NewDeleteCmd(f, ioStreams),
			},
		},
		{
			Message: "Authorization Commands:",
			Commands: []*cobra.Command{
				login.NewLoginCmd(ioStreams),
			},
		},
		{
			Message: "Settings Commands:",
			Commands: []*cobra.Command{
				completion.NewCmdCompletion(ioStreams),
			},
		},
	}
	groups.Add(cmd)

	filters := []string{"options"}
	templates.ActsAsRootCommand(cmd, filters, groups...)
	cmd.AddCommand(version.NewCmdVersion(f))
	return cmd
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

		// Search config in home directory with name ".gl" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".gl")
	}

	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvPrefix("GITLAB")
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		// NOTE: the config file is not required to exists
		// raise an error if error is other than config file not found
		if !strings.Contains(err.Error(), `Config File ".gl" Not Found`) {
			cmdutil.Error(err)
		}
	}

}
func runHelp(cmd *cobra.Command, args []string) {
	_ = cmd.Help()
}
