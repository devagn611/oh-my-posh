package cli

import (
	"fmt"
	stdruntime "runtime"
	"slices"

	"github.com/jandedobbeleer/oh-my-posh/src/config"
	"github.com/jandedobbeleer/oh-my-posh/src/runtime"
	"github.com/jandedobbeleer/oh-my-posh/src/terminal"
	"github.com/jandedobbeleer/oh-my-posh/src/upgrade"
	"github.com/spf13/cobra"
)

var force bool

// noticeCmd represents the get command
var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade when a new version is available.",
	Long:  "Upgrade when a new version is available.",
	Args:  cobra.NoArgs,
	Run: func(_ *cobra.Command, _ []string) {
		supportedPlatforms := []string{
			runtime.WINDOWS,
			runtime.DARWIN,
			runtime.LINUX,
		}

		if !slices.Contains(supportedPlatforms, stdruntime.GOOS) {
			fmt.Print("\n⚠️ upgrade is not supported on this platform\n\n")
			return
		}

		env := &runtime.Terminal{
			CmdFlags: &runtime.Flags{},
		}
		env.Init()
		defer env.Close()

		terminal.Init(env.Shell())
		fmt.Print(terminal.StartProgress())

		defer fmt.Print(terminal.StopProgress())

		if force {
			upgrade.Run(env)
			return
		}

		cfg := config.Load(env)

		if _, hasNotice := upgrade.Notice(env, true); !hasNotice {
			if !cfg.DisableNotice {
				fmt.Print("\n✅  no new version available\n\n")
			}
			return
		}

		upgrade.Run(env)
	},
}

func init() {
	upgradeCmd.Flags().BoolVarP(&force, "force", "f", false, "force the upgrade even if the version is up to date")
	RootCmd.AddCommand(upgradeCmd)
}