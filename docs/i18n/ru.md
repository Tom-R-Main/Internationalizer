> [English (original)](../../README.md)

<p align="center">
  <img src="../../assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

AI-нативный пайплайн интернационализации для программных проектов. Переводите, проверяйте и управляйте файлами i18n с помощью LLM.

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](LICENSE)

## Почему Internationalizer?

Большинство инструментов i18n — это либо библиотеки времени выполнения (i18next, react-intl), либо SaaS-платформы для управления ключами (Crowdin, Lokalise). Ни один из них не решает саму проблему перевода должным образом:

- **Ручной перевод** плохо масштабируется при работе с более чем несколькими языками
- **API машинного перевода** (Google Translate, DeepL) игнорируют вашу терминологию, тон и соглашения пользовательского интерфейса
- **Обычный перевод с помощью LLM** работает лучше, но без глоссариев и руководств по стилю вы получаете противоречивые результаты

Internationalizer устроен иначе. Это **CLI-пайплайн**, который объединяет LLM-перевод с:

- **Глоссариями для каждого языка** — обеспечивают единообразную терминологию во всем приложении
- **Руководствами по стилю для каждого языка** — контролируют тон, формальность, плюрализацию и типографику
- **Памятью переводов (Translation memory)** — пропускает неизмененные строки, экономя деньги на вызовах API
- **Валидацией ключей** — выявляет отсутствующие переводы и несовпадения интерполяции до релиза

## Установка

Установка из npm:

```bash
npm install -g internationalizer
```

Или запуск без глобальной установки:

```bash
npx internationalizer --help
```

Пакет npm устанавливает соответствующий предварительно собранный бинарный файл из npm через платформозависимые опциональные зависимости.

Установка с помощью Go:

```bash
go install github.com/Tom-R-Main/Internationalizer/cmd/internationalizer@latest
```

Или сборка из исходного кода:

```bash
git clone https://github.com/Tom-R-Main/Internationalizer.git
cd Internationalizer
go build -o internationalizer ./cmd/internationalizer
```

## Пакеты npm

- Теги Git и версии пакетов npm должны совпадать, например `v0.1.0` и `0.1.0`
- Корневой пакет `internationalizer` зависит от платформенных пакетов, таких как `internationalizer-darwin-arm64`
- Поддерживаемые целевые платформы npm: macOS arm64/x64, Linux arm64/x64, Windows x64
- Для публикации в CI требуется секрет GitHub с именем `NPM_TOKEN`

## Быстрый старт

1. Создайте файл конфигурации в корне вашего проекта:

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

2. Установите ваш API-ключ:

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

3. Предварительный просмотр того, что будет переведено:

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

## Команды

### `translate`

Поиск отсутствующих ключей и их перевод с помощью LLM.

```bash
internationalizer translate                    # перевести все локали
internationalizer translate -l fr              # перевести только на французский
internationalizer translate --dry-run          # предпросмотр без вызовов API
internationalizer translate --batch-size 20    # меньший размер пакета (батча)
internationalizer translate --concurrency 2    # меньше параллельных вызовов
```

### `validate`

Проверка всех файлов локалей на наличие отсутствующих ключей, лишних ключей и несовпадений интерполяции.

```bash
internationalizer validate                     # человекочитаемый вывод
internationalizer validate --json              # машиночитаемый JSON
internationalizer validate -q                  # только код возврата
```

### `detect`

Автоматическое определение фреймворка i18n и предложение конфигурации.

```bash
internationalizer detect
```

Поддерживает: react-i18next, next-intl, vue-i18n, обычный JSON, документацию в Markdown.

### `glossary`

Управление терминами глоссария для каждого языка, которые применяются во время перевода.

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

Управление памятью переводов (кэш ранее переведенных строк в формате JSONL).

```bash
internationalizer tm stats                     # показать количество записей
internationalizer tm export                    # выгрузить в формате JSON
internationalizer tm clear --force             # удалить все записи
```

## Справочник по конфигурации

```yaml
# .internationalizer.yml

# Исходный язык (по умолчанию: en)
source_locale: en

# Языки для перевода (обязательно)
target_locales: [fr, de, es, ja, zh-CN, ar]

# Путь к исходному файлу локали (обязательно)
source_path: locales/en.json

# Настройки провайдера LLM
llm:
  # Провайдер: "anthropic", "openai", "gemini" или "openrouter" (по умолчанию: gemini)
  provider: gemini

  # Имена моделей по умолчанию для провайдеров:
  #   anthropic:  claude-sonnet-4-6
  #   openai:     gpt-5.4
  #   gemini:     gemini-3.1-pro-preview
  #   openrouter: google/gemini-3-flash-preview
  model: gemini-3.1-pro-preview

  # Переменная окружения, содержащая API-ключ
  api_key_env: GOOGLE_AI_STUDIO_API_KEY

  # Базовый URL для OpenAI-совместимых эндпоинтов (опционально)
  # base_url: https://api.openai.com

# Количество ключей на один вызов LLM (по умолчанию: 40)
batch_size: 40

# Параллельные вызовы LLM (по умолчанию: 4)
concurrency: 4

# Директория с Markdown-файлами руководств по стилю для каждой локали (по умолчанию: style-guides)
style_guides_dir: style-guides

# Директория с JSON-файлами глоссариев для каждой локали (по умолчанию: glossary)
glossary_dir: glossary

# Путь к файлу памяти переводов (по умолчанию: .internationalizer/tm.jsonl)
tm_path: .internationalizer/tm.jsonl
```

## Руководства по стилю

Руководства по стилю — это Markdown-файлы, которые внедряются в промпт для LLM-перевода. Они контролируют тон, формальность, типографику и другие языковые соглашения.

```
style-guides/
  _conventions.md    # общие правила для всех языков
  fr.md              # правила для французского языка
  ja.md              # правила для японского языка
  ar.md              # правила для арабского языка
```

### Общие соглашения (`_conventions.md`)

Определяют правила, применимые ко всем языкам: синтаксис интерполяции, сохранение HTML, соглашения о типах строк (кнопки, метки, ошибки) и т. д.

### Руководства для конкретных языков (`{locale}.md`)

Определяют специфичные для языка правила: регистр формальности (ты/вы), пунктуацию (кавычки-елочки, перевернутые вопросительные знаки), формы множественного числа, форматирование дат/чисел и терминологический глоссарий.

Смотрите [`examples/react-app/style-guides/`](examples/react-app/style-guides/) для ознакомления с рабочим примером.

## Формат глоссария

Файлы глоссария представляют собой массивы JSON, хранящиеся в `{glossary_dir}/{locale}.json`:

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

Термины внедряются в промпт LLM в виде терминологической таблицы, обеспечивая единообразный перевод ключевых терминов во всем вашем приложении.

## Память переводов

Память переводов (Translation memory) хранится в виде файла JSONL (одна запись JSON на строку). Каждая запись содержит:

- Исходный ключ и значение
- Переведенное значение
- SHA-256 хэш исходного значения
- Временную метку

При последующих запусках неизмененные строки берутся из кэша TM без вызова LLM, что экономит время и затраты на API. Файл TM удобен для работы с git и может быть закоммичен вместе с файлами локалей.

## Поддерживаемые форматы

| Формат | Расширения | Режим |
|--------|-----------|------|
| JSON | `.json` | Ключ-значение (вложенные, плоские с точечной нотацией) |
| YAML | `.yml`, `.yaml` | Ключ-значение (сохраняет комментарии и порядок) |
| Markdown | `.md`, `.mdx` | Перевод документа целиком |

## Определение типа проекта

`internationalizer detect` определяет вашу конфигурацию i18n путем проверки:

- Зависимостей в `package.json` на наличие react-i18next, next-intl или vue-i18n
- Структуры каталогов, соответствующей общим паттернам локалей
- Расширений файлов и соглашений об именовании

## Архитектура

```
cmd/internationalizer/     Точка входа CLI и определения команд
internal/
  config/                  Загрузка YAML-конфигурации со значениями по умолчанию
  detect/                  Автоопределение типа проекта
  formats/                 Парсеры форматов (JSON, YAML, Markdown)
  glossary/                Управление глоссариями для каждой локали
  llm/                     Интерфейс провайдера LLM + реализации
    anthropic.go           Бэкенд Anthropic Claude
    openai.go              Бэкенд OpenAI / совместимый
    gemini.go              Бэкенд Google Gemini через AI Studio
                           OpenRouter использует openai.go с кастомным base_url
  styleguide/              Загрузчик руководств по стилю
  tm/                      Память переводов JSONL
  translate/               Оркестратор перевода
  validate/                Валидация локалей и сравнение
```

## Сравнение с альтернативами

| Функция | Internationalizer | i18next | Crowdin | Обычная LLM |
|---------|------------------|---------|---------|-------------|
| Перевод на базе LLM | Да | Нет | Частично | Да |
| Руководства по стилю для каждого языка | Да | Нет | Нет | Нет |
| Применение глоссария | Да | Нет | Да | Нет |
| Память переводов | Да | Нет | Да | Нет |
| CLI / локальное выполнение | Да | Н/Д | Нет | Вручную |
| Удобные для Git файлы | Да | Да | Частично | Вручную |
| Отсутствие зависимости от SaaS | Да | Да | Нет | Зависит |
| Открытый исходный код (AGPL-3.0) | Да | Да | Нет | Зависит |

## Лицензия

[AGPL-3.0](LICENSE)

## Участие в разработке

Смотрите [CONTRIBUTING.md](CONTRIBUTING.md) для получения инструкций по настройке среды разработки и руководств. Все вклады требуют подписания DCO (DCO sign-off).

