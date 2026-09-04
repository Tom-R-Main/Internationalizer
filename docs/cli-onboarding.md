# Configuration workflow and JSON contract

Run setup commands from the project root. Inspection resolves relative catalog
paths against that root, matching translation; an explicit config file does not
change the base directory for its paths.

## Discover and inspect

```sh
internationalizer commands --json
internationalizer commands --command 'config plan' --json
internationalizer detect --json
internationalizer config check --json
internationalizer config check --bundle web --locale fr --json
```

`detect` reports source candidates, configured bundle IDs, runtime evidence,
suggested syntax, and uncertainty. It never edits configuration. A catalog under
`tmp/` is a source-ownership question, not automatically an error.

`config check` adds resolved sources, target templates and locale-specific paths,
syntax provenance, diagnostics, and credential-presence checks. It is offline:
`provider_verified: false` means no provider request was made, even when a
credential is present. Neither command emits credential values.

i18next dependency evidence suggests `i18next` syntax. An `i18next-icu`
dependency or localization-module reference leaves the profile unresolved until
you confirm the plugin registration. Discovery does not infer grammar from JSON
or claim that an absent static reference proves there is no ICU integration.

Source-locale filenames such as `en.json` and paths such as
`locales/en/common.json` are catalog candidates. Other authoritative sources can
be supplied explicitly with `--add-bundle` and `--target`.

Scan limits are 50,000 entries, depth 14, and 8 MiB per inspected catalog.
Dependency, build, hidden, and data directories and symlink traversals are
excluded. Large localization modules are skipped. A truncated scan is not proof
that all catalogs or integrations have been found.

## Plan a change

This example extends an existing marketing-only config whose `source_path` is
`tmp/english-keys.json`:

```sh
internationalizer config plan --json \
  --add-bundle web=exf-app/web/src/i18n/locales/en.json \
  --syntax web=i18next \
  --syntax default=i18next \
  --confirm-source tmp/english-keys.json \
  --out config-plan.json
```

Each `--add-bundle ID=SOURCE` is an explicit catalog selection. For a discovered
source, the proposal uses its target template; use `--target ID=TEMPLATE` to
override it. New JSON bundles need an explicit `--syntax ID=PROFILE` decision.
Profiles are `plain`, `i18next`, and `icu`; an ambiguous `auto` choice remains
unresolved. Existing `source_path` configurations become a bundle named
`default`, preserving the existing translation-state identity.

For a new project, select target locales too:

```sh
internationalizer config plan --config .internationalizer.yml \
  --add-bundle web=src/locales/en.json \
  --syntax web=i18next --locale fr --locale ja --out config-plan.json --json
```

`--source-locale` and `--locale` also change existing locale settings when
explicitly supplied. Omitting them preserves the current locale set and overrides.

Planning does not modify the config. `--out` creates a new, owner-readable plan
file and refuses to overwrite an existing file. Without `--out`, the proposal is
returned only in the command output. A `needs_decision` status means inspect
`required_decisions` and create a new plan with the missing selections.

Plans contain proposed YAML, a diff, configuration fingerprints, and observations
of source/runtime evidence. Existing provider settings, locale overrides,
glossary paths, comments, and ordinary unknown settings are preserved. Configs
with YAML aliases/merge keys or possible inline credentials are rejected rather
than rewritten. Use environment-variable references for credentials. Treat
saved plans as private project configuration, not public attachments.

Inspection can report a home-directory fallback config. Plan/apply only edits
files inside the project root; it does not silently replace a home config with a
new local config. Choose an explicit local `--config` path to start fresh.

## Apply and verify

### Retarget an existing bundle

Use an explicit update decision to change an existing bundle's target template:

```sh
internationalizer config plan --update-bundle default \
  --target 'default=tmp/translations-{locale}.json' \
  --confirm-source tmp/english-keys.json --out repair-plan.json --json
```

The example assumes that the marketing runtime syntax is already configured.
If it remains ambiguous, also supply the appropriate `--syntax default=PROFILE`.
`--update-bundle` requires an existing ID and a matching `--target`; it cannot
be combined with `--add-bundle` for the same ID. Source, format, syntax, locale,
provider, and other settings remain unchanged unless separately selected by an
existing explicit flag. Legacy `source_path` configuration becomes the stable
`default` bundle when retargeted. Renaming or deleting bundles is not supported.

Review the proposed YAML and diff before applying. Application only updates the
configuration: it does not copy catalogs, replace links, translate messages, or
approve translations. A destination with different content must still satisfy
the existing validation and approval checks.

When a target is a symlink, planning reports the offending path, bundle, and
locale. A bounded metadata check may suggest an existing in-project destination;
it does not inspect outside-project contents or select a path automatically.
External, dangling, cyclic, and unsafe ancestor links remain rejected.

An explicit safe replacement can repair the config even when the old targets
are unsafe. The replacement is checked independently, so the old link is not
an input to the saved plan and its identity is not attested by the receipt.
The saved-plan schema is unchanged. Apply rechecks replacement paths and still
rejects symlinks introduced after planning; normal drift and lock checks apply.

### Apply a reviewed proposal

```sh
internationalizer config apply --plan config-plan.json --no-input --json
internationalizer config check --json
internationalizer translate --dry-run --json
internationalizer validate --bundle web --locale fr --json
```

Review the saved plan before invoking `config apply`. That invocation requests
the exact config mutation; `--no-input` only disables prompts. A plan hash checks
integrity but is not a signature or authorization grant. Never apply an
unreviewed plan from an untrusted source.

Apply checks fingerprints, rejects unresolved decisions and unsafe paths, and
uses a per-config lock and atomic replacement. Its receipt includes the plan ID,
changed paths, resulting config fingerprint, and `applied` or `already_applied`
status. It does not translate, approve, or edit catalog files.

An `already_applied` receipt confirms the config fingerprint only; its
`observations_revalidated: false` is not a fresh source or provider check. Run
`config check` and a dry-run again if sources changed after the first application.

After a stale-plan error, inspect current state and generate a new proposal.
Do not edit the saved JSON to bypass its fingerprint. After an interrupted apply,
inspect the lock and current config before retrying; do not delete another
process's lock. Cooperative locking protects concurrent CLI applications; it is
not a transaction against arbitrary filesystem writers.

`translate --dry-run --json` performs no writes or provider calls. It reports
planned keys separately from generated keys and blocked jobs. An explicit syntax
profile resolves grammar ambiguity, but it does not waive placeholder, plural,
or protected-code validation. Explicit ICU remains strict; malformed ICU is not
silently reinterpreted as plain text.

## JSON, filtering, and failures

### Catalog integrity

JSON catalogs must have unique object members and unambiguous flattened keys.
For example, `{"a.b":"First","a":{"b":"Second"}}` gives two values the
same catalog identity, `a.b`, and is rejected. Duplicate members are rejected
even when their values match or their names use different JSON escapes.
These checks apply without `--strict`; sorting keys or choosing the last value
would discard content, not repair it.

Errors use `json_duplicate_member` or `json_flattened_key_collision`. Source
errors include member locations under `errors[].details`; target validation
findings include `path` and `other_path`. Locations are JSON pointers into the
catalog; the containing error or report identifies the catalog file. Resolve
the conflicting members explicitly, preserving the intended messages and
placeholders. Internationalizer does not choose a surviving value for you.

Discovery keeps malformed catalog candidates visible with `parse_error_code`.
`detect` remains advisory; `config check` fails when an error diagnostic exists,
even if presentation filters hide it. The same integrity checks protect LLM
response parsing and JSON catalog rewrites. `pseudo --force` bypasses ownership
checks, not JSON integrity checks.

Translation jobs expose a structured `input_error` when a target catalog or
provider response fails JSON integrity checks. The run can still report
`translation_failed` because other jobs may have completed; inspect the job's
error code and persistence flags before retrying.

### Output envelope

The onboarding commands, `commands`, `translate`, and `validate` use this envelope
when `--json` is supplied:

```json
{
  "schema_version": 1,
  "status": "planned",
  "data": {},
  "errors": []
}
```

Parse `schema_version` before consuming command-specific `data`. Errors use the
same envelope, including argument and configuration failures. Stdout contains
JSON only; command diagnostics must not be parsed from human prose. Each error
has a stable `code`, a message, and recovery actions with `argv`, `side_effects`,
and `required_decisions`. Review those decisions before executing a recovery
action.

`detect` and `config check` put results under `data.inspection`, with `total` and
`matched` counts. Use `--bundle`, `--locale`, `--finding-code`, and `--limit` to
bound output. Their default limit is 50 entries per result section. Translation
defaults to 100 returned jobs; validation defaults to 100 reports/detail items.
Use `--limit 0` for unbounded presentation. Limits and finding filters never
turn an underlying failure into a pass. They do not reduce work already needed
to inspect selected catalogs.

`validate --json` previously returned an array. It now returns reports under
`data.reports`, with counts and filtering metadata. Update existing consumers
before using this version.

Exit code 0 means the command completed, not that every required decision was
made. A plan can exit 0 with `needs_decision`. Exit code 1 means execution failed,
configuration checks blocked, or validation failed. Inspect the structured
status and errors as well as the exit code.

Common errors include `invalid_arguments`, `unknown_bundle`, `unknown_locale`,
`config_invalid`, `stale_plan`, `invalid_plan`, `decisions_required`,
`unsafe_path`, `apply_locked`, `credentials_missing`, and `translation_failed`.
Configuration findings include `UNCOVERED_CATALOG`,
`SOURCE_CONFIRMATION_REQUIRED`, and `AUTO_SYNTAX_AMBIGUOUS`.

Configuration validity, planned translation, generated translation, structural
validation, and human approval are separate states. A successful provider call
does not establish approval. Translation can report `partial_failure` when some
work succeeded; inspect its jobs and retained state before retrying. Provider
requests are not inherently idempotent, even when completed local work can be
reused.

Translation jobs report `catalog_written` and `manifest_updated` only after
those writes succeed. `summary.persisted_jobs` counts jobs with either kind of
retained update; an error yields `partial_failure` only when this count is
positive. Generated and cached key counters can include staged work discarded
after a later batch fails. They are not persistence receipts. A catalog write
can succeed before a later state update fails, so inspect both flags before
retrying; the catalog, translation memory, and manifest are not one transaction.
