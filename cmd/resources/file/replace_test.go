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

package file

import (
	"testing"

	"github.com/huhouhua/glctl/pkg/cli/genericiooptions"

	"github.com/spf13/cobra"

	cmdtesting "github.com/huhouhua/glctl/cmd/testing"
	cmdutil "github.com/huhouhua/glctl/cmd/util"
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
		name: "replace all branch test/test.yaml file",
		args: []string{"test/test.yaml"},
		optionsFunc: func(opt *ReplaceOptions) {
			opt.Project = "Group2/SubGroup3/Project13"
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
		args: []string{"test/test-force.yaml"},
		optionsFunc: func(opt *ReplaceOptions) {
			opt.Project = "Group2/SubGroup3/Project13"
			opt.RefMatch = "*"
			opt.FileName = "../../../testdata/replace/force_replace.yaml"
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
	streams := genericiooptions.NewTestIOStreamsDiscard()
	factory := cmdutil.NewFactory(cmdtesting.NewFakeRESTClientGetter())
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmd := NewReplaceFileCmd(factory, streams)
			var cmdOptions = NewReplaceOptions(streams)
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
