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

package get

import (
	"github.com/huhouhua/gitlab-repo-operator/cmd/resources/branch"
	"github.com/huhouhua/gitlab-repo-operator/cmd/resources/file"
	"github.com/huhouhua/gitlab-repo-operator/cmd/resources/group"
	"github.com/huhouhua/gitlab-repo-operator/cmd/resources/project"
	cmdutil "github.com/huhouhua/gitlab-repo-operator/cmd/util"
	"github.com/huhouhua/gitlab-repo-operator/cmd/validate"
	"github.com/spf13/cobra"
)

var getDesc = "Get Gitlab resources"

func NewGetCmd(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "get",
		Aliases:               []string{"g"},
		Short:                 getDesc,
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return validate.ValidateOutFlagValue(cmd)
		},
	}
	cmd.AddCommand(project.NewGetProjectsCmd(f))
	cmd.AddCommand(group.NewGetGroupsCmd(f))
	cmd.AddCommand(branch.NewGetBranchesCmd(f))
	cmd.AddCommand(file.NewGetFilesCmd(f))
	return cmd
}
