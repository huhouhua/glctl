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

package validate

import (
	"fmt"
	"github.com/golang/glog"
	cmdutil "github.com/huhouhua/gitlab-repo-operator/cmd/util"
	"github.com/spf13/cobra"
	"strings"
)

func ValidateSortFlagValue(cmd *cobra.Command) error {
	return ValidateFlagStringValue([]string{"asc", "desc"},
		cmd, "sort")
}

func ValidateProjectOrderByFlagValue(cmd *cobra.Command) error {
	return ValidateFlagStringValue([]string{"id", "name", "path",
		"created_at", "updated_at", "last_activity_at"},
		cmd, "order-by")
}

func ValidateVisibilityFlagValue(cmd *cobra.Command) error {
	return ValidateFlagStringValue([]string{"public", "private", "internal"},
		cmd, "visibility")
}

func ValidateOutFlagValue(cmd *cobra.Command) error {
	return ValidateFlagStringValue([]string{cmdutil.JSON, cmdutil.YAML, "simple"},
		cmd, "out")
}

func ValidateGroupOrderByFlagValue(cmd *cobra.Command) error {
	return ValidateFlagStringValue([]string{"path", "name"},
		cmd, "order-by")
}

func ValidateFlagStringValue(stringSlice []string,
	cmd *cobra.Command, fName string) error {
	fValue := cmdutil.GetFlagString(cmd, fName)
	for _, v := range stringSlice {
		if fValue == v {
			return nil
		}
	}
	return fmt.Errorf("'%s' is not a recognized value of '%s' flag. "+
		"Please choose from: [%s]\n",
		fValue, fName, strings.Join(stringSlice, ", "))
}
func VerifyMarkFlagRequired(cmd *cobra.Command, fName string) {
	if err := cmd.MarkFlagRequired(fName); err != nil {
		glog.Fatalf("error marking %s flag as required for command %s: %v",
			fName, cmd.Name(), err)
	}
}
