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
	"github.com/xanzy/go-gitlab"
	"strings"
)

func VerifyMarkFlagRequired(cmd *cobra.Command, fName string) {
	if err := cmd.MarkFlagRequired(fName); err != nil {
		glog.Fatalf("error marking %s flag as required for command %s: %v",
			fName, cmd.Name(), err)
	}
}

func AddPaginationVarFlags(cmd *cobra.Command, page *gitlab.ListOptions) {
	flags := cmd.Flags()
	flags.IntVarP(&page.Page, "page", "p", page.Page, "Page of results to retrieve")
	flags.IntVarP(&page.PerPage, "per-page", "", page.PerPage, "The number of results to include per page")
}

func AddOutFlag(cmd *cobra.Command, p *string) {
	cmd.PersistentFlags().StringVarP(p, "out", "o", *p,
		"Print the command output to the "+
			"desired format. (json, yaml, simple)")
}
func AddFromGroupVarPFlag(cmd *cobra.Command, p *string) {
	cmd.Flags().StringVarP(p, "group", "G", "",
		"Use a group as the target namespace when performing the command")
}

func AddFromProjectVarFlag(cmd *cobra.Command, p *string) {
	cmd.Flags().StringVarP(p, "from-project", "P", "",
		"Use a project as the target namespace when performing the command")
}
func AddProjectOrderByVarFlag(cmd *cobra.Command, p *string) {
	cmd.Flags().StringVar(p, "order-by", *p,
		"Return projects ordered by id, name, path, created_at, updated_at, "+
			"or last_activity_at fields. Default is created_at")
}

func AddProjectVarPFlag(cmd *cobra.Command, p *string) {
	cmd.Flags().StringVarP(p, "project", "p", *p, "The name or ID of the project")
}

func AddDescriptionVarFlag(cmd *cobra.Command, p *string) {
	cmd.Flags().StringVar(p, "desc", "", "The description of the resource")
}

func AddLFSenabledVarPFlag(cmd *cobra.Command, p *bool) {
	cmd.Flags().BoolVar(p, "lfs-enabled", *p, "Enable LFS")
}

func AddRequestAccessEnabledVarFlag(cmd *cobra.Command, p *bool) {
	cmd.Flags().BoolVar(p, "request-access-enabled", *p, "Enable request access")
}

func AddSearchVarFlag(cmd *cobra.Command, p *string) {
	cmd.Flags().StringVar(p, "search", "",
		"Return the list of resources matching the search criteria")
}

func AddOwnedVarFlag(cmd *cobra.Command, p *bool) {
	cmd.Flags().BoolVar(p, "owned", *p,
		"Limit to resources owned by the current user")
}
func AddStatisticsVarFlag(cmd *cobra.Command, p *bool) {
	cmd.Flags().BoolVar(p, "statistics", *p,
		"Include resource statistics (admins only)")
}

func AddSortVarFlag(cmd *cobra.Command, p *string) {
	cmd.Flags().StringVar(p, "sort", *p,
		"Order resources in asc or desc order. Default is asc")
}
func AddVisibilityVarFlag(cmd *cobra.Command, p *string) {
	cmd.Flags().StringVar(p, "visibility", *p, "public, internal or private")
}

// WarnWordSepNormalizeFunc changes and warns for flags that contain "_" separators.
func WarnWordSepNormalizeFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	if strings.Contains(name, "_") {
		nname := strings.ReplaceAll(name, "_", "-")
		glog.Warningf("%s is DEPRECATED and will be removed in a future version. Use %s instead.", name, nname)

		return pflag.NormalizedName(nname)
	}
	return pflag.NormalizedName(name)
}

// WordSepNormalizeFunc changes all flags that contain "_" separators.
func WordSepNormalizeFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	if strings.Contains(name, "_") {
		return pflag.NormalizedName(strings.ReplaceAll(name, "_", "-"))
	}
	return pflag.NormalizedName(name)
}
