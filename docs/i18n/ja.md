> [English (original)](../../README.md)

<p align="center">
  <img src="../../assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

ソフトウェアプロジェクト向けのAIネイティブな国際化パイプラインです。LLMを使用してi18nファイルの翻訳、検証、管理を行います。

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](../../LICENSE)

<p align="center">
<a href="ar.md">العربية</a> · <a href="bn.md">বাংলা</a> · <a href="cs.md">Čeština</a> · <a href="da.md">Dansk</a> · <a href="de.md">Deutsch</a> · <a href="el.md">Ελληνικά</a> · <a href="es.md">Español</a> · <a href="fi.md">Suomi</a> · <a href="fr.md">Français</a> · <a href="he.md">עברית</a> · <a href="hi.md">हिन्दी</a> · <a href="id.md">Indonesia</a> · <a href="it.md">Italiano</a> · <a href="ja.md">日本語</a> · <a href="ko.md">한국어</a> · <a href="ms.md">Bahasa Melayu</a><br><a href="nl.md">Nederlands</a> · <a href="pa.md">ਪੰਜਾਬੀ</a> · <a href="pl.md">Polski</a> · <a href="pt-BR.md">Português</a> · <a href="ro.md">Română</a> · <a href="ru.md">Русский</a> · <a href="sv.md">Svenska</a> · <a href="te.md">తెలుగు</a> · <a href="th.md">ไทย</a> · <a href="tr.md">Türkçe</a> · <a href="uk.md">Українська</a> · <a href="vi.md">Tiếng Việt</a> · <a href="yue.md">粵語</a> · <a href="zh-CN.md">简体中文</a> · <a href="zh-TW.md">繁體中文</a>
</p>

---
<!-- internationalizer:unit markdown:why-internationalizer -->
## Internationalizerを選ぶ理由

多くのi18nツールは、ランタイムライブラリ（i18next、react-intl）か、キー管理SaaSプラットフォーム（Crowdin、Lokalise）のいずれかです。しかし、実際の翻訳課題を適切に解決できているツールはありません。

- **手動翻訳**は、数言語を超えるとスケールしません
- **機械翻訳API**（Google翻訳、DeepL）は、独自の用語集、トーン、UIの慣例を無視します
- **汎用LLM翻訳**はより高精度ですが、用語集やスタイルガイドがなければ出力結果にばらつきが生じます

Internationalizerは異なります。LLM翻訳と以下の機能を組み合わせた**CLIパイプライン**です。

- **言語ごとの用語集** — アプリケーション全体で一貫した用語を適用
- **言語ごとのスタイルガイド** — トーン、丁寧さの度合い、複数形、タイポグラフィを制御
- **翻訳メモリ** — 変更のない文字列をスキップし、API呼び出しコストを削減
- **決定論的な検証** — キーの欠落や余分なキー、保護対象構造の乖離、用語集違反、複数形やICUのエラーをリリース前に検出
<!-- internationalizer:unit markdown:installation -->
## インストール

npmからインストール:

```bash
npm install -g internationalizer
```

またはグローバルインストールなしで実行:

```bash
npx internationalizer --help
```

npmパッケージは、プラットフォーム固有のオプションの依存関係を介して、一致するビルド済みバイナリをnpmからインストールします。

Goでインストール:

```bash
go install github.com/Tom-R-Main/Internationalizer/cmd/internationalizer@latest
```

またはソースからビルド:

```bash
git clone https://github.com/Tom-R-Main/Internationalizer.git
cd Internationalizer
go build -o internationalizer ./cmd/internationalizer
```
<!-- internationalizer:unit markdown:npm-packages -->
## npmパッケージ

- Gitタグとnpmパッケージのバージョンは一致している必要があります（例: `v0.1.0`と`0.1.0`）
- ルートの`internationalizer`パッケージは、`internationalizer-darwin-arm64`などのプラットフォーム別パッケージに依存します
- サポート対象のnpmターゲット: macOS arm64/x64、Linux arm64/x64、Windows x64
- CIでのパッケージ公開には、`NPM_TOKEN`という名前のGitHubシークレットが必要です
<!-- internationalizer:unit markdown:quick-start -->
## クイックスタート

1. プロジェクトのルートに設定ファイルを作成します:

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

2. APIキーを設定します:

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

3. 翻訳対象をプレビューします:

```bash
internationalizer translate --dry-run
```

4. 翻訳を実行します:

```bash
internationalizer translate
```

5. すべてのロケールを検証します:

```bash
internationalizer validate
```
<!-- internationalizer:unit markdown:commands -->
## コマンド

### `translate`

欠落しているキーや古くなったキーを検出し、LLM経由で翻訳します。

```bash
internationalizer translate                    # translate all locales
internationalizer translate -l fr              # translate French only
internationalizer translate --dry-run          # preview without API calls
internationalizer translate --adopt-existing   # baseline existing translations without API calls
internationalizer translate --refresh-policy   # refresh prompt/style/model-stale entries
internationalizer translate --batch-size 20    # smaller batches
internationalizer translate --concurrency 2    # fewer parallel calls
```

翻訳状態は、欠落（missing）、ソース変更（source-stale）、ポリシー変更（policy-stale）、最新（current）、手動編集済み（manually edited）の状態を個別にレポートするため、手動編集によってソースやポリシーの変更が隠蔽されることはありません。ポリシー変更によって古くなった値はレポートされますが、`--refresh-policy`を指定した場合にのみ再翻訳されます。手動編集された値が自動的に上書きされることは決してありません。レビュー済みの翻訳にマニフェストを導入する場合や、レビュー済みの手動編集を新しいベースラインとして明示的に受け入れる場合は、`--adopt-existing`を使用してください。

### `validate`

すべてのロケールファイルをソースバンドルと照合します。デフォルトの検証では構造的な充足率（必要なターゲットキーが存在する割合）を検査し、余分なキーを警告として報告します。キーの欠落、補間パラメータの不一致、無効なICU MessageFormat構造がある場合は失敗します。

```bash
internationalizer validate                     # human-readable output
internationalizer validate --json              # machine-readable JSON
internationalizer validate -q                  # exit code only
internationalizer validate --strict             # enforce translation quality rules
internationalizer validate --require-state      # require current manifest provenance
```

`--strict`は翻訳済みの充足率も報告します。言語表現としての値がソースと同一の場合、値全体に対してソースとターゲットが完全に一致する項目が用語集に明示的に登録されていない限り、未翻訳とみなされます。`ignore_case`は考慮されますが、より長い値の一部を用語集の語が占めているだけでは除外対象になりません。strictモードでは、余分なキー、ソースと同一の値、補間・HTML・コード・Markdownリンクの構造変更、用語集違反、設定された複数形の不備がある場合に失敗します。

`--require-state`は、各ターゲットを`.internationalizer.lock`と照合します。キーが未追跡の場合や、記録されたソース、翻訳ポリシー、ターゲットのハッシュが古い場合に失敗します。`--strict`と併用できます。

人間向けおよびJSONレポートでは、安定した検出コードを使用します:

| コード | 意味 |
| --- | --- |
| `missing_key` / `extra_key` | ソースとターゲットのキー集合が異なる |
| `blank_translation` | 空でないソースに対してstrictモードのターゲットが空 |
| `source_identical` | strictモードの言語表現の値が未翻訳のまま |
| `protected_structure_mismatch` | 補間、HTML、コード、リンク構造が変更されている |
| `glossary_violation` | 承認済みのターゲット用語または異表記が見つからない |
| `plural_form_missing` | 設定されたロケールの複数形が存在しない |
| `icu_message_syntax` | ソースまたはターゲットのICUメッセージが不正 |
| `icu_argument_mismatch` | ICUの引数名、型、またはフォーマッタースタイルが異なる |
| `icu_selector_mismatch` | セレクターが異なる、または複数形カテゴリーがターゲットロケールで無効 |
| `untracked` | ターゲットに対応するマニフェストレコードが存在しない |
| `source_stale` | 記録された翻訳の後にソースの内容が変更された |
| `policy_stale` | 生成プロンプトまたはモデル設定が変更された |
| `target_modified` | ターゲットの内容がマニフェストの記録と異なる |

### `detect`

i18nフレームワークを自動検出し、設定を提案します。

```bash
internationalizer detect
```

サポート対象: react-i18next、next-intl、vue-i18n、プレーンなJSON、Markdownドキュメント。

### `glossary`

翻訳時に適用される言語ごとの用語集を管理します。

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

翻訳メモリ（過去に翻訳された文字列のJSONLキャッシュ）を管理します。

```bash
internationalizer tm stats                     # show record counts
internationalizer tm export                    # dump as JSON
internationalizer tm clear --force             # delete all records
```
<!-- internationalizer:unit markdown:configuration-reference -->
## 設定リファレンス

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

ロケール識別子には、`fr`、`pt-BR`、`sr-Latn-RS`などの適切な形式のBCP 47タグを使用する必要があります。正規同等なターゲットロケールは重複として拒否され、ロケール固有のプロバイダーオーバーライドも正規同等な表記に一致します。上の例では、オーバーライドのないロケール（日本語を含む）はグローバルのGemini設定を継承します。

ICU MessageFormatの値は構造的に解析されます。単純な引数、`select`、`plural`、`selectordinal`、`number`、`date`、`time`がサポートされており、ネストされたメッセージ、複数形のオフセット、厳密な数値セレクター、`#`が含まれます。検証では、構文、引数の型とフォーマッタースタイル、複数形のオフセット、select分岐の同一性、ターゲットロケールのCLDR複数形カテゴリーをチェックします。これらの不変条件に違反するプロバイダー出力は、ロケールファイルや翻訳メモリレコードに書き込まれる前に拒否されます。

`i18next-v4`を使用すると、認識されたソースの複数形ファミリーが翻訳中にターゲットロケールのCLDRカテゴリーに展開されます。ターゲットにのみ存在するカテゴリーは、ソースファミリーの`_other`値を翻訳テンプレートとして使用します。strict検証ではこれらのターゲットカテゴリーが必須となります。ターゲットロケールで使用されないソース固有のカテゴリーは任意です。
<!-- internationalizer:unit markdown:style-guides -->
## スタイルガイド

スタイルガイドは、LLM翻訳プロンプトに挿入されるMarkdownファイルです。トーン、丁寧さの度合い、タイポグラフィ、その他の言語固有の慣例を制御します。

```
style-guides/
  _conventions.md    # shared rules for all languages
  fr.md              # French-specific rules
  ja.md              # Japanese-specific rules
  ar.md              # Arabic-specific rules
```

### 共通規則（`_conventions.md`）

すべての言語に適用されるルールを定義します: 補間構文、HTMLの保持、文字列種別の規則（ボタン、ラベル、エラーなど）。

### 言語別ガイド（`{locale}.md`）

言語固有のルールを定義します: 敬体・常体の使い分け（tuとvousなど）、句読点（ギュメ、逆疑問符）、複数形、日付・数値の書式、用語集。

スタイルガイドは永続的なポリシー入力であり、生成された出力ではありません。Internationalizerはこれらを読み取りますが、書き換えることはありません。その内容は用語集やプロンプト規約とは別にハッシュ化されるため、アプリケーションコードの変更によって翻訳が古くなることはありません。ガイドを編集すると、そのロケールは意図的にポリシー再確認の対象としてマークされます。内部プロンプトの文言を変更しても、プロンプト規約のバージョンが変わらない限り、古くなったとはみなされません。

動作例については、[`examples/react-app/style-guides/`](../../examples/react-app/style-guides/)を参照してください。
<!-- internationalizer:unit markdown:glossary-format -->
## 用語集の形式

用語集ファイルは、`{glossary_dir}/{locale}.json`に格納されるJSON配列です:

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

`variants`には、承認されたその他のターゲット表現をリストします。`enforcement`には`error`、`warning`を指定でき、省略した場合はデフォルトのerror動作になります。用語は用語集テーブルとしてLLMプロンプトに挿入され、アプリケーション全体での一貫した翻訳を保証します。`{"source":"API","target":"API"}`のような完全一致エントリは、その完全なソース同一値をstrictモードの未翻訳検出から除外します。単に`API`を含むだけのより長い値は除外されません。
<!-- internationalizer:unit markdown:translation-memory -->
## 翻訳メモリ

翻訳メモリはJSONLファイル（1行に1つのJSONレコード）として保存されます。各レコードには以下が含まれます:

- バンドル、キー、ソース値、翻訳値、正規化されたターゲットロケール
- ソース、スタイルガイド、用語集、プロンプト規約、および統合ポリシーの各ハッシュ
- 翻訳を生成したプロバイダーおよびモデル
- タイムスタンプ

以降の実行では、ソースとポリシーのハッシュが同じ文字列は、LLMを呼び出さずにキャッシュから提供されます。デフォルトのパスは無視対象の`.internationalizer/`ディレクトリ配下にあるため、ローカルキャッシュのまま保持されます。プロジェクトで意図的に翻訳メモリを共有する場合は、`tm_path`を追跡対象の場所に変更してください。レビュー可能な`.internationalizer.lock`マニフェストは個別にバージョン管理されます。
<!-- internationalizer:unit markdown:supported-formats -->
## サポート対象の形式

| 形式 | 拡張子 | モード |
| --- | --- | --- |
| JSON | `.json` | キーバリュー（ネスト対応、ドット記法による平坦化） |
| YAML | `.yml`, `.yaml` | キーバリュー（コメントと順序を保持） |
| Markdown | `.md`, `.mdx` | プリアンブルおよびH2レベルのセクション |

Markdownターゲットには、H2セクションの前に不可視の`internationalizer:unit`コメントが含まれます。これらの安定したマーカーにより、Internationalizerは関係のないセクションを再翻訳することなく、1つのソースセクションを追加、移動、編集できます。マーカーのない既存のドキュメントには、次回の正常な更新時にマーカーが付与されます。
<!-- internationalizer:unit markdown:project-type-detection -->
## プロジェクトタイプの検出

`internationalizer detect`は、以下をチェックしてi18n設定を特定します:

- `package.json`の依存関係におけるreact-i18next、next-intl、vue-i18n
- 一般的なロケールパターンに一致するディレクトリ構造
- ファイル拡張子および命名規則
<!-- internationalizer:unit markdown:architecture -->
## アーキテクチャ

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
## 代替ツールとの比較

| 機能 | Internationalizer | i18next | Crowdin | 汎用LLM |
| --- | --- | --- | --- | --- |
| LLMを活用した翻訳 | はい | いいえ | 一部対応 | はい |
| 言語ごとのスタイルガイド | はい | いいえ | いいえ | いいえ |
| 用語集の強制適用 | はい | いいえ | はい | いいえ |
| 翻訳メモリ | はい | いいえ | はい | いいえ |
| CLI / ローカル実行 | はい | 該当なし | いいえ | 手動 |
| Gitに適したファイル | はい | はい | 一部対応 | 手動 |
| SaaSへの依存なし | はい | はい | いいえ | 状況による |
| オープンソース（AGPL-3.0） | はい | はい | いいえ | 状況による |
<!-- internationalizer:unit markdown:license -->
## ライセンス

[AGPL-3.0](../../LICENSE)

依存関係の通知については、[THIRD_PARTY_NOTICES.md](../../THIRD_PARTY_NOTICES.md)を参照してください。
<!-- internationalizer:unit markdown:contributing -->
## コントリビューション

開発のセットアップとガイドラインについては、[CONTRIBUTING.md](../../CONTRIBUTING.md)を参照してください。すべてのコントリビューションにはDCOのサインオフが必要です。
