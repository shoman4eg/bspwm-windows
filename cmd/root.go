package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/shoman4eg/bspwm-windows/cmd/actions"
)

var RootCmd = &cobra.Command{
	Use:   "bspcw",
	Short: "bpswm windows",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		bindEnv(cmd)
		return nil
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func bindEnv(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Determine the naming convention of the flags when represented in the config file
		configName := strings.ReplaceAll(f.Name, ".", "_")
		configValue := os.Getenv(strings.ToUpper(configName))
		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && configValue != "" {
			_ = cmd.Flags().Set(f.Name, fmt.Sprintf("%v", configValue))
		}
	})
}

func init() {
	RootCmd.AddCommand(watchCmd)
	RootCmd.AddCommand(actions.RootCmd)
}
