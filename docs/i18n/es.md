> [English (original)](../../README.md)

<p align="center">
  <img src="../../assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

Pipeline de internacionalización nativo para IA en proyectos de software. Traduzca, valide y gestione archivos i18n mediante LLM.

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](../../LICENSE)

<p align="center">
<a href="ar.md">العربية</a> · <a href="bn.md">বাংলা</a> · <a href="cs.md">Čeština</a> · <a href="da.md">Dansk</a> · <a href="de.md">Deutsch</a> · <a href="el.md">Ελληνικά</a> · <a href="es.md">Español</a> · <a href="fi.md">Suomi</a> · <a href="fr.md">Français</a> · <a href="he.md">עברית</a> · <a href="hi.md">हिन्दी</a> · <a href="id.md">Indonesia</a> · <a href="it.md">Italiano</a> · <a href="ja.md">日本語</a> · <a href="ko.md">한국어</a> · <a href="ms.md">Bahasa Melayu</a><br><a href="nl.md">Nederlands</a> · <a href="pa.md">ਪੰਜਾਬੀ</a> · <a href="pl.md">Polski</a> · <a href="pt-BR.md">Português</a> · <a href="ro.md">Română</a> · <a href="ru.md">Русский</a> · <a href="sv.md">Svenska</a> · <a href="te.md">తెలుగు</a> · <a href="th.md">ไทย</a> · <a href="tr.md">Türkçe</a> · <a href="uk.md">Українська</a> · <a href="vi.md">Tiếng Việt</a> · <a href="yue.md">粵語</a> · <a href="zh-CN.md">简体中文</a> · <a href="zh-TW.md">繁體中文</a>
</p>

---
<!-- internationalizer:unit markdown:why-internationalizer -->
## ¿Por qué Internationalizer?

La mayoría de las herramientas de i18n son bibliotecas de entorno de ejecución (i18next, react-intl) o plataformas SaaS de gestión de claves (Crowdin, Lokalise). Ninguna resuelve adecuadamente el problema real de la traducción:

- **La traducción manual** no escala más allá de unos pocos idiomas
- **Las API de traducción automática** (Google Translate, DeepL) ignoran la terminología, el tono y las convenciones de la interfaz de usuario de su proyecto
- **La traducción genérica con LLM** funciona mejor, pero sin glosarios ni guías de estilo produce resultados incoherentes

Internationalizer es diferente. Es una **canalización por línea de comandos (CLI)** que combina la traducción mediante LLM con:

- **Glosarios por idioma**: garantizan una terminología coherente en toda la aplicación
- **Guías de estilo por idioma**: controlan el tono, la formalidad, la pluralización y la tipografía
- **Memoria de traducción**: omite cadenas sin cambios para ahorrar costes en llamadas a la API
- **Validación determinista**: detecta claves faltantes o sobrantes, desviaciones en estructuras protegidas, problemas de glosario y errores de plurales o de sintaxis ICU antes de pasar a producción
<!-- internationalizer:unit markdown:installation -->
## Instalación

Instale el paquete desde npm:

```bash
npm install -g internationalizer
```

O ejecútelo directamente sin instalación global:

```bash
npx internationalizer --help
```

El paquete npm instala el binario precompilado correspondiente desde npm mediante dependencias opcionales específicas para cada plataforma.

Instalación con Go:

```bash
go install github.com/Tom-R-Main/Internationalizer/cmd/internationalizer@latest
```

O compile directamente desde el código fuente:

```bash
git clone https://github.com/Tom-R-Main/Internationalizer.git
cd Internationalizer
go build -o internationalizer ./cmd/internationalizer
```
<!-- internationalizer:unit markdown:npm-packages -->
## Paquetes npm

- Las etiquetas de Git y las versiones de los paquetes npm deben coincidir exactamente; por ejemplo, `v0.1.0` y `0.1.0`
- El paquete raíz `internationalizer` depende de paquetes específicos de plataforma como `internationalizer-darwin-arm64`
- Destinos de npm compatibles: macOS arm64/x64, Linux arm64/x64 y Windows x64
- La publicación en CI requiere un secreto de GitHub denominado `NPM_TOKEN`
<!-- internationalizer:unit markdown:quick-start -->
## Inicio rápido

1. Cree un archivo de configuración en la raíz de su proyecto:

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

2. Configure su clave de API:

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

3. Previsualice los elementos que se traducirán:

```bash
internationalizer translate --dry-run
```

4. Ejecute la traducción:

```bash
internationalizer translate
```

5. Valide todas las configuraciones regionales:

```bash
internationalizer validate
```
<!-- internationalizer:unit markdown:commands -->
## Comandos

### `translate`

Detecta claves que faltan o están desactualizadas y las traduce mediante un LLM.

```bash
internationalizer translate                    # translate all locales
internationalizer translate -l fr              # translate French only
internationalizer translate --dry-run          # preview without API calls
internationalizer translate --adopt-existing   # baseline existing translations without API calls
internationalizer translate --refresh-policy   # refresh prompt/style/model-stale entries
internationalizer translate --batch-size 20    # smaller batches
internationalizer translate --concurrency 2    # fewer parallel calls
```

El estado de la traducción notifica de forma independiente las condiciones de ausencia, origen desactualizado, directiva desactualizada, estado al día y ediciones manuales, de modo que una edición manual nunca oculte un cambio de origen o de directiva. Los valores con directiva desactualizada se notifican, pero solo se vuelven a traducir con `--refresh-policy`. Los valores modificados manualmente nunca se sobrescriben de forma automática. Use `--adopt-existing` al incorporar el manifiesto a traducciones ya revisadas o al aceptar explícitamente una edición manual revisada como nueva base de referencia.

### `validate`

Compara todos los archivos de configuración regional con sus paquetes de origen. La validación predeterminada comprueba la cobertura estructural (el porcentaje de claves de destino requeridas que están presentes), informa sobre claves sobrantes en forma de advertencias y genera un error si faltan claves, si hay discrepancias de interpolación o si la estructura ICU MessageFormat no es válida.

```bash
internationalizer validate                     # human-readable output
internationalizer validate --json              # machine-readable JSON
internationalizer validate -q                  # exit code only
internationalizer validate --strict             # enforce translation quality rules
internationalizer validate --require-state      # require current manifest provenance
```

`--strict` también informa sobre la cobertura traducida. Un valor lingüístico idéntico a su origen se considera no traducido a menos que el glosario contenga de forma explícita una entrada idéntica con el mismo origen y el mismo destino para el valor completo; se respeta `ignore_case`, pero un término de glosario insertado dentro de un valor más largo no constituye una excepción. El modo estricto falla ante claves sobrantes, valores idénticos al origen, alteraciones en la estructura de interpolación, HTML, código o vínculos Markdown, infracciones de glosario y formas plurales configuradas.

`--require-state` comprueba cada destino frente a `.internationalizer.lock`. Genera un error si una clave no está registrada o si su origen registrado, directiva de traducción o hash de destino están desactualizados. Se puede combinar con `--strict`.

Los informes legibles por humanos y en JSON utilizan códigos de resultado estables:

| Código | Significado |
| --- | --- |
| `missing_key` / `extra_key` | Los conjuntos de claves de origen y destino difieren |
| `blank_translation` | Un origen no vacío tiene un destino vacío en modo estricto |
| `source_identical` | Un valor lingüístico en modo estricto permanece sin traducir |
| `protected_structure_mismatch` | La estructura de interpolación, HTML, código o vínculo ha cambiado |
| `glossary_violation` | No se encontró ningún término o variante aprobado en el destino |
| `plural_form_missing` | Falta una forma plural configurada para el idioma de destino |
| `icu_message_syntax` | Un mensaje ICU de origen o de destino tiene un formato incorrecto |
| `icu_argument_mismatch` | Los nombres, tipos o estilos de formato de los argumentos ICU difieren |
| `icu_selector_mismatch` | Los selectores difieren o una categoría plural no es válida para la configuración regional de destino |
| `untracked` | No existe ningún registro en el manifiesto para el destino |
| `source_stale` | El contenido de origen cambió después de registrarse la traducción |
| `policy_stale` | La petición (prompt) generada o la configuración del modelo han cambiado |
| `target_modified` | El contenido de destino difiere del registro del manifiesto |

### `detect`

Detecta automáticamente el marco de i18n utilizado y sugiere una configuración.

```bash
internationalizer detect
```

Soporta: react-i18next, next-intl, vue-i18n, JSON sin dependencias (vanilla JSON), documentos Markdown.

### `glossary`

Gestiona términos de glosario por idioma cuyo cumplimiento se exige durante la traducción.

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

Gestiona la memoria de traducción (caché JSONL de cadenas traducidas previamente).

```bash
internationalizer tm stats                     # show record counts
internationalizer tm export                    # dump as JSON
internationalizer tm clear --force             # delete all records
```
<!-- internationalizer:unit markdown:configuration-reference -->
## Referencia de configuración

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

Los identificadores de configuración regional deben ser etiquetas BCP 47 con formato correcto, como `fr`, `pt-BR` o `sr-Latn-RS`. Las configuraciones regionales de destino canónicamente equivalentes se rechazan por duplicadas, y las invalidaciones de proveedor específicas por configuración regional coinciden con la grafía canónicamente equivalente. En el ejemplo anterior, las configuraciones regionales sin una invalidación (incluido el japonés) heredan la configuración global de Gemini.

Los valores de ICU MessageFormat se analizan estructuralmente. Se admiten argumentos simples, `select`, `plural`, `selectordinal`, `number`, `date` y `time`, incluidos mensajes anidados, desplazamientos de plural (offsets), selectores de números exactos y `#`. La validación comprueba la sintaxis, los tipos de argumentos y estilos de formateo, los desplazamientos de plural, la identidad de las ramas de selección y las categorías de plurales de CLDR para el idioma de destino. Las respuestas del proveedor que vulneren estas invariantes se rechazan antes de escribir cualquier archivo de configuración regional o registro en la memoria de traducción.

Con `i18next-v4`, las familias de plurales reconocidas en el origen se expanden durante la traducción hacia las categorías CLDR correspondientes a la configuración regional de destino. Las categorías exclusivas del destino toman el valor `_other` de la familia de origen como plantilla de traducción. La validación estricta exige la presencia de esas categorías de destino; las categorías exclusivas del origen son opcionales para las configuraciones regionales de destino que no las empleen.
<!-- internationalizer:unit markdown:style-guides -->
## Guías de estilo

Las guías de estilo son archivos Markdown que se inyectan en la petición (prompt) de traducción del LLM. Controlan el tono, la formalidad, la tipografía y otras convenciones específicas de cada idioma.

```
style-guides/
  _conventions.md    # shared rules for all languages
  fr.md              # French-specific rules
  ja.md              # Japanese-specific rules
  ar.md              # Arabic-specific rules
```

### Convenciones compartidas (`_conventions.md`)

Define las reglas aplicables a todos los idiomas: sintaxis de interpolación, preservación de HTML, convenciones de tipos de cadenas (botones frente a etiquetas o errores), etc.

### Guías por idioma (`{locale}.md`)

Define las reglas específicas para cada idioma: registro de formalidad (tú frente a usted), puntuación (comillas angulares, signos de apertura de interrogación), formas plurales, formatos de fecha y números, y un glosario terminológico.

Las guías de estilo constituyen entradas normativas y duraderas, no resultados generados. Internationalizer las lee, pero nunca las sobrescribe. Su contenido se procesa mediante un hash independiente del glosario y del contrato de petición, por lo que una modificación en el código de la aplicación no marca una traducción como desactualizada. La modificación deliberada de una guía marca ese idioma para revisión de directiva; los cambios en la redacción interna de la petición no lo hacen, a menos que también varíe la versión del contrato de petición.

Consulte [`examples/react-app/style-guides/`](../../examples/react-app/style-guides/) para ver un ejemplo funcional.
<!-- internationalizer:unit markdown:glossary-format -->
## Formato de glosario

Los archivos de glosario son matrices JSON almacenadas en `{glossary_dir}/{locale}.json`:

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

`variants` enumera otras formas de destino autorizadas. `enforcement` puede establecerse en `error`, `warning` o bien omitirse para conservar el comportamiento predeterminado de error. Los términos se incorporan al prompt del LLM en forma de tabla de terminología, lo que garantiza una traducción coherente en toda la aplicación. Una entrada exacta como `{"source":"API","target":"API"}` también exime a ese valor completo idéntico al origen de ser marcado como no traducido en el modo estricto; sin embargo, no exime a un valor más extenso que simplemente contenga `API`.
<!-- internationalizer:unit markdown:translation-memory -->
## Memoria de traducción

La memoria de traducción se almacena en un archivo JSONL (un registro JSON por línea). Cada registro contiene:

- El paquete, la clave, el valor de origen, el valor traducido y la configuración regional de destino canónica
- Los valores hash del origen, la guía de estilo, el glosario, el contrato de la petición y la directiva combinada
- El proveedor y el modelo que produjeron la traducción
- Una marca de tiempo

En las ejecuciones posteriores, las cadenas con los mismos hashes de origen y de directiva se recuperan de la memoria caché sin realizar llamadas al LLM. La ruta predeterminada se sitúa bajo el directorio omitido `.internationalizer/`, por lo que se mantiene como una caché local. Asigne a `tm_path` una ruta bajo control de versiones si en su proyecto se decide compartir deliberadamente la memoria de traducción. El manifiesto auditable `.internationalizer.lock` se gestiona por versiones de forma independiente.
<!-- internationalizer:unit markdown:supported-formats -->
## Formatos admitidos

| Formato | Extensiones | Modo |
| --- | --- | --- |
| JSON | `.json` | Clave-valor (anidado, aplanado con notación de puntos) |
| YAML | `.yml`, `.yaml` | Clave-valor (conserva comentarios y ordenación) |
| Markdown | `.md`, `.mdx` | Preámbulo y secciones de nivel H2 |

Los archivos de destino Markdown incluyen comentarios invisibles `internationalizer:unit` antes de las secciones H2. Estos marcadores estables permiten a Internationalizer añadir, mover o editar una sección de origen sin volver a traducir secciones no relacionadas. Los documentos existentes que carezcan de marcas las recibirán en su siguiente actualización correcta.
<!-- internationalizer:unit markdown:project-type-detection -->
## Detección del tipo de proyecto

`internationalizer detect` identifica la configuración de su framework de i18n examinando:

- Las dependencias del archivo `package.json` en busca de react-i18next, next-intl o vue-i18n
- Estructuras de directorios que coincidan con patrones habituales de localización
- Extensiones de archivo y convenciones de nomenclatura
<!-- internationalizer:unit markdown:architecture -->
## Arquitectura

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
## Comparación con alternativas

| Característica | Internationalizer | i18next | Crowdin | LLM genérico |
| --- | --- | --- | --- | --- |
| Traducción mediante LLM | Sí | No | Parcial | Sí |
| Guías de estilo por idioma | Sí | No | No | No |
| Aplicación obligatoria de glosario | Sí | No | Sí | No |
| Memoria de traducción | Sí | No | Sí | No |
| CLI / ejecución local | Sí | N/A | No | Manual |
| Archivos compatibles con Git | Sí | Sí | Parcial | Manual |
| Sin dependencia de servicios SaaS | Sí | Sí | No | Variable |
| Código abierto (AGPL-3.0) | Sí | Sí | No | Variable |
<!-- internationalizer:unit markdown:license -->
## Licencia

[AGPL-3.0](../../LICENSE)

Consulte [THIRD_PARTY_NOTICES.md](../../THIRD_PARTY_NOTICES.md) para ver los avisos sobre dependencias de terceros.
<!-- internationalizer:unit markdown:contributing -->
## Cómo contribuir

Consulte [CONTRIBUTING.md](../../CONTRIBUTING.md) para obtener información sobre la configuración del entorno de desarrollo y las directrices de contribución. Todas las contribuciones requieren la firma del DCO.
