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
	"fmt"
	"testing"

	"github.com/huhouhua/glctl/pkg/cli/genericiooptions"

	"github.com/stretchr/testify/assert"
	gitlab "gitlab.com/gitlab-org/api/client-go"

	"github.com/AlekSi/pointer"
	"github.com/spf13/cobra"

	cmdtesting "github.com/huhouhua/glctl/cmd/testing"
	cmdutil "github.com/huhouhua/glctl/cmd/util"
)

func TestDeleteFile(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		optionsFunc func(opt *DeleteOptions)
		validate    func(opt *DeleteOptions, cmd *cobra.Command, args []string) error
		run         func(opt *DeleteOptions, args []string) error
		wantError   error
	}{{
		name: "delete a file",
		args: []string{"delete.yaml"},
		optionsFunc: func(opt *DeleteOptions) {
			opt.Project = "Group2/SubGroup3/Project15"
			opt.file.Branch = pointer.ToString("main")
		},
		run: func(opt *DeleteOptions, args []string) error {
			err := createTestFile(opt.gitlabClient, opt.FileName, opt)
			if err != nil {
				return err
			}
			defer func() {
				clearTestFile(opt.gitlabClient, opt)
			}()
			out := cmdtesting.RunForStdout(opt.ioStreams, func() {
				err = opt.Run(args)
			})
			expectedOutput := fmt.Sprintf(
				"file (%s) for %s branch with project id (%s) has been deleted",
				opt.FileName,
				*opt.file.Branch,
				opt.Project,
			)
			assert.Containsf(
				t,
				out,
				expectedOutput,
				"delete a file: Unexpected output! Expected\n%s\ngot\n%s",
				expectedOutput,
				out,
			)
			return err
		},
		wantError: nil,
	}, {
		// to do
		name: "delete dir",
		args: []string{"/glctl-dir"},
		optionsFunc: func(opt *DeleteOptions) {
			opt.Project = "Group2/SubGroup3/Project14"
			opt.file.Branch = pointer.ToString("main")
		},
		run: func(opt *DeleteOptions, args []string) error {
			// err := createTestFile(opt.gitlabClient, "glctl-dir/glctl.yaml", opt)
			// if err != nil {
			//	return err
			// }
			// defer func() {
			//	clearTestFile(opt.gitlabClient, opt)
			// }()
			// out := cmdtesting.RunForStdout(opt.ioStreams, func() {
			//	err = opt.Run(args)
			// })
			// expectedOutput := fmt.Sprintf(
			//	"file (%s) for %s branch with project id (%s) has been deleted",
			//	opt.FileName,
			//	*opt.file.Branch,
			//	opt.Project,
			// )
			// if !strings.Contains(out, expectedOutput) {
			//	err = errors.New(
			//		fmt.Sprintf("delete dir : Unexpected output! Expected\n%s\ngot\n%s", expectedOutput, out),
			//	)
			//}
			// return err
			return nil
		},
		wantError: nil,
	}}
	factory := cmdutil.NewFactory(cmdtesting.NewFakeRESTClientGetter())
	for _, tc := range tests {
		streams := genericiooptions.NewTestIOStreamsForPipe()
		t.Run(tc.name, func(t *testing.T) {
			cmd := NewDeleteFilesCmd(factory, streams)
			var cmdOptions = NewDeleteOptions(streams)
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

func createTestFile(client *gitlab.Client, fileName string, opt *DeleteOptions) error {
	_, _, err := client.RepositoryFiles.CreateFile(opt.Project, fileName, &gitlab.CreateFileOptions{
		Branch:        opt.file.Branch,
		Content:       pointer.ToString("delete: true"),
		CommitMessage: pointer.ToString("test delete file"),
	})
	return err
}
func clearTestFile(client *gitlab.Client, opt *DeleteOptions) {
	_, _ = opt.gitlabClient.RepositoryFiles.DeleteFile(opt.Project, opt.FileName, &gitlab.DeleteFileOptions{
		Branch:        opt.file.Branch,
		CommitMessage: pointer.ToString("clear test file"),
	})
}
