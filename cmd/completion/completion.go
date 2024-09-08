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

package completion

import (
	"fmt"
	"github.com/huhouhua/glctl/pkg/cli/genericiooptions"
	"github.com/huhouhua/glctl/pkg/util/templates"
	"io"

	"github.com/spf13/cobra"

	cmdutil "github.com/huhouhua/glctl/cmd/util"
)

const defaultBoilerPlate = `
# Copyright 2024 The huhouhua Authors
# 
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
# 
# http:www.apache.org/licenses/LICENSE-2.0
# 
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
`

var (
	completionLong = templates.LongDesc(`
		Output shell completion code for the specified shell (bash, zsh, fish, or powershell).
		The shell code must be evaluated to provide interactive
		completion of glctl commands.  This can be done by sourcing it from
		the .bash_profile.

		Detailed instructions on how to do this are available here:

		Note for zsh users: [1] zsh completions are only supported in versions of zsh >= 5.2`)

	completionExample = templates.Examples(`
		# Installing bash completion on macOS using homebrew
		## If running Bash 3.2 included with macOS
		    brew install bash-completion
		## or, if running Bash 4.1+
		    brew install bash-completion@2
		## If glctl is installed via homebrew, this should start working immediately.
		## If you've installed via other means, you may need add the completion to your completion directory
		    glctl completion bash > $(brew --prefix)/etc/bash_completion.d/glctl


		# Installing bash completion on Linux
		## If bash-completion is not installed on Linux, please install the 'bash-completion' package
		## via your distribution's package manager.
		## Load the glctl completion code for bash into the current shell
		    source <(glctl completion bash)
		## Write bash completion code to a file and source if from .bash_profile
		    glctl completion bash > ~/.gl/completion.bash.inc
		    printf "
		      # gl shell completion
		      source '$HOME/.gl/completion.bash.inc'
		      " >> $HOME/.bash_profile
		    source $HOME/.bash_profile

		# Load the glctl completion code for zsh[1] into the current shell
		    source <(glctl completion zsh)
		# Set the glctl completion code for zsh[1] to autoload on startup
		    glctl completion zsh > "${fpath[1]}/_glctl"


		# Load the glctl completion code for fish[2] into the current shell
		    glctl completion fish | source
		# To load completions for each session, execute once:
		    glctl completion fish > ~/.config/fish/completions/glctl.fish

		# Load the glctl completion code for powershell into the current shell
		    glctl completion powershell | Out-String | Invoke-Expression
		# Set glctl completion code for powershell to run on startup
		## Save completion code to a script and execute in the profile
		    glctl completion powershell > $HOME\.gl\completion.ps1
		    Add-Content $PROFILE "$HOME\.gl\completion.ps1"
		## Execute completion code in the profile
		    Add-Content $PROFILE "if (Get-Command glctl -ErrorAction SilentlyContinue) {
		        glctl completion powershell | Out-String | Invoke-Expression
		    }"
		## Add completion code directly to the $PROFILE script
		    glctl completion powershell >> $PROFILE`)
)

var (
	completionShells = map[string]func(out io.Writer, boilerPlate string, cmd *cobra.Command) error{
		"bash":       runCompletionBash,
		"zsh":        runCompletionZsh,
		"fish":       runCompletionFish,
		"powershell": runCompletionPwsh,
	}
)

// NewCmdCompletion creates the `completion` command
func NewCmdCompletion(ioStreams genericiooptions.IOStreams, boilerPlate string) *cobra.Command {
	shells := []string{}
	for s := range completionShells {
		shells = append(shells, s)
	}
	cmd := &cobra.Command{
		Use:                   "completion SHELL",
		Aliases:               []string{"c"},
		Short:                 "Output shell completion code for the specified shell (bash, zsh, fish, or powershell)",
		Long:                  completionLong,
		Example:               completionExample,
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(RunCompletion(ioStreams.Out, boilerPlate, cmd, args))
		},
		ValidArgs: shells,
	}

	return cmd
}

// RunCompletion checks given arguments and executes command
func RunCompletion(out io.Writer, boilerPlate string, cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmdutil.UsageErrorf(cmd, "Shell not specified.")
	}
	if len(args) > 1 {
		return cmdutil.UsageErrorf(cmd, "Too many arguments. Expected only the shell type.")
	}
	run, found := completionShells[args[0]]
	if !found {
		return cmdutil.UsageErrorf(cmd, "Unsupported shell type %q.", args[0])
	}

	return run(out, boilerPlate, cmd.Parent())
}

func runCompletionBash(out io.Writer, boilerPlate string, glctl *cobra.Command) error {
	if len(boilerPlate) == 0 {
		boilerPlate = defaultBoilerPlate
	}
	if _, err := out.Write([]byte(boilerPlate)); err != nil {
		return err
	}

	return glctl.GenBashCompletionV2(out, true)
}

func runCompletionZsh(out io.Writer, boilerPlate string, glctl *cobra.Command) error {
	zshHead := fmt.Sprintf("#compdef %[1]s\ncompdef _%[1]s %[1]s\n", glctl.Name())
	_, err := out.Write([]byte(zshHead))
	if err != nil {
		return err
	}
	if len(boilerPlate) == 0 {
		boilerPlate = defaultBoilerPlate
	}
	if _, err = out.Write([]byte(boilerPlate)); err != nil {
		return err
	}

	return glctl.GenZshCompletion(out)
}

func runCompletionFish(out io.Writer, boilerPlate string, glctl *cobra.Command) error {
	if len(boilerPlate) == 0 {
		boilerPlate = defaultBoilerPlate
	}
	if _, err := out.Write([]byte(boilerPlate)); err != nil {
		return err
	}

	return glctl.GenFishCompletion(out, true)
}

func runCompletionPwsh(out io.Writer, boilerPlate string, glctl *cobra.Command) error {
	if len(boilerPlate) == 0 {
		boilerPlate = defaultBoilerPlate
	}

	if _, err := out.Write([]byte(boilerPlate)); err != nil {
		return err
	}

	return glctl.GenPowerShellCompletionWithDesc(out)
}
