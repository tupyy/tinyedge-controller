package set

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	rootCmd "github.com/tupyy/tinyedge-controller/client/cmd"
	adminGrpc "github.com/tupyy/tinyedge-controller/pkg/grpc/admin"
	"github.com/tupyy/tinyedge-controller/pkg/grpc/common"
)

var (
	setID           string
	namepsaceID     string
	configurationID string
)
var deviceToSet = &cobra.Command{
	Use:   "device",
	Short: "device [FLAGS]",
	Long:  "Update device set, namespace or configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("Please provide a device id")
		}
		deviceID := args[len(args)-1]
		if setID == "" && namepsaceID == "" && configurationID == "" {
			return fmt.Errorf("Please provide at least a set id, namespace id or configuration id")
		}

		fn := func(ctx context.Context, client adminGrpc.AdminServiceClient) (*common.Device, error) {
			req := &adminGrpc.UpdateDeviceRequest{
				Id:              deviceID,
				SetId:           setID,
				NamespaceId:     namepsaceID,
				ConfigurationId: configurationID,
			}
			return client.UpdateDevice(ctx, req)
		}

		return rootCmd.RunCmd(fn)
	},
}

func init() {
	setCmd.AddCommand(deviceToSet)
	deviceToSet.Flags().StringVarP(&setID, "set", "s", "", "set id")
	deviceToSet.Flags().StringVarP(&namepsaceID, "namespace", "n", "", "namespace id")
	deviceToSet.Flags().StringVarP(&configurationID, "configuration", "c", "", "configuration id")
}
