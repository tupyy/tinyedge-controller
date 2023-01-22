package delete

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	rootCmd "github.com/tupyy/tinyedge-controller/client/cmd"
	adminGrpc "github.com/tupyy/tinyedge-controller/pkg/grpc/admin"
)

var deleteNamespaceCmd = &cobra.Command{
	Use:   "namespace",
	Short: "namespace [name]",
	Long:  "Removes the namespace",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("namespace name is missing")
		}

		fn := func(ctx context.Context, client adminGrpc.AdminServiceClient) (*adminGrpc.Namespace, error) {
			req := &adminGrpc.IdRequest{
				Id: args[0],
			}
			return client.DeleteNamespace(ctx, req)
		}

		return rootCmd.RunCmd(fn)
	},
}

func init() {
	deleteCmd.AddCommand(deleteNamespaceCmd)
}
