> [English (original)](../../README.md)

<p align="center">
  <img src="../../assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

Tekoälypohjainen kansainvälistämisen työnkulku ohjelmistoprojekteille. Käännä, validoi ja hallitse i18n-tiedostoja LLM-mallien avulla.

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](../../LICENSE)

<p align="center">
<a href="ar.md">العربية</a> · <a href="bn.md">বাংলা</a> · <a href="cs.md">Čeština</a> · <a href="da.md">Dansk</a> · <a href="de.md">Deutsch</a> · <a href="el.md">Ελληνικά</a> · <a href="es.md">Español</a> · <a href="fi.md">Suomi</a> · <a href="fr.md">Français</a> · <a href="he.md">עברית</a> · <a href="hi.md">हिन्दी</a> · <a href="id.md">Indonesia</a> · <a href="it.md">Italiano</a> · <a href="ja.md">日本語</a> · <a href="ko.md">한국어</a> · <a href="ms.md">Bahasa Melayu</a><br><a href="nl.md">Nederlands</a> · <a href="pa.md">ਪੰਜਾਬੀ</a> · <a href="pl.md">Polski</a> · <a href="pt-BR.md">Português</a> · <a href="ro.md">Română</a> · <a href="ru.md">Русский</a> · <a href="sv.md">Svenska</a> · <a href="te.md">తెలుగు</a> · <a href="th.md">ไทย</a> · <a href="tr.md">Türkçe</a> · <a href="uk.md">Українська</a> · <a href="vi.md">Tiếng Việt</a> · <a href="yue.md">粵語</a> · <a href="zh-CN.md">简体中文</a> · <a href="zh-TW.md">繁體中文</a>
</p>

---

<!-- internationalizer:unit markdown:why-internationalizer -->
## Miksi Internationalizer?

Useimmat i18n-työkalut ovat joko ajonaikaisia kirjastoja (i18next, react-intl) tai avaintenhallintaan tarkoitettuja SaaS-alustoja (Crowdin, Lokalise). Mikään niistä ei ratkaise varsinaista käännösongelmaa kunnolla:

- **Manuaalinen kääntäminen** ei skaalaudu muutamaa kieltä pidemmälle
- **Konekäännösrajapinnat** (Google Translate, DeepL) jättävät huomiotta terminologian, sävyn ja käyttöliittymäkäytännöt
- **Yleinen LLM-kääntäminen** toimii paremmin, mutta ilman sanastoja ja tyylioppaita tulokset ovat epäjohdonmukaisia

Internationalizer tekee poikkeuksen. Se on **CLI-työnkulku**, joka yhdistää LLM-kääntämisen ja seuraavat ominaisuudet:

- **Kielikohtaiset sanastot** – valvo johdonmukaista terminologiaa koko sovelluksessa
- **Kielikohtaiset tyylioppaat** – hallitse sävyä, muodollisuutta, monikkomuotoja ja typografiaa
- **Käännösmuisti** – ohita muuttumattomat merkkijonot ja säästä rahaa API-kutsuissa
- **Deterministinen validointi** – havaitse puuttuvat tai ylimääräiset avaimet, suojatun rakenteen poikkeamat, sanasto-ongelmat sekä monikko- tai ICU-virheet ennen julkaisua

<!-- internationalizer:unit markdown:installation -->
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

Asenna Go-työkalulla:

```bash
go install github.com/Tom-R-Main/Internationalizer/cmd/internationalizer@latest
```

Tai rakenna lähdekoodista:

```bash
git clone https://github.com/Tom-R-Main/Internationalizer.git
cd Internationalizer
go build -o internationalizer ./cmd/internationalizer
```

<!-- internationalizer:unit markdown:npm-packages -->
## npm-paketit

- Git-tunnisteiden ja npm-pakettien versioiden on vastattava toisiaan, esimerkiksi `v0.1.0` ja `0.1.0`
- Juuritason `internationalizer`-paketti riippuu alustapaketeista, kuten `internationalizer-darwin-arm64`
- Tuetut npm-kohteet: macOS arm64/x64, Linux arm64/x64, Windows x64
- Julkaisu CI-ympäristöstä vaatii GitHub-salaisuuden nimeltä `NPM_TOKEN`

<!-- internationalizer:unit markdown:quick-start -->
## Pika-aloitus

1. Luo asetustiedosto projektin juureen:

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

5. Validoi kaikki kieliympäristöt:

```bash
internationalizer validate
```

<!-- internationalizer:unit markdown:commands -->
## Komennot

### `translate`

Etsi puuttuvat tai vanhentuneet avaimet ja käännä ne LLM:n avulla.

```bash
internationalizer translate                    # translate all locales
internationalizer translate -l fr              # translate French only
internationalizer translate --dry-run          # preview without API calls
internationalizer translate --adopt-existing   # baseline existing translations without API calls
internationalizer translate --refresh-policy   # refresh prompt/style/model-stale entries
internationalizer translate --batch-size 20    # smaller batches
internationalizer translate --concurrency 2    # fewer parallel calls
```

Käännöstila raportoi toisistaan riippumatta puuttuvat (missing), lähdetiedoltaan vanhentuneet (source-stale), käytännöltään vanhentuneet (policy-stale), ajan tasalla olevat (current) ja manuaalisesti muokatut (manually edited) tilat, joten manuaalinen muokkaus ei voi piilottaa lähteen tai käytännön muutosta. Käytännöltään vanhentuneet arvot raportoidaan, mutta ne käännetään uudelleen vain lipulla `--refresh-policy`. Manuaalisesti muokattuja arvoja ei koskaan korvata automaattisesti. Käytä lippua `--adopt-existing`, kun otat manifestin käyttöön aiemmin tarkastetuille käännöksille tai kun hyväksyt tarkastetun manuaalisen muokkauksen nimenomaisesti uudeksi perustasoksi.

### `validate`

Tarkista kaikki kielialuetiedostot suhteessa niiden lähdepaketteihin. Oletusarvoinen validointi tarkistaa rakenteellisen kattavuuden (vaadittujen kohdeavaimien läsnäolon prosenttiosuuden), raportoi ylimääräiset avaimet varoituksina ja epäonnistuu, jos avaimia puuttuu, interpolaatioissa on ristiriitoja tai ICU MessageFormat -rakenne on virheellinen.

```bash
internationalizer validate                     # human-readable output
internationalizer validate --json              # machine-readable JSON
internationalizer validate -q                  # exit code only
internationalizer validate --strict             # enforce translation quality rules
internationalizer validate --require-state      # require current manifest provenance
```

`--strict` raportoi myös käännöskattavuuden. Lähdearvon kanssa identtinen kielellinen arvo katsotaan kääntämättömäksi, ellei sanastossa ole nimenomaista koko arvoa koskevaa merkintää, jossa lähde ja kohde ovat täysin samat; `ignore_case`-asetusta noudatetaan, mutta pidempään arvoon sisältyvä sanastotermi ei muodosta poikkeusta. Tiukka tila epäonnistuu, jos löytyy ylimääräisiä avaimia, lähteen kanssa identtisiä arvoja, muuttuneita interpolaatio-, HTML-, koodi- tai Markdown-linkkirakenteita, sanastorikkomuksia tai määritettyjen monikkomuotojen puutteita.

`--require-state` tarkistaa kunkin kohteen `.internationalizer.lock`-tiedostoa vasten. Se epäonnistuu, jos avainta ei seurata tai jos sen tallennettu lähde, käännöskäytäntö tai kohteen tiiviste on vanhentunut. Se voidaan yhdistää lipun `--strict` kanssa.

Ihmisille luettavat ja JSON-raportit käyttävät pysyviä havaintokoodeja:

| Koodi | Merkitys |
| --- | --- |
| `missing_key` / `extra_key` | Lähde- ja kohdeavainten joukot eroavat toisistaan |
| `blank_translation` | Ei-tyhjällä lähteellä on tyhjä kohde tiukassa tilassa |
| `source_identical` | Tiukan tilan kielellinen arvo on kääntämätön |
| `protected_structure_mismatch` | Interpolaatio-, HTML-, koodi- tai linkkirakenne on muuttunut |
| `glossary_violation` | Hyväksyttyä kohdetermiä tai sen varianttia ei löytynyt |
| `plural_form_missing` | Kielialueelle määritetty monikkomuoto puuttuu |
| `icu_message_syntax` | Lähteen tai kohteen ICU-viesti on virheellisesti muotoiltu |
| `icu_argument_mismatch` | ICU-argumenttien nimet, tyypit tai muotoilijatyylit eroavat |
| `icu_selector_mismatch` | Valitsimet eroavat toisistaan tai monikkoluokka ei kelpaa kohdekielialueelle |
| `untracked` | Kohteelle ei löydy tietuetta manifestista |
| `source_stale` | Lähdesisältö on muuttunut tallennetun käännöksen jälkeen |
| `policy_stale` | Luotu kehote tai malliasetukset ovat muuttuneet |
| `target_modified` | Kohdesisältö poikkeaa manifestin tietueesta |

### `detect`

Tunnista i18n-kehys automaattisesti ja ehdota määrityksiä.

```bash
internationalizer detect
```

Tukee seuraavia: react-i18next, next-intl, vue-i18n, tavallinen JSON, Markdown-dokumentit.

### `glossary`

Hallitse kielikohtaisia sanastotermejä, joiden käyttöä valvotaan käännöksen aikana.

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

Hallitse käännösmuistia (aiemmin käännettyjen merkkijonojen JSONL-välimuistia).

```bash
internationalizer tm stats                     # show record counts
internationalizer tm export                    # dump as JSON
internationalizer tm clear --force             # delete all records
```

<!-- internationalizer:unit markdown:configuration-reference -->
## Asetusten viiteopas

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

Kielialuetunnisteiden on oltava oikein muodostettuja BCP 47 -tunnisteita, kuten `fr`, `pt-BR` tai `sr-Latn-RS`. Kanonisesti vastaavat kohdekielialueet hylätään kaksoiskappaleina, ja kielialuekohtaiset palveluntarjoajien ohitukset täsmäävät kanonisesti vastaavaan kirjoitusasuun. Yllä olevassa esimerkissä ohittamattomat kielialueet – mukaan lukien japani – perivät globaalit Gemini-asetukset.

ICU MessageFormat -arvot jäsennetään rakenteellisesti. Yksinkertaisia argumentteja sekä tyyppejä `select`, `plural`, `selectordinal`, `number`, `date` ja `time` tuetaan, mukaan lukien sisäkkäiset viestit, monikkojen siirtymät (plural offsets), tarkat lukumäärävalitsimet ja `#`. Validointi tarkistaa syntaksin, argumenttien tyypit ja muotoilijatyylit, monikon siirtymät, select-haarojen identiteetin sekä kohdekielialueen CLDR-monikkoluokat. Palveluntarjoajan tuottama tuloste, joka rikkoo näitä invariantteja, hylätään ennen kuin kielialuetiedostoa tai käännösmuistitietuetta kirjoitetaan.

Käytettäessä arvoa `i18next-v4` tunnistetut lähteen monikkoperheet laajennetaan käännöksen aikana kohdekielialueen CLDR-luokkiin. Vain kohteessa esiintyvä luokka käyttää lähdeperheen `_other`-arvoa käännösmallinaan. Tiukka validointi edellyttää näitä kohdeluokkia; vain lähteessä esiintyvät luokat ovat valinnaisia niille kohdekielialueille, jotka eivät käytä niitä.

<!-- internationalizer:unit markdown:style-guides -->
## Tyylioppaat

Tyylioppaat ovat Markdown-tiedostoja, jotka syötetään LLM-käännöskomentoon. Ne ohjaavat sävyä, muodollisuutta, typografiaa ja muita kielikohtaisia käytäntöjä.

```
style-guides/
  _conventions.md    # shared rules for all languages
  fr.md              # French-specific rules
  ja.md              # Japanese-specific rules
  ar.md              # Arabic-specific rules
```

### Jaetut käytännöt (`_conventions.md`)

Määritä kaikkiin kieliin sovellettavat säännöt: interpolaatiosyntaksi, HTML:n säilyttäminen, merkkijonotyyppien käytännöt (painikkeet vs. selitteet vs. virheet) jne.

### Kielikohtaiset oppaat (`{locale}.md`)

Määritä kielikohtaiset säännöt: muodollisuusrekisteri (sinuttelu vs. teitittely), välimerkit (lainausmerkit, ylösalaisin käännetyt kysymysmerkit), monikkomuodot, päivämäärien ja numeroiden muotoilu sekä terminologiasanasto.

Tyylioppaat ovat pysyviä käytäntösyötteitä, eivät luotua tulostetta. Internationalizer lukee ne, muttei koskaan kirjoita niitä uudelleen. Niiden sisältö tiivistetään erikseen sanastosta ja kehotesopimuksesta, joten sovelluskoodin muutos ei tee käännöksestä vanhentunutta. Oppaan muokkaaminen merkitsee kyseisen kielialueen tarkoituksellisesti käytäntötarkastukseen; sisäisen kehotteen sanamuodon muuttaminen ei tee niin, ellei myös kehotesopimuksen versio muutu.

Katso toimiva esimerkki kohdasta [`examples/react-app/style-guides/`](../../examples/react-app/style-guides/).

<!-- internationalizer:unit markdown:glossary-format -->
## Sanaston muoto

Sanastotiedostot ovat JSON-taulukoita, jotka tallennetaan polkuun `{glossary_dir}/{locale}.json`:

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

`variants` luettelee muut hyväksytyt kohdemuodot. `enforcement`-kentän arvo voi olla `error`, `warning` tai se voidaan jättää pois, jolloin käytetään oletusarvoista error-toimintaa. Termit syötetään LLM-kehotteeseen terminologiataulukkona, mikä varmistaa johdonmukaisen kääntämisen läpi koko sovelluksen. Tarkka merkintä, kuten `{"source":"API","target":"API"}`, myös vapauttaa kyseisen täysin lähdeidenttisen arvon tiukan tilan kääntämätön arvo -havainnoista; se ei vapauta pidempää arvoa, joka vain sisältää termin `API`.

<!-- internationalizer:unit markdown:translation-memory -->
## Käännösmuisti

Käännösmuisti tallennetaan JSONL-tiedostona (yksi JSON-tietue riviä kohden). Jokainen tietue sisältää seuraavat tiedot:

- Paketti (bundle), avain, lähdearvo, käännetty arvo ja kanoninen kohdekielialue
- Lähteen, tyylioppaan, sanaston, kehotesopimuksen ja yhdistetyn käytännön tiivisteet
- Palveluntarjoaja ja malli, joka tuotti käännöksen
- Aikaleima

Seuraavilla suorituskerroilla merkkijonot, joilla on sama lähde- ja käytäntötiiviste, palautetaan välimuistista kutsumatta LLM-mallia. Oletuspolku sijaitsee huomiotta jätetyssä `.internationalizer/`-hakemistossa, joten se pysyy paikallisena välimuistina. Aseta `tm_path` seurattuun sijaintiin, jos projektisi jakaa käännösmuistin tarkoituksellisesti. Katselmoitava `.internationalizer.lock`-manifesti versioidaan erikseen.

<!-- internationalizer:unit markdown:supported-formats -->
## Tuetut tiedostomuodot

| Muoto | Tiedostotunnisteet | Tila |
|--------|-----------|------|
| JSON | `.json` | Avain-arvo (sisäkkäinen, pistenotaatiolla litistetty) |
| YAML | `.yml`, `.yaml` | Avain-arvo (säilyttää kommentit ja järjestyksen) |
| Markdown | `.md`, `.mdx` | Johdanto-osa ja H2-tason osiot |

Markdown-kohteet sisältävät näkymättömiä `internationalizer:unit`-kommentteja ennen H2-osioita. Näiden pysyvien merkintöjen ansiosta Internationalizer voi lisätä, siirtää tai muokata yhtä lähdeosiota kääntämättä toisiinsa liittymättömiä osioita uudelleen. Olemassa olevat merkitsemättömät asiakirjat saavat merkinnät seuraavan onnistuneen päivityksen yhteydessä.

<!-- internationalizer:unit markdown:project-type-detection -->
## Projektityypin tunnistus

`internationalizer detect` tunnistaa i18n-määrityksesi tarkistamalla seuraavat asiat:

- `package.json`-riippuvuudet (react-i18next, next-intl tai vue-i18n)
- Yleisiä kielialuemalleja vastaavat hakemistorakenteet
- Tiedostopäätteet ja nimeämiskäytännöt

<!-- internationalizer:unit markdown:architecture -->
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
## Vertailu vaihtoehtoihin

| Ominaisuus | Internationalizer | i18next | Crowdin | Yleinen LLM |
|---------|------------------|---------|---------|-------------|
| LLM-pohjainen käännös | Kyllä | Ei | Osittainen | Kyllä |
| Kielikohtaiset tyylioppaat | Kyllä | Ei | Ei | Ei |
| Sanaston valvonta | Kyllä | Ei | Kyllä | Ei |
| Käännösmuisti | Kyllä | Ei | Kyllä | Ei |
| CLI- / paikallinen suoritus | Kyllä | Ei sovellu | Ei | Manuaalinen |
| Git-ystävälliset tiedostot | Kyllä | Kyllä | Osittainen | Manuaalinen |
| Ei SaaS-riippuvuutta | Kyllä | Kyllä | Ei | Vaihtelee |
| Avoin lähdekoodi (AGPL-3.0) | Kyllä | Kyllä | Ei | Vaihtelee |

<!-- internationalizer:unit markdown:license -->
## Lisenssi

[AGPL-3.0](../../LICENSE)

Katso riippuvuuksien ilmoitukset tiedostosta [THIRD_PARTY_NOTICES.md](../../THIRD_PARTY_NOTICES.md).

<!-- internationalizer:unit markdown:contributing -->
## Osallistuminen

Katso kehitysympäristön määritys ja ohjeet tiedostosta [CONTRIBUTING.md](../../CONTRIBUTING.md). Kaikki kontribuutiot vaativat DCO-hyväksynnän.
