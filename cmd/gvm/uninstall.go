package gvm

import (
	"fmt"

	"github.com/ibreez3/gvm/internal/core"
	"github.com/spf13/cobra"
)

var (
	uninstallBelow     string
	uninstallPattern   string
	uninstallKeep      int
	uninstallAll       bool
	uninstallKeepCurrent bool
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall [version]",
	Short: "Uninstall Go versions",
	Long: `Uninstall Go versions, supports single or batch uninstall.

Examples:
  gvm uninstall 1.21.0           # Uninstall a single version
  gvm uninstall --below 1.22     # Uninstall all versions below 1.22
  gvm uninstall --pattern "1.21.*"  # Uninstall all 1.21.x versions
  gvm uninstall --keep 2         # Keep only the latest 2 versions
  gvm uninstall --all            # Uninstall all versions

Note: Batch uninstall automatically skips the currently active version.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Batch mode: no version specified, use flags
		if len(args) == 0 {
			spec := &core.UninstallBatchSpec{
				Below:       uninstallBelow,
				Pattern:     uninstallPattern,
				Keep:        uninstallKeep,
				All:         uninstallAll,
				KeepCurrent: uninstallKeepCurrent,
			}

			// Validate that exactly one batch option is specified
			count := 0
			if uninstallBelow != "" {
				count++
			}
			if uninstallPattern != "" {
				count++
			}
			if uninstallKeep > 0 {
				count++
			}
			if uninstallAll {
				count++
			}

			if count != 1 {
				return fmt.Errorf("please specify one batch uninstall option: --below, --pattern, --keep, or --all")
			}

			uninstalled, err := core.UninstallBatch(spec)
			if err != nil {
				return err
			}
			fmt.Printf("\nSuccessfully uninstalled %d version(s)\n", len(uninstalled))
			return nil
		}

		// Single version mode
		if uninstallBelow != "" || uninstallPattern != "" || uninstallKeep > 0 || uninstallAll {
			return fmt.Errorf("cannot specify version and batch options at the same time")
		}

		return core.UninstallVersion(args[0])
	},
}

func init() {
	uninstallCmd.Flags().StringVar(&uninstallBelow, "below", "", "Uninstall all versions below this version")
	uninstallCmd.Flags().StringVar(&uninstallPattern, "pattern", "", "Uninstall versions matching this pattern (supports wildcards, e.g., 1.21.*)")
	uninstallCmd.Flags().IntVar(&uninstallKeep, "keep", 0, "Keep only the latest N versions")
	uninstallCmd.Flags().BoolVar(&uninstallAll, "all", false, "Uninstall all versions")
	uninstallCmd.Flags().BoolVarP(&uninstallKeepCurrent, "keep-current", "c", true, "Keep the currently active version")
	rootCmd.AddCommand(uninstallCmd)
}
