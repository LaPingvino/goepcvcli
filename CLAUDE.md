# goepcvcli — Europass CV CLI Tool

Go CLI for managing, tailoring, and generating Europass-format CVs.
CV data is stored as `output/cv.json` — all modifications go through the CLI, no editor needed.
All personal data lives in the `output/` directory (gitignored).

## Interactive Mode

```bash
goepcvcli -i                                # guided menu for humans
```

Launches a menu-driven interface that walks through all operations
with prompts, choices, defaults, and confirmations — no flags needed.

## Quick Reference

```bash
# View CV
goepcvcli show                              # full summary
goepcvcli show -s experience                # specific section
goepcvcli show --json                       # raw JSON for piping

# Modify fields
goepcvcli set headline "New Headline"
goepcvcli set phone "+351 913044570"

# Add entries
goepcvcli add work --title "Developer" --employer "Acme" --from "JAN 2024" \
  --description "Building things" --tags "dev,go"
goepcvcli add skill "Kubernetes" "Terraform"
goepcvcli add language --name Japanese --all B1
goepcvcli add contact Matrix "@joop:chat.kiefte.eu"

# Update existing entries (by index or name)
goepcvcli update work 0 --description "Updated description" --tags "dev,go,llm"
goepcvcli update language Portuguese --spoken-interaction C1

# Remove entries
goepcvcli remove work 3
goepcvcli remove skill Docker
goepcvcli remove language Afrikaans

# Tailor for a specific job (does NOT modify cv.json)
goepcvcli tailor --tags dev,go,architecture \
  --headline "Go Developer | Systems Architecture" \
  --output output/dev-cv.json --pdf output/dev-cv.pdf

# Generate output from any JSON
goepcvcli generate                                              # default: output/cv.json -> output/cv.pdf
goepcvcli generate -f output/dev-cv.json -o output/dev-cv.pdf   # specific files
goepcvcli generate -f output/cv.json -o output/cv.xml --format xml   # standalone Europass XML
goepcvcli generate -f output/cv.json -o output/cv.pdf --format plain # PDF without embedded XML
goepcvcli generate -f output/cv.json -o output/cv-de.pdf --lang de   # German labels
```

## Sections for show/set

`personal`, `headline`, `experience`, `education`, `languages`, `digital`, `skills`

## Tags System

Work entries have tags for filtering with `tailor`. Common tags:
- `dev`, `go`, `gcp`, `architecture` — development roles
- `support`, `legal-tech`, `b2b`, `microsoft` — support roles  
- `leadership`, `ngo`, `international` — leadership/org roles
- `content`, `web`, `process` — content/web roles
- `devops`, `sql` — ops roles

## i18n

Set `"lang": "de"` in the CV JSON or use `--lang de` on generate.
Supported: all 24 EU languages (bg, cs, da, de, el, en, es, et, fi, fr, ga, hr,
hu, it, lt, lv, mt, nl, pl, pt, ro, sk, sl, sv) plus eo and tok.

## PDF Generation

Uses DejaVu Sans Condensed (UTF-8) from `/usr/share/fonts/TTF/`.
Europass XML is embedded as a PDF attachment for machine readability.
Use `--format xml` to export standalone XML, `--format plain` for PDF without XML.

## LLM Workflow

This tool is designed for LLM-driven CV management. Typical flow:
1. `show --json` to read current state
2. `set`, `add`, `update`, `remove` to modify
3. `tailor` to create job-specific variants
4. `generate` to produce PDFs

All commands use flags (no interactive prompts needed).
Use `--json` on show for machine-readable output.
All files default to `output/` directory — use `-f` to override.
