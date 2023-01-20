package set

import (
	"github.com/spf13/cobra"
	"github.com/tupyy/tinyedge-controller/client/cmd"
)

// setCmd represents the associate command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "A brief description of your command",
}

func init() {
	cmd.AddCommand(setCmd)
}
