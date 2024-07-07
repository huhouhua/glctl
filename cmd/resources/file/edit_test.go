package file

import (
	"errors"
	"fmt"
	"github.com/AlekSi/pointer"
	cmdtesting "github.com/huhouhua/gl/cmd/testing"
	cmdutil "github.com/huhouhua/gl/cmd/util"
	"github.com/huhouhua/gl/util/cli"
	"github.com/spf13/cobra"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		optionsFunc func(opt *EditOptions)
		validate    func(opt *EditOptions, cmd *cobra.Command, args []string) error
		run         func(opt *EditOptions, args []string) error
		wantError   error
	}{{
		name: "compare content",
		args: []string{"dada.yml"},
		optionsFunc: func(opt *EditOptions) {
			opt.Project = "216"
			opt.file.Ref = pointer.ToString("test-ci")
		},
		run: func(opt *EditOptions, args []string) error {
			var err error
			out := cmdtesting.RunTestForStdout(func() {
				err = opt.Run(args)
			})
			expectedOutput := fmt.Sprintf("%s edited", opt.path)
			if !strings.Contains(out, expectedOutput) {
				err = errors.New(fmt.Sprintf("compare content : Unexpected output! Expected\n%s\ngot\n%s", expectedOutput, out))
			}
			return err
		},
		wantError: nil,
	}}
	streams := cli.NewTestIOStreamsDiscard()
	factory := cmdutil.NewFactory(cmdtesting.NewFakeRESTClientGetter())
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmd := NewEditFileCmd(factory, streams)
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
