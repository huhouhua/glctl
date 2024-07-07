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
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/AlekSi/pointer"
	"github.com/huhouhua/gl/cmd/require"
	cmdutil "github.com/huhouhua/gl/cmd/util"
	"github.com/huhouhua/gl/cmd/util/editor"
	"github.com/huhouhua/gl/cmd/util/editor/crlf"
	"github.com/huhouhua/gl/util/cli"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type EditOptions struct {
	gitlabClient       *gitlab.Client
	path               string
	Project            string
	file               *gitlab.GetFileOptions
	WindowsLineEndings bool
	ioStreams          cli.IOStreams
}

func NewEditOptions(ioStreams cli.IOStreams) *EditOptions {
	return &EditOptions{
		ioStreams: ioStreams,
		file: &gitlab.GetFileOptions{
			Ref: pointer.ToString("main"),
		},
		WindowsLineEndings: runtime.GOOS == "windows",
	}
}

var (
	editFileDesc = "edit file for project "

	editFileExample = `# edit file
gl edit files myfile -p project`
)

func NewEditFileCmd(f cmdutil.Factory, ioStreams cli.IOStreams) *cobra.Command {
	o := NewEditOptions(ioStreams)
	cmd := &cobra.Command{
		Use:                   "file",
		Aliases:               []string{"f"},
		Short:                 editFileDesc,
		Example:               editFileExample,
		Args:                  require.MinimumNArgs(1),
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

func (o *EditOptions) AddFlags(cmd *cobra.Command) {
	cmdutil.AddProjectVarPFlag(cmd, &o.Project)
	f := cmd.Flags()
	f.StringVar(o.file.Ref, "ref", *o.file.Ref, "The name of a repository branch or tag or, if not given, the default branch.")
	cmdutil.VerifyMarkFlagRequired(cmd, "project")
}

// Complete completes all the required options.
func (o *EditOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	var err error
	if len(args) > 0 {
		o.path = args[0]
	}
	o.gitlabClient, err = f.GitlabClient()
	return err
}

// Validate makes sure there is no discrepency in command options.
func (o *EditOptions) Validate(cmd *cobra.Command, args []string) error {
	if strings.TrimSpace(o.Project) == "" {
		_ = cmd.Usage()
		return fmt.Errorf("please enter project name and id")
	}
	return nil
}

// Run executes a list subcommand using the specified options.
func (o *EditOptions) Run(args []string) error {
	repoFile, _, err := o.gitlabClient.RepositoryFiles.GetFile(o.Project, o.path, o.file)
	if err != nil {
		return err
	}
	edit := editor.NewDefaultEditor([]string{"EDITOR"})
	// generate the file to edit
	decodeString, err := base64.StdEncoding.DecodeString(repoFile.Content)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(decodeString)
	var w io.Writer = buf
	if o.WindowsLineEndings {
		w = crlf.NewCRLFWriter(w)
	}
	edited, file, err := edit.LaunchTempFile(fmt.Sprintf("%s-edit-", filepath.Base(os.Args[0])), filepath.Ext(repoFile.FileName), buf)
	if err != nil {
		return err
	}
	fmt.Print(edited)
	defer func() {
		// cleanup any file from the previous pass
		if len(file) > 0 {
			_ = os.Remove(file)
		}
	}()

	return nil
}
