package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a starter cv.json from the built-in template",
	RunE: func(cmd *cobra.Command, args []string) error {
		path := "output/cv.json"
		if len(args) > 0 {
			path = args[0]
		}

		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf("%s already exists — use a different name or delete it first", path)
		}

		// Ensure output directory exists
		if dir := dirOf(path); dir != "." {
			os.MkdirAll(dir, 0755)
		}

		data, err := json.MarshalIndent(templateCV(), "", "  ")
		if err != nil {
			return err
		}

		if err := os.WriteFile(path, data, 0644); err != nil {
			return err
		}
		fmt.Printf("Created %s — edit it with your details, then run: goepcvcli generate\n", path)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
