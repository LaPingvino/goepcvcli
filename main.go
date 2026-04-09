package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "goepcvcli",
	Short: "Europass CV tool — structure, tailor, and generate CVs",
	Long: `A CLI tool for managing Europass-format CVs.

Store your CV as structured JSON, tailor it per job application,
and generate professional PDFs with embedded Europass XML.

Use -i for interactive mode with guided prompts.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if interactive {
			return interactiveMain()
		}
		return cmd.Help()
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&interactive, "interactive", "i", false, "interactive mode with guided prompts")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
