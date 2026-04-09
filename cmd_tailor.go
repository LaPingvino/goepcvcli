package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var tailorCmd = &cobra.Command{
	Use:   "tailor",
	Short: "Create a job-specific CV variant",
	Long: `Create a tailored CV by filtering and reordering experience entries by tags,
adjusting the headline, and optionally overriding skill descriptions.

The base cv.json is not modified — output goes to a new file.

Examples:
  # Keep only dev-tagged experience, new headline
  goepcvcli tailor --tags dev,go,architecture \
    --headline "Go Developer | Systems Architecture | LLM Workflows" \
    --output dev-cv.json

  # Keep support + leadership roles
  goepcvcli tailor --tags support,leadership,b2b \
    --headline "Technical Support Lead | B2B SaaS | Multilingual" \
    --output support-cv.json

  # Exclude specific tags instead
  goepcvcli tailor --exclude-tags microsoft \
    --output trimmed-cv.json

  # Generate PDF directly
  goepcvcli tailor --tags dev --output dev-cv.json --pdf dev-cv.pdf`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cv, err := loadCV(inputFile)
		if err != nil {
			return err
		}

		output, _ := cmd.Flags().GetString("output")
		pdfOut, _ := cmd.Flags().GetString("pdf")
		headline, _ := cmd.Flags().GetString("headline")
		tagsStr, _ := cmd.Flags().GetString("tags")
		excludeStr, _ := cmd.Flags().GetString("exclude-tags")
		orgSkills, _ := cmd.Flags().GetString("organisational-skills")
		commSkills, _ := cmd.Flags().GetString("communication-skills")
		jobSkills, _ := cmd.Flags().GetString("job-related-skills")

		if output == "" {
			return fmt.Errorf("--output is required")
		}

		// Filter experience by tags
		if tagsStr != "" {
			tags := strings.Split(tagsStr, ",")
			var filtered []Work
			for _, w := range cv.Experience {
				if hasAnyTag(w.Tags, tags) {
					filtered = append(filtered, w)
				}
			}
			cv.Experience = filtered
		}

		if excludeStr != "" {
			excludeTags := strings.Split(excludeStr, ",")
			var filtered []Work
			for _, w := range cv.Experience {
				if !hasAnyTag(w.Tags, excludeTags) {
					filtered = append(filtered, w)
				}
			}
			cv.Experience = filtered
		}

		// Override fields
		if headline != "" {
			cv.Headline = headline
		}
		if orgSkills != "" {
			cv.Org = orgSkills
		}
		if commSkills != "" {
			cv.Comm = commSkills
		}
		if jobSkills != "" {
			cv.JobRelated = jobSkills
		}

		if err := saveCV(output, cv); err != nil {
			return err
		}
		fmt.Printf("Tailored CV written to %s (%d experience entries)\n", output, len(cv.Experience))

		// Optionally generate PDF
		if pdfOut != "" {
			if err := generatePDF(cv, pdfOut); err != nil {
				return fmt.Errorf("PDF generation failed: %w", err)
			}
			fmt.Printf("PDF generated: %s\n", pdfOut)
		}

		return nil
	},
}

func hasAnyTag(entryTags, filterTags []string) bool {
	for _, ft := range filterTags {
		ft = strings.TrimSpace(ft)
		for _, et := range entryTags {
			if strings.EqualFold(et, ft) {
				return true
			}
		}
	}
	return false
}

func init() {
	tailorCmd.Flags().StringVarP(&inputFile, "input", "f", "cv.json", "base CV JSON file")
	tailorCmd.Flags().StringP("output", "o", "", "output JSON file (required)")
	tailorCmd.Flags().String("pdf", "", "also generate PDF")
	tailorCmd.Flags().String("headline", "", "override headline")
	tailorCmd.Flags().String("tags", "", "include experience with any of these tags (comma-separated)")
	tailorCmd.Flags().String("exclude-tags", "", "exclude experience with any of these tags")
	tailorCmd.Flags().String("organisational-skills", "", "override organisational skills text")
	tailorCmd.Flags().String("communication-skills", "", "override communication skills text")
	tailorCmd.Flags().String("job-related-skills", "", "override job-related skills text")
	rootCmd.AddCommand(tailorCmd)
}
