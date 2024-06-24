package project

import (
	"github.com/AlekSi/pointer"
	"github.com/huhouhua/gitlab-repo-operator/cmd/require"
	cmdutil "github.com/huhouhua/gitlab-repo-operator/cmd/util"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

type CreateOptions struct {
	gitlabClient *gitlab.Client
	Visibility   string
	project      *gitlab.CreateProjectOptions
	Out          string
}

var (
	createProjectDesc = "Create a new branch for a specified project"

	createProjectExample = `# create a develop branch from master branch for project group/myapp
grepo create branch develop --project=group/myapp --ref=master`
)

func NewCreateOptions() *CreateOptions {
	return &CreateOptions{
		project: &gitlab.CreateProjectOptions{
			Description:          pointer.ToString(""),
			LFSEnabled:           pointer.ToBool(false),
			RequestAccessEnabled: pointer.ToBool(false),
		},
		Out: "simple",
	}
}

func NewCreateProjectCmd(f cmdutil.Factory) *cobra.Command {
	o := NewDeleteOptions()
	cmd := &cobra.Command{
		Use:                   "project",
		Aliases:               []string{"p"},
		Short:                 createProjectDesc,
		Example:               createProjectExample,
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
	return cmd
}

// AddFlags registers flags for a cli
func (o *CreateOptions) AddFlags(cmd *cobra.Command) {
	cmdutil.AddDescriptionVarFlag(cmd, o.project.Description)
	cmdutil.AddLFSenabledVarPFlag(cmd, o.project.LFSEnabled)
	cmdutil.AddRequestAccessEnabledVarFlag(cmd, o.project.RequestAccessEnabled)
	cmdutil.AddVisibilityVarFlag(cmd, &o.Visibility)
	// unique flags for projects
	f := cmd.Flags()

	f.Bool("issues-enabled", true, "Enable issues")
	f.Bool("merge-requests-enabled", true, "Enable merge requests")
	f.Bool("jobs-enabled", true, "Enable jobs")
	f.Bool("wiki-enabled", true, "Enable wiki")
	f.Bool("snippets-enabled", true, "Enable snippets")
	f.Bool("resolve-outdated-diff-discussions", false,
		"Automatically resolve merge request diffs discussions on lines "+
			"changed with a push")
	f.Bool("container-registry-enabled", false,
		"Enable container registry for this project")
	f.Bool("shared-runners-enabled", false,
		"Enable shared runners for this project")
	f.Bool("public-jobs", false, "If true, jobs can be viewed "+
		"by non-project-members")
	f.Bool("only-allow-merge-if-pipeline-succeeds", false,
		"Set whether merge requests can only be merged with successful jobs")
	f.Bool("only-allow-merge-if-discussion-are-resolved", false,
		"Set whether merge requests can only be merged "+
			"when all the discussions are resolved")
	f.String("merge-method", "merge",
		"Set the merge method used. (available: 'merge', 'rebase_merge', 'ff')")
	f.StringSlice("tag-list", []string{},
		"The list of tags for a project; put array of tags, "+
			"that should be finally assigned to a project.\n"+
			"Example: --tag-list='tag1,tag2'")
	f.Bool("printing-merge-request-link-enabled", true,
		"Show link to create/view merge request "+
			"when pushing from the command line")
	f.String("ci-config-path", "", "The path to CI config file")
}
