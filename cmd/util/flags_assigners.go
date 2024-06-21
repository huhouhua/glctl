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

package util

import (
	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	gitlab "github.com/xanzy/go-gitlab"
)

const (
	// JSON is used as a constant of word "json" for out flag
	JSON = "json"
	// YAML is used as a constant of word "yaml" for out flag
	YAML = "yaml"
)

// AssignListProjectOptions assigns the validate' values to gitlab.ListProjectsOptions fields.
// If a flag's default value is not changed by the caller,
// it's value will not be assigned to the associated gitlab.ListProjectsOptions field.
func AssignListProjectOptions(cmd *cobra.Command) *gitlab.ListProjectsOptions {
	opts := new(gitlab.ListProjectsOptions)
	if isChanged(cmd.Flag("page")) {
		opts.Page = GetFlagInt(cmd, "page")
	}
	if isChanged(cmd.Flag("per-page")) {
		opts.PerPage = GetFlagInt(cmd, "per-page")
	}
	if isChanged(cmd.Flag("archived")) {
		opts.Archived = gitlab.Ptr(GetFlagBool(cmd, "archived"))
	}
	if isChanged(cmd.Flag("order-by")) {
		opts.OrderBy = gitlab.Ptr(GetFlagString(cmd, "order-by"))
	}
	if isChanged(cmd.Flag("sort")) {
		opts.Sort = gitlab.Ptr(GetFlagString(cmd, "sort"))
	}
	if isChanged(cmd.Flag("search")) {
		opts.Search = gitlab.Ptr(GetFlagString(cmd, "search"))
	}
	if isChanged(cmd.Flag("simple")) {
		opts.Simple = gitlab.Ptr(GetFlagBool(cmd, "simple"))
	}
	if isChanged(cmd.Flag("owned")) {
		opts.Owned = gitlab.Ptr(GetFlagBool(cmd, "owned"))
	}
	if isChanged(cmd.Flag("membership")) {
		opts.Membership = gitlab.Ptr(GetFlagBool(cmd, "membership"))
	}
	if isChanged(cmd.Flag("starred")) {
		opts.Starred = gitlab.Ptr(GetFlagBool(cmd, "starred"))
	}
	if isChanged(cmd.Flag("statistics")) {
		opts.Statistics = gitlab.Ptr(GetFlagBool(cmd, "statistics"))
	}
	if isChanged(cmd.Flag("with-merge-requests-enabled")) {
		v := GetFlagVisibility(cmd)
		opts.Visibility = v
	}
	if isChanged(cmd.Flag("with-issues-enabled")) {
		opts.WithIssuesEnabled = gitlab.Ptr(
			GetFlagBool(cmd, "with-issues-enabled"))
	}
	if isChanged(cmd.Flag("with-merge-requests-enabled")) {
		opts.WithMergeRequestsEnabled = gitlab.Ptr(GetFlagBool(cmd,
			"with-merge-requests-enabled"))
	}
	return opts
}

func isChanged(flag *pflag.Flag) bool {
	return flag != nil && flag.Changed
}

// GetFlagVisibility converts the string flag visiblity to gitlab.VisibilityValue.
func GetFlagVisibility(cmd *cobra.Command) *gitlab.VisibilityValue {
	v := GetFlagString(cmd, "visibility")
	return gitlab.Visibility(gitlab.VisibilityValue(v))
}

func GetFlagStringSlice(cmd *cobra.Command, flag string) []string {
	s, err := cmd.Flags().GetStringSlice(flag)
	if err != nil {
		glog.Fatalf("error accessing flag %s for command %s: %v",
			flag, cmd.Name(), err)
	}
	return s
}

func GetFlagBool(cmd *cobra.Command, flag string) bool {
	b, err := cmd.Flags().GetBool(flag)
	if err != nil {
		glog.Fatalf("error accessing flag %s for command %s: %v",
			flag, cmd.Name(), err)
	}
	return b
}

func GetFlagInt(cmd *cobra.Command, flag string) int {
	i, err := cmd.Flags().GetInt(flag)
	if err != nil {
		glog.Fatalf("error accessing flag %s for command %s: %v",
			flag, cmd.Name(), err)
	}
	return i
}
func GetFlagString(cmd *cobra.Command, flag string) string {
	s, err := cmd.Flags().GetString(flag)
	if err != nil {
		glog.Fatalf("error accessing flag %s for command %s: %v",
			flag, cmd.Name(), err)
	}
	return s
}

func AddOutFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP("out", "o", "simple",
		"Print the command output to the "+
			"desired format. (json, yaml, simple)")
}
