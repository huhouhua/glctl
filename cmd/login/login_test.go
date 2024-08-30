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

package login

import (
	"fmt"
	"strings"
	"testing"

	cmdtesting "github.com/huhouhua/glctl/cmd/testing"
	"github.com/huhouhua/glctl/util/cli"
)

func TestLogin(t *testing.T) {
	tests := []struct {
		name           string
		optionsFunc    func(opt *Options)
		args           []string
		expectedOutput string
	}{
		{
			name: "login incorrect password",
			args: []string{"http://localhost:8080"},
			optionsFunc: func(opt *Options) {
				opt.User = "root"
				opt.Password = "12345"
			},
			expectedOutput: "Login failed!\nThe provided authorization grant is invalid, expired, revoked, does not match the redirection URI used in the authorization request, or was issued to another client.",
		}, {
			name: "login incorrect address",
			args: []string{"http://1.3.4.5"},
			optionsFunc: func(opt *Options) {
				opt.User = "12345"
				opt.Password = "12345"
			},
			expectedOutput: "dial tcp 1.3.4.5:80: i/o timeout",
		},
		{
			name: "login success",
			args: []string{"http://localhost:8080"},
			optionsFunc: func(opt *Options) {
				opt.User = "john.doe"
				opt.Password = "123qwe123"
			},
			expectedOutput: "\nLogin Succeeded",
		}}
	for _, tc := range tests {
		streams := cli.NewTestIOStreamsForPipe()
		t.Run(tc.name, func(t *testing.T) {
			cmd := NewLoginCmd(streams)
			var cmdOptions = NewOptions(streams)
			if tc.optionsFunc != nil {
				tc.optionsFunc(cmdOptions)
			}
			out := cmdtesting.RunForStdout(streams, func() {
				if err := cmdOptions.Complete(cmd, tc.args); err != nil {
					_, _ = fmt.Fprint(streams.Out, err)
					return
				}
				if err := cmdOptions.Validate(cmd, tc.args); err != nil {
					_, _ = fmt.Fprint(streams.Out, err)
					return
				}
				if err := cmdOptions.Run(tc.args); err != nil {
					_, _ = fmt.Fprint(streams.Out, err)
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
	}{
		{
			name: "login address is empty",
			args: []string{""},
			optionsFunc: func(opt *Options) {
				opt.User = "john.doe"
				opt.Password = "123qwe123"
			},
			expectedOutput: "please enter the gitlab url",
		}, {
			name: "login username is empty",
			args: []string{"http://localhost:8080"},
			optionsFunc: func(opt *Options) {
				opt.User = ""
				opt.Password = "123qwe123"
			},
			expectedOutput: "please enter the username",
		},
		{
			name: "login password is empty",
			args: []string{"http://localhost:8080"},
			optionsFunc: func(opt *Options) {
				opt.User = "john.doe"
				opt.Password = ""
			},
			expectedOutput: "please enter the password",
		}}
	streams := cli.NewTestIOStreamsDiscard()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmd := NewLoginCmd(streams)
			var cmdOptions = NewOptions(streams)
			if tc.optionsFunc != nil {
				tc.optionsFunc(cmdOptions)
			}
			out := cmdtesting.Run(func() {
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

func TestRunLogin(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		flags          map[string]string
		expectedOutput string
	}{
		{
			name:           "login gitlab success",
			args:           []string{"http://localhost:8080"},
			flags:          map[string]string{"username": "john.smith", "password": "123qwe123"},
			expectedOutput: "Login Succeeded",
		},
	}
	for _, tc := range tests {
		streams := cli.NewTestIOStreamsForPipe()
		t.Run(tc.name, func(t *testing.T) {
			for i, arg := range tc.args {
				cmdtesting.TInfo(fmt.Sprintf("(%d) %s", i, arg))
			}
			cmd := NewLoginCmd(streams)
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
