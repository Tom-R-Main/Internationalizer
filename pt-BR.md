<p align="center">
  <img src="assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

Pipeline de internacionalização nativo de IA para projetos de software. Traduza, valide e gerencie arquivos i18n usando LLMs.

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](LICENSE)

## Por que o Internationalizer?

A maioria das ferramentas de i18n são bibliotecas de tempo de execução (i18next, react-intl) ou plataformas SaaS de gerenciamento de chaves (Crowdin, Lokalise). Nenhuma delas resolve bem o problema real de tradução:

- **Tradução manual** não escala além de alguns idiomas
- **APIs de tradução automática** (Google Translate, DeepL) ignoram sua terminologia, tom e convenções de interface de usuário (UI)
- **Tradução genérica por LLM** funciona melhor, mas sem glossários e guias de estilo, você obtém resultados inconsistentes

O Internationalizer é diferente. É um **pipeline de CLI** que combina tradução por LLM com:

- **Glossários por idioma** — impõe terminologia consistente em todo o seu aplicativo
- **Guias de estilo por idioma** — controla o tom, a formalidade, a pluralização e a tipografia
- **Memória de tradução** — ignora strings inalteradas, economizando dinheiro em chamadas de API
- **Validação de chaves** — detecta traduções ausentes e incompatibilidades de interpolação antes do lançamento

## Instalação

Instale a partir do npm:

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

## Pacotes npm

- As tags do Git e as versões do pacote npm devem corresponder, por exemplo, `v0.1.0` e `0.1.0`
- O pacote raiz `internationalizer` depende de pacotes de plataforma, como `internationalizer-darwin-arm64`
- Alvos npm suportados: macOS arm64/x64, Linux arm64/x64, Windows x64
- A publicação via CI requer um secret do GitHub chamado `NPM_TOKEN`

## Início Rápido

1. Crie um arquivo de configuração na raiz do seu projeto:

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

2. Defina sua chave de API:

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

3. Visualize o que será traduzido:

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

## Comandos

### `translate`

Encontre chaves ausentes e traduza-as via um LLM.

```bash
internationalizer translate                    # traduzir todos os locales
internationalizer translate -l fr              # traduzir apenas para francês
internationalizer translate --dry-run          # visualizar sem chamadas de API
internationalizer translate --batch-size 20    # lotes menores
internationalizer translate --concurrency 2    # menos chamadas paralelas
```

### `validate`

Verifique todos os arquivos de locale em busca de chaves ausentes, chaves extras e incompatibilidades de interpolação.

```bash
internationalizer validate                     # saída legível por humanos
internationalizer validate --json              # JSON legível por máquina
internationalizer validate -q                  # apenas código de saída
```

### `detect`

Detecte automaticamente o framework i18n e sugira uma configuração.

```bash
internationalizer detect
```

Suporta: react-i18next, next-intl, vue-i18n, JSON puro, documentos markdown.

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
internationalizer tm stats                     # mostrar contagem de registros
internationalizer tm export                    # exportar como JSON
internationalizer tm clear --force             # excluir todos os registros
```

## Referência de Configuração

```yaml
# .internationalizer.yml

# Idioma de origem (padrão: en)
source_locale: en

# Idiomas para os quais traduzir (obrigatório)
target_locales: [fr, de, es, ja, zh-CN, ar]

# Caminho para o arquivo de locale de origem (obrigatório)
source_path: locales/en.json

# Configurações do provedor de LLM
llm:
  # Provedor: "anthropic", "openai", "gemini" ou "openrouter" (padrão: gemini)
  provider: gemini

  # Padrões de nome de modelo por provedor:
  #   anthropic:  claude-sonnet-4-6
  #   openai:     gpt-5.4
  #   gemini:     gemini-3.1-pro-preview
  #   openrouter: google/gemini-3-flash-preview
  model: gemini-3.1-pro-preview

  # Variável de ambiente contendo a chave de API
  api_key_env: GOOGLE_AI_STUDIO_API_KEY

  # URL base para endpoints compatíveis com OpenAI (opcional)
  # base_url: https://api.openai.com

# Chaves por chamada de LLM (padrão: 40)
batch_size: 40

# Chamadas paralelas de LLM (padrão: 4)
concurrency: 4

# Diretório contendo arquivos Markdown de guia de estilo por locale (padrão: style-guides)
style_guides_dir: style-guides

# Diretório contendo arquivos JSON de glossário por locale (padrão: glossary)
glossary_dir: glossary

# Caminho para o arquivo de memória de tradução (padrão: .internationalizer/tm.jsonl)
tm_path: .internationalizer/tm.jsonl
```

## Guias de Estilo

Os guias de estilo são arquivos Markdown que são injetados no prompt de tradução do LLM. Eles controlam o tom, a formalidade, a tipografia e outras convenções específicas do idioma.

```
style-guides/
  _conventions.md    # regras compartilhadas para todos os idiomas
  fr.md              # regras específicas para francês
  ja.md              # regras específicas para japonês
  ar.md              # regras específicas para árabe
```

### Convenções compartilhadas (`_conventions.md`)

Defina regras que se aplicam a todos os idiomas: sintaxe de interpolação, preservação de HTML, convenções de tipo de string (botões vs. rótulos vs. erros), etc.

### Guias por idioma (`{locale}.md`)

Defina regras específicas do idioma: registro de formalidade (tu vs. vous), pontuação (aspas, pontos de interrogação invertidos), formas plurais, formatação de data/número e um glossário de terminologia.

Consulte [`examples/react-app/style-guides/`](examples/react-app/style-guides/) para ver um exemplo prático.

## Formato do Glossário

Os arquivos de glossário são arrays JSON armazenados em `{glossary_dir}/{locale}.json`:

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

Os termos são injetados no prompt do LLM como uma tabela de terminologia, garantindo a tradução consistente de termos-chave em todo o seu aplicativo.

## Memória de Tradução

A memória de tradução é armazenada como um arquivo JSONL (um registro JSON por linha). Cada registro contém:

- A chave e o valor de origem
- O valor traduzido
- Um hash SHA-256 do valor de origem
- Um carimbo de data/hora (timestamp)

Em execuções subsequentes, as strings inalteradas são servidas a partir do cache da TM (Memória de Tradução) sem chamar o LLM, economizando tempo e custos de API. O arquivo da TM é amigável ao git e pode ser comitado junto com seus arquivos de locale.

## Formatos Suportados

| Formato | Extensões | Modo |
|--------|-----------|------|
| JSON | `.json` | Chave-valor (aninhado, achatado por notação de ponto) |
| YAML | `.yml`, `.yaml` | Chave-valor (preserva comentários e ordenação) |
| Markdown | `.md`, `.mdx` | Tradução de documento inteiro |

## Detecção de Tipo de Projeto

O `internationalizer detect` identifica sua configuração de i18n verificando:

- Dependências no `package.json` para react-i18next, next-intl ou vue-i18n
- Estruturas de diretório que correspondem a padrões comuns de locale
- Extensões de arquivo e convenções de nomenclatura

## Arquitetura

```
cmd/internationalizer/     Ponto de entrada da CLI e definições de comando
internal/
  config/                  Carregamento de configuração YAML com padrões
  detect/                  Detecção automática de tipo de projeto
  formats/                 Analisadores de formato (JSON, YAML, Markdown)
  glossary/                Gerenciamento de glossário por locale
  llm/                     Interface do provedor de LLM + implementações
    anthropic.go           Backend do Anthropic Claude
    openai.go              Backend da OpenAI / compatível
    gemini.go              Backend do Google Gemini via AI Studio
                           OpenRouter usa openai.go com base_url personalizada
  styleguide/              Carregador de guia de estilo
  tm/                      Memória de tradução JSONL
  translate/               Orquestrador de tradução
  validate/                Validação de locale e diffing
```

## Comparação com Alternativas

| Recurso | Internationalizer | i18next | Crowdin | LLM Genérico |
|---------|------------------|---------|---------|-------------|
| Tradução baseada em LLM | Sim | Não | Parcial | Sim |
| Guias de estilo por idioma | Sim | Não | Não | Não |
| Aplicação de glossário | Sim | Não | Sim | Não |
| Memória de tradução | Sim | Não | Sim | Não |
| CLI / execução local | Sim | N/A | Não | Manual |
| Arquivos amigáveis ao Git | Sim | Sim | Parcial | Manual |
| Sem dependência de SaaS | Sim | Sim | Não | Varia |
| Código aberto (AGPL-3.0) | Sim | Sim | Não | Varia |

## Licença

[AGPL-3.0](LICENSE)

## Contribuição

Consulte [CONTRIBUTING.md](CONTRIBUTING.md) para obter as diretrizes e a configuração de desenvolvimento. Todas as contribuições exigem a assinatura do DCO.

