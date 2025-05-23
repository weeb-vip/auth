package commands

import (
	"github.com/spf13/cobra"
)

// Execute parses and executes commands based on argument.
func Execute() {
	var rootCmd = &cobra.Command{
		Use:   "auth",
		Short: "CLI for executing tasks tasks",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	configureServerCommand(rootCmd)
	configureMigrateCommand(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		rootCmd.PrintErr(err)
	}
}
