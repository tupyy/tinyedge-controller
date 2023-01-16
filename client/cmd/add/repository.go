package add

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	rootCmd "github.com/tupyy/tinyedge-controller/client/cmd"
	adminGrpc "github.com/tupyy/tinyedge-controller/pkg/grpc/admin"
)

var (
	repoUrl  string
	repoName string
)

var addRepository = &cobra.Command{
	Use:   "repository",
	Short: "repository",
	RunE: func(cmd *cobra.Command, args []string) error {
		if repoUrl == "" || repoName == "" {
			return fmt.Errorf("repository url or name is missing")
		}

		fn := func(ctx context.Context, client adminGrpc.AdminServiceClient) (*adminGrpc.AddRepositoryResponse, error) {
			req := &adminGrpc.AddRepositoryRequest{
				Url:  repoUrl,
				Name: repoName,
			}
			return client.AddRepository(ctx, req)
		}

		return rootCmd.RunCmd(fn)
	},
}

func init() {
	addCmd.AddCommand(addRepository)

	addRepository.Flags().StringVarP(&repoUrl, "repo-url", "r", "", "git repository url")
	addRepository.Flags().StringVarP(&repoName, "repo-name", "n", "", "git repository name")
}
