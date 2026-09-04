# Changelog

## Unreleased

- Reject duplicate JSON object members and flattened-key collisions before catalog content can be lost. These integrity checks apply with or without `--strict`, including LLM translation responses and catalog rewrite operations.
- Report malformed catalogs during discovery and configuration checks. JSON errors identify conflicting member paths without printing translation values.
- Refuse ambiguous existing JSON pseudolocale targets even with `--force`; repair the catalog before replacing it.

## 0.2.0 - 2026-09-04

### Added

- Discover configured and uncovered catalogs across nested apps with `detect --json`, including runtime evidence, syntax suggestions, and unresolved source choices.
- Inspect resolved bundles, locale targets, syntax provenance, and offline credential presence with `config check --json`.
- Review and apply saved configuration proposals with `config plan` and `config apply`. Plans preserve provider settings, locale overrides, glossary paths, and bundle identities; application checks drift and returns a receipt tied to the plan.
- Discover command arguments, schemas, side effects, and next steps with `commands --json`.
- Preview translation work with `translate --dry-run --json`, separating planned work, blocked jobs, generated translations, and provider calls.

### Changed

- **JSON compatibility:** `validate --json` now returns a versioned envelope. Read reports from `data.reports` instead of the previous top-level array. Check `schema_version` before consuming command-specific data.
- Discovery and configuration diagnostics, translation results, and validation reports support bounded JSON output and scope filters. Structured errors include recovery argument arrays.
- Framework detection is advisory: i18next ICU integrations and mixed frameworks require explicit decisions. Temporary marketing catalogs are flagged for confirmation rather than replaced.

### Fixed

- Explain when automatic syntax detection interprets literal code braces as ICU, and point to explicit runtime profiles. Explicit ICU parsing and placeholder/protected-code validation remain strict.
- Base partial-failure reporting on persisted catalog or manifest updates, not staged translations from a failed batch sequence. Job results distinguish catalog writes from manifest updates.

See [the configuration workflow and JSON contract](docs/cli-onboarding.md) for setup, migration, and retry behavior.
