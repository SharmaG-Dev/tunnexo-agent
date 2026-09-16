package cmd

import (
	"os/signal"
	"syscall"

	"github.com/SharmaG-Dev/Portune-Agent/internal/agent"
	"github.com/spf13/cobra"
)

var localTarget string

var agentCmd = &cobra.Command{
	Use:   "agent --target <local-url>",
	Short: "Start Tunnexo tunnel agent",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := agent.LoadConfig(".env", localTarget)
		if err != nil {
			return err
		}
		ctx, cancel := signal.NotifyContext(cmd.Context(), syscall.SIGINT, syscall.SIGTERM)
		defer cancel()
		return agent.New(config).Run(ctx)
	},
}

func init() {
	rootCmd.AddCommand(agentCmd)
	agentCmd.Flags().StringVar(&localTarget, "target", "", "Local application URL (required)")
	_ = agentCmd.MarkFlagRequired("target")
}
