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

package file

import (
	"github.com/AlekSi/pointer"
	cmdtesting "github.com/huhouhua/gitlab-repo-operator/cmd/testing"
	cmdutil "github.com/huhouhua/gitlab-repo-operator/cmd/util"
	"github.com/pkg/errors"
	"testing"
)

func TestGetFiles(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		optionsFunc    func() *ListOptions
		expectedOutput string
		wantError      error
	}{{
		name: "list all file",
		args: []string{"70"},
		optionsFunc: func() *ListOptions {
			opt := NewListOptions()
			opt.All = true
			return opt
		},
		wantError: nil,
	}, {
		name: "get specified branch",
		args: []string{
			"70",
		},
		optionsFunc: func() *ListOptions {
			opt := NewListOptions()
			opt.file.Ref = pointer.ToString("develop")
			return opt
		},
		wantError: nil,
	}, {
		name: "Get specified directory",
		args: []string{"70"},
		optionsFunc: func() *ListOptions {
			opt := NewListOptions()
			opt.file.Ref = pointer.ToString("develop")
			opt.file.Path = pointer.ToString("clusters")
			return opt
		},
		wantError: nil,
	}, {
		name: "list all projects with page",
		args: []string{"70"},
		optionsFunc: func() *ListOptions {
			opt := NewListOptions()
			opt.file.Ref = pointer.ToString("develop")
			opt.file.ListOptions.Page = 1
			opt.file.ListOptions.PerPage = 100
			return opt
		},
		wantError: nil,
	}, {
		name: "desc sort",
		args: []string{},
		optionsFunc: func() *ListOptions {
			opt := NewListOptions()
			opt.file.Sort = "desc"
			return opt
		},
		wantError: nil,
	}}
	factory := cmdutil.NewFactory(cmdtesting.NewFakeRESTClientGetter())
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmd := NewGetFilesCmd(factory)
			var cmdOptions *ListOptions
			if tc.optionsFunc != nil {
				cmdOptions = tc.optionsFunc()
			} else {
				cmdOptions = NewListOptions()
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
