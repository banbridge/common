package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/banbridge/common/cmd/petal/internal/proto"
	"github.com/banbridge/common/cmd/petal/internal/upgrade"
)

const release = "v0.0.1"

var rootCmd = &cobra.Command{
	Use:     "petal",
	Short:   "petal: An elegant toolkit for Go microservices.",
	Long:    `petal: An elegant toolkit for Go microservices.`,
	Version: release,
}

func init() {
	rootCmd.AddCommand(proto.CmdProto)
	rootCmd.AddCommand(upgrade.CmdUpgrade)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
