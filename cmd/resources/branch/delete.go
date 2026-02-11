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

package branch

import (
	"fmt"
	"strings"

	"github.com/huhouhua/glctl/pkg/cli/genericiooptions"
	"github.com/huhouhua/glctl/pkg/util/templates"

	"github.com/spf13/cobra"
	gitlab "gitlab.com/gitlab-org/api/client-go"

	"github.com/huhouhua/glctl/cmd/require"
	cmdutil "github.com/huhouhua/glctl/cmd/util"
	"github.com/huhouhua/glctl/cmd/validate"
)

type DeleteOptions struct {
	gitlabClient *gitlab.Client
	project      string
	branch       string
	ioStreams    genericiooptions.IOStreams
}

var (
	deleteBranchExample = templates.Examples(`
# delete a develop branch from project group/myapp
glctl delete branch develop --project=group/myapp`)
)

func NewDeleteOptions(ioStreams genericiooptions.IOStreams) *DeleteOptions {
	return &DeleteOptions{
		ioStreams: ioStreams,
	}
}

func NewDeleteBranchCmd(f cmdutil.Factory, ioStreams genericiooptions.IOStreams) *cobra.Command {
	o := NewDeleteOptions(ioStreams)
	cmd := &cobra.Command{
		Use:                   "branch",
		Aliases:               []string{"b"},
		Short:                 "Delete a project branch",
		Example:               deleteBranchExample,
		Args:                  require.ExactArgs(1),
		DisableFlagsInUseLine: true,
		TraverseChildren:      true,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Validate(cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
	}
	cmdutil.AddProjectVarPFlag(cmd, &o.project)
	validate.VerifyMarkFlagRequired(cmd, "project")
	return cmd
}

// Complete completes all the required options.
func (o *DeleteOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	var err error
	o.gitlabClient, err = f.GitlabClient()
	if len(args) > 0 {
		o.branch = args[0]
	}
	return err
}

// Validate makes sure there is no discrepency in command options.
func (o *DeleteOptions) Validate(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("please enter branch")
	}
	if strings.TrimSpace(args[0]) == "" {
		return fmt.Errorf("error from server (NotFound): project %s not found", args[0])
	}
	if strings.TrimSpace(o.project) == "" {
		return cmd.Usage()
	}
	return nil
}

// Run executes a list subcommand using the specified options.
func (o *DeleteOptions) Run(args []string) error {
	_, err := o.gitlabClient.Branches.DeleteBranch(o.project, o.branch)
	if err != nil {
		return err
	}
	_, _ = fmt.Fprintf(o.ioStreams.Out, "Branch (%s) from project (%s) has been deleted\n", o.branch, o.project)
	return nil
}
