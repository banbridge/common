package proto

import (
	"github.com/spf13/cobra"

	"github.com/banbridge/common/cmd/petal/internal/proto/add"
	"github.com/banbridge/common/cmd/petal/internal/proto/client"
	"github.com/banbridge/common/cmd/petal/internal/proto/server"
)

var CmdProto = &cobra.Command{
	Use:   "proto",
	Short: "Generate the proto files",
	Long:  "Generate the proto files.",
}

func init() {
	CmdProto.AddCommand(add.CmdAdd)
	CmdProto.AddCommand(server.CmdServer)
	CmdProto.AddCommand(client.CmdClient)
}
