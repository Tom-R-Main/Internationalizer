> [English (original)](../../README.md)

<p align="center">
  <img src="../../assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

專為軟體專案打造的 AI 原生國際化管線。使用 LLM 來翻譯、驗證與管理 i18n 檔案。

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](../../LICENSE)

<p align="center">
<a href="ar.md">العربية</a> · <a href="bn.md">বাংলা</a> · <a href="cs.md">Čeština</a> · <a href="da.md">Dansk</a> · <a href="de.md">Deutsch</a> · <a href="el.md">Ελληνικά</a> · <a href="es.md">Español</a> · <a href="fi.md">Suomi</a> · <a href="fr.md">Français</a> · <a href="he.md">עברית</a> · <a href="hi.md">हिन्दी</a> · <a href="id.md">Indonesia</a> · <a href="it.md">Italiano</a> · <a href="ja.md">日本語</a> · <a href="ko.md">한국어</a> · <a href="ms.md">Bahasa Melayu</a><br><a href="nl.md">Nederlands</a> · <a href="pa.md">ਪੰਜਾਬੀ</a> · <a href="pl.md">Polski</a> · <a href="pt-BR.md">Português</a> · <a href="ro.md">Română</a> · <a href="ru.md">Русский</a> · <a href="sv.md">Svenska</a> · <a href="te.md">తెలుగు</a> · <a href="th.md">ไทย</a> · <a href="tr.md">Türkçe</a> · <a href="uk.md">Українська</a> · <a href="vi.md">Tiếng Việt</a> · <a href="yue.md">粵語</a> · <a href="zh-CN.md">简体中文</a> · <a href="zh-TW.md">繁體中文</a>
</p>

---

<!-- internationalizer:unit markdown:why-internationalizer -->
## 為什麼選擇 Internationalizer？

大多數 i18n 工具不是執行階段函式庫（i18next、react-intl），就是鍵值管理 SaaS 平台（Crowdin、Lokalise）。它們都無法妥善解決核心的翻譯問題：

- **人工翻譯**在語言數量超過少數幾種後便無法有效擴展
- **機器翻譯 API**（Google Translate、DeepL）會忽略您的術語、語氣與 UI 慣例
- **通用 LLM 翻譯**效果較好，但若缺乏詞彙表與風格指南，產出的結果容易前後不一

Internationalizer 截然不同。它是一個結合了 LLM 翻譯與下列特點的 **CLI 管線**：

- **各語言專屬詞彙表** — 在整個應用程式中強制維持一致的術語
- **各語言專屬風格指南** — 控制語氣、正式程度、複數規則與文字排版
- **翻譯記憶庫** — 跳過未變更的字串，節省 API 呼叫費用
- **確定性驗證** — 在發布前找出缺少或多餘的鍵、受保護結構漂移、詞彙表問題，以及複數或 ICU 錯誤

<!-- internationalizer:unit markdown:installation -->
## 安裝

從 npm 安裝：

```bash
npm install -g internationalizer
```

或者無需全域安裝即可執行：

```bash
npx internationalizer --help
```

npm 套件會透過平台專屬的選用相依性，從 npm 安裝相對應的預先建置二進位檔。

使用 Go 安裝：

```bash
go install github.com/Tom-R-Main/Internationalizer/cmd/internationalizer@latest
```

或從原始碼建置：

```bash
git clone https://github.com/Tom-R-Main/Internationalizer.git
cd Internationalizer
go build -o internationalizer ./cmd/internationalizer
```

<!-- internationalizer:unit markdown:npm-packages -->
## npm 套件

- Git 標籤與 npm 套件版本必須相符，例如 `v0.1.0` 與 `0.1.0`
- 根目錄的 `internationalizer` 套件相依於平台套件，例如 `internationalizer-darwin-arm64`
- 支援的 npm 目標平台：macOS arm64/x64、Linux arm64/x64、Windows x64
- CI 發布需要名為 `NPM_TOKEN` 的 GitHub secret

<!-- internationalizer:unit markdown:quick-start -->
## 快速入門

1. 在您的專案根目錄建立設定檔：

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

2. 設定您的 API 金鑰：

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

3. 預覽將翻譯的內容：

```bash
internationalizer translate --dry-run
```

4. 執行翻譯：

```bash
internationalizer translate
```

5. 驗證所有語言地區：

```bash
internationalizer validate
```

<!-- internationalizer:unit markdown:commands -->
## 指令

### `translate`

尋找缺少或過期的鍵，並透過 LLM 進行翻譯。

```bash
internationalizer translate                    # translate all locales
internationalizer translate -l fr              # translate French only
internationalizer translate --dry-run          # preview without API calls
internationalizer translate --adopt-existing   # baseline existing translations without API calls
internationalizer translate --refresh-policy   # refresh prompt/style/model-stale entries
internationalizer translate --batch-size 20    # smaller batches
internationalizer translate --concurrency 2    # fewer parallel calls
```

翻譯狀態會獨立回報缺少、來源過期、政策過期、最新以及手動編輯等情況，因此手動編輯不會掩蓋來源或政策的變更。政策過期的值會被回報，但僅在附加 `--refresh-policy` 時才會重新翻譯。手動編輯的值絕不會被自動覆寫。當初次將資訊清單導入已審核的翻譯，或明確要將審核後的手動編輯採納為新基準時，請使用 `--adopt-existing`。

### `validate`

依據來源套件檢查所有語言地區檔案。預設驗證會檢查結構涵蓋率（現存必要目標鍵的百分比），將多餘的鍵回報為警告，並在缺少鍵、插值不符或 ICU MessageFormat 結構無效時判定為失敗。

```bash
internationalizer validate                     # human-readable output
internationalizer validate --json              # machine-readable JSON
internationalizer validate -q                  # exit code only
internationalizer validate --strict             # enforce translation quality rules
internationalizer validate --require-state      # require current manifest provenance
```

`--strict` 也會回報翻譯涵蓋率。語言內容若與其來源完全相同，除非詞彙表明確包含該完整值且來源與目標相同的精確項目，否則將視為未翻譯；驗證會遵循 `ignore_case`，但較長值中僅嵌入詞彙表術語並不會獲得豁免。嚴格模式會在出現多餘的鍵、與來源相同的值、插值／HTML／程式碼／Markdown 連結結構變更、詞彙表違規以及已設定的複數形式缺失時判定為失敗。

`--require-state` 會對照 `.internationalizer.lock` 驗證每個目標。當鍵未受追蹤，或其記錄的來源、翻譯政策或目標雜湊已過期時，驗證將判定為失敗。此旗標可與 `--strict` 搭配使用。

人類可讀與 JSON 報告均使用穩定的檢查結果代碼：

| 代碼 | 意義 |
| --- | --- |
| `missing_key` / `extra_key` | 來源與目標鍵集合不符 |
| `blank_translation` | 非空的來源在嚴格模式下對應了空的目標值 |
| `source_identical` | 嚴格模式下語言值仍未翻譯 |
| `protected_structure_mismatch` | 插值、HTML、程式碼或連結結構已變更 |
| `glossary_violation` | 找不到已核准的目標術語或變體 |
| `plural_form_missing` | 缺少已設定的語言地區複數形式 |
| `icu_message_syntax` | 來源或目標 ICU 訊息語法格式錯誤 |
| `icu_argument_mismatch` | ICU 引數名稱、類型或格式器樣式不符 |
| `icu_selector_mismatch` | 選擇器不符，或複數類別對目標語言地區無效 |
| `untracked` | 資訊清單中不存在該目標的記錄 |
| `source_stale` | 記錄翻譯後來源內容已變更 |
| `policy_stale` | 產生的提示詞或模型設定已變更 |
| `target_modified` | 目標內容與資訊清單記錄不符 |

### `detect`

自動偵測 i18n 框架並建議設定。

```bash
internationalizer detect
```

支援：react-i18next、next-intl、vue-i18n、原生 JSON、Markdown 文件。

### `glossary`

管理在翻譯期間強制套用的各語言專屬詞彙表術語。

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

管理翻譯記憶庫（先前翻譯字串的 JSONL 快取）。

```bash
internationalizer tm stats                     # show record counts
internationalizer tm export                    # dump as JSON
internationalizer tm clear --force             # delete all records
```

<!-- internationalizer:unit markdown:configuration-reference -->
## 設定參考

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

語言地區識別碼必須是格式正確的 BCP 47 標籤，例如 `fr`、`pt-BR` 或 `sr-Latn-RS`。規範等價的目標語言地區會視為重複項目而遭拒絕，且各語言地區專屬的供應商覆寫會比對規範等價拼寫。在上述範例中，未設定覆寫的語言地區（包括日文）會繼承全域 Gemini 設定。

ICU MessageFormat 值會進行結構化解析。系統支援簡易引數、`select`、`plural`、`selectordinal`、`number`、`date` 以及 `time`，包括巢狀訊息、複數偏移量、精確數字選擇器與 `#`。驗證會檢查語法、引數類型與格式器樣式、複數偏移量、select 分支一致性，以及目標語言地區 CLDR 複數類別。破壞這些不變條件的供應商輸出會在寫入語言地區檔案或翻譯記憶庫記錄之前遭到拒絕。

在 `i18next-v4` 下，已識別的來源複數群組會在翻譯期間擴充為目標語言地區的 CLDR 類別。僅目標端存在的類別會使用來源群組的 `_other` 值作為其翻譯範本。嚴格驗證要求具備這些目標類別；若目標語言地區不使用某些僅來源端具備的類別，則這些類別為選填。

<!-- internationalizer:unit markdown:style-guides -->
## 風格指南

風格指南是會被注入到 LLM 翻譯提示詞中的 Markdown 檔案。它們控制語氣、正式程度、文字排版與其他特定語言慣例。

```
style-guides/
  _conventions.md    # shared rules for all languages
  fr.md              # French-specific rules
  ja.md              # Japanese-specific rules
  ar.md              # Arabic-specific rules
```

### 共用慣例（`_conventions.md`）

定義適用於所有語言的規則：插值語法、HTML 保留、字串類型慣例（按鈕與標籤與錯誤訊息等）。

### 各語言指南（`{locale}.md`）

定義特定語言規則：正式程度語域（tu 與 vous）、標點符號（法文引號、倒問號）、複數形式、日期／數字格式，以及術語詞彙表。

風格指南屬於長期有效的政策輸入，而非產生的輸出。Internationalizer 會讀取它們但絕不重寫。其內容雜湊與詞彙表及提示詞協定分開計算，因此應用程式程式碼的變更不會使翻譯過期。編輯指南會特意將該語言地區標記為需進行政策審查；變更內部提示詞用詞則不會，除非提示詞協定版本同時變更。

實作範例請參閱 [`examples/react-app/style-guides/`](../../examples/react-app/style-guides/)。

<!-- internationalizer:unit markdown:glossary-format -->
## 詞彙表格式

詞彙表檔案是儲存在 `{glossary_dir}/{locale}.json` 中的 JSON 陣列：

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

`variants` 列出其他已核准的目標形式。`enforcement` 可為 `error`、`warning`，或省略以採用預設的 error 行為。術語會以術語表形式注入至 LLM 提示詞中，確保整個應用程式的翻譯一致。完全相同的精確項目（如 `{"source":"API","target":"API"}`）亦會讓該完整同源值免於嚴格模式的未翻譯值檢查結果；但若僅是較長的值中包含 `API`，則不符合豁免條件。

<!-- internationalizer:unit markdown:translation-memory -->
## 翻譯記憶庫

翻譯記憶庫儲存為 JSONL 檔案（每行一個 JSON 記錄）。每個記錄包含：

- 套件、鍵、來源值、翻譯後的值，以及標準目標語言地區
- 來源、風格指南、詞彙表、提示詞協定與合併政策雜湊
- 產生翻譯的供應商與模型
- 時間戳記

在後續執行時，具有相同來源與政策雜湊的字串會直接從快取提供，無需呼叫 LLM。預設路徑位於被忽略的 `.internationalizer/` 目錄下，因此它保持為本機快取。如果您的專案有意共用翻譯記憶庫，請將 `tm_path` 設為受版本控制的位置。可供審查的 `.internationalizer.lock` 資訊清單則分開進行版本控制。

<!-- internationalizer:unit markdown:supported-formats -->
## 支援的格式

| 格式 | 副檔名 | 模式 |
|--------|-----------|------|
| JSON | `.json` | 鍵值對（巢狀、點記法扁平化） |
| YAML | `.yml`, `.yaml` | 鍵值對（保留註解與排序） |
| Markdown | `.md`, `.mdx` | 前言與 H2 層級章節 |

Markdown 目標在 H2 章節前包含隱形的 `internationalizer:unit` 註解。這些穩定的標記讓 Internationalizer 能夠新增、移動或編輯單一來源章節，而無需重新翻譯不相關的章節。現有未標記的文件會在下次成功更新時獲得標記。

<!-- internationalizer:unit markdown:project-type-detection -->
## 專案類型偵測

`internationalizer detect` 透過檢查下列項目來識別您的 i18n 設定：

- `package.json` 中 react-i18next、next-intl 或 vue-i18n 的相依性
- 符合常見語言地區模式的目錄結構
- 檔案副檔名與命名慣例

<!-- internationalizer:unit markdown:architecture -->
## 架構

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
## 與替代方案的比較

| 功能 | Internationalizer | i18next | Crowdin | 通用 LLM |
|---------|------------------|---------|---------|-------------|
| LLM 驅動翻譯 | 是 | 否 | 部分支援 | 是 |
| 各語言專屬風格指南 | 是 | 否 | 否 | 否 |
| 詞彙表強制執行 | 是 | 否 | 是 | 否 |
| 翻譯記憶庫 | 是 | 否 | 是 | 否 |
| CLI／本機執行 | 是 | 不適用 | 否 | 手動 |
| Git 友善檔案 | 是 | 是 | 部分支援 | 手動 |
| 無 SaaS 相依性 | 是 | 是 | 否 | 視情況而定 |
| 開放原始碼（AGPL-3.0） | 是 | 是 | 否 | 視情況而定 |

<!-- internationalizer:unit markdown:license -->
## 授權條款

[AGPL-3.0](../../LICENSE)

相依性聲明請參閱 [THIRD_PARTY_NOTICES.md](../../THIRD_PARTY_NOTICES.md)。

<!-- internationalizer:unit markdown:contributing -->
## 貢獻指南

開發設定與指導原則請參閱 [CONTRIBUTING.md](../../CONTRIBUTING.md)。所有貢獻均需具備 DCO 簽名。
