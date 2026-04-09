package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var showSection string

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Display the CV data (or a specific section)",
	Long: `Display the full CV or a specific section as formatted text or JSON.

Sections: personal, headline, experience, education, languages, digital, skills, all

Examples:
  goepcvcli show                          # full CV summary
  goepcvcli show --section experience     # just work experience
  goepcvcli show --section languages      # language table
  goepcvcli show --json                   # raw JSON output
  goepcvcli show --section personal --json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := os.ReadFile(inputFile)
		if err != nil {
			return fmt.Errorf("cannot read %s: %w", inputFile, err)
		}
		var cv CV
		if err := json.Unmarshal(data, &cv); err != nil {
			return fmt.Errorf("invalid JSON: %w", err)
		}

		jsonOut, _ := cmd.Flags().GetBool("json")

		if jsonOut {
			return showJSON(&cv, showSection)
		}
		return showText(&cv, showSection)
	},
}

func showJSON(cv *CV, section string) error {
	var out any
	switch section {
	case "", "all":
		out = cv
	case "personal":
		out = cv.Personal
	case "headline":
		out = map[string]string{"headline": cv.Headline}
	case "experience":
		out = cv.Experience
	case "education":
		out = cv.Education
	case "languages":
		out = cv.Languages
	case "digital":
		out = cv.Digital
	case "skills":
		out = map[string]string{
			"organisational":  cv.Org,
			"communication":   cv.Comm,
			"job_related":     cv.JobRelated,
		}
	default:
		return fmt.Errorf("unknown section %q — use: personal, headline, experience, education, languages, digital, skills", section)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

func showText(cv *CV, section string) error {
	switch section {
	case "", "all":
		showPersonalText(cv)
		showExperienceText(cv)
		showEducationText(cv)
		showLanguagesText(cv)
		showDigitalText(cv)
		showSkillsText(cv)
	case "personal":
		showPersonalText(cv)
	case "headline":
		fmt.Printf("Headline: %s\n", cv.Headline)
	case "experience":
		showExperienceText(cv)
	case "education":
		showEducationText(cv)
	case "languages":
		showLanguagesText(cv)
	case "digital":
		showDigitalText(cv)
	case "skills":
		showSkillsText(cv)
	default:
		return fmt.Errorf("unknown section %q", section)
	}
	return nil
}

func showPersonalText(cv *CV) {
	p := cv.Personal
	fmt.Printf("=== %s %s ===\n", p.FirstName, p.Surname)
	fmt.Printf("Headline: %s\n", cv.Headline)
	fmt.Printf("DOB: %s  Nationality: %s\n", p.DateOfBirth, p.Nationality)
	fmt.Printf("Phone: %s  Email: %s\n", p.Phone, p.Email)
	if p.Website != "" {
		fmt.Printf("Web: %s\n", p.Website)
	}
	if p.GitHub != "" {
		fmt.Printf("GitHub: %s\n", p.GitHub)
	}
	if p.LinkedIn != "" {
		fmt.Printf("LinkedIn: %s\n", p.LinkedIn)
	}
	for _, kv := range p.Extra {
		fmt.Printf("%s: %s\n", kv.Key, kv.Value)
	}
	if p.Address != "" {
		fmt.Printf("Address: %s\n", p.Address)
	}
	fmt.Println()
}

func showExperienceText(cv *CV) {
	fmt.Println("=== WORK EXPERIENCE ===")
	for i, w := range cv.Experience {
		period := w.From
		if w.To != "" {
			period += " - " + w.To
		} else {
			period += " - Present"
		}
		loc := ""
		if w.Location != "" {
			loc = ", " + w.Location
			if w.Country != "" {
				loc += ", " + w.Country
			}
		}
		tags := ""
		if len(w.Tags) > 0 {
			tags = "  [" + strings.Join(w.Tags, ", ") + "]"
		}
		fmt.Printf("[%d] %s | %s @ %s%s%s\n", i, period, w.Title, w.Employer, loc, tags)
		if w.Description != "" {
			fmt.Printf("    %s\n", w.Description)
		}
		fmt.Println()
	}
}

func showEducationText(cv *CV) {
	fmt.Println("=== EDUCATION ===")
	for i, e := range cv.Education {
		period := e.From
		if e.To != "" {
			period += " - " + e.To
		}
		loc := ""
		if e.Location != "" {
			loc = ", " + e.Location
			if e.Country != "" {
				loc += ", " + e.Country
			}
		}
		fmt.Printf("[%d] %s | %s @ %s%s\n", i, period, e.Title, e.Institution, loc)
		if e.Level != "" {
			fmt.Printf("    Level: %s\n", e.Level)
		}
		if e.Description != "" {
			fmt.Printf("    %s\n", e.Description)
		}
		fmt.Println()
	}
}

func showLanguagesText(cv *CV) {
	fmt.Println("=== LANGUAGES ===")
	fmt.Printf("Mother tongue: %s\n", strings.Join(cv.Languages.MotherTongue, ", "))
	fmt.Printf("%-14s  %s  %s  %s  %s  %s\n", "", "Listen", "Read", "SpkProd", "SpkInt", "Write")
	for _, l := range cv.Languages.Foreign {
		fmt.Printf("%-14s  %-6s  %-4s  %-7s  %-6s  %s\n",
			l.Name, l.Listening, l.Reading, l.SpokenProduction, l.SpokenInteraction, l.Writing)
	}
	fmt.Println()
}

func showDigitalText(cv *CV) {
	fmt.Println("=== DIGITAL SKILLS ===")
	fmt.Println(strings.Join(cv.Digital, ", "))
	fmt.Println()
}

func showSkillsText(cv *CV) {
	fmt.Println("=== ADDITIONAL SKILLS ===")
	if cv.Org != "" {
		fmt.Printf("Organisational: %s\n\n", cv.Org)
	}
	if cv.Comm != "" {
		fmt.Printf("Communication: %s\n\n", cv.Comm)
	}
	if cv.JobRelated != "" {
		fmt.Printf("Job-related: %s\n\n", cv.JobRelated)
	}
}

func init() {
	showCmd.Flags().StringVarP(&showSection, "section", "s", "", "section to display")
	showCmd.Flags().StringVarP(&inputFile, "input", "i", "cv.json", "input JSON file")
	showCmd.Flags().Bool("json", false, "output as JSON")
	rootCmd.AddCommand(showCmd)
}
