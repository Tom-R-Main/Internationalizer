# Shared i18n Conventions

## Interpolation
- Syntax: `{{variableName}}` — never translate the variable name
- OK to reorder text around variables to match target language word order
- Variables may appear mid-sentence

## HTML in Strings
- Preserve all HTML tags exactly (tag name, attributes, self-closing)
- OK to move tags to match word order
- Never add or remove tags
- Paired tags must wrap semantically equivalent text

## String Types
| Type | Style | Example |
|------|-------|---------|
| Button/action | Imperative verb, no period | "Save", "Create Task" |
| Label | Noun phrase, no period | "Timezone", "Start time" |
| Status/toast | Past tense or state, no period | "Settings updated" |
| Error | Empathetic + actionable, period OK | "Something went wrong. Please try again." |
| Placeholder | Hint text, no period, lowercase OK | "Search..." |

## Pluralization
Uses CLDR-based suffixes: `_one`, `_other`, `_few`, `_many`, `_zero`, `_two`.
Provide all plural forms required by the target language.
