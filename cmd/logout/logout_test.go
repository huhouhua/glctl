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

package logout

import (
	"fmt"
	cmdtesting "github.com/huhouhua/glctl/cmd/testing"
	"github.com/huhouhua/glctl/util/cli"
	"strings"
	"testing"
)

func TestLogout(t *testing.T) {
	tests := []struct {
		name           string
		optionsFunc    func(opt *Options)
		args           []string
		expectedOutput string
	}{}
	streams := cli.NewTestIOStreamsForPipe()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmd := NewLogoutCmd(streams)
			var cmdOptions = NewOptions(streams)
			if tc.optionsFunc != nil {
				tc.optionsFunc(cmdOptions)
			}
			out := cmdtesting.RunTestForStdout(streams, func() {
				var err error
				if err = cmdOptions.Complete(cmd, tc.args); err != nil {
					fmt.Print(err)
					return
				}
				if err = cmdOptions.Validate(cmd, tc.args); err != nil {
					fmt.Print(err)
					return
				}
				if err = cmdOptions.Run(tc.args); err != nil {
					fmt.Print(err)
					return
				}
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

func TestValidate(t *testing.T) {
	tests := []struct {
		name           string
		optionsFunc    func(opt *Options)
		args           []string
		expectedOutput string
	}{}
	streams := cli.NewTestIOStreamsForPipe()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmd := NewLogoutCmd(streams)
			var cmdOptions = NewOptions(streams)
			if tc.optionsFunc != nil {
				tc.optionsFunc(cmdOptions)
			}
			out := cmdtesting.RunTestForStdout(streams, func() {
				var err error
				if err = cmdOptions.Complete(cmd, tc.args); err != nil {
					fmt.Print(err)
					return
				}
				if err = cmdOptions.Validate(cmd, tc.args); err != nil {
					fmt.Print(err)
					return
				}
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

func TestRunLogout(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		flags          map[string]string
		expectedOutput string
	}{}
	streams, _, buf, _ := cli.NewTestIOStreams()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for i, arg := range tc.args {
				cmdtesting.TInfo(fmt.Sprintf("(%d) %s", i, arg))
			}
			cmd := NewLogoutCmd(streams)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			for flag, value := range tc.flags {
				err := cmd.Flags().Set(flag, value)
				if err != nil {
					t.Errorf("set %s flag error", err.Error())
					return
				}
			}
			out := cmdtesting.RunTestForStdout(streams, func() {
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
