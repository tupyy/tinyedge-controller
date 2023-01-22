/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package add

import (
	"github.com/spf13/cobra"
	"github.com/tupyy/tinyedge-controller/client/cmd"
)

var (
	namespaceID     string
	repoUrl         string
	repoName        string
	name            string
	configurationID string
	isDefault       bool
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add",
}

func init() {
	cmd.AddCommand(addCmd)
}
