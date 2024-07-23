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
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	gitlab "github.com/xanzy/go-gitlab"
	"gopkg.in/yaml.v3"
	"io"
	"strconv"
	"strings"
)

var (
	noResultMsg = "The command returned no result. " +
		"Use the (-h) flag to see the command usage."
)

func PrintProjectsOut(format string, w io.Writer, projects ...*gitlab.Project) {
	switch format {
	case JSON:
		printJSON(w, projects)
	case YAML:
		printYAML(w, projects)
	default:
		if len(projects) == 0 {
			fmt.Fprintln(w, noResultMsg)
			return
		}
		header := []string{"ID", "PATH", "URL", "ISSUES COUNT", "TAGS"}
		var rows [][]string
		for _, v := range projects {
			rows = append(rows, []string{
				strconv.Itoa(v.ID),
				v.PathWithNamespace,
				v.HTTPURLToRepo,
				strconv.Itoa(v.OpenIssuesCount),
				strings.Join(v.TagList, ","),
			})
		}
		printTable(header, w, rows)
	}
}

func PrintGroupsOut(format string, w io.Writer, groups ...*gitlab.Group) {
	switch format {
	case JSON:
		printJSON(w, groups)
	case YAML:
		printYAML(w, groups)
	default:
		if len(groups) == 0 {
			fmt.Fprintln(w, noResultMsg)
			return
		}
		header := []string{"ID", "PATH", "URL", "PARENT ID"}
		var rows [][]string
		for _, v := range groups {
			rows = append(rows, []string{
				strconv.Itoa(v.ID),
				v.FullPath,
				v.WebURL,
				strconv.Itoa(v.ParentID),
			})
		}
		printTable(header, w, rows)
	}
}
func PrintBranchOut(format string, w io.Writer, branches ...*gitlab.Branch) {
	switch format {
	case YAML:
		printYAML(w, branches)
	case JSON:
		printJSON(w, branches)
	default:
		if len(branches) == 0 {
			fmt.Fprintln(w, noResultMsg)
			return
		}
		header := []string{"NAME", "PROTECTED", "DEVELOPERS CAN PUSH", "DEVELOPERS CAN MERGE"}
		var rows [][]string
		for _, v := range branches {
			rows = append(rows, []string{
				v.Name,
				strconv.FormatBool(v.Protected),
				strconv.FormatBool(v.DevelopersCanPush),
				strconv.FormatBool(v.DevelopersCanMerge),
			})
		}
		printTable(header, w, rows)
	}
}

func PrintFilesOut(format string, w io.Writer, trees ...*gitlab.TreeNode) {
	switch format {
	case JSON:
		printJSON(w, trees)
	case YAML:
		printYAML(w, trees)
	default:
		if len(trees) == 0 {
			fmt.Println(noResultMsg)
			return
		}
		header := []string{"PATH", "TYPE"}
		var rows [][]string
		for _, v := range trees {
			rows = append(rows, []string{
				v.Path,
				v.Type,
			})
		}
		printTable(header, w, rows)
	}
}

func printJSON(w io.Writer, v interface{}) {
	b, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		Error(w, fmt.Sprintf("failed printing to json: %v", err))
	}
	fmt.Fprintln(w, string(b))
}

func printYAML(w io.Writer, v interface{}) {
	b, err := yaml.Marshal(v)
	if err != nil {
		Error(w, fmt.Sprintf("failed printing to yaml: %v", err))
	}
	fmt.Fprintln(w, string(b))
}

func printTable(header []string, w io.Writer, rows [][]string) {
	if len(header) > 5 {
		panic("maximum allowed length of a table header is only 5.")
	}
	table := tablewriter.NewWriter(w)
	table.SetHeader(header)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetNoWhiteSpace(true)
	table.SetTablePadding("\t") // pad with tabs
	table.AppendBulk(rows)
	table.Render()
}
