<p align="center">
  <img src="assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

專為軟體專案打造的 AI 原生國際化管線。使用 LLM 來翻譯、驗證和管理 i18n 檔案。

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](LICENSE)

## 為什麼選擇 Internationalizer？

大多數的 i18n 工具不是執行階段函式庫 (i18next, react-intl)，就是金鑰管理 SaaS 平台 (Crowdin, Lokalise)。它們都沒有妥善解決實際的翻譯問題：

- **手動翻譯**在語言數量增加後難以擴展
- **機器翻譯 API** (Google Translate, DeepL) 會忽略您的術語、語氣和 UI 慣例
- **通用 LLM 翻譯**效果較好，但如果沒有詞彙表和風格指南，翻譯結果會不一致

Internationalizer 與眾不同。它是一個結合了 LLM 翻譯與以下功能的 **CLI 管線**：

- **各語言專屬詞彙表** — 確保應用程式中的術語保持一致
- **各語言專屬風格指南** — 控制語氣、正式程度、複數形式和排版
- **翻譯記憶庫** — 略過未變更的字串，節省 API 呼叫費用
- **金鑰驗證** — 在發布前捕捉遺漏的翻譯和插值不符的問題

## 安裝

透過 npm 安裝：

```bash
npm install -g internationalizer
```

或在不全域安裝的情況下執行：

```bash
npx internationalizer --help
```

npm 套件會透過特定平台的選用相依性，從 npm 安裝相符的預先建置二進位檔。

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

- Git 標籤和 npm 套件版本必須相符，例如 `v0.1.0` 和 `0.1.0`
- 根目錄的 `internationalizer` 套件相依於平台套件，例如 `internationalizer-darwin-arm64`
- 支援的 npm 目標平台：macOS arm64/x64、Linux arm64/x64、Windows x64
- CI 發布需要名為 `NPM_TOKEN` 的 GitHub secret

## 快速入門

1. 在您的專案根目錄建立設定檔：

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

5. 驗證所有地區設定：

```bash
internationalizer validate
```

## 指令

### `translate`

尋找遺漏的金鑰並透過 LLM 進行翻譯。

```bash
internationalizer translate                    # 翻譯所有地區設定
internationalizer translate -l fr              # 僅翻譯法文
internationalizer translate --dry-run          # 預覽而不呼叫 API
internationalizer translate --batch-size 20    # 較小的批次
internationalizer translate --concurrency 2    # 較少的平行呼叫
```

### `validate`

檢查所有地區設定檔案是否有遺漏的金鑰、多餘的金鑰以及插值不符的情況。

```bash
internationalizer validate                     # 人類可讀的輸出
internationalizer validate --json              # 機器可讀的 JSON
internationalizer validate -q                  # 僅輸出結束代碼
```

### `detect`

自動偵測 i18n 框架並建議設定。

```bash
internationalizer detect
```

支援：react-i18next、next-intl、vue-i18n、原生 JSON、Markdown 文件。

### `glossary`

管理在翻譯期間強制執行的各語言專屬詞彙表術語。

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

管理翻譯記憶庫 (先前翻譯字串的 JSONL 快取)。

```bash
internationalizer tm stats                     # 顯示記錄數量
internationalizer tm export                    # 匯出為 JSON
internationalizer tm clear --force             # 刪除所有記錄
```

## 設定參考

```yaml
# .internationalizer.yml

# 來源語言 (預設：en)
source_locale: en

# 要翻譯成的目標語言 (必填)
target_locales: [fr, de, es, ja, zh-CN, ar]

# 來源地區設定檔案的路徑 (必填)
source_path: locales/en.json

# LLM 供應商設定
llm:
  # 供應商："anthropic"、"openai"、"gemini" 或 "openrouter" (預設：gemini)
  provider: gemini

  # 各供應商的預設模型名稱：
  #   anthropic:  claude-sonnet-4-6
  #   openai:     gpt-5.4
  #   gemini:     gemini-3.1-pro-preview
  #   openrouter: google/gemini-3-flash-preview
  model: gemini-3.1-pro-preview

  # 包含 API 金鑰的環境變數
  api_key_env: GOOGLE_AI_STUDIO_API_KEY

  # 相容於 OpenAI 端點的基礎 URL (選填)
  # base_url: https://api.openai.com

# 每次 LLM 呼叫的金鑰數量 (預設：40)
batch_size: 40

# 平行 LLM 呼叫數量 (預設：4)
concurrency: 4

# 包含各語言專屬風格指南 Markdown 檔案的目錄 (預設：style-guides)
style_guides_dir: style-guides

# 包含各語言專屬詞彙表 JSON 檔案的目錄 (預設：glossary)
glossary_dir: glossary

# 翻譯記憶庫檔案的路徑 (預設：.internationalizer/tm.jsonl)
tm_path: .internationalizer/tm.jsonl
```

## 風格指南

風格指南是注入到 LLM 翻譯提示中的 Markdown 檔案。它們控制語氣、正式程度、排版以及其他特定語言的慣例。

```
style-guides/
  _conventions.md    # 所有語言的共用規則
  fr.md              # 法文專屬規則
  ja.md              # 日文專屬規則
  ar.md              # 阿拉伯文專屬規則
```

### 共用慣例 (`_conventions.md`)

定義適用於所有語言的規則：插值語法、HTML 保留、字串類型慣例 (按鈕 vs. 標籤 vs. 錯誤) 等。

### 各語言專屬指南 (`{locale}.md`)

定義特定語言的規則：正式程度 (tu vs. vous)、標點符號 (法文引號、倒問號)、複數形式、日期/數字格式以及術語詞彙表。

請參閱 [`examples/react-app/style-guides/`](examples/react-app/style-guides/) 以取得實際範例。

## 詞彙表格式

詞彙表檔案是儲存在 `{glossary_dir}/{locale}.json` 中的 JSON 陣列：

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

術語會作為術語表注入到 LLM 提示中，確保應用程式中關鍵術語的翻譯保持一致。

## 翻譯記憶庫

翻譯記憶庫儲存為 JSONL 檔案 (每行一筆 JSON 記錄)。每筆記錄包含：

- 來源金鑰和值
- 翻譯後的值
- 來源值的 SHA-256 雜湊
- 時間戳記

在後續執行時，未變更的字串會從 TM 快取中提供，而無需呼叫 LLM，從而節省時間和 API 成本。TM 檔案對 Git 友善，可以與您的地區設定檔案一起提交。

## 支援的格式

| 格式 | 副檔名 | 模式 |
|--------|-----------|------|
| JSON | `.json` | 鍵值對 (巢狀、點記號扁平化) |
| YAML | `.yml`, `.yaml` | 鍵值對 (保留註解和順序) |
| Markdown | `.md`, `.mdx` | 整份文件翻譯 |

## 專案類型偵測

`internationalizer detect` 會透過檢查以下項目來識別您的 i18n 設定：

- `package.json` 中 react-i18next、next-intl 或 vue-i18n 的相依性
- 符合常見地區設定模式的目錄結構
- 副檔名和命名慣例

## 架構

```
cmd/internationalizer/     CLI 進入點和指令定義
internal/
  config/                  載入 YAML 設定與預設值
  detect/                  專案類型自動偵測
  formats/                 格式解析器 (JSON、YAML、Markdown)
  glossary/                各語言專屬詞彙表管理
  llm/                     LLM 供應商介面與實作
    anthropic.go           Anthropic Claude 後端
    openai.go              OpenAI / 相容後端
    gemini.go              透過 AI Studio 的 Google Gemini 後端
                           OpenRouter 使用帶有自訂 base_url 的 openai.go
  styleguide/              風格指南載入器
  tm/                      JSONL 翻譯記憶庫
  translate/               翻譯協調器
  validate/                地區設定驗證與差異比對
```

## 與替代方案的比較

| 功能 | Internationalizer | i18next | Crowdin | 通用 LLM |
|---------|------------------|---------|---------|-------------|
| LLM 驅動翻譯 | 是 | 否 | 部分 | 是 |
| 各語言專屬風格指南 | 是 | 否 | 否 | 否 |
| 強制執行詞彙表 | 是 | 否 | 是 | 否 |
| 翻譯記憶庫 | 是 | 否 | 是 | 否 |
| CLI / 本機執行 | 是 | 不適用 | 否 | 手動 |
| 對 Git 友善的檔案 | 是 | 是 | 部分 | 手動 |
| 無 SaaS 相依性 | 是 | 是 | 否 | 視情況而定 |
| 開源 (AGPL-3.0) | 是 | 是 | 否 | 視情況而定 |

## 授權條款

[AGPL-3.0](LICENSE)

## 貢獻

請參閱 [CONTRIBUTING.md](CONTRIBUTING.md) 以了解開發設定和指南。所有貢獻都需要 DCO 簽署。

