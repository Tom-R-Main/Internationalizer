> [English (original)](../../README.md)

<p align="center">
  <img src="../../assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

專為軟件項目而設嘅 AI 原生國際化 pipeline。使用 LLM 翻譯、驗證同管理 i18n 檔案。

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](LICENSE)

## 點解要用 Internationalizer？

大部分 i18n 工具一係 runtime 庫（i18next、react-intl），一係 key 管理 SaaS 平台（Crowdin、Lokalise）。佢哋都冇辦法好好解決實際嘅翻譯問題：

- **人手翻譯** 喺語言數量多嗰陣好難擴展
- **機器翻譯 API**（Google Translate、DeepL）會忽略你嘅術語、語氣同 UI 慣例
- **通用 LLM 翻譯** 效果比較好，但如果冇詞彙表同風格指南，結果會唔一致

Internationalizer 唔同。佢係一個結合咗 LLM 翻譯嘅 **CLI pipeline**，並提供：

- **每種語言嘅詞彙表** — 確保成個應用程式嘅術語一致
- **每種語言嘅風格指南** — 控制語氣、正式程度、眾數同排版
- **翻譯記憶庫 (Translation memory)** — 跳過冇改過嘅字串，慳返 API 費用
- **Key 驗證** — 喺發佈之前搵出漏譯同插值 (interpolation) 唔匹配嘅問題

## 安裝

透過 npm 安裝：

```bash
npm install -g internationalizer
```

或者唔需要全域安裝直接執行：

```bash
npx internationalizer --help
```

npm 套件會透過特定平台嘅 optional dependencies 從 npm 安裝相應嘅預先編譯二進制檔案。

透過 Go 安裝：

```bash
go install github.com/Tom-R-Main/Internationalizer/cmd/internationalizer@latest
```

或者從源碼編譯：

```bash
git clone https://github.com/Tom-R-Main/Internationalizer.git
cd Internationalizer
go build -o internationalizer ./cmd/internationalizer
```

## npm 套件

- Git tag 同 npm 套件版本必須一致，例如 `v0.1.0` 同 `0.1.0`
- 根目錄嘅 `internationalizer` 套件依賴於平台套件，例如 `internationalizer-darwin-arm64`
- 支援嘅 npm 目標平台：macOS arm64/x64、Linux arm64/x64、Windows x64
- CI 發佈需要一個名為 `NPM_TOKEN` 嘅 GitHub secret

## 快速開始

1. 喺你嘅項目根目錄建立一個設定檔：

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

5. 驗證所有語言：

```bash
internationalizer validate
```

## 指令

### `translate`

搵出遺漏嘅 key 並透過 LLM 進行翻譯。

```bash
internationalizer translate                    # 翻譯所有語言
internationalizer translate -l fr              # 只翻譯法文
internationalizer translate --dry-run          # 預覽，唔呼叫 API
internationalizer translate --batch-size 20    # 較細嘅 batch
internationalizer translate --concurrency 2    # 減少並行呼叫
```

### `validate`

檢查所有語言檔案有冇遺漏嘅 key、多出嘅 key 同插值唔匹配嘅問題。

```bash
internationalizer validate                     # 人類可讀嘅輸出
internationalizer validate --json              # 機器可讀嘅 JSON
internationalizer validate -q                  # 只回傳 exit code
```

### `detect`

自動偵測 i18n 框架並建議設定。

```bash
internationalizer detect
```

支援：react-i18next、next-intl、vue-i18n、原生 JSON、markdown 文件。

### `glossary`

管理每種語言嘅詞彙表，喺翻譯期間強制執行。

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

管理翻譯記憶庫（之前翻譯過嘅字串嘅 JSONL 快取）。

```bash
internationalizer tm stats                     # 顯示記錄數量
internationalizer tm export                    # 匯出為 JSON
internationalizer tm clear --force             # 刪除所有記錄
```

## 設定參考

```yaml
# .internationalizer.yml

# 來源語言（預設：en）
source_locale: en

# 要翻譯成嘅語言（必填）
target_locales: [fr, de, es, ja, zh-CN, ar]

# 來源語言檔案嘅路徑（必填）
source_path: locales/en.json

# LLM 供應商設定
llm:
  # 供應商："anthropic"、"openai"、"gemini" 或 "openrouter"（預設：gemini）
  provider: gemini

  # 預設模型名稱（按供應商）：
  #   anthropic:  claude-sonnet-4-6
  #   openai:     gpt-5.4
  #   gemini:     gemini-3.1-pro-preview
  #   openrouter: google/gemini-3-flash-preview
  model: gemini-3.1-pro-preview

  # 包含 API key 嘅環境變數
  api_key_env: GOOGLE_AI_STUDIO_API_KEY

  # 兼容 OpenAI 嘅 endpoint 嘅 Base URL（選填）
  # base_url: https://api.openai.com

# 每次 LLM 呼叫嘅 key 數量（預設：40）
batch_size: 40

# 並行 LLM 呼叫數量（預設：4）
concurrency: 4

# 包含每種語言風格指南 Markdown 檔案嘅目錄（預設：style-guides）
style_guides_dir: style-guides

# 包含每種語言詞彙表 JSON 檔案嘅目錄（預設：glossary）
glossary_dir: glossary

# 翻譯記憶庫檔案嘅路徑（預設：.internationalizer/tm.jsonl）
tm_path: .internationalizer/tm.jsonl
```

## 風格指南

風格指南係 Markdown 檔案，會被注入到 LLM 翻譯 prompt 入面。佢哋控制語氣、正式程度、排版同其他特定語言嘅慣例。

```
style-guides/
  _conventions.md    # 所有語言共用嘅規則
  fr.md              # 法文專用規則
  ja.md              # 日文專用規則
  ar.md              # 阿拉伯文專用規則
```

### 共用慣例 (`_conventions.md`)

定義適用於所有語言嘅規則：插值語法、保留 HTML、字串類型慣例（按鈕 vs 標籤 vs 錯誤）等。

### 每種語言嘅指南 (`{locale}.md`)

定義特定語言嘅規則：正式程度（tu vs vous）、標點符號（書名號、倒問號）、眾數形式、日期/數字格式，以及術語詞彙表。

請參考 [`examples/react-app/style-guides/`](examples/react-app/style-guides/) 嘅實際例子。

## 詞彙表格式

詞彙表檔案係儲存喺 `{glossary_dir}/{locale}.json` 嘅 JSON 陣列：

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

術語會作為術語表注入到 LLM prompt 入面，確保成個應用程式嘅關鍵術語翻譯一致。

## 翻譯記憶庫

翻譯記憶庫儲存為 JSONL 檔案（每行一筆 JSON 記錄）。每筆記錄包含：

- 來源 key 同 value
- 翻譯後嘅 value
- 來源 value 嘅 SHA-256 hash
- 時間戳記

喺之後嘅執行中，冇改過嘅字串會從 TM 快取中提供，而唔需要呼叫 LLM，節省時間同 API 費用。TM 檔案對 git 友好，可以同你嘅語言檔案一齊 commit。

## 支援嘅格式

| 格式 | 副檔名 | 模式 |
|--------|-----------|------|
| JSON | `.json` | Key-value（巢狀，點號表示法攤平） |
| YAML | `.yml`, `.yaml` | Key-value（保留註解同順序） |
| Markdown | `.md`, `.mdx` | 成份文件翻譯 |

## 項目類型偵測

`internationalizer detect` 透過檢查以下項目嚟識別你嘅 i18n 設定：

- `package.json` 裡面有冇 react-i18next、next-intl 或 vue-i18n 嘅 dependencies
- 符合常見語言模式嘅目錄結構
- 副檔名同命名慣例

## 架構

```
cmd/internationalizer/     CLI 進入點同指令定義
internal/
  config/                  載入 YAML 設定同預設值
  detect/                  自動偵測項目類型
  formats/                 格式解析器（JSON、YAML、Markdown）
  glossary/                管理每種語言嘅詞彙表
  llm/                     LLM 供應商介面 + 實作
    anthropic.go           Anthropic Claude 後端
    openai.go              OpenAI / 兼容後端
    gemini.go              透過 AI Studio 嘅 Google Gemini 後端
                           OpenRouter 使用 openai.go 配合自訂 base_url
  styleguide/              風格指南載入器
  tm/                      JSONL 翻譯記憶庫
  translate/               翻譯協調器
  validate/                語言驗證同差異比對
```

## 替代方案比較

| 功能 | Internationalizer | i18next | Crowdin | 通用 LLM |
|---------|------------------|---------|---------|-------------|
| LLM 驅動翻譯 | 是 | 否 | 部分 | 是 |
| 每種語言風格指南 | 是 | 否 | 否 | 否 |
| 強制執行詞彙表 | 是 | 否 | 是 | 否 |
| 翻譯記憶庫 | 是 | 否 | 是 | 否 |
| CLI / 本機執行 | 是 | 不適用 | 否 | 手動 |
| Git 友好檔案 | 是 | 是 | 部分 | 手動 |
| 冇 SaaS 依賴 | 是 | 是 | 否 | 視情況而定 |
| 開源 (AGPL-3.0) | 是 | 是 | 否 | 視情況而定 |

## 授權條款

[AGPL-3.0](LICENSE)

## 貢獻

請參閱 [CONTRIBUTING.md](CONTRIBUTING.md) 了解開發設定同指南。所有貢獻都需要 DCO 簽署。

