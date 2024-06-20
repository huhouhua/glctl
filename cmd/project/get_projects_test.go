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

package project

import (
	"bytes"
	cmdutil "github.com/huhouhua/gitlab-repo-operator/cmd/util"
	"github.com/spf13/cobra"
	"testing"
)

func TestGetProjects(t *testing.T) {
	tests := []cmdutil.CmdTestCase{{
		Name:      "list all projects",
		Cmd:       "",
		WantError: true,
	}}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			out, err := cmdutil.ExecuteCommand(func(buffer *bytes.Buffer) (*cobra.Command, error) {
				return NewGetProjectsCmd(cmdutil.NewTestFactory()), nil
			}, tc.Cmd)

			cmdutil.TInfo(out)
			if tc.WantError && err == nil {
				t.Errorf("expected error, got success with the following output:\n%s", out)
			}
			if !tc.WantError && err != nil {
				t.Errorf("expected no error, got: '%v'", err)
			}
		})
	}
}
