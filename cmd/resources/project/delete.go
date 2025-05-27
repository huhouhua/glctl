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
	"fmt"
	"strings"

	"github.com/huhouhua/glctl/pkg/cli/genericiooptions"
	"github.com/huhouhua/glctl/pkg/util/templates"

	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"

	"github.com/huhouhua/glctl/cmd/require"
	cmdutil "github.com/huhouhua/glctl/cmd/util"
)

type DeleteOptions struct {
	gitlabClient *gitlab.Client
	project      string
	ioStreams    genericiooptions.IOStreams
}

var (
	deleteProjectExample = templates.Examples(`
# delete a project
glctl delete project ProjectX

# delete a project under a group
glctl delete project group/project`)
)

func NewDeleteOptions(ioStreams genericiooptions.IOStreams) *DeleteOptions {
	return &DeleteOptions{
		ioStreams: ioStreams,
	}
}

func NewDeleteProjectCmd(f cmdutil.Factory, ioStreams genericiooptions.IOStreams) *cobra.Command {
	o := NewDeleteOptions(ioStreams)
	cmd := &cobra.Command{
		Use:                   "project",
		Aliases:               []string{"p"},
		Short:                 "Delete a Gitlab project by specifying the full path",
		Example:               deleteProjectExample,
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
	return cmd
}

// Complete completes all the required options.
func (o *DeleteOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	var err error
	o.gitlabClient, err = f.GitlabClient()
	if len(args) > 0 {
		o.project = args[0]
	}
	return err
}

// Validate makes sure there is no discrepency in command options.
func (o *DeleteOptions) Validate(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("please enter project name or id")
	}
	if strings.TrimSpace(args[0]) == "" {
		return fmt.Errorf("error from server (NotFound): project %s not found", args[0])
	}
	return nil
}

// Run executes a list subcommand using the specified options.
func (o *DeleteOptions) Run(args []string) error {
	projectInfo, _, err := o.gitlabClient.Projects.GetProject(o.project, &gitlab.GetProjectOptions{})
	if err != nil {
		return err
	}
	_, err = o.gitlabClient.Projects.DeleteProject(projectInfo.ID)
	if err != nil {
		return err
	}
	_, _ = fmt.Fprintf(o.ioStreams.Out, "project (%s) with id (%d) has been deleted\n", o.project, projectInfo.ID)
	return nil
}
