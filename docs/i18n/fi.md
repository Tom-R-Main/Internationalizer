> [English (original)](../../README.md)

<p align="center">
  <img src="../../assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

Tekoälypohjainen kansainvälistämisen työnkulku ohjelmistoprojekteille. Käännä, vahvista ja hallitse i18n-tiedostoja LLM-mallien avulla.

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](LICENSE)

<p align="center">
<a href="docs/i18n/ar.md">العربية</a> · <a href="docs/i18n/bn.md">বাংলা</a> · <a href="docs/i18n/de.md">Deutsch</a> · <a href="docs/i18n/es.md">Español</a> · <a href="docs/i18n/fr.md">Français</a> · <a href="docs/i18n/hi.md">हिन्दी</a> · <a href="docs/i18n/id.md">Indonesia</a> · <a href="docs/i18n/ja.md">日本語</a> · <a href="docs/i18n/pa.md">ਪੰਜਾਬੀ</a> · <a href="docs/i18n/pt-BR.md">Português</a> · <a href="docs/i18n/ru.md">Русский</a> · <a href="docs/i18n/te.md">తెలుగు</a><br>
<a href="docs/i18n/th.md">ไทย</a> · <a href="docs/i18n/uk.md">Українська</a> · <a href="docs/i18n/yue.md">粵語</a> · <a href="docs/i18n/zh-CN.md">简体中文</a> · <a href="docs/i18n/zh-TW.md">繁體中文</a>
</p>

---

## Miksi Internationalizer?

Useimmat i18n-työkalut ovat joko ajonaikaisia kirjastoja (i18next, react-intl) tai avaintenhallintaan tarkoitettuja SaaS-alustoja (Crowdin, Lokalise). Mikään niistä ei ratkaise varsinaista käännösongelmaa hyvin:

- **Manuaalinen kääntäminen** ei skaalaudu muutamaa kieltä pidemmälle
- **Konekäännös-API:t** (Google Translate, DeepL) jättävät huomiotta terminologiasi, sävysi ja käyttöliittymäsi käytännöt
- **Yleinen LLM-kääntäminen** toimii paremmin, mutta ilman sanastoja ja tyylioppaita tulokset ovat epäjohdonmukaisia

Internationalizer on erilainen. Se on **CLI-työnkulku**, joka yhdistää LLM-kääntämisen seuraaviin ominaisuuksiin:

- **Kielikohtaiset sanastot** — pakota johdonmukainen terminologia koko sovelluksessasi
- **Kielikohtaiset tyylioppaat** — hallitse sävyä, muodollisuutta, monikkomuotoja ja typografiaa
- **Käännösmuisti** — ohita muuttumattomat merkkijonot ja säästä rahaa API-kutsuissa
- **Avainten vahvistus** — huomaa puuttuvat käännökset ja interpolaatiovirheet ennen julkaisua

## Asennus

Asenna npm:n kautta:

```bash
npm install -g internationalizer
```

Tai suorita ilman globaalia asennusta:

```bash
npx internationalizer --help
```

npm-paketti asentaa vastaavan valmiiksi käännetyn binäärin npm:stä alustakohtaisten valinnaisten riippuvuuksien kautta.

Asenna Go:n avulla:

```bash
go install github.com/Tom-R-Main/Internationalizer/cmd/internationalizer@latest
```

Tai käännä lähdekoodista:

```bash
git clone https://github.com/Tom-R-Main/Internationalizer.git
cd Internationalizer
go build -o internationalizer ./cmd/internationalizer
```

## npm-paketit

- Git-tunnisteiden ja npm-pakettien versioiden on vastattava toisiaan, esimerkiksi `v0.1.0` ja `0.1.0`
- Juuripaketti `internationalizer` riippuu alustapaketeista, kuten `internationalizer-darwin-arm64`
- Tuetut npm-kohteet: macOS arm64/x64, Linux arm64/x64, Windows x64
- CI-julkaisu vaatii GitHub-salaisuuden nimeltä `NPM_TOKEN`

## Pika-aloitus

1. Luo asetustiedosto projektisi juureen:

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

2. Aseta API-avaimesi:

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

3. Esikatsele, mitä käännetään:

```bash
internationalizer translate --dry-run
```

4. Suorita käännös:

```bash
internationalizer translate
```

5. Vahvista kaikki kielitiedostot:

```bash
internationalizer validate
```

## Komennot

### `translate`

Etsi puuttuvat avaimet ja käännä ne LLM:n avulla.

```bash
internationalizer translate                    # translate all locales
internationalizer translate -l fr              # translate French only
internationalizer translate --dry-run          # preview without API calls
internationalizer translate --batch-size 20    # smaller batches
internationalizer translate --concurrency 2    # fewer parallel calls
```

### `validate`

Tarkista kaikki kielitiedostot puuttuvien avainten, ylimääräisten avainten ja interpolaatiovirheiden varalta.

```bash
internationalizer validate                     # human-readable output
internationalizer validate --json              # machine-readable JSON
internationalizer validate -q                  # exit code only
```

### `detect`

Tunnista i18n-kehys automaattisesti ja ehdota asetuksia.

```bash
internationalizer detect
```

Tukee seuraavia: react-i18next, next-intl, vue-i18n, tavallinen JSON, markdown-asiakirjat.

### `glossary`

Hallitse kielikohtaisia sanastotermejä, joita pakotetaan käännöksen aikana.

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

Hallitse käännösmuistia (aiemmin käännettyjen merkkijonojen JSONL-välimuisti).

```bash
internationalizer tm stats                     # show record counts
internationalizer tm export                    # dump as JSON
internationalizer tm clear --force             # delete all records
```

## Asetusten viiteopas

```yaml
# .internationalizer.yml

# Source language (default: en)
source_locale: en

# Languages to translate into (required)
target_locales: [fr, de, es, ja, zh-CN, ar]

# Path to the source locale file (required)
source_path: locales/en.json

# LLM provider settings
llm:
  # Provider: "anthropic", "openai", "gemini", or "openrouter" (default: gemini)
  provider: gemini

  # Model name defaults by provider:
  #   anthropic:  claude-sonnet-4-6
  #   openai:     gpt-5.4
  #   gemini:     gemini-3.1-pro-preview
  #   openrouter: google/gemini-3-flash-preview
  model: gemini-3.1-pro-preview

  # Environment variable containing the API key
  api_key_env: GOOGLE_AI_STUDIO_API_KEY

  # Base URL for OpenAI-compatible endpoints (optional)
  # base_url: https://api.openai.com

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
```

## Tyylioppaat

Tyylioppaat ovat Markdown-tiedostoja, jotka syötetään LLM-käännöspromptiin. Ne ohjaavat sävyä, muodollisuutta, typografiaa ja muita kielikohtaisia käytäntöjä.

```
style-guides/
  _conventions.md    # shared rules for all languages
  fr.md              # French-specific rules
  ja.md              # Japanese-specific rules
  ar.md              # Arabic-specific rules
```

### Jaetut käytännöt (`_conventions.md`)

Määritä säännöt, jotka koskevat kaikkia kieliä: interpolaatiosyntaksi, HTML:n säilyttäminen, merkkijonotyyppien käytännöt (painikkeet vs. tunnisteet vs. virheet) jne.

### Kielikohtaiset oppaat (`{locale}.md`)

Määritä kielikohtaiset säännöt: muodollisuusrekisteri (sinuttelu vs. teitittely), välimerkit (lainausmerkit, ylösalaisin olevat kysymysmerkit), monikkomuodot, päivämäärien/numeroiden muotoilu ja terminologiasanasto.

Katso toimiva esimerkki kohdasta [`examples/react-app/style-guides/`](examples/react-app/style-guides/).

## Sanaston muoto

Sanastotiedostot ovat JSON-taulukoita, jotka on tallennettu sijaintiin `{glossary_dir}/{locale}.json`:

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

Termit syötetään LLM-promptiin terminologiataulukkona, mikä varmistaa avaintermien johdonmukaisen kääntämisen koko sovelluksessasi.

## Käännösmuisti

Käännösmuisti tallennetaan JSONL-tiedostona (yksi JSON-tietue per rivi). Jokainen tietue sisältää seuraavat:

- Lähdeavaimen ja -arvon
- Käännetyn arvon
- Lähdearvon SHA-256-tiivisteen
- Aikaleiman

Seuraavilla suorituskerroilla muuttumattomat merkkijonot tarjoillaan TM-välimuistista ilman LLM-kutsua, mikä säästää sekä aikaa että API-kustannuksia. TM-tiedosto on Git-yhteensopiva ja se voidaan kommitoida kielitiedostojesi mukana.

## Tuetut muodot

| Muoto | Tunnisteet | Tila |
|--------|-----------|------|
| JSON | `.json` | Avain-arvo (sisäkkäinen, pistenotaatiolla litistetty) |
| YAML | `.yml`, `.yaml` | Avain-arvo (säilyttää kommentit ja järjestyksen) |
| Markdown | `.md`, `.mdx` | Koko asiakirjan käännös |

## Projektityypin tunnistus

`internationalizer detect` tunnistaa i18n-asetuksesi tarkistamalla seuraavat:

- `package.json`-riippuvuudet (react-i18next, next-intl tai vue-i18n)
- Yleisiä kieliasetusten malleja vastaavat hakemistorakenteet
- Tiedostotunnisteet ja nimeämiskäytännöt

## Arkkitehtuuri

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
  styleguide/              Style guide loader
  tm/                      JSONL translation memory
  translate/               Translation orchestrator
  validate/                Locale validation and diffing
```

## Vertailu vaihtoehtoihin

| Ominaisuus | Internationalizer | i18next | Crowdin | Yleinen LLM |
|---------|------------------|---------|---------|-------------|
| LLM-pohjainen käännös | Kyllä | Ei | Osittainen | Kyllä |
| Kielikohtaiset tyylioppaat | Kyllä | Ei | Ei | Ei |
| Sanaston pakotus | Kyllä | Ei | Kyllä | Ei |
| Käännösmuisti | Kyllä | Ei | Kyllä | Ei |
| CLI / paikallinen suoritus | Kyllä | Ei sovellu | Ei | Manuaalinen |
| Git-yhteensopivat tiedostot | Kyllä | Kyllä | Osittainen | Manuaalinen |
| Ei SaaS-riippuvuutta | Kyllä | Kyllä | Ei | Vaihtelee |
| Avoin lähdekoodi (AGPL-3.0) | Kyllä | Kyllä | Ei | Vaihtelee |

## Lisenssi

[AGPL-3.0](LICENSE)

## Osallistuminen

Katso kehitysympäristön asennusohjeet ja suuntaviivat tiedostosta [CONTRIBUTING.md](CONTRIBUTING.md). Kaikki osallistuminen vaatii DCO-allekirjoituksen.

