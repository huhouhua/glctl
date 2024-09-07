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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xanzy/go-gitlab"

	cmdtesting "github.com/huhouhua/glctl/cmd/testing"
	cmdutil "github.com/huhouhua/glctl/cmd/util"
	"github.com/huhouhua/glctl/util/cli"
)

func TestDeleteGroup(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		optionsFunc func(opt *DeleteOptions)
		run         func(opt *DeleteOptions, args []string) error
		wantError   error
	}{
		{
			name:      "delete group with on parent",
			args:      []string{"Test_DELETE_Group_Without_Parent"},
			wantError: nil,
		}, {
			name:      "delete sub group",
			args:      []string{"Group2/Test_DELETE_Group_Under_Group1"},
			wantError: nil,
		},
	}
	factory := cmdutil.NewFactory(cmdtesting.NewFakeRESTClientGetter())
	for _, tc := range tests {
		streams := cli.NewTestIOStreamsForPipe()
		t.Run(tc.name, func(t *testing.T) {
			cmd := NewDeleteGroupCmd(factory, streams)
			var cmdOptions = NewDeleteOptions(streams)
			if tc.optionsFunc != nil {
				tc.optionsFunc(cmdOptions)
			}
			group, groupErr := createTestGroup(t, factory, tc.args[0], tc.wantError)
			if groupErr != nil {
				return
			}
			defer func() {
				_, _ = cmdOptions.gitlabClient.Groups.DeleteGroup(group.ID, &gitlab.DeleteGroupOptions{})
			}()
			var err error
			err = cmdOptions.Complete(factory, cmd, tc.args)
			cmdtesting.ErrorAssertionWithEqual(t, tc.wantError, err)
			if err != nil {
				return
			}
			err = cmdOptions.Validate(cmd, tc.args)
			cmdtesting.ErrorAssertionWithEqual(t, tc.wantError, err)
			if err != nil {
				return
			}
			if tc.run != nil {
				err = tc.run(cmdOptions, tc.args)
				cmdtesting.ErrorAssertionWithEqual(t, tc.wantError, err)
				return
			}
			out := cmdtesting.RunForStdout(streams, func() {
				err = cmdOptions.Run(tc.args)
			})
			cmdtesting.ErrorAssertionWithEqual(t, tc.wantError, err)
			expectedOutput := fmt.Sprintf("Group (%s) with id (%d) has been deleted", tc.args[0], group.ID)
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
func TestDeleteWithComplete(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		optionsFunc func(opt *DeleteOptions)
		handler     func(err error)
		wantError   error
	}{
		{
			name: "deleting a non existent group should fail",
			args: []string{"GroupNotFount"},
			handler: func(err error) {
				require.NotNil(t, err)
				assert.Contains(t, err.Error(), "{message: 404 Group Not Found}")
			},
			wantError: nil,
		}}
	factory := cmdutil.NewFactory(cmdtesting.NewFakeRESTClientGetter())
	for _, tc := range tests {
		streams := cli.NewTestIOStreamsForPipe()
		t.Run(tc.name, func(t *testing.T) {
			cmd := NewDeleteGroupCmd(factory, streams)
			var cmdOptions = NewDeleteOptions(streams)
			if tc.optionsFunc != nil {
				tc.optionsFunc(cmdOptions)
			}
			err := cmdOptions.Complete(factory, cmd, tc.args)
			if tc.handler != nil {
				tc.handler(err)
			} else {
				cmdtesting.ErrorAssertionWithEqual(t, tc.wantError, err)
			}
		})
	}
}

func createTestGroup(t *testing.T, factory cmdutil.Factory, groupName string, wantError error) (*gitlab.Group, error) {
	client, err := factory.GitlabClient()
	cmdtesting.ErrorAssertionWithEqual(t, wantError, err)

	opt := &gitlab.CreateGroupOptions{
		Path: pointer.ToString(groupName),
		Name: pointer.ToString(groupName),
	}
	groupPath := strings.SplitAfter(groupName, "/")
	if len(groupPath) > 1 {
		lastIndex := len(groupPath) - 1
		parentName := strings.Join(groupPath[:lastIndex], "/")
		group, _, groupErr := client.Groups.GetGroup(strings.TrimSuffix(parentName, "/"), &gitlab.GetGroupOptions{})
		cmdtesting.ErrorAssertionWithEqual(t, wantError, groupErr)
		if err != nil {
			return nil, err
		}
		subName := groupPath[lastIndex]
		opt.ParentID = pointer.ToInt(group.ID)
		opt.Name = pointer.ToString(subName)
		opt.Path = pointer.ToString(subName)
	}
	newGroup, _, groupErr := client.Groups.CreateGroup(opt)
	cmdtesting.ErrorAssertionWithEqual(t, wantError, groupErr)
	return newGroup, groupErr
}
