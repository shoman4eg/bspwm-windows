package actions

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/shoman4eg/bspwm-windows/bspc"
)

const Close = "close"

var closeActionCmd = &cobra.Command{
	Use:   Close,
	Short: "close an action",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := bspc.NewClient()
		if err != nil {
			return err
		}

		return client.Query(fmt.Sprintf("node %v -c", args[0]), nil)
	},
}
