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

package replace

import (
	"github.com/huhouhua/glctl/pkg/cli/genericiooptions"
	"github.com/huhouhua/glctl/pkg/util/templates"
	"github.com/spf13/cobra"

	"github.com/huhouhua/glctl/cmd/resources/file"
	cmdutil "github.com/huhouhua/glctl/cmd/util"
)

var (
	replaceLong = templates.LongDesc(`
		Replace a repository file by file name or stdin.

		all formats are accepted. If replacing an existing repository file, the
		complete repository file spec must be provided. This can be obtained by

		    $ glctl get files PROJECT --path=my.yml --ref=BRANCH --raw`)

	replaceExample = templates.Examples(`
		# Replace a single file using the data in my.yml
		glctl replace file app/my.yml -f ./my.yml --ref=main --project=myproject 

		# Replace all branch file using the data in my.yml
	    glctl replace file app/my.yml -f ./my.yml --ref-match=* --project=myproject`)
)

func NewReplaceCmd(f cmdutil.Factory, ioStreams genericiooptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "replace",
		Aliases:               []string{"r"},
		Short:                 replaceLong,
		Example:               replaceExample,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(file.NewReplaceFileCmd(f, ioStreams))
	return cmd
}
