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

package group

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/huhouhua/glctl/pkg/cli/genericiooptions"
	"github.com/huhouhua/glctl/pkg/util/templates"

	"github.com/AlekSi/pointer"
	"github.com/spf13/cobra"
	gitlab "gitlab.com/gitlab-org/api/client-go"

	"github.com/huhouhua/glctl/cmd/require"
	cmdutil "github.com/huhouhua/glctl/cmd/util"
	"github.com/huhouhua/glctl/cmd/validate"
)

type CreateOptions struct {
	gitlabClient *gitlab.Client
	Group        *gitlab.CreateGroupOptions
	Namespace    string
	ioStreams    genericiooptions.IOStreams
	Out          string
}

var (
	createGroupExample = templates.Examples(`
# create a new group
glctl create group myGroup

# create a subgroup using namespace
glctl create group myGroup --namespace=ParentGroupX`)
)

func NewCreateOptions(ioStreams genericiooptions.IOStreams) *CreateOptions {
	return &CreateOptions{
		ioStreams: ioStreams,
		Group: &gitlab.CreateGroupOptions{
			Description:          pointer.ToString(""),
			RequestAccessEnabled: pointer.ToBool(false),
			Visibility:           pointer.To(gitlab.PrivateVisibility),
			LFSEnabled:           pointer.ToBool(false),
		},
		Out: "simple",
	}
}

func NewCreateGroupCmd(f cmdutil.Factory, ioStreams genericiooptions.IOStreams) *cobra.Command {
	o := NewCreateOptions(ioStreams)
	cmd := &cobra.Command{
		Use:                   "group",
		Aliases:               []string{"g"},
		Short:                 "Create a new group by specifying the group name as the first argument",
		Example:               createGroupExample,
		Args:                  require.ExactArgs(1),
		DisableFlagsInUseLine: true,
		TraverseChildren:      true,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Validate(cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
	}
	o.AddFlags(cmd)
	return cmd
}

func (o *CreateOptions) AddFlags(cmd *cobra.Command) {
	cmdutil.AddDescriptionVarFlag(cmd, o.Group.Description)
	cmdutil.AddRequestAccessEnabledVarFlag(cmd, o.Group.RequestAccessEnabled)
	cmdutil.AddVisibilityVarFlag(cmd, (*string)(o.Group.Visibility))
	cmdutil.AddOutFlag(cmd, &o.Out)
	f := cmd.Flags()
	f.BoolVar(o.Group.LFSEnabled, "lfs-enabled", *o.Group.LFSEnabled, "Enable LFS")
	f.StringVarP(&o.Namespace, "namespace", "n", o.Namespace,
		"This can be the parent namespace ID, group path, or user path. "+
			"(defaults to current user namespace)")
}

// Complete completes all the required options.
func (o *CreateOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	client, err := f.GitlabClient()
	if err != nil {
		return err
	}
	o.gitlabClient = client
	o.Group.Path = pointer.ToString(args[0])
	o.Group.Name = pointer.ToString(args[0])
	if strings.TrimSpace(o.Namespace) != "" {
		id, convErr := strconv.ParseInt(o.Namespace, 10, 0)
		// if not nil take the given number
		if convErr == nil {
			o.Group.ParentID = &id
			// find the group as string and get it's id
		} else {
			groupInfo, _, errGroup := o.gitlabClient.Groups.GetGroup(o.Namespace, &gitlab.GetGroupOptions{})
			if errGroup != nil {
				return errGroup
			}
			o.Group.ParentID = pointer.ToInt64(groupInfo.ID)
		}
	}
	return nil
}

// Validate makes sure there is no discrepency in command options.
func (o *CreateOptions) Validate(cmd *cobra.Command, args []string) error {
	err := validate.ValidateVisibilityFlagValue(cmd)
	return err
}

// Run executes a list subcommand using the specified options.
func (o *CreateOptions) Run(args []string) error {
	group, _, err := o.gitlabClient.Groups.CreateGroup(o.Group)
	if err != nil {
		return err
	}
	_, _ = fmt.Fprintf(o.ioStreams.Out, "%s created \n", group.FullPath)
	return nil
}
