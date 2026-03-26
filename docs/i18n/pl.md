> [English (original)](../../README.md)

<p align="center">
  <img src="../../assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

Oparty na AI potok internacjonalizacji dla projektów programistycznych. Tłumacz, weryfikuj i zarządzaj plikami i18n przy użyciu LLM.

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](LICENSE)

<p align="center">
<a href="docs/i18n/ar.md">العربية</a> · <a href="docs/i18n/bn.md">বাংলা</a> · <a href="docs/i18n/de.md">Deutsch</a> · <a href="docs/i18n/es.md">Español</a> · <a href="docs/i18n/fr.md">Français</a> · <a href="docs/i18n/hi.md">हिन्दी</a> · <a href="docs/i18n/id.md">Indonesia</a> · <a href="docs/i18n/ja.md">日本語</a> · <a href="docs/i18n/pa.md">ਪੰਜਾਬੀ</a> · <a href="docs/i18n/pt-BR.md">Português</a> · <a href="docs/i18n/ru.md">Русский</a> · <a href="docs/i18n/te.md">తెలుగు</a><br><a href="docs/i18n/th.md">ไทย</a> · <a href="docs/i18n/uk.md">Українська</a> · <a href="docs/i18n/yue.md">粵語</a> · <a href="docs/i18n/zh-CN.md">简体中文</a> · <a href="docs/i18n/zh-TW.md">繁體中文</a>
</p>

---

## Dlaczego Internationalizer?

Większość narzędzi i18n to biblioteki wykonawcze (i18next, react-intl) lub platformy SaaS do zarządzania kluczami (Crowdin, Lokalise). Żadne z nich nie rozwiązuje dobrze rzeczywistego problemu tłumaczenia:

- **Ręczne tłumaczenie** nie skaluje się powyżej kilku języków
- **API tłumaczenia maszynowego** (Google Translate, DeepL) ignorują Twoją terminologię, ton i konwencje UI
- **Ogólne tłumaczenie LLM** działa lepiej, ale bez glosariuszy i przewodników po stylu otrzymujesz niespójne wyniki

Internationalizer jest inny. To **potok CLI**, który łączy tłumaczenie LLM z:

- **Glosariuszami dla poszczególnych języków** — wymuszają spójną terminologię w całej aplikacji
- **Przewodnikami po stylu dla poszczególnych języków** — kontrolują ton, formalność, liczby mnogie i typografię
- **Pamięcią tłumaczeniową** — pomijają niezmienione ciągi znaków, oszczędzając pieniądze na wywołaniach API
- **Weryfikacją kluczy** — wychwytują brakujące tłumaczenia i niezgodności interpolacji przed wdrożeniem

## Instalacja

Zainstaluj z npm:

```bash
npm install -g internationalizer
```

Lub uruchom bez globalnej instalacji:

```bash
npx internationalizer --help
```

Pakiet npm instaluje pasujący, wstępnie skompilowany plik binarny z npm poprzez opcjonalne zależności specyficzne dla platformy.

Zainstaluj za pomocą Go:

```bash
go install github.com/Tom-R-Main/Internationalizer/cmd/internationalizer@latest
```

Lub zbuduj ze źródeł:

```bash
git clone https://github.com/Tom-R-Main/Internationalizer.git
cd Internationalizer
go build -o internationalizer ./cmd/internationalizer
```

## Pakiety npm

- Tagi Git i wersje pakietów npm muszą być zgodne, na przykład `v0.1.0` i `0.1.0`
- Główny pakiet `internationalizer` zależy od pakietów platformowych, takich jak `internationalizer-darwin-arm64`
- Obsługiwane cele npm: macOS arm64/x64, Linux arm64/x64, Windows x64
- Publikowanie w CI wymaga sekretu GitHub o nazwie `NPM_TOKEN`

## Szybki start

1. Utwórz plik konfiguracyjny w katalogu głównym projektu:

```yaml
# .internationalizer.yml
source_locale: en
target_locales: [fr, de, es, ja]
source_path: locales/en.json

llm:
  provider: gemini
  model: gemini-3.1-pro-preview
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

5. Zweryfikuj wszystkie języki:

```bash
internationalizer validate
```

## Polecenia

### `translate`

Znajdź brakujące klucze i przetłumacz je za pomocą LLM.

```bash
internationalizer translate                    # przetłumacz wszystkie języki
internationalizer translate -l fr              # przetłumacz tylko język francuski
internationalizer translate --dry-run          # podgląd bez wywołań API
internationalizer translate --batch-size 20    # mniejsze partie
internationalizer translate --concurrency 2    # mniej równoległych wywołań
```

### `validate`

Sprawdź wszystkie pliki językowe pod kątem brakujących kluczy, dodatkowych kluczy i niezgodności interpolacji.

```bash
internationalizer validate                     # czytelne dla człowieka dane wyjściowe
internationalizer validate --json              # format JSON czytelny dla maszyn
internationalizer validate -q                  # tylko kod wyjścia
```

### `detect`

Automatycznie wykryj framework i18n i zasugeruj konfigurację.

```bash
internationalizer detect
```

Obsługuje: react-i18next, next-intl, vue-i18n, czysty JSON, dokumenty Markdown.

### `glossary`

Zarządzaj terminami w glosariuszu dla poszczególnych języków, które są wymuszane podczas tłumaczenia.

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

Zarządzaj pamięcią tłumaczeniową (pamięć podręczna JSONL wcześniej przetłumaczonych ciągów znaków).

```bash
internationalizer tm stats                     # pokaż liczbę rekordów
internationalizer tm export                    # zrzuć jako JSON
internationalizer tm clear --force             # usuń wszystkie rekordy
```

## Dokumentacja konfiguracji

```yaml
# .internationalizer.yml

# Język źródłowy (domyślnie: en)
source_locale: en

# Języki docelowe tłumaczenia (wymagane)
target_locales: [fr, de, es, ja, zh-CN, ar]

# Ścieżka do źródłowego pliku językowego (wymagane)
source_path: locales/en.json

# Ustawienia dostawcy LLM
llm:
  # Dostawca: "anthropic", "openai", "gemini" lub "openrouter" (domyślnie: gemini)
  provider: gemini

  # Domyślne nazwy modeli według dostawcy:
  #   anthropic:  claude-sonnet-4-6
  #   openai:     gpt-5.4
  #   gemini:     gemini-3.1-pro-preview
  #   openrouter: google/gemini-3-flash-preview
  model: gemini-3.1-pro-preview

  # Zmienna środowiskowa zawierająca klucz API
  api_key_env: GOOGLE_AI_STUDIO_API_KEY

  # Bazowy adres URL dla punktów końcowych zgodnych z OpenAI (opcjonalnie)
  # base_url: https://api.openai.com

# Liczba kluczy na wywołanie LLM (domyślnie: 40)
batch_size: 40

# Równoległe wywołania LLM (domyślnie: 4)
concurrency: 4

# Katalog zawierający pliki Markdown z przewodnikami po stylu dla poszczególnych języków (domyślnie: style-guides)
style_guides_dir: style-guides

# Katalog zawierający pliki JSON z glosariuszami dla poszczególnych języków (domyślnie: glossary)
glossary_dir: glossary

# Ścieżka do pliku pamięci tłumaczeniowej (domyślnie: .internationalizer/tm.jsonl)
tm_path: .internationalizer/tm.jsonl
```

## Przewodniki po stylu

Przewodniki po stylu to pliki Markdown, które są wstrzykiwane do promptu tłumaczenia LLM. Kontrolują one ton, formalność, typografię i inne konwencje specyficzne dla danego języka.

```
style-guides/
  _conventions.md    # wspólne reguły dla wszystkich języków
  fr.md              # reguły specyficzne dla języka francuskiego
  ja.md              # reguły specyficzne dla języka japońskiego
  ar.md              # reguły specyficzne dla języka arabskiego
```

### Wspólne konwencje (`_conventions.md`)

Zdefiniuj reguły, które mają zastosowanie do wszystkich języków: składnię interpolacji, zachowywanie tagów HTML, konwencje typów ciągów znaków (przyciski vs etykiety vs błędy) itp.

### Przewodniki dla poszczególnych języków (`{locale}.md`)

Zdefiniuj reguły specyficzne dla danego języka: rejestr formalności (np. ty vs Pan/Pani), interpunkcję (cudzysłowy francuskie, odwrócone znaki zapytania), formy liczby mnogiej, formatowanie dat/liczb oraz glosariusz terminologiczny.

Zobacz [`examples/react-app/style-guides/`](examples/react-app/style-guides/), aby zapoznać się z działającym przykładem.

## Format glosariusza

Pliki glosariusza to tablice JSON przechowywane w `{glossary_dir}/{locale}.json`:

```json
[
  {
    "source": "Dashboard",
    "target": "Tableau de bord",
    "ignore_case": false,
    "whole_word": true
  }
]
```

Terminy są wstrzykiwane do promptu LLM jako tabela terminologiczna, co zapewnia spójne tłumaczenie kluczowych terminów w całej aplikacji.

## Pamięć tłumaczeniowa

Pamięć tłumaczeniowa jest przechowywana jako plik JSONL (jeden rekord JSON na linię). Każdy rekord zawiera:

- Źródłowy klucz i wartość
- Przetłumaczoną wartość
- Skrót SHA-256 wartości źródłowej
- Znacznik czasu

Podczas kolejnych uruchomień niezmienione ciągi znaków są serwowane z pamięci podręcznej TM bez wywoływania LLM, co oszczędza czas i koszty API. Plik TM jest przyjazny dla systemu Git i może być commitowany wraz z plikami językowymi.

## Obsługiwane formaty

| Format | Rozszerzenia | Tryb |
|--------|-----------|------|
| JSON | `.json` | Klucz-wartość (zagnieżdżone, spłaszczone z notacją kropkową) |
| YAML | `.yml`, `.yaml` | Klucz-wartość (zachowuje komentarze i kolejność) |
| Markdown | `.md`, `.mdx` | Tłumaczenie całego dokumentu |

## Wykrywanie typu projektu

Polecenie `internationalizer detect` identyfikuje Twoją konfigurację i18n, sprawdzając:

- Zależności w `package.json` dla react-i18next, next-intl lub vue-i18n
- Struktury katalogów pasujące do popularnych wzorców językowych
- Rozszerzenia plików i konwencje nazewnictwa

## Architektura

```
cmd/internationalizer/     Punkt wejścia CLI i definicje poleceń
internal/
  config/                  Ładowanie konfiguracji YAML z wartościami domyślnymi
  detect/                  Automatyczne wykrywanie typu projektu
  formats/                 Parsery formatów (JSON, YAML, Markdown)
  glossary/                Zarządzanie glosariuszem dla poszczególnych języków
  llm/                     Interfejs dostawcy LLM + implementacje
    anthropic.go           Backend Anthropic Claude
    openai.go              Backend OpenAI / kompatybilny
    gemini.go              Backend Google Gemini przez AI Studio
                           OpenRouter używa openai.go z niestandardowym base_url
  styleguide/              Moduł ładujący przewodnik po stylu
  tm/                      Pamięć tłumaczeniowa JSONL
  translate/               Orkiestrator tłumaczeń
  validate/                Weryfikacja i porównywanie języków
```

## Porównanie z alternatywami

| Funkcja | Internationalizer | i18next | Crowdin | Ogólne LLM |
|---------|------------------|---------|---------|-------------|
| Tłumaczenie oparte na LLM | Tak | Nie | Częściowo | Tak |
| Przewodniki po stylu dla języków | Tak | Nie | Nie | Nie |
| Wymuszanie glosariusza | Tak | Nie | Tak | Nie |
| Pamięć tłumaczeniowa | Tak | Nie | Tak | Nie |
| CLI / wykonanie lokalne | Tak | N/D | Nie | Ręcznie |
| Pliki przyjazne dla Git | Tak | Tak | Częściowo | Ręcznie |
| Brak zależności od SaaS | Tak | Tak | Nie | Różnie |
| Open source (AGPL-3.0) | Tak | Tak | Nie | Różnie |

## Licencja

[AGPL-3.0](LICENSE)

## Współtworzenie

Zobacz [CONTRIBUTING.md](CONTRIBUTING.md), aby zapoznać się z konfiguracją środowiska programistycznego i wytycznymi. Wszystkie wkłady wymagają zatwierdzenia DCO.

