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

package create

import (
	"github.com/spf13/cobra"

	"github.com/huhouhua/glctl/pkg/util/templates"

	"github.com/huhouhua/glctl/pkg/cli/genericiooptions"

	"github.com/huhouhua/glctl/cmd/resources/branch"
	"github.com/huhouhua/glctl/cmd/resources/group"
	"github.com/huhouhua/glctl/cmd/resources/project"
	cmdutil "github.com/huhouhua/glctl/cmd/util"
)

var (
	createDesc = "Create a resource from a file or from stdin"
	createLong = templates.LongDesc(`
		Create a resource from a file or from stdin.

		JSON and YAML formats are accepted.`)

	createExample = templates.Examples(`
		# Create a project using the data in project.json
		glctl create -f ./project.json

		# Create a project based on the JSON passed into stdin
		cat project.json | glctl create -f -

		# Edit the data in registry.yaml in JSON then create the resource using the edited data
		glctl create -f registry.yaml --edit -o json`)
)

func NewCreateCmd(f cmdutil.Factory, ioStreams genericiooptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "create -f FILENAME",
		Aliases:               []string{"c"},
		Short:                 createDesc,
		Long:                  createLong,
		Example:               createExample,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(group.NewCreateGroupCmd(f, ioStreams))
	cmd.AddCommand(project.NewCreateProjectCmd(f, ioStreams))
	cmd.AddCommand(branch.NewCreateBranchCmd(f, ioStreams))
	return cmd
}
