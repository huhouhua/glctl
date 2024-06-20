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

package project

import (
	"encoding/json"
	cmdutil "github.com/huhouhua/gitlab-repo-operator/cmd/util"
	"github.com/huhouhua/gitlab-repo-operator/cmd/validate"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

var getProjectsDesc = "List projects of the authenticated user or of a group"

var getProjectsExample = `# get all projects
grepo get projects

# get all projects from a group
grepo get projects --from-group=Group1`

type ListOptions struct {
	Page         int
	PerPage      int
	FromGroup    string
	Out          string
	gitlabClient *gitlab.Client
	group        *gitlab.ListGroupProjectsOptions
	project      *gitlab.ListProjectsOptions
}

func NewListOptions() *ListOptions {
	return &ListOptions{
		Page:    0,
		PerPage: 0,
	}
}

func NewGetProjectsCmd(f cmdutil.Factory) *cobra.Command {
	o := NewListOptions()
	cmd := &cobra.Command{
		Use:                   "projects",
		Aliases:               []string{},
		Short:                 getProjectsDesc,
		Example:               getProjectsExample,
		DisableFlagsInUseLine: true,
		TraverseChildren:      true,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Validate(cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
	}
	flags := cmd.Flags()
	flags.IntVarP(&o.Page, "page", "p", o.Page, "Page of results to retrieve")
	flags.IntVarP(&o.PerPage, "per-page", "", o.PerPage, "The number of results to include per page")

	cmdutil.AddOutFlag(cmd)
	return cmd
}

// Complete completes all the required options.
func (o *ListOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	var err error
	o.gitlabClient, err = f.GitlabClient()
	if err != nil {
		return err
	}
	o.Out = cmdutil.GetFlagString(cmd, "out")
	o.project = cmdutil.AssignListProjectOptions(cmd)
	opt, err := json.Marshal(o.project)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(opt, o.group); err != nil {
		return err
	}
	return nil
}

// Validate makes sure there is no discrepency in command options.
func (o *ListOptions) Validate(cmd *cobra.Command, args []string) error {
	if err := validate.ValidateSortFlagValue(cmd); err != nil {
		return err
	}
	if err := validate.ValidateProjectOrderByFlagValue(cmd); err != nil {
		return err
	}
	return validate.ValidateVisibilityFlagValue(cmd)
}

// Run executes a list subcommand using the specified options.
func (o *ListOptions) Run(args []string) error {
	var projects []*gitlab.Project
	var err error
	if o.FromGroup != "" {
		projects, _, err = o.gitlabClient.Groups.ListGroupProjects(o.FromGroup, o.group)
	} else {
		projects, _, err = o.gitlabClient.Projects.ListProjects(o.project)
	}
	if err != nil {
		return nil
	}
	cmdutil.PrintProjectsOut(o.Out, projects...)
	return nil
}
