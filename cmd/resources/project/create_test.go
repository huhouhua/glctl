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
	"fmt"
	"strings"
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"

	cmdtesting "github.com/huhouhua/glctl/cmd/testing"
	cmdutil "github.com/huhouhua/glctl/cmd/util"
	"github.com/huhouhua/glctl/util/cli"
)

func TestCreateProject(t *testing.T) {
	streams := cli.NewTestIOStreamsForPipe()
	tests := []struct {
		name        string
		args        []string
		optionsFunc func(opt *CreateOptions)
		validate    func(opt *CreateOptions, cmd *cobra.Command, args []string) error
		run         func(opt *CreateOptions, args []string) error
		wantError   error
	}{{
		name: "create a new project using all flags",
		args: []string{"gitlab-repo-test"},
		optionsFunc: func(opt *CreateOptions) {
			opt.namespace = "test-project"
			opt.project.Description = pointer.ToString("Created by go test")
			opt.project.IssuesAccessLevel = pointer.To(gitlab.EnabledAccessControl)
			opt.project.MergeRequestsAccessLevel = pointer.To(gitlab.EnabledAccessControl)
			opt.project.BuildsAccessLevel = pointer.To(gitlab.EnabledAccessControl)
			opt.project.WikiAccessLevel = pointer.To(gitlab.EnabledAccessControl)
			opt.project.SnippetsAccessLevel = pointer.To(gitlab.EnabledAccessControl)
			opt.project.ResolveOutdatedDiffDiscussions = pointer.ToBool(true)
			opt.project.ContainerRegistryAccessLevel = pointer.To(gitlab.EnabledAccessControl)
			opt.project.SharedRunnersEnabled = pointer.ToBool(true)
			opt.project.Visibility = pointer.To(gitlab.PrivateVisibility)
			opt.project.PublicBuilds = pointer.ToBool(true)
			opt.project.OnlyAllowMergeIfPipelineSucceeds = pointer.ToBool(true)
			opt.project.OnlyAllowMergeIfAllDiscussionsAreResolved = pointer.ToBool(true)
			opt.project.MergeMethod = pointer.To(gitlab.RebaseMerge)
			opt.project.LFSEnabled = pointer.ToBool(true)
			opt.project.RequestAccessEnabled = pointer.ToBool(true)
			opt.project.Topics = pointer.To([]string{"gotest", "tdd"})
			opt.project.PrintingMergeRequestLinkEnabled = pointer.ToBool(true)
			opt.project.CIConfigPath = pointer.ToString("gitlabci.yml")
		},
		run: func(opt *CreateOptions, args []string) error {
			var err error
			out := cmdtesting.RunForStdout(streams, func() {
				err = opt.Run(args)
			})
			expectedOutput := fmt.Sprintf("%s/%s.git", opt.namespace, args[0])
			if !strings.Contains(out, expectedOutput) {
				err = errors.New(
					fmt.Sprintf("create a new project : Unexpected output! Expected\n%s\ngot\n%s", expectedOutput, out),
				)
			}
			return err
		},
		wantError: nil,
	}}
	factory := cmdutil.NewFactory(cmdtesting.NewFakeRESTClientGetter())
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmd := NewCreateProjectCmd(factory, streams)
			cmdOptions := NewCreateOptions(streams)
			var err error
			if tc.optionsFunc != nil {
				tc.optionsFunc(cmdOptions)
			}
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
