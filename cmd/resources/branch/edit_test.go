package branch

import (
	"errors"
	"fmt"
	cmdtesting "github.com/huhouhua/glctl/cmd/testing"
	cmdutil "github.com/huhouhua/glctl/cmd/util"
	"github.com/huhouhua/glctl/util/cli"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
	"strings"
	"testing"
)

func TestEditBranch(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		optionsFunc func(opt *EditOptions)
		validate    func(opt *EditOptions, cmd *cobra.Command, args []string) error
		run         func(opt *EditOptions, args []string) error
		wantError   error
	}{{
		name: "edit by name",
		args: []string{"develop"},
		optionsFunc: func(opt *EditOptions) {
			opt.project = "huhouhua/gitlab-repo-branch"
		},
		run: func(opt *EditOptions, args []string) error {
			var err error
			out := cmdtesting.Run(func() {
				err = opt.Run(args)
			})
			expectedOutput := fmt.Sprintf("Branch (%s) from project (%s) has been deleted", opt.branch, opt.project)
			if !strings.Contains(out, expectedOutput) {
				err = errors.New(fmt.Sprintf("delete by path : Unexpected output! Expected\n%s\ngot\n%s", expectedOutput, out))
			}
			return err
		},
		wantError: nil,
	}, {
		name: "change can push",
		args: []string{"not-found"},
		optionsFunc: func(opt *EditOptions) {
			opt.project = "huhouhua/gitlab-repo-branch"
		},
		run: func(opt *EditOptions, args []string) error {
			err := opt.Run(args)
			var repo *gitlab.ErrorResponse
			if errors.As(err, &repo) && repo.Message == "{message: 404 Branch Not Found}" {
				return nil
			}
			return err
		},
	}, {
		name: "set unprotect",
		args: []string{"not-found"},
		optionsFunc: func(opt *EditOptions) {
			opt.project = "not-found"
		},
		run: func(opt *EditOptions, args []string) error {
			err := opt.Run(args)
			var repo *gitlab.ErrorResponse
			if errors.As(err, &repo) && repo.Message == "{message: 404 Project Not Found}" {
				return nil
			}
			return err
		},
	}, {
		name: "not definition branch",
		args: []string{},
		validate: func(opt *EditOptions, cmd *cobra.Command, args []string) error {
			err := opt.Validate(cmd, args)
			if err.Error() == "please enter branch" {
				return err
			}
			return nil
		},
	}}
	streams := cli.NewTestIOStreamsDiscard()
	factory := cmdutil.NewFactory(cmdtesting.NewFakeRESTClientGetter())
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmd := NewEditBranchCmd(factory, streams)
			var cmdOptions = NewEditOptions(streams)
			if tc.optionsFunc != nil {
				tc.optionsFunc(cmdOptions)
			}
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
