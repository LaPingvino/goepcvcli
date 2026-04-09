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

	grayR = 100
	grayG = 100
	grayB = 100
)

// renderer bundles everything needed to render a CV page.
type renderer struct {
	pdf   *fpdf.Fpdf
	fonts pdfFonts
	style Style
	label Labels
}

func newRenderer(cv *CV) (*renderer, error) {
	style := cv.Style.resolved()
	l := getLabels(cv.Lang)

	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(marginLeft, marginTop, marginRight)
	pdf.SetAutoPageBreak(true, 20)

	fonts := setupFonts(pdf, style)

	return &renderer{pdf: pdf, fonts: fonts, style: style, label: l}, nil
}

func (r *renderer) accent() (int, int, int) {
	return r.style.Accent[0], r.style.Accent[1], r.style.Accent[2]
}

func (r *renderer) renderCV(cv *CV) {
	r.pdf.AddPage()
	l := r.label

	switch r.style.Layout {
	case "modern":
		r.renderPersonalModern(&cv.Personal, cv.Headline)
	default:
		r.renderPersonalClassic(&cv.Personal, cv.Headline)
	}

	r.renderSection(l.WorkExperience, func() {
		for _, w := range cv.Experience {
			r.renderWork(&w)
		}
	})
	r.renderSection(l.Education, func() {
		for _, e := range cv.Education {
			r.renderEducation(&e)
		}
	})
	r.renderSection(l.LanguageSkills, func() {
		r.renderLanguages(&cv.Languages)
	})
	r.renderSection(l.DigitalSkills, func() {
		r.setBody("", r.style.FontSize+1)
		r.pdf.SetTextColor(0, 0, 0)

		if r.style.Layout == "modern" {
			r.renderSkillChips(cv.Digital)
		} else {
			r.pdf.MultiCell(contentWidth, 5, strings.Join(cv.Digital, "  \u2022  "), "", "L", false)
		}
		r.pdf.Ln(3)
	})
	if cv.Org != "" || cv.Comm != "" || cv.JobRelated != "" {
		r.renderSection(l.AdditionalInfo, func() {
			if cv.Org != "" {
				r.renderSubSection(l.OrgSkills, cv.Org)
			}
			if cv.Comm != "" {
				r.renderSubSection(l.CommSkills, cv.Comm)
			}
			if cv.JobRelated != "" {
				r.renderSubSection(l.JobSkills, cv.JobRelated)
			}
		})
	}
}

// --- Header font helpers ---

func (r *renderer) setHeader(style string, size float64) {
	r.pdf.SetFont(r.fonts.header.family, style, size)
}

func (r *renderer) setBody(style string, size float64) {
	r.pdf.SetFont(r.fonts.body.family, style, size)
}

func (r *renderer) writeBody(h float64, text, style string) {
	r.fonts.writeWithFallback(r.pdf, h, text, r.fonts.body.family, style)
}

func (r *renderer) writeHeader(h float64, text, style string) {
	r.fonts.writeWithFallback(r.pdf, h, text, r.fonts.header.family, style)
}

// --- Classic layout ---

func (r *renderer) renderPersonalClassic(p *Personal, headline string) {
	aR, aG, aB := r.accent()
	pdf := r.pdf

	// Name
	r.setHeader("B", 22)
	pdf.SetTextColor(aR, aG, aB)
	r.writeHeader(10, p.FirstName+" "+p.Surname, "B")
	pdf.Ln(12)

	// Headline
	if headline != "" {
		r.setHeader("I", r.style.FontSize+2)
		pdf.SetTextColor(aR, aG, aB)
		r.writeHeader(6, headline, "I")
		pdf.Ln(8)
	}

	// Accent bar
	pdf.SetDrawColor(aR, aG, aB)
	pdf.SetLineWidth(0.8)
	pdf.Line(marginLeft, pdf.GetY(), pageWidth-marginRight, pdf.GetY())
	pdf.Ln(4)

	// Contact grid
	r.setBody("", r.style.FontSize)
	pdf.SetTextColor(60, 60, 60)
	l := r.label

	var row1, row2 []string
	if p.Phone != "" {
		row1 = append(row1, l.Phone+": "+p.Phone)
	}
	if p.Email != "" {
		row1 = append(row1, l.Email+": "+p.Email)
	}
	if p.DateOfBirth != "" {
		row1 = append(row1, l.DateOfBirth+": "+p.DateOfBirth)
	}
	if p.Nationality != "" {
		row1 = append(row1, l.Nationality+": "+p.Nationality)
	}
	if len(row1) > 0 {
		pdf.MultiCell(contentWidth, 4.5, strings.Join(row1, "  |  "), "", "L", false)
	}

	if p.Website != "" {
		row2 = append(row2, l.Website+": "+p.Website)
	}
	if p.GitHub != "" {
		row2 = append(row2, "GitHub: "+p.GitHub)
	}
	if p.LinkedIn != "" {
		row2 = append(row2, "LinkedIn: "+p.LinkedIn)
	}
	for _, kv := range p.Extra {
		row2 = append(row2, kv.Key+": "+kv.Value)
	}
	if len(row2) > 0 {
		pdf.MultiCell(contentWidth, 4.5, strings.Join(row2, "  |  "), "", "L", false)
	}
	if p.Address != "" {
		pdf.CellFormat(contentWidth, 4.5, l.Address+": "+p.Address, "", 1, "L", false, 0, "")
	}
	pdf.Ln(5)
}

// --- Modern layout ---

func (r *renderer) renderPersonalModern(p *Personal, headline string) {
	aR, aG, aB := r.accent()
	pdf := r.pdf

	// Colored accent band on the left
	bandW := 4.0
	startY := marginTop
	pdf.SetFillColor(aR, aG, aB)
	pdf.Rect(marginLeft, startY, bandW, 30, "F")

	// Name next to the band
	r.setHeader("B", 24)
	pdf.SetTextColor(aR, aG, aB)
	pdf.SetXY(marginLeft+bandW+4, startY+2)
	r.writeHeader(10, p.FirstName+" "+p.Surname, "B")
	pdf.Ln(10)

	// Headline
	if headline != "" {
		r.setBody("I", r.style.FontSize+1)
		pdf.SetTextColor(80, 80, 80)
		pdf.SetX(marginLeft + bandW + 4)
		r.writeBody(5, headline, "I")
		pdf.Ln(8)
	}

	pdf.SetY(startY + 34)

	// Contact in two columns
	r.setBody("", r.style.FontSize-0.5)
	pdf.SetTextColor(60, 60, 60)
	l := r.label
	halfW := contentWidth / 2

	type kv struct{ k, v string }
	var left, right []kv
	if p.Phone != "" {
		left = append(left, kv{l.Phone, p.Phone})
	}
	if p.Email != "" {
		left = append(left, kv{l.Email, p.Email})
	}
	if p.Website != "" {
		left = append(left, kv{l.Website, p.Website})
	}
	if p.Address != "" {
		left = append(left, kv{l.Address, p.Address})
	}
	if p.DateOfBirth != "" {
		right = append(right, kv{l.DateOfBirth, p.DateOfBirth})
	}
	if p.Nationality != "" {
		right = append(right, kv{l.Nationality, p.Nationality})
	}
	if p.GitHub != "" {
		right = append(right, kv{"GitHub", p.GitHub})
	}
	if p.LinkedIn != "" {
		right = append(right, kv{"LinkedIn", p.LinkedIn})
	}
	for _, e := range p.Extra {
		right = append(right, kv{e.Key, e.Value})
	}

	rows := len(left)
	if len(right) > rows {
		rows = len(right)
	}
	for i := 0; i < rows; i++ {
		y := pdf.GetY()
		if i < len(left) {
			pdf.SetXY(marginLeft, y)
			r.setBody("B", r.style.FontSize-0.5)
			pdf.CellFormat(22, 4.5, left[i].k+":", "", 0, "R", false, 0, "")
			r.setBody("", r.style.FontSize-0.5)
			pdf.CellFormat(halfW-22, 4.5, " "+left[i].v, "", 0, "L", false, 0, "")
		}
		if i < len(right) {
			pdf.SetXY(marginLeft+halfW, y)
			r.setBody("B", r.style.FontSize-0.5)
			pdf.CellFormat(22, 4.5, right[i].k+":", "", 0, "R", false, 0, "")
			r.setBody("", r.style.FontSize-0.5)
			pdf.CellFormat(halfW-22, 4.5, " "+right[i].v, "", 0, "L", false, 0, "")
		}
		pdf.Ln(4.5)
	}
	pdf.Ln(5)
}

// --- Shared rendering ---

func (r *renderer) renderSection(title string, content func()) {
	aR, aG, aB := r.accent()
	pdf := r.pdf

	switch r.style.Layout {
	case "modern":
		// Accent bar on left + title
		y := pdf.GetY()
		pdf.SetFillColor(aR, aG, aB)
		pdf.Rect(marginLeft, y+1, 3, 8, "F")

		r.setHeader("B", 12)
		pdf.SetTextColor(aR, aG, aB)
		pdf.SetX(marginLeft + 7)
		pdf.CellFormat(contentWidth-7, 10, title, "", 1, "L", false, 0, "")
		pdf.Ln(1)

	default: // classic
		r.setHeader("B", 13)
		pdf.SetTextColor(aR, aG, aB)

		y := pdf.GetY() + 4
		pdf.SetFillColor(aR, aG, aB)
		pdf.Circle(marginLeft+2, y, 2.5, "F")

		pdf.SetX(marginLeft + 8)
		pdf.CellFormat(contentWidth-8, 10, title, "", 1, "L", false, 0, "")

		pdf.SetDrawColor(grayR, grayG, grayB)
		pdf.SetLineWidth(0.3)
		pdf.Line(marginLeft+8, pdf.GetY(), pageWidth-marginRight, pdf.GetY())
		pdf.Ln(3)
	}

	content()
	pdf.Ln(2)
}

func (r *renderer) renderWork(w *Work) {
	pdf := r.pdf
	l := r.label
	fs := r.style.FontSize

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

	switch r.style.Layout {
	case "modern":
		// Two-column: date left, content right
		dateW := r.style.DateColumn
		contentW := contentWidth - dateW - 3

		y := pdf.GetY()
		r.setBody("", fs-0.5)
		pdf.SetTextColor(grayR, grayG, grayB)
		pdf.SetXY(marginLeft, y)
		pdf.MultiCell(dateW, 4.5, period, "", "R", false)
		dateEndY := pdf.GetY()
		if loc != "" {
			pdf.SetX(marginLeft)
			r.setBody("I", fs-1)
			pdf.MultiCell(dateW, 4, loc, "", "R", false)
			dateEndY = pdf.GetY()
		}

		pdf.SetXY(marginLeft+dateW+3, y)
		r.setBody("B", fs+1)
		pdf.SetTextColor(0, 0, 0)
		pdf.CellFormat(contentW, 5.5, w.Title, "", 2, "L", false, 0, "")
		r.setBody("", fs)
		aR, aG, aB := r.accent()
		pdf.SetTextColor(aR, aG, aB)
		pdf.CellFormat(contentW, 4.5, w.Employer, "", 2, "L", false, 0, "")
		pdf.Ln(1)

		if w.Description != "" {
			pdf.SetX(marginLeft + dateW + 3)
			r.setBody("", fs)
			pdf.SetTextColor(50, 50, 50)
			pdf.MultiCell(contentW, 4.5, w.Description, "", "L", false)
		}
		if pdf.GetY() < dateEndY {
			pdf.SetY(dateEndY)
		}
		pdf.Ln(4)

	default: // classic
		r.setBody("", fs)
		pdf.SetTextColor(grayR, grayG, grayB)
		line := period
		if loc != "" {
			line += "  \u2022  " + loc
		}
		pdf.CellFormat(contentWidth, 4.5, line, "", 1, "L", false, 0, "")

		r.setBody("B", fs+1)
		pdf.SetTextColor(0, 0, 0)
		titleW := pdf.GetStringWidth(w.Title) + 2
		pdf.CellFormat(titleW, 5.5, w.Title, "", 0, "L", false, 0, "")
		aR, aG, aB := r.accent()
		r.setBody("", fs+1)
		pdf.SetTextColor(aR, aG, aB)
		remaining := contentWidth - titleW
		if remaining < 20 {
			// Employer doesn't fit on same line — wrap
			pdf.Ln(5.5)
			remaining = contentWidth
		}
		pdf.CellFormat(remaining, 5.5, w.Employer, "", 1, "L", false, 0, "")
		pdf.Ln(1)

		if w.Description != "" {
			r.setBody("", fs)
			pdf.SetTextColor(50, 50, 50)
			pdf.MultiCell(contentWidth, 4.5, w.Description, "", "L", false)
		}
		pdf.Ln(3)
	}
}

func (r *renderer) renderEducation(e *Education) {
	pdf := r.pdf
	l := r.label
	fs := r.style.FontSize

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

	switch r.style.Layout {
	case "modern":
		dateW := r.style.DateColumn
		contentW := contentWidth - dateW - 3

		y := pdf.GetY()
		r.setBody("", fs-0.5)
		pdf.SetTextColor(grayR, grayG, grayB)
		pdf.SetXY(marginLeft, y)
		pdf.MultiCell(dateW, 4.5, period, "", "R", false)
		if loc != "" {
			pdf.SetX(marginLeft)
			r.setBody("I", fs-1)
			pdf.MultiCell(dateW, 4, loc, "", "R", false)
		}

		pdf.SetXY(marginLeft+dateW+3, y)
		r.setBody("B", fs+1)
		pdf.SetTextColor(0, 0, 0)
		pdf.CellFormat(contentW, 5.5, e.Title, "", 2, "L", false, 0, "")
		aR, aG, aB := r.accent()
		r.setBody("", fs)
		pdf.SetTextColor(aR, aG, aB)
		pdf.CellFormat(contentW, 4.5, e.Institution, "", 2, "L", false, 0, "")

		if e.Level != "" {
			r.setBody("I", fs-0.5)
			pdf.SetTextColor(grayR, grayG, grayB)
			pdf.CellFormat(contentW, 4, l.Level+": "+e.Level, "", 2, "L", false, 0, "")
		}
		if e.Description != "" {
			r.setBody("", fs)
			pdf.SetTextColor(50, 50, 50)
			pdf.SetX(marginLeft + dateW + 3)
			pdf.MultiCell(contentW, 4.5, e.Description, "", "L", false)
		}
		pdf.Ln(4)

	default: // classic
		r.setBody("", fs)
		pdf.SetTextColor(grayR, grayG, grayB)
		line := period
		if loc != "" {
			line += "  \u2022  " + loc
		}
		pdf.CellFormat(contentWidth, 4.5, line, "", 1, "L", false, 0, "")

		r.setBody("B", fs+1)
		pdf.SetTextColor(0, 0, 0)
		titleW := pdf.GetStringWidth(e.Title) + 2
		pdf.CellFormat(titleW, 5.5, e.Title, "", 0, "L", false, 0, "")
		aR, aG, aB := r.accent()
		r.setBody("", fs+1)
		pdf.SetTextColor(aR, aG, aB)
		remaining := contentWidth - titleW
		if remaining < 20 {
			pdf.Ln(5.5)
			remaining = contentWidth
		}
		pdf.CellFormat(remaining, 5.5, e.Institution, "", 1, "L", false, 0, "")

		if e.Level != "" {
			r.setBody("I", fs-0.5)
			pdf.SetTextColor(grayR, grayG, grayB)
			pdf.CellFormat(contentWidth, 4.5, l.Level+": "+e.Level, "", 1, "L", false, 0, "")
		}
		if e.Description != "" {
			r.setBody("", fs)
			pdf.SetTextColor(50, 50, 50)
			pdf.MultiCell(contentWidth, 4.5, e.Description, "", "L", false)
		}
		pdf.Ln(3)
	}
}

func (r *renderer) renderLanguages(lang *Languages) {
	aR, aG, aB := r.accent()
	pdf := r.pdf
	l := r.label
	fs := r.style.FontSize

	r.setBody("", fs+1)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(contentWidth, 5.5, l.MotherTongue+": "+strings.Join(lang.MotherTongue, ", "), "", 1, "L", false, 0, "")
	pdf.Ln(2)

	if len(lang.Foreign) == 0 {
		return
	}

	r.setBody("", fs)
	pdf.SetTextColor(grayR, grayG, grayB)
	pdf.CellFormat(contentWidth, 4.5, l.OtherLanguages+":", "", 1, "L", false, 0, "")
	pdf.Ln(2)

	colW := contentWidth / 6
	nameW := colW

	// Table header with accent color
	r.setBody("B", fs-1)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFillColor(aR, aG, aB)
	pdf.CellFormat(nameW, 6, "", "1", 0, "C", true, 0, "")
	pdf.CellFormat(colW, 6, l.Listening, "1", 0, "C", true, 0, "")
	pdf.CellFormat(colW, 6, l.Reading, "1", 0, "C", true, 0, "")
	pdf.CellFormat(colW, 6, l.SpokenProd, "1", 0, "C", true, 0, "")
	pdf.CellFormat(colW, 6, l.SpokenInt, "1", 0, "C", true, 0, "")
	pdf.CellFormat(colW, 6, l.Writing, "1", 1, "C", true, 0, "")

	r.setBody("", fs-1)
	for i, fl := range lang.Foreign {
		if i%2 == 0 {
			pdf.SetFillColor(245, 245, 250)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}
		fill := true
		r.setBody("B", fs-1)
		pdf.SetTextColor(0, 0, 0)
		pdf.CellFormat(nameW, 5.5, fl.Name, "1", 0, "L", fill, 0, "")
		r.setBody("", fs-1)
		pdf.CellFormat(colW, 5.5, fl.Listening, "1", 0, "C", fill, 0, "")
		pdf.CellFormat(colW, 5.5, fl.Reading, "1", 0, "C", fill, 0, "")
		pdf.CellFormat(colW, 5.5, fl.SpokenProduction, "1", 0, "C", fill, 0, "")
		pdf.CellFormat(colW, 5.5, fl.SpokenInteraction, "1", 0, "C", fill, 0, "")
		pdf.CellFormat(colW, 5.5, fl.Writing, "1", 1, "C", fill, 0, "")
	}
	pdf.Ln(2)

	r.setBody("I", fs-2)
	pdf.SetTextColor(grayR, grayG, grayB)
	pdf.MultiCell(contentWidth, 3.5, l.CEFRLegend, "", "L", false)
	pdf.Ln(3)
}

func (r *renderer) renderSkillChips(skills []string) {
	pdf := r.pdf
	aR, aG, aB := r.accent()
	fs := r.style.FontSize

	chipH := 5.5
	chipPad := 3.0
	chipGap := 2.0
	x := marginLeft
	y := pdf.GetY()

	r.setBody("", fs-0.5)
	for _, skill := range skills {
		w := pdf.GetStringWidth(skill) + chipPad*2
		if x+w > pageWidth-marginRight {
			x = marginLeft
			y += chipH + chipGap
		}
		// Light accent background, accent border
		pdf.SetFillColor(aR, aG, aB)
		pdf.SetDrawColor(aR, aG, aB)
		pdf.SetTextColor(255, 255, 255)
		pdf.RoundedRect(x, y, w, chipH, 2, "1234", "FD")
		pdf.SetXY(x+chipPad, y+0.5)
		pdf.CellFormat(w-chipPad*2, chipH-1, skill, "", 0, "C", false, 0, "")
		x += w + chipGap
	}
	pdf.SetY(y + chipH + 2)
}

func (r *renderer) renderSubSection(title, text string) {
	aR, aG, aB := r.accent()
	pdf := r.pdf
	fs := r.style.FontSize

	r.setBody("B", fs+0.5)
	pdf.SetTextColor(aR, aG, aB)
	pdf.CellFormat(contentWidth, 6, title, "", 1, "L", false, 0, "")

	r.setBody("", fs)
	pdf.SetTextColor(50, 50, 50)
	pdf.MultiCell(contentWidth, 4.5, text, "", "L", false)
	pdf.Ln(3)
}

// --- Public entry points ---

func generatePDF(cv *CV, output string) error {
	r, err := newRenderer(cv)
	if err != nil {
		return err
	}
	r.renderCV(cv)

	xmlData := toEuropassXML(cv)
	r.pdf.SetAttachments([]fpdf.Attachment{
		{Content: xmlData, Filename: "Europass_CV.xml", Description: "Europass CV XML data"},
	})
	return r.pdf.OutputFileAndClose(output)
}

func generatePlainPDF(cv *CV, output string) error {
	r, err := newRenderer(cv)
	if err != nil {
		return err
	}
	r.renderCV(cv)
	return r.pdf.OutputFileAndClose(output)
}

func toEuropassXML(cv *CV) []byte {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	b.WriteString(`<SkillsPassport xmlns="http://europass.cedefop.europa.eu/Europass" `)
	b.WriteString(`xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">` + "\n")
	b.WriteString("  <LearnerInfo>\n")

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

	if cv.Headline != "" {
		b.WriteString(fmt.Sprintf("    <Headline><Type><Label>%s</Label></Type></Headline>\n", xmlEsc(cv.Headline)))
	}

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

	b.WriteString("    <Skills>\n")
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
	if len(cv.Digital) > 0 {
		b.WriteString(fmt.Sprintf("      <Computer>%s</Computer>\n", xmlEsc(strings.Join(cv.Digital, ", "))))
	}
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
