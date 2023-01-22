package add

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	rootCmd "github.com/tupyy/tinyedge-controller/client/cmd"
	adminGrpc "github.com/tupyy/tinyedge-controller/pkg/grpc/admin"
	"github.com/tupyy/tinyedge-controller/pkg/grpc/common"
)

var addSet = &cobra.Command{
	Use:   "set",
	Short: "set [name] [options]",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("set name is missing")
		}

		if namespaceID == "" {
			return fmt.Errorf("namespace id is missing")
		}

		fn := func(ctx context.Context, client adminGrpc.AdminServiceClient) (*common.Set, error) {
			req := &adminGrpc.AddSetRequest{
				Id:          args[0],
				NamespaceId: namespaceID,
			}
			if configurationID != "" {
				req.ConfigurationId = &configurationID
			}
			return client.AddSet(ctx, req)
		}

		return rootCmd.RunCmd(fn)
	},
}

func init() {
	addCmd.AddCommand(addSet)

	addSet.Flags().StringVarP(&configurationID, "configuration", "", "", "configuration id")
	addSet.Flags().StringVarP(&namespaceID, "namespace", "", "", "namespace id")
}
