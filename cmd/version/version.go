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

package version

import (
	cmdutil "github.com/huhouhua/gl/cmd/util"
	"github.com/huhouhua/gl/util/templates"
	"github.com/spf13/cobra"
)

var versionExample = templates.Examples(`
		# Print the client and server versions for the current context
		gl version`)

// NewCmdVersion returns a cobra command for fetching versions.
func NewCmdVersion(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "version",
		Short:   "Print the client and server version information",
		Long:    "Print the client and server version information for the current context",
		Example: versionExample,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	return cmd
}
