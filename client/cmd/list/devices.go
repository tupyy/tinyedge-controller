package list

import (
	"context"

	"github.com/spf13/cobra"
	rootCmd "github.com/tupyy/tinyedge-controller/client/cmd"
	adminGrpc "github.com/tupyy/tinyedge-controller/pkg/grpc/admin"
)

var devicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "devices",
	Long:  "Print out information about devices.",
	RunE: func(cmd *cobra.Command, args []string) error {
		fn := func(ctx context.Context, client adminGrpc.AdminServiceClient) (*adminGrpc.DevicesListResponse, error) {
			return client.GetDevices(ctx, &adminGrpc.DevicesListRequest{})
		}
		return rootCmd.RunCmd(fn)
	},
}

func init() {
	listCmd.AddCommand(devicesCmd)
}
