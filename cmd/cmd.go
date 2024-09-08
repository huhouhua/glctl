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
	"github.com/huhouhua/glctl/pkg/cli/genericiooptions"
	"github.com/huhouhua/glctl/pkg/util/progress"
	templates2 "github.com/huhouhua/glctl/pkg/util/templates"
	"io"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/huhouhua/glctl/cmd/completion"
	"github.com/huhouhua/glctl/cmd/create"
	delete "github.com/huhouhua/glctl/cmd/delete"
	"github.com/huhouhua/glctl/cmd/edit"
	"github.com/huhouhua/glctl/cmd/get"
	"github.com/huhouhua/glctl/cmd/login"
	"github.com/huhouhua/glctl/cmd/logout"
	"github.com/huhouhua/glctl/cmd/replace"
	cmdutil "github.com/huhouhua/glctl/cmd/util"
	"github.com/huhouhua/glctl/cmd/version"
)

var AuthDoc = `
There are two options to authenticate the command-line client to Gitlab interface:

1.) Using the 'login' command by passing the host url, username and password.

$ glctl login

The login token will be saved in $HOME/.glctl.yaml file.

2.) Using Environment variables.

* Basic Authentication (if using a username and password)
    - GITLAB_USERNAME
    - GITLAB_PASSWORD
    - GITLAB_URL

* Private Token (if using a private token)
    - GITLAB_PRIVATE_TOKEN
    - GITLAB_URL

* OAuth Token (if using an oauth token)
    - GITLAB_OAUTH_TOKEN
    - GITLAB_URL`

var cfgFile string
var globalUsage = `The gitlab repository operator for the command-line.

This client helps you view, update, create, and delete Gitlab resources from the
command-line interface.
`

// NeDefaultGlCtlCommand creates the `glctl` command with default arguments.
func NeDefaultGlCtlCommand() *cobra.Command {
	return NeGlCtlCommand(os.Stdin, os.Stdout, os.Stderr)
}

func NeGlCtlCommand(in io.Reader, out, err io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "glctl",
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
	flags.StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.glctl.yaml)")
	flags.SetNormalizeFunc(cmdutil.WarnWordSepNormalizeFunc) // Warn for "_" flags

	// Normalize all flags that are coming from other packages or pre-configurations
	// a.k.a. change all "_" to "-". e.g. glog package
	flags.SetNormalizeFunc(cmdutil.WordSepNormalizeFunc)

	addProfilingFlags(flags)

	if noColor, ok := os.LookupEnv("NO_COLOR"); ok && noColor != "" {
		progress.NoColor()
	}

	cobra.OnInitialize(initConfig)
	flags.AddGoFlagSet(flag.CommandLine)

	configFlags := cmdutil.NewConfigFlags(false)
	configFlags.AddFlags(flags)
	f := cmdutil.NewFactory(configFlags)
	// From this point and forward we get warnings on flags that contain "_" separators
	cmd.SetGlobalNormalizationFunc(cmdutil.WarnWordSepNormalizeFunc)
	ioStreams := genericiooptions.IOStreams{In: in, Out: out, ErrOut: err}
	groups := templates2.CommandGroups{
		{
			Message: "Basic Commands:",
			Commands: []*cobra.Command{
				get.NewGetCmd(f, ioStreams),
				edit.NewEditCmd(f, ioStreams),
				delete.NewDeleteCmd(f, ioStreams),
				create.NewCreateCmd(f, ioStreams),
			},
		},
		{
			Message: "Authorization Commands:",
			Commands: []*cobra.Command{
				login.NewLoginCmd(ioStreams),
				logout.NewLogoutCmd(ioStreams),
			},
		},
		{
			Message: "Advanced Commands:",
			Commands: []*cobra.Command{
				replace.NewReplaceCmd(f, ioStreams),
			},
		},
		{
			Message: "Settings Commands:",
			Commands: []*cobra.Command{
				completion.NewCmdCompletion(ioStreams, ""),
			},
		},
	}
	groups.Add(cmd)

	filters := []string{"options"}
	templates2.ActsAsRootCommand(cmd, filters, groups...)
	cmd.AddCommand(version.NewCmdVersion(f, ioStreams))
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
			cmdutil.Error(os.Stdout, err)
		}

		// Search config in home directory with name ".gl" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".glctl")
	}
	viper.AutomaticEnv() // read in environment variables that match
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		// NOTE: the config file is not required to exists
		// raise an error if error is other than config file not found
		if !strings.Contains(err.Error(), `Config File ".glctl" Not Found`) {
			cmdutil.Error(os.Stdout, err)
		}
	}
}
func runHelp(cmd *cobra.Command, args []string) {
	_ = cmd.Help()
}
