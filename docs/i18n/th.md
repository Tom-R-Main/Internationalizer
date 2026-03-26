> [English (original)](../../README.md)

<p align="center">
  <img src="../../assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

ไปป์ไลน์ internationalization แบบ AI-native สำหรับโปรเจกต์ซอฟต์แวร์ แปล ตรวจสอบ และจัดการไฟล์ i18n โดยใช้ LLM

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](LICENSE)

## ทำไมต้อง Internationalizer?

เครื่องมือ i18n ส่วนใหญ่มักจะเป็น runtime libraries (i18next, react-intl) หรือแพลตฟอร์ม SaaS สำหรับจัดการคีย์ (Crowdin, Lokalise) ซึ่งไม่มีเครื่องมือใดที่แก้ปัญหาการแปลได้อย่างแท้จริง:

- **การแปลด้วยตนเอง (Manual translation)** ไม่สามารถขยายสเกลได้เมื่อมีหลายภาษา
- **Machine translation APIs** (Google Translate, DeepL) ละเลยคำศัพท์ โทนเสียง และรูปแบบ UI ของคุณ
- **การแปลด้วย LLM ทั่วไป** ทำงานได้ดีกว่า แต่หากไม่มีอภิธานศัพท์ (glossaries) และคู่มือสไตล์ (style guides) ผลลัพธ์ที่ได้ก็จะไม่สม่ำเสมอ

Internationalizer นั้นแตกต่างออกไป มันคือ **CLI pipeline** ที่ผสานรวมการแปลด้วย LLM เข้ากับ:

- **อภิธานศัพท์แยกตามภาษา** — บังคับใช้คำศัพท์ให้สอดคล้องกันทั่วทั้งแอปของคุณ
- **คู่มือสไตล์แยกตามภาษา** — ควบคุมโทนเสียง ระดับความเป็นทางการ การทำพหูพจน์ และการจัดรูปแบบตัวอักษร
- **หน่วยความจำการแปล (Translation memory)** — ข้ามข้อความที่ไม่มีการเปลี่ยนแปลง ช่วยประหยัดค่าใช้จ่ายในการเรียก API
- **การตรวจสอบคีย์ (Key validation)** — ตรวจจับการแปลที่ตกหล่นและตัวแปร interpolation ที่ไม่ตรงกันก่อนนำไปใช้งานจริง

## การติดตั้ง

ติดตั้งจาก npm:

```bash
npm install -g internationalizer
```

หรือรันโดยไม่ต้องติดตั้งแบบ global:

```bash
npx internationalizer --help
```

แพ็กเกจ npm จะติดตั้ง prebuilt binary ที่ตรงกันจาก npm ผ่าน optional dependencies ที่เฉพาะเจาะจงกับแต่ละแพลตฟอร์ม

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

## แพ็กเกจ npm

- Git tags และเวอร์ชันของแพ็กเกจ npm ต้องตรงกัน ตัวอย่างเช่น `v0.1.0` และ `0.1.0`
- แพ็กเกจ `internationalizer` หลักจะขึ้นอยู่กับแพ็กเกจของแพลตฟอร์ม เช่น `internationalizer-darwin-arm64`
- เป้าหมาย npm ที่รองรับ: macOS arm64/x64, Linux arm64/x64, Windows x64
- การเผยแพร่ผ่าน CI จำเป็นต้องใช้ GitHub secret ที่ชื่อว่า `NPM_TOKEN`

## เริ่มต้นใช้งานอย่างรวดเร็ว

1. สร้างไฟล์คอนฟิกใน root ของโปรเจกต์คุณ:

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

2. ตั้งค่า API key ของคุณ:

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

3. ดูตัวอย่างสิ่งที่จะถูกแปล:

```bash
internationalizer translate --dry-run
```

4. รันการแปล:

```bash
internationalizer translate
```

5. ตรวจสอบ locales ทั้งหมด:

```bash
internationalizer validate
```

## คำสั่ง

### `translate`

ค้นหาคีย์ที่หายไปและแปลผ่าน LLM

```bash
internationalizer translate                    # แปล locales ทั้งหมด
internationalizer translate -l fr              # แปลเฉพาะภาษาฝรั่งเศส
internationalizer translate --dry-run          # ดูตัวอย่างโดยไม่เรียก API
internationalizer translate --batch-size 20    # ลดขนาด batch
internationalizer translate --concurrency 2    # ลดการเรียกพร้อมกัน (parallel)
```

### `validate`

ตรวจสอบไฟล์ locale ทั้งหมดเพื่อหาคีย์ที่หายไป คีย์ที่เกินมา และตัวแปร interpolation ที่ไม่ตรงกัน

```bash
internationalizer validate                     # แสดงผลลัพธ์ที่อ่านง่ายสำหรับมนุษย์
internationalizer validate --json              # แสดงผลลัพธ์เป็น JSON สำหรับเครื่อง
internationalizer validate -q                  # แสดงเฉพาะ exit code
```

### `detect`

ตรวจจับเฟรมเวิร์ก i18n อัตโนมัติและแนะนำการตั้งค่า

```bash
internationalizer detect
```

รองรับ: react-i18next, next-intl, vue-i18n, vanilla JSON, markdown docs

### `glossary`

จัดการคำศัพท์ในอภิธานศัพท์แยกตามภาษา ซึ่งจะถูกบังคับใช้ในระหว่างการแปล

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

จัดการหน่วยความจำการแปล (แคช JSONL ของข้อความที่เคยแปลไปแล้ว)

```bash
internationalizer tm stats                     # แสดงจำนวนเรกคอร์ด
internationalizer tm export                    # ดัมป์เป็น JSON
internationalizer tm clear --force             # ลบเรกคอร์ดทั้งหมด
```

## ข้อมูลอ้างอิงการตั้งค่า (Configuration Reference)

```yaml
# .internationalizer.yml

# ภาษาต้นทาง (ค่าเริ่มต้น: en)
source_locale: en

# ภาษาที่ต้องการแปล (จำเป็น)
target_locales: [fr, de, es, ja, zh-CN, ar]

# พาธไปยังไฟล์ locale ต้นทาง (จำเป็น)
source_path: locales/en.json

# การตั้งค่าผู้ให้บริการ LLM
llm:
  # ผู้ให้บริการ: "anthropic", "openai", "gemini" หรือ "openrouter" (ค่าเริ่มต้น: gemini)
  provider: gemini

  # ชื่อโมเดลเริ่มต้นตามผู้ให้บริการ:
  #   anthropic:  claude-sonnet-4-6
  #   openai:     gpt-5.4
  #   gemini:     gemini-3.1-pro-preview
  #   openrouter: google/gemini-3-flash-preview
  model: gemini-3.1-pro-preview

  # ตัวแปรสภาพแวดล้อม (Environment variable) ที่เก็บ API key
  api_key_env: GOOGLE_AI_STUDIO_API_KEY

  # Base URL สำหรับ endpoint ที่เข้ากันได้กับ OpenAI (ไม่บังคับ)
  # base_url: https://api.openai.com

# จำนวนคีย์ต่อการเรียก LLM หนึ่งครั้ง (ค่าเริ่มต้น: 40)
batch_size: 40

# การเรียก LLM พร้อมกัน (ค่าเริ่มต้น: 4)
concurrency: 4

# ไดเรกทอรีที่เก็บไฟล์ Markdown คู่มือสไตล์แยกตามภาษา (ค่าเริ่มต้น: style-guides)
style_guides_dir: style-guides

# ไดเรกทอรีที่เก็บไฟล์ JSON อภิธานศัพท์แยกตามภาษา (ค่าเริ่มต้น: glossary)
glossary_dir: glossary

# พาธไปยังไฟล์หน่วยความจำการแปล (ค่าเริ่มต้น: .internationalizer/tm.jsonl)
tm_path: .internationalizer/tm.jsonl
```

## คู่มือสไตล์ (Style Guides)

คู่มือสไตล์คือไฟล์ Markdown ที่จะถูกแทรกเข้าไปใน prompt การแปลของ LLM ซึ่งจะช่วยควบคุมโทนเสียง ระดับความเป็นทางการ การจัดรูปแบบตัวอักษร และธรรมเนียมปฏิบัติเฉพาะของแต่ละภาษา

```
style-guides/
  _conventions.md    # กฎที่ใช้ร่วมกันสำหรับทุกภาษา
  fr.md              # กฎเฉพาะสำหรับภาษาฝรั่งเศส
  ja.md              # กฎเฉพาะสำหรับภาษาญี่ปุ่น
  ar.md              # กฎเฉพาะสำหรับภาษาอาหรับ
```

### ธรรมเนียมปฏิบัติที่ใช้ร่วมกัน (`_conventions.md`)

กำหนดกฎที่ใช้กับทุกภาษา: ไวยากรณ์ interpolation, การคงแท็ก HTML ไว้, ธรรมเนียมปฏิบัติของประเภทข้อความ (ปุ่ม vs. ป้ายกำกับ vs. ข้อผิดพลาด) เป็นต้น

### คู่มือแยกตามภาษา (`{locale}.md`)

กำหนดกฎเฉพาะของแต่ละภาษา: ระดับความเป็นทางการ (tu vs. vous), เครื่องหมายวรรคตอน (guillemets, เครื่องหมายคำถามกลับหัว), รูปพหูพจน์, การจัดรูปแบบวันที่/ตัวเลข และอภิธานศัพท์

ดูตัวอย่างการใช้งานจริงได้ที่ [`examples/react-app/style-guides/`](examples/react-app/style-guides/)

## รูปแบบอภิธานศัพท์ (Glossary Format)

ไฟล์อภิธานศัพท์เป็น JSON arrays ที่จัดเก็บไว้ใน `{glossary_dir}/{locale}.json`:

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

คำศัพท์ต่างๆ จะถูกแทรกเข้าไปใน prompt ของ LLM ในรูปแบบตารางคำศัพท์ เพื่อให้แน่ใจว่าคำศัพท์สำคัญต่างๆ จะถูกแปลอย่างสอดคล้องกันทั่วทั้งแอปพลิเคชันของคุณ

## หน่วยความจำการแปล (Translation Memory)

หน่วยความจำการแปลจะถูกจัดเก็บเป็นไฟล์ JSONL (หนึ่ง JSON record ต่อบรรทัด) แต่ละเรกคอร์ดประกอบด้วย:

- คีย์และค่าของต้นทาง
- ค่าที่แปลแล้ว
- SHA-256 hash ของค่าต้นทาง
- Timestamp

ในการรันครั้งถัดไป ข้อความที่ไม่มีการเปลี่ยนแปลงจะถูกดึงมาจากแคช TM โดยไม่ต้องเรียก LLM ซึ่งช่วยประหยัดทั้งเวลาและค่าใช้จ่าย API ไฟล์ TM นี้เป็นมิตรกับ git และสามารถ commit ไปพร้อมกับไฟล์ locale ของคุณได้

## รูปแบบที่รองรับ

| รูปแบบ | นามสกุลไฟล์ | โหมด |
|--------|-----------|------|
| JSON | `.json` | Key-value (ซ้อนทับกันได้, แบนราบด้วย dot-notation) |
| YAML | `.yml`, `.yaml` | Key-value (คงคอมเมนต์และลำดับไว้) |
| Markdown | `.md`, `.mdx` | แปลทั้งเอกสาร |

## การตรวจจับประเภทโปรเจกต์

`internationalizer detect` จะระบุการตั้งค่า i18n ของคุณโดยตรวจสอบ:

- dependencies ใน `package.json` สำหรับ react-i18next, next-intl หรือ vue-i18n
- โครงสร้างไดเรกทอรีที่ตรงกับรูปแบบ locale ทั่วไป
- นามสกุลไฟล์และธรรมเนียมการตั้งชื่อ

## สถาปัตยกรรม (Architecture)

```
cmd/internationalizer/     จุดเริ่มต้นของ CLI และการกำหนดคำสั่ง
internal/
  config/                  โหลดคอนฟิก YAML พร้อมค่าเริ่มต้น
  detect/                  ตรวจจับประเภทโปรเจกต์อัตโนมัติ
  formats/                 ตัวแยกวิเคราะห์รูปแบบ (JSON, YAML, Markdown)
  glossary/                การจัดการอภิธานศัพท์แยกตามภาษา
  llm/                     อินเทอร์เฟซผู้ให้บริการ LLM + การนำไปใช้งาน
    anthropic.go           แบ็กเอนด์ Anthropic Claude
    openai.go              แบ็กเอนด์ OpenAI / ที่เข้ากันได้
    gemini.go              แบ็กเอนด์ Google Gemini ผ่าน AI Studio
                           OpenRouter ใช้ openai.go พร้อม base_url แบบกำหนดเอง
  styleguide/              ตัวโหลดคู่มือสไตล์
  tm/                      หน่วยความจำการแปล JSONL
  translate/               ตัวจัดการการแปล
  validate/                การตรวจสอบ locale และการหาความแตกต่าง
```

## เปรียบเทียบกับทางเลือกอื่น

| ฟีเจอร์ | Internationalizer | i18next | Crowdin | LLM ทั่วไป |
|---------|------------------|---------|---------|-------------|
| การแปลด้วย LLM | ใช่ | ไม่ | บางส่วน | ใช่ |
| คู่มือสไตล์แยกตามภาษา | ใช่ | ไม่ | ไม่ | ไม่ |
| การบังคับใช้อภิธานศัพท์ | ใช่ | ไม่ | ใช่ | ไม่ |
| หน่วยความจำการแปล | ใช่ | ไม่ | ใช่ | ไม่ |
| CLI / รันบนเครื่อง | ใช่ | N/A | ไม่ | ทำเอง |
| ไฟล์ที่เป็นมิตรกับ Git | ใช่ | ใช่ | บางส่วน | ทำเอง |
| ไม่พึ่งพา SaaS | ใช่ | ใช่ | ไม่ | แตกต่างกันไป |
| โอเพนซอร์ส (AGPL-3.0) | ใช่ | ใช่ | ไม่ | แตกต่างกันไป |

## ไลเซนส์ (License)

[AGPL-3.0](LICENSE)

## การมีส่วนร่วม (Contributing)

ดู [CONTRIBUTING.md](CONTRIBUTING.md) สำหรับการตั้งค่าการพัฒนาและแนวทางปฏิบัติ การมีส่วนร่วมทั้งหมดจำเป็นต้องมีการลงนาม DCO (DCO sign-off)

