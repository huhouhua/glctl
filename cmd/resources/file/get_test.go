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

package file

import (
	"github.com/huhouhua/glctl/pkg/cli/genericiooptions"
	"testing"

	"github.com/AlekSi/pointer"

	cmdtesting "github.com/huhouhua/glctl/cmd/testing"
	cmdutil "github.com/huhouhua/glctl/cmd/util"
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
		args: []string{"Group2/SubGroup3/Project13"},
		optionsFunc: func(opt *ListOptions) {
			opt.All = true
		},
		wantError: nil,
	}, {
		name: "get specified branch",
		args: []string{
			"Group2/SubGroup3/Project14",
		},
		optionsFunc: func(opt *ListOptions) {
			opt.file.Ref = pointer.ToString("main")
		},
		wantError: nil,
	}, {
		name: "read to receive the raw file in repository",
		args: []string{
			"Group2/SubGroup3/Project15",
		},
		optionsFunc: func(opt *ListOptions) {
			opt.file.Ref = pointer.ToString("main")
			opt.file.Path = pointer.ToString("test/test.yaml")
			opt.Raw = true
		},
		wantError: nil,
	}, {
		name: "Get specified directory",
		args: []string{"Group2/SubGroup3/Project15"},
		optionsFunc: func(opt *ListOptions) {
			opt.file.Ref = pointer.ToString("main")
			opt.file.Path = pointer.ToString("test")
		},
		wantError: nil,
	}, {
		name: "list all projects with page",
		args: []string{"Group2/SubGroup3/Project15"},
		optionsFunc: func(opt *ListOptions) {
			opt.file.Ref = pointer.ToString("main")
			opt.file.ListOptions.Page = 1
			opt.file.ListOptions.PerPage = 100
		},
		wantError: nil,
	}, {
		name: "desc sort",
		args: []string{"Group2/SubGroup3/Project13"},
		optionsFunc: func(opt *ListOptions) {
			opt.file.Sort = "desc"
		},
		wantError: nil,
	}}
	streams := genericiooptions.NewTestIOStreamsDiscard()
	factory := cmdutil.NewFactory(cmdtesting.NewFakeRESTClientGetter())
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmd := NewGetFilesCmd(factory, streams)
			var cmdOptions = NewListOptions(streams)
			if tc.optionsFunc != nil {
				tc.optionsFunc(cmdOptions)
			}
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
			err = cmdOptions.Run(tc.args)
			cmdtesting.ErrorAssertionWithEqual(t, tc.wantError, err)
			if err != nil {
				return
			}
		})
	}
}
