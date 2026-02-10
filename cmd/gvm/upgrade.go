package gvm

import (
	"fmt"

	"github.com/ibreez3/gvm/internal/core"
	"github.com/spf13/cobra"
)

var (
	upgradeUse   bool
	upgradeYes   bool
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade [version]",
	Short: "Upgrade Go minor version to latest patch version",
	Long: `Upgrade a Go minor version to the latest patch version.

Examples:
  gvm upgrade 1.25    # Upgrade to latest 1.25.x
  gvm upgrade go1.25  # Same as above (supports go prefix)`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		version, err := core.UpgradeVersion(args[0])
		if err != nil {
			return err
		}

		// Automatically use the newly upgraded version if requested
		if upgradeUse {
			if err := core.UseVersion(version); err != nil {
				return err
			}
			fmt.Printf("Switched to go%s\n", version)
		}

		return nil
	},
}

func init() {
	upgradeCmd.Flags().BoolVarP(&upgradeUse, "use", "u", false, "Automatically switch to new version after upgrade")
	upgradeCmd.Flags().BoolVarP(&upgradeYes, "yes", "y", false, "Auto-confirm upgrade")
	rootCmd.AddCommand(upgradeCmd)
}
