package branch

import (
	"errors"
	"fmt"
	"github.com/AlekSi/pointer"
	cmdtesting "github.com/huhouhua/glctl/cmd/testing"
	cmdutil "github.com/huhouhua/glctl/cmd/util"
	"github.com/huhouhua/glctl/util/cli"
	"github.com/spf13/cobra"
	"strings"
	"testing"
)

func TestEditBranch(t *testing.T) {
	streams := cli.NewTestIOStreamsForPipe()
	tests := []struct {
		name        string
		args        []string
		optionsFunc func(opt *EditOptions)
		validate    func(opt *EditOptions, cmd *cobra.Command, args []string) error
		run         func(opt *EditOptions, args []string) error
		wantError   error
	}{{
		name: "set unprotect",
		args: []string{"main"},
		optionsFunc: func(opt *EditOptions) {
			opt.project = "huhouhua/gitlab-repo-branch"
			opt.Unprotect = true
		},
		run: func(opt *EditOptions, args []string) error {
			var err error
			out := cmdtesting.RunForStdout(streams, func() {
				err = opt.Run(args)
			})
			expectedOutput := fmt.Sprintf("branch %s un protect", args[0])
			if !strings.Contains(out, expectedOutput) {
				err = errors.New(fmt.Sprintf("unprotect main : Unexpected output! Expected\n%s\ngot\n%s", expectedOutput, out))
			}
			return err
		},
		wantError: nil,
	}, {
		name: "set protect",
		args: []string{"main"},
		optionsFunc: func(opt *EditOptions) {
			opt.project = "huhouhua/gitlab-repo-branch"
			opt.protect = true
			opt.protectBranch.DevelopersCanMerge = pointer.ToBool(true)
			opt.protectBranch.DevelopersCanPush = pointer.ToBool(true)
		},
		run: func(opt *EditOptions, args []string) error {
			var err error
			out := cmdtesting.RunForStdout(streams, func() {
				err = opt.Run(args)
			})
			expectedOutput := fmt.Sprintf("branch %s updated", args[0])
			if !strings.Contains(out, expectedOutput) {
				err = errors.New(fmt.Sprintf("unprotect main : Unexpected output! Expected\n%s\ngot\n%s", expectedOutput, out))
			}
			return err
		},
		wantError: nil,
	}}
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
