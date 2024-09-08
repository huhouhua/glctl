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
	"github.com/huhouhua/glctl/pkg/cli/genericiooptions"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AlekSi/pointer"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"

	cmdtesting "github.com/huhouhua/glctl/cmd/testing"
	cmdutil "github.com/huhouhua/glctl/cmd/util"
)

func TestCreateBranch(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		optionsFunc func(opt *CreateOptions)
		validate    func(opt *CreateOptions, cmd *cobra.Command, args []string) error
		run         func(opt *CreateOptions, args []string) error
		wantError   error
	}{{
		name: "create a new branch",
		args: []string{"create1"},
		optionsFunc: func(opt *CreateOptions) {
			opt.project = "Group1/Project2"
			opt.branch.Ref = pointer.ToString("main")
		},
		run: func(opt *CreateOptions, args []string) error {
			defer func() {
				_, _ = opt.gitlabClient.Branches.DeleteBranch(opt.project, *opt.branch.Branch)
			}()
			var err error
			out := cmdtesting.RunForStdout(opt.ioStreams, func() {
				err = opt.Run(args)
			})
			assert.Containsf(
				t,
				out,
				*opt.branch.Branch,
				"create a new branch: Unexpected output! Expected\n%s\ngot\n%s",
				*opt.branch.Branch,
				out,
			)
			return err
		},
		wantError: nil,
	}, {
		name: "create an existing branch",
		args: []string{"master"},
		optionsFunc: func(opt *CreateOptions) {
			opt.project = "Group1/Project2"
		},
		run: func(opt *CreateOptions, args []string) error {
			err := opt.Run(args)
			var repoErr *gitlab.ErrorResponse
			assert.ErrorAs(t, err, &repoErr)
			if assert.Equal(t, repoErr.Message, "{error: ref is empty}") {
				return nil
			}
			return err
		},
	}}
	factory := cmdutil.NewFactory(cmdtesting.NewFakeRESTClientGetter())
	streams := genericiooptions.NewTestIOStreamsForPipe()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmd := NewCreateBranchCmd(factory, streams)
			cmdOptions := NewCreateOptions(streams)
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
