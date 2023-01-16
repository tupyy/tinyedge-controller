package list

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/tupyy/tinyedge-controller/client/cmd"
	rootCmd "github.com/tupyy/tinyedge-controller/client/cmd"
	"github.com/tupyy/tinyedge-controller/client/internal/common"
	"github.com/tupyy/tinyedge-controller/pkg/grpc/admin"
	adminGrpc "github.com/tupyy/tinyedge-controller/pkg/grpc/admin"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List resources",
	Long:  `Use list command to list all resources that match the specified type.`,
}

func init() {
	cmd.AddCommand(listCmd)
}

func runCmd[T any](fn func(ctx context.Context, client admin.AdminServiceClient) (T, error)) error {
	conn, err := common.Dial(rootCmd.Url)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := adminGrpc.NewAdminServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	response, err := fn(ctx, client)
	if err != nil {
		return err
	}

	output, err := jsonOutput(response)
	if err != nil {
		return err
	}
	fmt.Println(output)
	return nil
}
func jsonOutput[T any](response T) (string, error) {
	data, err := json.Marshal(response)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
