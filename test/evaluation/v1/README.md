# Internationalizer evaluation fixtures v1

This directory freezes the shared evidence contract used by deterministic
validation and future translation-memory experiments.

`cases.json` is an append-only corpus within schema version 1. Every case has a
stable `id`, a `kind`, its relevant inputs, and expected stable finding or
eligibility codes. New case kinds may add fields, but existing meanings must not
change without a schema-version bump.

The fixture project is intentionally invalid. `baseline.json` records how the
CLI at commit `c95cfbdda62f569e353bac59f4a8dbd9b7d5834e` observed it before strict
validation existed. Provider fields are null because this baseline is fully
deterministic and makes no model calls.
