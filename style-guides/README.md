# Translation style-guide policy

Files in this directory are durable translation policy. Internationalizer reads
them when it builds a locale prompt, but a translation run never edits them.
Change a guide only when the language rule itself changes, not when application
code or source copy changes.

`_conventions.md` contains rules shared by every locale. Each `{locale}.md`
file contains only language-specific grammar, register, typography, plural, and
terminology decisions. Keep product facts in the source document instead of
copying them into a language guide.

Guide changes are intentionally visible: the manifest hashes each guide and
marks that locale's existing translations as policy-stale. Internal prompt
refactors do not have that effect unless the prompt-contract version changes.

## Reference sources

- Bangla, Czech, Danish, Dutch, Finnish, French, German, Greek, Hebrew, Hindi,
  Indonesian, Italian, Korean, Malay, Polish, Portuguese (Brazil), Punjabi,
  Romanian, Spanish, Swedish, Telugu, Thai, Turkish, Ukrainian, and Vietnamese:
  [Microsoft Localization Style Guides](https://learn.microsoft.com/en-us/globalization/reference/microsoft-style-guides)
- Japanese: JTF Japanese Standard Style Guide for Translation, plus the
  supplied Japanese professional-writing reference
- Russian: Microsoft Russian Style Guide and Nora Gal's writing guidance
- Arabic: Microsoft Arabic (Saudi Arabia) Style Guide, Mahdi Alosh's usage
  reference, and the supplied Arabic translation reference
- Simplified and Traditional Chinese: the supplied Chinese editing and
  translation references; regional terminology remains separate in each guide
- Cantonese: the repository's Cantonese guide; Microsoft does not publish a
  matching Cantonese guide in the collection above

External references inform these concise project rules; they are not copied
into the repository. Review model-proposed guide changes as policy changes,
then update the guide deliberately and refresh affected translations.

`scripts/generate-style-guides.mjs` writes review candidates under the ignored
`.internationalizer/` directory by default. It changes tracked guides only
when invoked with `--apply`.
