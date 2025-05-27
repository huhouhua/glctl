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
	"errors"
	"fmt"
	"github.com/huhouhua/glctl/pkg/cli/genericiooptions"
	"strconv"
	"testing"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/rand"

	cmdutil "github.com/huhouhua/glctl/cmd/util"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"

	cmdtesting "github.com/huhouhua/glctl/cmd/testing"
)

func TestDeleteProject(t *testing.T) {
	factory := cmdutil.NewFactory(cmdtesting.NewFakeRESTClientGetter())
	client, err := factory.GitlabClient()
	if err != nil {
		t.Errorf("create gitlab client fail! %v", err)
		return
	}
	tests := []struct {
		name      string
		args      func() []string
		validate  func(opt *DeleteOptions, cmd *cobra.Command, args []string) error
		run       func(opt *DeleteOptions, args []string) error
		wantError error
	}{{
		name: "delete by path",
		args: func() []string {
			return []string{"Group1/SubGroup1/Project1"}
		},
		run: func(opt *DeleteOptions, args []string) error {
			var err error
			_, err = forkProject("Group1/Project1", opt.gitlabClient)
			if err != nil {
				return err
			}
			defer func() {
				_, _ = opt.gitlabClient.Projects.DeleteProject(args[0])
			}()
			out := cmdtesting.RunForStdout(opt.ioStreams, func() {
				err = opt.Run(args)
			})
			expectedOutput := fmt.Sprintf("project (%s) with id", args[0])
			assert.Containsf(
				t,
				out,
				expectedOutput,
				"delete by path : Unexpected output! Expected\n%s\ngot\n%s",
				expectedOutput,
				out,
			)
			return err
		},
		wantError: nil,
	}, {
		name: "delete by id",
		args: func() []string {
			project, err := forkProject("Group1/Project2", client)
			if err != nil {
				t.Errorf("fork project fail")
				return nil
			}
			return []string{strconv.Itoa(project.ID)}
		},
		run: func(opt *DeleteOptions, args []string) error {
			defer func() {
				_, _ = opt.gitlabClient.Projects.DeleteProject(args[0])
			}()
			var err error
			out := cmdtesting.RunForStdout(opt.ioStreams, func() {
				err = opt.Run(args)
			})
			expectedOutput := fmt.Sprintf("with id (%s) has been deleted", args[0])
			assert.Containsf(
				t,
				out,
				expectedOutput,
				"delete by id : Unexpected output! Expected\n%s\ngot\n%s",
				expectedOutput,
				out,
			)
			return err
		},
		wantError: nil,
	}, {
		name: "example delete a nonexistent ID",
		args: func() []string {
			return []string{"100001"}
		},
		run: func(opt *DeleteOptions, args []string) error {
			err := opt.Run(args)
			var repoErr *gitlab.ErrorResponse
			assert.ErrorAs(t, err, &repoErr)
			if assert.Equal(t, repoErr.Message, "{message: 404 Project Not Found}") {
				return nil
			}
			return err
		},
	}, {
		name:      "no id",
		wantError: errors.New("please enter project name or id"),
	}}
	for _, tc := range tests {
		streams := genericiooptions.NewTestIOStreamsForPipe()
		t.Run(tc.name, func(t *testing.T) {
			cmd := NewDeleteProjectCmd(factory, streams)
			cmdOptions := NewDeleteOptions(streams)
			var (
				err  error
				args []string
			)
			if tc.args != nil {
				args = tc.args()
			}
			err = cmdOptions.Complete(factory, cmd, args)
			cmdtesting.ErrorAssertionWithEqual(t, tc.wantError, err)
			if err != nil {
				return
			}
			if tc.validate != nil {
				err = tc.validate(cmdOptions, cmd, args)
			} else {
				err = cmdOptions.Validate(cmd, args)
			}
			cmdtesting.ErrorAssertionWithEqual(t, tc.wantError, err)
			if err != nil {
				return
			}
			if tc.run != nil {
				err = tc.run(cmdOptions, args)
			} else {
				err = cmdOptions.Run(args)
			}
			cmdtesting.ErrorAssertionWithEqual(t, tc.wantError, err)
		})
	}
}

func forkProject(forkProject string, client *gitlab.Client) (*gitlab.Project, error) {
	rand.Seed(uint64(time.Now().UnixNano()))
	r := rand.Intn(1000)
	p, _, err := client.Projects.ForkProject(forkProject, &gitlab.ForkProjectOptions{
		NamespacePath: pointer.ToString("Group1/SubGroup1"),
		Name:          pointer.ToString(fmt.Sprintf("%s-%b", "glctl-from-delete", r)),
		Visibility:    pointer.To(gitlab.PublicVisibility),
	})
	return p, err
}
