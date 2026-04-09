package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an entry to a CV section",
	Long:  `Add work experience, education, language, digital skill, or contact info.`,
}

var addWorkCmd = &cobra.Command{
	Use:   "work",
	Short: "Add a work experience entry",
	Long: `Add a work experience entry to the CV.

Examples:
  goepcvcli add work --title "Software Engineer" --employer "Acme Corp" \
    --from "JAN 2024" --description "Building distributed systems in Go" \
    --location "Remote" --country "Portugal" --tags "dev,go,distributed"

  goepcvcli add work --title "President" --employer "NGO" \
    --from "SEP 2019" --to "SEP 2020" --description "Leadership role"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cv, err := loadCV(inputFile)
		if err != nil {
			return err
		}

		title, _ := cmd.Flags().GetString("title")
		employer, _ := cmd.Flags().GetString("employer")
		from, _ := cmd.Flags().GetString("from")
		to, _ := cmd.Flags().GetString("to")
		desc, _ := cmd.Flags().GetString("description")
		location, _ := cmd.Flags().GetString("location")
		country, _ := cmd.Flags().GetString("country")
		tagsStr, _ := cmd.Flags().GetString("tags")
		pos, _ := cmd.Flags().GetInt("position")

		if title == "" || from == "" {
			return fmt.Errorf("--title and --from are required")
		}

		var tags []string
		if tagsStr != "" {
			tags = strings.Split(tagsStr, ",")
		}

		w := Work{
			From:        from,
			To:          to,
			Title:       title,
			Employer:    employer,
			Location:    location,
			Country:     country,
			Description: desc,
			Tags:        tags,
		}

		if pos >= 0 && pos < len(cv.Experience) {
			// Insert at position
			cv.Experience = append(cv.Experience[:pos], append([]Work{w}, cv.Experience[pos:]...)...)
		} else {
			cv.Experience = append(cv.Experience, w)
		}

		if err := saveCV(inputFile, cv); err != nil {
			return err
		}
		fmt.Printf("Added work: %s @ %s (%s)\n", title, employer, from)
		return nil
	},
}

var addEducationCmd = &cobra.Command{
	Use:   "education",
	Short: "Add an education entry",
	Long: `Add an education/training entry to the CV.

Examples:
  goepcvcli add education --title "Computer Science" --institution "MIT" \
    --from "2020" --to "2024" --level "Bachelor" \
    --description "Focus on distributed systems"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cv, err := loadCV(inputFile)
		if err != nil {
			return err
		}

		title, _ := cmd.Flags().GetString("title")
		institution, _ := cmd.Flags().GetString("institution")
		from, _ := cmd.Flags().GetString("from")
		to, _ := cmd.Flags().GetString("to")
		desc, _ := cmd.Flags().GetString("description")
		location, _ := cmd.Flags().GetString("location")
		country, _ := cmd.Flags().GetString("country")
		level, _ := cmd.Flags().GetString("level")

		if title == "" || from == "" {
			return fmt.Errorf("--title and --from are required")
		}

		e := Education{
			From:        from,
			To:          to,
			Title:       title,
			Institution: institution,
			Location:    location,
			Country:     country,
			Level:       level,
			Description: desc,
		}

		cv.Education = append(cv.Education, e)
		if err := saveCV(inputFile, cv); err != nil {
			return err
		}
		fmt.Printf("Added education: %s @ %s\n", title, institution)
		return nil
	},
}

var addLanguageCmd = &cobra.Command{
	Use:   "language",
	Short: "Add a foreign language",
	Long: `Add a foreign language with CEFR levels.

Examples:
  goepcvcli add language --name "Japanese" \
    --listening B1 --reading B2 --spoken-production A2 \
    --spoken-interaction B1 --writing A2

  goepcvcli add language --name "Russian" --all B1`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cv, err := loadCV(inputFile)
		if err != nil {
			return err
		}

		name, _ := cmd.Flags().GetString("name")
		allLevel, _ := cmd.Flags().GetString("all")
		listening, _ := cmd.Flags().GetString("listening")
		reading, _ := cmd.Flags().GetString("reading")
		spProd, _ := cmd.Flags().GetString("spoken-production")
		spInt, _ := cmd.Flags().GetString("spoken-interaction")
		writing, _ := cmd.Flags().GetString("writing")

		if name == "" {
			return fmt.Errorf("--name is required")
		}

		if allLevel != "" {
			if listening == "" { listening = allLevel }
			if reading == "" { reading = allLevel }
			if spProd == "" { spProd = allLevel }
			if spInt == "" { spInt = allLevel }
			if writing == "" { writing = allLevel }
		}

		fl := ForeignLang{
			Name:              name,
			Listening:         listening,
			Reading:           reading,
			SpokenProduction:  spProd,
			SpokenInteraction: spInt,
			Writing:           writing,
		}

		cv.Languages.Foreign = append(cv.Languages.Foreign, fl)
		if err := saveCV(inputFile, cv); err != nil {
			return err
		}
		fmt.Printf("Added language: %s\n", name)
		return nil
	},
}

var addSkillCmd = &cobra.Command{
	Use:   "skill [skill-name]",
	Short: "Add a digital skill",
	Long: `Add one or more digital skills.

Examples:
  goepcvcli add skill Kubernetes
  goepcvcli add skill "Claude Code" "LLM Workflows" "Prompt Engineering"`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cv, err := loadCV(inputFile)
		if err != nil {
			return err
		}

		cv.Digital = append(cv.Digital, args...)
		if err := saveCV(inputFile, cv); err != nil {
			return err
		}
		fmt.Printf("Added digital skills: %s\n", strings.Join(args, ", "))
		return nil
	},
}

var addContactCmd = &cobra.Command{
	Use:   "contact [key] [value]",
	Short: "Add extra contact info (Telegram, Matrix, etc.)",
	Long: `Add an extra contact method to the personal section.

Examples:
  goepcvcli add contact Matrix "@joop:chat.kiefte.eu"
  goepcvcli add contact Mastodon "@joop@poliglota.social.br"`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cv, err := loadCV(inputFile)
		if err != nil {
			return err
		}

		cv.Personal.Extra = append(cv.Personal.Extra, KV{Key: args[0], Value: args[1]})
		if err := saveCV(inputFile, cv); err != nil {
			return err
		}
		fmt.Printf("Added contact: %s = %s\n", args[0], args[1])
		return nil
	},
}

func init() {
	// Work flags
	addWorkCmd.Flags().String("title", "", "job title (required)")
	addWorkCmd.Flags().String("employer", "", "employer name")
	addWorkCmd.Flags().String("from", "", "start date, e.g. 'JAN 2024' (required)")
	addWorkCmd.Flags().String("to", "", "end date (omit for current)")
	addWorkCmd.Flags().String("description", "", "role description")
	addWorkCmd.Flags().String("location", "", "city")
	addWorkCmd.Flags().String("country", "", "country")
	addWorkCmd.Flags().String("tags", "", "comma-separated tags for tailoring")
	addWorkCmd.Flags().IntP("position", "p", -1, "insert position (0=top, default=end)")
	addWorkCmd.Flags().StringVarP(&inputFile, "input", "f", "cv.json", "input JSON file")

	// Education flags
	addEducationCmd.Flags().String("title", "", "degree/certificate title (required)")
	addEducationCmd.Flags().String("institution", "", "school/institution name")
	addEducationCmd.Flags().String("from", "", "start date (required)")
	addEducationCmd.Flags().String("to", "", "end date")
	addEducationCmd.Flags().String("description", "", "description")
	addEducationCmd.Flags().String("location", "", "city")
	addEducationCmd.Flags().String("country", "", "country")
	addEducationCmd.Flags().String("level", "", "EQF level or national classification")
	addEducationCmd.Flags().StringVarP(&inputFile, "input", "f", "cv.json", "input JSON file")

	// Language flags
	addLanguageCmd.Flags().String("name", "", "language name (required)")
	addLanguageCmd.Flags().String("all", "", "set all CEFR levels at once")
	addLanguageCmd.Flags().String("listening", "", "CEFR level A1-C2")
	addLanguageCmd.Flags().String("reading", "", "CEFR level A1-C2")
	addLanguageCmd.Flags().String("spoken-production", "", "CEFR level A1-C2")
	addLanguageCmd.Flags().String("spoken-interaction", "", "CEFR level A1-C2")
	addLanguageCmd.Flags().String("writing", "", "CEFR level A1-C2")
	addLanguageCmd.Flags().StringVarP(&inputFile, "input", "f", "cv.json", "input JSON file")

	// Skill — inputFile via persistent flag on parent
	addSkillCmd.Flags().StringVarP(&inputFile, "input", "f", "cv.json", "input JSON file")

	// Contact
	addContactCmd.Flags().StringVarP(&inputFile, "input", "f", "cv.json", "input JSON file")

	addCmd.AddCommand(addWorkCmd, addEducationCmd, addLanguageCmd, addSkillCmd, addContactCmd)
	rootCmd.AddCommand(addCmd)
}
