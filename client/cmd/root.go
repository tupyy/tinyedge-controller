/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/tupyy/tinyedge-controller/client/internal/common"
	"github.com/tupyy/tinyedge-controller/pkg/grpc/admin"
)

var (
	Url string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:          "client",
	Short:        "CLI for tinyedge",
	Long:         `CLI for tinyedge`,
	SilenceUsage: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func AddCommand(cmd *cobra.Command) {
	rootCmd.AddCommand(cmd)
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&Url, "url", "u", "localhost:8081", "server url")
}

func RunCmd[T any](fn func(ctx context.Context, client admin.AdminServiceClient) (T, error)) error {
	conn, err := common.Dial(Url)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := admin.NewAdminServiceClient(conn)
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
