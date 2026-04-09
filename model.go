package main

// Europass CV data model — matches the Europass XML schema structure
// but stored as JSON on disk for easy editing and per-job tailoring.

type CV struct {
	Lang       string      `json:"lang,omitempty"`
	Style      Style       `json:"style,omitempty"`
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

// Style controls the visual appearance of the PDF output.
type Style struct {
	Layout     string   `json:"layout,omitempty"`      // "classic" (default), "modern", "compact"
	Accent     [3]int   `json:"accent,omitempty"`      // accent color [R,G,B] (default: [30,60,114])
	FontSize   float64  `json:"font_size,omitempty"`   // base body font size (default: 9)
	DateColumn float64  `json:"date_column,omitempty"` // date column width in mm (modern layout, default: 38)
	HeaderFont FontSpec `json:"header_font,omitempty"` // font for section headers and name
	BodyFont   FontSpec `json:"body_font,omitempty"`   // font for body text
	Fallback   []string `json:"fallback,omitempty"`    // fallback TTF paths for missing glyphs (e.g. emoji)
}

// FontSpec defines a font by name (resolved via fc-match) or explicit TTF path.
// Name is a fontconfig name like "DejaVu Sans Condensed" or "Noto Serif".
// If Path is set, it overrides Name and is used directly as the TTF file path.
type FontSpec struct {
	Name string `json:"name,omitempty"` // fontconfig name, resolved via fc-match
	Path string `json:"path,omitempty"` // explicit TTF path override
}

// resolved returns a Style with all defaults filled in.
func (s Style) resolved() Style {
	if s.Layout == "" {
		s.Layout = "classic"
	}
	if s.Accent == [3]int{} {
		s.Accent = [3]int{30, 60, 114}
	}
	if s.FontSize == 0 {
		s.FontSize = 9
	}
	if s.DateColumn == 0 {
		s.DateColumn = 38
	}
	if s.HeaderFont.Name == "" && s.HeaderFont.Path == "" {
		s.HeaderFont.Name = "DejaVu Sans Condensed"
	}
	if s.BodyFont.Name == "" && s.BodyFont.Path == "" {
		s.BodyFont.Name = "DejaVu Sans Condensed"
	}
	return s
}
