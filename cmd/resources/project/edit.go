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
	"github.com/AlekSi/pointer"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"

	"github.com/huhouhua/glctl/cmd/require"
	cmdutil "github.com/huhouhua/glctl/cmd/util"
	"github.com/huhouhua/glctl/cmd/validate"
	"github.com/huhouhua/glctl/util/cli"
	"github.com/huhouhua/glctl/util/templates"
)

type EditOptions struct {
	gitlabClient *gitlab.Client
	project      *gitlab.EditProjectOptions
	Out          string
	ioStreams    cli.IOStreams
}

var (
	editProjectExample = templates.Examples(`
# update a project by path
glctl edit project ProjectX --desc="A go project"

glctl edit project GroupX/ProjectX --merge-method=rebase_merge 

# update a project with id (23)
glctl edit project 3 --desc="A go project"`)
)

func NewEditOptions(ioStreams cli.IOStreams) *EditOptions {
	return &EditOptions{
		ioStreams: ioStreams,
		project: &gitlab.EditProjectOptions{
			Name:                 pointer.ToString(""),
			Path:                 pointer.ToString(""),
			DefaultBranch:        pointer.ToString("main"),
			Description:          pointer.ToString(""),
			LFSEnabled:           pointer.ToBool(false),
			RequestAccessEnabled: pointer.ToBool(false),
			Visibility:           pointer.To(gitlab.PrivateVisibility),
		},
		Out: "simple",
	}
}

func NewEditProjectCmd(f cmdutil.Factory, ioStreams cli.IOStreams) *cobra.Command {
	o := NewCreateOptions(ioStreams)
	cmd := &cobra.Command{
		Use:                   "project",
		Aliases:               []string{"p"},
		Short:                 "Edit a project by specifying the project id or path and using flags for fields to modify",
		Example:               editProjectExample,
		Args:                  require.ExactArgs(1),
		DisableFlagsInUseLine: true,
		TraverseChildren:      true,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Validate(cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{"projects"},
	}
	o.AddFlags(cmd)
	return cmd
}

// AddFlags registers flags for a cli
func (o *EditOptions) AddFlags(cmd *cobra.Command) {
	cmdutil.AddDescriptionVarFlag(cmd, o.project.Description)
	cmdutil.AddLFSenabledVarPFlag(cmd, o.project.LFSEnabled)
	cmdutil.AddRequestAccessEnabledVarFlag(cmd, o.project.RequestAccessEnabled)
	cmdutil.AddVisibilityVarFlag(cmd, (*string)(o.project.Visibility))
	// unique flags for projects
	f := cmd.Flags()
	f.StringVar(o.project.Name, "name", *o.project.Name, "New project name")
	f.StringVar(o.project.Path, "path", *o.project.Path, "New project path")
	f.StringVar(o.project.DefaultBranch, "default-branch", *o.project.DefaultBranch, "The default branch")

	f.String("issues_access_level", "", "issues access level "+
		"(disabled,enabled,private,public)")
	f.String("merge_requests_access_level", "", "merge requests access level "+
		"(disabled,enabled,private,public)")
	f.String("builds_access_level", "", "builds access level "+
		"(disabled,enabled,private,public)")
	f.String("wiki_access_level", "", "wiki access level "+
		"(disabled,enabled,private,public)")
	f.String("snippets_access_level", "", "snippets access level "+
		"(disabled,enabled,private,public)")
	f.Bool("resolve-outdated-diff-discussions", false,
		"Automatically resolve merge request diffs discussions on lines "+
			"changed with a push")
	f.String("container_registry_access_level", "", "container registry access level for this project "+
		"(disabled,enabled,private,public)")
	f.Bool("shared-runners-enabled", false,
		"Enable shared runners for this project")
	f.Bool("public_builds", false,
		"enable public builds")
	f.Bool("only-allow-merge-if-pipeline-succeeds", false,
		"Set whether merge requests can only be merged with successful jobs")
	f.Bool("only-allow-merge-if-discussion-are-resolved", false,
		"Set whether merge requests can only be merged "+
			"when all the discussions are resolved")
	f.String("merge-method", "",
		"Set the merge method used. (available: 'merge', 'rebase_merge', 'ff')")
	f.StringSlice("tag-list", []string{},
		"The list of tags for a project; put array of tags, "+
			"that should be finally assigned to a project.\n"+
			"Example: --tag-list='tag1,tag2'")
	f.Bool("printing-merge-request-link-enabled", false,
		"Show link to create/view merge request "+
			"when pushing from the command line")
	f.String("ci-config-path", "", "The path to CI config file")

}

// Complete completes all the required options.
func (o *EditOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	gitlabClient, err := f.GitlabClient()
	if err != nil {
		return err
	}
	o.gitlabClient = gitlabClient
	return o.assignOptions(cmd)
}

// Validate makes sure there is no discrepency in command options.
func (o *EditOptions) Validate(cmd *cobra.Command, args []string) error {
	if err := validate.ValidateVisibilityFlagValue(cmd); err != nil {
		return err
	}
	if err := validate.ValidateMergeMethodValue(cmd); err != nil {
		return err
	}
	return nil
}

// Run executes a list subcommand using the specified options.
func (o *EditOptions) Run(args []string) error {
	project, _, err := o.gitlabClient.Projects.EditProject(args[0], o.project)
	if err != nil {
		return err
	}
	cmdutil.PrintProjectsOut(o.Out, o.ioStreams.Out, project)
	return nil
}
