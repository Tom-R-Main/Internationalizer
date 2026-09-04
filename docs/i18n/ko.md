> [English (original)](../../README.md)

<p align="center">
  <img src="../../assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

소프트웨어 프로젝트를 위한 AI 네이티브 국제화 파이프라인입니다. LLM을 활용하여 i18n 파일을 번역, 검증 및 관리합니다.

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](../../LICENSE)

<p align="center">
<a href="ar.md">العربية</a> · <a href="bn.md">বাংলা</a> · <a href="cs.md">Čeština</a> · <a href="da.md">Dansk</a> · <a href="de.md">Deutsch</a> · <a href="el.md">Ελληνικά</a> · <a href="es.md">Español</a> · <a href="fi.md">Suomi</a> · <a href="fr.md">Français</a> · <a href="he.md">עברית</a> · <a href="hi.md">हिन्दी</a> · <a href="id.md">Indonesia</a> · <a href="it.md">Italiano</a> · <a href="ja.md">日本語</a> · <a href="ko.md">한국어</a> · <a href="ms.md">Bahasa Melayu</a><br><a href="nl.md">Nederlands</a> · <a href="pa.md">ਪੰਜਾਬੀ</a> · <a href="pl.md">Polski</a> · <a href="pt-BR.md">Português</a> · <a href="ro.md">Română</a> · <a href="ru.md">Русский</a> · <a href="sv.md">Svenska</a> · <a href="te.md">తెలుగు</a> · <a href="th.md">ไทย</a> · <a href="tr.md">Türkçe</a> · <a href="uk.md">Українська</a> · <a href="vi.md">Tiếng Việt</a> · <a href="yue.md">粵語</a> · <a href="zh-CN.md">简体中文</a> · <a href="zh-TW.md">繁體中文</a>
</p>

---

<!-- internationalizer:unit markdown:why-internationalizer -->
## 왜 Internationalizer인가요?

대부분의 i18n 도구는 런타임 라이브러리(i18next, react-intl)이거나 키 관리 SaaS 플랫폼(Crowdin, Lokalise)입니다. 두 유형 모두 실제 번역 문제를 만족스럽게 해결하지 못합니다.

- **수동 번역**은 지원 언어가 몇 개만 넘어가도 확장할 수 없습니다.
- **기계 번역 API**(Google Translate, DeepL)는 제품 전용 용어, 어조, UI 규칙을 반영하지 못합니다.
- **일반 LLM 번역**은 상대적으로 더 우수하지만, 용어집과 스타일 가이드가 없으면 일관성 없는 결과를 산출합니다.

Internationalizer는 다릅니다. LLM 번역에 다음 핵심 기능을 결합한 **CLI 파이프라인**입니다.

- **언어별 용어집** — 애플리케이션 전반에서 일관된 전문 용어를 강제합니다.
- **언어별 스타일 가이드** — 어조, 격식 수준, 복수형 처리, 타이포그래피 규칙을 제어합니다.
- **번역 메모리** — 변경되지 않은 문자열을 건너뛰어 API 호출 비용을 절감합니다.
- **결정론적 검증** — 누락되거나 불필요한 키, 보호된 구조 변형, 용어집 위반, 복수형 및 ICU 구문 오류를 배포 전에 선제적으로 감지합니다.

<!-- internationalizer:unit markdown:installation -->
## 설치

npm을 통해 설치할 수 있습니다.

```bash
npm install -g internationalizer
```

전역 설치 없이 바로 실행할 수도 있습니다.

```bash
npx internationalizer --help
```

npm 패키지는 플랫폼별 선택적 종속성(optional dependencies)을 통해 환경에 맞는 사전 빌드 바이너리를 설치합니다.

Go를 사용하여 설치할 수도 있습니다.

```bash
go install github.com/Tom-R-Main/Internationalizer/cmd/internationalizer@latest
```

소스 코드에서 직접 빌드하는 방법도 지원합니다.

```bash
git clone https://github.com/Tom-R-Main/Internationalizer.git
cd Internationalizer
go build -o internationalizer ./cmd/internationalizer
```

<!-- internationalizer:unit markdown:npm-packages -->
## npm 패키지

- Git 태그와 npm 패키지 버전은 반드시 일치해야 합니다(예: `v0.1.0` 및 `0.1.0`).
- 루트 `internationalizer` 패키지는 `internationalizer-darwin-arm64`와 같은 플랫폼별 바이너리 패키지에 종속됩니다.
- 지원되는 npm 플랫폼: macOS arm64/x64, Linux arm64/x64, Windows x64
- CI 자동 배포를 구성하려면 `NPM_TOKEN`이라는 이름의 GitHub Secret이 필요합니다.

<!-- internationalizer:unit markdown:quick-start -->
## 빠른 시작

1. 프로젝트 루트에 구성 파일을 생성합니다.

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

2. API 키 환경 변수를 설정합니다.

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

3. 번역될 대상을 미리 확인합니다.

```bash
internationalizer translate --dry-run
```

4. 번역을 실행합니다.

```bash
internationalizer translate
```

5. 모든 로캘 파일을 검증합니다.

```bash
internationalizer validate
```

<!-- internationalizer:unit markdown:commands -->
## 명령어

### `translate`

누락되었거나 변경이 필요한 키를 찾아 LLM으로 번역합니다.

```bash
internationalizer translate                    # translate all locales
internationalizer translate -l fr              # translate French only
internationalizer translate --dry-run          # preview without API calls
internationalizer translate --adopt-existing   # baseline existing translations without API calls
internationalizer translate --refresh-policy   # refresh prompt/style/model-stale entries
internationalizer translate --batch-size 20    # smaller batches
internationalizer translate --concurrency 2    # fewer parallel calls
```

번역 상태는 누락(missing), 소스 변경(source-stale), 정책 변경(policy-stale), 최신 상태(current), 수동 편집(manually edited) 조건을 각각 독립적으로 추적하고 보고합니다. 따라서 수동 편집 내역이 원본 소스 변경이나 정책 변경 사항을 은폐하지 않습니다. 정책 변경으로 인해 갱신이 필요한 항목은 기본 실행 시 상태만 보고되며 `--refresh-policy` 플래그를 지정할 때만 다시 번역됩니다. 수동으로 편집된 값은 자동으로 덮어쓰지 않습니다. 이미 검토가 완료된 기존 번역본에 매니페스트를 처음 도입하거나, 수동 편집 사항을 새로운 기준선으로 명시적으로 수용할 때는 `--adopt-existing` 옵션을 사용하세요.

### `validate`

모든 로캘 파일을 소스 번들과 비교하여 검증합니다. 기본 검증 모드는 구조적 커버리지(필수 타깃 키의 존재 비율)를 확인하고, 불필요한 추가 키는 경고로 보고하며, 키 누락, 보간(interpolation) 불일치, 유효하지 않은 ICU MessageFormat 구조가 발견되면 검증 실패로 처리합니다.

```bash
internationalizer validate                     # human-readable output
internationalizer validate --json              # machine-readable JSON
internationalizer validate -q                  # exit code only
internationalizer validate --strict             # enforce translation quality rules
internationalizer validate --require-state      # require current manifest provenance
```

`--strict` 플래그는 번역 완료 커버리지도 함께 검증합니다. 소스와 완전히 동일한 언어적 텍스트는 용어집에 정확히 동일한 소스 및 타깃 텍스트가 전체 값 단위로 등록되어 있지 않는 한 미번역 항목으로 간주됩니다(`ignore_case` 설정은 반영되지만, 긴 문장 안에 해당 용어집 단어가 포함되어 있다는 사실만으로는 예외 처리가 적용되지 않습니다). 엄격(strict) 모드는 추가 키, 소스와 동일한 미번역 값, 보간/HTML/코드/Markdown 링크 구조 변형, 용어집 위반, 설정된 복수형 미준수 사항이 감지되면 즉시 실패 처리합니다.

`--require-state` 플래그는 각 타깃 파일이 `.internationalizer.lock` 상태와 일치하는지 추적 상태를 확인합니다. 매니페스트에 등록되지 않은 키가 존재하거나 기록된 소스, 번역 정책, 타깃 해시가 최신 상태가 아니면 실패 처리됩니다. 이 옵션은 `--strict`와 결합하여 실행할 수 있습니다.

사람이 읽을 수 있는 보고서와 JSON 출력 모두 안정적인 검증 코드(finding code)를 사용합니다.

| 코드 | 의미 |
| --- | --- |
| `missing_key` / `extra_key` | 소스와 타깃의 키 세트가 서로 일치하지 않습니다. |
| `blank_translation` | 비어 있지 않은 소스에 대해 엄격 모드 타깃 값이 비어 있습니다. |
| `source_identical` | 엄격 모드에서 언어적 텍스트가 번역되지 않은 채 원본과 동일하게 유지되었습니다. |
| `protected_structure_mismatch` | 보간 변수, HTML, 인라인 코드 또는 링크 구조가 원본과 다릅니다. |
| `glossary_violation` | 승인된 타깃 용어 또는 대체 표기가 사용되지 않았습니다. |
| `plural_form_missing` | 로캘에 설정된 필수 복수형이 누락되었습니다. |
| `icu_message_syntax` | 소스 또는 타깃의 ICU 메시지 문법 오류가 존재합니다. |
| `icu_argument_mismatch` | ICU 인수의 이름, 타입 또는 포맷터 스타일이 다릅니다. |
| `icu_selector_mismatch` | 선택자(selector)가 다르거나 타깃 로캘에 적합하지 않은 복수형 범주가 사용되었습니다. |
| `untracked` | 타깃 키에 대한 매니페스트 기록이 존재하지 않습니다. |
| `source_stale` | 기록된 번역 시점 이후 원본 소스 내용이 변경되었습니다. |
| `policy_stale` | 생성된 프롬프트 또는 모델 설정이 변경되었습니다. |
| `target_modified` | 타깃 콘텐츠가 매니페스트에 기록된 내용과 다릅니다. |

### `detect`

프로젝트에서 사용 중인 i18n 프레임워크를 자동으로 감지하고 최적의 구성을 제안합니다.

```bash
internationalizer detect
```

지원 프레임워크: react-i18next, next-intl, vue-i18n, 일반 JSON, Markdown 문서.

### `glossary`

번역 시 엄격하게 적용할 언어별 용어집 항목을 관리합니다.

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

이전에 번역된 문자열을 캐시하는 번역 메모리(JSONL)를 관리합니다.

```bash
internationalizer tm stats                     # show record counts
internationalizer tm export                    # dump as JSON
internationalizer tm clear --force             # delete all records
```

<!-- internationalizer:unit markdown:configuration-reference -->
## 구성 참조

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

로캘 식별자는 `fr`, `pt-BR`, `sr-Latn-RS`와 같이 표준 BCP 47 형식이어야 합니다. 정규 동등(canonical-equivalent)한 타깃 로캘은 중복으로 처리되어 거부되며, 로캘별 공급자 재정의도 정규 동등 표기 기준으로 매칭됩니다. 위 예시 설정에서 일본어를 포함해 별도의 재정의가 지정되지 않은 로캘은 전역 Gemini 설정을 상속받습니다.

ICU MessageFormat 값은 구문 구조 단위로 파싱됩니다. 일반 인수 및 `select`, `plural`, `selectordinal`, `number`, `date`, `time`이 지원되며 중첩 메시지, 복수형 오프셋, 정밀 숫자 선택자, `#` 기호 처리도 지원합니다. 검증 과정에서는 구문, 인수 타입 및 포맷터 스타일, 복수형 오프셋, select 분기 일치 여부, 타깃 로캘의 CLDR 복수형 범주 유효성을 평가합니다. 이러한 불변 조건을 위반하는 공급자 응답은 로캘 파일이나 번역 메모리에 기록되기 전에 거부됩니다.

`i18next-v4` 설정 시, 번역 과정에서 인식된 원본 소스의 복수형 그룹이 타깃 로캘에 정의된 CLDR 범주로 확장됩니다. 타깃 언어에만 존재하는 복수형 범주는 소스 그룹의 `_other` 값을 번역 템플릿으로 사용합니다. 엄격 검증 모드에서는 이러한 타깃 범주를 필수로 요구하며, 소스 언어에만 있고 타깃 언어에는 사용되지 않는 범주는 선택 사항으로 취급됩니다.

<!-- internationalizer:unit markdown:style-guides -->
## 스타일 가이드

스타일 가이드는 LLM 번역 프롬프트에 직접 주입되는 Markdown 파일입니다. 번역 어조, 격식 수준, 문장 부호 규칙, 기타 언어별 세부 규약을 제어합니다.

```
style-guides/
  _conventions.md    # shared rules for all languages
  fr.md              # French-specific rules
  ja.md              # Japanese-specific rules
  ar.md              # Arabic-specific rules
```

### 공통 규칙 (`_conventions.md`)

모든 언어에 공통으로 적용되는 규칙을 정의합니다. 보간 변수 구문, HTML 태그 보존, 문자열 유형별 규약(버튼, 레이블, 오류 메시지 등)이 포함됩니다.

### 언어별 가이드 (`{locale}.md`)

특정 언어에 특화된 규칙을 정의합니다. 대화 상대에 따른 격식 수준(예: tu vs. vous), 문장 부호(길레메 따옴표, 역물음표 등), 복수형 규칙, 날짜 및 숫자 표기법, 용어집 등이 포함됩니다.

스타일 가이드는 시스템이 생성하는 산출물이 아니라 지속적으로 관리되는 정책 입력값입니다. Internationalizer는 이 파일들을 읽기만 하며 절대 임의로 수정하지 않습니다. 스타일 가이드 내용은 용어집 및 프롬프트 규약과 별도의 해시로 관리되므로, 일반 애플리케이션 코드가 변경되더라도 기존 번역이 정책 만료 상태로 바뀌지 않습니다. 스타일 가이드를 직접 편집하면 의도적으로 해당 로캘의 정책 재검토 플래그가 설정되지만, 프롬프트 규약 버전이 올라가지 않는 한 내부 프롬프트 문구가 약간 변경되는 것만으로는 정책 만료 상태가 발생하지 않습니다.

실제 동작하는 예시는 [`examples/react-app/style-guides/`](../../examples/react-app/style-guides/) 디렉터리에서 확인하실 수 있습니다.

<!-- internationalizer:unit markdown:glossary-format -->
## 용어집 형식

용어집 파일은 `{glossary_dir}/{locale}.json` 경로에 JSON 배열 형태로 저장됩니다.

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

`variants` 필드는 타깃 언어로 허용 가능한 대체 표기 목록을 지정합니다. `enforcement` 필드는 `error` 또는 `warning`으로 지정할 수 있으며 생략 시 기본 동작인 error로 처리됩니다. 등록된 용어는 용어집 테이블 형태로 LLM 프롬프트에 주입되어 애플리케이션 전반에서 일관된 번역을 보장합니다. `{"source":"API","target":"API"}`와 같이 원본과 번역이 동일한 완전 일치 항목을 등록하면 해당 텍스트 전체가 원본과 일치하더라도 엄격 모드에서 미번역 오류로 보고되지 않는 예외 처리를 받습니다. 단, 단순히 `API`라는 단어를 문장 일부로 포함하고 있는 더 긴 문자열에는 이러한 예외가 적용되지 않습니다.

<!-- internationalizer:unit markdown:translation-memory -->
## 번역 메모리

번역 메모리는 한 줄에 하나의 JSON 레코드가 기록되는 JSONL 파일 형태로 저장됩니다. 각 레코드에는 다음 정보가 포함됩니다.

- 번들 식별자, 키, 원본 소스 값, 번역된 값, 표준 타깃 로캘 태그
- 소스, 스타일 가이드, 용어집, 프롬프트 규약 및 결합 정책 해시
- 해당 번역을 생성한 LLM 공급자 및 모델 명칭
- 타임스탬프

이후 번역을 재실행할 때 원본 소스와 정책 해시가 일치하는 문자열은 LLM API를 다시 호출하지 않고 캐시에서 즉시 가져옵니다. 기본 경로는 버전 관리에서 제외되는 `.internationalizer/` 디렉터리 하위에 위치하여 로컬 캐시 역할을 합니다. 프로젝트 팀 전체가 번역 메모리를 공유하려는 경우 `tm_path`를 버전 관리 대상 경로로 설정하세요. 검토 가능한 매니페스트 파일인 `.internationalizer.lock`은 이와 별도로 버전 관리됩니다.

<!-- internationalizer:unit markdown:supported-formats -->
## 지원되는 형식

| 형식 | 확장자 | 모드 |
| --- | --- | --- |
| JSON | `.json` | 키-값 (중첩 객체 및 점 표기법 평탄화 지원) |
| YAML | `.yml`, `.yaml` | 키-값 (주석 및 순서 보존) |
| Markdown | `.md`, `.mdx` | 프리앰블 및 H2 단위 섹션 분할 |

Markdown 타깃 파일에는 H2 섹션 직전에 보이지 않는 `internationalizer:unit` 주석이 삽입됩니다. 이러한 안정적 마커를 통해 Internationalizer는 관련 없는 다른 섹션을 다시 번역하지 않고도 변경된 소스 섹션만 개별적으로 추가, 이동, 편집할 수 있습니다. 마커가 없는 기존 문서의 경우 다음 성공적인 업데이트 시 자동으로 마커가 주입됩니다.

<!-- internationalizer:unit markdown:project-type-detection -->
## 프로젝트 유형 감지

`internationalizer detect` 명령어는 다음 항목을 확인하여 프로젝트의 i18n 구성을 자동으로 식별합니다.

- `package.json`의 react-i18next, next-intl, vue-i18n 종속성 여부
- 일반적인 로캘 패턴과 일치하는 디렉터리 구조
- 파일 확장자 및 명명 규칙

<!-- internationalizer:unit markdown:architecture -->
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
## 다른 도구와의 비교

| 기능 | Internationalizer | i18next | Crowdin | 일반 LLM |
| --- | --- | --- | --- | --- |
| LLM 기반 번역 | 예 | 아니요 | 부분 지원 | 예 |
| 언어별 스타일 가이드 | 예 | 아니요 | 아니요 | 아니요 |
| 용어집 강제 적용 | 예 | 아니요 | 예 | 아니요 |
| 번역 메모리 | 예 | 아니요 | 예 | 아니요 |
| CLI / 로컬 실행 | 예 | 해당 없음 | 아니요 | 수동 |
| Git 친화적 파일 관리 | 예 | 예 | 부분 지원 | 수동 |
| SaaS 종속성 없음 | 예 | 예 | 아니요 | 다양함 |
| 오픈 소스 (AGPL-3.0) | 예 | 예 | 아니요 | 다양함 |

<!-- internationalizer:unit markdown:license -->
## 라이선스

[AGPL-3.0](../../LICENSE)

종속성 관련 고지 사항은 [THIRD_PARTY_NOTICES.md](../../THIRD_PARTY_NOTICES.md)를 참조하세요.

<!-- internationalizer:unit markdown:contributing -->
## 기여

개발 환경 설정 및 가이드라인은 [CONTRIBUTING.md](../../CONTRIBUTING.md)를 참조하세요. 모든 기여에는 DCO(Developer Certificate of Origin) 서명이 필요합니다.
