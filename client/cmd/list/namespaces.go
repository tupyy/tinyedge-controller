package list

import (
	"context"

	"github.com/spf13/cobra"
	rootCmd "github.com/tupyy/tinyedge-controller/client/cmd"
	adminGrpc "github.com/tupyy/tinyedge-controller/pkg/grpc/admin"
)

var namespaceCmd = &cobra.Command{
	Use:   "namespaces",
	Short: "namespaces",
	Long:  "Print out information about namespaces.",
	RunE: func(cmd *cobra.Command, args []string) error {
		fn := func(ctx context.Context, client adminGrpc.AdminServiceClient) (*adminGrpc.NamespaceListResponse, error) {
			return client.GetNamespaces(ctx, &adminGrpc.ListRequest{})
		}
		return rootCmd.RunCmd(fn)
	},
}

func init() {
	listCmd.AddCommand(namespaceCmd)
}
