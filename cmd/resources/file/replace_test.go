package file

import (
	"errors"
	cmdtesting "github.com/huhouhua/gl/cmd/testing"
	cmdutil "github.com/huhouhua/gl/cmd/util"
	"github.com/huhouhua/gl/util/cli"
	"github.com/spf13/cobra"
	"testing"
)

func TestRunReplace(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		optionsFunc func(opt *ReplaceOptions)
		validate    func(opt *ReplaceOptions, cmd *cobra.Command, args []string) error
		run         func(opt *ReplaceOptions, args []string) error
		wantError   error
	}{{
		name: "replace all branch",
		args: []string{"resources/pipeline/buildserver/package.yaml"},
		optionsFunc: func(opt *ReplaceOptions) {
			opt.Project = "224"
			opt.RefMatch = "*"
			opt.FileName = "../../../testdata/replace/new_test.yaml"
		},
		run: func(opt *ReplaceOptions, args []string) error {
			var err error
			_ = cmdtesting.RunTest(func() {
				err = opt.Run(args)
			})
			//expectedOutput := fmt.Sprintf("%s edited", opt.path)
			//if !strings.Contains(out, expectedOutput) {
			//	err = errors.New(fmt.Sprintf("compare content : Unexpected output! Expected\n%s\ngot\n%s", expectedOutput, out))
			//}
			return err
		},
		wantError: nil,
	}, {
		name: "force replace",
		args: []string{"package.yaml"},
		optionsFunc: func(opt *ReplaceOptions) {
			opt.Project = "224"
			opt.RefMatch = "*"
			opt.FileName = "../../../testdata/replace/new_test.yaml"
			opt.Force = true
		},
		run: func(opt *ReplaceOptions, args []string) error {
			var err error
			_ = cmdtesting.RunTest(func() {
				err = opt.Run(args)
			})
			//expectedOutput := fmt.Sprintf("%s edited", opt.path)
			//if !strings.Contains(out, expectedOutput) {
			//	err = errors.New(fmt.Sprintf("compare content : Unexpected output! Expected\n%s\ngot\n%s", expectedOutput, out))
			//}
			return err
		},
		wantError: nil,
	}}
	streams := cli.NewTestIOStreamsDiscard()
	factory := cmdutil.NewFactory(cmdtesting.NewFakeRESTClientGetter())
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmd := NewReplaceFileCmd(factory, streams)
			var cmdOptions = NewReplaceOptions(streams)
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
