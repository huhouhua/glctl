package cmd

import (
	"bytes"
	cmdutil "github.com/huhouhua/gitlab-repo-operator/cmd/testing"
	"github.com/spf13/cobra"
	"testing"
)

func TestRoot(t *testing.T) {
	tests := []struct {
		name      string
		cmd       string
		wantError bool
	}{{
		name:      "list all projects",
		cmd:       "get projects",
		wantError: false,
	}}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out, err := cmdutil.ExecuteCommand(func(buffer *bytes.Buffer) (*cobra.Command, error) {
				return NewRootCmd(buffer)
			}, tc.cmd)

			cmdutil.TInfo(out)
			if tc.wantError && err == nil {
				t.Errorf("expected error, got success with the following output:\n%s", out)
			}
			if !tc.wantError && err != nil {
				t.Errorf("expected no error, got: '%v'", err)
			}
		})
	}
}
