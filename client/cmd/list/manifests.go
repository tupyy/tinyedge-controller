package list

import (
	"context"

	"github.com/spf13/cobra"
	rootCmd "github.com/tupyy/tinyedge-controller/client/cmd"
	adminGrpc "github.com/tupyy/tinyedge-controller/pkg/grpc/admin"
)

var manifestsCmd = &cobra.Command{
	Use:   "manifests",
	Short: "manifests",
	Long:  "Print out information about manifests.",
	RunE: func(cmd *cobra.Command, args []string) error {
		fn := func(ctx context.Context, client adminGrpc.AdminServiceClient) (*adminGrpc.ManifestListResponse, error) {
			return client.GetManifests(ctx, &adminGrpc.ListRequest{})
		}
		return rootCmd.RunCmd(fn)
	},
}

func init() {
	listCmd.AddCommand(manifestsCmd)
}
