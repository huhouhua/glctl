package cmd

import (
	"github.com/spf13/cobra"
	"io"
)

var globalUsage = `the gitlab repository operator

Common actions for gitrepo

- gitrepo insert:   start insert file
`

func NewRootCmd(out io.Writer) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:           "gitrepo",
		Short:         "the gitlab repository operator",
		Long:          globalUsage,
		SilenceErrors: true,
	}
	//flags := cmd.PersistentFlags()
	cmd.AddCommand(newInsertCmd(out))
	return cmd, nil
}
