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
	"errors"
	"testing"

	"github.com/spf13/cobra"

	cmdtesting "github.com/huhouhua/glctl/cmd/testing"
	cmdutil "github.com/huhouhua/glctl/cmd/util"
	"github.com/huhouhua/glctl/util/cli"
)

func TestRunReplace(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		optionsFunc func(opt *ReplaceOptions)
		validate    func(opt *ReplaceOptions, cmd *cobra.Command, args []string) error
		run         func(opt *ReplaceOptions, args []string) error
		wantError   error
	}{{
		name: "replace all branch",
		args: []string{"resources/pipeline/buildserver/package.yaml"},
		optionsFunc: func(opt *ReplaceOptions) {
			opt.Project = "224"
			opt.RefMatch = "*"
			opt.FileName = "../../../testdata/replace/new_test.yaml"
		},
		run: func(opt *ReplaceOptions, args []string) error {
			var err error
			_ = cmdtesting.Run(func() {
				err = opt.Run(args)
			})
			return err
		},
		wantError: nil,
	}, {
		name: "force replace",
		args: []string{"package.yaml"},
		optionsFunc: func(opt *ReplaceOptions) {
			opt.Project = "224"
			opt.RefMatch = "*"
			opt.FileName = "../../../testdata/replace/new_test.yaml"
			opt.Force = true
		},
		run: func(opt *ReplaceOptions, args []string) error {
			var err error
			_ = cmdtesting.Run(func() {
				err = opt.Run(args)
			})
			return err
		},
		wantError: nil,
	}}
	streams := cli.NewTestIOStreamsDiscard()
	factory := cmdutil.NewFactory(cmdtesting.NewFakeRESTClientGetter())
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmd := NewReplaceFileCmd(factory, streams)
			var cmdOptions = NewReplaceOptions(streams)
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
