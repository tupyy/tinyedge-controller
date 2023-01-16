package list

import (
	"context"

	"github.com/spf13/cobra"
	adminGrpc "github.com/tupyy/tinyedge-controller/pkg/grpc/admin"
)

var setCmd = &cobra.Command{
	Use:   "sets",
	Short: "sets",
	Long:  "Print out information about sets.",
	RunE: func(cmd *cobra.Command, args []string) error {
		fn := func(ctx context.Context, client adminGrpc.AdminServiceClient) (*adminGrpc.SetsListResponse, error) {
			return client.GetSets(ctx, &adminGrpc.ListRequest{})
		}
		return runCmd(fn)
	},
}

func init() {
	listCmd.AddCommand(setCmd)
}
