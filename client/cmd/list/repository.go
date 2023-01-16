/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package list

import (
	"context"

	"github.com/spf13/cobra"
	rootCmd "github.com/tupyy/tinyedge-controller/client/cmd"
	adminGrpc "github.com/tupyy/tinyedge-controller/pkg/grpc/admin"
)

// repositoryCmd represents the repository command
var repositoryCmd = &cobra.Command{
	Use:   "repositories",
	Short: "repositories [options]",
	Long:  `Print out information about repositories.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fn := func(ctx context.Context, client adminGrpc.AdminServiceClient) (*adminGrpc.RepositoryListResponse, error) {
			return client.GetRepositories(ctx, &adminGrpc.ListRequest{})
		}
		return rootCmd.RunCmd(fn)
	},
}

func init() {
	listCmd.AddCommand(repositoryCmd)
}
