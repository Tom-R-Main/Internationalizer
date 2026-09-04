# Using Internationalizer with ExecuFunction

The 0.1.2 regression was a grammar mismatch: the marketing catalog contains
literal shell brace syntax inside HTML code elements. Automatic ICU inference
rejected `{.sift,.claude,.codex,.agents}` before any translation could run.

Configure the web catalog independently from marketing. From ExecuFunction's
repository root, the web mapping is:

```yaml
source_locale: en
# Populate this from the web application's supported locale manifest.
target_locales: [fr]
bundles:
  - id: web
    source: exf-app/web/src/i18n/locales/en.json
    target: exf-app/web/src/i18n/locales/{locale}.json
    format: json
    message_syntax: i18next
```

This profile uses i18next's default double-brace interpolation, including nested
paths and escaping modifiers, documented in the
[i18next interpolation reference](https://www.i18next.com/translation-function/interpolation).
Internationalizer also preserves formatter modifiers and checks v4 plural-key
families against the target locale.

The existing marketing config points at `tmp/english-keys.json` and sibling
locale files. Give that bundle its own explicit syntax, chosen for the runtime
that consumes those files: `plain` for literal strings or `i18next` if double-brace
interpolation is used. Coverage for these temporary catalogs measures those
catalogs only; it does not establish coverage of the web app or of the marketing
pipeline's authoritative inputs. Do not add a guessed production path. Verify
the marketing build's actual input and output mapping before adopting it.

Run `validate --json` and `translate --dry-run` first. Syntax selection removes
the ICU false positive; it can reveal actual translated code or missing
placeholders, which remain errors. Neither dry-run nor validation writes files.

The acceptance fixture `test/acceptance/testdata/execufunction-syntax` exercises
the exact shell-brace pattern alongside i18next interpolation and an explicit
ICU bundle. Its lifecycle covers both pseudo strategies, adoption without
approval, explicit approval, and invalidation after a syntax-policy change.
