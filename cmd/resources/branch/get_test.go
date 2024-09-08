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

package branch

import (
	"errors"
	"fmt"
	"github.com/huhouhua/glctl/pkg/cli/genericiooptions"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	cmdtesting "github.com/huhouhua/glctl/cmd/testing"
	cmdutil "github.com/huhouhua/glctl/cmd/util"
)

func TestGetBranch(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		optionsFunc func(opt *ListOptions)
		validate    func(opt *ListOptions, cmd *cobra.Command, args []string) error
		run         func(opt *ListOptions, args []string) error
		wantError   error
	}{{
		name:      "project name is an empty string",
		args:      []string{""},
		wantError: errors.New("error from server (NotFound): project  not found"),
	}}
	streams := genericiooptions.NewTestIOStreamsForPipe()
	factory := cmdutil.NewFactory(cmdtesting.NewFakeRESTClientGetter())
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmd := NewGetBranchesCmd(factory, streams)
			cmdOptions := NewListOptions(streams)
			if tc.optionsFunc != nil {
				tc.optionsFunc(cmdOptions)
			}
			var err error
			err = cmdOptions.Complete(factory, cmd, tc.args)
			cmdtesting.ErrorAssertionWithEqual(t, tc.wantError, err)
			if err != nil {
				return
			}
			if tc.validate != nil {
				err = tc.validate(cmdOptions, cmd, tc.args)
			} else {
				err = cmdOptions.Validate(cmd, tc.args)
			}
			cmdtesting.ErrorAssertionWithEqual(t, tc.wantError, err)
			if err != nil {
				return
			}
			if tc.run != nil {
				err = tc.run(cmdOptions, tc.args)
			} else {
				err = cmdOptions.Run(tc.args)
			}
			cmdtesting.ErrorAssertionWithEqual(t, tc.wantError, err)
		})
	}
}

func TestRunGetBranch(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		flags          map[string]string
		expectedOutput string
	}{{
		name:           "list all branch with project",
		args:           []string{"Group1/Project1"},
		flags:          map[string]string{"all": "true"},
		expectedOutput: "main",
	}, {
		name:           "list all branch with project id",
		args:           []string{"2"},
		flags:          map[string]string{"all": "true"},
		expectedOutput: "main",
	}, {
		name:           "list all branch default",
		args:           []string{"3"},
		flags:          map[string]string{},
		expectedOutput: "main",
	}}
	factory := cmdutil.NewFactory(cmdtesting.NewFakeRESTClientGetter())
	for _, tc := range tests {
		streams := genericiooptions.NewTestIOStreamsForPipe()
		t.Run(tc.name, func(t *testing.T) {
			for i, arg := range tc.args {
				cmdtesting.TInfo(fmt.Sprintf("(%d) %s", i, arg))
			}
			cmd := NewGetBranchesCmd(factory, streams)
			cmd.SetOut(streams.Out)
			cmd.SetErr(streams.ErrOut)
			for flag, value := range tc.flags {
				err := cmd.Flags().Set(flag, value)
				if err != nil {
					t.Errorf("set %s flag error", err.Error())
					return
				}
			}
			out := cmdtesting.RunForStdout(streams, func() {
				cmd.Run(cmd, tc.args)
			})
			cmdtesting.TInfo(out)
			if tc.expectedOutput == "" {
				t.Errorf("%s: Invalid test case. Specify expected result.\n", tc.name)
			}
			if !strings.Contains(out, tc.expectedOutput) {
				t.Errorf("%s: Unexpected output! Expected\n%s\ngot\n%s", tc.name, tc.expectedOutput, out)
			}
		})
	}

}
