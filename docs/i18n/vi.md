> [English (original)](../../README.md)

<p align="center">
  <img src="../../assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

Pipeline internationalization (i18n) gốc AI dành cho các dự án phần mềm. Dịch, xác thực và quản lý các tệp i18n bằng LLM.

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](LICENSE)

<p align="center">
<a href="docs/i18n/ar.md">العربية</a> · <a href="docs/i18n/bn.md">বাংলা</a> · <a href="docs/i18n/de.md">Deutsch</a> · <a href="docs/i18n/es.md">Español</a> · <a href="docs/i18n/fr.md">Français</a> · <a href="docs/i18n/hi.md">हिन्दी</a> · <a href="docs/i18n/id.md">Indonesia</a> · <a href="docs/i18n/ja.md">日本語</a> · <a href="docs/i18n/pa.md">ਪੰਜਾਬੀ</a> · <a href="docs/i18n/pt-BR.md">Português</a> · <a href="docs/i18n/ru.md">Русский</a> · <a href="docs/i18n/te.md">తెలుగు</a><br>
<a href="docs/i18n/th.md">ไทย</a> · <a href="docs/i18n/uk.md">Українська</a> · <a href="docs/i18n/yue.md">粵語</a> · <a href="docs/i18n/zh-CN.md">简体中文</a> · <a href="docs/i18n/zh-TW.md">繁體中文</a>
</p>

---

## Tại sao chọn Internationalizer?

Hầu hết các công cụ i18n đều là thư viện runtime (i18next, react-intl) hoặc nền tảng SaaS quản lý key (Crowdin, Lokalise). Không có công cụ nào giải quyết tốt vấn đề dịch thuật thực sự:

- **Dịch thủ công** không thể mở rộng khi có nhiều ngôn ngữ
- **API dịch máy** (Google Translate, DeepL) bỏ qua thuật ngữ, giọng văn và quy ước UI của bạn
- **Dịch bằng LLM thông thường** hoạt động tốt hơn, nhưng nếu không có bảng thuật ngữ và hướng dẫn văn phong, bạn sẽ nhận được kết quả không nhất quán

Internationalizer thì khác. Đây là một **CLI pipeline** kết hợp dịch thuật LLM với:

- **Bảng thuật ngữ theo ngôn ngữ** — đảm bảo tính nhất quán của thuật ngữ trên toàn bộ ứng dụng
- **Hướng dẫn văn phong theo ngôn ngữ** — kiểm soát giọng văn, mức độ trang trọng, số nhiều và cách trình bày
- **Bộ nhớ dịch thuật (Translation memory)** — bỏ qua các chuỗi không thay đổi, tiết kiệm chi phí gọi API
- **Xác thực key** — phát hiện các bản dịch bị thiếu và lỗi khớp biến nội suy (interpolation) trước khi phát hành

## Cài đặt

Cài đặt từ npm:

```bash
npm install -g internationalizer
```

Hoặc chạy mà không cần cài đặt global:

```bash
npx internationalizer --help
```

Gói npm sẽ cài đặt tệp nhị phân dựng sẵn tương ứng từ npm thông qua các dependency tùy chọn theo nền tảng.

Cài đặt bằng Go:

```bash
go install github.com/Tom-R-Main/Internationalizer/cmd/internationalizer@latest
```

Hoặc build từ mã nguồn:

```bash
git clone https://github.com/Tom-R-Main/Internationalizer.git
cd Internationalizer
go build -o internationalizer ./cmd/internationalizer
```

## Các gói npm

- Git tag và phiên bản gói npm phải khớp nhau, ví dụ `v0.1.0` và `0.1.0`
- Gói `internationalizer` gốc phụ thuộc vào các gói nền tảng như `internationalizer-darwin-arm64`
- Các mục tiêu npm được hỗ trợ: macOS arm64/x64, Linux arm64/x64, Windows x64
- Việc xuất bản qua CI yêu cầu một GitHub secret có tên là `NPM_TOKEN`

## Bắt đầu nhanh

1. Tạo tệp cấu hình trong thư mục gốc của dự án:

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

2. Thiết lập API key của bạn:

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

3. Xem trước những gì sẽ được dịch:

```bash
internationalizer translate --dry-run
```

4. Chạy quá trình dịch:

```bash
internationalizer translate
```

5. Xác thực tất cả các ngôn ngữ:

```bash
internationalizer validate
```

## Lệnh

### `translate`

Tìm các key bị thiếu và dịch chúng thông qua LLM.

```bash
internationalizer translate                    # translate all locales
internationalizer translate -l fr              # translate French only
internationalizer translate --dry-run          # preview without API calls
internationalizer translate --batch-size 20    # smaller batches
internationalizer translate --concurrency 2    # fewer parallel calls
```

### `validate`

Kiểm tra tất cả các tệp ngôn ngữ để tìm key bị thiếu, key thừa và lỗi khớp biến nội suy.

```bash
internationalizer validate                     # human-readable output
internationalizer validate --json              # machine-readable JSON
internationalizer validate -q                  # exit code only
```

### `detect`

Tự động phát hiện framework i18n và đề xuất cấu hình.

```bash
internationalizer detect
```

Hỗ trợ: react-i18next, next-intl, vue-i18n, vanilla JSON, markdown docs.

### `glossary`

Quản lý các thuật ngữ trong bảng thuật ngữ theo ngôn ngữ được áp dụng bắt buộc trong quá trình dịch.

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

Quản lý bộ nhớ dịch thuật (bộ nhớ cache JSONL của các chuỗi đã dịch trước đó).

```bash
internationalizer tm stats                     # show record counts
internationalizer tm export                    # dump as JSON
internationalizer tm clear --force             # delete all records
```

## Tham chiếu cấu hình

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

## Hướng dẫn văn phong

Hướng dẫn văn phong là các tệp Markdown được đưa vào prompt dịch của LLM. Chúng kiểm soát giọng văn, mức độ trang trọng, cách trình bày và các quy ước đặc thù khác của ngôn ngữ.

```
style-guides/
  _conventions.md    # shared rules for all languages
  fr.md              # French-specific rules
  ja.md              # Japanese-specific rules
  ar.md              # Arabic-specific rules
```

### Quy ước chung (`_conventions.md`)

Định nghĩa các quy tắc áp dụng cho tất cả ngôn ngữ: cú pháp nội suy, bảo toàn HTML, quy ước loại chuỗi (nút bấm so với nhãn so với lỗi), v.v.

### Hướng dẫn theo ngôn ngữ (`{locale}.md`)

Định nghĩa các quy tắc đặc thù của ngôn ngữ: mức độ trang trọng (tu so với vous), dấu câu (ngoặc kép, dấu hỏi ngược), các dạng số nhiều, định dạng ngày/số và bảng thuật ngữ.

Xem [`examples/react-app/style-guides/`](examples/react-app/style-guides/) để biết ví dụ thực tế.

## Định dạng bảng thuật ngữ

Các tệp bảng thuật ngữ là mảng JSON được lưu trữ trong `{glossary_dir}/{locale}.json`:

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

Các thuật ngữ được đưa vào prompt của LLM dưới dạng bảng thuật ngữ, đảm bảo việc dịch nhất quán các thuật ngữ chính trên toàn bộ ứng dụng của bạn.

## Bộ nhớ dịch thuật

Bộ nhớ dịch thuật được lưu trữ dưới dạng tệp JSONL (mỗi dòng là một bản ghi JSON). Mỗi bản ghi chứa:

- Key và giá trị nguồn
- Giá trị đã dịch
- Mã băm SHA-256 của giá trị nguồn
- Dấu thời gian

Trong các lần chạy tiếp theo, các chuỗi không thay đổi sẽ được lấy từ bộ nhớ cache TM mà không cần gọi LLM, giúp tiết kiệm cả thời gian và chi phí API. Tệp TM thân thiện với git và có thể được commit cùng với các tệp ngôn ngữ của bạn.

## Các định dạng được hỗ trợ

| Định dạng | Phần mở rộng | Chế độ |
|--------|-----------|------|
| JSON | `.json` | Key-value (lồng nhau, làm phẳng bằng dấu chấm) |
| YAML | `.yml`, `.yaml` | Key-value (bảo toàn chú thích và thứ tự) |
| Markdown | `.md`, `.mdx` | Dịch toàn bộ tài liệu |

## Phát hiện loại dự án

`internationalizer detect` xác định thiết lập i18n của bạn bằng cách kiểm tra:

- Các dependency trong `package.json` cho react-i18next, next-intl hoặc vue-i18n
- Cấu trúc thư mục khớp với các mẫu ngôn ngữ phổ biến
- Phần mở rộng tệp và quy ước đặt tên

## Kiến trúc

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

## So sánh với các giải pháp thay thế

| Tính năng | Internationalizer | i18next | Crowdin | LLM thông thường |
|---------|------------------|---------|---------|-------------|
| Dịch thuật bằng LLM | Có | Không | Một phần | Có |
| Hướng dẫn văn phong theo ngôn ngữ | Có | Không | Không | Không |
| Áp dụng bảng thuật ngữ | Có | Không | Có | Không |
| Bộ nhớ dịch thuật | Có | Không | Có | Không |
| CLI / thực thi cục bộ | Có | N/A | Không | Thủ công |
| Tệp thân thiện với Git | Có | Có | Một phần | Thủ công |
| Không phụ thuộc SaaS | Có | Có | Không | Tùy trường hợp |
| Mã nguồn mở (AGPL-3.0) | Có | Có | Không | Tùy trường hợp |

## Giấy phép

[AGPL-3.0](LICENSE)

## Đóng góp

Xem [CONTRIBUTING.md](CONTRIBUTING.md) để biết hướng dẫn và thiết lập phát triển. Tất cả các đóng góp đều yêu cầu xác nhận DCO.

