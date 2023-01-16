package list

import (
	"github.com/spf13/cobra"
	"github.com/tupyy/tinyedge-controller/client/cmd"
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
