package project

import (
	"fmt"
	cmdtesting "github.com/huhouhua/gitlab-repo-operator/cmd/testing"
	cmdutil "github.com/huhouhua/gitlab-repo-operator/cmd/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
	"strings"
	"testing"
)

func TestCreateProject(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		validate  func(opt *CreateOptions, cmd *cobra.Command, args []string) error
		run       func(opt *CreateOptions, args []string) error
		wantError error
	}{{
		name: "delete by path",
		args: []string{"huhouhua/gitlab-repo-test"},
		run: func(opt *CreateOptions, args []string) error {
			var err error
			out := cmdtesting.RunTestForStdout(func() {
				err = opt.Run(args)
			})
			expectedOutput := fmt.Sprintf("project (%s) with id", "huhouhua/gitlab-repo-test")
			if !strings.Contains(out, expectedOutput) {
				err = errors.New(fmt.Sprintf("delete by path : Unexpected output! Expected\n%s\ngot\n%s", expectedOutput, out))
			}
			return err
		},
		wantError: nil,
	}, {
		name: "delete by id",
		args: []string{"222"},
		run: func(opt *CreateOptions, args []string) error {
			var err error
			out := cmdtesting.RunTestForStdout(func() {
				err = opt.Run(args)
			})
			expectedOutput := fmt.Sprintf("with id (%d) has been deleted", 222)
			if !strings.Contains(out, expectedOutput) {
				err = errors.New(fmt.Sprintf("delete by path : Unexpected output! Expected\n%s\ngot\n%s", expectedOutput, out))
			}
			return err
		},
		wantError: nil,
	}, {
		name: "example delete a nonexistent ID",
		args: []string{"100001"},
		run: func(opt *CreateOptions, args []string) error {
			err := opt.Run(args)
			var repo *gitlab.ErrorResponse
			if errors.As(err, &repo) && repo.Message == "{message: 404 Project Not Found}" {
				return nil
			}
			return err
		},
	}, {
		name: "no id",
		args: []string{},
		validate: func(opt *CreateOptions, cmd *cobra.Command, args []string) error {
			err := opt.Validate(cmd, args)
			if err.Error() == "please enter project and id" {
				return err
			}
			return nil
		},
	}}
	factory := cmdutil.NewFactory(cmdtesting.NewFakeRESTClientGetter())
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmd := NewCreateProjectCmd(factory)
			cmdOptions := NewCreateOptions()
			var err error
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
