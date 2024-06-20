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
	if cmd.Flag("page").Changed {
		opts.Page = GetFlagInt(cmd, "page")
	}
	if cmd.Flag("per-page").Changed {
		opts.PerPage = GetFlagInt(cmd, "per-page")
	}
	if cmd.Flag("archived").Changed {
		opts.Archived = gitlab.Ptr(GetFlagBool(cmd, "archived"))
	}
	if cmd.Flag("order-by").Changed {
		opts.OrderBy = gitlab.Ptr(GetFlagString(cmd, "order-by"))
	}
	if cmd.Flag("sort").Changed {
		opts.Sort = gitlab.Ptr(GetFlagString(cmd, "sort"))
	}
	if cmd.Flag("search").Changed {
		opts.Search = gitlab.Ptr(GetFlagString(cmd, "search"))
	}
	if cmd.Flag("simple").Changed {
		opts.Simple = gitlab.Ptr(GetFlagBool(cmd, "simple"))
	}
	if cmd.Flag("owned").Changed {
		opts.Owned = gitlab.Ptr(GetFlagBool(cmd, "owned"))
	}
	if cmd.Flag("membership").Changed {
		opts.Membership = gitlab.Ptr(GetFlagBool(cmd, "membership"))
	}
	if cmd.Flag("starred").Changed {
		opts.Starred = gitlab.Ptr(GetFlagBool(cmd, "starred"))
	}
	if cmd.Flag("statistics").Changed {
		opts.Statistics = gitlab.Ptr(GetFlagBool(cmd, "statistics"))
	}
	if cmd.Flag("visibility").Changed {
		v := GetFlagVisibility(cmd)
		opts.Visibility = v
	}
	if cmd.Flag("with-issues-enabled").Changed {
		opts.WithIssuesEnabled = gitlab.Ptr(
			GetFlagBool(cmd, "with-issues-enabled"))
	}
	if cmd.Flag("with-merge-requests-enabled").Changed {
		opts.WithMergeRequestsEnabled = gitlab.Ptr(GetFlagBool(cmd,
			"with-merge-requests-enabled"))
	}
	return opts
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
