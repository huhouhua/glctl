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

package project

import (
	"github.com/AlekSi/pointer"
	cmdutil "github.com/huhouhua/glctl/cmd/util"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

// assignOptions assigns the flags' values to gitlab.EditProjectOptions fields.
// If a flag's default value is not changed by the caller,
// it's value will not be assigned to the associated gitlab.EditProjectOptions field.
func (o *EditOptions) assignOptions(cmd *cobra.Command) error {
	if cmd.Flag("desc").Changed {
		o.project.Description = pointer.ToString(cmdutil.GetFlagString(cmd, "desc"))
	}
	if cmd.Flag("issues_access_level").Changed {
		o.project.IssuesAccessLevel = pointer.To(gitlab.AccessControlValue(cmdutil.GetFlagString(cmd, "issues_access_level")))
	}
	if cmd.Flag("merge_requests_access_level").Changed {
		o.project.MergeRequestsAccessLevel = pointer.To(gitlab.AccessControlValue(cmdutil.GetFlagString(cmd, "merge_requests_access_level")))
	}
	if cmd.Flag("builds_access_level").Changed {
		o.project.BuildsAccessLevel = pointer.To(gitlab.AccessControlValue(cmdutil.GetFlagString(cmd, "builds_access_level")))
	}
	if cmd.Flag("wiki_access_level").Changed {
		o.project.WikiAccessLevel = pointer.To(gitlab.AccessControlValue(cmdutil.GetFlagString(cmd, "wiki_access_level")))
	}
	if cmd.Flag("snippets_access_level").Changed {
		o.project.SnippetsAccessLevel = pointer.To(gitlab.AccessControlValue(cmdutil.GetFlagString(cmd, "snippets_access_level")))
	}
	if cmd.Flag("resolve-outdated-diff-discussions").Changed {
		o.project.ResolveOutdatedDiffDiscussions = pointer.ToBool(cmdutil.GetFlagBool(cmd, "resolve-outdated-diff-discussions"))
	}
	if cmd.Flag("container_registry_access_level").Changed {
		o.project.ContainerRegistryAccessLevel = pointer.To(gitlab.AccessControlValue(cmdutil.GetFlagString(cmd, "container_registry_access_level")))
	}
	if cmd.Flag("shared-runners-enabled").Changed {
		o.project.SharedRunnersEnabled = pointer.ToBool(cmdutil.GetFlagBool(cmd, "shared-runners-enabled"))
	}
	if cmd.Flag("public_builds").Changed {
		o.project.PublicBuilds = pointer.ToBool(cmdutil.GetFlagBool(cmd, "public_builds"))
	}
	if cmd.Flag("only-allow-merge-if-pipeline-succeeds").Changed {
		o.project.OnlyAllowMergeIfPipelineSucceeds = pointer.ToBool(cmdutil.GetFlagBool(cmd, "only-allow-merge-if-pipeline-succeeds"))
	}
	if cmd.Flag("only-allow-merge-if-discussion-are-resolved").Changed {
		o.project.OnlyAllowMergeIfAllDiscussionsAreResolved = pointer.ToBool(cmdutil.GetFlagBool(cmd, "only-allow-merge-if-discussion-are-resolved"))
	}
	if cmd.Flag("merge-method").Changed {
		o.project.MergeMethod = pointer.To(gitlab.MergeMethodValue(cmdutil.GetFlagString(cmd, "merge-method")))
	}
	if cmd.Flag("tag-list").Changed {
		p := new([]string)
		*p = cmdutil.GetFlagStringSlice(cmd, "tag-list")
		o.project.TagList = p
	}
	if cmd.Flag("printing-merge-request-link-enabled").Changed {
		o.project.PrintingMergeRequestLinkEnabled = pointer.ToBool(cmdutil.GetFlagBool(cmd, "printing-merge-request-link-enabled"))
	}
	if cmd.Flag("ci-config-path").Changed {
		o.project.CIConfigPath = pointer.ToString(cmdutil.GetFlagString(cmd, "ci-config-path"))
	}
	return nil
}
