> [English (original)](../../README.md)

<p align="center">
  <img src="../../assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

Natywny dla AI potok internacjonalizacji dla projektów programistycznych. Tłumacz, weryfikuj i zarządzaj plikami i18n przy użyciu LLM.

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](../../LICENSE)

<p align="center">
<a href="ar.md">العربية</a> · <a href="bn.md">বাংলা</a> · <a href="cs.md">Čeština</a> · <a href="da.md">Dansk</a> · <a href="de.md">Deutsch</a> · <a href="el.md">Ελληνικά</a> · <a href="es.md">Español</a> · <a href="fi.md">Suomi</a> · <a href="fr.md">Français</a> · <a href="he.md">עברית</a> · <a href="hi.md">हिन्दी</a> · <a href="id.md">Indonesia</a> · <a href="it.md">Italiano</a> · <a href="ja.md">日本語</a> · <a href="ko.md">한국어</a> · <a href="ms.md">Bahasa Melayu</a><br><a href="nl.md">Nederlands</a> · <a href="pa.md">ਪੰਜਾਬੀ</a> · <a href="pl.md">Polski</a> · <a href="pt-BR.md">Português</a> · <a href="ro.md">Română</a> · <a href="ru.md">Русский</a> · <a href="sv.md">Svenska</a> · <a href="te.md">తెలుగు</a> · <a href="th.md">ไทย</a> · <a href="tr.md">Türkçe</a> · <a href="uk.md">Українська</a> · <a href="vi.md">Tiếng Việt</a> · <a href="yue.md">粵語</a> · <a href="zh-CN.md">简体中文</a> · <a href="zh-TW.md">繁體中文</a>
</p>

---

<!-- internationalizer:unit markdown:why-internationalizer -->
## Dlaczego Internationalizer?

Większość narzędzi i18n to biblioteki uruchomieniowe (i18next, react-intl) albo platformy SaaS do zarządzania kluczami (Crowdin, Lokalise). Żadne z nich nie rozwiązuje dobrze rzeczywistego problemu tłumaczenia:

- **Ręczne tłumaczenie** nie skaluje się powyżej kilku języków
- **Interfejsy API tłumaczenia maszynowego** (Google Translate, DeepL) ignorują Twoją terminologię, ton i konwencje interfejsu użytkownika
- **Ogólne tłumaczenie za pomocą LLM** działa lepiej, ale bez glosariuszy i przewodników po stylu daje niespójne wyniki

Internationalizer działa inaczej. To **potok CLI**, który łączy tłumaczenie LLM z:

- **Glosariuszami dla poszczególnych języków** — wymuszają spójną terminologię w całej aplikacji
- **Przewodnikami po stylu dla poszczególnych języków** — kontrolują ton, formalność, formy liczby mnogiej i typografię
- **Pamięcią tłumaczeniową** — pomija niezmienione ciągi znaków, obniżając koszty wywołań API
- **Deterministyczną walidacją** — wychwytuje brakujące lub nadmiarowe klucze, dryf chronionej struktury, niezgodności z glosariuszem oraz błędy form liczby mnogiej lub formatu ICU przed wdrożeniem

<!-- internationalizer:unit markdown:installation -->
## Instalacja

Zainstaluj z npm:

```bash
npm install -g internationalizer
```

Lub uruchom bez instalacji globalnej:

```bash
npx internationalizer --help
```

Pakiet npm instaluje odpowiedni, prekompilowany plik binarny z npm za pośrednictwem opcjonalnych zależności specyficznych dla danej platformy.

Zainstaluj za pomocą Go:

```bash
go install github.com/Tom-R-Main/Internationalizer/cmd/internationalizer@latest
```

Lub skompiluj ze źródeł:

```bash
git clone https://github.com/Tom-R-Main/Internationalizer.git
cd Internationalizer
go build -o internationalizer ./cmd/internationalizer
```

<!-- internationalizer:unit markdown:npm-packages -->
## Pakiety npm

- Tagi Git i wersje pakietów npm muszą być zgodne, na przykład `v0.1.0` i `0.1.0`
- Główny pakiet `internationalizer` zależy od pakietów platformowych, takich jak `internationalizer-darwin-arm64`
- Obsługiwane platformy docelowe npm: macOS arm64/x64, Linux arm64/x64, Windows x64
- Publikowanie w ramach CI wymaga sekretu GitHub o nazwie `NPM_TOKEN`

<!-- internationalizer:unit markdown:quick-start -->
## Szybki start

1. Utwórz plik konfiguracyjny w katalogu głównym projektu:

```yaml
# .internationalizer.yml
source_locale: en
target_locales: [fr, de, es, ja]
bundles:
  - id: app
    source: locales/en.json
    target: locales/{locale}.json
    format: json

llm:
  provider: gemini
  model: gemini-3.8-flash
  api_key_env: GOOGLE_AI_STUDIO_API_KEY
```

2. Ustaw swój klucz API:

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

3. Wyświetl podgląd tego, co zostanie przetłumaczone:

```bash
internationalizer translate --dry-run
```

4. Uruchom tłumaczenie:

```bash
internationalizer translate
```

5. Zweryfikuj wszystkie wersje językowe:

```bash
internationalizer validate
```

<!-- internationalizer:unit markdown:commands -->
## Polecenia

### `translate`

Znajdź brakujące lub nieaktualne klucze i przetłumacz je za pomocą LLM.

```bash
internationalizer translate                    # translate all locales
internationalizer translate -l fr              # translate French only
internationalizer translate --dry-run          # preview without API calls
internationalizer translate --adopt-existing   # baseline existing translations without API calls
internationalizer translate --refresh-policy   # refresh prompt/style/model-stale entries
internationalizer translate --batch-size 20    # smaller batches
internationalizer translate --concurrency 2    # fewer parallel calls
```

Stan tłumaczenia niezależnie raportuje stany: brakujące, nieaktualne względem źródła (source-stale), nieaktualne względem zasad (policy-stale), aktualne oraz zmodyfikowane ręcznie, dzięki czemu ręczna edycja nie może zamaskować zmiany w źródle lub zasadach. Wartości nieaktualne względem zasad są raportowane, lecz zostają przetłumaczone ponownie wyłącznie przy użyciu opcji `--refresh-policy`. Ręcznie zmodyfikowane wartości nigdy nie są nadpisywane automatycznie. Użyj opcji `--adopt-existing` przy wprowadzaniu manifestu do przejrzanych tłumaczeń lub przy wyraźnym zaakceptowaniu przejrzanej edycji ręcznej jako nowego punktu bazowego.

### `validate`

Sprawdź wszystkie pliki językowe względem ich pakietów źródłowych. Domyślna walidacja sprawdza pokrycie strukturalne (procent obecnych wymaganych kluczy docelowych), raportuje nadmiarowe klucze jako ostrzeżenia i kończy się niepowodzeniem w przypadku brakujących kluczy, niezgodności interpolacji lub nieprawidłowej struktury ICU MessageFormat.

```bash
internationalizer validate                     # human-readable output
internationalizer validate --json              # machine-readable JSON
internationalizer validate -q                  # exit code only
internationalizer validate --strict             # enforce translation quality rules
internationalizer validate --require-state      # require current manifest provenance
```

Opcja `--strict` raportuje również pokrycie przetłumaczenia. Wartość językowa identyczna ze źródłem jest uznawana za nieprzetłumaczoną, chyba że glosariusz zawiera jednoznaczny wpis o identycznym źródle i celu dla całej wartości; parametr `ignore_case` jest uwzględniany, ale termin z glosariusza osadzony w dłuższej wartości nie stanowi wyjątku. Tryb ścisły kończy się niepowodzeniem w przypadku nadmiarowych kluczy, wartości identycznych ze źródłem, zmienionej struktury interpolacji/HTML/kodu/łączy Markdown, naruszeń glosariusza oraz brakujących skonfigurowanych form liczby mnogiej.

Opcja `--require-state` weryfikuje każdy plik docelowy względem `.internationalizer.lock`. Zgłasza błąd, gdy klucz nie jest śledzony lub gdy jego zapisany skrót źródła, zasad tłumaczenia bądź wartości docelowej jest nieaktualny. Opcję tę można łączyć z `--strict`.

Raporty czytelne dla człowieka oraz JSON korzystają ze stabilnych kodów wyników:

| Kod | Znaczenie |
| --- | --- |
| `missing_key` / `extra_key` | Zestawy kluczy źródłowych i docelowych różnią się |
| `blank_translation` | Niepusta wartość źródłowa ma pustą wartość docelową w trybie ścisłym |
| `source_identical` | Wartość językowa w trybie ścisłym pozostaje nieprzetłumaczona |
| `protected_structure_mismatch` | Zmiana struktury interpolacji, HTML, kodu lub łącza |
| `glossary_violation` | Nie znaleziono zatwierdzonego terminu docelowego ani jego wariantu |
| `plural_form_missing` | Brak skonfigurowanej formy liczby mnogiej dla danego języka |
| `icu_message_syntax` | Niepoprawna składnia komunikatu ICU w źródle lub pliku docelowym |
| `icu_argument_mismatch` | Różnice w nazwach, typach lub stylach formatowania argumentów ICU |
| `icu_selector_mismatch` | Selektory różnią się lub kategoria liczby mnogiej jest nieprawidłowa dla docelowego języka |
| `untracked` | Brak wpisu w manifeście dla elementu docelowego |
| `source_stale` | Zawartość źródłowa zmieniła się po zarejestrowaniu tłumaczenia |
| `policy_stale` | Wygenerowany prompt lub ustawienia modelu uległy zmianie |
| `target_modified` | Zawartość docelowa różni się od wpisu w manifeście |

### `detect`

Automatycznie wykryj framework i18n i zasugeruj konfigurację.

```bash
internationalizer detect
```

Obsługuje: react-i18next, next-intl, vue-i18n, czysty JSON, dokumenty Markdown.

### `glossary`

Zarządzaj terminami glosariusza dla poszczególnych języków, które są egzekwowane podczas tłumaczenia.

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

Zarządzaj pamięcią tłumaczeniową (pamięć podręczna JSONL wcześniej przetłumaczonych ciągów znaków).

```bash
internationalizer tm stats                     # show record counts
internationalizer tm export                    # dump as JSON
internationalizer tm clear --force             # delete all records
```

<!-- internationalizer:unit markdown:configuration-reference -->
## Dokumentacja konfiguracji

```yaml
# .internationalizer.yml

# Source language (default: en)
source_locale: en

# Languages to translate into (required)
target_locales: [fr, de, es, ja, yue, zh-CN, zh-TW, ar]

# One or more source-to-target mappings (required).
# {locale} is replaced with each configured target locale.
bundles:
  - id: app
    source: locales/en.json
    target: locales/{locale}.json
    format: json
  - id: docs
    source: README.md
    target: docs/i18n/{locale}.md
    format: markdown

# Backward compatibility: source_path still maps targets to sibling files
# such as locales/fr.json. Prefer bundles for new projects.
# source_path: locales/en.json

# LLM provider settings
llm:
  # Provider: "anthropic", "openai", "gemini", or "openrouter" (default: gemini)
  provider: gemini

  # Model name defaults by provider:
  #   anthropic:  claude-opus-5
  #   openai:     gpt-5.6-luna (reasoning effort defaults to max)
  #   gemini:     gemini-3.8-flash
  #   openrouter: deepseek/deepseek-v4-pro-0813
  model: gemini-3.8-flash

  # Environment variable containing the API key
  api_key_env: GOOGLE_AI_STUDIO_API_KEY

  # Base URL for OpenAI-compatible endpoints (optional)
  # base_url: https://api.openai.com

  # OpenAI GPT-5-series Responses API reasoning effort
  # (default: max for the OpenAI provider)
  reasoning_effort: max

  # Optional LLM settings for individual target locales. An override using the
  # global provider inherits unspecified global settings. A different provider
  # uses that provider's defaults for unspecified settings.
  locale_overrides:
    yue:
      provider: openrouter
      model: deepseek/deepseek-v4-flash-0731
      api_key_env: OPENROUTER_API_KEY
    zh-CN:
      provider: openrouter
      model: deepseek/deepseek-v4-flash-0731
      api_key_env: OPENROUTER_API_KEY
    zh-TW:
      provider: openrouter
      model: deepseek/deepseek-v4-flash-0731
      api_key_env: OPENROUTER_API_KEY

# Keys per LLM call (default: 40)
batch_size: 40

# Parallel LLM calls (default: 4)
concurrency: 4

# Directory containing per-locale style guide Markdown files (default: style-guides)
style_guides_dir: style-guides

# Directory containing per-locale glossary JSON files (default: glossary)
glossary_dir: glossary

# Path to translation memory file (default: .internationalizer/tm.jsonl)
tm_path: .internationalizer/tm.jsonl

# Versioned source, policy, target, and provenance state
# (default: .internationalizer.lock; commit this file)
manifest_path: .internationalizer.lock

# Optional translation and strict-validation rules
validation:
  plural_style: i18next-v4 # generate and validate target-locale plural forms
```

Identyfikatory ustawień regionalnych muszą być prawidłowymi znacznikami BCP 47, takimi jak `fr`, `pt-BR` lub `sr-Latn-RS`. Równoważne kanonicznie docelowe ustawienia regionalne są odrzucane jako duplikaty, a specyficzne dla ustawień regionalnych nadpisania dostawców są dopasowywane według pisowni kanonicznej. W powyższym przykładzie języki docelowe bez nadpisania — w tym japoński — dziedziczą globalną konfigurację Gemini.

Wartości ICU MessageFormat są analizowane pod kątem struktury. Obsługiwane są proste argumenty, `select`, `plural`, `selectordinal`, `number`, `date` oraz `time`, w tym zagnieżdżone komunikaty, przesunięcia liczby mnogiej (plural offsets), selektory dokładnych liczb oraz `#`. Walidacja weryfikuje składnię, typy argumentów i style formatowania, przesunięcia liczby mnogiej, zgodność gałęzi instrukcji select oraz kategorie liczb mnogich CLDR dla docelowego języka. Wyniki od dostawcy, które naruszają te niezmienniki, są odrzucane przed zapisaniem pliku językowego lub rekordu w pamięci tłumaczeniowej.

W przypadku `i18next-v4` rozpoznane źródłowe rodziny liczb mnogich są rozwijane podczas tłumaczenia do kategorii CLDR docelowego języka. Kategoria występująca wyłącznie w języku docelowym używa wartości `_other` rodziny źródłowej jako szablonu tłumaczenia. Walidacja ścisła wymaga tych kategorii docelowych; kategorie występujące wyłącznie w źródle są opcjonalne dla języków docelowych, które ich nie używają.

<!-- internationalizer:unit markdown:style-guides -->
## Przewodniki po stylu

Przewodniki po stylu to pliki Markdown wstrzykiwane do promptu tłumaczeniowego LLM. Sterują one tonem, formalnością, typografią i innymi konwencjami specyficznymi dla danego języka.

```
style-guides/
  _conventions.md    # shared rules for all languages
  fr.md              # French-specific rules
  ja.md              # Japanese-specific rules
  ar.md              # Arabic-specific rules
```

### Wspólne konwencje (`_conventions.md`)

Definiują reguły mające zastosowanie do wszystkich języków: składnię interpolacji, zachowywanie znaczników HTML, konwencje typów ciągów znaków (przyciski a etykiety a błędy) itp.

### Przewodniki dla poszczególnych języków (`{locale}.md`)

Definiują reguły specyficzne dla danego języka: rejestr formalności (np. ty a Pan/Pani), interpunkcję (cudzysłowy francuskie, odwrócone znaki zapytania), formy liczby mnogiej, formatowanie dat i liczb oraz glosariusz terminologiczny.

Przewodniki po stylu stanowią trwałe dane wejściowe zasad, a nie generowane dane wyjściowe. Internationalizer odczytuje je, lecz nigdy ich nie modyfikuje. Ich zawartość jest haszowana niezależnie od glosariusza i kontraktu promptu, dzięki czemu zmiana kodu aplikacji nie powoduje dezaktualizacji tłumaczenia. Edycja przewodnika celowo oznacza dany język do przeglądu zasad; zmiana wewnętrznego sformułowania promptu tego nie powoduje, chyba że zmieni się również wersja kontraktu promptu.

Zobacz [`examples/react-app/style-guides/`](../../examples/react-app/style-guides/), aby zapoznać się z działającym przykładem.

<!-- internationalizer:unit markdown:glossary-format -->
## Format glosariusza

Pliki glosariusza to tablice JSON zapisane w `{glossary_dir}/{locale}.json`:

```json
[
  {
    "source": "Dashboard",
    "target": "Tableau de bord",
    "variants": ["Panneau de contrôle"],
    "enforcement": "error",
    "ignore_case": false,
    "whole_word": true
  }
]
```

Pole `variants` zawiera listę innych zatwierdzonych form docelowych. Pole `enforcement` może przyjmować wartość `error`, `warning` lub zostać pominięte (wtedy domyślnym zachowaniem jest błąd). Terminy są wstrzykiwane do promptu LLM jako tabela terminologiczna, co zapewnia spójne tłumaczenie w całej aplikacji. Dokładny wpis, taki jak `{"source":"API","target":"API"}`, zwalnia również całą wartość identyczną ze źródłem z ustaleń o nieprzetłumaczonych wartościach w trybie ścisłym; nie zwalnia to jednak dłuższej wartości, która jedynie zawiera `API`.

<!-- internationalizer:unit markdown:translation-memory -->
## Pamięć tłumaczeniowa

Pamięć tłumaczeniowa jest przechowywana jako plik JSONL (jeden rekord JSON na wiersz). Każdy rekord zawiera:

- Pakiet, klucz, wartość źródłową, przetłumaczoną wartość oraz kanoniczne docelowe ustawienia regionalne
- Skróty źródła, przewodnika po stylu, glosariusza, kontraktu promptu oraz połączonych zasad
- Nazwę dostawcy i model, który wygenerował tłumaczenie
- Znacznik czasu

Podczas kolejnych uruchomień ciągi znaków o identycznych skrótach źródła i zasad są pobierane z pamięci podręcznej bez wywoływania LLM. Domyślna ścieżka znajduje się w ignorowanym katalogu `.internationalizer/`, dzięki czemu pozostaje lokalną pamięcią podręczną. Ustaw `tm_path` na śledzoną lokalizację, jeśli projekt celowo współdzieli pamięć tłumaczeniową. Przeznaczony do przeglądu manifest `.internationalizer.lock` jest wersjonowany niezależnie.

<!-- internationalizer:unit markdown:supported-formats -->
## Obsługiwane formaty

| Format | Rozszerzenia | Tryb |
|--------|-----------|------|
| JSON | `.json` | Klucz-wartość (zagnieżdżone, spłaszczone notacją kropkową) |
| YAML | `.yml`, `.yaml` | Klucz-wartość (zachowuje komentarze i kolejność) |
| Markdown | `.md`, `.mdx` | Preambuła i sekcje poziomu H2 |

Docelowe pliki Markdown zawierają niewidoczne komentarze `internationalizer:unit` przed sekcjami H2. Te stabilne znaczniki pozwalają narzędziu Internationalizer dodawać, przenosić lub modyfikować pojedynczą sekcję źródłową bez ponownego tłumaczenia niepowiązanych sekcji. Istniejące dokumenty bez znaczników otrzymają je przy kolejnej udanej aktualizacji.

<!-- internationalizer:unit markdown:project-type-detection -->
## Wykrywanie typu projektu

Polecenie `internationalizer detect` identyfikuje Twoją konfigurację i18n, sprawdzając:

- Zależności w `package.json` pod kątem react-i18next, next-intl lub vue-i18n
- Struktury katalogów pasujące do popularnych wzorców ustawień regionalnych
- Rozszerzenia plików i konwencje nazewnictwa

<!-- internationalizer:unit markdown:architecture -->
## Architektura

```
cmd/internationalizer/     CLI entry point and command definitions
internal/
  config/                  YAML config loading with defaults
  detect/                  Project type auto-detection
  formats/                 Format parsers (JSON, YAML, Markdown)
  glossary/                Per-locale glossary management
  llm/                     LLM provider interface + implementations
    anthropic.go           Anthropic Claude backend
    openai.go              OpenAI / compatible backend
    gemini.go              Google Gemini via AI Studio backend
                           OpenRouter uses openai.go with custom base_url
  locale/                  BCP 47 identity and CLDR plural categories
  message/                 ICU MessageFormat parser and structural comparison
  policy/                  Stable translation-policy hashing
  state/                   Versioned translation manifest
  styleguide/              Style guide loader
  tm/                      JSONL translation memory
  translate/               Translation orchestrator
  validate/                Locale validation and diffing
```

<!-- internationalizer:unit markdown:comparison-to-alternatives -->
## Porównanie z alternatywami

| Funkcja | Internationalizer | i18next | Crowdin | Ogólny LLM |
|---------|------------------|---------|---------|-------------|
| Tłumaczenie oparte na LLM | Tak | Nie | Częściowo | Tak |
| Przewodniki po stylu dla języków | Tak | Nie | Nie | Nie |
| Wymuszanie glosariusza | Tak | Nie | Tak | Nie |
| Pamięć tłumaczeniowa | Tak | Nie | Tak | Nie |
| CLI / uruchamianie lokalne | Tak | N/D | Nie | Ręcznie |
| Pliki przyjazne dla Git | Tak | Tak | Częściowo | Ręcznie |
| Brak zależności od SaaS | Tak | Tak | Nie | Różnie |
| Otwarte oprogramowanie (AGPL-3.0) | Tak | Tak | Nie | Różnie |

<!-- internationalizer:unit markdown:license -->
## Licencja

[AGPL-3.0](../../LICENSE)

Zobacz [THIRD_PARTY_NOTICES.md](../../THIRD_PARTY_NOTICES.md), aby zapoznać się z informacjami o zależnościach.

<!-- internationalizer:unit markdown:contributing -->
## Współtworzenie

Zobacz [CONTRIBUTING.md](../../CONTRIBUTING.md), aby zapoznać się z konfiguracją środowiska programistycznego i wytycznymi. Wszystkie wkłady wymagają zatwierdzenia DCO.
