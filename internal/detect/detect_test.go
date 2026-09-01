package detect

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectProjectFixtures(t *testing.T) {
	tests := []struct {
		name           string
		files          map[string]string
		wantType       ProjectType
		wantSource     string
		wantTargets    []string
		wantConfidence float64
	}{
		{
			name: "react i18next",
			files: map[string]string{
				"package.json":                  `{"dependencies":{"react-i18next":"1"}}`,
				"public/locales/en.json":        `{}`,
				"public/locales/fr.json":        `{}`,
				"public/locales/not-locale.txt": "ignored",
			},
			wantType: ReactI18Next, wantSource: filepath.Join("public", "locales", "en.json"), wantTargets: []string{"fr"}, wantConfidence: 0.9,
		},
		{
			name: "next intl",
			files: map[string]string{
				"package.json":     `{"devDependencies":{"next-intl":"1"}}`,
				"messages/de.json": `{}`,
				"messages/en.json": `{}`,
			},
			wantType: NextIntl, wantSource: filepath.Join("messages", "en.json"), wantTargets: []string{"de"}, wantConfidence: 0.9,
		},
		{
			name: "vue i18n",
			files: map[string]string{
				"package.json":        `{"dependencies":{"vue-i18n":"1"}}`,
				"src/locales/en.json": `{}`,
				"src/locales/ja.json": `{}`,
			},
			wantType: VueI18n, wantSource: filepath.Join("src", "locales", "en.json"), wantTargets: []string{"ja"}, wantConfidence: 0.9,
		},
		{
			name: "vanilla locale directories",
			files: map[string]string{
				"en/common.json": `{}`,
				"es/common.json": `{}`,
			},
			wantType: VanillaJSON, wantSource: "en/", wantTargets: []string{"es"}, wantConfidence: 0.7,
		},
		{
			name: "markdown locale directories",
			files: map[string]string{
				"en/readme.md": "Hello",
				"fr/readme.md": "Bonjour",
			},
			wantType: MarkdownDocs, wantSource: "en/", wantTargets: []string{"fr"}, wantConfidence: 0.7,
		},
		{
			name: "malformed package",
			files: map[string]string{
				"package.json": `{`,
			},
			wantType: Unknown,
		},
		{
			name:     "empty project",
			files:    map[string]string{},
			wantType: Unknown,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dir := t.TempDir()
			for name, contents := range test.files {
				path := filepath.Join(dir, name)
				if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
					t.Fatal(err)
				}
			}

			got := Detect(dir)
			if got.Type != test.wantType || got.SourcePath != test.wantSource || got.Confidence != test.wantConfidence {
				t.Fatalf("Detect() = %#v, want type=%q source=%q confidence=%v", got, test.wantType, test.wantSource, test.wantConfidence)
			}
			if !equalStrings(got.TargetLocales, test.wantTargets) {
				t.Fatalf("target locales = %#v, want %#v", got.TargetLocales, test.wantTargets)
			}
		})
	}
}

func equalStrings(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}
