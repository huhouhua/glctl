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

package branch

import (
	"fmt"
	"github.com/AlekSi/pointer"
	"github.com/huhouhua/glctl/cmd/require"
	cmdutil "github.com/huhouhua/glctl/cmd/util"
	"github.com/huhouhua/glctl/cmd/validate"
	"github.com/huhouhua/glctl/util/cli"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
	"strings"
)

type CreateOptions struct {
	gitlabClient *gitlab.Client
	branch       *gitlab.CreateBranchOptions
	project      string
	Out          string
	ioStreams    cli.IOStreams
}

var (
	createBranchDesc = "Create a new branch for a specified project"

	createBranchExample = `# create a develop branch from master branch for project group/myapp
glctl create branch develop --project=group/myapp --ref=master`
)

func NewCreateOptions(ioStreams cli.IOStreams) *CreateOptions {
	return &CreateOptions{
		ioStreams: ioStreams,
		branch: &gitlab.CreateBranchOptions{
			Ref:    pointer.ToString(""),
			Branch: pointer.ToString(""),
		},
		Out: "simple",
	}
}

func NewCreateBranchCmd(f cmdutil.Factory, ioStreams cli.IOStreams) *cobra.Command {
	o := NewCreateOptions(ioStreams)
	cmd := &cobra.Command{
		Use:                   "branch",
		Aliases:               []string{"b"},
		Short:                 createBranchDesc,
		Example:               createBranchExample,
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
	o.AddFlags(cmd)
	return cmd
}

// AddFlags registers flags for a cli
func (o *CreateOptions) AddFlags(cmd *cobra.Command) {
	cmdutil.AddProjectVarPFlag(cmd, &o.project)
	cmdutil.AddOutFlag(cmd, &o.Out)
	validate.VerifyMarkFlagRequired(cmd, "project")
	f := cmd.Flags()
	f.StringVarP(o.branch.Ref, "ref", "r", *o.branch.Ref,
		"The branch name or commit SHA to create branch from")
	validate.VerifyMarkFlagRequired(cmd, "ref")
}

// Complete completes all the required options.
func (o *CreateOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	var err error
	o.gitlabClient, err = f.GitlabClient()
	if len(args) > 0 {
		o.branch.Branch = pointer.ToString(args[0])
	}
	return err
}

// Validate makes sure there is no discrepency in command options.
func (o *CreateOptions) Validate(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("please enter branch")
	}
	if strings.TrimSpace(args[0]) == "" {
		return fmt.Errorf("error from server (NotFound): project %s not found", args[0])
	}
	if strings.TrimSpace(o.project) == "" || strings.TrimSpace(*o.branch.Ref) == "" {
		return cmd.Usage()
	}
	return nil
}

// Run executes a list subcommand using the specified options.
func (o *CreateOptions) Run(args []string) error {
	branch, _, err := o.gitlabClient.Branches.CreateBranch(o.project, o.branch)
	if err != nil {
		return err
	}
	cmdutil.PrintBranchOut(o.Out, branch)
	return nil
}
