package main

// Europass CV data model — matches the Europass XML schema structure
// but stored as JSON on disk for easy editing and per-job tailoring.

type CV struct {
	Personal   Personal    `json:"personal" xml:"learnerInfo>identification"`
	Headline   string      `json:"headline" xml:"learnerInfo>headline>description"`
	Experience []Work      `json:"experience" xml:"learnerInfo>workExperience"`
	Education  []Education `json:"education" xml:"learnerInfo>education"`
	Languages  Languages   `json:"languages" xml:"learnerInfo>skills>linguistic"`
	Digital    []string    `json:"digital_skills" xml:"learnerInfo>skills>computer"`
	Org        string      `json:"organisational_skills"`
	Comm       string      `json:"communication_skills"`
	JobRelated string      `json:"job_related_skills"`
}

type Personal struct {
	FirstName   string   `json:"first_name" xml:"personName>firstName"`
	Surname     string   `json:"surname" xml:"personName>surname"`
	DateOfBirth string   `json:"date_of_birth" xml:"demographics>birthdate"`
	Nationality string   `json:"nationality" xml:"demographics>nationality>label"`
	Phone       string   `json:"phone" xml:"contactInfo>telephone>contact"`
	Email       string   `json:"email" xml:"contactInfo>email>contact"`
	Website     string   `json:"website"`
	GitHub      string   `json:"github"`
	LinkedIn    string   `json:"linkedin"`
	Address     string   `json:"address" xml:"contactInfo>address>contact>addressLine"`
	Extra       []KV     `json:"extra,omitempty"` // Telegram, Google Chat, etc.
}

type KV struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Work struct {
	From        string   `json:"from" xml:"period>from"`
	To          string   `json:"to,omitempty" xml:"period>to"`
	Title       string   `json:"title" xml:"position>label"`
	Employer    string   `json:"employer" xml:"employer>name"`
	Location    string   `json:"location"`
	Country     string   `json:"country"`
	Description string   `json:"description" xml:"activities"`
	Tags        []string `json:"tags,omitempty"` // for tailoring: filter by tag
}

type Education struct {
	From        string `json:"from" xml:"period>from"`
	To          string `json:"to,omitempty" xml:"period>to"`
	Title       string `json:"title" xml:"title>label"`
	Institution string `json:"institution" xml:"organisation>name"`
	Location    string `json:"location"`
	Country     string `json:"country"`
	Level       string `json:"level,omitempty"` // EQF / national classification
	Description string `json:"description" xml:"activities"`
}

type Languages struct {
	MotherTongue []string  `json:"mother_tongue"`
	Foreign      []ForeignLang `json:"foreign"`
}

type ForeignLang struct {
	Name              string `json:"name"`
	Listening         string `json:"listening"`
	Reading           string `json:"reading"`
	SpokenProduction  string `json:"spoken_production"`
	SpokenInteraction string `json:"spoken_interaction"`
	Writing           string `json:"writing"`
}
