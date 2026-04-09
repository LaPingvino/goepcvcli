package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove an entry from a CV section",
	Long:  `Remove work experience, education, language, digital skill, or contact entry by index.`,
}

var removeWorkCmd = &cobra.Command{
	Use:   "work [index]",
	Short: "Remove a work experience entry by index",
	Long: `Remove a work experience entry. Use 'show --section experience' to see indices.

Examples:
  goepcvcli remove work 3     # remove entry at index 3`,
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
		removed := cv.Experience[idx]
		cv.Experience = append(cv.Experience[:idx], cv.Experience[idx+1:]...)
		if err := saveCV(inputFile, cv); err != nil {
			return err
		}
		fmt.Printf("Removed work[%d]: %s @ %s\n", idx, removed.Title, removed.Employer)
		return nil
	},
}

var removeEducationCmd = &cobra.Command{
	Use:   "education [index]",
	Short: "Remove an education entry by index",
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
		removed := cv.Education[idx]
		cv.Education = append(cv.Education[:idx], cv.Education[idx+1:]...)
		if err := saveCV(inputFile, cv); err != nil {
			return err
		}
		fmt.Printf("Removed education[%d]: %s @ %s\n", idx, removed.Title, removed.Institution)
		return nil
	},
}

var removeLanguageCmd = &cobra.Command{
	Use:   "language [name]",
	Short: "Remove a foreign language by name",
	Long: `Remove a foreign language entry by name (case-insensitive).

Examples:
  goepcvcli remove language Afrikaans`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cv, err := loadCV(inputFile)
		if err != nil {
			return err
		}
		name := strings.ToLower(args[0])
		found := false
		for i, l := range cv.Languages.Foreign {
			if strings.ToLower(l.Name) == name {
				cv.Languages.Foreign = append(cv.Languages.Foreign[:i], cv.Languages.Foreign[i+1:]...)
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("language %q not found", args[0])
		}
		if err := saveCV(inputFile, cv); err != nil {
			return err
		}
		fmt.Printf("Removed language: %s\n", args[0])
		return nil
	},
}

var removeSkillCmd = &cobra.Command{
	Use:   "skill [name]",
	Short: "Remove a digital skill by name",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cv, err := loadCV(inputFile)
		if err != nil {
			return err
		}
		name := strings.ToLower(args[0])
		found := false
		for i, s := range cv.Digital {
			if strings.ToLower(s) == name {
				cv.Digital = append(cv.Digital[:i], cv.Digital[i+1:]...)
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("skill %q not found", args[0])
		}
		if err := saveCV(inputFile, cv); err != nil {
			return err
		}
		fmt.Printf("Removed skill: %s\n", args[0])
		return nil
	},
}

var removeContactCmd = &cobra.Command{
	Use:   "contact [key]",
	Short: "Remove extra contact info by key",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cv, err := loadCV(inputFile)
		if err != nil {
			return err
		}
		key := strings.ToLower(args[0])
		found := false
		for i, kv := range cv.Personal.Extra {
			if strings.ToLower(kv.Key) == key {
				cv.Personal.Extra = append(cv.Personal.Extra[:i], cv.Personal.Extra[i+1:]...)
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("contact %q not found", args[0])
		}
		if err := saveCV(inputFile, cv); err != nil {
			return err
		}
		fmt.Printf("Removed contact: %s\n", args[0])
		return nil
	},
}

func init() {
	removeWorkCmd.Flags().StringVarP(&inputFile, "input", "f", "cv.json", "input JSON file")
	removeEducationCmd.Flags().StringVarP(&inputFile, "input", "f", "cv.json", "input JSON file")
	removeLanguageCmd.Flags().StringVarP(&inputFile, "input", "f", "cv.json", "input JSON file")
	removeSkillCmd.Flags().StringVarP(&inputFile, "input", "f", "cv.json", "input JSON file")
	removeContactCmd.Flags().StringVarP(&inputFile, "input", "f", "cv.json", "input JSON file")

	removeCmd.AddCommand(removeWorkCmd, removeEducationCmd, removeLanguageCmd, removeSkillCmd, removeContactCmd)
	rootCmd.AddCommand(removeCmd)
}
