package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing entry in a CV section",
	Long:  `Modify fields on an existing work experience or education entry by index.`,
}

var updateWorkCmd = &cobra.Command{
	Use:   "work [index]",
	Short: "Update a work experience entry",
	Long: `Update fields on an existing work entry. Only specified flags are changed.

Examples:
  goepcvcli update work 0 --description "New description text"
  goepcvcli update work 0 --to "DEC 2025" --tags "dev,go,llm"
  goepcvcli update work 0 --title "Senior Developer" --employer "New Name"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cv, err := loadCV(inputFile)
		if err != nil {
			return err
		}
		idx, err := strconv.Atoi(args[0])
		if err != nil || idx < 0 || idx >= len(cv.Experience) {
			return fmt.Errorf("invalid index %q — valid range: 0-%d", args[0], len(cv.Experience)-1)
		}

		w := &cv.Experience[idx]
		if cmd.Flags().Changed("title") {
			w.Title, _ = cmd.Flags().GetString("title")
		}
		if cmd.Flags().Changed("employer") {
			w.Employer, _ = cmd.Flags().GetString("employer")
		}
		if cmd.Flags().Changed("from") {
			w.From, _ = cmd.Flags().GetString("from")
		}
		if cmd.Flags().Changed("to") {
			w.To, _ = cmd.Flags().GetString("to")
		}
		if cmd.Flags().Changed("description") {
			w.Description, _ = cmd.Flags().GetString("description")
		}
		if cmd.Flags().Changed("location") {
			w.Location, _ = cmd.Flags().GetString("location")
		}
		if cmd.Flags().Changed("country") {
			w.Country, _ = cmd.Flags().GetString("country")
		}
		if cmd.Flags().Changed("tags") {
			tagsStr, _ := cmd.Flags().GetString("tags")
			w.Tags = strings.Split(tagsStr, ",")
		}

		if err := saveCV(inputFile, cv); err != nil {
			return err
		}
		fmt.Printf("Updated work[%d]: %s @ %s\n", idx, w.Title, w.Employer)
		return nil
	},
}

var updateEducationCmd = &cobra.Command{
	Use:   "education [index]",
	Short: "Update an education entry",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cv, err := loadCV(inputFile)
		if err != nil {
			return err
		}
		idx, err := strconv.Atoi(args[0])
		if err != nil || idx < 0 || idx >= len(cv.Education) {
			return fmt.Errorf("invalid index %q — valid range: 0-%d", args[0], len(cv.Education)-1)
		}

		e := &cv.Education[idx]
		if cmd.Flags().Changed("title") {
			e.Title, _ = cmd.Flags().GetString("title")
		}
		if cmd.Flags().Changed("institution") {
			e.Institution, _ = cmd.Flags().GetString("institution")
		}
		if cmd.Flags().Changed("from") {
			e.From, _ = cmd.Flags().GetString("from")
		}
		if cmd.Flags().Changed("to") {
			e.To, _ = cmd.Flags().GetString("to")
		}
		if cmd.Flags().Changed("description") {
			e.Description, _ = cmd.Flags().GetString("description")
		}
		if cmd.Flags().Changed("location") {
			e.Location, _ = cmd.Flags().GetString("location")
		}
		if cmd.Flags().Changed("country") {
			e.Country, _ = cmd.Flags().GetString("country")
		}
		if cmd.Flags().Changed("level") {
			e.Level, _ = cmd.Flags().GetString("level")
		}

		if err := saveCV(inputFile, cv); err != nil {
			return err
		}
		fmt.Printf("Updated education[%d]: %s @ %s\n", idx, e.Title, e.Institution)
		return nil
	},
}

var updateLanguageCmd = &cobra.Command{
	Use:   "language [name]",
	Short: "Update a foreign language's CEFR levels",
	Long: `Update CEFR levels for an existing language (by name, case-insensitive).

Examples:
  goepcvcli update language Portuguese --listening C2 --spoken-interaction C1
  goepcvcli update language Japanese --all B2`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cv, err := loadCV(inputFile)
		if err != nil {
			return err
		}

		name := strings.ToLower(args[0])
		var fl *ForeignLang
		for i := range cv.Languages.Foreign {
			if strings.ToLower(cv.Languages.Foreign[i].Name) == name {
				fl = &cv.Languages.Foreign[i]
				break
			}
		}
		if fl == nil {
			return fmt.Errorf("language %q not found", args[0])
		}

		allLevel, _ := cmd.Flags().GetString("all")
		if allLevel != "" {
			fl.Listening = allLevel
			fl.Reading = allLevel
			fl.SpokenProduction = allLevel
			fl.SpokenInteraction = allLevel
			fl.Writing = allLevel
		}
		if cmd.Flags().Changed("listening") {
			fl.Listening, _ = cmd.Flags().GetString("listening")
		}
		if cmd.Flags().Changed("reading") {
			fl.Reading, _ = cmd.Flags().GetString("reading")
		}
		if cmd.Flags().Changed("spoken-production") {
			fl.SpokenProduction, _ = cmd.Flags().GetString("spoken-production")
		}
		if cmd.Flags().Changed("spoken-interaction") {
			fl.SpokenInteraction, _ = cmd.Flags().GetString("spoken-interaction")
		}
		if cmd.Flags().Changed("writing") {
			fl.Writing, _ = cmd.Flags().GetString("writing")
		}

		if err := saveCV(inputFile, cv); err != nil {
			return err
		}
		fmt.Printf("Updated language: %s\n", fl.Name)
		return nil
	},
}

func init() {
	// Work update flags
	updateWorkCmd.Flags().String("title", "", "job title")
	updateWorkCmd.Flags().String("employer", "", "employer name")
	updateWorkCmd.Flags().String("from", "", "start date")
	updateWorkCmd.Flags().String("to", "", "end date")
	updateWorkCmd.Flags().String("description", "", "role description")
	updateWorkCmd.Flags().String("location", "", "city")
	updateWorkCmd.Flags().String("country", "", "country")
	updateWorkCmd.Flags().String("tags", "", "comma-separated tags")
	updateWorkCmd.Flags().StringVarP(&inputFile, "input", "f", "cv.json", "input JSON file")

	// Education update flags
	updateEducationCmd.Flags().String("title", "", "degree/certificate title")
	updateEducationCmd.Flags().String("institution", "", "institution name")
	updateEducationCmd.Flags().String("from", "", "start date")
	updateEducationCmd.Flags().String("to", "", "end date")
	updateEducationCmd.Flags().String("description", "", "description")
	updateEducationCmd.Flags().String("location", "", "city")
	updateEducationCmd.Flags().String("country", "", "country")
	updateEducationCmd.Flags().String("level", "", "qualification level")
	updateEducationCmd.Flags().StringVarP(&inputFile, "input", "f", "cv.json", "input JSON file")

	// Language update flags
	updateLanguageCmd.Flags().String("all", "", "set all CEFR levels at once")
	updateLanguageCmd.Flags().String("listening", "", "CEFR level")
	updateLanguageCmd.Flags().String("reading", "", "CEFR level")
	updateLanguageCmd.Flags().String("spoken-production", "", "CEFR level")
	updateLanguageCmd.Flags().String("spoken-interaction", "", "CEFR level")
	updateLanguageCmd.Flags().String("writing", "", "CEFR level")
	updateLanguageCmd.Flags().StringVarP(&inputFile, "input", "f", "cv.json", "input JSON file")

	updateCmd.AddCommand(updateWorkCmd, updateEducationCmd, updateLanguageCmd)
	rootCmd.AddCommand(updateCmd)
}
