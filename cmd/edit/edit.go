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

package edit

import (
	"github.com/spf13/cobra"

	"github.com/huhouhua/glctl/pkg/cli/genericiooptions"

	"github.com/huhouhua/glctl/cmd/resources/branch"
	"github.com/huhouhua/glctl/cmd/resources/file"
	"github.com/huhouhua/glctl/cmd/resources/group"
	cmdutil "github.com/huhouhua/glctl/cmd/util"
)

var editDesc = "Edit a Gitlab resource"

func NewEditCmd(f cmdutil.Factory, ioStreams genericiooptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "edit",
		Aliases:               []string{"e"},
		Short:                 editDesc,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(group.NewEditGroupCmd(f, ioStreams))
	cmd.AddCommand(branch.NewEditBranchCmd(f, ioStreams))
	cmd.AddCommand(file.NewEditFileCmd(f, ioStreams))
	return cmd
}
