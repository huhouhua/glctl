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

package file

import (
	"fmt"
	"github.com/AlekSi/pointer"
	cmdtesting "github.com/huhouhua/glctl/cmd/testing"
	cmdutil "github.com/huhouhua/glctl/cmd/util"
	"github.com/huhouhua/glctl/util/cli"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"strings"
	"testing"
)

func TestDeleteFile(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		optionsFunc func(opt *DeleteOptions)
		validate    func(opt *DeleteOptions, cmd *cobra.Command, args []string) error
		run         func(opt *DeleteOptions, args []string) error
		wantError   error
	}{{
		name: "delete by file",
		args: []string{"delete.yaml"},
		optionsFunc: func(opt *DeleteOptions) {
			opt.Project = "223"
			opt.file.Branch = pointer.ToString("main")
		},
		run: func(opt *DeleteOptions, args []string) error {
			var err error
			out := cmdtesting.RunTest(func() {
				err = opt.Run(args)
			})
			expectedOutput := fmt.Sprintf("file (%s) for %s branch with project id (%s) has been deleted", opt.FileName, *opt.file.Branch, opt.Project)
			if !strings.Contains(out, expectedOutput) {
				err = errors.New(fmt.Sprintf("delete by path : Unexpected output! Expected\n%s\ngot\n%s", expectedOutput, out))
			}
			return err
		},
		wantError: nil,
	}, { //to do
		name: "delete by dir",
		args: []string{"/weqwee"},
		optionsFunc: func(opt *DeleteOptions) {
			opt.Project = "223"
			opt.file.Branch = pointer.ToString("main")
		},
		run: func(opt *DeleteOptions, args []string) error {
			var err error
			out := cmdtesting.RunTest(func() {
				err = opt.Run(args)
			})
			expectedOutput := fmt.Sprintf("file (%s) for %s branch with project id (%s) has been deleted", opt.FileName, *opt.file.Branch, opt.Project)
			if !strings.Contains(out, expectedOutput) {
				err = errors.New(fmt.Sprintf("delete by path : Unexpected output! Expected\n%s\ngot\n%s", expectedOutput, out))
			}
			return err
		},
		wantError: nil,
	}}
	streams := cli.NewTestIOStreamsDiscard()
	factory := cmdutil.NewFactory(cmdtesting.NewFakeRESTClientGetter())
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmd := NewDeleteFilesCmd(factory, streams)
			var cmdOptions = NewDeleteOptions(streams)
			if tc.optionsFunc != nil {
				tc.optionsFunc(cmdOptions)
			}
			var err error
			if err = cmdOptions.Complete(factory, cmd, tc.args); err != nil && !errors.Is(err, tc.wantError) {
				t.Errorf("expected %v, got: '%v'", tc.wantError, err)
				return
			}
			if tc.validate != nil {
				err = tc.validate(cmdOptions, cmd, tc.args)
				if err != nil {
					return
				}
			} else {
				if err = cmdOptions.Validate(cmd, tc.args); err != nil && !errors.Is(err, tc.wantError) {
					t.Errorf("expected %v, got: '%v'", tc.wantError, err)
					return
				}
			}
			if tc.run != nil {
				err = tc.run(cmdOptions, tc.args)
				if err != nil {
					t.Error(err)
				}
				return
			}
			if err = cmdOptions.Run(tc.args); !errors.Is(err, tc.wantError) {
				t.Errorf("expected %v, got: '%v'", tc.wantError, err)
				return
			}
		})
	}
}
