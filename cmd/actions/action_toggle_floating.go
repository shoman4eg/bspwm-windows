package actions

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/shoman4eg/bspwm-windows/bspc"
)

const Floating = "floating"

var floatingActionCmd = &cobra.Command{
	Use:   Floating,
	Short: "action floating window",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := bspc.NewClient()
		if err != nil {
			return err
		}

		return client.Query(fmt.Sprintf("node %v -t ~floating", args[0]), nil)
	},
}
