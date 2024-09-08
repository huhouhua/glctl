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
	"github.com/huhouhua/glctl/pkg/cli/genericiooptions"
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xanzy/go-gitlab"

	cmdtesting "github.com/huhouhua/glctl/cmd/testing"
	cmdutil "github.com/huhouhua/glctl/cmd/util"
)

func TestEditGroup(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		optionsFunc func(opt *EditOptions)
		validate    func(opt *EditOptions, cmd *cobra.Command, args []string) error
		run         func(opt *EditOptions, args []string) error
		wantError   error
	}{
		{
			name: "edit an existing group",
			args: []string{"Group5"},
			optionsFunc: func(opt *EditOptions) {
				opt.Group.Description = pointer.ToString("Updated by go test")
				opt.Group.Visibility = pointer.To(gitlab.InternalVisibility)
			},
			wantError: nil,
		}, {
			name: "edit an existing subgroup",
			args: []string{"Group2/SubGroup4"},
			optionsFunc: func(opt *EditOptions) {
				opt.Group.Description = pointer.ToString("Updated by go test")
			},
			wantError: nil,
		}, {
			name: "edit an existing group by id",
			args: []string{"27"}, // is Group4/SubGroup8
			optionsFunc: func(opt *EditOptions) {
				opt.Group.Description = pointer.ToString("Updated by go test by id")
			},
			wantError: nil,
		}, {
			name: "change a group name and path",
			args: []string{"Group4/SubGroup7"},
			optionsFunc: func(opt *EditOptions) {
				opt.Group.Description = pointer.ToString("Updated by go test")
				opt.Group.Name = pointer.ToString("Renamed")
			},
			wantError: nil,
		}}
	factory := cmdutil.NewFactory(cmdtesting.NewFakeRESTClientGetter())
	for _, tc := range tests {
		streams := genericiooptions.NewTestIOStreamsForPipe()
		t.Run(tc.name, func(t *testing.T) {
			cmd := NewEditGroupCmd(factory, streams)
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
				cmdtesting.ErrorAssertionWithEqual(t, tc.wantError, err)
				return
			}
			group, _, groupErr := cmdOptions.gitlabClient.Groups.GetGroup(tc.args[0], &gitlab.GetGroupOptions{})
			cmdtesting.ErrorAssertionWithEqual(t, tc.wantError, groupErr)
			if err != nil {
				return
			}
			defer func() {
				revertGroup(cmdOptions.gitlabClient, tc.args[0], group)
			}()
			out := cmdtesting.RunForStdout(streams, func() {
				err = cmdOptions.Run(tc.args)
			})
			cmdtesting.ErrorAssertionWithEqual(t, tc.wantError, err)

			expectedOutput := fmt.Sprintf("%s configured", group.FullPath)
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

func TestEditWithComplete(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		optionsFunc func(opt *EditOptions)
		handler     func(err error)
		wantError   error
	}{
		{
			name: "editing a non existent group must fail",
			args: []string{"GroupNotFount"},
			optionsFunc: func(opt *EditOptions) {
				opt.Group.Description = pointer.ToString("Updated by go test")
			},
			handler: func(err error) {
				require.NotNil(t, err)
				assert.Contains(t, err.Error(), "{message: 404 Group Not Found}")
			},
			wantError: nil,
		}}
	factory := cmdutil.NewFactory(cmdtesting.NewFakeRESTClientGetter())
	for _, tc := range tests {
		streams := genericiooptions.NewTestIOStreamsForPipe()
		t.Run(tc.name, func(t *testing.T) {
			cmd := NewEditGroupCmd(factory, streams)
			var cmdOptions = NewEditOptions(streams)
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

func revertGroup(client *gitlab.Client, groupName string, group *gitlab.Group) {
	_, _, _ = client.Groups.UpdateGroup(groupName, &gitlab.UpdateGroupOptions{
		Name:                 pointer.ToString(group.Name),
		Path:                 pointer.ToString(group.Path),
		Description:          pointer.ToString(group.Description),
		RequestAccessEnabled: pointer.ToBool(group.RequestAccessEnabled),
		Visibility:           pointer.To(group.Visibility),
		LFSEnabled:           pointer.ToBool(group.LFSEnabled),
	})
}
