package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	inputFile  string
	outputFile string
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a PDF from cv.json",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := os.ReadFile(inputFile)
		if err != nil {
			return fmt.Errorf("cannot read %s: %w", inputFile, err)
		}

		var cv CV
		if err := json.Unmarshal(data, &cv); err != nil {
			return fmt.Errorf("invalid JSON in %s: %w", inputFile, err)
		}

		if err := generatePDF(&cv, outputFile); err != nil {
			return fmt.Errorf("PDF generation failed: %w", err)
		}

		fmt.Printf("Generated %s\n", outputFile)
		return nil
	},
}

func init() {
	generateCmd.Flags().StringVarP(&inputFile, "input", "i", "cv.json", "input JSON file")
	generateCmd.Flags().StringVarP(&outputFile, "output", "o", "cv.pdf", "output PDF file")
	rootCmd.AddCommand(generateCmd)
}
