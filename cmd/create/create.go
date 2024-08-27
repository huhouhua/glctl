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

package create

import (
	"github.com/spf13/cobra"

	"github.com/huhouhua/glctl/cmd/resources/branch"
	"github.com/huhouhua/glctl/cmd/resources/group"
	"github.com/huhouhua/glctl/cmd/resources/project"
	cmdutil "github.com/huhouhua/glctl/cmd/util"
	cli "github.com/huhouhua/glctl/util/cli"
)

var createDesc = "Create a Gitlab resource"

func NewCreateCmd(f cmdutil.Factory, ioStreams cli.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "create",
		Aliases:               []string{"c"},
		Short:                 createDesc,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(group.NewCreateGroupCmd(f, ioStreams))
	cmd.AddCommand(project.NewCreateProjectCmd(f, ioStreams))
	cmd.AddCommand(branch.NewCreateBranchCmd(f, ioStreams))
	return cmd
}
