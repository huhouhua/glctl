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
	"strconv"

	"github.com/AlekSi/pointer"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"

	"github.com/huhouhua/glctl/cmd/require"
	cmdutil "github.com/huhouhua/glctl/cmd/util"
	"github.com/huhouhua/glctl/cmd/validate"
	"github.com/huhouhua/glctl/util/cli"
	"github.com/huhouhua/glctl/util/templates"
)

type CreateOptions struct {
	gitlabClient *gitlab.Client
	project      *gitlab.CreateProjectOptions
	namespace    string
	Out          string
	ioStreams    cli.IOStreams
}

var (
	createProjectExample = templates.Examples(`
# create a new project
glctl create project myProject --desc="Project X is party!"

# create a new project under a group
glctl create project myProject --namespace=GroupY`)
)

func NewCreateOptions(ioStreams cli.IOStreams) *CreateOptions {
	return &CreateOptions{
		ioStreams: ioStreams,
		project: &gitlab.CreateProjectOptions{
			Description:                               pointer.ToString(""),
			LFSEnabled:                                pointer.ToBool(false),
			RequestAccessEnabled:                      pointer.ToBool(false),
			ResolveOutdatedDiffDiscussions:            pointer.ToBool(false),
			SharedRunnersEnabled:                      pointer.ToBool(false),
			PublicBuilds:                              pointer.ToBool(false),
			OnlyAllowMergeIfPipelineSucceeds:          pointer.ToBool(false),
			OnlyAllowMergeIfAllDiscussionsAreResolved: pointer.ToBool(false),
			Visibility:                                pointer.To(gitlab.PrivateVisibility),
			IssuesAccessLevel:                         pointer.To(gitlab.EnabledAccessControl),
			MergeRequestsAccessLevel:                  pointer.To(gitlab.EnabledAccessControl),
			BuildsAccessLevel:                         pointer.To(gitlab.EnabledAccessControl),
			WikiAccessLevel:                           pointer.To(gitlab.EnabledAccessControl),
			SnippetsAccessLevel:                       pointer.To(gitlab.EnabledAccessControl),
			ContainerRegistryAccessLevel:              pointer.To(gitlab.DisabledAccessControl),
			MergeMethod:                               pointer.To(gitlab.NoFastForwardMerge),
			Topics:                                    pointer.To([]string{}),
			PrintingMergeRequestLinkEnabled:           pointer.ToBool(false),
			CIConfigPath:                              pointer.ToString(""),
		},
		Out: "simple",
	}
}

func NewCreateProjectCmd(f cmdutil.Factory, ioStreams cli.IOStreams) *cobra.Command {
	o := NewCreateOptions(ioStreams)
	cmd := &cobra.Command{
		Use:                   "project",
		Aliases:               []string{"p"},
		Short:                 "Create a new project by specifying the project name as the first argument",
		Example:               createProjectExample,
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
	cmdutil.AddDescriptionVarFlag(cmd, o.project.Description)
	cmdutil.AddLFSenabledVarPFlag(cmd, o.project.LFSEnabled)
	cmdutil.AddRequestAccessEnabledVarFlag(cmd, o.project.RequestAccessEnabled)
	cmdutil.AddVisibilityVarFlag(cmd, (*string)(o.project.Visibility))
	// unique flags for projects
	f := cmd.Flags()
	f.StringVar(
		(*string)(o.project.IssuesAccessLevel),
		"issues_access_level",
		(string)(*o.project.IssuesAccessLevel),
		"issues access level "+
			"(disabled,enabled,private,public)",
	)
	f.StringVar(
		(*string)(o.project.MergeRequestsAccessLevel),
		"merge_requests_access_level",
		(string)(*o.project.MergeRequestsAccessLevel),
		"merge requests access level "+
			"(disabled,enabled,private,public)",
	)
	f.StringVar(
		(*string)(o.project.BuildsAccessLevel),
		"builds_access_level",
		(string)(*o.project.BuildsAccessLevel),
		"builds access level "+
			"(disabled,enabled,private,public)",
	)
	f.StringVar(
		(*string)(o.project.WikiAccessLevel),
		"wiki_access_level",
		(string)(*o.project.WikiAccessLevel),
		"wiki access level "+
			"(disabled,enabled,private,public)",
	)
	f.StringVar(
		(*string)(o.project.SnippetsAccessLevel),
		"snippets_access_level",
		(string)(*o.project.SnippetsAccessLevel),
		"snippets access level "+
			"(disabled,enabled,private,public)",
	)
	f.BoolVar(
		o.project.ResolveOutdatedDiffDiscussions,
		"resolve-outdated-diff-discussions",
		*o.project.ResolveOutdatedDiffDiscussions,
		"Automatically resolve merge request diffs discussions on lines "+
			"changed with a push",
	)
	f.StringVar(
		(*string)(o.project.ContainerRegistryAccessLevel),
		"container_registry_access_level",
		(string)(*o.project.ContainerRegistryAccessLevel),
		"container registry access level for this project "+
			"(disabled,enabled,private,public)",
	)
	f.BoolVar(o.project.SharedRunnersEnabled, "shared-runners-enabled", *o.project.SharedRunnersEnabled,
		"Enable shared runners for this project")
	f.BoolVar(o.project.PublicBuilds, "public_builds", *o.project.PublicBuilds,
		"enable public builds")
	f.BoolVar(
		o.project.OnlyAllowMergeIfPipelineSucceeds,
		"only-allow-merge-if-pipeline-succeeds",
		*o.project.OnlyAllowMergeIfPipelineSucceeds,
		"Set whether merge requests can only be merged with successful jobs",
	)
	f.BoolVar(
		o.project.OnlyAllowMergeIfAllDiscussionsAreResolved,
		"only-allow-merge-if-discussion-are-resolved",
		*o.project.OnlyAllowMergeIfAllDiscussionsAreResolved,
		"Set whether merge requests can only be merged "+
			"when all the discussions are resolved",
	)
	f.StringVar((*string)(o.project.MergeMethod), "merge-method", (string)(*o.project.MergeMethod),
		"Set the merge method used. (available: 'merge', 'rebase_merge', 'ff')")
	f.StringSliceVar(o.project.Topics, "tag-list", *o.project.Topics,
		"The list of tags for a project; put array of tags, "+
			"that should be finally assigned to a project.\n"+
			"Example: --tag-list='tag1,tag2'")
	f.BoolVar(
		o.project.PrintingMergeRequestLinkEnabled,
		"printing-merge-request-link-enabled",
		*o.project.PrintingMergeRequestLinkEnabled,
		"Show link to create/view merge request "+
			"when pushing from the command line",
	)
	f.StringVar(o.project.CIConfigPath, "ci-config-path", *o.project.CIConfigPath, "The path to CI config file")

	f.StringVarP(&o.namespace, "namespace", "n",
		o.namespace, "This can be the parent namespace ID, group path, or user path. "+
			"(defaults to current user namespace)")

}

// Complete completes all the required options.
func (o *CreateOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	var err error
	o.gitlabClient, err = f.GitlabClient()
	o.project.Name = pointer.ToString(args[0])
	o.project.Path = pointer.ToString(args[0])
	gid, convErr := strconv.Atoi(o.namespace)
	// if namespace is not a number,
	// get the namespace's group id and assign it to gid
	if convErr != nil {
		ns, _, err := o.gitlabClient.Namespaces.GetNamespace(o.namespace)
		if err == nil {
			gid = ns.ID
		}
	}
	o.project.NamespaceID = pointer.ToInt(gid)
	return err
}

// Validate makes sure there is no discrepency in command options.
func (o *CreateOptions) Validate(cmd *cobra.Command, args []string) error {
	if err := validate.ValidateVisibilityFlagValue(cmd); err != nil {
		return err
	}
	if err := validate.ValidateMergeMethodValue(cmd); err != nil {
		return err
	}
	return nil
}

// Run executes a list subcommand using the specified options.
func (o *CreateOptions) Run(args []string) error {
	project, _, err := o.gitlabClient.Projects.CreateProject(o.project)
	if err != nil {
		return err
	}
	cmdutil.PrintProjectsOut(o.Out, o.ioStreams.Out, project)
	return nil
}
