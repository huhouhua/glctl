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

	"github.com/huhouhua/glctl/pkg/cli/genericiooptions"
	"github.com/huhouhua/glctl/pkg/util/templates"

	"github.com/AlekSi/pointer"
	"github.com/spf13/cobra"
	gitlab "gitlab.com/gitlab-org/api/client-go"

	"github.com/huhouhua/glctl/cmd/require"
	cmdutil "github.com/huhouhua/glctl/cmd/util"
	"github.com/huhouhua/glctl/cmd/validate"
)

type EditOptions struct {
	gitlabClient *gitlab.Client
	groupId      int
	Group        *gitlab.UpdateGroupOptions
	ioStreams    genericiooptions.IOStreams
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

func NewEditOptions(ioStreams genericiooptions.IOStreams) *EditOptions {
	return &EditOptions{
		ioStreams: ioStreams,
		Group:     &gitlab.UpdateGroupOptions{},
		Out:       "simple",
	}
}

func NewEditGroupCmd(f cmdutil.Factory, ioStreams genericiooptions.IOStreams) *cobra.Command {
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
	cmdutil.AddDescriptionFlag(cmd)
	cmdutil.AddRequestAccessEnabledFlag(cmd, false)
	cmdutil.AddVisibilityFlag(cmd)
	cmdutil.AddOutFlag(cmd, &o.Out)
	f := cmd.Flags()
	f.String("name", "",
		"New group name")
	f.String("path", "",
		"New group path")
	f.Bool("lfs-enabled", false, "Enable LFS")
}

// Complete completes all the required options.
func (o *EditOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	client, err := f.GitlabClient()
	if err != nil {
		return err
	}
	o.gitlabClient = client
	o.assignOptions(cmd)
	gid, err := GroupNameFromCommandArgs(cmd, args)
	if err != nil {
		return err
	}
	o.groupId, err = strconv.Atoi(gid)
	// if group is not a number,
	// search for the group path's id and assign it to gid
	if err != nil {
		groupInfo, _, errGroup := o.gitlabClient.Groups.GetGroup(gid, &gitlab.GetGroupOptions{})
		if errGroup != nil {
			return fmt.Errorf("couldn't find the id of group %s, got error: %v",
				gid, errGroup)
		}
		o.groupId = groupInfo.ID
	}
	return nil
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
	_, _ = fmt.Fprintf(o.ioStreams.Out, "%s configured \n", group.FullPath)
	return nil
}

// assign cmd flag to options
func (o *EditOptions) assignOptions(cmd *cobra.Command) {
	if cmd.Flag("desc").Changed {
		o.Group.Description = pointer.ToString(cmdutil.GetFlagString(cmd, "desc"))
	}
	if cmd.Flag("request-access-enabled").Changed {
		o.Group.RequestAccessEnabled = pointer.ToBool(cmdutil.GetFlagBool(cmd, "request-access-enabled"))
	}
	if cmd.Flag("visibility").Changed {
		o.Group.Visibility = pointer.To(gitlab.VisibilityValue(cmdutil.GetFlagString(cmd, "visibility")))
	}
	if cmd.Flag("name").Changed {
		o.Group.Name = pointer.ToString(cmdutil.GetFlagString(cmd, "name"))
	}
	if cmd.Flag("path").Changed {
		o.Group.Path = pointer.ToString(cmdutil.GetFlagString(cmd, "path"))
	}
	if cmd.Flag("lfs-enabled").Changed {
		o.Group.LFSEnabled = pointer.ToBool(cmdutil.GetFlagBool(cmd, "lfs-enabled"))
	}
}

// GroupNameFromCommandArgs is a utility function for commands that assume the first argument is a group name
func GroupNameFromCommandArgs(cmd *cobra.Command, args []string) (string, error) {
	argsLen := cmd.ArgsLenAtDash()
	// ArgsLenAtDash returns -1 when -- was not specified
	if argsLen == -1 {
		argsLen = len(args)
	}
	if argsLen != 1 {
		return "", cmdutil.UsageErrorf(cmd, "exactly one NAME is required, got %d", argsLen)
	}
	return args[0], nil
}
