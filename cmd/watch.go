package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var cfg watchConfig

type watchConfig struct {
	format     string
	monitor    string
	configPath string
}

func (cfg *watchConfig) Flags() *pflag.FlagSet {
	f := pflag.NewFlagSet("WatchConfig", pflag.PanicOnError)

	f.StringVar(&cfg.monitor, "monitor", "", "set monitor")
	f.StringVar(&cfg.configPath, "config_path", "", "set config_path")
	f.StringVar(&cfg.format, "format", "polybar", "set format for output")

	return f
}

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch for changes",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	watchCmd.Flags().AddFlagSet(cfg.Flags())
}
