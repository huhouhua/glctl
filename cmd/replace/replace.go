package replace

import (
	"github.com/huhouhua/gl/cmd/resources/file"
	cmdutil "github.com/huhouhua/gl/cmd/util"
	"github.com/huhouhua/gl/util/cli"
	"github.com/huhouhua/gl/util/templates"
	"github.com/spf13/cobra"
)

var (
	replaceLong = templates.LongDesc(`
		Replace a repository file by file name or stdin.

		all formats are accepted. If replacing an existing repository file, the
		complete repository file spec must be provided. This can be obtained by

		    $ gl get files PROJECT --path=my.yml --ref=BRANCH --raw`)

	replaceExample = templates.Examples(`
		# Replace a single file using the data in my.yml
		gl replace file app/my.yml -f ./my.yml --ref=main --project=myproject 

		# Replace all branch file using the data in my.yml
	    gl replace file app/my.yml -f ./my.yml --ref-match=* --project=myproject`)
)

func NewReplaceCmd(f cmdutil.Factory, ioStreams cli.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "replace",
		Aliases:               []string{"r"},
		Short:                 replaceLong,
		Example:               replaceExample,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(file.NewReplaceFileCmd(f, ioStreams))
	return cmd
}
