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
	cmdtesting "github.com/huhouhua/gl/cmd/testing"
	cmdutil "github.com/huhouhua/gl/cmd/util"
	"github.com/pkg/errors"
	"testing"
)

func TestGetFiles(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		optionsFunc    func(opt *ListOptions)
		expectedOutput string
		wantError      error
	}{{
		name: "list all file",
		args: []string{"70"},
		optionsFunc: func(opt *ListOptions) {
			opt.All = true
		},
		wantError: nil,
	}, {
		name: "get specified branch",
		args: []string{
			"70",
		},
		optionsFunc: func(opt *ListOptions) {
			opt.file.Ref = pointer.ToString("develop")
		},
		wantError: nil,
	}, {
		name: "read to receive the raw file in repository",
		args: []string{
			"70",
		},
		optionsFunc: func(opt *ListOptions) {
			opt.file.Ref = pointer.ToString("develop")
			opt.path = "clusters/devops/manifests/local-path-provisioner.yaml"
			opt.Raw = true
		},
		wantError: nil,
	}, {
		name: "Get specified directory",
		args: []string{"70"},
		optionsFunc: func(opt *ListOptions) {
			opt.file.Ref = pointer.ToString("develop")
			opt.file.Path = pointer.ToString("clusters")
		},
		wantError: nil,
	}, {
		name: "list all projects with page",
		args: []string{"70"},
		optionsFunc: func(opt *ListOptions) {
			opt.file.Ref = pointer.ToString("develop")
			opt.file.ListOptions.Page = 1
			opt.file.ListOptions.PerPage = 100
		},
		wantError: nil,
	}, {
		name: "desc sort",
		args: []string{},
		optionsFunc: func(opt *ListOptions) {
			opt.file.Sort = "desc"
		},
		wantError: nil,
	}}
	ioStreams := cmdutil.NewTestIOStreamsDiscard()
	factory := cmdutil.NewFactory(cmdtesting.NewFakeRESTClientGetter())
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmd := NewGetFilesCmd(factory, ioStreams)
			var cmdOptions = NewListOptions(ioStreams)
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
