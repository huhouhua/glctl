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
	"fmt"
	"testing"

	"github.com/huhouhua/glctl/pkg/cli/genericiooptions"

	"github.com/stretchr/testify/assert"

	"github.com/AlekSi/pointer"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"

	cmdtesting "github.com/huhouhua/glctl/cmd/testing"
	cmdutil "github.com/huhouhua/glctl/cmd/util"
)

func TestCreateProject(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		optionsFunc func(opt *CreateOptions)
		validate    func(opt *CreateOptions, cmd *cobra.Command, args []string) error
		run         func(opt *CreateOptions, args []string) error
		wantError   error
	}{{
		name: "create a new project using all flags",
		args: []string{"glctl-from-create"},
		optionsFunc: func(opt *CreateOptions) {
			opt.namespace = "Group1"
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
			projectPath := fmt.Sprintf("%s/%s", opt.namespace, args[0])
			defer func() {
				_, _ = opt.gitlabClient.Projects.DeleteProject(projectPath)
			}()
			out := cmdtesting.RunForStdout(opt.ioStreams, func() {
				err = opt.Run(args)
			})
			expectedOutput := fmt.Sprintf("%s.git", projectPath)
			assert.Containsf(
				t,
				out,
				expectedOutput,
				"create a new project: Unexpected output! Expected\n%s\ngot\n%s",
				expectedOutput,
				out,
			)
			return err
		},
		wantError: nil,
	}}
	factory := cmdutil.NewFactory(cmdtesting.NewFakeRESTClientGetter())
	streams := genericiooptions.NewTestIOStreamsForPipe()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmd := NewCreateProjectCmd(factory, streams)
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
			} else {
				err = cmdOptions.Run(tc.args)
			}
			cmdtesting.ErrorAssertionWithEqual(t, tc.wantError, err)
		})
	}
}
