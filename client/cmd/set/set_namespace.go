package set

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	rootCmd "github.com/tupyy/tinyedge-controller/client/cmd"
	adminGrpc "github.com/tupyy/tinyedge-controller/pkg/grpc/admin"
)

var (
	isDefault bool
)

var updateNamespace = &cobra.Command{
	Use:   "namespace",
	Short: "namespace [FLAGS]",
	Long:  "Update namespace configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("Please provide a namespace id")
		}
		if !isDefault {
			return fmt.Errorf("Please provide configuration id or set is-default flat")
		}

		fn := func(ctx context.Context, client adminGrpc.AdminServiceClient) (*adminGrpc.Namespace, error) {
			req := &adminGrpc.UpdateNamespaceRequest{
				Id:        args[0],
				IsDefault: isDefault,
			}
			return client.UpdateNamespace(ctx, req)
		}

		return rootCmd.RunCmd(fn)
	},
}

func init() {
	setCmd.AddCommand(updateNamespace)
	updateNamespace.Flags().BoolVarP(&isDefault, "is-default", "d", true, "set the namespace be the default one")
}
