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

package cmd

import (
	"encoding/json"
	cmdutil "github.com/huhouhua/gitlab-repo-operator/cmd/util"
	"github.com/huhouhua/gitlab-repo-operator/cmd/validate"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

var getProjectsDesc = "List projects of the authenticated user or of a group"

var getProjectsExample = `# get all projects
gitlabctl get projects

# get all projects from a group
gitlabctl get projects --from-group=Group1`

func newGetProjectsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "projects",
		Short:             getProjectsDesc,
		SilenceErrors:     true,
		SilenceUsage:      true,
		DisableAutoGenTag: true,
		Example:           getProjectsExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			if cmdutil.GetFlagString(cmd, "from-group") != "" {
				return runGetProjectsFromGroup(cmd)
			}
			return runGetProjects(cmd)
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := validate.ValidateSortFlagValue(cmd); err != nil {
				return err
			}
			if err := validate.ValidateProjectOrderByFlagValue(cmd); err != nil {
				return err
			}
			return validate.ValidateVisibilityFlagValue(cmd)
		},
	}
	cmdutil.AddOutFlag(cmd)
	cmdutil.AddPaginationFlags(cmd)
	return cmd
}

func runGetProjects(cmd *cobra.Command) error {
	opts := cmdutil.AssignListProjectOptions(cmd)
	git, err := newClient()
	if err != nil {
		return nil
	}
	projects, _, err := git.Projects.ListProjects(opts)
	if err != nil {
		return nil
	}
	cmdutil.PrintProjectsOut(cmdutil.GetFlagString(cmd, "out"), projects...)
	return nil
}

func runGetProjectsFromGroup(cmd *cobra.Command) error {
	optstr, err := json.Marshal(cmdutil.AssignListProjectOptions(cmd))
	if err != nil {
		return err
	}
	opt := &gitlab.ListGroupProjectsOptions{}
	if err = json.Unmarshal(optstr, opt); err != nil {
		return err
	}
	git, err := newClient()
	if err != nil {
		return err
	}
	projects, _, err := git.Groups.ListGroupProjects(cmdutil.GetFlagString(cmd, "from-group"), opt)
	if err != nil {
		return err
	}

	cmdutil.PrintProjectsOut(cmdutil.GetFlagString(cmd, "out"), projects...)
	return nil
}
