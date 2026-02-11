// Copyright 2024 The Kevin Berger <huhouhuam@gmail.com> Authors
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

	cmdutil "github.com/huhouhua/glctl/cmd/util"
	"github.com/huhouhua/glctl/cmd/validate"
)

type ListOptions struct {
	gitlabClient *gitlab.Client
	group        *gitlab.ListGroupsOptions
	subGroup     *gitlab.ListSubGroupsOptions
	groupId      *int
	FromGroup    string
	Out          string
	AllGroups    bool
	ioStreams    genericiooptions.IOStreams
}

var (
	getGroupsExample = templates.Examples(`
# list all groups
glctl get groups

# list all subgroups of GroupX
glctl get groups --all-groups=GroupX`)
)

func NewListOptions(ioStreams genericiooptions.IOStreams) *ListOptions {
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

func NewGetGroupsCmd(f cmdutil.Factory, ioStreams genericiooptions.IOStreams) *cobra.Command {
	o := NewListOptions(ioStreams)
	cmd := &cobra.Command{
		Use:                   "groups",
		Aliases:               []string{"g"},
		Short:                 "List groups and subgroups",
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
		return cmdutil.PrintGroupsOut(o.Out, o.ioStreams.Out, group)
	}
	var groups []*gitlab.Group
	var err error
	if strings.TrimSpace(o.FromGroup) != "" {
		groups, _, err = o.gitlabClient.Groups.ListSubGroups(o.FromGroup, o.subGroup)
	} else {
		for {
			var portion []*gitlab.Group
			portion, _, err = o.gitlabClient.Groups.ListGroups(o.group)
			if err != nil {
				return nil
			}
			groups = append(groups, portion...)
			if cap(portion) == 0 || !o.AllGroups {
				break
			}
			o.group.Page++
		}
	}
	if err != nil {
		return nil
	}
	return cmdutil.PrintGroupsOut(o.Out, o.ioStreams.Out, groups...)
}
