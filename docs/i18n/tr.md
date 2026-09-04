> [English (original)](../../README.md)

<p align="center">
  <img src="../../assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

Yazılım projeleri için AI tabanlı uluslararasılaştırma ardışık düzeni. LLM'leri kullanarak i18n dosyalarını çevirin, doğrulayın ve yönetin.

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](../../LICENSE)

<p align="center">
<a href="ar.md">العربية</a> · <a href="bn.md">বাংলা</a> · <a href="cs.md">Čeština</a> · <a href="da.md">Dansk</a> · <a href="de.md">Deutsch</a> · <a href="el.md">Ελληνικά</a> · <a href="es.md">Español</a> · <a href="fi.md">Suomi</a> · <a href="fr.md">Français</a> · <a href="he.md">עברית</a> · <a href="hi.md">हिन्दी</a> · <a href="id.md">Indonesia</a> · <a href="it.md">Italiano</a> · <a href="ja.md">日本語</a> · <a href="ko.md">한국어</a> · <a href="ms.md">Bahasa Melayu</a><br><a href="nl.md">Nederlands</a> · <a href="pa.md">ਪੰਜਾਬੀ</a> · <a href="pl.md">Polski</a> · <a href="pt-BR.md">Português</a> · <a href="ro.md">Română</a> · <a href="ru.md">Русский</a> · <a href="sv.md">Svenska</a> · <a href="te.md">తెలుగు</a> · <a href="th.md">ไทย</a> · <a href="tr.md">Türkçe</a> · <a href="uk.md">Українська</a> · <a href="vi.md">Tiếng Việt</a> · <a href="yue.md">粵語</a> · <a href="zh-CN.md">简体中文</a> · <a href="zh-TW.md">繁體中文</a>
</p>

---

<!-- internationalizer:unit markdown:why-internationalizer -->
## Neden Internationalizer?

Çoğu i18n aracı ya çalışma zamanı kütüphaneleridir (i18next, react-intl) ya da anahtar yönetimi SaaS platformlarıdır (Crowdin, Lokalise). Hiçbiri asıl çeviri sorununu iyi bir şekilde çözmez:

- **Manuel çeviri** birkaç dilden sonra ölçeklenemez
- **Makine çevirisi API'leri** (Google Translate, DeepL) terminolojinizi, üslubunuzu ve kullanıcı arayüzü kurallarınızı göz ardı eder
- **Genel LLM çevirisi** daha iyi sonuç verir; ancak sözlükler ve stil kılavuzları olmadan tutarsız sonuçlar elde edersiniz

Internationalizer farklıdır. LLM çevirisini şunlarla birleştiren bir **CLI ardışık düzenidir**:

- **Dile özel sözlükler** — uygulamanız genelinde tutarlı terminolojiyi zorunlu kılar
- **Dile özel stil kılavuzları** — üslubu, resmiyeti, çoğullaştırmayı ve tipografiyi kontrol eder
- **Çeviri belleği** — değişmeyen dizeleri atlar, API çağrılarından tasarruf sağlar
- **Belirleyici doğrulama** — eksik veya fazladan anahtarları, korunan yapı kaymalarını, sözlük sorunlarını ve çoğul ya da ICU hatalarını yayına girmeden önce yakalar

<!-- internationalizer:unit markdown:installation -->
## Kurulum

npm üzerinden kurun:

```bash
npm install -g internationalizer
```

Veya genel bir kurulum yapmadan çalıştırın:

```bash
npx internationalizer --help
```

npm paketi, platforma özel isteğe bağlı bağımlılıklar aracılığıyla npm'den eşleşen önceden derlenmiş ikili dosyayı kurar.

Go ile kurun:

```bash
go install github.com/Tom-R-Main/Internationalizer/cmd/internationalizer@latest
```

Veya kaynaktan derleyin:

```bash
git clone https://github.com/Tom-R-Main/Internationalizer.git
cd Internationalizer
go build -o internationalizer ./cmd/internationalizer
```

<!-- internationalizer:unit markdown:npm-packages -->
## npm paketleri

- Git etiketleri ve npm paket sürümleri eşleşmelidir, örneğin `v0.1.0` ve `0.1.0`
- Kök `internationalizer` paketi, `internationalizer-darwin-arm64` gibi platform paketlerine bağlıdır
- Desteklenen npm hedefleri: macOS arm64/x64, Linux arm64/x64, Windows x64
- CI yayınlaması için `NPM_TOKEN` adlı bir GitHub gizli dizisi (secret) gerekir

<!-- internationalizer:unit markdown:quick-start -->
## Hızlı başlangıç

1. Projenizin kök dizininde bir yapılandırma dosyası oluşturun:

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

2. API anahtarınızı ayarlayın:

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

3. Nelerin çevrileceğini önizleyin:

```bash
internationalizer translate --dry-run
```

4. Çeviriyi çalıştırın:

```bash
internationalizer translate
```

5. Tüm yerel ayarları doğrulayın:

```bash
internationalizer validate
```

<!-- internationalizer:unit markdown:commands -->
## Komutlar

### `translate`

Eksik veya güncelliğini yitirmiş anahtarları bulun ve bunları bir LLM aracılığıyla çevirin.

```bash
internationalizer translate                    # translate all locales
internationalizer translate -l fr              # translate French only
internationalizer translate --dry-run          # preview without API calls
internationalizer translate --adopt-existing   # baseline existing translations without API calls
internationalizer translate --refresh-policy   # refresh prompt/style/model-stale entries
internationalizer translate --batch-size 20    # smaller batches
internationalizer translate --concurrency 2    # fewer parallel calls
```

Çeviri durumu; eksik, kaynak açısından eski, ilke açısından eski, güncel ve el ile düzenlenmiş durumları bağımsız olarak bildirir; böylece el ile yapılan bir düzenleme kaynak veya ilke değişikliğini gizleyemez. İlke açısından eski değerler bildirilir ancak yalnızca `--refresh-policy` ile yeniden çevrilir. El ile düzenlenen değerlerin üzerine hiçbir zaman otomatik olarak yazılmaz. Bildirimi (manifest) incelenmiş çevirilerle ilk kez tanıştırırken veya incelenmiş el ile düzenlemeyi açıkça yeni temel durum (baseline) olarak kabul ederken `--adopt-existing` seçeneğini kullanın.

### `validate`

Tüm yerel ayar dosyalarını kaynak paketlerine göre denetleyin. Varsayılan doğrulama; yapısal kapsamı (gerekli hedef anahtarların mevcut olma yüzdesini) denetler, fazladan anahtarları uyarı olarak bildirir ve eksik anahtarlar, interpolasyon uyuşmazlıkları veya geçersiz ICU MessageFormat yapısı durumunda işlemi başarısız sayar.

```bash
internationalizer validate                     # human-readable output
internationalizer validate --json              # machine-readable JSON
internationalizer validate -q                  # exit code only
internationalizer validate --strict             # enforce translation quality rules
internationalizer validate --require-state      # require current manifest provenance
```

`--strict` aynı zamanda çevrilmiş kapsamı da bildirir. Kaynağıyla birebir aynı olan dilsel bir değer, sözlükte değerin tamamı için tam olarak aynı kaynak ve aynı hedefe sahip bir girdi açıkça yer almadıkça çevrilmemiş kabul edilir; `ignore_case` dikkate alınır, ancak daha uzun bir değerin içine gömülmüş bir sözlük terimi muafiyet sağlamaz. Katı mod (strict mode); fazladan anahtarlar, kaynakla aynı değerler, değişen interpolasyon/HTML/kod/Markdown bağlantı yapısı, sözlük ihlalleri ve yapılandırılmış çoğul biçimler durumunda başarısız olur.

`--require-state`, her hedefi `.internationalizer.lock` dosyasına göre doğrular. Bir anahtar izlenmediğinde ya da kaydedilen kaynağı, çeviri ilkesi veya hedef karması (hash) güncelliğini yitirdiğinde başarısız olur. `--strict` ile birlikte kullanılabilir.

Kullanıcı dostu ve JSON raporları kararlı bulgu kodları kullanır:

| Kod | Anlamı |
| --- | --- |
| `missing_key` / `extra_key` | Kaynak ve hedef anahtar kümeleri farklı |
| `blank_translation` | Boş olmayan bir kaynağın katı mod hedefi boş |
| `source_identical` | Katı mod dilsel değeri çevrilmeden kalmış |
| `protected_structure_mismatch` | İnterpolasyon, HTML, kod veya bağlantı yapısı değişmiş |
| `glossary_violation` | Onaylanmış hiçbir hedef terim veya varyant bulunamadı |
| `plural_form_missing` | Yapılandırılmış yerel ayar çoğul biçimi eksik |
| `icu_message_syntax` | Kaynak veya hedef ICU iletisi bozuk |
| `icu_argument_mismatch` | ICU bağımsız değişken adları, türleri veya biçimlendirici stilleri farklı |
| `icu_selector_mismatch` | Seçiciler farklı veya bir çoğul kategorisi hedef yerel ayar için geçersiz |
| `untracked` | Hedef için bildirim kaydı yok |
| `source_stale` | Kaynak içerik, kaydedilen çeviriden sonra değişti |
| `policy_stale` | Oluşturulan istem veya model ayarları değişti |
| `target_modified` | Hedef içerik bildirim kaydından farklı |

### `detect`

i18n çerçevesini otomatik olarak algılayın ve bir yapılandırma önerin.

```bash
internationalizer detect
```

Desteklenenler: react-i18next, next-intl, vue-i18n, yalın JSON, markdown belgeleri.

### `glossary`

Çeviri sırasında zorunlu tutulan dile özel sözlük terimlerini yönetin.

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

Çeviri belleğini (daha önce çevrilmiş dizelerin JSONL önbelleği) yönetin.

```bash
internationalizer tm stats                     # show record counts
internationalizer tm export                    # dump as JSON
internationalizer tm clear --force             # delete all records
```

<!-- internationalizer:unit markdown:configuration-reference -->
## Yapılandırma referansı

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

Yerel ayar tanımlayıcıları, `fr`, `pt-BR` veya `sr-Latn-RS` gibi düzgün biçimlendirilmiş BCP 47 etiketleri olmalıdır. Kurallı eşdeğer hedef yerel ayarlar yinelenen olarak reddedilir ve yerel ayara özel sağlayıcı geçersiz kılmaları kurallı eşdeğer yazımla eşleşir. Yukarıdaki örnekte, geçersiz kılma içermeyen yerel ayarlar —Japonca dahil— genel Gemini yapılandırmasını devralır.

ICU MessageFormat değerleri yapısal olarak ayrıştırılır. Basit bağımsız değişkenler, `select`, `plural`, `selectordinal`, `number`, `date` ve `time`; iç içe geçmiş iletiler, çoğul uzaklıkları (offsets), tam sayı seçicileri ve `#` dahil olmak üzere desteklenir. Doğrulama; sözdizimini, bağımsız değişken türlerini ve biçimlendirici stillerini, çoğul uzaklıklarını, select dalı kimliğini ve hedef yerel ayarın CLDR çoğul kategorilerini denetler. Bu değişmezleri bozan sağlayıcı çıktısı, bir yerel ayar dosyası veya çeviri belleği kaydı yazılmadan önce reddedilir.

`i18next-v4` ile, tanınan kaynak çoğul aileleri çeviri sırasında hedef yerel ayarın CLDR kategorilerine genişletilir. Yalnızca hedefte bulunan bir kategori, kaynak ailesinin `_other` değerini çeviri şablonu olarak kullanır. Katı doğrulama bu hedef kategorilerini zorunlu kılar; yalnızca kaynakta bulunan kategoriler ise bunları kullanmayan hedef yerel ayarlar için isteğe bağlıdır.

<!-- internationalizer:unit markdown:style-guides -->
## Stil kılavuzları

Stil kılavuzları, LLM çeviri istemine eklenen Markdown dosyalarıdır. Üslubu, resmiyeti, tipografiyi ve dile özgü diğer kuralları kontrol ederler.

```
style-guides/
  _conventions.md    # shared rules for all languages
  fr.md              # French-specific rules
  ja.md              # Japanese-specific rules
  ar.md              # Arabic-specific rules
```

### Paylaşılan kurallar (`_conventions.md`)

Tüm diller için geçerli kuralları tanımlayın: interpolasyon sözdizimi, HTML'in korunması, dize türü kuralları (butonlar, etiketler, hatalar vb.).

### Dile özel kılavuzlar (`{locale}.md`)

Dile özgü kuralları tanımlayın: resmiyet düzeyi (sen/siz), noktalama işaretleri (açılı tırnaklar, ters soru işaretleri), çoğul biçimler, tarih/sayı biçimlendirmesi ve terminoloji sözlüğü.

Stil kılavuzları kalıcı ilke girdileridir, oluşturulan çıktılar değildir. Internationalizer bunları okur ancak asla üzerlerine yazmaz. İçerikleri sözlükten ve istem sözleşmesinden ayrı olarak karma işlemine tabi tutulur, bu nedenle bir uygulama kodu değişikliği çeviriyi geçersiz kılmaz. Bir kılavuzu düzenlemek, ilgili yerel ayarı kasıtlı olarak ilke incelemesi için işaretler; istem sözleşmesi sürümü değişmedikçe dahili istem metninin değiştirilmesi bunu tetiklemez.

Çalışan bir örnek için [`examples/react-app/style-guides/`](../../examples/react-app/style-guides/) dizinine bakın.

<!-- internationalizer:unit markdown:glossary-format -->
## Sözlük biçimi

Sözlük dosyaları, `{glossary_dir}/{locale}.json` yolunda saklanan JSON dizileridir:

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

`variants`, onaylanmış diğer hedef biçimleri listeler. `enforcement`; `error`, `warning` olabilir veya varsayılan hata davranışı için belirtilmeyebilir. Terimler, LLM istemine bir terminoloji tablosu olarak eklenir ve uygulamanız genelinde tutarlı çeviri yapılmasını sağlar. `{"source":"API","target":"API"}` gibi tam bir girdi, aynı zamanda kaynakla birebir aynı olan bu değerin tamamını katı çevrilmemiş değer bulgularından muaf tutar; yalnızca `API` içeren daha uzun bir değeri muaf tutmaz.

<!-- internationalizer:unit markdown:translation-memory -->
## Çeviri belleği

Çeviri belleği bir JSONL dosyası olarak saklanır (satır başına bir JSON kaydı). Her kayıt şunları içerir:

- Paket, anahtar, kaynak değer, çevrilmiş değer ve kurallı hedef yerel ayar
- Kaynak, stil kılavuzu, sözlük, istem sözleşmesi ve birleşik ilke karmaları
- Çeviriyi üreten sağlayıcı ve model
- Bir zaman damgası

Sonraki çalıştırmalarda, aynı kaynak ve ilke karmalarına sahip dizeler LLM çağrılmadan önbellekten sunulur. Varsayılan yol yoksayılan `.internationalizer/` dizini altındadır, bu nedenle yerel bir önbellek olarak kalır. Projeniz çeviri belleğini kasıtlı olarak paylaşıyorsa `tm_path` değerini izlenen bir konuma ayarlayın. İncelenebilir `.internationalizer.lock` bildirimi ayrı olarak sürümlendirilir.

<!-- internationalizer:unit markdown:supported-formats -->
## Desteklenen biçimler

| Biçim | Uzantılar | Mod |
| --- | --- | --- |
| JSON | `.json` | Anahtar-değer (iç içe, nokta gösterimiyle düzleştirilmiş) |
| YAML | `.yml`, `.yaml` | Anahtar-değer (yorumları ve sıralamayı korur) |
| Markdown | `.md`, `.mdx` | Başlangıç metni (preamble) ve H2 düzeyinde bölümler |

Markdown hedefleri, H2 bölümlerinden önce görünmeyen `internationalizer:unit` yorumları içerir. Bu kararlı işaretleyiciler, Internationalizer'ın ilişkisiz bölümleri yeniden çevirmeden tek bir kaynak bölümü eklemesine, taşımasına veya düzenlemesine olanak tanır. İşaretsiz mevcut belgeler, bir sonraki başarılı güncellemelerinde bu işaretleyicileri alır.

<!-- internationalizer:unit markdown:project-type-detection -->
## Proje türü algılama

`internationalizer detect` şunları denetleyerek i18n kurulumunuzu tanımlar:

- react-i18next, next-intl veya vue-i18n için `package.json` bağımlılıkları
- Yaygın yerel ayar desenleriyle eşleşen dizin yapıları
- Dosya uzantıları ve adlandırma kuralları

<!-- internationalizer:unit markdown:architecture -->
## Mimari

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
## Alternatiflerle karşılaştırma

| Özellik | Internationalizer | i18next | Crowdin | Genel LLM |
| --- | --- | --- | --- | --- |
| LLM destekli çeviri | Evet | Hayır | Kısmen | Evet |
| Dile özel stil kılavuzları | Evet | Hayır | Hayır | Hayır |
| Sözlük zorunluluğu | Evet | Hayır | Evet | Hayır |
| Çeviri belleği | Evet | Hayır | Evet | Hayır |
| CLI / yerel yürütme | Evet | Yok | Hayır | Manuel |
| Git dostu dosyalar | Evet | Evet | Kısmen | Manuel |
| SaaS bağımlılığı yok | Evet | Evet | Hayır | Değişir |
| Açık kaynak (AGPL-3.0) | Evet | Evet | Hayır | Değişir |

<!-- internationalizer:unit markdown:license -->
## Lisans

[AGPL-3.0](../../LICENSE)

Bağımlılık bildirimleri için [THIRD_PARTY_NOTICES.md](../../THIRD_PARTY_NOTICES.md) dosyasına bakın.

<!-- internationalizer:unit markdown:contributing -->
## Katkıda bulunma

Geliştirme kurulumu ve yönergeler için [CONTRIBUTING.md](../../CONTRIBUTING.md) dosyasına bakın. Tüm katkılar için DCO onayı (sign-off) gereklidir.
