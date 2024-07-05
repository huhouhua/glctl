// Copyright 2024 The huhouhua Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http:www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package file

import (
	"fmt"
	"github.com/AlekSi/pointer"
	cmdutil "github.com/huhouhua/gl/cmd/util"
	"github.com/huhouhua/gl/util/cli"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
	"strings"
)

type DeleteOptions struct {
	gitlabClient *gitlab.Client
	file         *gitlab.DeleteFileOptions
	project      string
	FileName     string
	ioStreams    cli.IOStreams
}

func NewDeleteOptions(ioStreams cli.IOStreams) *DeleteOptions {
	return &DeleteOptions{
		ioStreams: ioStreams,
		file: &gitlab.DeleteFileOptions{
			Branch:        pointer.ToString("main"),
			CommitMessage: pointer.ToString(""),
			AuthorName:    pointer.ToString(""),
			AuthorEmail:   pointer.ToString(""),
			LastCommitID:  pointer.ToString(""),
			StartBranch:   pointer.ToString(""),
		},
	}
}

var (
	deleteFilesDesc = "delete file for project "

	deleteFilesExample = `# delete project file
delete get files myProject`
)

func NewDeleteFilesCmd(f cmdutil.Factory, ioStreams cli.IOStreams) *cobra.Command {
	o := NewDeleteOptions(ioStreams)
	cmd := &cobra.Command{
		Use:                   "files",
		Aliases:               []string{"f"},
		Short:                 deleteFilesDesc,
		Example:               deleteFilesExample,
		DisableFlagsInUseLine: true,
		TraverseChildren:      true,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Validate(cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{"file"},
	}
	o.AddFlags(cmd)
	return cmd
}
func (o *DeleteOptions) AddFlags(cmd *cobra.Command) {
	cmdutil.AddProjectVarPFlag(cmd, &o.project)
	f := cmd.Flags()
	f.StringVarP(o.file.Branch, "branch", "b", *o.file.Branch, "Name of the new branch to create. The commit is added to this branch.(default main)")
	f.StringVarP(o.file.CommitMessage, "message", "m", *o.file.CommitMessage, "The commit message.(default delete file_path)")
	f.StringVar(o.file.AuthorEmail, "author_email", *o.file.AuthorEmail, "The commit author’s email address.")
	f.StringVar(o.file.AuthorName, "author_name", *o.file.AuthorName, "The commit author’s name.")
	f.StringVar(o.file.LastCommitID, "last_commit_id", *o.file.LastCommitID, "Last known file commit ID.")
	f.StringVar(o.file.StartBranch, "start_branch", *o.file.StartBranch, "Name of the base branch to create the new branch from.")
}

// Complete completes all the required options.
func (o *DeleteOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	var err error
	if len(args) > 0 {
		o.FileName = args[0]
	}
	if strings.TrimSpace(*o.file.CommitMessage) == "" {
		o.file.CommitMessage = pointer.ToString(fmt.Sprintf("delete %s file!", o.FileName))
	}
	o.gitlabClient, err = f.GitlabClient()
	return err
}

// Validate makes sure there is no discrepency in command options.
func (o *DeleteOptions) Validate(cmd *cobra.Command, args []string) error {
	if len(args) == 0 || strings.TrimSpace(o.FileName) == "" {
		return fmt.Errorf("please enter file name")
	}
	if strings.TrimSpace(o.project) == "" || strings.TrimSpace(*o.file.Branch) == "" {
		return cmd.Usage()
	}
	return nil
}

// Run executes a list subcommand using the specified options.
func (o *DeleteOptions) Run(args []string) error {
	//o.gitlabClient.RepositoryFiles.GetFileMetaData()
	_, err := o.gitlabClient.RepositoryFiles.DeleteFile(o.project, o.FileName, o.file)
	if err != nil {
		return err
	}
	_, _ = fmt.Fprintf(o.ioStreams.Out, "file (%s) for %s branch with project id (%s) has been deleted\n", o.FileName, *o.file.Branch, o.project)
	return nil
}
