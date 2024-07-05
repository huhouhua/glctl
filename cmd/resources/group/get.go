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
	"github.com/AlekSi/pointer"
	cmdutil "github.com/huhouhua/gl/cmd/util"
	"github.com/huhouhua/gl/cmd/validate"
	"github.com/huhouhua/gl/util/cli"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
	"strconv"
	"strings"
)

type ListOptions struct {
	gitlabClient *gitlab.Client
	group        *gitlab.ListGroupsOptions
	subGroup     *gitlab.ListSubGroupsOptions
	groupId      *int
	FromGroup    string
	Out          string
	AllGroups    bool
	ioStreams    cli.IOStreams
}

var (
	getGroupsDesc = "List groups and subgroups"

	getGroupsExample = `# list all groups
gl get groups

# list all subgroups of GroupX
gl get groups --all-groups=GroupX`
)

func NewListOptions(ioStreams cli.IOStreams) *ListOptions {
	return &ListOptions{
		ioStreams: ioStreams,
		group: &gitlab.ListGroupsOptions{
			AllAvailable: pointer.ToBool(false),
			OrderBy:      pointer.ToString("name"),
			Sort:         pointer.ToString("asc"),
			Search:       pointer.ToString(""),
			Statistics:   pointer.ToBool(false),
			Owned:        pointer.ToBool(false),
			ListOptions: gitlab.ListOptions{
				Page:    1,
				PerPage: 100,
			},
		},
		groupId:   nil,
		AllGroups: false,
		Out:       "simple",
	}
}

func NewGetGroupsCmd(f cmdutil.Factory, ioStreams cli.IOStreams) *cobra.Command {
	o := NewListOptions(ioStreams)
	cmd := &cobra.Command{
		Use:                   "groups",
		Aliases:               []string{"g"},
		Short:                 getGroupsDesc,
		Example:               getGroupsExample,
		DisableFlagsInUseLine: true,
		TraverseChildren:      true,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Validate(cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{"group"},
	}
	o.AddFlags(cmd)
	return cmd
}

func (o *ListOptions) AddFlags(cmd *cobra.Command) {
	cmdutil.AddPaginationVarFlags(cmd, &o.group.ListOptions)
	cmdutil.AddOwnedVarFlag(cmd, o.group.Owned)
	cmdutil.AddSortVarFlag(cmd, o.group.Sort)
	cmdutil.AddStatisticsVarFlag(cmd, o.group.Statistics)
	cmdutil.AddSearchVarFlag(cmd, o.group.Search)
	cmdutil.AddFromGroupVarPFlag(cmd, &o.FromGroup)
	cmdutil.AddOutFlag(cmd, &o.Out)
	f := cmd.Flags()
	f.BoolVar(o.group.AllAvailable, "all-available", *o.group.AllAvailable, "Show all the groups you have access to "+
		"(defaults to false for authenticated users, true for admin)")
	f.StringVar(o.group.OrderBy, "order-by", *o.group.OrderBy,
		"Order groups by name or path. Default is name")
	f.BoolVarP(&o.AllGroups, "all-groups", "A", o.AllGroups, "If present, list the requested object(s) across all groups. group in current context is ignored even if specified with --group.")
}

// Complete completes all the required options.
func (o *ListOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	var err error
	if len(args) > 0 {
		var id int
		id, err = strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("error from server (NotFound): group %s not found", args[0])
		}
		o.groupId = pointer.ToInt(id)
	}
	o.subGroup = (*gitlab.ListSubGroupsOptions)(o.group)
	o.gitlabClient, err = f.GitlabClient()
	return err
}

// Validate makes sure there is no discrepency in command options.
func (o *ListOptions) Validate(cmd *cobra.Command, args []string) error {
	if err := validate.ValidateGroupOrderByFlagValue(cmd); err != nil {
		return err
	}
	return validate.ValidateSortFlagValue(cmd)
}

// Run executes a list subcommand using the specified options.
func (o *ListOptions) Run(args []string) error {
	if o.groupId != nil {
		group, _, err := o.gitlabClient.Groups.GetGroup(*o.groupId, &gitlab.GetGroupOptions{})
		if err != nil {
			return err
		}
		cmdutil.PrintGroupsOut(o.Out, group)
		return nil
	}
	var groups []*gitlab.Group
	var err error
	if strings.TrimSpace(o.FromGroup) != "" {
		groups, _, err = o.gitlabClient.Groups.ListSubGroups(o.FromGroup, o.subGroup)
	} else {
		for {
			list, _, err := o.gitlabClient.Groups.ListGroups(o.group)
			if err != nil {
				return nil
			}
			groups = append(groups, list...)
			if cap(list) == 0 || !o.AllGroups {
				break
			}
			o.group.Page++
		}
	}
	if err != nil {
		return nil
	}
	cmdutil.PrintGroupsOut(o.Out, groups...)
	return nil
}
