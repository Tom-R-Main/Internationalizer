> [English (original)](../../README.md)

<p align="center">
  <img src="../../assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

Pipeline de internacionalização nativo de IA para projetos de software. Traduza, valide e gerencie arquivos de i18n usando LLMs.

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](../../LICENSE)

<p align="center">
<a href="ar.md">العربية</a> · <a href="bn.md">বাংলা</a> · <a href="cs.md">Čeština</a> · <a href="da.md">Dansk</a> · <a href="de.md">Deutsch</a> · <a href="el.md">Ελληνικά</a> · <a href="es.md">Español</a> · <a href="fi.md">Suomi</a> · <a href="fr.md">Français</a> · <a href="he.md">עברית</a> · <a href="hi.md">हिन्दी</a> · <a href="id.md">Indonesia</a> · <a href="it.md">Italiano</a> · <a href="ja.md">日本語</a> · <a href="ko.md">한국어</a> · <a href="ms.md">Bahasa Melayu</a><br><a href="nl.md">Nederlands</a> · <a href="pa.md">ਪੰਜਾਬੀ</a> · <a href="pl.md">Polski</a> · <a href="pt-BR.md">Português</a> · <a href="ro.md">Română</a> · <a href="ru.md">Русский</a> · <a href="sv.md">Svenska</a> · <a href="te.md">తెలుగు</a> · <a href="th.md">ไทย</a> · <a href="tr.md">Türkçe</a> · <a href="uk.md">Українська</a> · <a href="vi.md">Tiếng Việt</a> · <a href="yue.md">粵語</a> · <a href="zh-CN.md">简体中文</a> · <a href="zh-TW.md">繁體中文</a>
</p>

---
<!-- internationalizer:unit markdown:why-internationalizer -->
## Por que o Internationalizer?

A maioria das ferramentas de i18n é composta por bibliotecas de runtime (i18next, react-intl) ou plataformas SaaS de gerenciamento de chaves (Crowdin, Lokalise). Nenhuma delas resolve bem o verdadeiro problema de tradução:

- **A tradução manual** não escala além de alguns idiomas
- **APIs de tradução automática** (Google Translate, DeepL) ignoram sua terminologia, seu tom e as convenções de UI
- **A tradução genérica por LLM** funciona melhor, mas sem glossários e guias de estilo você obtém resultados inconsistentes

O Internationalizer é diferente. É um **pipeline de CLI** que combina a tradução por LLM com:

- **Glossários por idioma** — garantem terminologia consistente em todo o seu aplicativo
- **Guias de estilo por idioma** — controlam o tom, a formalidade, a pluralização e a tipografia
- **Memória de tradução** — ignora strings inalteradas, economizando dinheiro em chamadas de API
- **Validação determinística** — detecta chaves ausentes ou excedentes, desvios de estrutura protegida, problemas de glossário e erros de plural ou ICU antes do lançamento
<!-- internationalizer:unit markdown:installation -->
## Instalação

Instale via npm:

```bash
npm install -g internationalizer
```

Ou execute sem uma instalação global:

```bash
npx internationalizer --help
```

O pacote npm instala o binário pré-compilado correspondente do npm por meio de dependências opcionais específicas da plataforma.

Instale com Go:

```bash
go install github.com/Tom-R-Main/Internationalizer/cmd/internationalizer@latest
```

Ou compile a partir do código-fonte:

```bash
git clone https://github.com/Tom-R-Main/Internationalizer.git
cd Internationalizer
go build -o internationalizer ./cmd/internationalizer
```
<!-- internationalizer:unit markdown:npm-packages -->
## Pacotes npm

- As tags do Git e as versões do pacote npm devem corresponder; por exemplo, `v0.1.0` e `0.1.0`
- O pacote raiz `internationalizer` depende de pacotes de plataforma, como `internationalizer-darwin-arm64`
- Destinos npm compatíveis: macOS arm64/x64, Linux arm64/x64, Windows x64
- A publicação via CI requer um secret do GitHub com o nome `NPM_TOKEN`
<!-- internationalizer:unit markdown:quick-start -->
## Início rápido

1. Crie um arquivo de configuração na raiz do seu projeto:

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

2. Defina sua chave de API:

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

3. Faça uma prévia do que será traduzido:

```bash
internationalizer translate --dry-run
```

4. Execute a tradução:

```bash
internationalizer translate
```

5. Valide todos os locales:

```bash
internationalizer validate
```
<!-- internationalizer:unit markdown:commands -->
## Comandos

### `translate`

Encontre chaves ausentes ou desatualizadas e traduza-as por meio de um LLM.

```bash
internationalizer translate                    # translate all locales
internationalizer translate -l fr              # translate French only
internationalizer translate --dry-run          # preview without API calls
internationalizer translate --adopt-existing   # baseline existing translations without API calls
internationalizer translate --refresh-policy   # refresh prompt/style/model-stale entries
internationalizer translate --batch-size 20    # smaller batches
internationalizer translate --concurrency 2    # fewer parallel calls
```

O estado de tradução relata de forma independente as condições de ausente, com origem desatualizada (source-stale), com política desatualizada (policy-stale), atual e editado manualmente; assim, uma edição manual não consegue ocultar uma alteração na origem ou na política. Valores com política desatualizada são relatados, mas só são retraduzidos com `--refresh-policy`. Valores editados manualmente nunca são substituídos de forma automática. Use `--adopt-existing` ao introduzir o manifesto em traduções já revisadas ou ao aceitar explicitamente uma edição manual revisada como a nova linha de base.

### `validate`

Verifique todos os arquivos de locale em relação aos respectivos bundles de origem. A validação padrão verifica a cobertura estrutural (a porcentagem de chaves de destino necessárias presentes), relata chaves extras como avisos e falha em caso de chaves ausentes, incompatibilidades de interpolação ou estrutura de ICU MessageFormat inválida.

```bash
internationalizer validate                     # human-readable output
internationalizer validate --json              # machine-readable JSON
internationalizer validate -q                  # exit code only
internationalizer validate --strict             # enforce translation quality rules
internationalizer validate --require-state      # require current manifest provenance
```

O parâmetro `--strict` também relata a cobertura traduzida. Um valor linguístico idêntico à sua origem é considerado não traduzido, a menos que o glossário contenha explicitamente uma entrada idêntica de mesma origem e mesmo destino para o valor completo; `ignore_case` é respeitado, mas um termo de glossário inserido em um valor mais longo não constitui isenção. O modo estrito falha em caso de chaves extras, valores idênticos à origem, alterações na estrutura de interpolação/HTML/código/links de Markdown, violações de glossário e formas plurais configuradas.

O parâmetro `--require-state` verifica cada destino em relação ao arquivo `.internationalizer.lock`. Ele falha quando uma chave não é rastreada ou quando a origem gravada, a política de tradução ou o hash de destino estão desatualizados. Ele pode ser combinado com `--strict`.

Os relatórios legíveis por humanos e em JSON usam códigos de constatação estáveis:

| Código | Significado |
| --- | --- |
| `missing_key` / `extra_key` | Os conjuntos de chaves de origem e de destino diferem |
| `blank_translation` | Uma origem não vazia tem um destino vazio no modo estrito |
| `source_identical` | Um valor linguístico no modo estrito permanece não traduzido |
| `protected_structure_mismatch` | A estrutura de interpolação, HTML, código ou links foi alterada |
| `glossary_violation` | Nenhum termo ou variante de destino aprovado foi encontrado |
| `plural_form_missing` | Uma forma plural configurada para o locale está ausente |
| `icu_message_syntax` | Uma mensagem ICU de origem ou de destino está malformada |
| `icu_argument_mismatch` | Nomes, tipos ou estilos de formatador dos argumentos ICU diferem |
| `icu_selector_mismatch` | Os seletores diferem ou uma categoria de plural é inválida para o locale de destino |
| `untracked` | Não há registro no manifesto para o destino |
| `source_stale` | O conteúdo de origem foi alterado após a tradução gravada |
| `policy_stale` | As configurações do modelo ou o prompt gerado foram alterados |
| `target_modified` | O conteúdo de destino difere do registro no manifesto |

### `detect`

Detecte automaticamente o framework de i18n e sugira uma configuração.

```bash
internationalizer detect
```

Compatível com: react-i18next, next-intl, vue-i18n, JSON puro, documentos Markdown.

### `glossary`

Gerencie termos de glossário por idioma que são aplicados durante a tradução.

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

Gerencie a memória de tradução (cache JSONL de strings traduzidas anteriormente).

```bash
internationalizer tm stats                     # show record counts
internationalizer tm export                    # dump as JSON
internationalizer tm clear --force             # delete all records
```
<!-- internationalizer:unit markdown:configuration-reference -->
## Referência de configuração

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

Os identificadores de locale devem ser tags BCP 47 válidas, como `fr`, `pt-BR` ou `sr-Latn-RS`. Locales de destino canonicamente equivalentes são rejeitados como duplicatas, e as substituições de provedor específicas por locale correspondem à grafia canonicamente equivalente. No exemplo acima, os locales sem substituição — incluindo o japonês — herdam a configuração global do Gemini.

Valores em formato ICU MessageFormat são analisados estruturalmente. Há suporte para argumentos simples, `select`, `plural`, `selectordinal`, `number`, `date` e `time`, incluindo mensagens aninhadas, deslocamentos de plural (offsets), seletores de número exato e `#`. A validação verifica sintaxe, tipos de argumentos e estilos de formatadores, deslocamentos de plural, identidade de ramificações de select e categorias de plural do CLDR para o locale de destino. A resposta do provedor que violar essas invariantes é rejeitada antes que um arquivo de locale ou registro de memória de tradução seja gravado.

Com `i18next-v4`, famílias de plural reconhecidas na origem são expandidas durante a tradução para as categorias CLDR do locale de destino. Uma categoria exclusiva do destino usa o valor `_other` da família de origem como seu modelo de tradução. A validação estrita exige essas categorias de destino; categorias exclusivas da origem são opcionais para locales de destino que não as utilizam.
<!-- internationalizer:unit markdown:style-guides -->
## Guias de estilo

Guias de estilo são arquivos Markdown injetados no prompt de tradução do LLM. Eles controlam o tom, a formalidade, a tipografia e outras convenções específicas do idioma.

```
style-guides/
  _conventions.md    # shared rules for all languages
  fr.md              # French-specific rules
  ja.md              # Japanese-specific rules
  ar.md              # Arabic-specific rules
```

### Convenções compartilhadas (`_conventions.md`)

Defina regras que se aplicam a todos os idiomas: sintaxe de interpolação, preservação de HTML, convenções de tipos de strings (botões vs. rótulos vs. erros), etc.

### Guias por idioma (`{locale}.md`)

Defina regras específicas do idioma: registro de formalidade (tu vs. vous), pontuação (aspas angulares, pontos de interrogação invertidos), formas plurais, formatação de data/número e um glossário de terminologia.

Os guias de estilo são entradas de política permanentes, não saídas geradas. O Internationalizer lê esses guias, mas nunca os reescreve. O conteúdo deles é submetido a hash separadamente do glossário e do contrato de prompt, de modo que uma alteração no código do aplicativo não invalida uma tradução existente. Editar um guia marca intencionalmente esse locale para revisão de política; alterar o texto interno do prompt não faz isso, a menos que a versão do contrato de prompt também mude.

Consulte [`examples/react-app/style-guides/`](../../examples/react-app/style-guides/) para ver um exemplo funcional.
<!-- internationalizer:unit markdown:glossary-format -->
## Formato do glossário

Os arquivos de glossário são arrays JSON armazenados em `{glossary_dir}/{locale}.json`:

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

O campo `variants` lista outras formas de destino aprovadas. O campo `enforcement` pode ser `error`, `warning` ou omitido para o comportamento padrão de erro. Os termos são injetados no prompt do LLM como uma tabela de terminologia, garantindo uma tradução consistente em todo o seu aplicativo. Uma entrada exata como `{"source":"API","target":"API"}` também isenta esse valor idêntico à origem completo de constatações de valor não traduzido no modo estrito; ela não isenta um valor mais longo que apenas contenha `API`.
<!-- internationalizer:unit markdown:translation-memory -->
## Memória de tradução

A memória de tradução é armazenada como um arquivo JSONL (um registro JSON por linha). Cada registro contém:

- O bundle, a chave, o valor de origem, o valor traduzido e o locale canônico de destino
- Hashes da origem, do guia de estilo, do glossário, do contrato de prompt e da política combinada
- O provedor e o modelo que geraram a tradução
- Um carimbo de data/hora (timestamp)

Em execuções seguintes, strings com os mesmos hashes de origem e de política são recuperadas do cache sem chamar o LLM. O caminho padrão fica dentro do diretório ignorado `.internationalizer/`, portanto permanece como um cache local. Defina `tm_path` para um local rastreado se o seu projeto compartilha a memória de tradução intencionalmente. O manifesto `.internationalizer.lock` revisável é versionado separadamente.
<!-- internationalizer:unit markdown:supported-formats -->
## Formatos compatíveis

| Formato | Extensões | Modo |
| --- | --- | --- |
| JSON | `.json` | Chave-valor (aninhado, nivelado por notação de ponto) |
| YAML | `.yml`, `.yaml` | Chave-valor (preserva comentários e ordenação) |
| Markdown | `.md`, `.mdx` | Preâmbulo e seções de nível H2 |

Os arquivos Markdown de destino contêm comentários invisíveis `internationalizer:unit` antes das seções H2. Esses marcadores estáveis permitem que o Internationalizer adicione, mova ou edite uma seção de origem sem retraduzir seções não relacionadas. Documentos existentes sem marcação recebem marcadores na próxima atualização bem-sucedida.
<!-- internationalizer:unit markdown:project-type-detection -->
## Detecção de tipo de projeto

O comando `internationalizer detect` identifica sua configuração de i18n verificando:

- Dependências no `package.json` para react-i18next, next-intl ou vue-i18n
- Estruturas de diretórios correspondentes aos padrões comuns de locale
- Extensões de arquivo e convenções de nomenclatura
<!-- internationalizer:unit markdown:architecture -->
## Arquitetura

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
## Comparação com alternativas

| Recurso | Internationalizer | i18next | Crowdin | LLM genérico |
| --- | --- | --- | --- | --- |
| Tradução com tecnologia de LLM | Sim | Não | Parcial | Sim |
| Guias de estilo por idioma | Sim | Não | Não | Não |
| Aplicação de glossário | Sim | Não | Sim | Não |
| Memória de tradução | Sim | Não | Sim | Não |
| CLI / execução local | Sim | N/D | Não | Manual |
| Arquivos compatíveis com o Git | Sim | Sim | Parcial | Manual |
| Sem dependência de SaaS | Sim | Sim | Não | Varia |
| Código aberto (AGPL-3.0) | Sim | Sim | Não | Varia |
<!-- internationalizer:unit markdown:license -->
## Licença

[AGPL-3.0](../../LICENSE)

Consulte [THIRD_PARTY_NOTICES.md](../../THIRD_PARTY_NOTICES.md) para avisos sobre dependências.
<!-- internationalizer:unit markdown:contributing -->
## Contribuição

Consulte [CONTRIBUTING.md](../../CONTRIBUTING.md) para ver as diretrizes e a configuração de desenvolvimento. Todas as contribuições exigem sign-off do DCO.
