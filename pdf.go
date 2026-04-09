package main

import (
	"fmt"
	"strings"

	"github.com/go-pdf/fpdf"
)

const (
	marginLeft   = 15.0
	marginRight  = 15.0
	marginTop    = 15.0
	pageWidth    = 210.0
	contentWidth = pageWidth - marginLeft - marginRight

	// Europass-style colors
	headerR = 30
	headerG = 60
	headerB = 114

	accentR = 0
	accentG = 102
	accentB = 153

	grayR = 100
	grayG = 100
	grayB = 100
)

const (
	fontDir    = "/usr/share/fonts/TTF"
	fontFamily = "DejaVuSans"
)

func generatePDF(cv *CV, output string) error {
	l := getLabels(cv.Lang)

	pdf := fpdf.New("P", "mm", "A4", fontDir)
	pdf.SetMargins(marginLeft, marginTop, marginRight)
	pdf.SetAutoPageBreak(true, 20)

	// Register UTF-8 fonts
	pdf.AddUTF8Font(fontFamily, "", "DejaVuSansCondensed.ttf")
	pdf.AddUTF8Font(fontFamily, "B", "DejaVuSansCondensed-Bold.ttf")
	pdf.AddUTF8Font(fontFamily, "I", "DejaVuSansCondensed-Oblique.ttf")
	pdf.AddUTF8Font(fontFamily, "BI", "DejaVuSansCondensed-BoldOblique.ttf")

	pdf.AddPage()

	renderPersonal(pdf, &cv.Personal, l)
	renderSection(pdf, l.WorkExperience, func() {
		for _, w := range cv.Experience {
			renderWork(pdf, &w, l)
		}
	})
	renderSection(pdf, l.Education, func() {
		for _, e := range cv.Education {
			renderEducation(pdf, &e, l)
		}
	})
	renderSection(pdf, l.LanguageSkills, func() {
		renderLanguages(pdf, &cv.Languages, l)
	})
	renderSection(pdf, l.DigitalSkills, func() {
		pdf.SetFont(fontFamily, "", 10)
		pdf.SetTextColor(0, 0, 0)
		pdf.MultiCell(contentWidth, 5, strings.Join(cv.Digital, "  |  "), "", "L", false)
		pdf.Ln(3)
	})
	if cv.Org != "" || cv.Comm != "" || cv.JobRelated != "" {
		renderSection(pdf, l.AdditionalInfo, func() {
			if cv.Org != "" {
				renderSubSection(pdf, l.OrgSkills, cv.Org)
			}
			if cv.Comm != "" {
				renderSubSection(pdf, l.CommSkills, cv.Comm)
			}
			if cv.JobRelated != "" {
				renderSubSection(pdf, l.JobSkills, cv.JobRelated)
			}
		})
	}

	// Embed Europass XML as attachment
	xmlData := toEuropassXML(cv)
	pdf.SetAttachments([]fpdf.Attachment{
		{
			Content:     xmlData,
			Filename:    "Europass_CV.xml",
			Description: "Europass CV XML data",
		},
	})

	return pdf.OutputFileAndClose(output)
}

func generatePlainPDF(cv *CV, output string) error {
	l := getLabels(cv.Lang)

	pdf := fpdf.New("P", "mm", "A4", fontDir)
	pdf.SetMargins(marginLeft, marginTop, marginRight)
	pdf.SetAutoPageBreak(true, 20)

	pdf.AddUTF8Font(fontFamily, "", "DejaVuSansCondensed.ttf")
	pdf.AddUTF8Font(fontFamily, "B", "DejaVuSansCondensed-Bold.ttf")
	pdf.AddUTF8Font(fontFamily, "I", "DejaVuSansCondensed-Oblique.ttf")
	pdf.AddUTF8Font(fontFamily, "BI", "DejaVuSansCondensed-BoldOblique.ttf")

	pdf.AddPage()

	renderPersonal(pdf, &cv.Personal, l)
	renderSection(pdf, l.WorkExperience, func() {
		for _, w := range cv.Experience {
			renderWork(pdf, &w, l)
		}
	})
	renderSection(pdf, l.Education, func() {
		for _, e := range cv.Education {
			renderEducation(pdf, &e, l)
		}
	})
	renderSection(pdf, l.LanguageSkills, func() {
		renderLanguages(pdf, &cv.Languages, l)
	})
	renderSection(pdf, l.DigitalSkills, func() {
		pdf.SetFont(fontFamily, "", 10)
		pdf.SetTextColor(0, 0, 0)
		pdf.MultiCell(contentWidth, 5, strings.Join(cv.Digital, "  |  "), "", "L", false)
		pdf.Ln(3)
	})
	if cv.Org != "" || cv.Comm != "" || cv.JobRelated != "" {
		renderSection(pdf, l.AdditionalInfo, func() {
			if cv.Org != "" {
				renderSubSection(pdf, l.OrgSkills, cv.Org)
			}
			if cv.Comm != "" {
				renderSubSection(pdf, l.CommSkills, cv.Comm)
			}
			if cv.JobRelated != "" {
				renderSubSection(pdf, l.JobSkills, cv.JobRelated)
			}
		})
	}

	return pdf.OutputFileAndClose(output)
}

func renderPersonal(pdf *fpdf.Fpdf, p *Personal, l Labels) {
	// Name
	pdf.SetFont(fontFamily, "B", 20)
	pdf.SetTextColor(headerR, headerG, headerB)
	pdf.CellFormat(contentWidth, 10, p.FirstName+" "+p.Surname, "", 1, "L", false, 0, "")
	pdf.Ln(1)

	// Horizontal rule
	pdf.SetDrawColor(headerR, headerG, headerB)
	pdf.SetLineWidth(0.8)
	pdf.Line(marginLeft, pdf.GetY(), pageWidth-marginRight, pdf.GetY())
	pdf.Ln(3)

	// Contact details
	pdf.SetFont(fontFamily, "", 9)
	pdf.SetTextColor(0, 0, 0)

	var details []string
	if p.DateOfBirth != "" {
		details = append(details, l.DateOfBirth+": "+p.DateOfBirth)
	}
	if p.Nationality != "" {
		details = append(details, l.Nationality+": "+p.Nationality)
	}
	if p.Phone != "" {
		details = append(details, l.Phone+": "+p.Phone)
	}
	if p.Email != "" {
		details = append(details, l.Email+": "+p.Email)
	}
	if len(details) > 0 {
		pdf.MultiCell(contentWidth, 4.5, strings.Join(details, "  |  "), "", "L", false)
	}

	var links []string
	if p.Website != "" {
		links = append(links, l.Website+": "+p.Website)
	}
	if p.GitHub != "" {
		links = append(links, "GitHub: "+p.GitHub)
	}
	if p.LinkedIn != "" {
		links = append(links, "LinkedIn: "+p.LinkedIn)
	}
	for _, kv := range p.Extra {
		links = append(links, kv.Key+": "+kv.Value)
	}
	if len(links) > 0 {
		pdf.MultiCell(contentWidth, 4.5, strings.Join(links, "  |  "), "", "L", false)
	}

	if p.Address != "" {
		pdf.CellFormat(contentWidth, 4.5, l.Address+": "+p.Address, "", 1, "L", false, 0, "")
	}
	pdf.Ln(5)
}

func renderSection(pdf *fpdf.Fpdf, title string, content func()) {
	// Bullet + title
	pdf.SetFont(fontFamily, "B", 13)
	pdf.SetTextColor(headerR, headerG, headerB)

	// Small filled circle as bullet
	y := pdf.GetY() + 4
	pdf.SetFillColor(headerR, headerG, headerB)
	pdf.Circle(marginLeft+2, y, 2.5, "F")

	pdf.SetX(marginLeft + 8)
	pdf.CellFormat(contentWidth-8, 10, title, "", 1, "L", false, 0, "")

	// Rule under section header
	pdf.SetDrawColor(grayR, grayG, grayB)
	pdf.SetLineWidth(0.3)
	pdf.Line(marginLeft+8, pdf.GetY(), pageWidth-marginRight, pdf.GetY())
	pdf.Ln(3)

	content()
	pdf.Ln(2)
}

func renderWork(pdf *fpdf.Fpdf, w *Work, l Labels) {
	// Date range + location
	pdf.SetFont(fontFamily, "", 9)
	pdf.SetTextColor(grayR, grayG, grayB)
	period := w.From
	if w.To != "" {
		period += " \u2013 " + w.To
	} else {
		period += " \u2013 " + l.Present
	}
	loc := ""
	if w.Location != "" {
		loc = w.Location
		if w.Country != "" {
			loc += ", " + w.Country
		}
	}
	if loc != "" {
		period += "  " + loc
	}
	pdf.CellFormat(contentWidth, 4.5, period, "", 1, "L", false, 0, "")

	// Title on one line, employer on same line right-aligned won't fit — put employer after title
	pdf.SetFont(fontFamily, "B", 10)
	pdf.SetTextColor(0, 0, 0)
	titleStr := strings.ToUpper(w.Title)
	pdf.CellFormat(0, 5.5, titleStr+" "+w.Employer, "", 1, "L", false, 0, "")
	pdf.Ln(1)

	// Description
	if w.Description != "" {
		pdf.SetFont(fontFamily, "", 9)
		pdf.SetTextColor(50, 50, 50)
		pdf.MultiCell(contentWidth, 4.5, w.Description, "", "L", false)
	}
	pdf.Ln(3)
}

func renderEducation(pdf *fpdf.Fpdf, e *Education, l Labels) {
	pdf.SetFont(fontFamily, "", 9)
	pdf.SetTextColor(grayR, grayG, grayB)
	period := e.From
	if e.To != "" {
		period += " \u2013 " + e.To
	}
	loc := ""
	if e.Location != "" {
		loc = e.Location
		if e.Country != "" {
			loc += ", " + e.Country
		}
	}
	if loc != "" {
		period += "  " + loc
	}
	pdf.CellFormat(contentWidth, 4.5, period, "", 1, "L", false, 0, "")

	pdf.SetFont(fontFamily, "B", 10)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(0, 5.5, strings.ToUpper(e.Title)+" "+e.Institution, "", 1, "L", false, 0, "")

	if e.Level != "" {
		pdf.SetFont(fontFamily, "", 9)
		pdf.SetTextColor(grayR, grayG, grayB)
		pdf.CellFormat(contentWidth, 4.5, l.Level+": "+e.Level, "", 1, "L", false, 0, "")
	}

	if e.Description != "" {
		pdf.SetFont(fontFamily, "", 9)
		pdf.SetTextColor(50, 50, 50)
		pdf.MultiCell(contentWidth, 4.5, e.Description, "", "L", false)
	}
	pdf.Ln(3)
}

func renderLanguages(pdf *fpdf.Fpdf, lang *Languages, l Labels) {
	// Mother tongue
	pdf.SetFont(fontFamily, "", 10)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(contentWidth, 5.5, l.MotherTongue+": "+strings.Join(lang.MotherTongue, ", "), "", 1, "L", false, 0, "")
	pdf.Ln(2)

	if len(lang.Foreign) == 0 {
		return
	}

	pdf.SetFont(fontFamily, "", 9)
	pdf.SetTextColor(grayR, grayG, grayB)
	pdf.CellFormat(contentWidth, 4.5, l.OtherLanguages+":", "", 1, "L", false, 0, "")
	pdf.Ln(2)

	// Table header
	colW := contentWidth / 6
	nameW := colW

	pdf.SetFont(fontFamily, "B", 8)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFillColor(230, 230, 230)
	pdf.CellFormat(nameW, 6, "", "1", 0, "C", true, 0, "")
	pdf.CellFormat(colW, 6, l.Listening, "1", 0, "C", true, 0, "")
	pdf.CellFormat(colW, 6, l.Reading, "1", 0, "C", true, 0, "")
	pdf.CellFormat(colW, 6, l.SpokenProd, "1", 0, "C", true, 0, "")
	pdf.CellFormat(colW, 6, l.SpokenInt, "1", 0, "C", true, 0, "")
	pdf.CellFormat(colW, 6, l.Writing, "1", 1, "C", true, 0, "")

	pdf.SetFont(fontFamily, "", 8)
	for _, fl := range lang.Foreign {
		pdf.SetFont(fontFamily, "B", 8)
		pdf.CellFormat(nameW, 5.5, fl.Name, "1", 0, "L", false, 0, "")
		pdf.SetFont(fontFamily, "", 8)
		pdf.CellFormat(colW, 5.5, fl.Listening, "1", 0, "C", false, 0, "")
		pdf.CellFormat(colW, 5.5, fl.Reading, "1", 0, "C", false, 0, "")
		pdf.CellFormat(colW, 5.5, fl.SpokenProduction, "1", 0, "C", false, 0, "")
		pdf.CellFormat(colW, 5.5, fl.SpokenInteraction, "1", 0, "C", false, 0, "")
		pdf.CellFormat(colW, 5.5, fl.Writing, "1", 1, "C", false, 0, "")
	}
	pdf.Ln(2)

	pdf.SetFont(fontFamily, "I", 7)
	pdf.SetTextColor(grayR, grayG, grayB)
	pdf.MultiCell(contentWidth, 3.5, l.CEFRLegend, "", "L", false)
	pdf.Ln(3)
}

func renderSubSection(pdf *fpdf.Fpdf, title, text string) {
	pdf.SetFont(fontFamily, "B", 10)
	pdf.SetTextColor(headerR, headerG, headerB)
	pdf.CellFormat(contentWidth, 6, title, "", 1, "L", false, 0, "")

	pdf.SetFont(fontFamily, "", 9)
	pdf.SetTextColor(50, 50, 50)
	pdf.MultiCell(contentWidth, 4.5, text, "", "L", false)
	pdf.Ln(3)
}

func toEuropassXML(cv *CV) []byte {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	b.WriteString(`<SkillsPassport xmlns="http://europass.cedefop.europa.eu/Europass" `)
	b.WriteString(`xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">` + "\n")
	b.WriteString("  <LearnerInfo>\n")

	// Identification
	b.WriteString("    <Identification>\n")
	b.WriteString("      <PersonName>\n")
	b.WriteString(fmt.Sprintf("        <FirstName>%s</FirstName>\n", xmlEsc(cv.Personal.FirstName)))
	b.WriteString(fmt.Sprintf("        <Surname>%s</Surname>\n", xmlEsc(cv.Personal.Surname)))
	b.WriteString("      </PersonName>\n")
	b.WriteString("      <ContactInfo>\n")
	if cv.Personal.Email != "" {
		b.WriteString(fmt.Sprintf("        <Email><Contact>%s</Contact></Email>\n", xmlEsc(cv.Personal.Email)))
	}
	if cv.Personal.Phone != "" {
		b.WriteString(fmt.Sprintf("        <Telephone><Contact>%s</Contact></Telephone>\n", xmlEsc(cv.Personal.Phone)))
	}
	if cv.Personal.Address != "" {
		b.WriteString(fmt.Sprintf("        <Address><Contact><AddressLine>%s</AddressLine></Contact></Address>\n", xmlEsc(cv.Personal.Address)))
	}
	if cv.Personal.Website != "" {
		b.WriteString(fmt.Sprintf("        <Website><Contact>%s</Contact></Website>\n", xmlEsc(cv.Personal.Website)))
	}
	b.WriteString("      </ContactInfo>\n")
	b.WriteString("      <Demographics>\n")
	if cv.Personal.DateOfBirth != "" {
		b.WriteString(fmt.Sprintf("        <Birthdate>%s</Birthdate>\n", xmlEsc(cv.Personal.DateOfBirth)))
	}
	if cv.Personal.Nationality != "" {
		b.WriteString(fmt.Sprintf("        <Nationality><Label>%s</Label></Nationality>\n", xmlEsc(cv.Personal.Nationality)))
	}
	b.WriteString("      </Demographics>\n")
	b.WriteString("    </Identification>\n")

	// Headline
	if cv.Headline != "" {
		b.WriteString(fmt.Sprintf("    <Headline><Type><Label>%s</Label></Type></Headline>\n", xmlEsc(cv.Headline)))
	}

	// Work experience
	for _, w := range cv.Experience {
		b.WriteString("    <WorkExperience>\n")
		b.WriteString("      <Period>\n")
		b.WriteString(fmt.Sprintf("        <From>%s</From>\n", xmlEsc(w.From)))
		if w.To != "" {
			b.WriteString(fmt.Sprintf("        <To>%s</To>\n", xmlEsc(w.To)))
		}
		b.WriteString("      </Period>\n")
		b.WriteString(fmt.Sprintf("      <Position><Label>%s</Label></Position>\n", xmlEsc(w.Title)))
		b.WriteString(fmt.Sprintf("      <Employer><Name>%s</Name></Employer>\n", xmlEsc(w.Employer)))
		if w.Location != "" || w.Country != "" {
			loc := w.Location
			if w.Country != "" {
				if loc != "" {
					loc += ", "
				}
				loc += w.Country
			}
			b.WriteString(fmt.Sprintf("      <Employer><Address>%s</Address></Employer>\n", xmlEsc(loc)))
		}
		if w.Description != "" {
			b.WriteString(fmt.Sprintf("      <Activities>%s</Activities>\n", xmlEsc(w.Description)))
		}
		b.WriteString("    </WorkExperience>\n")
	}

	// Education
	for _, e := range cv.Education {
		b.WriteString("    <Education>\n")
		b.WriteString("      <Period>\n")
		b.WriteString(fmt.Sprintf("        <From>%s</From>\n", xmlEsc(e.From)))
		if e.To != "" {
			b.WriteString(fmt.Sprintf("        <To>%s</To>\n", xmlEsc(e.To)))
		}
		b.WriteString("      </Period>\n")
		b.WriteString(fmt.Sprintf("      <Title><Label>%s</Label></Title>\n", xmlEsc(e.Title)))
		b.WriteString(fmt.Sprintf("      <Organisation><Name>%s</Name></Organisation>\n", xmlEsc(e.Institution)))
		if e.Level != "" {
			b.WriteString(fmt.Sprintf("      <Level><Label>%s</Label></Level>\n", xmlEsc(e.Level)))
		}
		if e.Description != "" {
			b.WriteString(fmt.Sprintf("      <Activities>%s</Activities>\n", xmlEsc(e.Description)))
		}
		b.WriteString("    </Education>\n")
	}

	// Skills
	b.WriteString("    <Skills>\n")

	// Languages
	b.WriteString("      <Linguistic>\n")
	for _, mt := range cv.Languages.MotherTongue {
		b.WriteString(fmt.Sprintf("        <MotherTongue><Description><Label>%s</Label></Description></MotherTongue>\n", xmlEsc(mt)))
	}
	for _, fl := range cv.Languages.Foreign {
		b.WriteString("        <ForeignLanguage>\n")
		b.WriteString(fmt.Sprintf("          <Description><Label>%s</Label></Description>\n", xmlEsc(fl.Name)))
		b.WriteString(fmt.Sprintf("          <Listening>%s</Listening>\n", fl.Listening))
		b.WriteString(fmt.Sprintf("          <Reading>%s</Reading>\n", fl.Reading))
		b.WriteString(fmt.Sprintf("          <SpokenProduction>%s</SpokenProduction>\n", fl.SpokenProduction))
		b.WriteString(fmt.Sprintf("          <SpokenInteraction>%s</SpokenInteraction>\n", fl.SpokenInteraction))
		b.WriteString(fmt.Sprintf("          <Writing>%s</Writing>\n", fl.Writing))
		b.WriteString("        </ForeignLanguage>\n")
	}
	b.WriteString("      </Linguistic>\n")

	// Digital skills
	if len(cv.Digital) > 0 {
		b.WriteString(fmt.Sprintf("      <Computer>%s</Computer>\n", xmlEsc(strings.Join(cv.Digital, ", "))))
	}

	// Soft skills
	if cv.Org != "" {
		b.WriteString(fmt.Sprintf("      <Organisational>%s</Organisational>\n", xmlEsc(cv.Org)))
	}
	if cv.Comm != "" {
		b.WriteString(fmt.Sprintf("      <Communication>%s</Communication>\n", xmlEsc(cv.Comm)))
	}
	if cv.JobRelated != "" {
		b.WriteString(fmt.Sprintf("      <JobRelated>%s</JobRelated>\n", xmlEsc(cv.JobRelated)))
	}

	b.WriteString("    </Skills>\n")
	b.WriteString("  </LearnerInfo>\n")
	b.WriteString("</SkillsPassport>\n")
	return []byte(b.String())
}

func xmlEsc(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return s
}
