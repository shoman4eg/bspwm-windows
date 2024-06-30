package cmd

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
	"os/signal"
	"syscall"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/shoman4eg/bspwm-windows/bspc"
	"github.com/shoman4eg/bspwm-windows/cmd/actions"
	"github.com/shoman4eg/bspwm-windows/polybar"
)

var cfg watchConfig

type watchConfig struct {
	format     string
	monitor    string
	configPath string
}

func (cfg *watchConfig) Flags() *pflag.FlagSet {
	f := pflag.NewFlagSet("WatchConfig", pflag.PanicOnError)

	f.StringVarP(&cfg.monitor, "monitor", "m", "", "set monitor")
	f.StringVar(&cfg.configPath, "config", "config.toml", "set config_path")
	f.StringVar(&cfg.format, "format", "polybar", "set format for output")

	return f
}

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch for changes",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, stop := signal.NotifyContext(cmd.Context(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()
		c, err := bspc.NewClient()
		if err != nil {
			return errors.WithMessage(err, "failed to create bspc client")
		}

		var config polybar.Config

		_, err = toml.DecodeFile(cfg.configPath, &config)
		if err != nil {
			return errors.WithMessage(err, "failed to load config file")
		}

		err = c.Subscribe(ctx, "node_add node_remove node_activate node_focus node_flag desktop_focus", func(bytes []byte) {
			q := "query -N -n .window --desktop"

			if cfg.monitor != "" {
				q += fmt.Sprintf(" %v:focused", cfg.monitor)
			}

			var allIDs []bspc.ID
			if err0 := c.Query(q, bspc.ToIDSlice(&allIDs)); err0 != nil {
				log.Fatal(err0)
			}
			var hiddenIDs []bspc.ID
			if err0 := c.Query("query -N -n .hidden", bspc.ToIDSlice(&hiddenIDs)); err0 != nil {
				log.Fatal(err0)
			}

			var activeID bspc.ID
			if err0 := c.Query("query -N -n", bspc.ToID(&activeID)); err0 != nil {
				log.Fatal(err0)
			}

			if config.MaxWindows > 0 && len(allIDs) > config.MaxWindows {
				allIDs = allIDs[0:config.MaxWindows]
			}
			if len(allIDs) == 0 {
				fmt.Println(config.EmptyDesktopString)
			}
			for _, id := range allIDs {
				var win bspc.Node
				if err := c.Query(fmt.Sprintf("query -T -n %v", id), bspc.ToStruct(&win)); err != nil {
					continue
				}

				isActive := win.ID == activeID

				title := win.Client.ClassName
				if t, ok := config.WindowNicknames[title]; ok {
					title = t
				}
				if config.NameMaxLength > 0 && len(title) > config.NameMaxLength {
					title = title[:config.NameMaxLength]
				}
				title = fmt.Sprintf(" %v ", title)

				hasFlag := false
				if win.Locked && config.Flags.LockedFlag != "" {
					title += config.Flags.LockedFlag
					hasFlag = true
				}
				if win.Private && config.Flags.PrivateFlag != "" {
					title += config.Flags.PrivateFlag
					hasFlag = true
				}
				if win.Sticky && config.Flags.StickyFlag != "" {
					title += config.Flags.StickyFlag
					hasFlag = true
				}
				if win.Marked && config.Flags.MarkedFlag != "" {
					title += config.Flags.MarkedFlag
					hasFlag = true
				}

				if hasFlag {
					title += " "
				}

				if color := config.GetBgColor(isActive, win.Hidden); color != "" {
					title = fmt.Sprintf("%%{B%v}%v%%{B-}", color, title)
				}
				if color := config.GetFgColor(isActive, win.Hidden); color != "" {
					title = fmt.Sprintf("%%{F%v}%v%%{F-}", color, title)
				}
				if action := config.GetActionLeftClick(isActive, win.Hidden); action != "" {
					title = fmt.Sprintf("%%{A1:%v %v:}%v%%{A}", actions.NormalizeCommand(action), win.ID, title)
				}
				if action := config.GetActionMiddleClick(isActive, win.Hidden); action != "" {
					title = fmt.Sprintf("%%{A2:%v %v:}%v%%{A}", actions.NormalizeCommand(action), win.ID, title)
				}
				if action := config.GetActionRightClick(isActive, win.Hidden); action != "" {
					title = fmt.Sprintf("%%{A3:%v %v:}%v%%{A}", actions.NormalizeCommand(action), win.ID, title)
				}
				if action := config.GetActionScrollUp(isActive, win.Hidden); action != "" {
					title = fmt.Sprintf("%%{A4:%v %v:}%v%%{A}", actions.NormalizeCommand(action), win.ID, title)
				}
				if action := config.GetActionScrollDown(isActive, win.Hidden); action != "" {
					title = fmt.Sprintf("%%{A5:%v %v:}%v%%{A}", actions.NormalizeCommand(action), win.ID, title)
				}

				fmt.Print(title)
			}
			fmt.Println("")
		})
		if err != nil {
			return err
		}
		<-ctx.Done()

		stop()
		return nil
	},
}

func init() {
	watchCmd.Flags().AddFlagSet(cfg.Flags())
}
