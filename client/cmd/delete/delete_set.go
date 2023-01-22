package delete

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	rootCmd "github.com/tupyy/tinyedge-controller/client/cmd"
	adminGrpc "github.com/tupyy/tinyedge-controller/pkg/grpc/admin"
	"github.com/tupyy/tinyedge-controller/pkg/grpc/common"
)

var deleteSetCmd = &cobra.Command{
	Use:   "set",
	Short: "set [name]",
	Long:  "Removes the set",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("set name is missing")
		}

		fn := func(ctx context.Context, client adminGrpc.AdminServiceClient) (*common.Set, error) {
			req := &adminGrpc.IdRequest{
				Id: args[0],
			}
			return client.DeleteSet(ctx, req)
		}

		return rootCmd.RunCmd(fn)
	},
}

func init() {
	deleteCmd.AddCommand(deleteSetCmd)
}
