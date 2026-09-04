// Package onboarding exposes offline, evidence-backed configuration discovery.
package onboarding

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Tom-R-Main/Internationalizer/internal/config"
	"github.com/Tom-R-Main/Internationalizer/internal/formats"
	"github.com/Tom-R-Main/Internationalizer/internal/message"
	"github.com/Tom-R-Main/Internationalizer/internal/protectedtext"
	"github.com/Tom-R-Main/Internationalizer/internal/validate"
	"gopkg.in/yaml.v3"
)

type Evidence struct {
	Path   string `json:"path"`
	Kind   string `json:"kind"`
	Detail string `json:"detail"`
}

type Recovery struct {
	Argv              []string `json:"argv"`
	SideEffects       []string `json:"side_effects"`
	RequiredDecisions []string `json:"required_decisions,omitempty"`
}

type Diagnostic struct {
	Code              string     `json:"code"`
	Severity          string     `json:"severity"`
	Message           string     `json:"message"`
	Bundle            string     `json:"bundle,omitempty"`
	Key               string     `json:"key,omitempty"`
	Evidence          []Evidence `json:"evidence,omitempty"`
	RequiredDecisions []string   `json:"required_decisions,omitempty"`
	Recovery          []Recovery `json:"recovery,omitempty"`
}

type Candidate struct {
	ID                   string         `json:"id"`
	Source               string         `json:"source"`
	Target               string         `json:"target"`
	Format               string         `json:"format"`
	Framework            string         `json:"framework"`
	SuggestedSyntax      message.Syntax `json:"suggested_syntax,omitempty"`
	Evidence             []Evidence     `json:"evidence"`
	Uncertainty          string         `json:"uncertainty"`
	ConfiguredBundles    []string       `json:"configured_bundles"`
	RequiresConfirmation bool           `json:"requires_confirmation"`
}

type ResolvedBundle struct {
	ID            string            `json:"id"`
	Source        string            `json:"source"`
	Target        string            `json:"target"`
	Format        string            `json:"format"`
	Framework     string            `json:"framework"`
	MessageSyntax message.Syntax    `json:"message_syntax"`
	Locales       []string          `json:"locales"`
	Targets       map[string]string `json:"targets"`
	Provenance    map[string]string `json:"provenance"`
	Evidence      []Evidence        `json:"evidence"`
}

type Credential struct {
	Locale              string `json:"locale"`
	Provider            string `json:"provider"`
	EnvironmentVariable string `json:"environment_variable"`
	Present             bool   `json:"present"`
	ProviderVerified    bool   `json:"provider_verified"`
}

type Inspection struct {
	Root          string           `json:"root"`
	ConfigPath    string           `json:"config_path"`
	ConfigExists  bool             `json:"config_exists"`
	SourceLocale  string           `json:"source_locale"`
	TargetLocales []string         `json:"target_locales"`
	Candidates    []Candidate      `json:"candidates"`
	Bundles       []ResolvedBundle `json:"bundles"`
	Diagnostics   []Diagnostic     `json:"diagnostics"`
	Credentials   []Credential     `json:"credentials"`
	Truncated     bool             `json:"truncated"`
}

type runtimeEvidence struct {
	i18next  bool
	icu      bool
	nextIntl bool
	vueI18n  bool
	evidence []Evidence
}

const maxDiscoveryBytes = 8 << 20

// Scan is read-only and never changes the process working directory. Relative
// config and bundle paths resolve against root, matching CLI execution there.
// Framework evidence is static; absence of an ICU import is not proof of absence.
func Scan(root, configPath string) (*Inspection, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	info, err := os.Stat(root)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("discovery root is not a directory")
	}
	root, err = filepath.EvalSymlinks(root)
	if err != nil {
		return nil, err
	}
	result := &Inspection{Root: root, Candidates: []Candidate{}, Bundles: []ResolvedBundle{}, Diagnostics: []Diagnostic{}, Credentials: []Credential{}}
	result.ConfigPath = resolveConfigPath(root, configPath)
	var raw config.Config
	data, err := safeRead(result.ConfigPath)
	if err == nil {
		result.ConfigExists = true
		if yaml.Unmarshal(data, &raw) != nil {
			result.Diagnostics = append(result.Diagnostics, Diagnostic{Code: "CONFIG_INVALID", Severity: "error", Message: "Configuration YAML could not be parsed; inspect the configuration file."})
			return result, fmt.Errorf("configuration YAML could not be parsed")
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return result, fmt.Errorf("reading configuration: %w", err)
	} else {
		result.Diagnostics = append(result.Diagnostics, Diagnostic{Code: "CONFIG_NOT_FOUND", Severity: "warning", Message: "No configuration found; choose authoritative catalogs and target locales before planning configuration."})
	}
	cfg := raw
	cfg.ApplyDefaults()
	result.SourceLocale, result.TargetLocales = cfg.SourceLocale, append([]string{}, cfg.TargetLocales...)
	if result.ConfigExists {
		if err := cfg.ValidateProject(); err != nil {
			result.Diagnostics = append(result.Diagnostics, Diagnostic{Code: "CONFIG_INVALID", Severity: "error", Message: err.Error()})
		}
	}
	files, runtimes := discoverFiles(root, result)
	for _, path := range files {
		if supportCatalogPath(root, path, cfg) {
			continue
		}
		target, ok := catalogTarget(path, cfg.SourceLocale)
		if !ok {
			continue
		}
		format, formatErr := formats.FormatForFile(path)
		if formatErr != nil {
			continue
		}
		content, readErr := safeRead(filepath.Join(root, path))
		if readErr != nil {
			result.Truncated = true
			continue
		}
		if _, parseErr := formats.ParseUnits(format, content); parseErr != nil {
			continue
		}
		candidate := Candidate{ID: path, Source: path, Target: target, Format: format.Name(), ConfiguredBundles: []string{}, RequiresConfirmation: temporaryPath(path)}
		candidate.Evidence = []Evidence{{Path: path, Kind: "catalog_path", Detail: "Source-locale filename or directory; storage format does not establish message grammar."}}
		applyRuntime(&candidate, nearestRuntime(filepath.Dir(path), runtimes))
		result.Candidates = append(result.Candidates, candidate)
	}
	for i, bundle := range cfg.EffectiveBundles() {
		resolved := inspectBundle(root, bundle, cfg, raw, i, runtimes, result)
		result.Bundles = append(result.Bundles, resolved)
		source := relativePath(root, resolved.Source)
		found := false
		for j := range result.Candidates {
			if result.Candidates[j].Source == source {
				result.Candidates[j].ConfiguredBundles = append(result.Candidates[j].ConfiguredBundles, bundle.ID)
				found = true
			}
		}
		if !found {
			candidate := Candidate{ID: source, Source: source, Target: relativePath(root, resolved.Target), Format: resolved.Format, ConfiguredBundles: []string{bundle.ID}, RequiresConfirmation: temporaryPath(source), Evidence: []Evidence{{Path: relativePath(root, result.ConfigPath), Kind: "configuration", Detail: "Existing configured source."}}}
			applyRuntime(&candidate, nearestRuntime(filepath.Dir(source), runtimes))
			result.Candidates = append(result.Candidates, candidate)
		}
	}
	for _, candidate := range result.Candidates {
		if len(candidate.ConfiguredBundles) == 0 {
			result.Diagnostics = append(result.Diagnostics, Diagnostic{Code: "UNCOVERED_CATALOG", Severity: "warning", Message: "Detected catalog is not covered by the current configuration: " + candidate.Source, Evidence: candidate.Evidence, RequiredDecisions: []string{"select_authoritative_source", "select_message_syntax"}, Recovery: []Recovery{{Argv: []string{"internationalizer", "config", "plan", "--help"}, SideEffects: []string{}}}})
		}
	}
	for _, locale := range cfg.TargetLocales {
		provider := cfg.LLMForLocale(locale)
		result.Credentials = append(result.Credentials, Credential{Locale: locale, Provider: provider.Provider, EnvironmentVariable: provider.APIKeyEnv, Present: os.Getenv(provider.APIKeyEnv) != "", ProviderVerified: false})
	}
	sort.Slice(result.Candidates, func(i, j int) bool { return result.Candidates[i].ID < result.Candidates[j].ID })
	if result.Truncated {
		result.Diagnostics = append(result.Diagnostics, Diagnostic{Code: "DISCOVERY_TRUNCATED", Severity: "warning", Message: "Discovery reached its file, depth, or size limit; narrow the project root. Unobserved integrations may exist."})
	}
	for i := range result.Diagnostics {
		diagnostic := &result.Diagnostics[i]
		if len(diagnostic.RequiredDecisions) > 0 {
			diagnostic.Recovery = []Recovery{{Argv: []string{"internationalizer", "config", "plan", "--help"}, SideEffects: []string{}, RequiredDecisions: append([]string{}, diagnostic.RequiredDecisions...)}}
		}
	}
	return result, nil
}

// Support documents inform translation policy; locale-shaped filenames alone
// do not make them translation catalogs. Explicit bundles are added separately.
func supportCatalogPath(root, path string, cfg config.Config) bool {
	for _, dir := range []string{cfg.StyleGuidesDir, cfg.GlossaryDir} {
		if dir == "" {
			continue
		}
		relative, err := filepath.Rel(absolutePath(root, dir), absolutePath(root, path))
		if err == nil && relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
			return true
		}
	}
	return false
}

func resolveConfigPath(root, explicit string) string {
	if explicit != "" {
		return absolutePath(root, explicit)
	}
	paths := []string{filepath.Join(root, ".internationalizer.yml"), filepath.Join(root, ".internationalizer.yaml")}
	if home, err := os.UserHomeDir(); err == nil {
		paths = append(paths, filepath.Join(home, ".internationalizer.yml"), filepath.Join(home, ".internationalizer.yaml"))
	}
	for _, path := range paths {
		if _, err := os.Lstat(path); err == nil {
			return path
		}
	}
	return paths[0]
}

func absolutePath(root, path string) string {
	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}
	return filepath.Join(root, path)
}

func relativePath(root, path string) string {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return filepath.ToSlash(path)
	}
	return filepath.ToSlash(rel)
}

func safeRead(path string) ([]byte, error) {
	name := strings.ToLower(filepath.Base(path))
	extension := strings.ToLower(filepath.Ext(path))
	if strings.HasPrefix(name, ".env") || extension == ".pem" || extension == ".p12" || extension == ".key" {
		return nil, fmt.Errorf("secret-shaped files are not inspected: %s", path)
	}
	// Check every component: a regular leaf beneath a symlink is still a symlink
	// traversal, and discovery must not escape through linked directories.
	for current := filepath.Clean(path); ; current = filepath.Dir(current) {
		info, err := os.Lstat(current)
		if err != nil {
			return nil, err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return nil, fmt.Errorf("symlink paths are not inspected: %s", path)
		}
		if current == filepath.Dir(current) {
			break
		}
	}
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("not a regular file: %s", path)
	}
	if info.Size() > maxDiscoveryBytes {
		return nil, fmt.Errorf("file exceeds discovery size limit: %s", path)
	}
	return os.ReadFile(path)
}

func skippedDirectory(name string) bool {
	switch name {
	case "node_modules", "build", "dist", "vendor", "coverage", "__pycache__", "data":
		return true
	}
	return strings.HasPrefix(name, ".")
}

func discoverFiles(root string, result *Inspection) ([]string, map[string]runtimeEvidence) {
	var files []string
	runtimes := map[string]runtimeEvidence{}
	var imports []Evidence
	count := 0
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			result.Diagnostics = append(result.Diagnostics, Diagnostic{Code: "DISCOVERY_PATH_UNREADABLE", Severity: "warning", Message: "A project path could not be inspected.", Evidence: []Evidence{{Path: relativePath(root, path), Kind: "unreadable", Detail: "Check filesystem access."}}})
			return nil
		}
		if path == root {
			return nil
		}
		count++
		if count > 50000 {
			result.Truncated = true
			return fs.SkipAll
		}
		rel := relativePath(root, path)
		if entry.Type()&os.ModeSymlink != 0 {
			return nil
		}
		if entry.IsDir() {
			if skippedDirectory(entry.Name()) {
				return fs.SkipDir
			}
			if strings.Count(rel, "/") >= 14 {
				result.Truncated = true
				return fs.SkipDir
			}
			return nil
		}
		if strings.HasPrefix(entry.Name(), ".") {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		switch ext {
		case ".json", ".yaml", ".yml", ".ftl", ".md", ".mdx":
			files = append(files, rel)
		}
		if entry.Name() == "package.json" {
			content, readErr := safeRead(path)
			if readErr != nil {
				result.Truncated = true
				return nil
			}
			var pkg struct {
				Dependencies     map[string]string `json:"dependencies"`
				DevDependencies  map[string]string `json:"devDependencies"`
				PeerDependencies map[string]string `json:"peerDependencies"`
			}
			if json.Unmarshal(content, &pkg) != nil {
				return nil
			}
			runtime := runtimeEvidence{evidence: []Evidence{{Path: rel, Kind: "package_manifest", Detail: "Nearest package dependency declarations inspected; runtime registration may be dynamic."}}}
			for _, deps := range []map[string]string{pkg.Dependencies, pkg.DevDependencies, pkg.PeerDependencies} {
				for _, name := range []string{"i18next", "react-i18next", "i18next-icu", "next-intl", "vue-i18n"} {
					if _, ok := deps[name]; ok {
						switch name {
						case "i18next-icu":
							runtime.icu = true
						case "next-intl":
							runtime.nextIntl = true
						case "vue-i18n":
							runtime.vueI18n = true
						default:
							runtime.i18next = true
						}
						runtime.evidence = append(runtime.evidence, Evidence{Path: rel, Kind: "dependency", Detail: name})
					}
				}
			}
			runtimes[filepath.ToSlash(filepath.Dir(rel))] = runtime
		}
		// Keep source inspection bounded and targeted to localization modules.
		localizationPath := strings.Contains(strings.ToLower(rel), "i18n") || strings.Contains(strings.ToLower(rel), "locale")
		if localizationPath && (ext == ".ts" || ext == ".tsx" || ext == ".js" || ext == ".jsx" || ext == ".mjs") {
			info, infoErr := entry.Info()
			if infoErr != nil || info.Size() > 256<<10 {
				result.Truncated = true
				return nil
			}
			content, readErr := safeRead(path)
			if readErr == nil && strings.Contains(string(content), "i18next-icu") {
				imports = append(imports, Evidence{Path: rel, Kind: "runtime_reference", Detail: "i18next-icu reference; verify plugin registration for this catalog."})
			}
		}
		return nil
	})
	if err != nil {
		result.Truncated = true
	}
	for _, evidence := range imports {
		dir := nearestRuntimeDirectory(filepath.Dir(evidence.Path), runtimes)
		if dir != "" {
			runtime := runtimes[dir]
			runtime.icu = true
			runtime.evidence = append(runtime.evidence, evidence)
			runtimes[dir] = runtime
		}
	}
	return files, runtimes
}

func nearestRuntimeDirectory(dir string, runtimes map[string]runtimeEvidence) string {
	if filepath.IsAbs(dir) || dir == ".." || strings.HasPrefix(filepath.ToSlash(dir), "../") {
		return ""
	}
	for {
		dir = filepath.ToSlash(dir)
		if _, ok := runtimes[dir]; ok {
			return dir
		}
		parent := filepath.ToSlash(filepath.Dir(dir))
		if parent == dir || dir == "." {
			return ""
		}
		dir = parent
	}
}

func nearestRuntime(dir string, runtimes map[string]runtimeEvidence) runtimeEvidence {
	return runtimes[nearestRuntimeDirectory(dir, runtimes)]
}

func applyRuntime(candidate *Candidate, runtime runtimeEvidence) {
	candidate.Framework = "unknown"
	candidate.SuggestedSyntax = ""
	candidate.Uncertainty = "high: storage and path evidence do not identify the runtime message grammar"
	candidate.Evidence = append(candidate.Evidence, runtime.evidence...)
	var frameworks []string
	if runtime.i18next || runtime.icu {
		frameworks = append(frameworks, "i18next")
	}
	if runtime.nextIntl {
		frameworks = append(frameworks, "next-intl")
	}
	if runtime.vueI18n {
		frameworks = append(frameworks, "vue-i18n")
	}
	if len(frameworks) > 1 {
		candidate.Framework = "multiple"
		candidate.Uncertainty = "high: multiple frameworks detected (" + strings.Join(frameworks, ", ") + "); explicitly choose this catalog's runtime and a compatible message syntax profile"
		return
	}
	if runtime.nextIntl {
		candidate.Framework = "next-intl"
		candidate.SuggestedSyntax = message.ICU
		candidate.Uncertainty = "medium: next-intl dependency suggests ICU messages; static dependency evidence does not prove that this catalog is consumed by next-intl"
	}
	if runtime.vueI18n {
		candidate.Framework = "vue-i18n"
		candidate.Uncertainty = "high: vue-i18n has its own message grammar; explicitly confirm compatibility and select a supported syntax profile before translation"
	}
	if runtime.i18next || runtime.icu {
		candidate.Framework = "i18next"
		candidate.SuggestedSyntax = message.I18next
		candidate.Uncertainty = "medium: i18next dependency detected; static inspection cannot prove that ICU integration is absent"
		if runtime.icu {
			candidate.SuggestedSyntax = ""
			candidate.Uncertainty = "high: i18next-icu evidence found; select i18next or icu after confirming runtime plugin registration"
		}
	}
}

func catalogTarget(path, locale string) (string, bool) {
	ext := filepath.Ext(path)
	if strings.TrimSuffix(filepath.Base(path), ext) == locale {
		return filepath.ToSlash(filepath.Join(filepath.Dir(path), "{locale}"+ext)), true
	}
	parts := strings.Split(filepath.ToSlash(path), "/")
	for i := len(parts) - 2; i >= 0; i-- {
		if parts[i] == locale {
			parts[i] = "{locale}"
			return strings.Join(parts, "/"), true
		}
	}
	return "", false
}

func temporaryPath(path string) bool {
	for _, part := range strings.Split(filepath.ToSlash(path), "/") {
		if part == "tmp" || part == "temp" {
			return true
		}
	}
	return false
}

func inspectBundle(root string, bundle config.Bundle, cfg, raw config.Config, index int, runtimes map[string]runtimeEvidence, result *Inspection) ResolvedBundle {
	source, target := absolutePath(root, bundle.Source), absolutePath(root, bundle.Target)
	format, err := formats.FormatForFile(source)
	if bundle.Format != "" {
		format, err = formats.FormatByName(bundle.Format)
	}
	formatName := bundle.Format
	if err == nil {
		formatName = format.Name()
	}
	provenance := map[string]string{"source": "bundle.source", "target": "bundle.target", "locales": "target_locales", "message_syntax": "default:auto", "format": "source_extension", "path_base": "working_directory"}
	if len(raw.Bundles) == 0 {
		provenance["source"] = "source_path"
		provenance["target"] = "source_path sibling convention"
	}
	if raw.MessageSyntax != "" {
		provenance["message_syntax"] = "message_syntax"
	}
	if index < len(raw.Bundles) && raw.Bundles[index].MessageSyntax != "" {
		provenance["message_syntax"] = "bundle.message_syntax"
	}
	if bundle.Format != "" {
		provenance["format"] = "bundle.format"
	}
	candidate := Candidate{}
	applyRuntime(&candidate, nearestRuntime(filepath.Dir(relativePath(root, source)), runtimes))
	resolved := ResolvedBundle{ID: bundle.ID, Source: source, Target: target, Format: formatName, Framework: candidate.Framework, MessageSyntax: bundle.MessageSyntax, Locales: append([]string{}, cfg.TargetLocales...), Targets: map[string]string{}, Provenance: provenance, Evidence: candidate.Evidence}
	for _, locale := range cfg.TargetLocales {
		if path, pathErr := bundle.TargetPath(locale); pathErr == nil {
			resolved.Targets[locale] = absolutePath(root, path)
		}
	}
	if temporaryPath(relativePath(root, source)) {
		result.Diagnostics = append(result.Diagnostics, Diagnostic{Code: "SOURCE_CONFIRMATION_REQUIRED", Severity: "warning", Bundle: bundle.ID, Message: "Configured source is under tmp/ or temp/; confirm which pipeline owns this artifact before changing the configuration.", RequiredDecisions: []string{"confirm_authoritative_source"}, Evidence: []Evidence{{Path: relativePath(root, source), Kind: "temporary_path", Detail: "Temporary location requires confirmation; it is not automatically incorrect."}}})
	}
	if err != nil {
		result.Diagnostics = append(result.Diagnostics, Diagnostic{Code: "SOURCE_FORMAT_INVALID", Severity: "error", Bundle: bundle.ID, Message: err.Error()})
		return resolved
	}
	data, err := safeRead(source)
	if err != nil {
		result.Diagnostics = append(result.Diagnostics, Diagnostic{Code: "SOURCE_UNREADABLE", Severity: "error", Bundle: bundle.ID, Message: err.Error()})
		return resolved
	}
	units, err := formats.ParseSourceUnits(format, data, bundle.MessageSyntax)
	if err != nil {
		result.Diagnostics = append(result.Diagnostics, Diagnostic{Code: "SOURCE_PARSE_FAILED", Severity: "error", Bundle: bundle.ID, Message: "Configured source could not be parsed using its storage format."})
		return resolved
	}
	for _, unit := range units {
		findings := validate.SyntaxSourceFindings(unit.ID, unit.Value, cfg.SourceLocale, unit.Syntax)
		codeBraces := false
		for _, span := range protectedtext.HTMLCode(unit.Value) {
			if strings.Contains(span, "{") {
				codeBraces = true
				break
			}
		}
		auto := bundle.MessageSyntax == message.Auto || bundle.MessageSyntax == ""
		if len(findings) == 0 && (!auto || unit.Syntax != message.ICU || !codeBraces) {
			continue
		}
		code, text := "SOURCE_SYNTAX_INVALID", "Source does not satisfy its explicit message syntax."
		severity := "error"
		if len(findings) == 0 {
			severity = "warning"
		}
		var decisions []string
		if auto {
			code = "AUTO_SYNTAX_AMBIGUOUS"
			text = "With message_syntax: auto, brace syntax was interpreted as ICU. Select the bundle's runtime syntax: plain, i18next, or icu."
			if codeBraces {
				text = "Brace syntax inside HTML code was interpreted as ICU with message_syntax: auto. Select the bundle's runtime syntax: plain, i18next, or icu."
			}
			decisions = []string{"select_message_syntax"}
		}
		result.Diagnostics = append(result.Diagnostics, Diagnostic{Code: code, Severity: severity, Bundle: bundle.ID, Key: unit.ID, Message: text, RequiredDecisions: decisions, Evidence: []Evidence{{Path: relativePath(root, source), Kind: "source_syntax", Detail: "Source-only check; not repeated for each target locale."}}})
	}
	return resolved
}
