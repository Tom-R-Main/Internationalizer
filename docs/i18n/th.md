> [English (original)](../../README.md)

<p align="center">
  <img src="../../assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

ไปป์ไลน์ internationalization แบบ AI-native สำหรับโปรเจกต์ซอฟต์แวร์ แปล ตรวจสอบความถูกต้อง และจัดการไฟล์ i18n โดยใช้ LLM

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](../../LICENSE)

<p align="center">
<a href="ar.md">العربية</a> · <a href="bn.md">বাংলা</a> · <a href="cs.md">Čeština</a> · <a href="da.md">Dansk</a> · <a href="de.md">Deutsch</a> · <a href="el.md">Ελληνικά</a> · <a href="es.md">Español</a> · <a href="fi.md">Suomi</a> · <a href="fr.md">Français</a> · <a href="he.md">עברית</a> · <a href="hi.md">हिन्दी</a> · <a href="id.md">Indonesia</a> · <a href="it.md">Italiano</a> · <a href="ja.md">日本語</a> · <a href="ko.md">한국어</a> · <a href="ms.md">Bahasa Melayu</a><br><a href="nl.md">Nederlands</a> · <a href="pa.md">ਪੰਜਾਬੀ</a> · <a href="pl.md">Polski</a> · <a href="pt-BR.md">Português</a> · <a href="ro.md">Română</a> · <a href="ru.md">Русский</a> · <a href="sv.md">Svenska</a> · <a href="te.md">తెలుగు</a> · <a href="th.md">ไทย</a> · <a href="tr.md">Türkçe</a> · <a href="uk.md">Українська</a> · <a href="vi.md">Tiếng Việt</a> · <a href="yue.md">粵語</a> · <a href="zh-CN.md">简体中文</a> · <a href="zh-TW.md">繁體中文</a>
</p>

---
<!-- internationalizer:unit markdown:why-internationalizer -->
## ทำไมต้อง Internationalizer?

เครื่องมือ i18n ส่วนใหญ่เป็นไลบรารีระดับรันไทม์ (i18next, react-intl) หรือไม่ก็แพลตฟอร์ม SaaS สำหรับจัดการคีย์ (Crowdin, Lokalise) ซึ่งไม่มีเครื่องมือใดที่แก้ปัญหาด้านการแปลได้อย่างแท้จริง:

- **การแปลด้วยตนเอง** ไม่สามารถขยายขนาดได้เมื่อมีมากกว่าสองสามภาษา
- **API การแปลภาษาด้วยเครื่อง** (Google Translate, DeepL) ละเลยคำศัพท์เฉพาะ โทนเสียง และแบบแผน UI ของคุณ
- **การแปลด้วย LLM ทั่วไป** ให้ผลลัพธ์ที่ดีกว่า แต่หากปราศจากอภิธานศัพท์และคู่มือสไตล์ ผลลัพธ์ที่ได้จะขาดความสม่ำเสมอ

Internationalizer แตกต่างออกไป เพราะนี่คือ **CLI pipeline** ที่ผสานการแปลด้วย LLM เข้ากับ:

- **อภิธานศัพท์เฉพาะแต่ละภาษา** — บังคับใช้คำศัพท์ให้สอดคล้องกันทั่วทั้งแอปของคุณ
- **คู่มือสไตล์เฉพาะแต่ละภาษา** — ควบคุมโทนเสียง ระดับความเป็นทางการ การผันพหูพจน์ และการจัดรูปแบบตัวพิมพ์
- **หน่วยความจำการแปล** — ข้ามสตริงที่ไม่มีการเปลี่ยนแปลง ช่วยประหยัดค่าใช้จ่ายในการเรียกใช้ API
- **การตรวจสอบความถูกต้องแบบกำหนดแน่นอน** — ดักจับคีย์ที่หายไปหรือเกินมา โครงสร้างที่มีการป้องกันคลาดเคลื่อน ปัญหาอภิธานศัพท์ ตลอดจนข้อผิดพลาดของพหูพจน์หรือ ICU ก่อนที่จะนำไปใช้งานจริง
<!-- internationalizer:unit markdown:installation -->
## การติดตั้ง

ติดตั้งจาก npm:

```bash
npm install -g internationalizer
```

หรือรันโดยไม่ต้องติดตั้งแบบส่วนกลาง:

```bash
npx internationalizer --help
```

แพ็กเกจ npm จะติดตั้งไบนารีที่สร้างไว้ล่วงหน้าซึ่งตรงกับระบบจาก npm ผ่าน optional dependencies เฉพาะของแต่ละแพลตฟอร์ม

ติดตั้งด้วย Go:

```bash
go install github.com/Tom-R-Main/Internationalizer/cmd/internationalizer@latest
```

หรือบิลด์จากซอร์สโค้ด:

```bash
git clone https://github.com/Tom-R-Main/Internationalizer.git
cd Internationalizer
go build -o internationalizer ./cmd/internationalizer
```
<!-- internationalizer:unit markdown:npm-packages -->
## แพ็กเกจ npm

- แท็ก Git และเวอร์ชันของแพ็กเกจ npm ต้องตรงกัน เช่น `v0.1.0` และ `0.1.0`
- แพ็กเกจราก `internationalizer` ขึ้นอยู่กับแพ็กเกจประจำแต่ละแพลตฟอร์ม เช่น `internationalizer-darwin-arm64`
- แพลตฟอร์มเป้าหมาย npm ที่รองรับ: macOS arm64/x64, Linux arm64/x64, Windows x64
- การเผยแพร่ผ่าน CI ต้องใช้ GitHub secret ที่ชื่อว่า `NPM_TOKEN`
<!-- internationalizer:unit markdown:quick-start -->
## เริ่มต้นใช้งานด่วน

1. สร้างไฟล์คอนฟิกที่ไดเรกทอรีรากของโปรเจกต์:

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

2. กำหนดค่า API key ของคุณ:

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

3. ดูตัวอย่างรายการที่จะถูกแปล:

```bash
internationalizer translate --dry-run
```

4. ดำเนินการแปล:

```bash
internationalizer translate
```

5. ตรวจสอบความถูกต้องของทุกภาษา:

```bash
internationalizer validate
```
<!-- internationalizer:unit markdown:commands -->
## คำสั่ง

### `translate`

ค้นหาคีย์ที่หายไปหรือล้าสมัย แล้วแปลผ่าน LLM

```bash
internationalizer translate                    # translate all locales
internationalizer translate -l fr              # translate French only
internationalizer translate --dry-run          # preview without API calls
internationalizer translate --adopt-existing   # baseline existing translations without API calls
internationalizer translate --refresh-policy   # refresh prompt/style/model-stale entries
internationalizer translate --batch-size 20    # smaller batches
internationalizer translate --concurrency 2    # fewer parallel calls
```

สถานะการแปลจะรายงานเงื่อนไข missing, source-stale, policy-stale, current และการแก้ไขด้วยตนเองอย่างเป็นอิสระต่อกัน ดังนั้นการแก้ไขด้วยตนเองจึงไม่สามารถบดบังการเปลี่ยนแปลงของต้นทางหรือนโยบายได้ ค่าที่เข้าข่าย policy-stale จะได้รับการรายงานแต่จะถูกแปลใหม่เมื่อระบุแฟล็ก `--refresh-policy` เท่านั้น ส่วนค่าที่มีการแก้ไขด้วยตนเองจะไม่ถูกเขียนทับโดยอัตโนมัติ ให้ใช้ `--adopt-existing` เมื่อนำ manifest เข้าสู่ชุดคำแปลที่ผ่านการตรวจสอบแล้ว หรือเมื่อต้องการยอมรับการแก้ไขด้วยตนเองที่ตรวจสอบแล้วให้เป็นเกณฑ์มาตรฐานใหม่อย่างชัดเจน

### `validate`

ตรวจสอบไฟล์ภาษาทั้งหมดเทียบกับชุดบันเดิลต้นทาง การตรวจสอบความถูกต้องตามค่าเริ่มต้นจะตรวจเช็กความครอบคลุมเชิงโครงสร้าง (สัดส่วนร้อยละของคีย์ปลายทางที่จำเป็นทั้งหมดที่มีอยู่) รายงานคีย์ส่วนเกินเป็นคำเตือน และจะล้มเหลวทันทีหากพบคีย์ที่หายไป ตัวแปร interpolation ไม่ตรงกัน หรือโครงสร้าง ICU MessageFormat ไม่ถูกต้อง

```bash
internationalizer validate                     # human-readable output
internationalizer validate --json              # machine-readable JSON
internationalizer validate -q                  # exit code only
internationalizer validate --strict             # enforce translation quality rules
internationalizer validate --require-state      # require current manifest provenance
```

`--strict` จะรายงานความครอบคลุมของการแปลด้วย ค่าเชิงภาษาที่เหมือนกับต้นทางจะถือว่ายังไม่ได้แปล เว้นแต่อภิธานศัพท์จะมีการระบุรายการสำหรับค่านั้นอย่างชัดเจนที่ตรงกันทั้งค่าต้นทางและปลายทางที่เหมือนกันอย่างสมบูรณ์ ทั้งนี้จะยังคงเคารพการตั้งค่า `ignore_case` แต่คำศัพท์ในอภิธานศัพท์ที่ฝังอยู่ภายในค่าที่ยาวกว่าจะไม่ได้รับการยกเว้น โหมดเข้มงวดจะล้มเหลวเมื่อพบคีย์ส่วนเกิน ค่าที่เหมือนกับต้นทาง โครงสร้าง interpolation/HTML/โค้ด/ลิงก์ Markdown ที่เปลี่ยนแปลง การละเมิดอภิธานศัพท์ ตลอดจนรูปพหูพจน์ที่มีการกำหนดค่าไว้

`--require-state` จะตรวจสอบความถูกต้องของแต่ละปลายทางเทียบกับ `.internationalizer.lock` โดยจะล้มเหลวเมื่อคีย์ไม่ได้ถูกติดตาม หรือเมื่อแฮชของต้นทาง นโยบายการแปล หรือปลายทางที่บันทึกไว้ล้าสมัย แฟล็กนี้สามารถใช้ร่วมกับ `--strict` ได้

รายงานทั้งแบบสำหรับมนุษย์และแบบ JSON จะใช้รหัสผลการตรวจพบที่คงที่:

| รหัส | ความหมาย |
| --- | --- |
| `missing_key` / `extra_key` | ชุดคีย์ต้นทางและปลายทางแตกต่างกัน |
| `blank_translation` | ต้นทางมีเนื้อหา แต่ปลายทางในโหมดเข้มงวดว่างเปล่า |
| `source_identical` | ค่าเชิงภาษาในโหมดเข้มงวดยังคงไม่ได้รับการแปล |
| `protected_structure_mismatch` | โครงสร้าง interpolation, HTML, โค้ด หรือลิงก์เกิดการเปลี่ยนแปลง |
| `glossary_violation` | ไม่พบคำหรือรูปแบบย่อยปลายทางที่ได้รับอนุมัติ |
| `plural_form_missing` | ขาดรูปพหูพจน์ตามที่กำหนดค่าไว้สำหรับภาษานั้น |
| `icu_message_syntax` | โครงสร้างข้อความ ICU ของต้นทางหรือปลายทางมีรูปแบบไม่ถูกต้อง |
| `icu_argument_mismatch` | ชื่อ ชนิด หรือสไตล์ตัวจัดรูปแบบของอาร์กิวเมนต์ ICU แตกต่างกัน |
| `icu_selector_mismatch` | ตัวเลือกแตกต่างกัน หรือหมวดหมู่พหูพจน์ไม่ถูกต้องสำหรับภาษาปลายทาง |
| `untracked` | ไม่มีบันทึก manifest สำหรับปลายทางนี้ |
| `source_stale` | เนื้อหาต้นทางเปลี่ยนแปลงไปหลังจากการบันทึกการแปล |
| `policy_stale` | พรอมต์ที่สร้างขึ้นหรือการตั้งค่าโมเดลมีการเปลี่ยนแปลง |
| `target_modified` | เนื้อหาปลายทางแตกต่างจากบันทึกใน manifest |

### `detect`

ตรวจหาเฟรมเวิร์ก i18n อัตโนมัติและแนะนำการกำหนดค่าที่เหมาะสม

```bash
internationalizer detect
```

รองรับ: react-i18next, next-intl, vue-i18n, vanilla JSON, markdown docs

### `glossary`

จัดการคำศัพท์ในอภิธานศัพท์เฉพาะแต่ละภาษาที่ถูกบังคับใช้ระหว่างการแปล

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

จัดการหน่วยความจำการแปล (แคช JSONL สำหรับสตริงที่เคยแปลแล้ว)

```bash
internationalizer tm stats                     # show record counts
internationalizer tm export                    # dump as JSON
internationalizer tm clear --force             # delete all records
```
<!-- internationalizer:unit markdown:configuration-reference -->
## ข้อมูลอ้างอิงการกำหนดค่า

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

ตัวระบุภาษาต้องเป็นแท็ก BCP 47 ที่มีรูปแบบถูกต้อง เช่น `fr`, `pt-BR` หรือ `sr-Latn-RS` ค่าภาษาเป้าหมายที่เทียบเท่าเชิงมาตรฐาน (canonical-equivalent) จะถูกปฏิเสธเป็นรายการซ้ำซ้อน และการเขียนทับผู้ให้บริการเฉพาะภาษาจะจับคู่ตามการสะกดที่เทียบเท่าเชิงมาตรฐาน ในตัวอย่างข้างต้น ภาษาที่ไม่มีการกำหนดค่าเขียนทับ—ซึ่งรวมถึงภาษาญี่ปุ่น—จะสืบทอดการกำหนดค่าสากลของ Gemini

ค่าของ ICU MessageFormat จะได้รับการแจงส่วนเชิงโครงสร้าง รองรับอาร์กิวเมนต์แบบพื้นฐาน, `select`, `plural`, `selectordinal`, `number`, `date` และ `time` ตลอดจนข้อความซ้อนในตัว, ออฟเซ็ตของพหูพจน์, ตัวเลือกตัวเลขที่แน่นอน และ `#` การตรวจสอบความถูกต้องจะตรวจสอบไวยากรณ์ ชนิดอาร์กิวเมนต์ สไตล์ตัวจัดรูปแบบ ออฟเซ็ตพหูพจน์ ความตรงกันของสาขา select และหมวดหมู่พหูพจน์ CLDR ของภาษาเป้าหมาย หากผลลัพธ์จากผู้ให้บริการละเมิดข้อกำหนดเหล่านี้ ระบบจะปฏิเสธก่อนที่จะเขียนไฟล์ภาษาหรือเรกคอร์ดลงหน่วยความจำการแปล

สำหรับ `i18next-v4` กลุ่มพหูพจน์ต้นทางที่ระบบรู้จักจะถูกขยายระหว่างการแปลให้สอดคล้องกับหมวดหมู่ CLDR ของภาษาเป้าหมาย หมวดหมู่ที่มีเฉพาะในภาษาเป้าหมายจะใช้ค่า `_other` ของกลุ่มต้นทางเป็นเทมเพลตสำหรับการแปล ทั้งนี้การตรวจสอบความถูกต้องแบบเข้มงวดจะกำหนดให้ต้องมีหมวดหมู่ปลายทางเหล่านั้นครบถ้วน ส่วนหมวดหมู่ที่มีเฉพาะในภาษาต้นทางจะถือเป็นทางเลือกสำหรับภาษาเป้าหมายที่ไม่ได้ใช้งาน
<!-- internationalizer:unit markdown:style-guides -->
## คู่มือสไตล์

คู่มือสไตล์คือไฟล์ Markdown ที่จะถูกแทรกลงในพรอมต์การแปลของ LLM เพื่อควบคุมโทนเสียง ระดับความเป็นทางการ การจัดรูปแบบตัวพิมพ์ และแบบแผนเฉพาะภาษาอื่นๆ

```
style-guides/
  _conventions.md    # shared rules for all languages
  fr.md              # French-specific rules
  ja.md              # Japanese-specific rules
  ar.md              # Arabic-specific rules
```

### ข้อตกลงร่วมกัน (`_conventions.md`)

กำหนดกฎที่ใช้ร่วมกันกับทุกภาษา: ไวยากรณ์ interpolation, การคงโครงสร้าง HTML, แบบแผนประเภทของสตริง (ปุ่ม เทียบกับ ป้ายกำกับ เทียบกับ ข้อผิดพลาด) เป็นต้น

### คู่มือเฉพาะแต่ละภาษา (`{locale}.md`)

กำหนดกฎเฉพาะของแต่ละภาษา: ระดับความเป็นทางการ (tu เทียบกับ vous), เครื่องหมายวรรคตอน (guillemets, เครื่องหมายคำถามกลับหัว), รูปพหูพจน์, การจัดรูปแบบวันที่/ตัวเลข และอภิธานศัพท์คำศัพท์

คู่มือสไตล์เป็นอินพุตเชิงนโยบายที่คงทน ไม่ใช่ผลลัพธ์ที่สร้างขึ้น Internationalizer จะอ่านคู่มือเหล่านี้แต่จะไม่เขียนทับอย่างเด็ดขาด เนื้อหาจะถูกสร้างแฮชแยกต่างหากจากอภิธานศัพท์และข้อกำหนดของพรอมต์ ดังนั้นการเปลี่ยนแปลงโค้ดของแอปพลิเคชันจะไม่ทำให้การแปลกลายเป็นสถานะล้าสมัย การแก้ไขคู่มือถือเป็นความตั้งใจทำเครื่องหมายภาษานั้นเพื่อตรวจสอบนโยบาย แต่การเปลี่ยนถ้อยคำภายในพรอมต์จะไม่ส่งผล เว้นแต่เวอร์ชันข้อกำหนดของพรอมต์จะมีการเปลี่ยนแปลงด้วยเช่นกัน

ดูตัวอย่างการใช้งานจริงได้ที่ [`examples/react-app/style-guides/`](../../examples/react-app/style-guides/)
<!-- internationalizer:unit markdown:glossary-format -->
## รูปแบบอภิธานศัพท์

ไฟล์อภิธานศัพท์เป็นอาร์เรย์ JSON ที่จัดเก็บไว้ใน `{glossary_dir}/{locale}.json`:

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

`variants` แสดงรายการรูปแบบปลายทางอื่นๆ ที่ได้รับอนุมัติ ส่วน `enforcement` อาจระบุเป็น `error`, `warning` หรือละเว้นไว้เพื่อใช้พฤติกรรมข้อผิดพลาดตามค่าเริ่มต้น คำศัพท์จะถูกส่งเข้าไปในพรอมต์ของ LLM ในรูปแบบตารางคำศัพท์ เพื่อให้มั่นใจได้ว่าการแปลทั่วทั้งแอปพลิเคชันจะสอดคล้องกัน รายการที่ระบุค่าตรงกันทุกประการ เช่น `{"source":"API","target":"API"}` จะช่วยยกเว้นค่านั้นทั้งค่าที่เหมือนกับต้นทางจากการตรวจพบว่าเป็นค่าที่ยังไม่ได้แปลในโหมดเข้มงวดด้วย แต่จะไม่ยกเว้นข้อความที่ยาวกว่าซึ่งเพียงแค่มีคำว่า `API` อยู่ภายใน
<!-- internationalizer:unit markdown:translation-memory -->
## หน่วยความจำการแปล

หน่วยความจำการแปลจะจัดเก็บในรูปแบบไฟล์ JSONL (หนึ่งเรกคอร์ด JSON ต่อหนึ่งบรรทัด) แต่ละเรกคอร์ดประกอบด้วย:

- บันเดิล คีย์ ค่าต้นทาง ค่าที่แปลแล้ว และภาษาปลายทางตามรูปแบบมาตรฐาน
- แฮชของต้นทาง คู่มือสไตล์ อภิธานศัพท์ ข้อกำหนดของพรอมต์ และนโยบายโดยรวม
- ผู้ให้บริการและโมเดลที่ทำการแปล
- การประทับเวลา

ในการรันครั้งต่อๆ ไป ข้อความที่มีแฮชของต้นทางและนโยบายตรงกันจะถูกดึงมาจากแคชทันทีโดยไม่ต้องเรียกใช้ LLM พาธเริ่มต้นจะอยู่ภายใต้ไดเรกทอรี `.internationalizer/` ซึ่งถูกละเว้นไว้ จึงทำงานเป็นแคชเฉพาะเครื่อง หากโปรเจกต์ของคุณต้องการแชร์หน่วยความจำการแปลร่วมกัน ให้กำหนด `tm_path` ไปยังตำแหน่งที่มีการติดตามในระบบควบคุมเวอร์ชัน ส่วน manifest `.internationalizer.lock` ที่สามารถตรวจสอบได้จะมีการกำหนดเวอร์ชันแยกต่างหาก
<!-- internationalizer:unit markdown:supported-formats -->
## รูปแบบที่รองรับ

| รูปแบบ | นามสกุลไฟล์ | โหมด |
|--------|-----------|------|
| JSON | `.json` | Key-value (ซ้อนกันได้, แบนราบด้วย dot-notation) |
| YAML | `.yml`, `.yaml` | Key-value (คงคอมเมนต์และลำดับไว้) |
| Markdown | `.md`, `.mdx` | ส่วนนำและส่วนย่อยระดับ H2 |

ไฟล์ปลายทาง Markdown จะมีคอมเมนต์ `internationalizer:unit` ที่มองไม่เห็นกำกับไว้ก่อนส่วนย่อยระดับ H2 เครื่องหมายที่เสถียรเหล่านี้ช่วยให้ Internationalizer สามารถเพิ่ม ย้าย หรือแก้ไขเนื้อหาต้นทางส่วนใดส่วนหนึ่งได้โดยไม่ต้องแปลส่วนย่อยอื่นที่ไม่เกี่ยวข้องซ้ำ เอกสารเดิมที่ยังไม่มีเครื่องหมายกำกับจะได้รับเครื่องหมายนี้ในการอัปเดตที่สำเร็จครั้งถัดไป
<!-- internationalizer:unit markdown:project-type-detection -->
## การตรวจหาประเภทโปรเจกต์

คำสั่ง `internationalizer detect` จะระบุการตั้งค่า i18n ของคุณโดยตรวจสอบจาก:

- dependencies ใน `package.json` สำหรับ react-i18next, next-intl หรือ vue-i18n
- โครงสร้างไดเรกทอรีที่ตรงกับแบบแผนทั่วไปของภาษา
- นามสกุลไฟล์และแบบแผนการตั้งชื่อ
<!-- internationalizer:unit markdown:architecture -->
## สถาปัตยกรรม

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
## เปรียบเทียบกับทางเลือกอื่น

| ฟีเจอร์ | Internationalizer | i18next | Crowdin | LLM ทั่วไป |
|---------|------------------|---------|---------|-------------|
| การแปลด้วยพลัง LLM | ใช่ | ไม่ | บางส่วน | ใช่ |
| คู่มือสไตล์เฉพาะแต่ละภาษา | ใช่ | ไม่ | ไม่ | ไม่ |
| การบังคับใช้อภิธานศัพท์ | ใช่ | ไม่ | ใช่ | ไม่ |
| หน่วยความจำการแปล | ใช่ | ไม่ | ใช่ | ไม่ |
| ทำงานผ่าน CLI / บนเครื่อง | ใช่ | N/A | ไม่ | ทำเอง |
| ไฟล์ที่เข้ากับ Git ได้ดี | ใช่ | ใช่ | บางส่วน | ทำเอง |
| ไม่ต้องพึ่งพา SaaS | ใช่ | ใช่ | ไม่ | แตกต่างกันไป |
| โอเพนซอร์ส (AGPL-3.0) | ใช่ | ใช่ | ไม่ | แตกต่างกันไป |
<!-- internationalizer:unit markdown:license -->
## สัญญาอนุญาต

[AGPL-3.0](../../LICENSE)

ดู [THIRD_PARTY_NOTICES.md](../../THIRD_PARTY_NOTICES.md) สำหรับประกาศเกี่ยวกับบุคคลที่สาม
<!-- internationalizer:unit markdown:contributing -->
## การมีส่วนร่วม

ดู [CONTRIBUTING.md](../../CONTRIBUTING.md) สำหรับการตั้งค่าสภาพแวดล้อมการพัฒนาและแนวทางปฏิบัติ การมีส่วนร่วมทั้งหมดจำเป็นต้องมีการลงนามรับรอง DCO (DCO sign-off)
