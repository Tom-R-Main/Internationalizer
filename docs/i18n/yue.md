> [English (original)](../../README.md)

<p align="center">
  <img src="../../assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

專為軟件項目而設嘅 AI 原生國際化 pipeline。使用 LLM 翻譯、驗證同管理 i18n 檔案。

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](../../LICENSE)

<p align="center">
<a href="ar.md">العربية</a> · <a href="bn.md">বাংলা</a> · <a href="cs.md">Čeština</a> · <a href="da.md">Dansk</a> · <a href="de.md">Deutsch</a> · <a href="el.md">Ελληνικά</a> · <a href="es.md">Español</a> · <a href="fi.md">Suomi</a> · <a href="fr.md">Français</a> · <a href="he.md">עברית</a> · <a href="hi.md">हिन्दी</a> · <a href="id.md">Indonesia</a> · <a href="it.md">Italiano</a> · <a href="ja.md">日本語</a> · <a href="ko.md">한국어</a> · <a href="ms.md">Bahasa Melayu</a><br><a href="nl.md">Nederlands</a> · <a href="pa.md">ਪੰਜਾਬੀ</a> · <a href="pl.md">Polski</a> · <a href="pt-BR.md">Português</a> · <a href="ro.md">Română</a> · <a href="ru.md">Русский</a> · <a href="sv.md">Svenska</a> · <a href="te.md">తెలుగు</a> · <a href="th.md">ไทย</a> · <a href="tr.md">Türkçe</a> · <a href="uk.md">Українська</a> · <a href="vi.md">Tiếng Việt</a> · <a href="yue.md">粵語</a> · <a href="zh-CN.md">简体中文</a> · <a href="zh-TW.md">繁體中文</a>
</p>

---

<!-- internationalizer:unit markdown:why-internationalizer -->
## 點解要揀 Internationalizer？

大部分 i18n 工具一係 runtime 函式庫（i18next、react-intl），一係 key 管理 SaaS 平台（Crowdin、Lokalise）。但佢哋都無法好好解決實際嘅翻譯問題：

- **人手翻譯** 喺語言數量多咗嗰陣根本應付唔嚟
- **機器翻譯 API**（Google Translate、DeepL）會忽略你嘅術語、語氣同 UI 慣例
- **通用 LLM 翻譯** 雖然好啲，但如果冇詞彙表同風格指南，翻譯出嚟嘅結果就會好唔一致

Internationalizer 截然不同。佢係一個結合 LLM 翻譯同以下功能嘅 **CLI pipeline**：

- **各語言專屬詞彙表** — 喺成個應用程式中強制維持一致嘅術語
- **各語言專屬風格指南** — 控制語氣、正式程度、眾數同排版規則
- **翻譯記憶庫** — 跳過冇修改過嘅字串，慳返 API 呼叫費用
- **確定性驗證** — 喺發佈前即時攔截缺少或多餘嘅 key、受保護結構跑位、詞彙表違規，以及眾數或 ICU 錯誤

<!-- internationalizer:unit markdown:installation -->
## 安裝

透過 npm 安裝：

```bash
npm install -g internationalizer
```

或者唔使全域安裝直接執行：

```bash
npx internationalizer --help
```

npm 套件會透過特定平台嘅 optional dependencies，從 npm 安裝對應嘅預先編譯二進制檔案。

透過 Go 安裝：

```bash
go install github.com/Tom-R-Main/Internationalizer/cmd/internationalizer@latest
```

或者從原始碼編譯：

```bash
git clone https://github.com/Tom-R-Main/Internationalizer.git
cd Internationalizer
go build -o internationalizer ./cmd/internationalizer
```

<!-- internationalizer:unit markdown:npm-packages -->
## npm 套件

- Git tag 同 npm 套件版本必須一致，例如 `v0.1.0` 同 `0.1.0`
- 根目錄嘅 `internationalizer` 套件依賴特定平台套件，例如 `internationalizer-darwin-arm64`
- 支援嘅 npm 目標平台：macOS arm64/x64、Linux arm64/x64、Windows x64
- CI 發佈流程需要一個名為 `NPM_TOKEN` 嘅 GitHub secret

<!-- internationalizer:unit markdown:quick-start -->
## 快速上手

1. 喺你嘅項目根目錄建立設定檔：

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

2. 設定你嘅 API key：

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

3. 預覽將會翻譯嘅內容：

```bash
internationalizer translate --dry-run
```

4. 執行翻譯：

```bash
internationalizer translate
```

5. 驗證所有語言地區檔案：

```bash
internationalizer validate
```

<!-- internationalizer:unit markdown:commands -->
## 指令

### `translate`

搵出缺少或者過期嘅 key，並透過 LLM 進行翻譯。

```bash
internationalizer translate                    # translate all locales
internationalizer translate -l fr              # translate French only
internationalizer translate --dry-run          # preview without API calls
internationalizer translate --adopt-existing   # baseline existing translations without API calls
internationalizer translate --refresh-policy   # refresh prompt/style/model-stale entries
internationalizer translate --batch-size 20    # smaller batches
internationalizer translate --concurrency 2    # fewer parallel calls
```

翻譯狀態會獨立報告缺少、來源過期、策略過期、最新同手動修改過嘅狀況，所以手動編輯唔會遮蔽來源變更或策略更新。策略過期嘅項目會被報告出嚟，但只有加咗 `--refresh-policy` 嗰陣先會重新翻譯。手動修改過嘅值永遠唔會被自動覆蓋。當第一次為已經過審查嘅翻譯引入 manifest，或者要明確將審查過嘅手動編輯採納為新嘅基準嗰陣，請使用 `--adopt-existing`。

### `validate`

對照來源 bundle 檢查所有語言地區檔案。預設驗證會檢查結構覆蓋率（所需目標 key 存在嘅百分比），將多出嘅 key 列為警告，並喺發現缺少 key、插值不匹配或無效嘅 ICU MessageFormat 結構時報錯。

```bash
internationalizer validate                     # human-readable output
internationalizer validate --json              # machine-readable JSON
internationalizer validate -q                  # exit code only
internationalizer validate --strict             # enforce translation quality rules
internationalizer validate --require-state      # require current manifest provenance
```

`--strict` 亦會報告已翻譯覆蓋率。除咗詞彙表明確包含完全相同來源、相同目標嘅整句完整詞條之外，同來源完全一樣嘅語言文字會被視為未翻譯；雖然系統支援 `ignore_case`，但只係喺較長文字入面嵌入詞彙表詞條並唔能夠獲得豁免。嚴格模式會喺出現多餘 key、與來源完全一致嘅值、插值/HTML/代碼/Markdown 連結結構變更、詞彙表違規同設定嘅眾數形式缺少時報錯。

`--require-state` 會對照 `.internationalizer.lock` 驗證每個目標。當 key 未被追蹤，或者記錄嘅來源、翻譯策略或目標 hash 過期嗰陣，驗證就會失敗。呢個參數可以同 `--strict` 一齊用。

人類可讀報告同 JSON 報告都使用穩定嘅問題代碼：

| 代碼 | 意思 |
| --- | --- |
| `missing_key` / `extra_key` | 來源同目標嘅 key 組合不一致 |
| `blank_translation` | 非空白嘅來源喺嚴格模式下對應咗空白嘅目標 |
| `source_identical` | 嚴格模式下語言文字依然保持未翻譯狀態 |
| `protected_structure_mismatch` | 插值、HTML、代碼或連結結構有所改變 |
| `glossary_violation` | 搵唔到認可嘅目標詞條或變體 |
| `plural_form_missing` | 設定嘅語言地區眾數形式缺少咗 |
| `icu_message_syntax` | 來源或目標嘅 ICU 訊息格式有誤 |
| `icu_argument_mismatch` | ICU 引數名稱、類型或格式化樣式不一致 |
| `icu_selector_mismatch` | 選擇器不一致，或者眾數分類對該目標語言地區無效 |
| `untracked` | 該目標喺 manifest 中完全冇記錄 |
| `source_stale` | 來源內容喺記錄嘅翻譯之後發生咗變更 |
| `policy_stale` | 產生嘅 prompt 或模型設定發生咗變更 |
| `target_modified` | 目標內容同 manifest 記錄唔同 |

### `detect`

自動偵測 i18n 框架並建議設定。

```bash
internationalizer detect
```

支援：react-i18next、next-intl、vue-i18n、一般 JSON、Markdown 文件。

### `glossary`

管理翻譯過程中強制執行嘅各語言詞彙表詞條。

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

管理翻譯記憶庫（儲存之前翻譯字串嘅 JSONL 快取）。

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

語言地區識別碼必須係格式正確嘅 BCP 47 標籤，例如 `fr`、`pt-BR` 或 `sr-Latn-RS`。規範等價嘅目標語言地區會被當成重複項而拒絕，而特定語言地區嘅供應商覆寫亦會以規範等價嘅寫法進行比對。喺上面嘅範例入面，冇特別設定覆寫嘅語言地區（包括日文）會直接繼承全域嘅 Gemini 設定。

系統會對 ICU MessageFormat 嘅值進行結構化解析。支援簡單引數、`select`、`plural`、`selectordinal`、`number`、`date` 同 `time`，包含巢狀訊息、眾數偏移（plural offsets）、精確數字選擇器同 `#`。驗證過程會檢查語法、引數類型同格式化樣式、眾數偏移、select 分支一致性，以及目標語言地區嘅 CLDR 眾數分類。只要供應商輸出破壞咗呢啲不變量，喺寫入語言地區檔案或翻譯記憶庫記錄之前就會被拒絕。

配合 `i18next-v4`，已識別嘅來源眾數家族喺翻譯期間會展開為目標語言地區嘅 CLDR 分類。目標特有嘅分類會使用來源家族嘅 `_other` 數值作為翻譯範本。嚴格驗證要求必須存在呢啲目標分類；至於來源特有嘅分類，對於唔使用佢哋嘅目標語言地區嚟講係選填嘅。

<!-- internationalizer:unit markdown:style-guides -->
## 風格指南

風格指南係 Markdown 檔案，會被注入到 LLM 翻譯 prompt 入面。佢哋負責控制語氣、正式程度、排版以及其他特定語言嘅慣例。

```
style-guides/
  _conventions.md    # shared rules for all languages
  fr.md              # French-specific rules
  ja.md              # Japanese-specific rules
  ar.md              # Arabic-specific rules
```

### 共用慣例 (`_conventions.md`)

定義適用於所有語言嘅規則：插值語法、保留 HTML、字串類型慣例（按鈕 vs 標籤 vs 錯誤訊息）等。

### 各語言專屬指南 (`{locale}.md`)

定義特定語言嘅規則：敬語程度（tu vs. vous）、標點符號（法文引號、倒問號）、眾數形式、日期/數字格式同術語詞彙表。

風格指南係持久嘅策略輸入，而唔係產生出嚟嘅輸出。Internationalizer 只會讀取佢哋，永遠唔會重寫。佢哋嘅內容 hash 係獨立於詞彙表同 prompt 協定之外分開計算嘅，所以應用程式代碼變更唔會導致翻譯過期。修改風格指南會特登將嗰個語言地區標記為需要策略審查；修改內部 prompt 字眼則唔會，除非 prompt 協定版本同時改變。

完整運作範例請參考 [`examples/react-app/style-guides/`](../../examples/react-app/style-guides/)。

<!-- internationalizer:unit markdown:glossary-format -->
## 詞彙表格式

詞彙表檔案係儲存喺 `{glossary_dir}/{locale}.json` 嘅 JSON 陣列：

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

`variants` 列出其他認可嘅目標形式。`enforcement` 可以係 `error`、`warning`，或者省略以使用預設嘅 error 行為。詞條會作為術語表注入到 LLM prompt 入面，確保應用程式嘅翻譯一致。好似 `{"source":"API","target":"API"}` 咁樣嘅完全相符項目，亦可以豁免該完整與來源相同嘅值喺嚴格模式下被列為未翻譯問題；但如果只係較長文字入面咁啱包含 `API`，就唔會獲得豁免。

<!-- internationalizer:unit markdown:translation-memory -->
## 翻譯記憶庫

翻譯記憶庫儲存為 JSONL 檔案（每行一筆 JSON 記錄）。每筆記錄包含：

- bundle、key、來源值、翻譯值同規範目標語言地區
- 來源、風格指南、詞彙表、prompt 協定同合併策略 hash
- 產生該翻譯嘅供應商同模型
- 時間戳記

喺後續嘅執行中，只要來源同策略 hash 相同嘅字串，就可以直接從快取讀取，唔需要呼叫 LLM。預設路徑位於被忽略嘅 `.internationalizer/` 目錄入面，所以佢只係一個本機快取。如果你嘅項目打算共享翻譯記憶庫，請將 `tm_path` 設定為納入版本控制嘅位置。可供審查嘅 `.internationalizer.lock` manifest 則係分開進行版本控制嘅。

<!-- internationalizer:unit markdown:supported-formats -->
## 支援嘅格式

| 格式 | 副檔名 | 模式 |
| --- | --- | --- |
| JSON | `.json` | Key-value（巢狀，點標記法扁平化） |
| YAML | `.yml`, `.yaml` | Key-value（保留註解同順序） |
| Markdown | `.md`, `.mdx` | 前言（preamble）同 H2 級段落 |

Markdown 目標檔案喺 H2 段落前會包含睇唔見嘅 `internationalizer:unit` 註解。呢啲穩定標記令 Internationalizer 可以新增、移動或編輯單一來源段落，而唔需要重新翻譯其他無關嘅段落。現有未加入標記嘅文件會喺下一次成功更新嗰陣自動補上標記。

<!-- internationalizer:unit markdown:project-type-detection -->
## 項目類型偵測

`internationalizer detect` 會透過檢查以下項目嚟識別你嘅 i18n 設定：

- `package.json` 裡面嘅 react-i18next、next-intl 或 vue-i18n 依賴套件
- 符合常見語言地區模式嘅目錄結構
- 檔案副檔名同命名慣例

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
## 替代方案比較

| 功能 | Internationalizer | i18next | Crowdin | 通用 LLM |
| --- | --- | --- | --- | --- |
| LLM 驅動翻譯 | 是 | 否 | 部分 | 是 |
| 各語言專屬風格指南 | 是 | 否 | 否 | 否 |
| 強制執行詞彙表 | 是 | 否 | 是 | 否 |
| 翻譯記憶庫 | 是 | 否 | 是 | 否 |
| CLI / 本地執行 | 是 | 不適用 | 否 | 手動 |
| Git 友好檔案 | 是 | 是 | 部分 | 手動 |
| 無 SaaS 依賴 | 是 | 是 | 否 | 視情況而定 |
| 開源 (AGPL-3.0) | 是 | 是 | 否 | 視情況而定 |

<!-- internationalizer:unit markdown:license -->
## 授權條款

[AGPL-3.0](../../LICENSE)

請參閱 [THIRD_PARTY_NOTICES.md](../../THIRD_PARTY_NOTICES.md) 查看依賴套件聲明。

<!-- internationalizer:unit markdown:contributing -->
## 參與貢獻

請參閱 [CONTRIBUTING.md](../../CONTRIBUTING.md) 了解開發環境設定同指引。所有貢獻都需要符合 DCO 簽署要求。
