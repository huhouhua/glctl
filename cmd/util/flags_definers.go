package util

import (
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func AddPaginationVarFlags(cmd *cobra.Command, page *gitlab.ListOptions) {
	flags := cmd.Flags()
	flags.IntVarP(&page.Page, "page", "p", page.Page, "Page of results to retrieve")
	flags.IntVarP(&page.PerPage, "per-page", "", page.PerPage, "The number of results to include per page")
}

func AddOutFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP("out", "o", "simple",
		"Print the command output to the "+
			"desired format. (json, yaml, simple)")
}
func AddFromGroupVarFlag(cmd *cobra.Command, p *string) {
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

func AddDescriptionVarFlag(cmd *cobra.Command, p *string) {
	cmd.Flags().StringVar(p, "desc", "", "The description of the resource")
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
