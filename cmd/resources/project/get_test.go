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

package project

import (
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/pkg/errors"

	cmdtesting "github.com/huhouhua/glctl/cmd/testing"
	cmdutil "github.com/huhouhua/glctl/cmd/util"
	"github.com/huhouhua/glctl/util/cli"
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
			"6",
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
	streams := cli.NewTestIOStreamsDiscard()
	factory := cmdutil.NewFactory(cmdtesting.NewFakeRESTClientGetter())
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmd := NewGetProjectsCmd(factory, streams)
			var cmdOptions = NewListOptions(streams)
			if tc.optionsFunc != nil {
				tc.optionsFunc(cmdOptions)
			}
			var err error
			if err = cmdOptions.Complete(factory, cmd, tc.args); !errors.Is(err, tc.wantError) {
				t.Errorf("expected %v, got: '%v'", tc.wantError, err)
				return
			}
			if err = cmdOptions.Validate(cmd, tc.args); !errors.Is(err, tc.wantError) {
				t.Errorf("expected %v, got: '%v'", tc.wantError, err)
				return
			}
			if err = cmdOptions.Run(tc.args); !errors.Is(err, tc.wantError) {
				t.Errorf("expected %v, got: '%v'", tc.wantError, err)
				return
			}
		})
	}
}
