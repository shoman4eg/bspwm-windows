package actions

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/shoman4eg/bspwm-windows/bspc"
)

const RaiseOrMinimize = "raise_or_minimize"

var raiseOrMinimizeActionCmd = &cobra.Command{
	Use:   RaiseOrMinimize,
	Short: "raise or minimize window",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := bspc.NewClient()
		if err != nil {
			return err
		}

		var win bspc.Node
		if err := client.Query(fmt.Sprintf("query -T -n %v", args[0]), bspc.ToStruct(&win)); err != nil {
			return err
		}

		state := "on"
		if win.Hidden {
			state = "off"
		}

		if err := client.Query(fmt.Sprintf("node %v -g hidden=%v", args[0], state), nil); err != nil {
			return err
		}

		return client.Query(fmt.Sprintf("node -f %v", args[0]), nil)
	},
}
