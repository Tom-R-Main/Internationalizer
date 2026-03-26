<p align="center">
  <img src="assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

面向软件项目的 AI 原生国际化流水线。使用 LLM 翻译、验证和管理 i18n 文件。

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](LICENSE)

## 为什么选择 Internationalizer？

大多数 i18n 工具要么是运行时库（i18next、react-intl），要么是键值管理 SaaS 平台（Crowdin、Lokalise）。它们都没有很好地解决实际的翻译问题：

- **人工翻译** 在语言数量增加后难以扩展
- **机器翻译 API**（Google Translate、DeepL）会忽略你的术语、语调和 UI 约定
- **通用 LLM 翻译** 效果更好，但如果没有术语表和样式指南，翻译结果会不一致

Internationalizer 则不同。它是一个 **CLI 流水线**，将 LLM 翻译与以下功能结合：

- **单语言术语表** — 确保整个应用中的术语一致
- **单语言样式指南** — 控制语调、正式程度、复数形式和排版
- **翻译记忆库** — 跳过未更改的字符串，节省 API 调用成本
- **键值验证** — 在发布前捕获缺失的翻译和插值不匹配问题

## 安装

通过 npm 安装：

```bash
npm install -g internationalizer
```

或者不进行全局安装直接运行：

```bash
npx internationalizer --help
```

npm 包会通过特定平台的预构建可选依赖项，从 npm 安装匹配的预构建二进制文件。

通过 Go 安装：

```bash
go install github.com/Tom-R-Main/Internationalizer/cmd/internationalizer@latest
```

或者从源码构建：

```bash
git clone https://github.com/Tom-R-Main/Internationalizer.git
cd Internationalizer
go build -o internationalizer ./cmd/internationalizer
```

## npm 包

- Git 标签和 npm 包版本必须匹配，例如 `v0.1.0` 和 `0.1.0`
- 根 `internationalizer` 包依赖于平台包，例如 `internationalizer-darwin-arm64`
- 支持的 npm 目标平台：macOS arm64/x64、Linux arm64/x64、Windows x64
- CI 发布需要一个名为 `NPM_TOKEN` 的 GitHub secret

## 快速开始

1. 在项目根目录创建一个配置文件：

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

2. 设置你的 API 密钥：

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

3. 预览将要翻译的内容：

```bash
internationalizer translate --dry-run
```

4. 运行翻译：

```bash
internationalizer translate
```

5. 验证所有语言环境：

```bash
internationalizer validate
```

## 命令

### `translate`

查找缺失的键值并通过 LLM 进行翻译。

```bash
internationalizer translate                    # 翻译所有语言环境
internationalizer translate -l fr              # 仅翻译法语
internationalizer translate --dry-run          # 预览但不调用 API
internationalizer translate --batch-size 20    # 较小的批处理大小
internationalizer translate --concurrency 2    # 较少的并发调用
```

### `validate`

检查所有语言环境文件是否存在缺失键、多余键以及插值不匹配的问题。

```bash
internationalizer validate                     # 人类可读的输出
internationalizer validate --json              # 机器可读的 JSON
internationalizer validate -q                  # 仅返回退出码
```

### `detect`

自动检测 i18n 框架并建议配置。

```bash
internationalizer detect
```

支持：react-i18next、next-intl、vue-i18n、原生 JSON、Markdown 文档。

### `glossary`

管理在翻译过程中强制执行的单语言术语表。

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

管理翻译记忆库（先前翻译字符串的 JSONL 缓存）。

```bash
internationalizer tm stats                     # 显示记录数
internationalizer tm export                    # 导出为 JSON
internationalizer tm clear --force             # 删除所有记录
```

## 配置参考

```yaml
# .internationalizer.yml

# 源语言（默认：en）
source_locale: en

# 目标翻译语言（必填）
target_locales: [fr, de, es, ja, zh-CN, ar]

# 源语言环境文件路径（必填）
source_path: locales/en.json

# LLM 提供商设置
llm:
  # 提供商："anthropic"、"openai"、"gemini" 或 "openrouter"（默认：gemini）
  provider: gemini

  # 各提供商的默认模型名称：
  #   anthropic:  claude-sonnet-4-6
  #   openai:     gpt-5.4
  #   gemini:     gemini-3.1-pro-preview
  #   openrouter: google/gemini-3-flash-preview
  model: gemini-3.1-pro-preview

  # 包含 API 密钥的环境变量
  api_key_env: GOOGLE_AI_STUDIO_API_KEY

  # 兼容 OpenAI 的端点基础 URL（可选）
  # base_url: https://api.openai.com

# 每次 LLM 调用的键数量（默认：40）
batch_size: 40

# 并行 LLM 调用数（默认：4）
concurrency: 4

# 包含各语言环境样式指南 Markdown 文件的目录（默认：style-guides）
style_guides_dir: style-guides

# 包含各语言环境术语表 JSON 文件的目录（默认：glossary）
glossary_dir: glossary

# 翻译记忆库文件路径（默认：.internationalizer/tm.jsonl）
tm_path: .internationalizer/tm.jsonl
```

## 样式指南

样式指南是注入到 LLM 翻译提示词中的 Markdown 文件。它们控制语调、正式程度、排版以及其他特定语言的约定。

```
style-guides/
  _conventions.md    # 适用于所有语言的共享规则
  fr.md              # 法语特定规则
  ja.md              # 日语特定规则
  ar.md              # 阿拉伯语特定规则
```

### 共享约定 (`_conventions.md`)

定义适用于所有语言的规则：插值语法、HTML 保留、字符串类型约定（按钮 vs 标签 vs 错误）等。

### 单语言指南 (`{locale}.md`)

定义特定语言的规则：正式程度（如 tu 与 vous）、标点符号（如法文引号、倒问号）、复数形式、日期/数字格式以及术语表。

有关实际示例，请参阅 [`examples/react-app/style-guides/`](examples/react-app/style-guides/)。

## 术语表格式

术语表文件是存储在 `{glossary_dir}/{locale}.json` 中的 JSON 数组：

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

术语会作为术语表注入到 LLM 提示词中，以确保整个应用程序中关键术语的翻译保持一致。

## 翻译记忆库

翻译记忆库存储为 JSONL 文件（每行一个 JSON 记录）。每条记录包含：

- 源键和值
- 翻译后的值
- 源值的 SHA-256 哈希
- 时间戳

在后续运行中，未更改的字符串将直接从 TM 缓存中提供，而无需调用 LLM，从而节省时间和 API 成本。TM 文件对 Git 友好，可以与你的语言环境文件一起提交。

## 支持的格式

| 格式 | 扩展名 | 模式 |
|--------|-----------|------|
| JSON | `.json` | 键值对（嵌套、点号表示法扁平化） |
| YAML | `.yml`, `.yaml` | 键值对（保留注释和顺序） |
| Markdown | `.md`, `.mdx` | 全文档翻译 |

## 项目类型检测

`internationalizer detect` 通过检查以下内容来识别你的 i18n 设置：

- `package.json` 中对 react-i18next、next-intl 或 vue-i18n 的依赖
- 匹配常见语言环境模式的目录结构
- 文件扩展名和命名约定

## 架构

```
cmd/internationalizer/     CLI 入口点和命令定义
internal/
  config/                  带默认值的 YAML 配置加载
  detect/                  项目类型自动检测
  formats/                 格式解析器（JSON、YAML、Markdown）
  glossary/                单语言术语表管理
  llm/                     LLM 提供商接口 + 实现
    anthropic.go           Anthropic Claude 后端
    openai.go              OpenAI / 兼容后端
    gemini.go              通过 AI Studio 的 Google Gemini 后端
                           OpenRouter 使用带有自定义 base_url 的 openai.go
  styleguide/              样式指南加载器
  tm/                      JSONL 翻译记忆库
  translate/               翻译编排器
  validate/                语言环境验证和差异对比
```

## 与替代方案的比较

| 功能 | Internationalizer | i18next | Crowdin | 通用 LLM |
|---------|------------------|---------|---------|-------------|
| LLM 驱动翻译 | 是 | 否 | 部分 | 是 |
| 单语言样式指南 | 是 | 否 | 否 | 否 |
| 强制执行术语表 | 是 | 否 | 是 | 否 |
| 翻译记忆库 | 是 | 否 | 是 | 否 |
| CLI / 本地执行 | 是 | 不适用 | 否 | 手动 |
| Git 友好文件 | 是 | 是 | 部分 | 手动 |
| 无 SaaS 依赖 | 是 | 是 | 否 | 视情况而定 |
| 开源 (AGPL-3.0) | 是 | 是 | 否 | 视情况而定 |

## 许可证

[AGPL-3.0](LICENSE)

## 贡献

有关开发设置和指南，请参阅 [CONTRIBUTING.md](CONTRIBUTING.md)。所有贡献都需要 DCO 签名。

