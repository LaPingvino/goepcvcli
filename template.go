package main

// templateCV returns Joop's current CV data as the default template.
// This serves as both a real CV and a structural example.
func templateCV() *CV {
	return &CV{
		Personal: Personal{
			FirstName:   "Abel Johannes",
			Surname:     "Kiefte",
			DateOfBirth: "25 Jan 1989",
			Nationality: "Dutch",
			Phone:       "(+351) 913044570",
			Email:       "joop@kiefte.net",
			Website:     "https://joop.kiefte.eu",
			GitHub:      "https://github.com/lapingvino",
			LinkedIn:    "https://linkedin.com/in/kiefte",
			Address:     "Portugal",
			Extra: []KV{
				{Key: "Google Chat", Value: "ikojba@gmail.com"},
				{Key: "Telegram", Value: "https://t.me/lapingvino"},
			},
		},
		Headline: "Systems thinker, natural high level troubleshooter, hyperpolyglot",
		Experience: []Work{
			{
				From:     "FEB 2023",
				Title:    "Second Line Support",
				Employer: "Teleperformance (Wolters Kluwer)",
				Location: "Lisbon",
				Country:  "Portugal",
				Description: "Returning to bring in the missing technical know-how on the project. " +
					"B2B technical support for legal practice management software (CRM, billing, document management, case management). " +
					"Troubleshooting complex integrations and escalated issues.",
				Tags: []string{"support", "legal-tech", "b2b"},
			},
			{
				From:     "JUL 2015",
				To:       "2018",
				Title:    "Computer Programmer",
				Employer: "Claire Automotive Support B.V.",
				Location: "Ede",
				Country:  "Netherlands",
				Description: "Planning the overhaul of the legacy system and getting to an understanding of the true needs for the new system. " +
					"Go programming. Working with Google products like GCP. " +
					"Tech support. Finding new technically adept recruits where previous attempts were fruitless.",
				Tags: []string{"dev", "go", "gcp", "architecture"},
			},
			{
				From:     "SEP 2018",
				To:       "NOV 2019",
				Title:    "Web Content Manager",
				Employer: "Wageningen University & Research",
				Location: "Maastricht",
				Country:  "Netherlands",
				Description: "Editing online contents for the different sections. " +
					"Streamlining and correcting processes, or creating new ones where they were lacking. " +
					"Assisting the section responsible to get a grip on the infrastructure.",
				Tags: []string{"content", "web", "process"},
			},
			{
				From:     "SEP 2019",
				To:       "SEP 2020",
				Title:    "President",
				Employer: "World Esperanto Youth Organization",
				Description: "Non-paid function with international NGO. " +
					"Solving issues when they arise, allowing the organization to be bold.",
				Tags: []string{"leadership", "ngo", "international"},
			},
			{
				From:     "MAY 2021",
				To:       "NOV 2022",
				Title:    "Technical Support Agent",
				Employer: "Teleperformance (Wolters Kluwer)",
				Location: "Lisbon",
				Country:  "Portugal",
				Description: "Technical support by phone and mail, business to business, for legal software.",
				Tags: []string{"support", "legal-tech", "b2b"},
			},
			{
				From:     "NOV 2022",
				To:       "JAN 2023",
				Title:    "Inhouse Tech Support",
				Employer: "Delcom International B.V.",
				Location: "Remote (Amsterdam, NL / Estoril, PT)",
				Country:  "Portugal",
				Description: "Tech Support for a software product for sport schools. " +
					"Includes writing queries for reports and some DevOps tasks.",
				Tags: []string{"support", "devops", "sql"},
			},
			{
				From:     "JAN 2012",
				To:       "DEC 2013",
				Title:    "Technical Support Agent",
				Employer: "Teleperformance (Microsoft)",
				Location: "Maastricht",
				Country:  "Netherlands",
				Description: "Tech support phone and mail. " +
					"Windows, Office, virus removal, providing support for blind users and people with linguistically complex systems.",
				Tags: []string{"support", "microsoft"},
			},
			{
				From:     "MAY 2012",
				To:       "OCT 2012",
				Title:    "Software Engineer",
				Employer: "LaunchIT",
				Location: "Heerlen",
				Country:  "Netherlands",
				Description: "Hunting bugs and improving the development procedures so software quality goes up.",
				Tags: []string{"dev", "qa"},
			},
			{
				From:     "DEC 2009",
				To:       "FEB 2010",
				Title:    "Programmer",
				Employer: "Bronboek Software Development BV",
				Description: "Software development.",
				Tags: []string{"dev"},
			},
		},
		Education: []Education{
			{
				From:        "2006",
				To:          "2009",
				Title:       "ICT Management Assistant",
				Institution: "RijnIJssel",
				Location:    "Arnhem",
				Country:     "Netherlands",
				Level:       "MBO ICT 3",
				Description: "Windows, Linux, networking competencies. Microsoft and Cisco certification.",
			},
		},
		Languages: Languages{
			MotherTongue: []string{"Dutch"},
			Foreign: []ForeignLang{
				{Name: "Portuguese", Listening: "C1", Reading: "C2", SpokenProduction: "C1", SpokenInteraction: "B2", Writing: "C1"},
				{Name: "English", Listening: "C2", Reading: "C2", SpokenProduction: "C1", SpokenInteraction: "C2", Writing: "C1"},
				{Name: "Esperanto", Listening: "C2", Reading: "C2", SpokenProduction: "C2", SpokenInteraction: "C2", Writing: "C2"},
				{Name: "Spanish", Listening: "B2", Reading: "B2", SpokenProduction: "B2", SpokenInteraction: "B2", Writing: "B2"},
				{Name: "French", Listening: "B1", Reading: "B1", SpokenProduction: "B2", SpokenInteraction: "B2", Writing: "B1"},
				{Name: "German", Listening: "C1", Reading: "C1", SpokenProduction: "B1", SpokenInteraction: "B2", Writing: "B1"},
				{Name: "Italian", Listening: "B2", Reading: "B2", SpokenProduction: "B1", SpokenInteraction: "B2", Writing: "A2"},
				{Name: "Afrikaans", Listening: "B1", Reading: "C1", SpokenProduction: "A2", SpokenInteraction: "A2", Writing: "A2"},
			},
		},
		Digital: []string{
			"Go", "TypeScript", "Rust", "Python", "Linux", "Git", "Bash",
			"Google Cloud Platform", "Docker", "SQL",
			"Microsoft Office", "Google Workspace",
			"IT Troubleshooting", "Programming and Scripting", "Debugging",
		},
		Org: "Understands well what needs and doesn't need to be arranged. " +
			"Capable of understanding the goals of a subject, and focus on having results even if a full resolution turns out to be impossible. " +
			"Capable of changing the way one works.",
		Comm: "Good communication capabilities going across cultures and different ways of thinking. " +
			"Good capabilities to listen and understand the hidden needs of the customer. " +
			"Knows to answer the unsaid questions instead of focusing on the details of the words used. " +
			"Manages to learn quickly the language and communication quirks of others.",
		JobRelated: "Very strong in troubleshooting of any kind, also outside of programming. " +
			"Strong capacity in teaching technology to colleague and user alike. " +
			"Strong capacity in learning new technologies that turn out to be needed for a job. " +
			"Strong focus on solving a problem in a forward-thinking matter that minimizes technical debt.",
	}
}
