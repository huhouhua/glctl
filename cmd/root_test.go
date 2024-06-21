package cmd

import (
	"bytes"
	cmdutil "github.com/huhouhua/gitlab-repo-operator/cmd/util"
	"github.com/spf13/cobra"
	"testing"
)

func TestRoot(t *testing.T) {
	tests := []cmdutil.CmdTestCase{{
		Name:      "list all projects",
		Cmd:       "get projects",
		WantError: false,
	}}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			out, err := cmdutil.ExecuteCommand(func(buffer *bytes.Buffer) (*cobra.Command, error) {
				return NewRootCmd(buffer)
			}, tc.Cmd)

			cmdutil.TInfo(out)
			if tc.WantError && err == nil {
				t.Errorf("expected error, got success with the following output:\n%s", out)
			}
			if !tc.WantError && err != nil {
				t.Errorf("expected no error, got: '%v'", err)
			}
		})
	}
}
