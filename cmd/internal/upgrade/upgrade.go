package upgrade

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/banbridge/common/cmd/petal/internal/base"
)

// CmdUpgrade represents the upgrade command.
var CmdUpgrade = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade the yakt tools",
	Long:  "Upgrade the yakt tools. Example: yakt upgrade",
	Run:   Run,
}

// Run upgrade the kratos tools.
func Run(cmd *cobra.Command, args []string) {
	err := base.GoInstall(

		"google.golang.org/protobuf/cmd/protoc-gen-go@latest",
		"github.com/google/wire/cmd/wire@latest",
	)
	if err != nil {
		fmt.Println(err)
	}
}
