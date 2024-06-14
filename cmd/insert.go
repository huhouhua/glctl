package cmd

import (
	"github.com/spf13/cobra"
	"io"
)

const createDesc = `dada`

func newInsertCmd(out io.Writer) *cobra.Command {

	cmd := &cobra.Command{
		Use: "insert ",
	}
	return cmd
}
