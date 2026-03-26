> [English (original)](../../README.md)

<p align="center">
  <img src="../../assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

ソフトウェアプロジェクト向けのAIネイティブな国際化パイプラインです。LLMを使用してi18nファイルの翻訳、検証、管理を行います。

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](LICENSE)

## なぜ Internationalizer なのか？

ほとんどのi18nツールは、ランタイムライブラリ（i18next、react-intl）か、キー管理SaaSプラットフォーム（Crowdin、Lokalise）のいずれかです。しかし、どれも実際の翻訳問題をうまく解決できていません。

- **手動翻訳**は、少数の言語を超えるとスケールしません。
- **機械翻訳API**（Google Translate、DeepL）は、専門用語、トーン、UIの規則を無視します。
- **一般的なLLM翻訳**はより良く機能しますが、用語集やスタイルガイドがないと、一貫性のない結果になります。

Internationalizerは違います。LLM翻訳と以下を組み合わせた**CLIパイプライン**です。

- **言語ごとの用語集** — アプリ全体で一貫した専門用語を強制します。
- **言語ごとのスタイルガイド** — トーン、フォーマルさ、複数形、タイポグラフィを制御します。
- **翻訳メモリ** — 変更されていない文字列をスキップし、API呼び出しのコストを節約します。
- **キーの検証** — リリース前に、翻訳の欠落や補間の不一致を検出します。

## インストール

npmからのインストール:

```bash
npm install -g internationalizer
```

または、グローバルインストールせずに実行する場合:

```bash
npx internationalizer --help
```

npmパッケージは、プラットフォーム固有のオプションの依存関係を介して、npmから一致するビルド済みバイナリをインストールします。

Goでのインストール:

```bash
go install github.com/Tom-R-Main/Internationalizer/cmd/internationalizer@latest
```

または、ソースからのビルド:

```bash
git clone https://github.com/Tom-R-Main/Internationalizer.git
cd Internationalizer
go build -o internationalizer ./cmd/internationalizer
```

## npmパッケージ

- Gitタグとnpmパッケージのバージョンは一致している必要があります（例: `v0.1.0` と `0.1.0`）。
- ルートの `internationalizer` パッケージは、`internationalizer-darwin-arm64` などのプラットフォームパッケージに依存しています。
- サポートされているnpmターゲット: macOS arm64/x64、Linux arm64/x64、Windows x64
- CIでの公開には、`NPM_TOKEN` という名前のGitHubシークレットが必要です。

## クイックスタート

1. プロジェクトのルートに設定ファイルを作成します:

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

2. APIキーを設定します:

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

3. 翻訳される内容をプレビューします:

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

## コマンド

### `translate`

欠落しているキーを見つけ、LLMを介して翻訳します。

```bash
internationalizer translate                    # すべてのロケールを翻訳
internationalizer translate -l fr              # フランス語のみ翻訳
internationalizer translate --dry-run          # APIを呼び出さずにプレビュー
internationalizer translate --batch-size 20    # バッチサイズを小さくする
internationalizer translate --concurrency 2    # 並行呼び出しを減らす
```

### `validate`

すべてのロケールファイルをチェックし、キーの欠落、余分なキー、補間の不一致がないか確認します。

```bash
internationalizer validate                     # 人間が読める形式で出力
internationalizer validate --json              # 機械可読なJSON形式で出力
internationalizer validate -q                  # 終了コードのみ出力
```

### `detect`

i18nフレームワークを自動検出し、設定を提案します。

```bash
internationalizer detect
```

サポート対象: react-i18next、next-intl、vue-i18n、プレーンなJSON、Markdownドキュメント。

### `glossary`

翻訳時に強制される言語ごとの用語集を管理します。

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

翻訳メモリ（過去に翻訳された文字列のJSONLキャッシュ）を管理します。

```bash
internationalizer tm stats                     # レコード数を表示
internationalizer tm export                    # JSONとしてダンプ
internationalizer tm clear --force             # すべてのレコードを削除
```

## 設定リファレンス

```yaml
# .internationalizer.yml

# ソース言語（デフォルト: en）
source_locale: en

# 翻訳先の言語（必須）
target_locales: [fr, de, es, ja, zh-CN, ar]

# ソースロケールファイルへのパス（必須）
source_path: locales/en.json

# LLMプロバイダーの設定
llm:
  # プロバイダー: "anthropic"、"openai"、"gemini"、または "openrouter"（デフォルト: gemini）
  provider: gemini

  # プロバイダーごとのデフォルトモデル名:
  #   anthropic:  claude-sonnet-4-6
  #   openai:     gpt-5.4
  #   gemini:     gemini-3.1-pro-preview
  #   openrouter: google/gemini-3-flash-preview
  model: gemini-3.1-pro-preview

  # APIキーを含む環境変数
  api_key_env: GOOGLE_AI_STUDIO_API_KEY

  # OpenAI互換エンドポイントのベースURL（オプション）
  # base_url: https://api.openai.com

# LLM呼び出しあたりのキー数（デフォルト: 40）
batch_size: 40

# 並行LLM呼び出し数（デフォルト: 4）
concurrency: 4

# 言語ごとのスタイルガイドMarkdownファイルを含むディレクトリ（デフォルト: style-guides）
style_guides_dir: style-guides

# 言語ごとの用語集JSONファイルを含むディレクトリ（デフォルト: glossary）
glossary_dir: glossary

# 翻訳メモリファイルへのパス（デフォルト: .internationalizer/tm.jsonl）
tm_path: .internationalizer/tm.jsonl
```

## スタイルガイド

スタイルガイドは、LLMの翻訳プロンプトに注入されるMarkdownファイルです。トーン、フォーマルさ、タイポグラフィ、その他の言語固有の規則を制御します。

```
style-guides/
  _conventions.md    # すべての言語に共通のルール
  fr.md              # フランス語固有のルール
  ja.md              # 日本語固有のルール
  ar.md              # アラビア語固有のルール
```

### 共通の規則 (`_conventions.md`)

すべての言語に適用されるルールを定義します。補間構文、HTMLの保持、文字列タイプの規則（ボタン、ラベル、エラーなど）が含まれます。

### 言語ごとのガイド (`{locale}.md`)

言語固有のルールを定義します。フォーマルさのレベル（tu と vous など）、句読点（ギュメ、逆疑問符など）、複数形、日付/数値のフォーマット、専門用語の用語集が含まれます。

実際の例については、[`examples/react-app/style-guides/`](examples/react-app/style-guides/) を参照してください。

## 用語集のフォーマット

用語集ファイルは、`{glossary_dir}/{locale}.json` に保存されるJSON配列です:

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

用語は用語表としてLLMプロンプトに注入され、アプリケーション全体で重要な用語が一貫して翻訳されるようにします。

## 翻訳メモリ

翻訳メモリはJSONLファイル（1行に1つのJSONレコード）として保存されます。各レコードには以下が含まれます:

- ソースのキーと値
- 翻訳された値
- ソース値のSHA-256ハッシュ
- タイムスタンプ

次回以降の実行では、変更されていない文字列はLLMを呼び出すことなくTMキャッシュから提供されるため、時間とAPIコストの両方を節約できます。TMファイルはGitと相性が良く、ロケールファイルと一緒にコミットできます。

## サポートされているフォーマット

| フォーマット | 拡張子 | モード |
|--------|-----------|------|
| JSON | `.json` | キーバリュー（ネスト、ドット記法によるフラット化） |
| YAML | `.yml`, `.yaml` | キーバリュー（コメントと順序を保持） |
| Markdown | `.md`, `.mdx` | ドキュメント全体の翻訳 |

## プロジェクトタイプの検出

`internationalizer detect` は以下をチェックしてi18nの設定を特定します:

- `package.json` の依存関係（react-i18next、next-intl、vue-i18n）
- 一般的なロケールパターンに一致するディレクトリ構造
- ファイルの拡張子と命名規則

## アーキテクチャ

```
cmd/internationalizer/     CLIエントリポイントとコマンド定義
internal/
  config/                  デフォルト値付きのYAML設定の読み込み
  detect/                  プロジェクトタイプの自動検出
  formats/                 フォーマットパーサー（JSON、YAML、Markdown）
  glossary/                言語ごとの用語集管理
  llm/                     LLMプロバイダーのインターフェースと実装
    anthropic.go           Anthropic Claudeバックエンド
    openai.go              OpenAI / 互換バックエンド
    gemini.go              AI Studio経由のGoogle Geminiバックエンド
                           OpenRouterはカスタムbase_urlでopenai.goを使用
  styleguide/              スタイルガイドローダー
  tm/                      JSONL翻訳メモリ
  translate/               翻訳オーケストレーター
  validate/                ロケールの検証と差分確認
```

## 代替ツールとの比較

| 機能 | Internationalizer | i18next | Crowdin | 一般的なLLM |
|---------|------------------|---------|---------|-------------|
| LLMによる翻訳 | はい | いいえ | 一部 | はい |
| 言語ごとのスタイルガイド | はい | いいえ | いいえ | いいえ |
| 用語集の強制 | はい | いいえ | はい | いいえ |
| 翻訳メモリ | はい | いいえ | はい | いいえ |
| CLI / ローカル実行 | はい | N/A | いいえ | 手動 |
| Gitと相性の良いファイル | はい | はい | 一部 | 手動 |
| SaaS依存なし | はい | はい | いいえ | ツールによる |
| オープンソース (AGPL-3.0) | はい | はい | いいえ | ツールによる |

## ライセンス

[AGPL-3.0](LICENSE)

## コントリビューション

開発のセットアップとガイドラインについては、[CONTRIBUTING.md](CONTRIBUTING.md) を参照してください。すべてのコントリビューションにはDCOの署名が必要です。

