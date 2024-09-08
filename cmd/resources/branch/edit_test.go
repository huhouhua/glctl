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
	"fmt"
	"github.com/huhouhua/glctl/pkg/cli/genericiooptions"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xanzy/go-gitlab"

	"github.com/AlekSi/pointer"
	"github.com/spf13/cobra"

	cmdtesting "github.com/huhouhua/glctl/cmd/testing"
	cmdutil "github.com/huhouhua/glctl/cmd/util"
)

func TestEditBranch(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		optionsFunc func(opt *EditOptions)
		validate    func(opt *EditOptions, cmd *cobra.Command, args []string) error
		run         func(opt *EditOptions, args []string) error
		wantError   error
	}{{
		name: "set unprotect",
		args: []string{"unprotect"},
		optionsFunc: func(opt *EditOptions) {
			opt.project = "Group1/Project1"
			opt.Unprotect = true
		},
		run: func(opt *EditOptions, args []string) error {
			var err error
			_, _, err = opt.gitlabClient.ProtectedBranches.ProtectRepositoryBranches(
				opt.project,
				&gitlab.ProtectRepositoryBranchesOptions{
					Name: pointer.ToString(args[0]),
				},
			)
			if err != nil {
				return err
			}
			defer func() {
				_, _ = opt.gitlabClient.ProtectedBranches.UnprotectRepositoryBranches(opt.project, args[0])
			}()
			out := cmdtesting.RunForStdout(opt.ioStreams, func() {
				err = opt.Run(args)
			})
			expectedOutput := fmt.Sprintf("branch %s un protect", args[0])
			assert.Containsf(
				t,
				out,
				expectedOutput,
				"set \"unprotect\" branch s unprotect : Unexpected output! Expected\n%s\ngot\n%s",
				expectedOutput,
				out,
			)
			return err
		},
		wantError: nil,
	}, {
		name: "set protect",
		args: []string{"protect"},
		optionsFunc: func(opt *EditOptions) {
			opt.project = "Group1/Project2"
			opt.protect = true
			opt.protectBranch.DevelopersCanMerge = pointer.ToBool(true)
			opt.protectBranch.DevelopersCanPush = pointer.ToBool(true)
		},
		run: func(opt *EditOptions, args []string) error {
			defer func() {
				_, _ = opt.gitlabClient.ProtectedBranches.UnprotectRepositoryBranches(opt.project, args[0])
			}()
			var err error
			out := cmdtesting.RunForStdout(opt.ioStreams, func() {
				err = opt.Run(args)
			})
			expectedOutput := fmt.Sprintf("branch %s updated", args[0])
			assert.Containsf(
				t,
				out,
				expectedOutput,
				"set \"updated\" branch s updated : Unexpected output! Expected\n%s\ngot\n%s",
				expectedOutput,
				out,
			)
			return err
		},
		wantError: nil,
	}}
	factory := cmdutil.NewFactory(cmdtesting.NewFakeRESTClientGetter())
	for _, tc := range tests {
		streams := genericiooptions.NewTestIOStreamsForPipe()
		t.Run(tc.name, func(t *testing.T) {
			cmd := NewEditBranchCmd(factory, streams)
			var cmdOptions = NewEditOptions(streams)
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
