package get

import (
	"context"
	"errors"

	"github.com/spf13/cobra"
	rootCmd "github.com/tupyy/tinyedge-controller/client/cmd"
	adminGrpc "github.com/tupyy/tinyedge-controller/pkg/grpc/admin"
)

var getNamespaceCmd = &cobra.Command{
	Use:   "manifest",
	Short: "manifest [id]",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("Please provide manifest id")
		}
		fn := func(ctx context.Context, client adminGrpc.AdminServiceClient) (*adminGrpc.Manifest, error) {
			return client.GetManifest(ctx, &adminGrpc.IdRequest{Id: args[0]})
		}
		return rootCmd.RunCmd(fn)
	},
}

func init() {
	getCmd.AddCommand(getNamespaceCmd)
}
