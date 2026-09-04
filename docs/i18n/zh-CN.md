> [English (original)](../../README.md)

<p align="center">
  <img src="../../assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

面向软件项目的 AI 原生国际化流水线。基于大语言模型（LLM）实现 i18n 文件的翻译、校验与管理。

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](../../LICENSE)

<p align="center">
<a href="ar.md">العربية</a> · <a href="bn.md">বাংলা</a> · <a href="cs.md">Čeština</a> · <a href="da.md">Dansk</a> · <a href="de.md">Deutsch</a> · <a href="el.md">Ελληνικά</a> · <a href="es.md">Español</a> · <a href="fi.md">Suomi</a> · <a href="fr.md">Français</a> · <a href="he.md">עברית</a> · <a href="hi.md">हिन्दी</a> · <a href="id.md">Indonesia</a> · <a href="it.md">Italiano</a> · <a href="ja.md">日本語</a> · <a href="ko.md">한국어</a> · <a href="ms.md">Bahasa Melayu</a><br><a href="nl.md">Nederlands</a> · <a href="pa.md">ਪੰਜਾਬੀ</a> · <a href="pl.md">Polski</a> · <a href="pt-BR.md">Português</a> · <a href="ro.md">Română</a> · <a href="ru.md">Русский</a> · <a href="sv.md">Svenska</a> · <a href="te.md">తెలుగు</a> · <a href="th.md">ไทย</a> · <a href="tr.md">Türkçe</a> · <a href="uk.md">Українська</a> · <a href="vi.md">Tiếng Việt</a> · <a href="yue.md">粵語</a> · <a href="zh-CN.md">简体中文</a> · <a href="zh-TW.md">繁體中文</a>
</p>

---

## 为什么选择 Internationalizer？

多数 i18n 工具要么属于运行时库（如 i18next、react-intl），要么属于翻译键管理 SaaS 平台（如 Crowdin、Lokalise）。它们均未能从根本上解决实际翻译问题：

- **人工翻译**：支持语言一旦增多便难以扩展
- **机器翻译 API**（Google Translate、DeepL）：无视术语库、语气要求与 UI 规范
- **通用 LLM 翻译**：效果虽有提升，但缺少术语表与样式指南时，译文风格容易脱节

Internationalizer 则另辟蹊径。它是一套 **CLI 流水线**，将 LLM 翻译与以下能力深度结合：

- **各语言术语表** — 确保全系统专业术语统一
- **各语言风格指南** — 规范语气、正式程度、复数规则及排版格式
- **翻译记忆库** — 自动跳过未修改字段，节省 API 调用开销
- **确定性校验** — 在发布前发现翻译键缺失或多余、受保护结构变更、术语表违规，以及复数或 ICU 错误

## 安装

通过 npm 安装：

```bash
npm install -g internationalizer
```

无需全局安装直接运行：

```bash
npx internationalizer --help
```

npm 包会借助对应平台的可选依赖，从 npm 拉取并安装匹配的预构建二进制文件。

使用 Go 安装：

```bash
go install github.com/Tom-R-Main/Internationalizer/cmd/internationalizer@latest
```

从源码编译：

```bash
git clone https://github.com/Tom-R-Main/Internationalizer.git
cd Internationalizer
go build -o internationalizer ./cmd/internationalizer
```

## npm 软件包

- Git 标签版本必须与 npm 包版本严格一致，例如 `v0.1.0` 与 `0.1.0`
- 顶层 `internationalizer` 包依赖各平台分发包，例如 `internationalizer-darwin-arm64`
- 支持的 npm 目标平台：macOS arm64/x64、Linux arm64/x64、Windows x64
- CI 发布流水线需在 GitHub Secret 中配置 `NPM_TOKEN`

## 快速上手

1. 在项目根目录下创建配置文件：

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

2. 配置 API 密钥：

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

3. 预览待翻译内容：

```bash
internationalizer translate --dry-run
```

4. 执行翻译：

```bash
internationalizer translate
```

5. 校验所有语言环境文件：

```bash
internationalizer validate
```

## 命令列表

### `translate`

检索缺失或已过期的键，并通过 LLM 自动完成翻译。

```bash
internationalizer translate                    # 翻译全部语言环境
internationalizer translate -l fr              # 仅翻译法语
internationalizer translate --dry-run          # 预览变更，不产生 API 调用
internationalizer translate --adopt-existing   # 将现有翻译纳为基线，不产生 API 调用
internationalizer translate --refresh-policy   # 重新翻译因提示词/样式指南/模型变更而过期的条目
internationalizer translate --batch-size 20    # 减小批处理大小
internationalizer translate --concurrency 2    # 降低并发请求数
```

翻译状态会分别独立跟踪缺失（missing）、源文过期（source-stale）、策略过期（policy-stale）、最新（current）以及手动修改（manually edited）等情况，手动编辑不会掩盖源文本或策略的变动。策略过期的词条会被列出，但仅在传入 `--refresh-policy` 时才会触发重新翻译。系统绝不会自动覆盖手动编辑的内容。初次为已审校的译文引入清单文件，或明确将人工修改确认采纳为新基线时，请使用 `--adopt-existing`。

### `validate`

将所有语言环境文件与源文件包进行比对。默认校验会计算必需目标翻译键的覆盖率，将多余翻译键列为警告，并在翻译键缺失、插值变量不匹配或 ICU MessageFormat 结构无效时失败。

```bash
internationalizer validate                     # 人类可读文本输出
internationalizer validate --json              # 机器可读 JSON 输出
internationalizer validate -q                  # 静默模式，仅返回退出码
internationalizer validate --strict             # 强制执行翻译质量规则
internationalizer validate --require-state      # 要求清单中的溯源状态为最新
```

`--strict` 还会报告实际翻译覆盖率。语言内容若与源文完全相同，会被视为未翻译；只有术语表为完整值明确配置了源文与译文相同的精确条目时才会豁免。校验会遵守 `ignore_case`，但较长文本中仅包含某个术语表词条并不能获得豁免。严格模式会在存在多余翻译键、与源文相同的值、插值／HTML／代码／Markdown 链接结构变更、术语表违规或缺少已配置复数形式时失败。

`--require-state` 会将每个目标值与 `.internationalizer.lock` 比对。翻译键未被跟踪，或清单记录的源文、翻译策略、目标哈希已过期时，校验都会失败。该选项可与 `--strict` 同时使用。

人工可读报告与 JSON 报告使用以下稳定的问题代码：

| 代码 | 含义 |
| --- | --- |
| `missing_key` / `extra_key` | 源文件与目标文件的翻译键集合不一致 |
| `blank_translation` | 非空源文对应的严格模式译文为空 |
| `source_identical` | 严格模式下语言内容仍与源文相同 |
| `protected_structure_mismatch` | 插值、HTML、代码或链接结构发生变化 |
| `glossary_violation` | 未找到获准使用的目标术语或变体 |
| `plural_form_missing` | 缺少为该语言配置的复数形式 |
| `icu_message_syntax` | 源文或译文的 ICU 消息格式有误 |
| `icu_argument_mismatch` | ICU 参数名、参数类型或格式化样式不一致 |
| `icu_selector_mismatch` | 选择器不一致，或复数类别不适用于目标语言 |
| `untracked` | 清单中没有对应的目标记录 |
| `source_stale` | 清单记录生成后源文发生变化 |
| `policy_stale` | 生成提示词或模型设置发生变化 |
| `target_modified` | 目标内容与清单记录不一致 |

### `detect`

自动检测项目使用的 i18n 框架并生成推荐配置。

```bash
internationalizer detect
```

现已支持：react-i18next、next-intl、vue-i18n、纯 JSON 及 Markdown 文档。

### `glossary`

管理在翻译期间强制遵循的单语言术语表。

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

管理翻译记忆库（基于 JSONL 存储的历史翻译缓存）。

```bash
internationalizer tm stats                     # 查看记录统计
internationalizer tm export                    # 导出为 JSON 格式
internationalizer tm clear --force             # 清空全部记录
```

## 配置参考

```yaml
# .internationalizer.yml

# 源语言（默认：en）
source_locale: en

# 目标翻译语言列表（必填）
target_locales: [fr, de, es, ja, yue, zh-CN, zh-TW, ar]

# 一个或多个源到目标的映射绑定（必填）。
# 占位符 {locale} 会被自动替换为各个已配置的目标语言标识。
bundles:
  - id: app
    source: locales/en.json
    target: locales/{locale}.json
    format: json
  - id: docs
    source: README.md
    target: docs/i18n/{locale}.md
    format: markdown

# 向后兼容配置：source_path 仍可将目标映射至同级文件，
# 例如 locales/fr.json。新项目建议优先使用 bundles。
# source_path: locales/en.json

# LLM 服务商配置
llm:
  # 服务商选项："anthropic"、"openai"、"gemini" 或 "openrouter"（默认：gemini）
  provider: gemini

  # 各服务商的默认模型：
  #   anthropic:  claude-opus-5
  #   openai:     gpt-5.6-luna（思考力度默认设为 max）
  #   gemini:     gemini-3.8-flash
  #   openrouter: deepseek/deepseek-v4-pro-0813
  model: gemini-3.8-flash

  # 存储 API 密钥的环境变量名称
  api_key_env: GOOGLE_AI_STUDIO_API_KEY

  # OpenAI 兼容接口的 Base URL（可选）
  # base_url: https://api.openai.com

  # OpenAI GPT-5 系列 Responses API 的思考力度（reasoning effort）
  #（OpenAI 服务商默认值：max）
  reasoning_effort: max

  # 针对特定目标语言的独立 LLM 设置（可选）。覆盖配置若沿用全局服务商，
  # 会自动继承全局中未显式指定的参数；若切换为其他服务商，
  # 未指定项则回退至该服务商的默认设置。
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

# 每次 LLM 请求包含的键数量（默认：40）
batch_size: 40

# 并行 LLM 请求数（默认：4）
concurrency: 4

# 各语言风格指南 Markdown 文件存放目录（默认：style-guides）
style_guides_dir: style-guides

# 各语言术语表 JSON 文件存放目录（默认：glossary）
glossary_dir: glossary

# 翻译记忆库文件路径（默认：.internationalizer/tm.jsonl）
tm_path: .internationalizer/tm.jsonl

# 记录源文、策略、目标及版本溯源状态的清单文件
#（默认：.internationalizer.lock；请将此文件提交至版本控制）
manifest_path: .internationalizer.lock

# 可选的翻译及严格校验规则
validation:
  plural_style: i18next-v4 # 生成并校验目标语言的复数形式
```

语言标识必须是格式正确的 BCP 47 标签，例如 `fr`、`pt-BR` 或 `sr-Latn-RS`。规范化后等价的目标语言会被判定为重复项；各语言专属的服务商覆写配置也会按规范化形式匹配。上述示例中，日语等没有单独覆写配置的语言均继承全局 Gemini 设置。

ICU MessageFormat 内容会按结构解析。系统支持简单参数、`select`、`plural`、`selectordinal`、`number`、`date` 和 `time`，并支持嵌套消息、复数偏移量、精确数字选择器以及 `#`。校验会检查语法、参数类型与格式化样式、复数偏移量、select 分支标识，以及目标语言适用的 CLDR 复数类别。若服务商返回的内容破坏这些约束，系统会在写入语言文件或翻译记忆库前拒绝该结果。

启用 `i18next-v4` 后，系统会在翻译期间将识别到的源复数词族扩展为目标语言所需的 CLDR 类别。仅目标语言需要的类别会使用源词族的 `_other` 值作为翻译模板。严格校验要求目标语言所需的类别齐全，但不会强制保留目标语言不使用的源语言专属类别。

## 风格指南

风格指南是以 Markdown 格式编写的文件，会在翻译时直接注入 LLM 提示词中。借此可精准把控译文的语调、正式程度、排版规范以及特定语言的行文习惯。

```
style-guides/
  _conventions.md    # 面向所有语言的通用规范
  fr.md              # 法语专属规则
  ja.md              # 日语专属规则
  ar.md              # 阿拉伯语专属规则
```

### 通用规范（`_conventions.md`）

用于定义适用于所有语言的全局规则，包括：插值语法、HTML 标签保护策略、不同文案类型的表述规范（如按钮、表单标签与错误提示）等。

### 各语言指南（`{locale}.md`）

用于定义特定语言的细分规则，包括：正式程度（如法语 tu 与 vous）、特定标点规范（如法文引号 « »、西班牙语倒问号 ¿）、复数形式定义、日期/数字格式化及特定术语表。

完整可运行范例参见 [`examples/react-app/style-guides/`](../../examples/react-app/style-guides/)。

## 术语表格式

术语表以 JSON 数组形式保存在 `{glossary_dir}/{locale}.json` 中：

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

`variants` 用于列出其他获准使用的译法。`enforcement` 可设为 `error` 或 `warning`；省略时默认为 `error`。这些词条会作为术语对照表注入 LLM 提示词，确保整个应用中的术语保持一致。对于 `{"source":"API","target":"API"}` 这类源文与译文完全相同的精确条目，严格校验不会将完整值判定为未翻译；较长文本中仅包含 `API` 则不会获得豁免。

## 翻译记忆库

翻译记忆库采用 JSONL 文件格式存储（每行对应一条独立 JSON 记录）。各记录包含以下字段：

- 文件包、翻译键、源文、译文及规范化后的目标语言
- 源文哈希与翻译策略哈希
- 生成译文的服务商与模型
- 时间戳

后续翻译时，源文哈希与策略哈希均相同的文本会直接从缓存复用，无需调用 LLM。默认路径位于已被忽略的 `.internationalizer/` 目录中，因此它是本地缓存。如需由项目成员共享翻译记忆库，请将 `tm_path` 明确设为受版本控制的路径。可审查的 `.internationalizer.lock` 清单则单独纳入版本控制。

## 支持的文件格式

| 格式 | 文件扩展名 | 处理模式 |
|--------|-----------|------|
| JSON | `.json` | 键值对（支持深层嵌套，按点号表示法扁平化） |
| YAML | `.yml`、`.yaml` | 键值对（完整保留注释与键序） |
| Markdown | `.md`、`.mdx` | 整篇文档全文翻译 |

## 项目类型检测机制

执行 `internationalizer detect` 时，工具将通过检查以下特征定位您的 i18n 环境：

- `package.json` 中的 react-i18next、next-intl 或 vue-i18n 依赖
- 符合通用规范的语言目录架构
- 文件扩展名与命名约定

## 系统架构

```
cmd/internationalizer/     CLI 入口点与命令定义
internal/
  config/                  YAML 配置加载与默认值填充
  detect/                  项目类型自动检测
  formats/                 格式解析器（JSON、YAML、Markdown）
  glossary/                单语言术语表管理
  llm/                     LLM 服务商通用接口与后端实现
    anthropic.go           Anthropic Claude 后端
    openai.go              OpenAI 及兼容后端
    gemini.go              基于 Google AI Studio 的 Gemini 后端
                           OpenRouter 复用 openai.go 并指定自定义 base_url
  locale/                  BCP 47 语言标识与 CLDR 复数类别
  message/                 ICU MessageFormat 解析器与结构比对
  policy/                  稳定的翻译策略哈希
  state/                   受版本控制的翻译清单
  styleguide/              风格指南加载器
  tm/                      JSONL 翻译记忆库实现
  translate/               翻译流水线核心编排器
  validate/                语言包校验与差异分析
```

## 与其他方案对比

| 功能特性 | Internationalizer | i18next | Crowdin | 通用 LLM |
|---------|------------------|---------|---------|-------------|
| 基于 LLM 驱动翻译 | 是 | 否 | 部分 | 是 |
| 支持各语言风格指南 | 是 | 否 | 否 | 否 |
| 强制约束专业术语表 | 是 | 否 | 是 | 否 |
| 内置翻译记忆库 | 是 | 否 | 是 | 否 |
| 支持 CLI / 本地直接执行 | 是 | 不适用 | 否 | 需手动操作 |
| 生成 Git 友好文件 | 是 | 是 | 部分 | 需手动操作 |
| 摆脱商业 SaaS 依赖 | 是 | 是 | 否 | 视情况而定 |
| 完全开源（AGPL-3.0） | 是 | 是 | 否 | 视情况而定 |

## 许可证

本项目遵循 [AGPL-3.0](../../LICENSE) 协议开源。

第三方依赖声明请参阅 [THIRD_PARTY_NOTICES.md](../../THIRD_PARTY_NOTICES.md)。

## 参与贡献

本地开发环境搭建与贡献指引请参阅 [CONTRIBUTING.md](../../CONTRIBUTING.md)。所有代码提交均需包含 DCO 签署。
