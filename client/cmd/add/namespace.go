package add

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	rootCmd "github.com/tupyy/tinyedge-controller/client/cmd"
	adminGrpc "github.com/tupyy/tinyedge-controller/pkg/grpc/admin"
)

var addNamespace = &cobra.Command{
	Use:   "namespace",
	Short: "namespace",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("namespace name is missing")
		}

		if configurationID == "" {
			return fmt.Errorf("configuration id is missing")
		}

		fn := func(ctx context.Context, client adminGrpc.AdminServiceClient) (*adminGrpc.Namespace, error) {
			req := &adminGrpc.AddNamespaceRequest{
				Id:              args[0],
				ConfigurationId: configurationID,
				IsDefault:       isDefault,
			}
			return client.AddNamespace(ctx, req)
		}

		return rootCmd.RunCmd(fn)
	},
}

func init() {
	addCmd.AddCommand(addNamespace)

	addNamespace.Flags().StringVarP(&configurationID, "configuration", "", "", "configuration id")
	addNamespace.Flags().BoolVarP(&isDefault, "is-default", "", false, "true if the namespace is the default one")
}
