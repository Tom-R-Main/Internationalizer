package detect

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// ProjectType identifies the i18n framework in use.
type ProjectType string

const (
	ReactI18Next ProjectType = "react-i18next"
	NextIntl     ProjectType = "next-intl"
	VueI18n      ProjectType = "vue-i18n"
	VanillaJSON  ProjectType = "vanilla-json"
	MarkdownDocs ProjectType = "markdown-docs"
	Unknown      ProjectType = "unknown"
)

// Detection holds the result of project type detection.
type Detection struct {
	Type           ProjectType `json:"type"`
	Confidence     float64     `json:"confidence"`
	SourcePath     string      `json:"source_path"`
	SourceLocale   string      `json:"source_locale"`
	TargetLocales  []string    `json:"target_locales,omitempty"`
	SuggestedPaths []string    `json:"suggested_paths,omitempty"`
}

// Detect scans the given directory to identify the i18n project type.
func Detect(dir string) Detection {
	// Check for package.json-based detection first.
	if d := detectFromPackageJSON(dir); d.Type != Unknown {
		return d
	}

	// Check for locale directory structures.
	if d := detectLocaleDirectories(dir); d.Type != Unknown {
		return d
	}

	return Detection{Type: Unknown, Confidence: 0}
}

func detectFromPackageJSON(dir string) Detection {
	data, err := os.ReadFile(filepath.Join(dir, "package.json"))
	if err != nil {
		return Detection{Type: Unknown}
	}

	var pkg struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return Detection{Type: Unknown}
	}

	allDeps := make(map[string]string)
	for k, v := range pkg.Dependencies {
		allDeps[k] = v
	}
	for k, v := range pkg.DevDependencies {
		allDeps[k] = v
	}

	if _, ok := allDeps["react-i18next"]; ok {
		return detectReactI18Next(dir)
	}
	if _, ok := allDeps["next-intl"]; ok {
		return detectNextIntl(dir)
	}
	if _, ok := allDeps["vue-i18n"]; ok {
		return detectVueI18n(dir)
	}

	return Detection{Type: Unknown}
}

func detectReactI18Next(dir string) Detection {
	d := Detection{
		Type:       ReactI18Next,
		Confidence: 0.9,
	}

	// Common locale paths for react-i18next.
	candidates := []string{
		"public/locales",
		"src/locales",
		"src/i18n/locales",
		"locales",
		"src/i18n",
	}
	for _, c := range candidates {
		p := filepath.Join(dir, c)
		if info, err := os.Stat(p); err == nil && info.IsDir() {
			locales := findLocalesInDir(p)
			if len(locales) > 0 {
				d.SourcePath = filepath.Join(c, "en.json")
				d.SourceLocale = "en"
				d.TargetLocales = filterLocales(locales, "en")
				d.SuggestedPaths = []string{c}
				return d
			}
		}
	}

	// Check for flat locale files (en.json in src/i18n/).
	for _, c := range candidates {
		enPath := filepath.Join(dir, c, "en.json")
		if _, err := os.Stat(enPath); err == nil {
			d.SourcePath = filepath.Join(c, "en.json")
			d.SourceLocale = "en"
			d.SuggestedPaths = []string{c}
			return d
		}
	}

	return d
}

func detectNextIntl(dir string) Detection {
	d := Detection{
		Type:       NextIntl,
		Confidence: 0.9,
	}
	candidates := []string{"messages", "locales", "src/messages"}
	for _, c := range candidates {
		enPath := filepath.Join(dir, c, "en.json")
		if _, err := os.Stat(enPath); err == nil {
			d.SourcePath = filepath.Join(c, "en.json")
			d.SourceLocale = "en"
			locales := findJSONLocalesInDir(filepath.Join(dir, c))
			d.TargetLocales = filterLocales(locales, "en")
			return d
		}
	}
	return d
}

func detectVueI18n(dir string) Detection {
	d := Detection{
		Type:       VueI18n,
		Confidence: 0.9,
	}
	candidates := []string{"src/locales", "locales", "src/i18n"}
	for _, c := range candidates {
		enPath := filepath.Join(dir, c, "en.json")
		if _, err := os.Stat(enPath); err == nil {
			d.SourcePath = filepath.Join(c, "en.json")
			d.SourceLocale = "en"
			locales := findJSONLocalesInDir(filepath.Join(dir, c))
			d.TargetLocales = filterLocales(locales, "en")
			return d
		}
	}
	return d
}

func detectLocaleDirectories(dir string) Detection {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return Detection{Type: Unknown}
	}

	// Look for directories named by locale codes containing translatable files.
	var localeDirs []string
	var mdDirs []string
	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		name := entry.Name()
		if isLocaleCode(name) {
			subPath := filepath.Join(dir, name)
			if hasJSONFiles(subPath) {
				localeDirs = append(localeDirs, name)
			}
			if hasMDFiles(subPath) {
				mdDirs = append(mdDirs, name)
			}
		}
	}

	if len(localeDirs) > 0 {
		d := Detection{
			Type:         VanillaJSON,
			Confidence:   0.7,
			SourceLocale: "en",
		}
		if contains(localeDirs, "en") {
			d.SourcePath = "en/"
			d.TargetLocales = filterLocales(localeDirs, "en")
		} else {
			d.SourcePath = localeDirs[0] + "/"
			d.TargetLocales = filterLocales(localeDirs, localeDirs[0])
		}
		return d
	}

	if len(mdDirs) > 0 {
		d := Detection{
			Type:         MarkdownDocs,
			Confidence:   0.7,
			SourceLocale: "en",
		}
		if contains(mdDirs, "en") {
			d.SourcePath = "en/"
			d.TargetLocales = filterLocales(mdDirs, "en")
		}
		return d
	}

	return Detection{Type: Unknown}
}

var commonLocales = map[string]bool{
	"en": true, "es": true, "fr": true, "de": true, "it": true,
	"pt": true, "ja": true, "ko": true, "zh": true, "ar": true,
	"ru": true, "hi": true, "nl": true, "sv": true, "pl": true,
	"da": true, "fi": true, "nb": true, "tr": true, "th": true,
	"id": true, "vi": true, "uk": true, "cs": true, "el": true,
	"he": true, "ro": true, "hu": true, "bn": true, "pa": true,
	"te": true, "pt-BR": true, "zh-CN": true, "zh-TW": true,
	"en-US": true, "en-GB": true, "es-MX": true, "fr-CA": true,
	"yue": true,
}

func isLocaleCode(name string) bool {
	return commonLocales[name]
}

func findLocalesInDir(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var locales []string
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() && isLocaleCode(name) {
			locales = append(locales, name)
		} else if !e.IsDir() && strings.HasSuffix(name, ".json") {
			locale := strings.TrimSuffix(name, ".json")
			if isLocaleCode(locale) {
				locales = append(locales, locale)
			}
		}
	}
	return locales
}

func findJSONLocalesInDir(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var locales []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".json") {
			locale := strings.TrimSuffix(e.Name(), ".json")
			if isLocaleCode(locale) {
				locales = append(locales, locale)
			}
		}
	}
	return locales
}

func hasJSONFiles(dir string) bool {
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".json") {
			return true
		}
	}
	return false
}

func hasMDFiles(dir string) bool {
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".md") || strings.HasSuffix(e.Name(), ".mdx") {
			return true
		}
	}
	return false
}

func filterLocales(locales []string, exclude string) []string {
	var result []string
	for _, l := range locales {
		if l != exclude {
			result = append(result, l)
		}
	}
	return result
}

func contains(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}
