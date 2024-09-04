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

package group

import (
	"fmt"
	"strings"
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/xanzy/go-gitlab"

	cmdtesting "github.com/huhouhua/glctl/cmd/testing"
	cmdutil "github.com/huhouhua/glctl/cmd/util"
	"github.com/huhouhua/glctl/util/cli"
)

func TestCreateGroup(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		optionsFunc func(opt *CreateOptions)
		validate    func(opt *CreateOptions, cmd *cobra.Command, args []string) error
		run         func(opt *CreateOptions, args []string) error
		wantError   error
	}{{
		name: "create sub group",
		args: []string{"Test_New_Group_Under_Group1"},
		optionsFunc: func(opt *CreateOptions) {
			opt.Group.Description = pointer.ToString("Created by to test")
			opt.Namespace = "Group1"
		},
		wantError: nil,
	}, {
		name: "create group with on namespace",
		args: []string{"Test_New_Group_Without_Namespace"},
		optionsFunc: func(opt *CreateOptions) {
			opt.Group.Description = pointer.ToString("Created by to test")
		},
		wantError: nil,
	}, {
		name: "create group using id in namespace",
		args: []string{"Test_New_Group_Using_Namespace"},
		optionsFunc: func(opt *CreateOptions) {
			opt.Group.Description = pointer.ToString("Created by to test")
			opt.Namespace = "2" // is GitLab Instance
		},
		wantError: nil,
	}, {
		name: "create an existing group",
		args: []string{"Group1"},
		run: func(opt *CreateOptions, args []string) error {
			err := opt.Run(args)
			var repoErr *gitlab.ErrorResponse
			assert.ErrorAs(t, err, &repoErr)
			if assert.Equal(
				t,
				repoErr.Message,
				"{message: Failed to save group {:path=>[\"has already been taken\"]}}",
			) {
				return nil
			}
			return err
		},
	}}
	factory := cmdutil.NewFactory(cmdtesting.NewFakeRESTClientGetter())
	for _, tc := range tests {
		streams := cli.NewTestIOStreamsForPipe()
		t.Run(tc.name, func(t *testing.T) {
			cmd := NewCreateGroupCmd(factory, streams)
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
				cmdtesting.ErrorAssertionWithEqual(t, tc.wantError, err)
				return
			}
			namespace := cmdOptions.Namespace
			if strings.TrimSpace(namespace) != "" {
				group, _, groupErr := cmdOptions.gitlabClient.Groups.GetGroup(namespace, &gitlab.GetGroupOptions{})
				cmdtesting.ErrorAssertionWithEqual(t, tc.wantError, groupErr)
				if groupErr != nil {
					return
				}
				namespace = group.FullPath + "/"
			}
			pathFull := fmt.Sprintf("%s%s", namespace, *cmdOptions.Group.Path)
			defer func() {
				_, _ = cmdOptions.gitlabClient.Groups.DeleteGroup(pathFull, &gitlab.DeleteGroupOptions{})
			}()
			out := cmdtesting.RunForStdout(streams, func() {
				err = cmdOptions.Run(tc.args)
			})
			cmdtesting.ErrorAssertionWithEqual(t, tc.wantError, err)

			expectedOutput := fmt.Sprintf("%s created", pathFull)
			assert.Containsf(
				t,
				out,
				expectedOutput,
				"%s: Unexpected output! Expected\n%s\ngot\n%s",
				tc.name,
				expectedOutput,
				out)

		})
	}
}
