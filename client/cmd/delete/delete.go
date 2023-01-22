package delete

import (
	"github.com/spf13/cobra"
	"github.com/tupyy/tinyedge-controller/client/cmd"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete a resource",
}

func init() {
	cmd.AddCommand(deleteCmd)
}
