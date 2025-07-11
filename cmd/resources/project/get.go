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

package project

import (
	"encoding/json"
	"strings"

	"github.com/huhouhua/glctl/pkg/cli/genericiooptions"
	"github.com/huhouhua/glctl/pkg/util/templates"

	"github.com/AlekSi/pointer"
	"github.com/spf13/cobra"
	gitlab "gitlab.com/gitlab-org/api/client-go"

	cmdutil "github.com/huhouhua/glctl/cmd/util"
	"github.com/huhouhua/glctl/cmd/validate"
)

type ListOptions struct {
	gitlabClient *gitlab.Client
	ioStreams    genericiooptions.IOStreams
	Visibility   string
	FromGroup    string
	Out          string
	group        *gitlab.ListGroupProjectsOptions
	project      *gitlab.ListProjectsOptions
	ProjectId    *string
	AllGroups    bool
}

var (
	getProjectsExample = templates.Examples(`
# get all projects
glctl get projects

# get all projects from a group
glctl get projects --all-groups=true`)
)

func NewListOptions(ioStreams genericiooptions.IOStreams) *ListOptions {
	return &ListOptions{
		ioStreams:  ioStreams,
		Visibility: string(gitlab.PrivateVisibility),
		group:      &gitlab.ListGroupProjectsOptions{},
		project: &gitlab.ListProjectsOptions{
			OrderBy:                  pointer.ToString("created_at"),
			Sort:                     pointer.ToString("asc"),
			Search:                   pointer.ToString(""),
			Statistics:               pointer.ToBool(false),
			Owned:                    pointer.ToBool(false),
			Archived:                 pointer.ToBool(false),
			Simple:                   pointer.ToBool(false),
			Membership:               pointer.ToBool(false),
			Starred:                  pointer.ToBool(false),
			WithIssuesEnabled:        pointer.ToBool(false),
			WithMergeRequestsEnabled: pointer.ToBool(false),
			ListOptions: gitlab.ListOptions{
				Page:    0,
				PerPage: 0,
			},
		},
		AllGroups: false,
		Out:       "simple",
	}
}
func NewGetProjectsCmd(f cmdutil.Factory, ioStreams genericiooptions.IOStreams) *cobra.Command {
	o := NewListOptions(ioStreams)
	cmd := &cobra.Command{
		Use:                   "projects",
		Aliases:               []string{"p"},
		Short:                 "List projects of the authenticated user or of a group",
		Example:               getProjectsExample,
		DisableFlagsInUseLine: true,
		TraverseChildren:      true,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Validate(cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{"project"},
	}
	o.AddFlags(cmd)
	return cmd
}

// AddFlags registers flags for a cli
func (o *ListOptions) AddFlags(cmd *cobra.Command) {
	cmdutil.AddFromGroupVarPFlag(cmd, &o.FromGroup)
	cmdutil.AddProjectOrderByVarFlag(cmd, o.project.OrderBy)
	cmdutil.AddSortVarFlag(cmd, o.project.Sort)
	cmdutil.AddSearchVarFlag(cmd, o.project.Search)
	cmdutil.AddStatisticsVarFlag(cmd, o.project.Statistics)
	cmdutil.AddVisibilityVarFlag(cmd, &o.Visibility)
	cmdutil.AddOwnedVarFlag(cmd, o.project.Owned)
	cmdutil.AddPaginationVarFlags(cmd, &o.project.ListOptions)
	cmdutil.AddOutFlag(cmd, &o.Out)
	f := cmd.Flags()
	f.BoolVar(o.project.Archived, "archived", *o.project.Archived,
		"Limit by archived status")
	f.BoolVar(o.project.Simple, "simple", *o.project.Simple,
		"Return only the ID, URL, name, and path of each project")
	f.BoolVar(o.project.Membership, "membership", *o.project.Membership,
		"Limit by projects that the current user is a member of")
	f.BoolVar(o.project.Starred, "starred", *o.project.Starred,
		"Limit by projects starred by the current user")
	f.BoolVar(o.project.WithIssuesEnabled, "with-issues-enabled", *o.project.WithIssuesEnabled,
		"Limit by enabled issues feature")
	f.BoolVar(o.project.WithMergeRequestsEnabled, "with-merge-requests-enabled", *o.project.WithMergeRequestsEnabled,
		"Limit by enabled merge requests feature")
	f.BoolVarP(
		&o.AllGroups,
		"all-groups",
		"A",
		o.AllGroups,
		"If present, list the requested object(s) across all groups. group in current context is ignored even if specified with --group.",
	)
}

// Complete completes all the required options.
func (o *ListOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	var err error
	if len(args) > 0 {
		o.ProjectId = pointer.ToString(args[0])
	}
	o.gitlabClient, err = f.GitlabClient()
	if err != nil {
		return err
	}
	o.project.Visibility = gitlab.Ptr(gitlab.VisibilityValue(o.Visibility))
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
	if o.ProjectId != nil {
		project, _, err := o.gitlabClient.Projects.GetProject(*o.ProjectId, &gitlab.GetProjectOptions{})
		if err != nil {
			return err
		}
		return cmdutil.PrintProjectsOut(o.Out, o.ioStreams.Out, project)
	}
	var projects []*gitlab.Project
	var err error

	if strings.TrimSpace(o.FromGroup) != "" {
		projects, _, err = o.gitlabClient.Groups.ListGroupProjects(o.FromGroup, o.group)
	} else {
		if o.AllGroups {
			o.project.ListOptions.PerPage = 100
			o.project.ListOptions.Page = 1
		}
		for {
			var portion []*gitlab.Project
			portion, _, err = o.gitlabClient.Projects.ListProjects(o.project)
			if err != nil {
				return nil
			}
			projects = append(projects, portion...)
			if cap(portion) == 0 || !o.AllGroups {
				break
			}
			o.project.Page++
		}
	}
	if err != nil {
		return nil
	}
	return cmdutil.PrintProjectsOut(o.Out, o.ioStreams.Out, projects...)
}
