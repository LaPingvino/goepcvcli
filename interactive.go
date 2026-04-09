package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var interactive bool

var scanner = bufio.NewScanner(os.Stdin)

// prompt prints a label and reads a line from stdin. Returns empty string on EOF.
func prompt(label string) string {
	fmt.Printf("%s: ", label)
	if scanner.Scan() {
		return strings.TrimSpace(scanner.Text())
	}
	return ""
}

// promptDefault prints a label with a default value hint and reads a line.
// Returns the default if the user enters nothing.
func promptDefault(label, def string) string {
	if def != "" {
		fmt.Printf("%s [%s]: ", label, def)
	} else {
		fmt.Printf("%s: ", label)
	}
	if scanner.Scan() {
		v := strings.TrimSpace(scanner.Text())
		if v == "" {
			return def
		}
		return v
	}
	return def
}

// promptRequired keeps asking until the user provides a non-empty value.
func promptRequired(label string) string {
	for {
		v := prompt(label + " (required)")
		if v != "" {
			return v
		}
		fmt.Println("  This field is required, please enter a value.")
	}
}

// promptChoice shows numbered options and returns the index chosen. Returns -1 on cancel.
func promptChoice(label string, options []string) int {
	fmt.Println(label)
	for i, opt := range options {
		fmt.Printf("  %d) %s\n", i+1, opt)
	}
	fmt.Printf("Choice (1-%d, or 'q' to cancel): ", len(options))
	if scanner.Scan() {
		v := strings.TrimSpace(scanner.Text())
		if v == "q" || v == "" {
			return -1
		}
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 || n > len(options) {
			fmt.Println("  Invalid choice.")
			return -1
		}
		return n - 1
	}
	return -1
}

// promptYesNo asks a yes/no question, defaulting to the given value.
func promptYesNo(label string, def bool) bool {
	hint := "y/N"
	if def {
		hint = "Y/n"
	}
	fmt.Printf("%s [%s]: ", label, hint)
	if scanner.Scan() {
		v := strings.ToLower(strings.TrimSpace(scanner.Text()))
		switch v {
		case "y", "yes":
			return true
		case "n", "no":
			return false
		}
	}
	return def
}

// promptCEFR prompts for a CEFR level with validation.
func promptCEFR(label, def string) string {
	valid := map[string]bool{"A1": true, "A2": true, "B1": true, "B2": true, "C1": true, "C2": true}
	for {
		v := promptDefault(label+" (A1/A2/B1/B2/C1/C2)", def)
		if v == "" {
			return def
		}
		upper := strings.ToUpper(v)
		if valid[upper] {
			return upper
		}
		fmt.Println("  Please enter a valid CEFR level: A1, A2, B1, B2, C1, or C2")
	}
}

// interactiveMain runs the top-level interactive menu when no subcommand is given.
func interactiveMain() error {
	fmt.Println("=== Europass CV Manager — Interactive Mode ===")
	fmt.Println()

	for {
		options := []string{
			"Show CV",
			"Set a field (name, headline, phone, etc.)",
			"Add entry (work, education, language, skill, contact)",
			"Update entry (work, education, language)",
			"Remove entry",
			"Tailor CV for a job application",
			"Generate PDF",
			"Quit",
		}
		choice := promptChoice("\nWhat would you like to do?", options)

		switch choice {
		case 0:
			interactiveShow()
		case 1:
			interactiveSet()
		case 2:
			interactiveAdd()
		case 3:
			interactiveUpdate()
		case 4:
			interactiveRemove()
		case 5:
			interactiveTailor()
		case 6:
			interactiveGenerate()
		case 7, -1:
			fmt.Println("Bye!")
			return nil
		}
	}
}

func interactiveShow() {
	sections := []string{"All (full CV)", "Personal info", "Headline", "Experience", "Education", "Languages", "Digital skills", "Other skills"}
	sectionKeys := []string{"all", "personal", "headline", "experience", "education", "languages", "digital", "skills"}

	choice := promptChoice("Which section?", sections)
	if choice < 0 {
		return
	}

	cv, err := loadCV(inputFile)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println()
	showText(cv, sectionKeys[choice])
}

func interactiveSet() {
	fields := []string{
		"headline", "first_name", "surname", "phone", "email",
		"address", "website", "github", "linkedin",
		"date_of_birth", "nationality",
		"organisational_skills", "communication_skills", "job_related_skills",
	}

	choice := promptChoice("Which field?", fields)
	if choice < 0 {
		return
	}
	field := fields[choice]

	cv, err := loadCV(inputFile)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Show current value
	var current string
	switch field {
	case "headline":
		current = cv.Headline
	case "first_name":
		current = cv.Personal.FirstName
	case "surname":
		current = cv.Personal.Surname
	case "phone":
		current = cv.Personal.Phone
	case "email":
		current = cv.Personal.Email
	case "address":
		current = cv.Personal.Address
	case "website":
		current = cv.Personal.Website
	case "github":
		current = cv.Personal.GitHub
	case "linkedin":
		current = cv.Personal.LinkedIn
	case "date_of_birth":
		current = cv.Personal.DateOfBirth
	case "nationality":
		current = cv.Personal.Nationality
	case "organisational_skills":
		current = cv.Org
	case "communication_skills":
		current = cv.Comm
	case "job_related_skills":
		current = cv.JobRelated
	}

	value := promptDefault(fmt.Sprintf("New value for %s", field), current)

	switch field {
	case "headline":
		cv.Headline = value
	case "first_name":
		cv.Personal.FirstName = value
	case "surname":
		cv.Personal.Surname = value
	case "phone":
		cv.Personal.Phone = value
	case "email":
		cv.Personal.Email = value
	case "address":
		cv.Personal.Address = value
	case "website":
		cv.Personal.Website = value
	case "github":
		cv.Personal.GitHub = value
	case "linkedin":
		cv.Personal.LinkedIn = value
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
	}

	if err := saveCV(inputFile, cv); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Set %s = %q\n", field, value)
}

func interactiveAdd() {
	types := []string{"Work experience", "Education", "Language", "Digital skill", "Contact info"}
	choice := promptChoice("What to add?", types)
	if choice < 0 {
		return
	}

	cv, err := loadCV(inputFile)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	switch choice {
	case 0: // work
		fmt.Println("\n--- Add Work Experience ---")
		title := promptRequired("Job title")
		employer := prompt("Employer")
		from := promptRequired("Start date (e.g. JAN 2024)")
		to := prompt("End date (leave empty for current)")
		desc := prompt("Description")
		location := prompt("City")
		country := prompt("Country")
		tagsStr := prompt("Tags (comma-separated, for tailoring)")

		var tags []string
		if tagsStr != "" {
			tags = strings.Split(tagsStr, ",")
			for i := range tags {
				tags[i] = strings.TrimSpace(tags[i])
			}
		}

		cv.Experience = append(cv.Experience, Work{
			From: from, To: to, Title: title, Employer: employer,
			Location: location, Country: country, Description: desc, Tags: tags,
		})
		if err := saveCV(inputFile, cv); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Printf("Added work: %s @ %s (%s)\n", title, employer, from)

	case 1: // education
		fmt.Println("\n--- Add Education ---")
		title := promptRequired("Degree/certificate title")
		institution := prompt("Institution")
		from := promptRequired("Start date")
		to := prompt("End date")
		level := prompt("Level (e.g. Bachelor, EQF 6)")
		desc := prompt("Description")
		location := prompt("City")
		country := prompt("Country")

		cv.Education = append(cv.Education, Education{
			From: from, To: to, Title: title, Institution: institution,
			Location: location, Country: country, Level: level, Description: desc,
		})
		if err := saveCV(inputFile, cv); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Printf("Added education: %s @ %s\n", title, institution)

	case 2: // language
		fmt.Println("\n--- Add Foreign Language ---")
		name := promptRequired("Language name")
		fmt.Println("Enter CEFR levels (A1-C2). Leave empty to skip.")
		allLevel := prompt("Set all levels at once (or leave empty to set individually)")

		var listening, reading, spProd, spInt, writing string
		if allLevel != "" {
			upper := strings.ToUpper(allLevel)
			listening, reading, spProd, spInt, writing = upper, upper, upper, upper, upper
		} else {
			listening = promptCEFR("Listening", "")
			reading = promptCEFR("Reading", "")
			spProd = promptCEFR("Spoken production", "")
			spInt = promptCEFR("Spoken interaction", "")
			writing = promptCEFR("Writing", "")
		}

		cv.Languages.Foreign = append(cv.Languages.Foreign, ForeignLang{
			Name: name, Listening: listening, Reading: reading,
			SpokenProduction: spProd, SpokenInteraction: spInt, Writing: writing,
		})
		if err := saveCV(inputFile, cv); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Printf("Added language: %s\n", name)

	case 3: // skill
		fmt.Println("\n--- Add Digital Skills ---")
		fmt.Println("Enter skill names one per line. Empty line to finish.")
		var skills []string
		for {
			s := prompt("Skill")
			if s == "" {
				break
			}
			skills = append(skills, s)
		}
		if len(skills) > 0 {
			cv.Digital = append(cv.Digital, skills...)
			if err := saveCV(inputFile, cv); err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			fmt.Printf("Added digital skills: %s\n", strings.Join(skills, ", "))
		}

	case 4: // contact
		fmt.Println("\n--- Add Contact Info ---")
		key := promptRequired("Contact type (e.g. Matrix, Telegram, Mastodon)")
		value := promptRequired("Value")
		cv.Personal.Extra = append(cv.Personal.Extra, KV{Key: key, Value: value})
		if err := saveCV(inputFile, cv); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Printf("Added contact: %s = %s\n", key, value)
	}
}

func interactiveUpdate() {
	types := []string{"Work experience", "Education", "Language"}
	choice := promptChoice("What to update?", types)
	if choice < 0 {
		return
	}

	cv, err := loadCV(inputFile)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	switch choice {
	case 0: // work
		if len(cv.Experience) == 0 {
			fmt.Println("No work experience entries.")
			return
		}
		fmt.Println("\n--- Update Work Experience ---")
		var labels []string
		for i, w := range cv.Experience {
			labels = append(labels, fmt.Sprintf("[%d] %s @ %s (%s)", i, w.Title, w.Employer, w.From))
		}
		idx := promptChoice("Which entry?", labels)
		if idx < 0 {
			return
		}
		w := &cv.Experience[idx]
		fmt.Println("Press Enter to keep current value.")
		w.Title = promptDefault("Title", w.Title)
		w.Employer = promptDefault("Employer", w.Employer)
		w.From = promptDefault("From", w.From)
		w.To = promptDefault("To", w.To)
		w.Description = promptDefault("Description", w.Description)
		w.Location = promptDefault("Location", w.Location)
		w.Country = promptDefault("Country", w.Country)
		tagsStr := promptDefault("Tags (comma-separated)", strings.Join(w.Tags, ","))
		if tagsStr != "" {
			w.Tags = strings.Split(tagsStr, ",")
			for i := range w.Tags {
				w.Tags[i] = strings.TrimSpace(w.Tags[i])
			}
		} else {
			w.Tags = nil
		}

		if err := saveCV(inputFile, cv); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Printf("Updated work[%d]: %s @ %s\n", idx, w.Title, w.Employer)

	case 1: // education
		if len(cv.Education) == 0 {
			fmt.Println("No education entries.")
			return
		}
		fmt.Println("\n--- Update Education ---")
		var labels []string
		for i, e := range cv.Education {
			labels = append(labels, fmt.Sprintf("[%d] %s @ %s (%s)", i, e.Title, e.Institution, e.From))
		}
		idx := promptChoice("Which entry?", labels)
		if idx < 0 {
			return
		}
		e := &cv.Education[idx]
		fmt.Println("Press Enter to keep current value.")
		e.Title = promptDefault("Title", e.Title)
		e.Institution = promptDefault("Institution", e.Institution)
		e.From = promptDefault("From", e.From)
		e.To = promptDefault("To", e.To)
		e.Level = promptDefault("Level", e.Level)
		e.Description = promptDefault("Description", e.Description)
		e.Location = promptDefault("Location", e.Location)
		e.Country = promptDefault("Country", e.Country)

		if err := saveCV(inputFile, cv); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Printf("Updated education[%d]: %s @ %s\n", idx, e.Title, e.Institution)

	case 2: // language
		if len(cv.Languages.Foreign) == 0 {
			fmt.Println("No foreign language entries.")
			return
		}
		fmt.Println("\n--- Update Language ---")
		var labels []string
		for _, l := range cv.Languages.Foreign {
			labels = append(labels, l.Name)
		}
		idx := promptChoice("Which language?", labels)
		if idx < 0 {
			return
		}
		fl := &cv.Languages.Foreign[idx]
		fmt.Println("Press Enter to keep current value.")
		fl.Listening = promptCEFR("Listening", fl.Listening)
		fl.Reading = promptCEFR("Reading", fl.Reading)
		fl.SpokenProduction = promptCEFR("Spoken production", fl.SpokenProduction)
		fl.SpokenInteraction = promptCEFR("Spoken interaction", fl.SpokenInteraction)
		fl.Writing = promptCEFR("Writing", fl.Writing)

		if err := saveCV(inputFile, cv); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Printf("Updated language: %s\n", fl.Name)
	}
}

func interactiveRemove() {
	types := []string{"Work experience", "Education", "Language", "Digital skill", "Contact"}
	choice := promptChoice("What to remove?", types)
	if choice < 0 {
		return
	}

	cv, err := loadCV(inputFile)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	switch choice {
	case 0: // work
		if len(cv.Experience) == 0 {
			fmt.Println("No work experience entries.")
			return
		}
		var labels []string
		for i, w := range cv.Experience {
			labels = append(labels, fmt.Sprintf("[%d] %s @ %s (%s)", i, w.Title, w.Employer, w.From))
		}
		idx := promptChoice("Which entry to remove?", labels)
		if idx < 0 {
			return
		}
		removed := cv.Experience[idx]
		if !promptYesNo(fmt.Sprintf("Remove %q @ %s?", removed.Title, removed.Employer), false) {
			return
		}
		cv.Experience = append(cv.Experience[:idx], cv.Experience[idx+1:]...)
		if err := saveCV(inputFile, cv); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Printf("Removed work: %s @ %s\n", removed.Title, removed.Employer)

	case 1: // education
		if len(cv.Education) == 0 {
			fmt.Println("No education entries.")
			return
		}
		var labels []string
		for i, e := range cv.Education {
			labels = append(labels, fmt.Sprintf("[%d] %s @ %s", i, e.Title, e.Institution))
		}
		idx := promptChoice("Which entry to remove?", labels)
		if idx < 0 {
			return
		}
		removed := cv.Education[idx]
		if !promptYesNo(fmt.Sprintf("Remove %q @ %s?", removed.Title, removed.Institution), false) {
			return
		}
		cv.Education = append(cv.Education[:idx], cv.Education[idx+1:]...)
		if err := saveCV(inputFile, cv); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Printf("Removed education: %s @ %s\n", removed.Title, removed.Institution)

	case 2: // language
		if len(cv.Languages.Foreign) == 0 {
			fmt.Println("No foreign languages.")
			return
		}
		var labels []string
		for _, l := range cv.Languages.Foreign {
			labels = append(labels, l.Name)
		}
		idx := promptChoice("Which language to remove?", labels)
		if idx < 0 {
			return
		}
		name := cv.Languages.Foreign[idx].Name
		if !promptYesNo(fmt.Sprintf("Remove %s?", name), false) {
			return
		}
		cv.Languages.Foreign = append(cv.Languages.Foreign[:idx], cv.Languages.Foreign[idx+1:]...)
		if err := saveCV(inputFile, cv); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Printf("Removed language: %s\n", name)

	case 3: // skill
		if len(cv.Digital) == 0 {
			fmt.Println("No digital skills.")
			return
		}
		idx := promptChoice("Which skill to remove?", cv.Digital)
		if idx < 0 {
			return
		}
		name := cv.Digital[idx]
		if !promptYesNo(fmt.Sprintf("Remove %q?", name), false) {
			return
		}
		cv.Digital = append(cv.Digital[:idx], cv.Digital[idx+1:]...)
		if err := saveCV(inputFile, cv); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Printf("Removed skill: %s\n", name)

	case 4: // contact
		if len(cv.Personal.Extra) == 0 {
			fmt.Println("No extra contacts.")
			return
		}
		var labels []string
		for _, kv := range cv.Personal.Extra {
			labels = append(labels, fmt.Sprintf("%s: %s", kv.Key, kv.Value))
		}
		idx := promptChoice("Which contact to remove?", labels)
		if idx < 0 {
			return
		}
		removed := cv.Personal.Extra[idx]
		if !promptYesNo(fmt.Sprintf("Remove %s?", removed.Key), false) {
			return
		}
		cv.Personal.Extra = append(cv.Personal.Extra[:idx], cv.Personal.Extra[idx+1:]...)
		if err := saveCV(inputFile, cv); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Printf("Removed contact: %s\n", removed.Key)
	}
}

func interactiveTailor() {
	fmt.Println("\n--- Tailor CV for a Job Application ---")

	cv, err := loadCV(inputFile)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Show available tags
	tagSet := make(map[string]bool)
	for _, w := range cv.Experience {
		for _, t := range w.Tags {
			tagSet[t] = true
		}
	}
	if len(tagSet) > 0 {
		var allTags []string
		for t := range tagSet {
			allTags = append(allTags, t)
		}
		fmt.Printf("Available tags: %s\n", strings.Join(allTags, ", "))
	}

	output := promptRequired("Output JSON filename (e.g. dev-cv.json)")
	tagsStr := prompt("Include tags (comma-separated, or leave empty for all)")
	excludeStr := prompt("Exclude tags (comma-separated, or leave empty)")
	headline := promptDefault("Headline override", cv.Headline)

	// Filter
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

	cv.Headline = headline

	if err := saveCV(output, cv); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Tailored CV written to %s (%d experience entries)\n", output, len(cv.Experience))

	if promptYesNo("Also generate a PDF?", false) {
		pdfOut := promptDefault("PDF filename", strings.TrimSuffix(output, ".json")+".pdf")
		if err := generatePDF(cv, pdfOut); err != nil {
			fmt.Printf("PDF error: %v\n", err)
			return
		}
		fmt.Printf("PDF generated: %s\n", pdfOut)
	}
}

func interactiveGenerate() {
	fmt.Println("\n--- Generate PDF ---")
	inFile := promptDefault("Input JSON file", inputFile)
	outFile := promptDefault("Output PDF file", strings.TrimSuffix(inFile, ".json")+".pdf")

	cv, err := loadCV(inFile)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if err := generatePDF(cv, outFile); err != nil {
		fmt.Printf("PDF error: %v\n", err)
		return
	}
	fmt.Printf("Generated %s\n", outFile)
}
