package set

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	rootCmd "github.com/tupyy/tinyedge-controller/client/cmd"
	adminGrpc "github.com/tupyy/tinyedge-controller/pkg/grpc/admin"
	"github.com/tupyy/tinyedge-controller/pkg/grpc/common"
)

var updateNamespace = &cobra.Command{
	Use:   "namespace",
	Short: "namespace [FLAGS]",
	Long:  "Update namespace configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("Please provide a namespace id")
		}
		if configurationID == "" {
			return fmt.Errorf("Please provide at least a set id, namespace id or configuration id")
		}

		fn := func(ctx context.Context, client adminGrpc.AdminServiceClient) (*common.Device, error) {
			return nil, nil
		}

		return rootCmd.RunCmd(fn)
	},
}

func init() {
	setCmd.AddCommand(updateNamespace)
	updateNamespace.Flags().StringVarP(&configurationID, "configuration", "c", "", "configuration id")
}
