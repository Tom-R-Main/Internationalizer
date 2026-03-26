<p align="center">
  <img src="assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

AI-нативний конвеєр інтернаціоналізації для програмних проєктів. Перекладайте, перевіряйте та керуйте файлами i18n за допомогою LLM.

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](LICENSE)

## Чому Internationalizer?

Більшість інструментів i18n — це або бібліотеки часу виконання (i18next, react-intl), або SaaS-платформи для керування ключами (Crowdin, Lokalise). Жоден із них не вирішує саму проблему перекладу належним чином:

- **Ручний переклад** не масштабується більше ніж на кілька мов
- **API машинного перекладу** (Google Translate, DeepL) ігнорують вашу термінологію, тон та угоди щодо інтерфейсу користувача
- **Звичайний переклад за допомогою LLM** працює краще, але без глосаріїв та посібників зі стилю ви отримуєте непослідовні результати

Internationalizer відрізняється. Це **CLI-конвеєр**, який поєднує переклад за допомогою LLM із:

- **Глосаріями для кожної мови** — забезпечують узгоджену термінологію у всьому вашому застосунку
- **Посібниками зі стилю для кожної мови** — контролюють тон, формальність, множину та типографіку
- **Пам'яттю перекладів (Translation memory)** — пропускає незмінені рядки, заощаджуючи гроші на викликах API
- **Перевіркою ключів** — виявляє відсутні переклади та невідповідності інтерполяції до релізу

## Встановлення

Встановіть через npm:

```bash
npm install -g internationalizer
```

Або запустіть без глобального встановлення:

```bash
npx internationalizer --help
```

Пакет npm встановлює відповідний попередньо зібраний бінарний файл з npm через специфічні для платформи необов'язкові залежності.

Встановіть за допомогою Go:

```bash
go install github.com/Tom-R-Main/Internationalizer/cmd/internationalizer@latest
```

Або зберіть із вихідного коду:

```bash
git clone https://github.com/Tom-R-Main/Internationalizer.git
cd Internationalizer
go build -o internationalizer ./cmd/internationalizer
```

## Пакети npm

- Теги Git та версії пакетів npm мають збігатися, наприклад `v0.1.0` та `0.1.0`
- Кореневий пакет `internationalizer` залежить від пакетів платформи, таких як `internationalizer-darwin-arm64`
- Підтримувані цільові платформи npm: macOS arm64/x64, Linux arm64/x64, Windows x64
- Для публікації через CI потрібен секрет GitHub із назвою `NPM_TOKEN`

## Швидкий старт

1. Створіть файл конфігурації в корені вашого проєкту:

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

2. Встановіть ваш API-ключ:

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

3. Попередній перегляд того, що буде перекладено:

```bash
internationalizer translate --dry-run
```

4. Запустіть переклад:

```bash
internationalizer translate
```

5. Перевірте всі локалі:

```bash
internationalizer validate
```

## Команди

### `translate`

Знаходить відсутні ключі та перекладає їх за допомогою LLM.

```bash
internationalizer translate                    # перекласти всі локалі
internationalizer translate -l fr              # перекласти лише французьку
internationalizer translate --dry-run          # попередній перегляд без викликів API
internationalizer translate --batch-size 20    # менші пакети
internationalizer translate --concurrency 2    # менше паралельних викликів
```

### `validate`

Перевіряє всі файли локалей на наявність відсутніх ключів, зайвих ключів та невідповідностей інтерполяції.

```bash
internationalizer validate                     # зручний для читання вивід
internationalizer validate --json              # машинозчитуваний JSON
internationalizer validate -q                  # лише код завершення
```

### `detect`

Автоматично визначає фреймворк i18n та пропонує конфігурацію.

```bash
internationalizer detect
```

Підтримує: react-i18next, next-intl, vue-i18n, звичайний JSON, документацію у Markdown.

### `glossary`

Керує термінами глосарія для кожної мови, які застосовуються під час перекладу.

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

Керує пам'яттю перекладів (кеш JSONL попередньо перекладених рядків).

```bash
internationalizer tm stats                     # показати кількість записів
internationalizer tm export                    # експортувати як JSON
internationalizer tm clear --force             # видалити всі записи
```

## Довідник з конфігурації

```yaml
# .internationalizer.yml

# Вихідна мова (за замовчуванням: en)
source_locale: en

# Мови для перекладу (обов'язково)
target_locales: [fr, de, es, ja, zh-CN, ar]

# Шлях до файлу вихідної локалі (обов'язково)
source_path: locales/en.json

# Налаштування провайдера LLM
llm:
  # Провайдер: "anthropic", "openai", "gemini" або "openrouter" (за замовчуванням: gemini)
  provider: gemini

  # Назви моделей за замовчуванням для провайдерів:
  #   anthropic:  claude-sonnet-4-6
  #   openai:     gpt-5.4
  #   gemini:     gemini-3.1-pro-preview
  #   openrouter: google/gemini-3-flash-preview
  model: gemini-3.1-pro-preview

  # Змінна середовища, що містить API-ключ
  api_key_env: GOOGLE_AI_STUDIO_API_KEY

  # Базова URL-адреса для сумісних з OpenAI кінцевих точок (необов'язково)
  # base_url: https://api.openai.com

# Кількість ключів на один виклик LLM (за замовчуванням: 40)
batch_size: 40

# Паралельні виклики LLM (за замовчуванням: 4)
concurrency: 4

# Каталог, що містить файли Markdown посібників зі стилю для кожної локалі (за замовчуванням: style-guides)
style_guides_dir: style-guides

# Каталог, що містить файли JSON глосаріїв для кожної локалі (за замовчуванням: glossary)
glossary_dir: glossary

# Шлях до файлу пам'яті перекладів (за замовчуванням: .internationalizer/tm.jsonl)
tm_path: .internationalizer/tm.jsonl
```

## Посібники зі стилю

Посібники зі стилю — це файли Markdown, які додаються до промпту перекладу LLM. Вони контролюють тон, формальність, типографіку та інші специфічні для мови угоди.

```
style-guides/
  _conventions.md    # спільні правила для всіх мов
  fr.md              # специфічні правила для французької
  ja.md              # специфічні правила для японської
  ar.md              # специфічні правила для арабської
```

### Спільні угоди (`_conventions.md`)

Визначають правила, які застосовуються до всіх мов: синтаксис інтерполяції, збереження HTML, угоди щодо типів рядків (кнопки, мітки, помилки тощо).

### Посібники для кожної мови (`{locale}.md`)

Визначають специфічні для мови правила: регістр формальності (tu чи vous), пунктуація (лапки-ялинки, перевернуті знаки питання), форми множини, форматування дат/чисел та глосарій термінології.

Дивіться [`examples/react-app/style-guides/`](examples/react-app/style-guides/) для робочого прикладу.

## Формат глосарія

Файли глосарія — це масиви JSON, що зберігаються у `{glossary_dir}/{locale}.json`:

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

Терміни додаються до промпту LLM у вигляді таблиці термінології, що забезпечує узгоджений переклад ключових термінів у всьому вашому застосунку.

## Пам'ять перекладів

Пам'ять перекладів зберігається як файл JSONL (один запис JSON на рядок). Кожен запис містить:

- Вихідний ключ та значення
- Перекладене значення
- Хеш SHA-256 вихідного значення
- Часову мітку

Під час наступних запусків незмінені рядки беруться з кешу TM без виклику LLM, що економить час і витрати на API. Файл TM є дружнім до git і може бути закомічений разом із вашими файлами локалей.

## Підтримувані формати

| Формат | Розширення | Режим |
|--------|-----------|------|
| JSON | `.json` | Ключ-значення (вкладені, зведені за допомогою крапкової нотації) |
| YAML | `.yml`, `.yaml` | Ключ-значення (зберігає коментарі та порядок) |
| Markdown | `.md`, `.mdx` | Переклад усього документа |

## Визначення типу проєкту

`internationalizer detect` ідентифікує ваші налаштування i18n шляхом перевірки:

- Залежностей `package.json` для react-i18next, next-intl або vue-i18n
- Структур каталогів, що відповідають загальним шаблонам локалей
- Розширень файлів та угод щодо найменування

## Архітектура

```
cmd/internationalizer/     Точка входу CLI та визначення команд
internal/
  config/                  Завантаження конфігурації YAML зі значеннями за замовчуванням
  detect/                  Автоматичне визначення типу проєкту
  formats/                 Парсери форматів (JSON, YAML, Markdown)
  glossary/                Керування глосарієм для кожної локалі
  llm/                     Інтерфейс провайдера LLM + реалізації
    anthropic.go           Бекенд Anthropic Claude
    openai.go              Бекенд OpenAI / сумісний
    gemini.go              Бекенд Google Gemini через AI Studio
                           OpenRouter використовує openai.go з власним base_url
  styleguide/              Завантажувач посібника зі стилю
  tm/                      Пам'ять перекладів JSONL
  translate/               Оркестратор перекладу
  validate/                Перевірка локалей та порівняння
```

## Порівняння з альтернативами

| Функція | Internationalizer | i18next | Crowdin | Звичайний LLM |
|---------|------------------|---------|---------|-------------|
| Переклад за допомогою LLM | Так | Ні | Частково | Так |
| Посібники зі стилю для кожної мови | Так | Ні | Ні | Ні |
| Застосування глосарія | Так | Ні | Так | Ні |
| Пам'ять перекладів | Так | Ні | Так | Ні |
| CLI / локальне виконання | Так | Н/Д | Ні | Вручну |
| Дружні до Git файли | Так | Так | Частково | Вручну |
| Відсутність залежності від SaaS | Так | Так | Ні | Варіюється |
| Відкритий вихідний код (AGPL-3.0) | Так | Так | Ні | Варіюється |

## Ліцензія

[AGPL-3.0](LICENSE)

## Внесок

Дивіться [CONTRIBUTING.md](CONTRIBUTING.md) для отримання інструкцій з налаштування середовища розробки та правил. Усі внески вимагають підписання DCO.

