/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"github.com/tupyy/tinyedge-controller/client/cmd"
	_ "github.com/tupyy/tinyedge-controller/client/cmd/add"
	_ "github.com/tupyy/tinyedge-controller/client/cmd/delete"
	_ "github.com/tupyy/tinyedge-controller/client/cmd/get"
	_ "github.com/tupyy/tinyedge-controller/client/cmd/list"
	_ "github.com/tupyy/tinyedge-controller/client/cmd/set"
)

func main() {
	cmd.Execute()
}
