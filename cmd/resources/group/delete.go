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
	"github.com/huhouhua/glctl/cmd/require"
	cmdutil "github.com/huhouhua/glctl/cmd/util"
	"github.com/huhouhua/glctl/util/cli"
	"github.com/huhouhua/glctl/util/templates"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
	"strconv"
)

type DeleteOptions struct {
	gitlabClient *gitlab.Client
	groupId      int
	ioStreams    cli.IOStreams
	Out          string
}

var (
	deleteGroupExample = templates.Examples(`
# delete a Group named GroupX
glctl delete group GroupX

# delete a Subgroup named GroupY under GroupX
glctl delete group GroupX/GroupY

# delete a group with id (3)
glctl delete group 3`)
)

func NewDeleteOptions(ioStreams cli.IOStreams) *DeleteOptions {
	return &DeleteOptions{
		ioStreams: ioStreams,
		Out:       "simple",
	}
}

func NewDeleteGroupCmd(f cmdutil.Factory, ioStreams cli.IOStreams) *cobra.Command {
	o := NewDeleteOptions(ioStreams)
	cmd := &cobra.Command{
		Use:                   "group",
		Aliases:               []string{"g"},
		Short:                 "Delete a Gitlab group by specifying the id or group path",
		Example:               deleteGroupExample,
		Args:                  require.ExactArgs(1),
		DisableFlagsInUseLine: true,
		TraverseChildren:      true,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Validate(cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
	}
	cmdutil.AddOutFlag(cmd, &o.Out)
	return cmd
}

// Complete completes all the required options.
func (o *DeleteOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	var err error
	o.gitlabClient, err = f.GitlabClient()
	if len(args) > 0 {
		o.groupId, err = strconv.Atoi(args[0])
		// if group is not a number,
		// search for the group path's id and assign it to gid
		if err == nil {
			groupInfo, _, err := o.gitlabClient.Groups.GetGroup("namespace", &gitlab.GetGroupOptions{})
			if err != nil {
				return fmt.Errorf("couldn't find the id of group %s, got error: %v",
					args[0], err)
			}
			o.groupId = groupInfo.ID
		}
	}
	return err
}

// Validate makes sure there is no discrepency in command options.
func (o *DeleteOptions) Validate(cmd *cobra.Command, args []string) error {
	return nil
}

// Run executes a list subcommand using the specified options.
func (o *DeleteOptions) Run(args []string) error {
	_, err := o.gitlabClient.Groups.DeleteGroup(o.groupId, &gitlab.DeleteGroupOptions{})
	if err != nil {
		return err
	}
	_, _ = fmt.Fprintf(o.ioStreams.Out, "Group (%s) with id (%d) has been deleted\n", args[0], o.groupId)
	return nil
}
