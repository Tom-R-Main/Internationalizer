> [English (original)](../../README.md)

<p align="center">
  <img src="../../assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

AI-нативный конвейер интернационализации для программных проектов. Переводите, проверяйте и обновляйте файлы i18n с помощью LLM.

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](../../LICENSE)

<p align="center">
<a href="ar.md">العربية</a> · <a href="bn.md">বাংলা</a> · <a href="cs.md">Čeština</a> · <a href="da.md">Dansk</a> · <a href="de.md">Deutsch</a> · <a href="el.md">Ελληνικά</a> · <a href="es.md">Español</a> · <a href="fi.md">Suomi</a> · <a href="fr.md">Français</a> · <a href="he.md">עברית</a> · <a href="hi.md">हिन्दी</a> · <a href="id.md">Indonesia</a> · <a href="it.md">Italiano</a> · <a href="ja.md">日本語</a> · <a href="ko.md">한국어</a> · <a href="ms.md">Bahasa Melayu</a><br><a href="nl.md">Nederlands</a> · <a href="pa.md">ਪੰਜਾਬੀ</a> · <a href="pl.md">Polski</a> · <a href="pt-BR.md">Português</a> · <a href="ro.md">Română</a> · <a href="ru.md">Русский</a> · <a href="sv.md">Svenska</a> · <a href="te.md">తెలుగు</a> · <a href="th.md">ไทย</a> · <a href="tr.md">Türkçe</a> · <a href="uk.md">Українська</a> · <a href="vi.md">Tiếng Việt</a> · <a href="yue.md">粵語</a> · <a href="zh-CN.md">简体中文</a> · <a href="zh-TW.md">繁體中文</a>
</p>

---
<!-- internationalizer:unit markdown:why-internationalizer -->
## Почему Internationalizer?

Большинство инструментов i18n — это либо runtime-библиотеки (i18next, react-intl), либо SaaS-платформы для управления ключами (Crowdin, Lokalise). Ни один из них не решает задачу качественного перевода:

- **Ручной перевод** перестает масштабироваться, когда языков становится больше нескольких
- **API машинного перевода** (Google Translate, DeepL) не учитывают терминологию, тональность и интерфейсные соглашения продукта
- **Обычный перевод через LLM** справляется лучше, но без глоссариев и руководств по стилю выдает несогласованные результаты

Internationalizer устроен иначе. Это **CLI-конвейер**, сочетающий перевод через LLM со следующими возможностями:

- **Глоссарии для каждого языка** — единая и выверенная терминология во всем приложении
- **Руководства по стилю для каждого языка** — контроль тональности, регистра вежливости, форм множественного числа и типографики
- **Память переводов (Translation memory)** — пропуск неизмененных строк и прямая экономия на вызовах API
- **Детерминированная валидация** — выявление отсутствующих или лишних ключей, искажений защищенных структур, нарушений глоссария и ошибок в ICU или формах множественного числа до релиза
<!-- internationalizer:unit markdown:installation -->
## Установка

Установка через npm:

```bash
npm install -g internationalizer
```

Запуск без глобальной установки:

```bash
npx internationalizer --help
```

Пакет npm автоматически устанавливает подходящий предсобранный бинарный файл из реестра через платформозависимые опциональные зависимости.

Установка с помощью Go:

```bash
go install github.com/Tom-R-Main/Internationalizer/cmd/internationalizer@latest
```

Сборка из исходного кода:

```bash
git clone https://github.com/Tom-R-Main/Internationalizer.git
cd Internationalizer
go build -o internationalizer ./cmd/internationalizer
```
<!-- internationalizer:unit markdown:npm-packages -->
## Пакеты npm

- Теги Git и версии пакетов npm должны совпадать, например `v0.1.0` и `0.1.0`
- Корневой пакет `internationalizer` зависит от платформенных пакетов, таких как `internationalizer-darwin-arm64`
- Поддерживаемые платформы npm: macOS arm64/x64, Linux arm64/x64, Windows x64
- Для публикации через CI требуется секрет GitHub с именем `NPM_TOKEN`
<!-- internationalizer:unit markdown:quick-start -->
## Быстрый старт

1. Создайте файл конфигурации в корне проекта:

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

2. Задайте ключ API в переменной окружения:

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

3. Посмотрите, какие строки будут переведены:

```bash
internationalizer translate --dry-run
```

4. Запустите перевод:

```bash
internationalizer translate
```

5. Проверьте все локали:

```bash
internationalizer validate
```
<!-- internationalizer:unit markdown:commands -->
## Команды

### `translate`

Поиск отсутствующих или устаревших ключей и их перевод с помощью LLM.

```bash
internationalizer translate                    # translate all locales
internationalizer translate -l fr              # translate French only
internationalizer translate --dry-run          # preview without API calls
internationalizer translate --adopt-existing   # baseline existing translations without API calls
internationalizer translate --refresh-policy   # refresh prompt/style/model-stale entries
internationalizer translate --batch-size 20    # smaller batches
internationalizer translate --concurrency 2    # fewer parallel calls
```

Состояние перевода раздельно отслеживает пять условий: отсутствие перевода, устаревание исходника, устаревание политики, актуальное состояние и ручное редактирование. Благодаря этому ручная правка не может скрыть изменение исходного текста или правил перевода. Устаревшие по политике значения выводятся в отчете, но повторно переводятся только при указании флага `--refresh-policy`. Отредактированные вручную значения никогда не перезаписываются автоматически. Используйте `--adopt-existing`, чтобы впервые зафиксировать проверенные переводы в манифесте или явно утвердить выверенную ручную правку как новый базовый уровень.

### `validate`

Сверка всех файлов локалей с исходными бандлами. По умолчанию валидация проверяет структурное покрытие (процент присутствующих обязательных целевых ключей), сообщает о лишних ключах в виде предупреждений и завершается с ошибкой при пропущенных ключах, несовпадении переменных интерполяции или неверной структуре ICU MessageFormat.

```bash
internationalizer validate                     # human-readable output
internationalizer validate --json              # machine-readable JSON
internationalizer validate -q                  # exit code only
internationalizer validate --strict             # enforce translation quality rules
internationalizer validate --require-state      # require current manifest provenance
```

Флаг `--strict` дополнительно сообщает процент переведенного содержимого. Любой текстовый фрагмент, идентичный исходнику, считается непереведенным, за исключением случаев, когда в глоссарии явно прописано точное совпадение исходного и целевого значения для всей строки целиком; параметр `ignore_case` учитывается, но вхождение термина глоссария в более длинный текст исключением не признается. Строгий режим завершается ошибкой при наличии лишних ключей, совпадении перевода с исходником, изменении структуры интерполяции, HTML, кода или ссылок Markdown, а также при нарушениях глоссария и отсутствии настроенных форм множественного числа.

Флаг `--require-state` сверяет каждый целевой ключ с манифестом `.internationalizer.lock`. Команда завершается ошибкой, если ключ не отслеживается либо если устарел его записанный исходник, правила перевода или хеш переведенного значения. Этот режим можно сочетать с `--strict`.

В текстовых отчетах и JSON используются стабильные коды находок:

| Код | Значение |
| --- | --- |
| `missing_key` / `extra_key` | Наборы ключей в исходном файле и целевом переводе различаются |
| `blank_translation` | Непустому исходному тексту соответствует пустой перевод в строгом режиме |
| `source_identical` | Текстовое значение в строгом режиме осталось непереведенным |
| `protected_structure_mismatch` | Изменилась структура интерполяции, HTML, кода или ссылок |
| `glossary_violation` | Не найден утвержденный термин глоссария или его вариант |
| `plural_form_missing` | Отсутствует настроенная форма множественного числа для локали |
| `icu_message_syntax` | Исходное или переведенное сообщение ICU сформировано неверно |
| `icu_argument_mismatch` | Различаются имена, типы или стили форматирования аргументов ICU |
| `icu_selector_mismatch` | Различаются селекторы или категория множественного числа недопустима для целевой локали |
| `untracked` | В манифесте отсутствует запись для целевого значения |
| `source_stale` | Исходный текст изменился после зафиксированного перевода |
| `policy_stale` | Изменился сформированный промпт или параметры модели |
| `target_modified` | Переведенное значение отличается от записи в манифесте |

### `detect`

Автоматическое определение используемого фреймворка i18n и подготовка рекомендуемой конфигурации.

```bash
internationalizer detect
```

Поддерживаются: react-i18next, next-intl, vue-i18n, стандартный JSON, документация в Markdown.

### `glossary`

Управление языковыми терминами глоссария, соблюдение которых обеспечивается при переводе.

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

Управление памятью переводов (кэш ранее переведенных строк в формате JSONL).

```bash
internationalizer tm stats                     # show record counts
internationalizer tm export                    # dump as JSON
internationalizer tm clear --force             # delete all records
```
<!-- internationalizer:unit markdown:configuration-reference -->
## Справочник по конфигурации

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

Идентификаторы локалей должны быть корректными тегами BCP 47, такими как `fr`, `pt-BR` или `sr-Latn-RS`. Канонически эквивалентные целевые локали отклоняются как дубликаты, а переопределения провайдера для конкретной локали сопоставляются с учетом канонической эквивалентности. В приведенном выше примере локали без индивидуального переопределения (включая японскую) наследуют общую конфигурацию Gemini.

Значения в формате ICU MessageFormat разбираются структурно. Поддерживаются простые аргументы, `select`, `plural`, `selectordinal`, `number`, `date` и `time`, включая вложенные сообщения, смещения во множественном числе (plural offsets), точные числовые селекторы и символ `#`. Валидация проверяет синтаксис, типы аргументов и стили форматирования, смещения во множественном числе, совпадение веток select и категории множественного числа CLDR для целевой локали. Ответ провайдера, нарушающий эти инварианты, отклоняется до записи в файл локали или в память переводов.

В режиме `i18next-v4` распознанные группы множественного числа из исходника при переводе расширяются до категорий CLDR целевой локали. Если категория существует только в целевом языке, шаблоном для ее перевода служит значение `_other` из исходной группы. Строгая валидация требует наличия всех целевых категорий; категории, присутствующие только в исходнике, необязательны для тех целевых локалей, где они не используются.
<!-- internationalizer:unit markdown:style-guides -->
## Руководства по стилю

Руководства по стилю представляют собой файлы Markdown, которые передаются в системный промпт при переводе через LLM. Они задают тональность, регистр вежливости, типографику и другие языковые нормы.

```
style-guides/
  _conventions.md    # shared rules for all languages
  fr.md              # French-specific rules
  ja.md              # Japanese-specific rules
  ar.md              # Arabic-specific rules
```

### Общие соглашения (`_conventions.md`)

Задают правила, действующие для всех языков: синтаксис интерполяции, сохранение тегов HTML, соглашения по типам строк (кнопки, метки, сообщения об ошибках) и т. д.

### Руководства для отдельных языков (`{locale}.md`)

Задают нормы конкретного языка: регистр вежливости («ты» или «вы»), пунктуацию («кавычки-елочки», перевернутые вопросительные знаки), формы множественного числа, форматы дат и чисел, а также терминологический глоссарий.

Руководства по стилю — это долговременные входные правила, а не генерируемый результат. Internationalizer считывает их, но никогда не перезаписывает. Их содержимое хешируется отдельно от глоссария и контракта промпта, поэтому изменения в коде приложения не приводят к устареванию переводов. Редактирование руководства намеренно помечает соответствующую локаль для пересмотра правил; изменение внутренних формулировок промпта к этому не приводит, если только не изменилась версия самого контракта промпта.

Рабочий пример структуры приведен в каталоге [`examples/react-app/style-guides/`](../../examples/react-app/style-guides/).
<!-- internationalizer:unit markdown:glossary-format -->
## Формат глоссария

Файлы глоссария представляют собой массивы JSON, хранящиеся по пути `{glossary_dir}/{locale}.json`:

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

В массиве `variants` перечисляются другие допустимые формы перевода. Поле `enforcement` может принимать значения `error`, `warning` либо опускаться (по умолчанию действует режим ошибки). Термины подставляются в промпт LLM в виде таблицы соответствий, гарантируя единообразный перевод во всем приложении. Точная запись вида `{"source":"API","target":"API"}` также освобождает полностью совпадающее исходное значение от ошибок непереведенного текста в строгом режиме; при этом частичное вхождение `API` в более длинную строку исключением не является.
<!-- internationalizer:unit markdown:translation-memory -->
## Память переводов

Память переводов хранится в файле формата JSONL (одна запись JSON на строку). Каждая запись содержит:

- Бандл, ключ, исходный текст, переведенный текст и каноническую целевую локаль
- Хеши исходного текста, руководства по стилю, глоссария, контракта промпта и общий хеш политики
- Провайдера и модель, выполнившие перевод
- Временную метку

При повторных запусках строки с неизменившимися хешами исходника и правил берутся из кэша без обращения к LLM. По умолчанию файл хранится в игнорируемом каталоге `.internationalizer/` в качестве локального кэша. Если проект предполагает совместное использование памяти переводов в команде, укажите в `tm_path` путь к отслеживаемому файлу. Манифест `.internationalizer.lock`, предназначенный для код-ревью, версионируется независимо.
<!-- internationalizer:unit markdown:supported-formats -->
## Поддерживаемые форматы

| Формат | Расширения | Режим |
| --- | --- | --- |
| JSON | `.json` | Ключ-значение (вложенные структуры, плоские ключи с точечной нотацией) |
| YAML | `.yml`, `.yaml` | Ключ-значение (сохранение комментариев и порядка следования) |
| Markdown | `.md`, `.mdx` | Преамбула и разделы уровня H2 |

В целевых файлах Markdown перед разделами H2 расставляются невидимые комментарии `internationalizer:unit`. Эти стабильные метки позволяют добавлять, перемещать или редактировать отдельные разделы оригинала без повторного перевода остального документа. В существующие файлы без разметки маркеры добавляются при очередном успешном обновлении.
<!-- internationalizer:unit markdown:project-type-detection -->
## Определение типа проекта

Команда `internationalizer detect` определяет структуру i18n в проекте по следующим признакам:

- Наличие зависимостей react-i18next, next-intl или vue-i18n в файле `package.json`
- Структура каталогов, соответствующая типовым шаблонам локализации
- Расширения файлов и соглашения об именовании
<!-- internationalizer:unit markdown:architecture -->
## Архитектура

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
## Сравнение с альтернативами

| Возможность | Internationalizer | i18next | Crowdin | Обычные LLM |
| --- | --- | --- | --- | --- |
| Перевод на базе LLM | Да | Нет | Частично | Да |
| Руководства по стилю для каждого языка | Да | Нет | Нет | Нет |
| Контроль соблюдения глоссария | Да | Нет | Да | Нет |
| Память переводов | Да | Нет | Да | Нет |
| CLI и локальный запуск | Да | Неприменимо | Нет | Вручную |
| Удобство хранения файлов в Git | Да | Да | Частично | Вручную |
| Без привязки к SaaS | Да | Да | Нет | Зависит от случая |
| Открытый исходный код (AGPL-3.0) | Да | Да | Нет | Зависит от случая |
<!-- internationalizer:unit markdown:license -->
## Лицензия

[AGPL-3.0](../../LICENSE)

Сведения о зависимостях см. в [THIRD_PARTY_NOTICES.md](../../THIRD_PARTY_NOTICES.md).
<!-- internationalizer:unit markdown:contributing -->
## Участие в разработке

Инструкции по настройке среды разработки и правила участия см. в [CONTRIBUTING.md](../../CONTRIBUTING.md). Для всех вкладов обязательно подписание DCO (DCO sign-off).
