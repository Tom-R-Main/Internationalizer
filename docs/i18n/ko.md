> [English (original)](../../README.md)

<p align="center">
  <img src="../../assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

소프트웨어 프로젝트를 위한 AI 네이티브 국제화 파이프라인입니다. LLM을 사용하여 i18n 파일을 번역, 검증 및 관리합니다.

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](LICENSE)

<p align="center">
<a href="docs/i18n/ar.md">العربية</a> · <a href="docs/i18n/bn.md">বাংলা</a> · <a href="docs/i18n/de.md">Deutsch</a> · <a href="docs/i18n/es.md">Español</a> · <a href="docs/i18n/fr.md">Français</a> · <a href="docs/i18n/hi.md">हिन्दी</a> · <a href="docs/i18n/id.md">Indonesia</a> · <a href="docs/i18n/ja.md">日本語</a> · <a href="docs/i18n/pa.md">ਪੰਜਾਬੀ</a> · <a href="docs/i18n/pt-BR.md">Português</a> · <a href="docs/i18n/ru.md">Русский</a> · <a href="docs/i18n/te.md">తెలుగు</a><br><a href="docs/i18n/th.md">ไทย</a> · <a href="docs/i18n/uk.md">Українська</a> · <a href="docs/i18n/yue.md">粵語</a> · <a href="docs/i18n/zh-CN.md">简体中文</a> · <a href="docs/i18n/zh-TW.md">繁體中文</a>
</p>

---

## 왜 Internationalizer를 사용해야 할까요?

대부분의 i18n 도구는 런타임 라이브러리(i18next, react-intl)이거나 키 관리 SaaS 플랫폼(Crowdin, Lokalise)입니다. 하지만 이들 중 어느 것도 실제 번역 문제를 제대로 해결하지 못합니다.

- **수동 번역**은 지원 언어가 몇 개만 넘어가도 확장하기 어렵습니다.
- **기계 번역 API**(Google Translate, DeepL)는 용어, 어조 및 UI 규칙을 무시합니다.
- **일반적인 LLM 번역**은 더 나은 결과를 보여주지만, 용어집과 스타일 가이드가 없으면 일관성 없는 결과를 얻게 됩니다.

Internationalizer는 다릅니다. LLM 번역과 다음 기능을 결합한 **CLI 파이프라인**입니다.

- **언어별 용어집** — 앱 전체에서 일관된 용어를 적용합니다.
- **언어별 스타일 가이드** — 어조, 격식, 복수형 및 타이포그래피를 제어합니다.
- **번역 메모리** — 변경되지 않은 문자열을 건너뛰어 API 호출 비용을 절약합니다.
- **키 검증** — 배포 전에 누락된 번역과 보간 불일치를 찾아냅니다.

## 설치

npm에서 설치합니다.

```bash
npm install -g internationalizer
```

또는 전역 설치 없이 실행합니다.

```bash
npx internationalizer --help
```

npm 패키지는 플랫폼별 선택적 종속성을 통해 npm에서 일치하는 사전 빌드된 바이너리를 설치합니다.

Go를 사용하여 설치합니다.

```bash
go install github.com/Tom-R-Main/Internationalizer/cmd/internationalizer@latest
```

또는 소스에서 빌드합니다.

```bash
git clone https://github.com/Tom-R-Main/Internationalizer.git
cd Internationalizer
go build -o internationalizer ./cmd/internationalizer
```

## npm 패키지

- Git 태그와 npm 패키지 버전은 일치해야 합니다. (예: `v0.1.0` 및 `0.1.0`)
- 루트 `internationalizer` 패키지는 `internationalizer-darwin-arm64`와 같은 플랫폼 패키지에 종속됩니다.
- 지원되는 npm 대상: macOS arm64/x64, Linux arm64/x64, Windows x64
- CI 게시에 `NPM_TOKEN`이라는 GitHub 시크릿이 필요합니다.

## 빠른 시작

1. 프로젝트 루트에 구성 파일을 생성합니다.

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

2. API 키를 설정합니다.

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

3. 번역될 내용을 미리 봅니다.

```bash
internationalizer translate --dry-run
```

4. 번역을 실행합니다.

```bash
internationalizer translate
```

5. 모든 로캘을 검증합니다.

```bash
internationalizer validate
```

## 명령어

### `translate`

누락된 키를 찾아 LLM을 통해 번역합니다.

```bash
internationalizer translate                    # translate all locales
internationalizer translate -l fr              # translate French only
internationalizer translate --dry-run          # preview without API calls
internationalizer translate --batch-size 20    # smaller batches
internationalizer translate --concurrency 2    # fewer parallel calls
```

### `validate`

모든 로캘 파일에서 누락된 키, 추가된 키 및 보간 불일치가 있는지 확인합니다.

```bash
internationalizer validate                     # human-readable output
internationalizer validate --json              # machine-readable JSON
internationalizer validate -q                  # exit code only
```

### `detect`

i18n 프레임워크를 자동 감지하고 구성을 제안합니다.

```bash
internationalizer detect
```

지원 항목: react-i18next, next-intl, vue-i18n, 일반 JSON, 마크다운 문서.

### `glossary`

번역 중에 적용되는 언어별 용어집을 관리합니다.

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

번역 메모리(이전에 번역된 문자열의 JSONL 캐시)를 관리합니다.

```bash
internationalizer tm stats                     # show record counts
internationalizer tm export                    # dump as JSON
internationalizer tm clear --force             # delete all records
```

## 구성 참조

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

## 스타일 가이드

스타일 가이드는 LLM 번역 프롬프트에 주입되는 마크다운 파일입니다. 어조, 격식, 타이포그래피 및 기타 언어별 규칙을 제어합니다.

```
style-guides/
  _conventions.md    # shared rules for all languages
  fr.md              # French-specific rules
  ja.md              # Japanese-specific rules
  ar.md              # Arabic-specific rules
```

### 공통 규칙 (`_conventions.md`)

모든 언어에 적용되는 규칙을 정의합니다. 보간 구문, HTML 보존, 문자열 유형 규칙(버튼, 레이블, 오류 등)이 포함됩니다.

### 언어별 가이드 (`{locale}.md`)

언어별 규칙을 정의합니다. 격식 수준(예: 반말과 존댓말), 문장 부호(예: 겹화살괄호, 역물음표), 복수형, 날짜/숫자 형식 및 용어집이 포함됩니다.

실제 작동 예시는 [`examples/react-app/style-guides/`](examples/react-app/style-guides/)를 참조하세요.

## 용어집 형식

용어집 파일은 `{glossary_dir}/{locale}.json`에 저장되는 JSON 배열입니다.

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

용어는 용어 표 형태로 LLM 프롬프트에 주입되어 애플리케이션 전체에서 주요 용어가 일관되게 번역되도록 보장합니다.

## 번역 메모리

번역 메모리는 JSONL 파일(줄당 하나의 JSON 레코드)로 저장됩니다. 각 레코드에는 다음이 포함됩니다.

- 소스 키 및 값
- 번역된 값
- 소스 값의 SHA-256 해시
- 타임스탬프

이후 실행 시 변경되지 않은 문자열은 LLM을 호출하지 않고 TM 캐시에서 제공되므로 시간과 API 비용이 모두 절약됩니다. TM 파일은 Git과 호환되며 로캘 파일과 함께 커밋할 수 있습니다.

## 지원되는 형식

| 형식 | 확장자 | 모드 |
|--------|-----------|------|
| JSON | `.json` | 키-값 (중첩, 점 표기법 평탄화) |
| YAML | `.yml`, `.yaml` | 키-값 (주석 및 순서 보존) |
| Markdown | `.md`, `.mdx` | 전체 문서 번역 |

## 프로젝트 유형 감지

`internationalizer detect`는 다음을 확인하여 i18n 설정을 식별합니다.

- react-i18next, next-intl 또는 vue-i18n에 대한 `package.json` 종속성
- 일반적인 로캘 패턴과 일치하는 디렉터리 구조
- 파일 확장자 및 명명 규칙

## 아키텍처

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

## 대안과의 비교

| 기능 | Internationalizer | i18next | Crowdin | 일반 LLM |
|---------|------------------|---------|---------|-------------|
| LLM 기반 번역 | 예 | 아니요 | 부분적 | 예 |
| 언어별 스타일 가이드 | 예 | 아니요 | 아니요 | 아니요 |
| 용어집 적용 | 예 | 아니요 | 예 | 아니요 |
| 번역 메모리 | 예 | 아니요 | 예 | 아니요 |
| CLI / 로컬 실행 | 예 | 해당 없음 | 아니요 | 수동 |
| Git 호환 파일 | 예 | 예 | 부분적 | 수동 |
| SaaS 종속성 없음 | 예 | 예 | 아니요 | 다양함 |
| 오픈 소스 (AGPL-3.0) | 예 | 예 | 아니요 | 다양함 |

## 라이선스

[AGPL-3.0](LICENSE)

## 기여

개발 설정 및 가이드라인은 [CONTRIBUTING.md](CONTRIBUTING.md)를 참조하세요. 모든 기여에는 DCO 서명이 필요합니다.

