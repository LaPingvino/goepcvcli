package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"unicode"

	"github.com/go-pdf/fpdf"
)

// resolvedFont holds the resolved TTF paths for a font family.
type resolvedFont struct {
	family  string // internal fpdf family name
	dir     string // directory for fpdf font path
	regular string // regular weight filename
	bold    string // bold weight filename
	italic  string // italic weight filename
}

// fcMatch resolves a fontconfig name + style to a TTF file path.
// Falls back to platform directory search if fc-match is unavailable.
func fcMatch(name, style string) string {
	query := name
	if style != "" {
		query += ":" + style
	}
	out, err := exec.Command("fc-match", "--format=%{file}", query).Output()
	if err != nil {
		return findFontFile(name, style)
	}
	path := strings.TrimSpace(string(out))
	if !strings.HasSuffix(strings.ToLower(path), ".ttf") {
		return ""
	}
	return path
}

// fontSearchDirs returns platform-specific directories to search for fonts.
func fontSearchDirs() []string {
	home, _ := os.UserHomeDir()
	switch runtime.GOOS {
	case "darwin":
		dirs := []string{"/System/Library/Fonts", "/Library/Fonts"}
		if home != "" {
			dirs = append(dirs, filepath.Join(home, "Library/Fonts"))
		}
		return dirs
	case "windows":
		windir := os.Getenv("WINDIR")
		if windir == "" {
			windir = `C:\Windows`
		}
		return []string{filepath.Join(windir, "Fonts")}
	default:
		dirs := []string{"/usr/share/fonts", "/usr/local/share/fonts"}
		if home != "" {
			dirs = append(dirs, filepath.Join(home, ".local/share/fonts"))
			dirs = append(dirs, filepath.Join(home, ".fonts"))
		}
		return dirs
	}
}

// findFontFile searches platform font directories for a TTF matching name+style.
// Fallback for when fc-match is not available (macOS, Windows).
func findFontFile(name, style string) string {
	nameLC := strings.ToLower(strings.ReplaceAll(name, " ", ""))
	nameHyphen := strings.ToLower(strings.ReplaceAll(name, " ", "-"))

	styleSuffix := ""
	switch style {
	case "bold":
		styleSuffix = "bold"
	case "italic":
		styleSuffix = "italic"
	}

	var result string
	for _, dir := range fontSearchDirs() {
		filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() || result != "" {
				return nil
			}
			if !strings.HasSuffix(strings.ToLower(path), ".ttf") {
				return nil
			}
			base := strings.ToLower(strings.TrimSuffix(d.Name(), ".ttf"))
			clean := strings.ReplaceAll(strings.ReplaceAll(base, " ", ""), "-", "")

			if !strings.Contains(clean, nameLC) && !strings.Contains(base, nameHyphen) {
				return nil
			}
			if styleSuffix == "" {
				if strings.Contains(clean, "bold") || strings.Contains(clean, "italic") || strings.Contains(clean, "oblique") {
					return nil
				}
			} else if !strings.Contains(clean, styleSuffix) {
				return nil
			}
			result = path
			return fs.SkipAll
		})
		if result != "" {
			return result
		}
	}
	return ""
}

// resolveFont resolves a FontSpec to actual TTF paths for regular, bold, italic.
func resolveFont(spec FontSpec, familyID string) resolvedFont {
	rf := resolvedFont{family: familyID}

	if spec.Path != "" {
		rf.dir = filepath.Dir(spec.Path)
		base := filepath.Base(spec.Path)
		rf.regular, rf.bold, rf.italic = base, base, base
		return rf
	}

	regular := fcMatch(spec.Name, "")
	bold := fcMatch(spec.Name, "bold")
	italic := fcMatch(spec.Name, "italic")

	if regular == "" {
		regular = fcMatch("DejaVu Sans Condensed", "")
		bold = fcMatch("DejaVu Sans Condensed", "bold")
		italic = fcMatch("DejaVu Sans Condensed", "italic")
	}

	if regular != "" {
		rf.dir = filepath.Dir(regular)
		rf.regular = filepath.Base(regular)
	}
	if bold != "" {
		rf.bold = filepath.Base(bold)
	} else {
		rf.bold = rf.regular
	}
	if italic != "" {
		rf.italic = filepath.Base(italic)
	} else {
		rf.italic = rf.regular
	}
	return rf
}

// registerFont registers a resolved font family with fpdf using absolute paths.
func registerFont(pdf *fpdf.Fpdf, rf resolvedFont) {
	regular := filepath.Join(rf.dir, rf.regular)
	bold := filepath.Join(rf.dir, rf.bold)
	italic := filepath.Join(rf.dir, rf.italic)

	// fpdf joins fontpath + filename, so set fontpath to "/" and use absolute paths
	pdf.SetFontLocation("/")
	pdf.AddUTF8Font(rf.family, "", regular)
	pdf.AddUTF8Font(rf.family, "B", bold)
	pdf.AddUTF8Font(rf.family, "I", italic)

	// Bold italic: try common naming patterns, fall back to bold
	biBase := strings.Replace(rf.bold, "-Bold.", "-BoldOblique.", 1)
	biBase = strings.Replace(biBase, "-Bold.ttf", "-BoldOblique.ttf", 1)
	if biBase == rf.bold {
		biBase = strings.Replace(rf.bold, "-Bold.ttf", "-BoldItalic.ttf", 1)
	}
	biPath := filepath.Join(rf.dir, biBase)
	if _, err := os.Stat(biPath); err != nil {
		biPath = bold // fall back to bold
	}
	pdf.AddUTF8Font(rf.family, "BI", biPath)
}

// pdfFonts holds all resolved and registered font families for a PDF.
type pdfFonts struct {
	header   resolvedFont
	body     resolvedFont
	fallback []resolvedFont
}

// setupFonts resolves and registers all fonts for a CV's style.
func setupFonts(pdf *fpdf.Fpdf, style Style) pdfFonts {
	pf := pdfFonts{}

	pf.header = resolveFont(style.HeaderFont, "header")
	registerFont(pdf, pf.header)

	if style.BodyFont.Name == style.HeaderFont.Name && style.BodyFont.Path == style.HeaderFont.Path {
		pf.body = resolvedFont{
			family:  "header",
			dir:     pf.header.dir,
			regular: pf.header.regular,
			bold:    pf.header.bold,
			italic:  pf.header.italic,
		}
	} else {
		pf.body = resolveFont(style.BodyFont, "body")
		registerFont(pdf, pf.body)
	}

	for i, fb := range style.Fallback {
		spec := FontSpec{}
		if strings.HasSuffix(strings.ToLower(fb), ".ttf") || strings.Contains(fb, "/") {
			spec.Path = fb
		} else {
			spec.Name = fb
		}
		rf := resolveFont(spec, fmt.Sprintf("fb%d", i))
		if rf.dir != "" {
			registerFont(pdf, rf)
			pf.fallback = append(pf.fallback, rf)
		}
	}

	return pf
}

// segmentText splits text into runs by font coverage.
// Characters needing fallback are grouped into separate segments.
func (pf *pdfFonts) segmentText(text, primaryFamily string) []textSegment {
	if len(pf.fallback) == 0 {
		return []textSegment{{text: text, family: primaryFamily}}
	}

	var segments []textSegment
	var current strings.Builder
	currentFamily := primaryFamily

	for _, r := range text {
		target := primaryFamily
		if needsFallbackFont(r) {
			target = pf.fallback[0].family
		}
		if target != currentFamily && current.Len() > 0 {
			segments = append(segments, textSegment{text: current.String(), family: currentFamily})
			current.Reset()
		}
		currentFamily = target
		current.WriteRune(r)
	}
	if current.Len() > 0 {
		segments = append(segments, textSegment{text: current.String(), family: currentFamily})
	}
	return segments
}

type textSegment struct {
	text   string
	family string
}

// needsFallbackFont returns true for characters likely missing from standard text fonts.
func needsFallbackFont(r rune) bool {
	return r > 0xFFFF || unicode.Is(unicode.So, r)
}

// writeWithFallback renders text using Write(), switching fonts for fallback segments.
func (pf *pdfFonts) writeWithFallback(pdf *fpdf.Fpdf, h float64, text, primaryFamily, style string) {
	segments := pf.segmentText(text, primaryFamily)
	for _, seg := range segments {
		sz, _ := pdf.GetFontSize()
		pdf.SetFont(seg.family, style, sz)
		pdf.Write(h, seg.text)
	}
	sz, _ := pdf.GetFontSize()
	pdf.SetFont(primaryFamily, style, sz)
}
