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

package logout

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"

	cmdutil "github.com/huhouhua/glctl/cmd/util"
	"github.com/huhouhua/glctl/util/cli"
	"github.com/huhouhua/glctl/util/templates"
)

var logoutLong = templates.LongDesc(`
login command. 

Log out from gitlab "delete the $HOME/.glctl.yaml.`)

type Options struct {
	ioStreams cli.IOStreams
	path      string
	file      os.FileInfo
}

func NewOptions(ioStreams cli.IOStreams) *Options {
	return &Options{
		ioStreams: ioStreams,
	}
}
func NewLogoutCmd(ioStreams cli.IOStreams) *cobra.Command {
	o := NewOptions(ioStreams)
	cmd := &cobra.Command{
		Use:                   "logout",
		Short:                 "logout to gitlab",
		Long:                  logoutLong,
		DisableFlagsInUseLine: true,
		Example:               `glctl logout`,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(cmd, args))
			cmdutil.CheckErr(o.Validate(cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
	}
	return cmd
}

// Complete completes all the required options.
func (o *Options) Complete(cmd *cobra.Command, args []string) error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}
	o.path = fmt.Sprintf("%s/.glctl.yaml", home)
	return nil
}

// Validate makes sure there is no discrepency in command options.
func (o *Options) Validate(cmd *cobra.Command, args []string) error {
	fi, err := os.Stat(o.path)
	if os.IsNotExist(err) {
		return fmt.Errorf("%s not exist", o.path)
	}
	if fi.IsDir() {
		return fmt.Errorf("%s cannot be directory", o.path)
	}
	o.file = fi
	return nil
}

// Run executes a create subcommand using the specified options.
func (o *Options) Run(args []string) error {
	err := os.Remove(o.path)
	if err != nil {
		return err
	}
	_, _ = fmt.Fprintf(o.ioStreams.Out, "%s file has been delete by logout command \n", o.path)
	_, _ = fmt.Fprintf(o.ioStreams.Out, "\nlogout Succeeded \n")
	return nil
}
