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

package branch

import (
	"errors"
	"fmt"
	"testing"

	"github.com/huhouhua/glctl/pkg/cli/genericiooptions"

	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/assert"

	cmdutil "github.com/huhouhua/glctl/cmd/util"

	"github.com/spf13/cobra"
	gitlab "gitlab.com/gitlab-org/api/client-go"

	cmdtesting "github.com/huhouhua/glctl/cmd/testing"
)

func TestDeleteBranch(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		optionsFunc func(opt *DeleteOptions)
		validate    func(opt *DeleteOptions, cmd *cobra.Command, args []string) error
		run         func(opt *DeleteOptions, args []string) error
		wantError   error
	}{{
		name: "delete branch success",
		args: []string{"develop-by-name"},
		optionsFunc: func(opt *DeleteOptions) {
			opt.project = "Group1/Project3"
		},
		run: func(opt *DeleteOptions, args []string) error {
			var err error
			_, _, err = opt.gitlabClient.Branches.CreateBranch(opt.project, &gitlab.CreateBranchOptions{
				Branch: pointer.ToString(opt.branch),
				Ref:    pointer.ToString("main"),
			})
			if err != nil {
				return err
			}
			defer func() {
				_, _ = opt.gitlabClient.Branches.DeleteBranch(opt.project, opt.branch)
			}()
			out := cmdtesting.RunForStdout(opt.ioStreams, func() {
				err = opt.Run(args)
			})
			expectedOutput := fmt.Sprintf("Branch (%s) from project (%s) has been deleted", opt.branch, opt.project)
			assert.Containsf(
				t,
				out,
				expectedOutput,
				"delete branch success: Unexpected output! Expected\n%s\ngot\n%s",
				expectedOutput,
				out,
			)
			return err
		},
		wantError: nil,
	}, {
		name: "branch not found",
		args: []string{"not-found"},
		optionsFunc: func(opt *DeleteOptions) {
			opt.project = "Group1/Project3"
		},
		run: func(opt *DeleteOptions, args []string) error {
			err := opt.Run(args)
			var repoErr error
			assert.ErrorAs(t, err, &repoErr)
			if assert.Equal(t, repoErr.Error(), "404 Not Found") {
				return nil
			}
			return err
		},
		wantError: nil,
	}, {
		name: "project not found",
		args: []string{"not-found"},
		optionsFunc: func(opt *DeleteOptions) {
			opt.project = "not-found"
		},
		run: func(opt *DeleteOptions, args []string) error {
			err := opt.Run(args)
			var repoErr error
			assert.ErrorAs(t, err, &repoErr)
			if assert.Equal(t, repoErr.Error(), "404 Not Found") {
				return nil
			}
			return err
		},
		wantError: nil,
	}, {
		name:      "not definition branch",
		wantError: errors.New("please enter branch"),
	}}
	factory := cmdutil.NewFactory(cmdtesting.NewFakeRESTClientGetter())
	streams := genericiooptions.NewTestIOStreamsForPipe()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmd := NewDeleteBranchCmd(factory, streams)
			var cmdOptions = NewDeleteOptions(streams)
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
