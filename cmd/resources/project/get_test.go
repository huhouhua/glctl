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

package project

import (
	"github.com/huhouhua/glctl/pkg/cli/genericiooptions"
	"testing"

	"github.com/AlekSi/pointer"

	cmdtesting "github.com/huhouhua/glctl/cmd/testing"
	cmdutil "github.com/huhouhua/glctl/cmd/util"
)

func TestGetProjects(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		optionsFunc    func(opt *ListOptions)
		expectedOutput string
		wantError      error
	}{{
		name: "list all projects",
		args: []string{},
		optionsFunc: func(opt *ListOptions) {
			opt.AllGroups = true
		},
		wantError: nil,
	}, {
		name: "project by id",
		args: []string{
			"1",
		},
		wantError: nil,
	}, {
		name: "list all projects with page",
		args: []string{},
		optionsFunc: func(opt *ListOptions) {
			opt.project.ListOptions.Page = 1
			opt.project.ListOptions.PerPage = 100
		},
		wantError: nil,
	}, {
		name: "desc sort",
		args: []string{},
		optionsFunc: func(opt *ListOptions) {
			opt.project.Sort = pointer.ToString("desc")
		},
		wantError: nil,
	}}
	streams := genericiooptions.NewTestIOStreamsDiscard()
	factory := cmdutil.NewFactory(cmdtesting.NewFakeRESTClientGetter())
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmd := NewGetProjectsCmd(factory, streams)
			var cmdOptions = NewListOptions(streams)
			if tc.optionsFunc != nil {
				tc.optionsFunc(cmdOptions)
			}
			err := cmdOptions.Complete(factory, cmd, tc.args)
			cmdtesting.ErrorAssertionWithEqual(t, tc.wantError, err)
			if err != nil {
				return
			}
			err = cmdOptions.Validate(cmd, tc.args)
			cmdtesting.ErrorAssertionWithEqual(t, tc.wantError, err)
			if err != nil {
				return
			}
			err = cmdOptions.Run(tc.args)
			cmdtesting.ErrorAssertionWithEqual(t, tc.wantError, err)
			if err != nil {
				return
			}
		})
	}
}
