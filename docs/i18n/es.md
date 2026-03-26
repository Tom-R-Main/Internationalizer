> [English (original)](../../README.md)

<p align="center">
  <img src="../../assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

Pipeline de internacionalización nativo de IA para proyectos de software. Traduce, valida y gestiona archivos i18n usando LLMs.

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](LICENSE)

## ¿Por qué Internationalizer?

La mayoría de las herramientas i18n son bibliotecas de tiempo de ejecución (i18next, react-intl) o plataformas SaaS de gestión de claves (Crowdin, Lokalise). Ninguna de ellas resuelve bien el problema real de la traducción:

- **La traducción manual** no escala más allá de unos pocos idiomas
- **Las API de traducción automática** (Google Translate, DeepL) ignoran tu terminología, tono y convenciones de la interfaz de usuario
- **La traducción genérica con LLM** funciona mejor, pero sin glosarios ni guías de estilo, obtienes resultados inconsistentes

Internationalizer es diferente. Es un **pipeline CLI** que combina la traducción con LLM con:

- **Glosarios por idioma**: aplican una terminología coherente en toda tu aplicación
- **Guías de estilo por idioma**: controlan el tono, la formalidad, la pluralización y la tipografía
- **Memoria de traducción**: omite las cadenas sin cambios, ahorrando dinero en llamadas a la API
- **Validación de claves**: detecta traducciones faltantes y desajustes de interpolación antes de su publicación

## Instalación

Instalar desde npm:

```bash
npm install -g internationalizer
```

O ejecutar sin una instalación global:

```bash
npx internationalizer --help
```

El paquete npm instala el binario precompilado correspondiente desde npm a través de dependencias opcionales específicas de la plataforma.

Instalar con Go:

```bash
go install github.com/Tom-R-Main/Internationalizer/cmd/internationalizer@latest
```

O compilar desde el código fuente:

```bash
git clone https://github.com/Tom-R-Main/Internationalizer.git
cd Internationalizer
go build -o internationalizer ./cmd/internationalizer
```

## Paquetes npm

- Las etiquetas de Git y las versiones de los paquetes npm deben coincidir, por ejemplo `v0.1.0` y `0.1.0`
- El paquete raíz `internationalizer` depende de paquetes de plataforma como `internationalizer-darwin-arm64`
- Objetivos npm compatibles: macOS arm64/x64, Linux arm64/x64, Windows x64
- La publicación en CI requiere un secreto de GitHub llamado `NPM_TOKEN`

## Inicio rápido

1. Crea un archivo de configuración en la raíz de tu proyecto:

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

2. Configura tu clave API:

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

3. Previsualiza lo que se traducirá:

```bash
internationalizer translate --dry-run
```

4. Ejecuta la traducción:

```bash
internationalizer translate
```

5. Valida todos los locales:

```bash
internationalizer validate
```

## Comandos

### `translate`

Encuentra claves faltantes y tradúcelas a través de un LLM.

```bash
internationalizer translate                    # traduce todos los locales
internationalizer translate -l fr              # traduce solo al francés
internationalizer translate --dry-run          # previsualiza sin llamadas a la API
internationalizer translate --batch-size 20    # lotes más pequeños
internationalizer translate --concurrency 2    # menos llamadas en paralelo
```

### `validate`

Comprueba todos los archivos de locales en busca de claves faltantes, claves adicionales y desajustes de interpolación.

```bash
internationalizer validate                     # salida legible para humanos
internationalizer validate --json              # JSON legible por máquina
internationalizer validate -q                  # solo código de salida
```

### `detect`

Detecta automáticamente el framework i18n y sugiere una configuración.

```bash
internationalizer detect
```

Soporta: react-i18next, next-intl, vue-i18n, JSON puro, documentos markdown.

### `glossary`

Gestiona los términos del glosario por idioma que se aplican durante la traducción.

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

Gestiona la memoria de traducción (caché JSONL de cadenas traducidas previamente).

```bash
internationalizer tm stats                     # muestra el recuento de registros
internationalizer tm export                    # vuelca como JSON
internationalizer tm clear --force             # elimina todos los registros
```

## Referencia de configuración

```yaml
# .internationalizer.yml

# Idioma de origen (por defecto: en)
source_locale: en

# Idiomas a los que traducir (requerido)
target_locales: [fr, de, es, ja, zh-CN, ar]

# Ruta al archivo de locale de origen (requerido)
source_path: locales/en.json

# Configuración del proveedor de LLM
llm:
  # Proveedor: "anthropic", "openai", "gemini" u "openrouter" (por defecto: gemini)
  provider: gemini

  # Nombres de modelo por defecto según el proveedor:
  #   anthropic:  claude-sonnet-4-6
  #   openai:     gpt-5.4
  #   gemini:     gemini-3.1-pro-preview
  #   openrouter: google/gemini-3-flash-preview
  model: gemini-3.1-pro-preview

  # Variable de entorno que contiene la clave API
  api_key_env: GOOGLE_AI_STUDIO_API_KEY

  # URL base para endpoints compatibles con OpenAI (opcional)
  # base_url: https://api.openai.com

# Claves por llamada al LLM (por defecto: 40)
batch_size: 40

# Llamadas paralelas al LLM (por defecto: 4)
concurrency: 4

# Directorio que contiene los archivos Markdown de la guía de estilo por locale (por defecto: style-guides)
style_guides_dir: style-guides

# Directorio que contiene los archivos JSON del glosario por locale (por defecto: glossary)
glossary_dir: glossary

# Ruta al archivo de memoria de traducción (por defecto: .internationalizer/tm.jsonl)
tm_path: .internationalizer/tm.jsonl
```

## Guías de estilo

Las guías de estilo son archivos Markdown que se inyectan en el prompt de traducción del LLM. Controlan el tono, la formalidad, la tipografía y otras convenciones específicas del idioma.

```
style-guides/
  _conventions.md    # reglas compartidas para todos los idiomas
  fr.md              # reglas específicas para francés
  ja.md              # reglas específicas para japonés
  ar.md              # reglas específicas para árabe
```

### Convenciones compartidas (`_conventions.md`)

Define reglas que se aplican a todos los idiomas: sintaxis de interpolación, preservación de HTML, convenciones de tipos de cadenas (botones vs. etiquetas vs. errores), etc.

### Guías por idioma (`{locale}.md`)

Define reglas específicas del idioma: registro de formalidad (tú vs. usted), puntuación (comillas angulares, signos de interrogación invertidos), formas plurales, formato de fechas/números y un glosario terminológico.

Consulta [`examples/react-app/style-guides/`](examples/react-app/style-guides/) para ver un ejemplo funcional.

## Formato del glosario

Los archivos de glosario son arrays JSON almacenados en `{glossary_dir}/{locale}.json`:

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

Los términos se inyectan en el prompt del LLM como una tabla terminológica, asegurando una traducción coherente de los términos clave en toda tu aplicación.

## Memoria de traducción

La memoria de traducción se almacena como un archivo JSONL (un registro JSON por línea). Cada registro contiene:

- La clave y el valor de origen
- El valor traducido
- Un hash SHA-256 del valor de origen
- Una marca de tiempo

En ejecuciones posteriores, las cadenas sin cambios se sirven desde la caché de la memoria de traducción (TM) sin llamar al LLM, ahorrando tiempo y costes de API. El archivo TM es compatible con git y se puede confirmar junto con tus archivos de locales.

## Formatos soportados

| Formato | Extensiones | Modo |
|--------|-----------|------|
| JSON | `.json` | Clave-valor (anidado, aplanado con notación de puntos) |
| YAML | `.yml`, `.yaml` | Clave-valor (preserva comentarios y orden) |
| Markdown | `.md`, `.mdx` | Traducción de documento completo |

## Detección del tipo de proyecto

`internationalizer detect` identifica tu configuración i18n comprobando:

- Las dependencias de `package.json` para react-i18next, next-intl o vue-i18n
- Estructuras de directorios que coinciden con patrones de locales comunes
- Extensiones de archivo y convenciones de nomenclatura

## Arquitectura

```
cmd/internationalizer/     Punto de entrada de la CLI y definiciones de comandos
internal/
  config/                  Carga de configuración YAML con valores por defecto
  detect/                  Detección automática del tipo de proyecto
  formats/                 Analizadores de formato (JSON, YAML, Markdown)
  glossary/                Gestión del glosario por locale
  llm/                     Interfaz de proveedor de LLM + implementaciones
    anthropic.go           Backend de Anthropic Claude
    openai.go              Backend de OpenAI / compatible
    gemini.go              Backend de Google Gemini vía AI Studio
                           OpenRouter usa openai.go con base_url personalizada
  styleguide/              Cargador de guías de estilo
  tm/                      Memoria de traducción JSONL
  translate/               Orquestador de traducción
  validate/                Validación de locales y diferencias
```

## Comparación con alternativas

| Característica | Internationalizer | i18next | Crowdin | LLM genérico |
|---------|------------------|---------|---------|-------------|
| Traducción impulsada por LLM | Sí | No | Parcial | Sí |
| Guías de estilo por idioma | Sí | No | No | No |
| Aplicación de glosario | Sí | No | Sí | No |
| Memoria de traducción | Sí | No | Sí | No |
| CLI / ejecución local | Sí | N/A | No | Manual |
| Archivos compatibles con Git | Sí | Sí | Parcial | Manual |
| Sin dependencia de SaaS | Sí | Sí | No | Varía |
| Código abierto (AGPL-3.0) | Sí | Sí | No | Varía |

## Licencia

[AGPL-3.0](LICENSE)

## Contribución

Consulta [CONTRIBUTING.md](CONTRIBUTING.md) para la configuración de desarrollo y las directrices. Todas las contribuciones requieren la firma del DCO.

