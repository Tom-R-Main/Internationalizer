> [English (original)](../../README.md)

<p align="center">
  <img src="../../assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

面向软件项目的 AI 原生国际化流水线。基于 LLM 翻译、校验和管理 i18n 文件。

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](../../LICENSE)

<p align="center">
<a href="ar.md">العربية</a> · <a href="bn.md">বাংলা</a> · <a href="cs.md">Čeština</a> · <a href="da.md">Dansk</a> · <a href="de.md">Deutsch</a> · <a href="el.md">Ελληνικά</a> · <a href="es.md">Español</a> · <a href="fi.md">Suomi</a> · <a href="fr.md">Français</a> · <a href="he.md">עברית</a> · <a href="hi.md">हिन्दी</a> · <a href="id.md">Indonesia</a> · <a href="it.md">Italiano</a> · <a href="ja.md">日本語</a> · <a href="ko.md">한국어</a> · <a href="ms.md">Bahasa Melayu</a><br><a href="nl.md">Nederlands</a> · <a href="pa.md">ਪੰਜਾਬੀ</a> · <a href="pl.md">Polski</a> · <a href="pt-BR.md">Português</a> · <a href="ro.md">Română</a> · <a href="ru.md">Русский</a> · <a href="sv.md">Svenska</a> · <a href="te.md">తెలుగు</a> · <a href="th.md">ไทย</a> · <a href="tr.md">Türkçe</a> · <a href="uk.md">Українська</a> · <a href="vi.md">Tiếng Việt</a> · <a href="yue.md">粵語</a> · <a href="zh-CN.md">简体中文</a> · <a href="zh-TW.md">繁體中文</a>
</p>

---
<!-- internationalizer:unit markdown:why-internationalizer -->
## 为什么选择 Internationalizer？

多数 i18n 工具要么是运行时库（i18next、react-intl），要么是键管理 SaaS 平台（Crowdin、Lokalise）。它们都未能真正解决实际的翻译问题：

- **人工翻译**在超出几种语言后便难以扩展
- **机器翻译 API**（Google Translate、DeepL）会忽略术语表、语气和 UI 规范
- **通用 LLM 翻译**效果更好，但缺少术语表与风格指南会导致译文风格不一致

Internationalizer 与众不同。它是一套 **CLI 流水线**，将 LLM 翻译与以下能力深度结合：

- **各语言专属术语表** — 在全应用范围内统一专业术语
- **各语言专属风格指南** — 把控语气、正式度、复数规则与排版格式
- **翻译记忆库** — 跳过未修改的字符串，节省 API 调用成本
- **确定性校验** — 在代码发布前拦截缺失或多余的键、受保护结构偏离、术语表问题，以及复数或 ICU 错误
<!-- internationalizer:unit markdown:installation -->
## 安装

通过 npm 安装：

```bash
npm install -g internationalizer
```

或无需全局安装直接运行：

```bash
npx internationalizer --help
```

npm 包会通过特定平台的可选依赖项从 npm 安装匹配的预构建二进制文件。

使用 Go 安装：

```bash
go install github.com/Tom-R-Main/Internationalizer/cmd/internationalizer@latest
```

或从源码构建：

```bash
git clone https://github.com/Tom-R-Main/Internationalizer.git
cd Internationalizer
go build -o internationalizer ./cmd/internationalizer
```
<!-- internationalizer:unit markdown:npm-packages -->
## npm 包

- Git 标签与 npm 包版本必须一致，例如 `v0.1.0` 与 `0.1.0`
- 根包 `internationalizer` 依赖各平台分发包，例如 `internationalizer-darwin-arm64`
- 支持的 npm 目标平台：macOS arm64/x64、Linux arm64/x64、Windows x64
- CI 发布需要配置名为 `NPM_TOKEN` 的 GitHub Secret
<!-- internationalizer:unit markdown:quick-start -->
## 快速开始

1. 在项目根目录创建配置文件：

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

2. 设置 API 密钥：

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

3. 预览待翻译内容：

```bash
internationalizer translate --dry-run
```

4. 运行翻译：

```bash
internationalizer translate
```

5. 校验所有语言环境文件：

```bash
internationalizer validate
```
<!-- internationalizer:unit markdown:commands -->
## 命令

### `translate`

查找缺失或过期的键，并通过 LLM 进行翻译。

```bash
internationalizer translate                    # translate all locales
internationalizer translate -l fr              # translate French only
internationalizer translate --dry-run          # preview without API calls
internationalizer translate --adopt-existing   # baseline existing translations without API calls
internationalizer translate --refresh-policy   # refresh prompt/style/model-stale entries
internationalizer translate --batch-size 20    # smaller batches
internationalizer translate --concurrency 2    # fewer parallel calls
```

翻译状态会独立报告缺失（missing）、源内容过期（source-stale）、策略过期（policy-stale）、最新（current）和手动编辑（manually edited）状态，因此手动编辑不会掩盖源内容或策略变更。策略过期的项会被报告，但仅在传入 `--refresh-policy` 时才会重新翻译。系统绝不会自动覆盖手动编辑的值。在初次将清单引入已审校的译文，或明确采纳已审校的手动编辑作为新基线时，请使用 `--adopt-existing`。

### `validate`

将所有语言环境文件与源资源包进行比对。默认校验会检查结构覆盖率（所需目标键存在的百分比），将多余的键报告为警告，并在出现键缺失、插值不匹配或 ICU MessageFormat 结构无效时失败退出。

```bash
internationalizer validate                     # human-readable output
internationalizer validate --json              # machine-readable JSON
internationalizer validate -q                  # exit code only
internationalizer validate --strict             # enforce translation quality rules
internationalizer validate --require-state      # require current manifest provenance
```

`--strict` 还会报告已翻译覆盖率。如果语言文本与源文本相同，则被视为未翻译，除非术语表中明确包含一条针对该完整文本、且源文与目标文本完全一致的精确词条；校验支持 `ignore_case`，但较长文本中仅包含某术语词条并不能获得豁免。严格模式会在存在多余键、与源文本相同的值、插值/HTML/代码/Markdown 链接结构发生变更、术语表违规以及缺少已配置的复数形式时失败。

`--require-state` 会将每个目标项与 `.internationalizer.lock` 进行比对验证。当键未被跟踪，或者记录的源内容、翻译策略或目标哈希已过期时，校验将失败。该参数可与 `--strict` 结合使用。

人类可读报告与 JSON 报告均使用稳定的发现项代码：

| 代码 | 含义 |
| --- | --- |
| `missing_key` / `extra_key` | 源键集与目标键集不一致 |
| `blank_translation` | 非空源内容对应的严格模式目标内容为空 |
| `source_identical` | 严格模式下语言文本仍未翻译 |
| `protected_structure_mismatch` | 插值、HTML、代码或链接结构发生变更 |
| `glossary_violation` | 未找到获准的目标术语或变体 |
| `plural_form_missing` | 缺少目标语言环境已配置的复数形式 |
| `icu_message_syntax` | 源或目标 ICU 消息格式有误 |
| `icu_argument_mismatch` | ICU 参数名、类型或格式化器样式不一致 |
| `icu_selector_mismatch` | 选择器不一致，或复数类别不适用于目标语言环境 |
| `untracked` | 清单中不存在该目标的记录 |
| `source_stale` | 记录翻译后源内容发生变化 |
| `policy_stale` | 生成的提示词或模型设置发生变化 |
| `target_modified` | 目标内容与清单记录不一致 |

### `detect`

自动检测项目所使用的 i18n 框架并给出配置建议。

```bash
internationalizer detect
```

支持：react-i18next、next-intl、vue-i18n、普通 JSON、Markdown 文档。

### `glossary`

管理在翻译过程中强制执行的各语言专属术语表。

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

管理翻译记忆库（以 JSONL 缓存以往翻译的字符串）。

```bash
internationalizer tm stats                     # show record counts
internationalizer tm export                    # dump as JSON
internationalizer tm clear --force             # delete all records
```
<!-- internationalizer:unit markdown:configuration-reference -->
## 配置参考

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

语言环境标识符必须是格式规范的 BCP 47 标签，如 `fr`、`pt-BR` 或 `sr-Latn-RS`。规范等价的目标语言环境会被作为重复项拒绝，特定语言环境的服务商覆盖项也会按规范等价拼写进行匹配。在上述示例中，未配置覆盖项的语言环境（包括日语）将继承全局 Gemini 配置。

ICU MessageFormat 文本会按结构解析。支持简单参数、`select`、`plural`、`selectordinal`、`number`、`date` 以及 `time`，包括嵌套消息、复数偏移量、精确数字选择器和 `#`。校验会检查语法、参数类型与格式化器样式、复数偏移量、select 分支一致性，以及目标语言环境的 CLDR 复数类别。破坏这些不变量的服务商输出会在写入语言环境文件或翻译记忆库记录前被直接拒收。

使用 `i18next-v4` 时，识别到的源复数族在翻译期间会扩展为目标语言环境对应的 CLDR 类别。仅目标语言环境需要的类别会使用源复数族的 `_other` 值作为其翻译模板。严格校验要求这些目标类别齐全；对于不使用这些类别的目标语言环境，仅源语言存在的类别则是可选的。
<!-- internationalizer:unit markdown:style-guides -->
## 风格指南

风格指南是注入到 LLM 翻译提示词中的 Markdown 文件。它们控制语气、正式度、排版格式以及其他特定语言的规范。

```
style-guides/
  _conventions.md    # shared rules for all languages
  fr.md              # French-specific rules
  ja.md              # Japanese-specific rules
  ar.md              # Arabic-specific rules
```

### 通用规范（`_conventions.md`）

定义适用于所有语言的规则：插值语法、HTML 保留规则、字符串类型约定（按钮与标签及错误的表述规范）等。

### 各语言专属指南（`{locale}.md`）

定义特定语言规则：正式度语体（tu 与 vous）、标点符号（书名号/引号、倒问号）、复数形式、日期/数字格式化，以及术语表。

风格指南是持久的策略输入，而非生成的输出。Internationalizer 读取它们但绝不重写。其内容与术语表和提示词契约分开哈希，因此应用程序代码的更改不会导致翻译过期。修改指南会有意将该语言环境标记为策略待复查；更改内部提示词措辞则不会，除非提示词契约版本也同时变更。

完整示例请参见 [`examples/react-app/style-guides/`](../../examples/react-app/style-guides/)。
<!-- internationalizer:unit markdown:glossary-format -->
## 术语表格式

术语表文件是保存在 `{glossary_dir}/{locale}.json` 中的 JSON 数组：

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

`variants` 列出其他批准的目标形式。`enforcement` 可以是 `error`、`warning`，或者省略以采用默认的错误行为。术语会作为术语表注入到 LLM 提示词中，以确保全应用翻译的一致性。类似 `{"source":"API","target":"API"}` 这样源文与目标完全一致的精确词条，也会免除该完整相同值在严格模式下的未翻译发现项；但仅包含 `API` 的较长文本并不能获得豁免。
<!-- internationalizer:unit markdown:translation-memory -->
## 翻译记忆库

翻译记忆库存储为 JSONL 文件（每行一个 JSON 记录）。每条记录包含：

- 资源包、键、源值、翻译值及规范的目标语言环境
- 源内容、风格指南、术语表、提示词契约以及合并策略哈希
- 生成翻译的服务商与模型
- 时间戳

在随后的运行中，具有相同源内容和策略哈希的字符串将直接从缓存中获取，而无需调用 LLM。默认路径位于被忽略的 `.internationalizer/` 目录中，因此它作为本地缓存存在。如果项目需要显式共享翻译记忆库，请将 `tm_path` 设置为版本控制跟踪的路径。可供审查的 `.internationalizer.lock` 清单则单独进行版本控制。
<!-- internationalizer:unit markdown:supported-formats -->
## 支持的格式

| 格式 | 扩展名 | 模式 |
|--------|-----------|------|
| JSON | `.json` | 键值对（嵌套结构，按点号表示法扁平化） |
| YAML | `.yml`, `.yaml` | 键值对（保留注释与顺序） |
| Markdown | `.md`, `.mdx` | 序言与二级标题（H2）级别章节 |

Markdown 目标文件在二级标题前包含不可见的 `internationalizer:unit` 注释。这些稳定标记使 Internationalizer 能够新增、移动或编辑单个源章节，而无需重新翻译不相关的章节。现有未标记的文档将在下一次成功更新时添加标记。
<!-- internationalizer:unit markdown:project-type-detection -->
## 项目类型检测

`internationalizer detect` 通过检查以下内容来识别您的 i18n 配置：

- `package.json` 中是否存在 react-i18next、next-intl 或 vue-i18n 依赖
- 是否存在匹配通用区域设置模式的目录结构
- 文件扩展名与命名约定
<!-- internationalizer:unit markdown:architecture -->
## 架构

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
## 同类方案对比

| 功能特性 | Internationalizer | i18next | Crowdin | 通用 LLM |
|---------|------------------|---------|---------|-------------|
| LLM 驱动翻译 | 是 | 否 | 部分 | 是 |
| 各语言专属风格指南 | 是 | 否 | 否 | 否 |
| 术语表强制执行 | 是 | 否 | 是 | 否 |
| 翻译记忆库 | 是 | 否 | 是 | 否 |
| CLI / 本地执行 | 是 | 不适用 | 否 | 手动 |
| Git 友好文件 | 是 | 是 | 部分 | 手动 |
| 无 SaaS 依赖 | 是 | 是 | 否 | 视情况而定 |
| 开源（AGPL-3.0） | 是 | 是 | 否 | 视情况而定 |
<!-- internationalizer:unit markdown:license -->
## 许可证

[AGPL-3.0](../../LICENSE)

有关依赖项声明，请参见 [THIRD_PARTY_NOTICES.md](../../THIRD_PARTY_NOTICES.md)。
<!-- internationalizer:unit markdown:contributing -->
## 贡献指南

有关开发环境搭建与贡献指南，请参见 [CONTRIBUTING.md](../../CONTRIBUTING.md)。所有贡献均需要 DCO 签署。
