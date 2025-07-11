// Copyright 2024 The Kevin Berger <huhouhuam@outlook.com> Authors
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
)

type ListOptions struct {
	gitlabClient *gitlab.Client
	Out          string
	branch       *gitlab.ListBranchesOptions
	All          bool
	ioStreams    genericiooptions.IOStreams
}

var (
	getBranchsExample = templates.Examples(`
# get branch from project

# get branch with project name
glctl get branch group1/devops

# get branch with project id
glctl get branch 100

# get all branch
glctl get branch 100 -A
`)
)

func NewListOptions(ioStreams genericiooptions.IOStreams) *ListOptions {
	return &ListOptions{
		ioStreams: ioStreams,
		branch: &gitlab.ListBranchesOptions{
			ListOptions: gitlab.ListOptions{
				Page:    1,
				PerPage: 10,
			},
		},
		All: false,
		Out: "simple",
	}
}
func NewGetBranchesCmd(f cmdutil.Factory, ioStreams genericiooptions.IOStreams) *cobra.Command {
	o := NewListOptions(ioStreams)
	cmd := &cobra.Command{
		Use:                   "branch [Project]",
		Aliases:               []string{"b"},
		Short:                 "List all branches of a repository",
		Example:               getBranchsExample,
		DisableFlagsInUseLine: true,
		TraverseChildren:      true,
		Args:                  require.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Validate(cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{"branch"},
	}
	o.AddFlags(cmd)
	return cmd
}

// AddFlags registers flags for a cli
func (o *ListOptions) AddFlags(cmd *cobra.Command) {
	cmdutil.AddPaginationVarFlags(cmd, &o.branch.ListOptions)
	cmdutil.AddOutFlag(cmd, &o.Out)
	f := cmd.Flags()
	f.BoolVarP(
		&o.All,
		"all",
		"A",
		o.All,
		"If present, list the across all project branch. branch in current context is ignored even if specified with --all.",
	)
}

// Complete completes all the required options.
func (o *ListOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	var err error
	o.gitlabClient, err = f.GitlabClient()
	return err
}

// Validate makes sure there is no discrepency in command options.
func (o *ListOptions) Validate(cmd *cobra.Command, args []string) error {
	if len(args) > 0 && strings.TrimSpace(args[0]) == "" {
		return fmt.Errorf("error from server (NotFound): project %s not found", args[0])
	}
	return nil
}

// Run executes a list subcommand using the specified options.
func (o *ListOptions) Run(args []string) error {
	var branches []*gitlab.Branch
	if o.All {
		o.branch.ListOptions.PerPage = 100
		o.branch.ListOptions.Page = 1
	}
	for {
		list, _, err := o.gitlabClient.Branches.ListBranches(args[0], o.branch)
		if err != nil {
			return nil
		}
		branches = append(branches, list...)
		if cap(list) == 0 || !o.All {
			break
		}
		o.branch.Page++
	}
	return cmdutil.PrintBranchOut(o.Out, o.ioStreams.Out, branches...)
}
