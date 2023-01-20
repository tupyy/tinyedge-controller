package get

import (
	"context"
	"errors"

	"github.com/spf13/cobra"
	rootCmd "github.com/tupyy/tinyedge-controller/client/cmd"
	adminGrpc "github.com/tupyy/tinyedge-controller/pkg/grpc/admin"
	"github.com/tupyy/tinyedge-controller/pkg/grpc/common"
)

var getSetCmd = &cobra.Command{
	Use:   "set",
	Short: "set [id]",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("Please provide set id")
		}
		fn := func(ctx context.Context, client adminGrpc.AdminServiceClient) (*common.Set, error) {
			return client.GetSet(ctx, &adminGrpc.IdRequest{Id: args[0]})
		}
		return rootCmd.RunCmd(fn)
	},
}

func init() {
	getCmd.AddCommand(getSetCmd)
}
