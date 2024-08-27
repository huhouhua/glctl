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

package branch

import (
	"fmt"

	"github.com/AlekSi/pointer"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"

	"github.com/huhouhua/glctl/cmd/require"
	cmdutil "github.com/huhouhua/glctl/cmd/util"
	"github.com/huhouhua/glctl/cmd/validate"
	"github.com/huhouhua/glctl/util/cli"
	"github.com/huhouhua/glctl/util/templates"
)

type EditOptions struct {
	gitlabClient      *gitlab.Client
	protectBranch     *gitlab.ProtectBranchOptions
	protectRepository *gitlab.ProtectRepositoryBranchesOptions
	branch            string
	project           string
	protect           bool
	Unprotect         bool
	Out               string
	ioStreams         cli.IOStreams
}

var (
	editBranchExample = templates.Examples(`
# protect a branch
glctl edit branch master -p test/glctl --protect

glctl edit branch master -p test/glctl --unprotect`)
)

func NewEditOptions(ioStreams cli.IOStreams) *EditOptions {
	return &EditOptions{
		ioStreams: ioStreams,
		protectBranch: &gitlab.ProtectBranchOptions{
			DevelopersCanPush:  pointer.ToBool(false),
			DevelopersCanMerge: pointer.ToBool(false),
		},
		protectRepository: &gitlab.ProtectRepositoryBranchesOptions{},
		Out:               "simple",
	}
}

func NewEditBranchCmd(f cmdutil.Factory, ioStreams cli.IOStreams) *cobra.Command {
	o := NewEditOptions(ioStreams)
	cmd := &cobra.Command{
		Use:                   "branch",
		Aliases:               []string{"b"},
		Short:                 "Protect or unprotect a repositort branch",
		Example:               editBranchExample,
		Args:                  require.ExactArgs(1),
		DisableFlagsInUseLine: true,
		TraverseChildren:      true,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Validate(cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
	}
	o.AddFlags(cmd)
	return cmd
}

// AddFlags registers flags for a cli
func (o *EditOptions) AddFlags(cmd *cobra.Command) {
	cmdutil.AddProjectVarPFlag(cmd, &o.project)
	cmdutil.AddOutFlag(cmd, &o.Out)
	validate.VerifyMarkFlagRequired(cmd, "project")
	f := cmd.Flags()
	f.BoolVar(&o.Unprotect, "unprotect", o.Unprotect,
		"Remove protection of a branch")
	f.BoolVar(&o.protect, "protect",
		o.protect, "Protect a branch")
	f.BoolVar(o.protectBranch.DevelopersCanPush, "dev-can-push", *o.protectBranch.DevelopersCanPush,
		"Used with '--protect'. Flag if developers can push to the branch")
	f.BoolVar(o.protectBranch.DevelopersCanMerge, "dev-can-merge", *o.protectBranch.DevelopersCanMerge,
		"Used with '--protect'. Flag if developers can merge to the branch")
}

// Complete completes all the required options.
func (o *EditOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	var err error
	o.gitlabClient, err = f.GitlabClient()
	if len(args) > 0 {
		o.protectRepository.Name = pointer.ToString(args[0])
	}
	if *o.protectBranch.DevelopersCanMerge {
		o.protectRepository.MergeAccessLevel = pointer.To(gitlab.DeveloperPermissions)
	}
	if *o.protectBranch.DevelopersCanPush {
		o.protectRepository.PushAccessLevel = pointer.To(gitlab.DeveloperPermissions)
	}
	return err
}

// Validate makes sure there is no discrepency in command options.
func (o *EditOptions) Validate(cmd *cobra.Command, args []string) error {
	if o.protect {
		if !*o.protectBranch.DevelopersCanMerge {
			return fmt.Errorf("'--%s' flag can only be used with '--%s' flag",
				"dev-can-push", "protect")
		}
		if !*o.protectBranch.DevelopersCanPush {
			return fmt.Errorf("'--%s' flag can only be used with '--%s' flag",
				"dev-can-merge", "protect")
		}
		if o.Unprotect {
			return fmt.Errorf("use only (1) flag from (protect)")
		}
	}
	return nil
}

// Run executes a list subcommand using the specified options.
func (o *EditOptions) Run(args []string) error {
	if o.protect {
		_, _, err := o.gitlabClient.ProtectedBranches.ProtectRepositoryBranches(o.project, o.protectRepository)
		if err != nil {
			return err
		}
		_, _ = fmt.Fprintf(o.ioStreams.Out, "branch %s updated\n", args[0])
		return nil
	}
	_, err := o.gitlabClient.ProtectedBranches.UnprotectRepositoryBranches(o.project, args[0])
	if err != nil {
		return err
	}
	_, _ = fmt.Fprintf(o.ioStreams.Out, "branch %s un protect\n", args[0])
	return nil
}
