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
	"github.com/AlekSi/pointer"
	cmdtesting "github.com/huhouhua/gl/cmd/testing"
	cmdutil "github.com/huhouhua/gl/cmd/util"
	"github.com/huhouhua/gl/util/cli"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
	"strings"
	"testing"
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
			opt.project = "huhouhua/gitlab-repo-branch"
			opt.branch.Ref = pointer.ToString("main")
		},
		run: func(opt *CreateOptions, args []string) error {
			var err error
			out := cmdtesting.RunTestForStdout(func() {
				err = opt.Run(args)
			})
			if !strings.Contains(out, *opt.branch.Branch) {
				err = errors.New(fmt.Sprintf("delete by path : Unexpected output! Expected\n%s\ngot\n%s", *opt.branch.Branch, out))
			}
			return err
		},
		wantError: nil,
	}, {
		name: "create an existing branch",
		args: []string{"create1"},
		optionsFunc: func(opt *CreateOptions) {
			opt.project = "huhouhua/gitlab-repo-branch"
		},
		run: func(opt *CreateOptions, args []string) error {
			err := opt.Run(args)
			var repo *gitlab.ErrorResponse
			if errors.As(err, &repo) && repo.Message == "{error: ref is empty}" {
				return nil
			}
			return err
		},
	}}
	streams := cli.NewTestIOStreamsDiscard()
	factory := cmdutil.NewFactory(cmdtesting.NewFakeRESTClientGetter())
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmd := NewCreateBranchCmd(factory, streams)
			cmdOptions := NewCreateOptions(streams)
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
