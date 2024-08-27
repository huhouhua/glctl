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

package group

import (
	"fmt"
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

type EditOptions struct {
	gitlabClient *gitlab.Client
	groupId      int
	Group        *gitlab.UpdateGroupOptions
	ioStreams    cli.IOStreams
	Out          string
}

var (
	editGroupExample = templates.Examples(`
# edit a group
gctl edit group myGroupAZ --desc="Updated group"

# edit a subgroup
gctl edit group myGroupX/myGroupZ --desc="Updated group"

# edit a group with id (23)
gctl edit group 23 --visibility="public`)
)

func NewEditOptions(ioStreams cli.IOStreams) *EditOptions {
	return &EditOptions{
		ioStreams: ioStreams,
		Group: &gitlab.UpdateGroupOptions{
			Description:          pointer.ToString(""),
			RequestAccessEnabled: pointer.ToBool(false),
			Visibility:           pointer.To(gitlab.PrivateVisibility),
			LFSEnabled:           pointer.ToBool(false),
		},
		Out: "simple",
	}
}

func NewEditGroupCmd(f cmdutil.Factory, ioStreams cli.IOStreams) *cobra.Command {
	o := NewEditOptions(ioStreams)
	cmd := &cobra.Command{
		Use:                   "group",
		Aliases:               []string{"g"},
		Short:                 "Update a group by specifying the group id or path and using flags for fields to modify",
		Example:               editGroupExample,
		Args:                  require.ExactArgs(1),
		DisableFlagsInUseLine: true,
		TraverseChildren:      true,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Validate(cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{"groups"},
	}
	o.AddFlags(cmd)
	return cmd
}

func (o *EditOptions) AddFlags(cmd *cobra.Command) {
	cmdutil.AddDescriptionVarFlag(cmd, o.Group.Description)
	cmdutil.AddRequestAccessEnabledVarFlag(cmd, o.Group.RequestAccessEnabled)
	cmdutil.AddVisibilityVarFlag(cmd, (*string)(o.Group.Visibility))
	cmdutil.AddOutFlag(cmd, &o.Out)
	f := cmd.Flags()
	f.StringVar(o.Group.Name, "name", *o.Group.Name,
		"New group name")
	f.StringVar(o.Group.Path, "path", *o.Group.Path,
		"New group path")
	f.BoolVar(o.Group.LFSEnabled, "lfs-enabled", *o.Group.LFSEnabled, "Enable LFS")
}

// Complete completes all the required options.
func (o *EditOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	var err error
	o.gitlabClient, err = f.GitlabClient()
	if len(args) > 0 {
		gid := args[0]
		o.groupId, err = strconv.Atoi(gid)
		// if group is not a number,
		// search for the group path's id and assign it to gid
		if err != nil {
			groupInfo, _, err := o.gitlabClient.Groups.GetGroup(gid, &gitlab.GetGroupOptions{})
			if err != nil {
				return fmt.Errorf("couldn't find the id of group %s, got error: %v",
					gid, err)
			}
			o.groupId = groupInfo.ID
		}
	}
	return err
}

// Validate makes sure there is no discrepency in command options.
func (o *EditOptions) Validate(cmd *cobra.Command, args []string) error {
	err := validate.ValidateVisibilityFlagValue(cmd)
	return err
}

// Run executes a list subcommand using the specified options.
func (o *EditOptions) Run(args []string) error {
	group, _, err := o.gitlabClient.Groups.UpdateGroup(o.groupId, o.Group)
	if err != nil {
		return err
	}
	cmdutil.PrintGroupsOut(o.Out, o.ioStreams.Out, group)
	return nil
}
