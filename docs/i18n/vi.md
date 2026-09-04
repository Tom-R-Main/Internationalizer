> [English (original)](../../README.md)

<p align="center">
  <img src="../../assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

Pipeline quốc tế hóa (i18n) native AI cho các dự án phần mềm. Dịch, xác thực và quản lý các tệp i18n bằng LLM.

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](../../LICENSE)

<p align="center">
<a href="ar.md">العربية</a> · <a href="bn.md">বাংলা</a> · <a href="cs.md">Čeština</a> · <a href="da.md">Dansk</a> · <a href="de.md">Deutsch</a> · <a href="el.md">Ελληνικά</a> · <a href="es.md">Español</a> · <a href="fi.md">Suomi</a> · <a href="fr.md">Français</a> · <a href="he.md">עברית</a> · <a href="hi.md">हिन्दी</a> · <a href="id.md">Indonesia</a> · <a href="it.md">Italiano</a> · <a href="ja.md">日本語</a> · <a href="ko.md">한국어</a> · <a href="ms.md">Bahasa Melayu</a><br><a href="nl.md">Nederlands</a> · <a href="pa.md">ਪੰਜਾਬੀ</a> · <a href="pl.md">Polski</a> · <a href="pt-BR.md">Português</a> · <a href="ro.md">Română</a> · <a href="ru.md">Русский</a> · <a href="sv.md">Svenska</a> · <a href="te.md">తెలుగు</a> · <a href="th.md">ไทย</a> · <a href="tr.md">Türkçe</a> · <a href="uk.md">Українська</a> · <a href="vi.md">Tiếng Việt</a> · <a href="yue.md">粵語</a> · <a href="zh-CN.md">简体中文</a> · <a href="zh-TW.md">繁體中文</a>
</p>

---

<!-- internationalizer:unit markdown:why-internationalizer -->
## Tại sao chọn Internationalizer?

Hầu hết các công cụ i18n đều là thư viện runtime (i18next, react-intl) hoặc nền tảng SaaS quản lý khóa (Crowdin, Lokalise). Không có công cụ nào giải quyết tốt vấn đề dịch thuật thực tế:

- **Dịch thủ công** không thể mở rộng quy mô quá một vài ngôn ngữ
- **API dịch máy** (Google Translate, DeepL) bỏ qua thuật ngữ, giọng văn và các quy ước UI của bạn
- **Dịch bằng LLM thông thường** hoạt động tốt hơn, nhưng nếu không có bảng thuật ngữ và hướng dẫn văn phong, kết quả nhận được sẽ thiếu nhất quán

Internationalizer thì khác. Đây là một **CLI pipeline** kết hợp dịch thuật bằng LLM với:

- **Bảng thuật ngữ theo từng ngôn ngữ** — thực thi thuật ngữ nhất quán trên toàn bộ ứng dụng của bạn
- **Hướng dẫn văn phong theo từng ngôn ngữ** — kiểm soát giọng văn, mức độ trang trọng, dạng số nhiều và quy tắc kiểu chữ
- **Bộ nhớ dịch thuật** — bỏ qua các chuỗi không thay đổi, tiết kiệm chi phí gọi API
- **Xác thực tất định** — phát hiện các key bị thiếu hoặc thừa, sai lệch cấu trúc được bảo vệ, vấn đề thuật ngữ, cùng các lỗi số nhiều hoặc lỗi ICU trước khi phát hành

<!-- internationalizer:unit markdown:installation -->
## Cài đặt

Cài đặt từ npm:

```bash
npm install -g internationalizer
```

Hoặc chạy mà không cần cài đặt global:

```bash
npx internationalizer --help
```

Gói npm sẽ cài đặt tệp nhị phân dựng sẵn phù hợp từ npm thông qua các dependency tùy chọn theo từng nền tảng.

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

<!-- internationalizer:unit markdown:npm-packages -->
## Các gói npm

- Git tag và phiên bản gói npm phải khớp nhau, ví dụ `v0.1.0` và `0.1.0`
- Gói `internationalizer` gốc phụ thuộc vào các gói nền tảng như `internationalizer-darwin-arm64`
- Các target npm được hỗ trợ: macOS arm64/x64, Linux arm64/x64, Windows x64
- Việc xuất bản qua CI yêu cầu một GitHub secret có tên là `NPM_TOKEN`

<!-- internationalizer:unit markdown:quick-start -->
## Bắt đầu nhanh

1. Tạo một tệp cấu hình trong thư mục gốc của dự án:

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

5. Xác thực tất cả các locale:

```bash
internationalizer validate
```

<!-- internationalizer:unit markdown:commands -->
## Lệnh

### `translate`

Tìm các key bị thiếu hoặc đã cũ và dịch chúng thông qua LLM.

```bash
internationalizer translate                    # translate all locales
internationalizer translate -l fr              # translate French only
internationalizer translate --dry-run          # preview without API calls
internationalizer translate --adopt-existing   # baseline existing translations without API calls
internationalizer translate --refresh-policy   # refresh prompt/style/model-stale entries
internationalizer translate --batch-size 20    # smaller batches
internationalizer translate --concurrency 2    # fewer parallel calls
```

Trạng thái dịch báo cáo độc lập các điều kiện missing (bị thiếu), source-stale (nguồn cũ), policy-stale (chính sách cũ), current (hiện hành) và manually edited (chỉnh sửa thủ công), nhờ đó một chỉnh sửa thủ công không thể che giấu thay đổi từ nguồn hoặc chính sách. Các giá trị policy-stale được báo cáo nhưng chỉ được dịch lại khi có cờ `--refresh-policy`. Các giá trị chỉnh sửa thủ công không bao giờ bị ghi đè tự động. Hãy dùng `--adopt-existing` khi áp dụng manifest vào các bản dịch đã được duyệt hoặc khi chấp thuận rõ ràng một chỉnh sửa thủ công đã qua rà soát làm mốc cơ sở mới.

### `validate`

Kiểm tra tất cả các tệp locale so với bundle nguồn tương ứng. Quá trình xác thực mặc định sẽ kiểm tra độ bao phủ cấu trúc (tỷ lệ phần trăm key đích bắt buộc hiện diện), báo cáo các key thừa dưới dạng cảnh báo, và báo lỗi nếu thiếu key, không khớp nội suy hoặc cấu trúc ICU MessageFormat không hợp lệ.

```bash
internationalizer validate                     # human-readable output
internationalizer validate --json              # machine-readable JSON
internationalizer validate -q                  # exit code only
internationalizer validate --strict             # enforce translation quality rules
internationalizer validate --require-state      # require current manifest provenance
```

`--strict` cũng báo cáo độ bao phủ bản dịch. Một giá trị ngôn ngữ trùng khớp với nguồn sẽ bị coi là chưa dịch trừ khi bảng thuật ngữ nêu rõ một mục từ nguồn và đích giống hệt nhau cho toàn bộ giá trị đó; `ignore_case` vẫn được áp dụng, nhưng một thuật ngữ trong bảng thuật ngữ nằm trong một giá trị dài hơn sẽ không được miễn trừ. Chế độ strict báo lỗi khi có key thừa, giá trị trùng khớp với nguồn, cấu trúc nội suy/HTML/mã/liên kết Markdown bị thay đổi, vi phạm bảng thuật ngữ và các dạng số nhiều đã cấu hình.

`--require-state` xác minh từng target so với `.internationalizer.lock`. Lệnh sẽ báo lỗi khi một key chưa được theo dõi, hoặc khi mã băm nguồn, chính sách dịch hay target đã ghi nhận bị cũ. Cờ này có thể kết hợp với `--strict`.

Báo cáo cho con người đọc và báo cáo JSON sử dụng các mã phát hiện ổn định:

| Mã | Ý nghĩa |
| --- | --- |
| `missing_key` / `extra_key` | Tập hợp key nguồn và đích khác nhau |
| `blank_translation` | Nguồn không rỗng nhưng target ở chế độ strict lại rỗng |
| `source_identical` | Giá trị ngôn ngữ ở chế độ strict vẫn chưa được dịch |
| `protected_structure_mismatch` | Cấu trúc nội suy, HTML, mã hoặc liên kết đã thay đổi |
| `glossary_violation` | Không tìm thấy thuật ngữ đích hoặc biến thể được phê duyệt |
| `plural_form_missing` | Thiếu một dạng số nhiều đã cấu hình của locale |
| `icu_message_syntax` | Thông điệp ICU nguồn hoặc đích bị sai định dạng |
| `icu_argument_mismatch` | Tên đối số, kiểu hoặc kiểu bộ định dạng của ICU khác nhau |
| `icu_selector_mismatch` | Bộ chọn khác nhau hoặc danh mục số nhiều không hợp lệ cho locale đích |
| `untracked` | Không có bản ghi manifest nào cho target |
| `source_stale` | Nội dung nguồn đã thay đổi sau lần dịch được ghi nhận |
| `policy_stale` | Prompt được tạo hoặc cài đặt mô hình đã thay đổi |
| `target_modified` | Nội dung target khác với bản ghi manifest |

### `detect`

Tự động phát hiện framework i18n và gợi ý cấu hình.

```bash
internationalizer detect
```

Hỗ trợ: react-i18next, next-intl, vue-i18n, vanilla JSON, tài liệu markdown.

### `glossary`

Quản lý các thuật ngữ trong bảng thuật ngữ theo từng ngôn ngữ được áp dụng bắt buộc trong quá trình dịch.

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

<!-- internationalizer:unit markdown:configuration-reference -->
## Tham chiếu cấu hình

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

Mã định danh locale phải là các tag BCP 47 đúng chuẩn như `fr`, `pt-BR` hoặc `sr-Latn-RS`. Các locale đích tương đương chuẩn hóa (canonical-equivalent) sẽ bị từ chối do trùng lặp, và cấu hình ghi đè nhà cung cấp theo từng locale sẽ khớp theo cách viết tương đương chuẩn hóa. Trong ví dụ trên, các locale không có cấu hình ghi đè — bao gồm cả tiếng Nhật — sẽ kế thừa cấu hình Gemini toàn cục.

Các giá trị ICU MessageFormat được phân tích cú pháp theo cấu trúc. Hỗ trợ các đối số đơn giản, `select`, `plural`, `selectordinal`, `number`, `date` và `time`, bao gồm cả thông điệp lồng nhau, độ dời số nhiều (plural offset), bộ chọn số chính xác và `#`. Quá trình xác thực kiểm tra cú pháp, kiểu đối số và kiểu bộ định dạng, độ dời số nhiều, sự đồng nhất nhánh select và các danh mục số nhiều CLDR của locale đích. Đầu ra từ nhà cung cấp phá vỡ các bất biến này sẽ bị từ chối trước khi tệp locale hoặc bản ghi bộ nhớ dịch thuật được ghi.

Với `i18next-v4`, các họ số nhiều nguồn được nhận diện sẽ được mở rộng trong quá trình dịch thành các danh mục CLDR của locale đích. Danh mục chỉ có ở đích sẽ sử dụng giá trị `_other` của họ nguồn làm mẫu dịch thuật. Xác thực strict bắt buộc phải có các danh mục đích đó; các danh mục chỉ có ở nguồn là tùy chọn đối với các locale đích không sử dụng chúng.

<!-- internationalizer:unit markdown:style-guides -->
## Hướng dẫn văn phong

Hướng dẫn văn phong là các tệp Markdown được đưa vào prompt dịch của LLM. Chúng kiểm soát giọng văn, mức độ trang trọng, quy tắc kiểu chữ và các quy ước đặc thù ngôn ngữ khác.

```
style-guides/
  _conventions.md    # shared rules for all languages
  fr.md              # French-specific rules
  ja.md              # Japanese-specific rules
  ar.md              # Arabic-specific rules
```

### Quy ước chung (`_conventions.md`)

Định nghĩa các quy tắc áp dụng cho mọi ngôn ngữ: cú pháp nội suy, bảo toàn HTML, quy ước loại chuỗi (nút bấm so với nhãn so với lỗi), v.v.

### Hướng dẫn theo từng ngôn ngữ (`{locale}.md`)

Định nghĩa các quy tắc theo từng ngôn ngữ: sắc thái trang trọng (tu so với vous), dấu câu (ngoặc kép, dấu hỏi ngược), các dạng số nhiều, định dạng ngày/số và bảng thuật ngữ.

Hướng dẫn văn phong là dữ liệu đầu vào chính sách bền vững chứ không phải đầu ra được tạo tự động. Internationalizer đọc nhưng không bao giờ ghi đè chúng. Nội dung của chúng được băm riêng biệt với bảng thuật ngữ và hợp đồng prompt, do đó thay đổi mã nguồn ứng dụng sẽ không làm bản dịch bị cũ. Việc chủ động chỉnh sửa hướng dẫn sẽ đánh dấu locale đó cần xem xét lại chính sách; thay đổi diễn đạt prompt nội bộ sẽ không gây ra điều này, trừ khi phiên bản hợp đồng prompt cũng thay đổi.

Xem [`examples/react-app/style-guides/`](../../examples/react-app/style-guides/) để biết ví dụ thực tế.

<!-- internationalizer:unit markdown:glossary-format -->
## Định dạng bảng thuật ngữ

Các tệp bảng thuật ngữ là mảng JSON được lưu trữ trong `{glossary_dir}/{locale}.json`:

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

`variants` liệt kê các dạng đích khác đã được duyệt. `enforcement` có thể là `error`, `warning` hoặc được bỏ qua để nhận hành vi mặc định là error. Các thuật ngữ được đưa vào prompt của LLM dưới dạng bảng thuật ngữ, đảm bảo tính nhất quán của bản dịch trên toàn bộ ứng dụng của bạn. Một mục từ chính xác như `{"source":"API","target":"API"}` cũng giúp miễn trừ toàn bộ giá trị trùng khớp với nguồn đó khỏi các phát hiện chưa dịch ở chế độ strict; mục từ này không miễn trừ một giá trị dài hơn chỉ đơn thuần chứa `API`.

<!-- internationalizer:unit markdown:translation-memory -->
## Bộ nhớ dịch thuật

Bộ nhớ dịch thuật được lưu trữ dưới dạng tệp JSONL (mỗi bản ghi JSON trên một dòng). Mỗi bản ghi bao gồm:

- Bundle, key, giá trị nguồn, giá trị đã dịch và locale đích chuẩn hóa
- Mã băm nguồn, hướng dẫn văn phong, bảng thuật ngữ, hợp đồng prompt và chính sách kết hợp
- Nhà cung cấp và mô hình đã tạo bản dịch
- Dấu thời gian

Trong những lần chạy tiếp theo, các chuỗi có cùng mã băm nguồn và chính sách sẽ được phân phát từ bộ nhớ cache mà không cần gọi LLM. Đường dẫn mặc định nằm trong thư mục bị bỏ qua `.internationalizer/`, do đó nó hoạt động như một bộ nhớ cache cục bộ. Hãy đặt `tm_path` tới một vị trí được theo dõi nếu dự án của bạn chủ đích chia sẻ bộ nhớ dịch thuật. Tệp manifest `.internationalizer.lock` có thể rà soát được quản lý phiên bản riêng biệt.

<!-- internationalizer:unit markdown:supported-formats -->
## Các định dạng được hỗ trợ

| Định dạng | Phần mở rộng | Chế độ |
|--------|-----------|------|
| JSON | `.json` | Key-value (lồng nhau, làm phẳng theo ký hiệu dấu chấm) |
| YAML | `.yml`, `.yaml` | Key-value (bảo toàn chú thích và thứ tự) |
| Markdown | `.md`, `.mdx` | Phần mở đầu (preamble) và các phần cấp H2 |

Các target Markdown chứa chú thích `internationalizer:unit` ẩn trước các phần H2. Những đánh dấu ổn định này cho phép Internationalizer thêm, di chuyển hoặc chỉnh sửa một phần nguồn mà không cần dịch lại các phần không liên quan khác. Các tài liệu hiện có chưa được đánh dấu sẽ nhận các điểm đánh dấu này trong lần cập nhật thành công tiếp theo.

<!-- internationalizer:unit markdown:project-type-detection -->
## Phát hiện loại dự án

`internationalizer detect` xác định thiết lập i18n của bạn bằng cách kiểm tra:

- Các dependency trong `package.json` để tìm react-i18next, next-intl hoặc vue-i18n
- Cấu trúc thư mục khớp với các khuôn mẫu locale phổ biến
- Phần mở rộng tệp và quy ước đặt tên

<!-- internationalizer:unit markdown:architecture -->
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
## So sánh với các giải pháp thay thế

| Tính năng | Internationalizer | i18next | Crowdin | LLM thông thường |
|---------|------------------|---------|---------|-------------|
| Dịch thuật bằng LLM | Có | Không | Một phần | Có |
| Hướng dẫn văn phong theo từng ngôn ngữ | Có | Không | Không | Không |
| Áp dụng bảng thuật ngữ | Có | Không | Có | Không |
| Bộ nhớ dịch thuật | Có | Không | Có | Không |
| CLI / thực thi cục bộ | Có | N/A | Không | Thủ công |
| Tệp thân thiện với Git | Có | Có | Một phần | Thủ công |
| Không phụ thuộc SaaS | Có | Có | Không | Tùy trường hợp |
| Mã nguồn mở (AGPL-3.0) | Có | Có | Không | Tùy trường hợp |

<!-- internationalizer:unit markdown:license -->
## Giấy phép

[AGPL-3.0](../../LICENSE)

Xem [THIRD_PARTY_NOTICES.md](../../THIRD_PARTY_NOTICES.md) để biết thông báo về dependency.

<!-- internationalizer:unit markdown:contributing -->
## Đóng góp

Xem [CONTRIBUTING.md](../../CONTRIBUTING.md) để biết hướng dẫn và thiết lập phát triển. Tất cả các đóng góp đều yêu cầu ký xác nhận DCO.
