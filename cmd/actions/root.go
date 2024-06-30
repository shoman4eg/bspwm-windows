package actions

import (
	"os"

	"github.com/spf13/cobra"
)

var currentBin string

var RootCmd = &cobra.Command{
	Use:   "actions",
	Short: "window click actions",
}

var availableCommands = map[string]string{
	Floating:        "actions " + Floating,
	RaiseOrMinimize: "actions " + RaiseOrMinimize,
	Close:           "actions " + Close,
}

func NormalizeCommand(command string) string {
	if cmd, ok := availableCommands[command]; ok {
		return currentBin + " " + cmd
	}

	return command
}

func init() {
	RootCmd.AddCommand(closeActionCmd)
	RootCmd.AddCommand(raiseOrMinimizeActionCmd)
	RootCmd.AddCommand(floatingActionCmd)

	var err error
	if currentBin, err = os.Executable(); err != nil {
		panic(err)
	}
}
