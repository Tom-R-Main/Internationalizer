package validate

import "sort"

// FindingCode is a stable machine-readable validation outcome.
type FindingCode string

const (
	CodeMissingKey                 FindingCode = "missing_key"
	CodeExtraKey                   FindingCode = "extra_key"
	CodeBlankTranslation           FindingCode = "blank_translation"
	CodeSourceIdentical            FindingCode = "source_identical"
	CodeProtectedStructureMismatch FindingCode = "protected_structure_mismatch"
	CodeGlossaryViolation          FindingCode = "glossary_violation"
	CodePluralFormMissing          FindingCode = "plural_form_missing"
	CodeUntracked                  FindingCode = "untracked"
	CodeSourceStale                FindingCode = "source_stale"
	CodePolicyStale                FindingCode = "policy_stale"
	CodeTargetModified             FindingCode = "target_modified"
)

// Severity determines whether a finding fails validation.
type Severity string

const (
	SeverityWarning Severity = "warning"
	SeverityError   Severity = "error"
)

// Finding describes one stable validation outcome for automation and humans.
type Finding struct {
	Code     FindingCode `json:"code"`
	Severity Severity    `json:"severity"`
	Key      string      `json:"key,omitempty"`
	Message  string      `json:"message"`
	Expected []string    `json:"expected,omitempty"`
	Actual   []string    `json:"actual,omitempty"`
}

func sortFindings(findings []Finding) {
	sort.SliceStable(findings, func(i, j int) bool {
		if findings[i].Key != findings[j].Key {
			return findings[i].Key < findings[j].Key
		}
		if findings[i].Code != findings[j].Code {
			return findings[i].Code < findings[j].Code
		}
		return findings[i].Message < findings[j].Message
	})
}
