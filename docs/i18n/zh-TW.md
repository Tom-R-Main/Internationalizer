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

## 為什麼選擇 Internationalizer？

大多數的 i18n 工具不是執行階段函式庫（i18next、react-intl），就是翻譯鍵管理 SaaS 平台（Crowdin、Lokalise）。它們都無法妥善解決核心的翻譯問題：

- **手動翻譯**在語言數量增加後難以擴展
- **機器翻譯 API**（Google Translate、DeepL）會忽略您的專屬術語、語氣與 UI 慣例
- **通用 LLM 翻譯**效果較佳，但若缺乏詞彙表與風格指南，產出的結果容易前後不一

Internationalizer 截然不同。它是一個結合了 LLM 翻譯與下列功能的 **CLI 管線**：

- **各語言專屬詞彙表** — 確保整個應用程式中的術語維持一致
- **各語言專屬風格指南** — 精確控制語氣、正式程度、複數規則與文字排版
- **翻譯記憶庫** — 自動跳過未變更的字串，節省 API 呼叫費用
- **確定性驗證** — 在發布前找出翻譯鍵遺漏或多餘、受保護結構變更、詞彙表違規，以及複數或 ICU 錯誤

## 安裝

透過 npm 安裝：

```bash
npm install -g internationalizer
```

或者無需全域安裝，直接執行：

```bash
npx internationalizer --help
```

npm 套件會藉由平台專屬的選用相依性，從 npm 安裝相對應的預先建置二進位檔。

透過 Go 安裝：

```bash
go install github.com/Tom-R-Main/Internationalizer/cmd/internationalizer@latest
```

或從原始碼建置：

```bash
git clone https://github.com/Tom-R-Main/Internationalizer.git
cd Internationalizer
go build -o internationalizer ./cmd/internationalizer
```

## npm 套件

- Git 標籤與 npm 套件版本必須完全一致，例如 `v0.1.0` 與 `0.1.0`
- 根目錄的 `internationalizer` 套件相依於各平台套件，例如 `internationalizer-darwin-arm64`
- 支援的 npm 目標平台：macOS arm64/x64、Linux arm64/x64、Windows x64
- CI 發布流程需要名為 `NPM_TOKEN` 的 GitHub secret

## 快速入門

1. 在專案根目錄建立設定檔：

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

3. 預覽將要翻譯的內容：

```bash
internationalizer translate --dry-run
```

4. 執行翻譯：

```bash
internationalizer translate
```

5. 驗證所有語系：

```bash
internationalizer validate
```

## 指令

### `translate`

找出遺漏或已過期的翻譯鍵並透過 LLM 進行翻譯。

```bash
internationalizer translate                    # 翻譯所有語系
internationalizer translate -l fr              # 僅翻譯法文
internationalizer translate --dry-run          # 預覽翻譯內容，不呼叫 API
internationalizer translate --adopt-existing   # 將現有翻譯設為基準，不呼叫 API
internationalizer translate --refresh-policy   # 重新整理因提示詞／風格／模型變更而過期的項目
internationalizer translate --batch-size 20    # 縮小每批次翻譯量
internationalizer translate --concurrency 2    # 降低平行呼叫數
```

翻譯狀態會獨立回報遺漏、來源過期、政策過期、最新與手動編輯等情況，因此手動編輯不會掩蓋來源或政策的變更。政策過期的內容會被列出回報，但僅在附加 `--refresh-policy` 旗標時才會重新翻譯。系統絕不會自動覆寫手動編輯的內容。初次將資訊清單（manifest）導入已審核的翻譯，或明確要將審核後的手動編輯採納為新基準時，請使用 `--adopt-existing`。

### `validate`

將所有語系檔案與來源套件進行比對。預設驗證會計算必要目標翻譯鍵的涵蓋率，將多餘的翻譯鍵列為警告，並在翻譯鍵遺漏、插值不符或 ICU MessageFormat 結構無效時失敗。

```bash
internationalizer validate                     # 人類可讀的輸出
internationalizer validate --json              # 機器可讀的 JSON
internationalizer validate -q                  # 僅傳回結束代碼
internationalizer validate --strict             # 強制執行翻譯品質規則
internationalizer validate --require-state      # 要求資訊清單中的來源追溯狀態為最新
```

`--strict` 也會回報實際翻譯涵蓋率。語言內容若與來源完全相同，便會視為未翻譯；只有詞彙表針對完整值明確設定來源與譯文相同的精確項目時才會豁免。驗證會遵守 `ignore_case`，但較長內容中只包含某個詞彙表項目並不會獲得豁免。嚴格模式會在出現多餘翻譯鍵、與來源相同的值、插值／HTML／程式碼／Markdown 連結結構變更、詞彙表違規或缺少已設定的複數形式時失敗。

`--require-state` 會將每個目標值與 `.internationalizer.lock` 比對。翻譯鍵未受追蹤，或資訊清單記錄的來源、翻譯政策、目標雜湊已過期時，驗證都會失敗。此選項可與 `--strict` 同時使用。

人類可讀報告與 JSON 報告使用下列穩定的問題代碼：

| 代碼 | 含義 |
| --- | --- |
| `missing_key` / `extra_key` | 來源與目標的翻譯鍵集合不一致 |
| `blank_translation` | 非空來源所對應的嚴格模式譯文為空 |
| `source_identical` | 嚴格模式下語言內容仍與來源相同 |
| `protected_structure_mismatch` | 插值、HTML、程式碼或連結結構發生變更 |
| `glossary_violation` | 找不到獲准使用的目標術語或變體 |
| `plural_form_missing` | 缺少為該語言設定的複數形式 |
| `icu_message_syntax` | 來源或譯文的 ICU 訊息格式有誤 |
| `icu_argument_mismatch` | ICU 引數名稱、類型或格式樣式不一致 |
| `icu_selector_mismatch` | 選擇器不一致，或複數類別不適用於目標語言 |
| `untracked` | 資訊清單中沒有對應的目標記錄 |
| `source_stale` | 資訊清單記錄建立後來源內容發生變更 |
| `policy_stale` | 產生提示詞或模型設定發生變更 |
| `target_modified` | 目標內容與資訊清單記錄不一致 |

### `detect`

自動偵測 i18n 框架並提供建議設定。

```bash
internationalizer detect
```

支援：react-i18next、next-intl、vue-i18n、純 JSON、Markdown 文件。

### `glossary`

管理在翻譯期間強制套用的各語言專屬詞彙表術語。

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

管理翻譯記憶庫（儲存先前翻譯字串的 JSONL 快取）。

```bash
internationalizer tm stats                     # 顯示記錄總數
internationalizer tm export                    # 匯出為 JSON
internationalizer tm clear --force             # 刪除所有記錄
```

## 設定參考

```yaml
# .internationalizer.yml

# 來源語言（預設：en）
source_locale: en

# 要翻譯成的目標語言（必填）
target_locales: [fr, de, es, ja, yue, zh-CN, zh-TW, ar]

# 一或多個來源到目標的對應設定（必填）。
# {locale} 會替換為每個設定的目標語言代碼。
bundles:
  - id: app
    source: locales/en.json
    target: locales/{locale}.json
    format: json
  - id: docs
    source: README.md
    target: docs/i18n/{locale}.md
    format: markdown

# 向後相容設定：source_path 仍可將目標對應至同層檔案，
# 例如 locales/fr.json。新專案建議優先使用 bundles。
# source_path: locales/en.json

# LLM 供應商設定
llm:
  # 供應商："anthropic"、"openai"、"gemini" 或 "openrouter"（預設：gemini）
  provider: gemini

  # 各供應商的預設模型名稱：
  #   anthropic:  claude-opus-5
  #   openai:     gpt-5.6-luna（推理程度預設為 max）
  #   gemini:     gemini-3.8-flash
  #   openrouter: deepseek/deepseek-v4-pro-0813
  model: gemini-3.8-flash

  # 存放 API 金鑰的環境變數名稱
  api_key_env: GOOGLE_AI_STUDIO_API_KEY

  # 相容於 OpenAI 端點的基礎 URL（選填）
  # base_url: https://api.openai.com

  # OpenAI GPT-5 系列 Responses API 的推理程度
  #（OpenAI 供應商預設值：max）
  reasoning_effort: max

  # 個別目標語言的選用 LLM 設定。若覆寫時沿用全域供應商，
  # 未指定的項目將繼承全域設定；若改用其他供應商，
  # 未指定的項目則採用該供應商的預設值。
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

# 單次 LLM 呼叫處理的鍵數量（預設：40）
batch_size: 40

# 平行 LLM 呼叫數量（預設：4）
concurrency: 4

# 存放各語言風格指南 Markdown 檔案的目錄（預設：style-guides）
style_guides_dir: style-guides

# 存放各語言詞彙表 JSON 檔案的目錄（預設：glossary）
glossary_dir: glossary

# 翻譯記憶庫檔案路徑（預設：.internationalizer/tm.jsonl）
tm_path: .internationalizer/tm.jsonl

# 來源、政策、目標及來源追溯資訊的版本化狀態
#（預設：.internationalizer.lock；請將此檔案提交至版本控制）
manifest_path: .internationalizer.lock

# 選用的翻譯及嚴格驗證規則
validation:
  plural_style: i18next-v4 # 產生並驗證目標語言的複數形式
```

語言識別碼必須是格式正確的 BCP 47 標籤，例如 `fr`、`pt-BR` 或 `sr-Latn-RS`。正規化後等價的目標語言會被判定為重複項目；各語言專屬的供應商覆寫設定也會依正規化形式比對。上述範例中，日文等沒有個別覆寫設定的語言都會繼承全域 Gemini 設定。

ICU MessageFormat 內容會依結構解析。系統支援簡單引數、`select`、`plural`、`selectordinal`、`number`、`date` 與 `time`，以及巢狀訊息、複數偏移量、精確數字選擇器和 `#`。驗證會檢查語法、引數類型與格式樣式、複數偏移量、select 分支識別，以及目標語言適用的 CLDR 複數類別。若供應商回傳的內容破壞這些條件，系統會在寫入語系檔案或翻譯記憶庫前拒絕該結果。

啟用 `i18next-v4` 後，系統會在翻譯期間將辨識出的來源複數詞組擴充為目標語言需要的 CLDR 類別。只有目標語言需要的類別會使用來源詞組的 `_other` 值作為翻譯範本。嚴格驗證要求目標語言所需的類別齊全，但不強制保留目標語言不使用的來源語言專屬類別。

## 風格指南

風格指南是會直接注入到 LLM 翻譯提示詞中的 Markdown 檔案。可用來控制語氣、正式程度、文字排版與其他語言專屬的慣例。

```
style-guides/
  _conventions.md    # 適用於所有語言的共用規則
  fr.md              # 法文專屬規則
  ja.md              # 日文專屬規則
  ar.md              # 阿拉伯文專屬規則
```

### 共用慣例（`_conventions.md`）

定義通用於所有語言的規範：變數插值語法、HTML 標籤保留、字串類型慣例（按鈕、標籤與錯誤訊息）等。

### 各語言專屬指南（`{locale}.md`）

定義特定語言的細部規則：正式程度（例如 tu 與 vous）、標點符號規範（法文引號、倒問號）、複數形式、日期／數字格式，以及術語對照清單。

實際運作範例請參閱 [`examples/react-app/style-guides/`](../../examples/react-app/style-guides/)。

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

`variants` 用來列出其他獲准使用的譯法。`enforcement` 可設為 `error` 或 `warning`；省略時預設為 `error`。這些項目會以術語對照表的形式注入 LLM 提示詞，確保整個應用程式中的術語維持一致。對於 `{"source":"API","target":"API"}` 這類來源與譯文完全相同的精確項目，嚴格驗證不會將完整值判定為未翻譯；較長內容中只包含 `API` 則不會獲得豁免。

## 翻譯記憶庫

翻譯記憶庫儲存為 JSONL 格式（每行一筆 JSON 記錄）。每筆記錄包含：

- 套件、翻譯鍵、來源值、譯文及正規化後的目標語言
- 來源雜湊與翻譯政策雜湊
- 產生譯文的供應商與模型
- 時間戳記

後續執行時，來源雜湊與政策雜湊都相同的內容會直接從快取重複使用，不必呼叫 LLM。預設路徑位於已被忽略的 `.internationalizer/` 目錄中，因此它是本機快取。如需讓專案成員共用翻譯記憶庫，請將 `tm_path` 明確設為受版本控制的路徑。可供審查的 `.internationalizer.lock` 資訊清單則另行納入版本控制。

## 支援的格式

| 格式 | 副檔名 | 模式 |
|--------|-----------|------|
| JSON | `.json` | 鍵值對（支援巢狀、以點號表示法扁平化） |
| YAML | `.yml`、`.yaml` | 鍵值對（保留註解與欄位順序） |
| Markdown | `.md`、`.mdx` | 全文件翻譯 |

## 專案類型偵測

執行 `internationalizer detect` 時，會透過檢查下列項目來識別您的 i18n 架構：

- `package.json` 中對 react-i18next、next-intl 或 vue-i18n 的相依性
- 符合常見語系結構的目錄配置
- 檔案副檔名與命名慣例

## 架構

```
cmd/internationalizer/     CLI 進入點與指令定義
internal/
  config/                  YAML 設定載入與預設值處理
  detect/                  專案類型自動偵測
  formats/                 格式解析器（JSON、YAML、Markdown）
  glossary/                各語言專屬詞彙表管理
  llm/                     LLM 供應商介面與實作
    anthropic.go           Anthropic Claude 後端
    openai.go              OpenAI／相容端點後端
    gemini.go              Google Gemini（透過 AI Studio 後端）
                           OpenRouter 採用自訂 base_url 的 openai.go
  locale/                  BCP 47 語言識別與 CLDR 複數類別
  message/                 ICU MessageFormat 解析器與結構比對
  policy/                  穩定的翻譯政策雜湊
  state/                   受版本控制的翻譯資訊清單
  styleguide/              風格指南載入器
  tm/                      JSONL 翻譯記憶庫
  translate/               翻譯流程協調器
  validate/                語系驗證與差異比對
```

## 與替代方案的比較

| 功能 | Internationalizer | i18next | Crowdin | 通用 LLM |
|---------|------------------|---------|---------|-------------|
| LLM 驅動翻譯 | 是 | 否 | 部分支援 | 是 |
| 各語言專屬風格指南 | 是 | 否 | 否 | 否 |
| 強制套用詞彙表 | 是 | 否 | 是 | 否 |
| 翻譯記憶庫 | 是 | 否 | 是 | 否 |
| CLI／本機執行 | 是 | 不適用 | 否 | 手動處理 |
| 對 Git 友善的檔案 | 是 | 是 | 部分支援 | 手動處理 |
| 無需依賴 SaaS | 是 | 是 | 否 | 視情況而定 |
| 開放原始碼（AGPL-3.0） | 是 | 是 | 否 | 視情況而定 |

## 授權條款

[AGPL-3.0](../../LICENSE)

第三方相依套件聲明請參閱 [THIRD_PARTY_NOTICES.md](../../THIRD_PARTY_NOTICES.md)。

## 貢獻指南

開發環境設定與規範請參閱 [CONTRIBUTING.md](../../CONTRIBUTING.md)。所有貢獻均需完成 DCO 簽署。
