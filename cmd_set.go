package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func dirOf(path string) string {
	return filepath.Dir(path)
}

var setCmd = &cobra.Command{
	Use:   "set [field] [value]",
	Short: "Set a top-level CV field",
	Long: `Set a simple text field on the CV.

Fields: headline, address, phone, email, website, github, linkedin,
        first_name, surname, date_of_birth, nationality,
        organisational_skills, communication_skills, job_related_skills

Examples:
  goepcvcli set headline "Developer Tooling & LLM Workflows | Go, TypeScript, Rust"
  goepcvcli set phone "+351 913044570"
  goepcvcli set organisational_skills "Understands what needs to be arranged..."`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		field, value := args[0], args[1]

		cv, err := loadCV(inputFile)
		if err != nil {
			return err
		}

		switch field {
		case "headline":
			cv.Headline = value
		case "address":
			cv.Personal.Address = value
		case "phone":
			cv.Personal.Phone = value
		case "email":
			cv.Personal.Email = value
		case "website":
			cv.Personal.Website = value
		case "github":
			cv.Personal.GitHub = value
		case "linkedin":
			cv.Personal.LinkedIn = value
		case "first_name":
			cv.Personal.FirstName = value
		case "surname":
			cv.Personal.Surname = value
		case "date_of_birth":
			cv.Personal.DateOfBirth = value
		case "nationality":
			cv.Personal.Nationality = value
		case "organisational_skills":
			cv.Org = value
		case "communication_skills":
			cv.Comm = value
		case "job_related_skills":
			cv.JobRelated = value
		default:
			return fmt.Errorf("unknown field %q — use: headline, address, phone, email, website, github, linkedin, first_name, surname, date_of_birth, nationality, organisational_skills, communication_skills, job_related_skills", field)
		}

		if err := saveCV(inputFile, cv); err != nil {
			return err
		}
		fmt.Printf("Set %s = %q\n", field, value)
		return nil
	},
}

func loadCV(path string) (*CV, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read %s: %w", path, err)
	}
	var cv CV
	if err := json.Unmarshal(data, &cv); err != nil {
		return nil, fmt.Errorf("invalid JSON in %s: %w", path, err)
	}
	return &cv, nil
}

func saveCV(path string, cv *CV) error {
	if dir := dirOf(path); dir != "." {
		os.MkdirAll(dir, 0755)
	}
	data, err := json.MarshalIndent(cv, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func init() {
	setCmd.Flags().StringVarP(&inputFile, "input", "f", "output/cv.json", "input JSON file")
	rootCmd.AddCommand(setCmd)
}
