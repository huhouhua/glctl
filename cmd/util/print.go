// Copyright 2024 The Kevin Berger <huhouhuam@gmail.com> Authors
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
	"io"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter/tw"

	"github.com/olekukonko/tablewriter"
	gitlab "gitlab.com/gitlab-org/api/client-go"
	"gopkg.in/yaml.v3"
)

var (
	noResultMsg = "The command returned no result. " +
		"Use the (-h) flag to see the command usage."
)

func PrintProjectsOut(format string, w io.Writer, projects ...*gitlab.Project) error {
	switch format {
	case JSON:
		return printJSON(w, projects)
	case YAML:
		return printYAML(w, projects)
	default:
		if len(projects) == 0 {
			_, err := fmt.Fprintln(w, noResultMsg)
			return err
		}
		header := []string{"ID", "PATH", "URL", "ISSUES COUNT", "TAGS"}
		var rows [][]string
		for _, v := range projects {
			rows = append(rows, []string{
				strconv.FormatInt(v.ID, 10),
				v.PathWithNamespace,
				v.HTTPURLToRepo,
				strconv.FormatInt(v.OpenIssuesCount, 10),
				strings.Join(v.Topics, ","),
			})
		}
		return printTable(header, w, rows)
	}
}

func PrintGroupsOut(format string, w io.Writer, groups ...*gitlab.Group) error {
	switch format {
	case JSON:
		return printJSON(w, groups)
	case YAML:
		return printYAML(w, groups)
	default:
		if len(groups) == 0 {
			_, err := fmt.Fprintln(w, noResultMsg)
			return err
		}
		header := []string{"ID", "PATH", "URL", "PARENT ID"}
		var rows [][]string
		for _, v := range groups {
			rows = append(rows, []string{
				strconv.FormatInt(v.ID, 10),
				v.FullPath,
				v.WebURL,
				strconv.FormatInt(v.ParentID, 10),
			})
		}
		return printTable(header, w, rows)
	}
}
func PrintBranchOut(format string, w io.Writer, branches ...*gitlab.Branch) error {
	switch format {
	case YAML:
		return printYAML(w, branches)
	case JSON:
		return printJSON(w, branches)
	default:
		if len(branches) == 0 {
			_, err := fmt.Fprintln(w, noResultMsg)
			return err
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
		return printTable(header, w, rows)
	}
}

func PrintFilesOut(format string, w io.Writer, trees ...*gitlab.TreeNode) error {
	switch format {
	case JSON:
		return printJSON(w, trees)
	case YAML:
		return printYAML(w, trees)
	default:
		if len(trees) == 0 {
			_, err := fmt.Fprintln(w, noResultMsg)
			return err
		}
		header := []string{"PATH", "TYPE"}
		var rows [][]string
		for _, v := range trees {
			rows = append(rows, []string{
				v.Path,
				v.Type,
			})
		}
		return printTable(header, w, rows)
	}
}

func printJSON(w io.Writer, v interface{}) error {
	b, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		Error(w, fmt.Sprintf("failed printing to json: %v", err))
	}
	_, err = fmt.Fprintln(w, string(b))
	return err
}

func printYAML(w io.Writer, v interface{}) error {
	b, err := yaml.Marshal(v)
	if err != nil {
		Error(w, fmt.Sprintf("failed printing to yaml: %v", err))
	}
	_, err = fmt.Fprintln(w, string(b))
	return err
}

func printTable(header []string, w io.Writer, rows [][]string) error {
	if len(header) > 5 {
		panic("maximum allowed length of a table header is only 5.")
	}
	table := tablewriter.NewTable(w,
		tablewriter.WithTrimSpace(tw.Off),
		tablewriter.WithAlignment(tw.Alignment{tw.AlignLeft, tw.AlignLeft, tw.AlignLeft}),
		tablewriter.WithRendition(tw.Rendition{
			Borders: tw.BorderNone,
			Settings: tw.Settings{
				Separators: tw.Separators{
					ShowHeader:     tw.Off,
					ShowFooter:     tw.Off,
					BetweenRows:    tw.Off,
					BetweenColumns: tw.Off,
				},
				Lines: tw.Lines{
					ShowHeaderLine: tw.Off,
				},
			},
		}),
	)
	table.Header(header)
	if err := table.Bulk(rows); err != nil {
		return err
	}
	if err := table.Render(); err != nil {
		return err
	}
	return nil
}
