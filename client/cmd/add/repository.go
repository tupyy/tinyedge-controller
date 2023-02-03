package add

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	rootCmd "github.com/tupyy/tinyedge-controller/client/cmd"
	adminGrpc "github.com/tupyy/tinyedge-controller/pkg/grpc/admin"
)

var (
	authMethod     string
	authSecretPath string
)

var addRepository = &cobra.Command{
	Use:   "repository",
	Short: "repository",
	RunE: func(cmd *cobra.Command, args []string) error {
		if repoUrl == "" || repoName == "" {
			return fmt.Errorf("repository url or name is missing")
		}

		if authMethod != "" && authSecretPath == "" {
			return fmt.Errorf("secret path is missing")
		}

		switch authMethod {
		case "ssh":
		case "token":
		case "basic":
		default:
			return fmt.Errorf("auth method unknown. Please choose one of \"ssh\", \"token\", \"basic\"")
		}

		fn := func(ctx context.Context, client adminGrpc.AdminServiceClient) (*adminGrpc.AddRepositoryResponse, error) {
			req := &adminGrpc.AddRepositoryRequest{
				Url:            repoUrl,
				Name:           repoName,
				AuthMethod:     authMethod,
				AuthSecretPath: authSecretPath,
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
	addRepository.Flags().StringVar(&authMethod, "auth-method", "", "auth method")
	addRepository.Flags().StringVar(&authSecretPath, "auth-secret-path", "", "auth vault secret path")
}
