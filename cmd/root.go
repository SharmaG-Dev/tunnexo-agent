package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tunnexo",
	Short: "Tunnexo secure tunnel agent",
	Long: `
Tunnexo Agent connects local applications
to the Tunnexo public tunnel network.
`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
