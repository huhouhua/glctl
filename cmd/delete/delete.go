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

package delete

import (
	"github.com/huhouhua/glctl/cmd/resources/branch"
	"github.com/huhouhua/glctl/cmd/resources/file"
	"github.com/huhouhua/glctl/cmd/resources/project"
	cmdutil "github.com/huhouhua/glctl/cmd/util"
	"github.com/huhouhua/glctl/util/cli"
	"github.com/spf13/cobra"
)

var deleteDesc = "Delete a Gitlab resource"

func NewDeleteCmd(f cmdutil.Factory, ioStreams cli.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "delete",
		Aliases:               []string{"d"},
		Short:                 deleteDesc,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(project.NewDeleteProjectCmd(f, ioStreams))
	cmd.AddCommand(branch.NewDeleteBranchCmd(f, ioStreams))
	cmd.AddCommand(file.NewDeleteFilesCmd(f, ioStreams))
	return cmd
}
