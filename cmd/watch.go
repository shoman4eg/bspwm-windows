package cmd

import (
	"fmt"
	"log"
	"os/signal"
	"slices"
	"strings"
	"syscall"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/shoman4eg/bspwm-windows/bspc"
	"github.com/shoman4eg/bspwm-windows/cmd/actions"
	"github.com/shoman4eg/bspwm-windows/config"
)

var cfg watchConfig

type watchConfig struct {
	monitor    string
	configPath string
}

func (cfg *watchConfig) Flags() *pflag.FlagSet {
	f := pflag.NewFlagSet("WatchConfig", pflag.PanicOnError)

	f.StringVarP(&cfg.monitor, "monitor", "m", "", "set monitor")
	f.StringVar(&cfg.configPath, "config", "config.toml", "set config_path")

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

		var mCfg config.Config

		_, err = toml.DecodeFile(cfg.configPath, &mCfg)
		if err != nil {
			return errors.WithMessage(err, "failed to load mCfg file")
		}

		err = c.Subscribe(ctx, "monitor_focus node_add node_remove node_activate node_focus node_flag desktop_focus", func(bytes []byte) {
			q := "query -N -n .window --desktop"

			if cfg.monitor != "" {
				q += fmt.Sprintf(" %v:focused", cfg.monitor)
			}

			var allIDs []bspc.ID
			if err0 := c.Query(q, bspc.ToIDSlice(&allIDs)); err0 != nil {
				log.Fatal(err0)
			}

			var activeID bspc.ID
			if err0 := c.Query("query -N -n", bspc.ToID(&activeID)); err0 != nil {
				log.Fatal(err0)
			}

			if mCfg.MaxWindows > 0 && len(allIDs) > mCfg.MaxWindows {
				allIDs = allIDs[0:mCfg.MaxWindows]
			}

			if len(allIDs) == 0 {
				fmt.Println(config.FormatStringColors(mCfg.EmptyDesktopString, mCfg.EmptyDesktopBgColor, mCfg.EmptyDesktopFgColor, mCfg.EmptyDesktopUlColor))
				return
			}

			separator := config.FormatStringColors(mCfg.SeparatorString, mCfg.SeparatorBgColor, mCfg.SeparatorFgColor, mCfg.SeparatorUlColor)

			labels := make([]string, 0, len(allIDs))
			for _, id := range allIDs {
				var win bspc.Node
				if err := c.Query(fmt.Sprintf("query -T -n %v", id), bspc.ToStruct(&win)); err != nil {
					continue
				}

				if slices.Contains(mCfg.IgnoredClasses, win.Client.ClassName) {
					continue
				}

				isActive := win.ID == activeID

				label := win.Client.ClassName
				if t, ok := mCfg.WindowNicknames[label]; ok {
					label = t
				}
				if mCfg.NameMaxLength > 0 && len(label) > mCfg.NameMaxLength {
					label = label[:mCfg.NameMaxLength]
				}

				flags := ""
				if win.Locked && mCfg.Flags.LockedFlag != "" {
					flags += mCfg.Flags.LockedFlag
				}
				if win.Private && mCfg.Flags.PrivateFlag != "" {
					flags += mCfg.Flags.PrivateFlag
				}
				if win.Sticky && mCfg.Flags.StickyFlag != "" {
					flags += mCfg.Flags.StickyFlag
				}
				if win.Marked && mCfg.Flags.MarkedFlag != "" {
					flags += mCfg.Flags.MarkedFlag
				}

				if len(flags) > 0 {
					label += " " + flags
				}

				label = fmt.Sprintf("%[1]v%v%[1]v", strings.Repeat(" ", mCfg.NamePadding), label)

				label = config.FormatStringColors(
					label,
					mCfg.GetBgColor(isActive, win.Hidden),
					mCfg.GetFgColor(isActive, win.Hidden),
					mCfg.GetUlColor(isActive, win.Hidden),
				)

				if action := mCfg.GetActionLeftClick(isActive, win.Hidden); action != "" {
					label = fmt.Sprintf("%%{A1:%v %v:}%v%%{A}", actions.NormalizeCommand(action), win.ID, label)
				}
				if action := mCfg.GetActionMiddleClick(isActive, win.Hidden); action != "" {
					label = fmt.Sprintf("%%{A2:%v %v:}%v%%{A}", actions.NormalizeCommand(action), win.ID, label)
				}
				if action := mCfg.GetActionRightClick(isActive, win.Hidden); action != "" {
					label = fmt.Sprintf("%%{A3:%v %v:}%v%%{A}", actions.NormalizeCommand(action), win.ID, label)
				}
				if action := mCfg.GetActionScrollUp(isActive, win.Hidden); action != "" {
					label = fmt.Sprintf("%%{A4:%v %v:}%v%%{A}", actions.NormalizeCommand(action), win.ID, label)
				}
				if action := mCfg.GetActionScrollDown(isActive, win.Hidden); action != "" {
					label = fmt.Sprintf("%%{A5:%v %v:}%v%%{A}", actions.NormalizeCommand(action), win.ID, label)
				}

				labels = append(labels, label)
			}
			fmt.Println(strings.Join(labels, separator))
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
