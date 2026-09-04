# Localization Core v2 contract

Internationalizer remains a local-first translation compiler for applications
that own their runtime. This contract does not add a hosted review service,
framework runtime, translation CDN, or source-rewriting setup command.

## Locale identity

- Configuration accepts well-formed BCP 47 language tags and preserves their
  spelling in file paths.
- Equality, provider overrides, validation, and state identity use canonical
  BCP 47 form. Canonical-equivalent target locales are duplicates.
- Private-use language tags remain valid. Path separators and malformed tags
  are rejected before resolving target paths.

## Message semantics

- ICU messages are parsed structurally rather than treated as brace-shaped
  interpolation strings.
- Supported arguments are simple interpolation, `select`, `plural`,
  `selectordinal`, `number`, `date`, and `time`.
- Selectors may contain nested messages. Plural arguments support exact-number
  selectors, offsets, and `#`.
- Every `select`, `plural`, and `selectordinal` requires `other`.
- A target must retain source argument names, argument kinds, formatter styles,
  and plural offsets. Select branch identities are semantic and must be
  preserved. Target plurals may add exact selectors and locale-applicable CLDR
  categories.
- Parsing and printing are deterministic so later source-unit hashing can use
  canonical message structure.

## Review transfer

Manifest v2 records provenance origin and human review independently. Provider,
translation-memory, adoption, and pseudo output begin in `needs_review`; only an
exact current source, policy, and target artifact may be approved. Changed
source context, policy, message structure, or target content requires review.

## Source units

Structured documents pass through a format-neutral semantic-unit boundary.
Each unit has a stable semantic ID, kind, translator context, value, and
adapter-owned structure signature; formatting offsets are not identity. Fluent
resources use this boundary for message values, terms, and attributes while
preserving comments and ordering during serialization.

Markdown uses the preamble and each H2 section as independent units, with
invisible target-side markers that preserve identity when sections move or are
inserted.

## Test locales and rich text

- `pseudo` creates deterministic accented (`en-XA`) and bidirectional (`ar-XB`)
  artifacts without a provider or translation-memory lookup.
- ICU and Fluent runtime expressions remain intact while linguistic text is
  transformed.
- `data-l10n-name` identifies semantic rich-text slots. Translators may reorder
  named slots, but element identity, protected attributes, nesting, and
  contained markup must remain compatible.

## Policy identity

The manifest records style-guide, glossary, and prompt-contract components
separately. The combined policy hash includes those components and the provider
settings, but not the rendered prompt text. Refactoring prompt construction
therefore leaves current translations alone; a semantic prompt change requires
an explicit prompt-contract version bump. Style guides are read-only inputs and
are never generated or rewritten by a translation run.

## Compatibility boundary

Existing JSON/YAML key workflows and i18next-v4 plural suffixes remain
supported. Canonical locale comparison may reject configurations that spell the
same locale more than once. ICU or Fluent structural failures are deterministic
validation errors and block provider output from being written.
