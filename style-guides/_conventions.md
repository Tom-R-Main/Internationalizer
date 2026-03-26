# Shared Translation Conventions

These rules apply to ALL target languages. Per-language guides may override specific points.

## Interpolation Variables
- Syntax: `{{variableName}}` — never translate the variable name inside braces
- OK to reorder text around variables to match target language word order
- Variables may appear mid-sentence

## HTML Tags
- Preserve all HTML tags exactly (tag name, attributes, self-closing)
- OK to move tags to match target word order
- Never add or remove tags
- Never translate HTML attribute values
- Paired tags must wrap semantically equivalent text

## Brand & Technical Terms (keep in English)
- Internationalizer (product name)
- CLI, API, JSON, YAML, JSONL, npm, Go
- GitHub, GitHub Actions
- i18n, i18next, next-intl, vue-i18n
- LLM, AI
- CLDR, AGPL-3.0

## String Type Conventions
| Type | Style | Example |
|------|-------|---------|
| Button/action | Imperative verb, no period | "Save", "Delete" |
| Label | Noun phrase, no period | "Settings", "Language" |
| Status/toast | Past tense or state, no period | "Settings updated" |
| Error | Empathetic + actionable, period OK | "Something went wrong. Please try again." |
| Heading | Sentence case | "Why Internationalizer?" |
| Body text | Complete sentences | "Most i18n tools are either..." |
| Code/commands | Never translate | `npm install`, `internationalizer translate` |

## Markdown Preservation
- Preserve all Markdown formatting (headings, links, code blocks, lists, tables)
- Do not translate URLs, file paths, or code blocks
- Do not translate text inside backticks (inline code)
- Preserve badge image URLs and link targets exactly

## Quality Checklist
- All `{{variables}}` preserved with identical names
- All HTML tags preserved
- All Markdown formatting intact
- Brand/technical terms in English
- Terminology matches the per-language glossary
- String type conventions followed
- No untranslated English (except brand terms, code, variables)
- Appropriate plural forms for the target language
