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
	Short: "Generate PDF, XML, or plain PDF from cv.json",
	Long: `Generate output from a CV JSON file.

Formats:
  pdf      Europass PDF with embedded XML attachment (default)
  xml      Standalone Europass XML
  plain    PDF without embedded XML

Examples:
  goepcvcli generate -f cv.json -o cv.pdf
  goepcvcli generate -f cv.json -o cv.xml --format xml
  goepcvcli generate -f cv.json -o cv.pdf --format plain
  goepcvcli generate -f helptech-cv-de.json -o helptech-cv-de.pdf --lang de`,
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := os.ReadFile(inputFile)
		if err != nil {
			return fmt.Errorf("cannot read %s: %w", inputFile, err)
		}

		var cv CV
		if err := json.Unmarshal(data, &cv); err != nil {
			return fmt.Errorf("invalid JSON in %s: %w", inputFile, err)
		}

		if lang, _ := cmd.Flags().GetString("lang"); lang != "" {
			cv.Lang = lang
		}

		format, _ := cmd.Flags().GetString("format")
		switch format {
		case "", "pdf":
			if err := generatePDF(&cv, outputFile); err != nil {
				return fmt.Errorf("PDF generation failed: %w", err)
			}
		case "xml":
			xmlData := toEuropassXML(&cv)
			if err := os.WriteFile(outputFile, xmlData, 0644); err != nil {
				return fmt.Errorf("XML write failed: %w", err)
			}
		case "plain":
			if err := generatePlainPDF(&cv, outputFile); err != nil {
				return fmt.Errorf("PDF generation failed: %w", err)
			}
		default:
			return fmt.Errorf("unknown format %q — use: pdf, xml, plain", format)
		}

		fmt.Printf("Generated %s\n", outputFile)
		return nil
	},
}

func init() {
	generateCmd.Flags().StringVarP(&inputFile, "input", "f", "output/cv.json", "input JSON file")
	generateCmd.Flags().StringVarP(&outputFile, "output", "o", "output/cv.pdf", "output PDF file")
	generateCmd.Flags().String("lang", "", "override CV language (en, de, nl, pt, fr, es)")
	generateCmd.Flags().String("format", "pdf", "output format: pdf (with XML), xml, plain (PDF without XML)")
	rootCmd.AddCommand(generateCmd)
}
