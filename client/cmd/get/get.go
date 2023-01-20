/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package get

import (
	"github.com/spf13/cobra"
	rootCmd "github.com/tupyy/tinyedge-controller/client/cmd"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "A brief description of your command",
}

func init() {
	rootCmd.AddCommand(getCmd)
}
