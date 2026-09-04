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

## Internationalizerを選ぶ理由

一般的なi18nツールの多くは、ランタイムライブラリ（i18next、react-intl）か、キー管理SaaSプラットフォーム（Crowdin、Lokalise）のいずれかです。しかし、実際の翻訳課題を適切に解決できているツールはありません。

- **手動翻訳**は、数言語を超えると対応しきれなくなります
- **機械翻訳API**（Google翻訳、DeepL）は、専門用語、トーン、UIの規則を無視します
- **汎用LLM翻訳**は精度が高いものの、用語集やスタイルガイドがなければ出力結果にばらつきが生じます

Internationalizerは異なります。LLM翻訳と以下の機能を組み合わせた**CLIパイプライン**です。

- **言語ごとの用語集** — アプリケーション全体で一貫した用語を適用
- **言語ごとのスタイルガイド** — トーン、丁寧さの度合い、複数形、タイポグラフィを制御
- **翻訳メモリ** — 変更のない文字列をスキップし、API呼び出しコストを削減
- **決定論的な検証** — リリース前にキーの欠落や余分なキー、保護対象の構造の差異、用語集違反、複数形やICUのエラーを検出

## インストール

npmでのインストール

```bash
npm install -g internationalizer
```

グローバルインストールせずに実行

```bash
npx internationalizer --help
```

npmパッケージは、プラットフォーム固有のオプション依存関係を介して、対応するビルド済みバイナリをnpmから自動的にインストールします。

Goでのインストール

```bash
go install github.com/Tom-R-Main/Internationalizer/cmd/internationalizer@latest
```

ソースからのビルド

```bash
git clone https://github.com/Tom-R-Main/Internationalizer.git
cd Internationalizer
go build -o internationalizer ./cmd/internationalizer
```

## npmパッケージ

- Gitタグとnpmパッケージのバージョンは一致している必要があります（例：`v0.1.0`と`0.1.0`）
- ルートの`internationalizer`パッケージは、`internationalizer-darwin-arm64`などのプラットフォーム別パッケージに依存します
- サポート対象のnpmターゲット：macOS arm64/x64、Linux arm64/x64、Windows x64
- CIでのパッケージ公開には、`NPM_TOKEN`という名前のGitHubシークレットが必要です

## クイックスタート

1. プロジェクトのルートに設定ファイルを作成します。

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

2. APIキーを設定します。

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

3. 翻訳対象をプレビューします。

```bash
internationalizer translate --dry-run
```

4. 翻訳を実行します。

```bash
internationalizer translate
```

5. すべてのロケールを検証します。

```bash
internationalizer validate
```

## コマンド

### `translate`

欠落しているキーや古くなったキーを検出し、LLM経由で翻訳します。

```bash
internationalizer translate                    # すべてのロケールを翻訳
internationalizer translate -l fr              # フランス語のみ翻訳
internationalizer translate --dry-run          # APIを呼び出さずにプレビュー
internationalizer translate --adopt-existing   # APIを呼び出さずに既存の翻訳をベースライン化
internationalizer translate --refresh-policy   # プロンプト/スタイル/モデルの変更により古くなったエントリを更新
internationalizer translate --batch-size 20    # より小さいバッチサイズで実行
internationalizer translate --concurrency 2    # 並行呼び出し数を抑えて実行
```

翻訳状態は、欠落（missing）、ソース変更（source-stale）、ポリシー変更（policy-stale）、最新（current）、手動編集済み（manually edited）の状態を個別にレポートするため、手動編集によってソースやポリシーの変更が覆い隠されることはありません。ポリシー変更によって古くなった値はレポートされますが、`--refresh-policy`を指定した場合にのみ再翻訳されます。手動編集された値が自動的に上書きされることは決してありません。レビュー済みの翻訳にマニフェストを導入する場合や、レビュー済みの手動編集を新しいベースラインとして明示的に受け入れる場合は、`--adopt-existing`を使用してください。

### `validate`

すべてのロケールファイルをソースバンドルと照合します。デフォルトでは、必須ターゲットキーの充足率を検査し、余分なキーを警告として報告します。キーの欠落、補間パラメータの不一致、ICU MessageFormatの構文エラーがある場合は失敗します。

```bash
internationalizer validate                     # 人間が読みやすい形式で出力
internationalizer validate --json              # 機械可読なJSON形式で出力
internationalizer validate -q                  # 終了コードのみ出力
internationalizer validate --strict             # 翻訳品質に関する規則を適用
internationalizer validate --require-state      # マニフェストが最新であることを必須化
```

`--strict`では翻訳済みの割合も報告します。言語表現としての値がソースと同一の場合、用語集にその値全体についてソースとターゲットが完全に同じ項目がない限り、未翻訳と判定されます。`ignore_case`は考慮されますが、長い値の一部に用語集の語が含まれるだけでは除外されません。strictモードでは、余分なキー、ソースと同一の値、補間・HTML・コード・Markdownリンクの構造変更、用語集違反、設定された複数形の不足があると失敗します。

`--require-state`は、各ターゲットを`.internationalizer.lock`と照合します。キーが未記録の場合や、記録されたソース、翻訳ポリシー、ターゲットのハッシュが古い場合は失敗します。`--strict`と併用できます。

人間向けレポートとJSONレポートでは、次の安定した検出コードを使用します。

| コード | 意味 |
| --- | --- |
| `missing_key` / `extra_key` | ソースとターゲットのキー集合が一致していない |
| `blank_translation` | 空でないソースに対するstrictモードのターゲットが空 |
| `source_identical` | 言語表現としての値がstrictモードでもソースと同一 |
| `protected_structure_mismatch` | 補間、HTML、コード、リンクの構造が変更されている |
| `glossary_violation` | 承認済みのターゲット用語または異表記が見つからない |
| `plural_form_missing` | 設定されたロケールの複数形が不足している |
| `icu_message_syntax` | ソースまたはターゲットのICUメッセージが不正 |
| `icu_argument_mismatch` | ICUの引数名、種類、フォーマッタースタイルが一致していない |
| `icu_selector_mismatch` | セレクターが一致していない、または複数形カテゴリーがターゲットロケールで無効 |
| `untracked` | ターゲットに対応するマニフェストレコードがない |
| `source_stale` | 記録後にソースの内容が変更された |
| `policy_stale` | 生成プロンプトまたはモデル設定が変更された |
| `target_modified` | ターゲットの内容がマニフェストの記録と異なる |

### `detect`

使用中のi18nフレームワークを自動検出し、推奨設定を提案します。

```bash
internationalizer detect
```

サポート対象：react-i18next、next-intl、vue-i18n、プレーンなJSON、Markdownドキュメント。

### `glossary`

翻訳時に適用される言語ごとの用語集を管理します。

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

翻訳メモリ（過去に翻訳された文字列を保存するJSONLキャッシュ）を管理します。

```bash
internationalizer tm stats                     # レコード数を表示
internationalizer tm export                    # JSON形式でダンプ出力
internationalizer tm clear --force             # すべてのレコードを削除
```

## 設定リファレンス

```yaml
# .internationalizer.yml

# ソース言語（デフォルト：en）
source_locale: en

# 翻訳先言語（必須）
target_locales: [fr, de, es, ja, yue, zh-CN, zh-TW, ar]

# 1つ以上のソースからターゲットへのマッピング（必須）。
# {locale}は設定された各ターゲットロケールに置換されます。
bundles:
  - id: app
    source: locales/en.json
    target: locales/{locale}.json
    format: json
  - id: docs
    source: README.md
    target: docs/i18n/{locale}.md
    format: markdown

# 下位互換性：source_pathも引き続きlocales/fr.jsonなどの兄弟ファイルへターゲットをマッピングします。
# 新規プロジェクトではbundlesの使用を推奨します。
# source_path: locales/en.json

# LLMプロバイダー設定
llm:
  # プロバイダー：「anthropic」、「openai」、「gemini」、または「openrouter」（デフォルト：gemini）
  provider: gemini

  # プロバイダーごとのデフォルトモデル名
  #   anthropic:  claude-opus-5
  #   openai:     gpt-5.6-luna（推論強度のデフォルトはmax）
  #   gemini:     gemini-3.8-flash
  #   openrouter: deepseek/deepseek-v4-pro-0813
  model: gemini-3.8-flash

  # APIキーを格納した環境変数名
  api_key_env: GOOGLE_AI_STUDIO_API_KEY

  # OpenAI互換エンドポイントのベースURL（オプション）
  # base_url: https://api.openai.com

  # OpenAI GPT-5シリーズ Responses APIの推論強度
  # （OpenAIプロバイダーのデフォルト：max）
  reasoning_effort: max

  # 特定のターゲットロケールに対するオプションのLLM設定。
  # グローバル設定と同一のプロバイダーを使用する言語オーバーライドは、未指定の項目にグローバル設定を継承します。
  # 異なるプロバイダーを指定した場合、未指定の項目にはそのプロバイダーのデフォルト値が適用されます。
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

# LLM呼び出しあたりのキー数（デフォルト：40）
batch_size: 40

# LLMの並行呼び出し数（デフォルト：4）
concurrency: 4

# ロケール別スタイルガイドMarkdownファイルを格納するディレクトリ（デフォルト：style-guides）
style_guides_dir: style-guides

# ロケール別用語集JSONファイルを格納するディレクトリ（デフォルト：glossary）
glossary_dir: glossary

# 翻訳メモリファイルへのパス（デフォルト：.internationalizer/tm.jsonl）
tm_path: .internationalizer/tm.jsonl

# ソース、ポリシー、ターゲット、および来歴情報のバージョン管理状態
# （デフォルト：.internationalizer.lock。このファイルをコミットしてください）
manifest_path: .internationalizer.lock

# 翻訳とstrict検証に関するオプション規則
validation:
  plural_style: i18next-v4 # ターゲットロケールの複数形を生成して検証
```

ロケール識別子には、`fr`、`pt-BR`、`sr-Latn-RS`など、正しい形式のBCP 47タグを指定してください。正規化すると同一になるターゲットロケールは重複として拒否され、ロケール別のプロバイダーオーバーライドも正規化後の表記で照合されます。上の例では、日本語を含むオーバーライドのないロケールは、グローバルなGemini設定を継承します。

ICU MessageFormatの値は構造として解析されます。単純な引数のほか、`select`、`plural`、`selectordinal`、`number`、`date`、`time`に対応し、メッセージのネスト、複数形のオフセット、数値セレクター、`#`も使用できます。検証では、構文、引数の種類とフォーマッタースタイル、複数形のオフセット、selectの分岐、ターゲットロケールのCLDR複数形カテゴリーを確認します。これらの条件を破るプロバイダー出力は、ロケールファイルや翻訳メモリへ書き込まれる前に拒否されます。

`i18next-v4`を指定すると、認識されたソースの複数形ファミリーが、翻訳時にターゲットロケールのCLDRカテゴリーへ展開されます。ターゲットにしかないカテゴリーでは、ソースファミリーの`_other`値を翻訳テンプレートとして使用します。strict検証ではターゲットに必要なカテゴリーを必須とし、ターゲットロケールで使用しないソース側だけのカテゴリーは任意として扱います。

## スタイルガイド

スタイルガイドは、LLMの翻訳プロンプトに挿入されるMarkdownファイルです。トーン、丁寧さの度合い、タイポグラフィ、その他の言語固有の規則を制御します。

```
style-guides/
  _conventions.md    # すべての言語に共通のルール
  fr.md              # フランス語固有のルール
  ja.md              # 日本語固有のルール
  ar.md              # アラビア語固有のルール
```

### 共通規則 (`_conventions.md`)

すべての言語に適用されるルールを定義します。補間構文、HTMLの保持、文字列種別ごとの規則（ボタン、ラベル、エラーメッセージなど）を指定します。

### 言語別ガイド (`{locale}.md`)

言語固有のルールを定義します。丁寧さの度合い（「tu」と「vous」の使い分けなど）、句読点（ギュメ、逆疑問符など）、複数形、日付や数値の書式、用語集を指定します。

実際の構成例については、[`examples/react-app/style-guides/`](../../examples/react-app/style-guides/)を参照してください。

## 用語集の形式

用語集ファイルは、`{glossary_dir}/{locale}.json`に配置するJSON配列です。

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

`variants`には、承認済みの別表記を指定します。`enforcement`には`error`または`warning`を指定でき、省略時は`error`です。用語は対照表としてLLMプロンプトに挿入され、アプリケーション全体で一貫した翻訳に使用されます。`{"source":"API","target":"API"}`のようにソースとターゲットが完全に同じ項目を登録すると、その値全体はstrict検証の未翻訳判定から除外されます。長い値の一部に`API`が含まれるだけでは除外されません。

## 翻訳メモリ

翻訳メモリはJSONLファイル（1行につき1件のJSONレコード）として保存されます。各レコードには以下の情報が含まれます。

- バンドル、キー、ソース値、翻訳後の値、正規化されたターゲットロケール
- ソースと翻訳ポリシーのハッシュ
- 翻訳に使用したプロバイダーとモデル
- タイムスタンプ

次回以降の実行時には、ソースとポリシーのハッシュが同じ文字列がLLMを呼び出さずにキャッシュから取得されます。デフォルトの保存先はGit管理から除外された`.internationalizer/`ディレクトリ内なので、ローカルキャッシュとして扱われます。翻訳メモリを意図的に共有する場合は、`tm_path`をGit管理対象のパスへ変更してください。レビュー可能な`.internationalizer.lock`マニフェストは別途バージョン管理されます。

## サポート対象の形式

| 形式 | 拡張子 | 処理モード |
|--------|-----------|------|
| JSON | `.json` | キーバリュー（ネスト対応、ドット記法による平坦化） |
| YAML | `.yml`, `.yaml` | キーバリュー（コメントと記述順序を保持） |
| Markdown | `.md`, `.mdx` | ドキュメント全体の翻訳 |

## プロジェクトタイプの検出

`internationalizer detect`は、以下の項目を検査してi18n構成を特定します。

- `package.json`内のreact-i18next、next-intl、vue-i18nなどの依存関係
- 一般的なロケール配置パターンに一致するディレクトリ構造
- ファイル拡張子および命名規則

## アーキテクチャ

```
cmd/internationalizer/     CLIエントリポイントと各コマンドの定義
internal/
  config/                  デフォルト値を適用したYAML設定の読み込み
  detect/                  プロジェクトタイプの自動検出
  formats/                 フォーマットパーサー（JSON、YAML、Markdown）
  glossary/                ロケール別用語集の管理
  llm/                     LLMプロバイダーのインターフェースと実装
    anthropic.go           Anthropic Claudeバックエンド
    openai.go              OpenAI / 互換エンドポイントバックエンド
    gemini.go              Google AI Studio経由のGeminiバックエンド
                           OpenRouterはカスタムbase_urlを指定してopenai.goを使用
  locale/                  BCP 47ロケールIDとCLDR複数形カテゴリー
  message/                 ICU MessageFormatのパーサーと構造比較
  policy/                  安定した翻訳ポリシーのハッシュ化
  state/                   バージョン管理される翻訳マニフェスト
  styleguide/              スタイルガイドの読み込み
  tm/                      JSONL形式の翻訳メモリ
  translate/               翻訳処理のオーケストレーター
  validate/                ロケールの検証と差分抽出
```

## 代替ツールとの比較

| 機能 | Internationalizer | i18next | Crowdin | 汎用LLM |
|---------|------------------|---------|---------|-------------|
| LLMを活用した翻訳 | はい | いいえ | 一部対応 | はい |
| 言語ごとのスタイルガイド | はい | いいえ | いいえ | いいえ |
| 用語集の強制適用 | はい | いいえ | はい | いいえ |
| 翻訳メモリ | はい | いいえ | はい | いいえ |
| CLI / ローカル環境での実行 | はい | 該当なし | いいえ | 手動 |
| Git管理に適したファイル形式 | はい | はい | 一部対応 | 手動 |
| SaaSへの依存なし | はい | はい | いいえ | ツールによる |
| オープンソース（AGPL-3.0） | はい | はい | いいえ | ツールによる |

## ライセンス

[AGPL-3.0](../../LICENSE)

依存関係に関する通知は、[THIRD_PARTY_NOTICES.md](../../THIRD_PARTY_NOTICES.md)を参照してください。

## コントリビューション

開発環境のセットアップおよびガイドラインについては、[CONTRIBUTING.md](../../CONTRIBUTING.md)を参照してください。すべてのコントリビューションにはDCO（開発者原産性証明）への署名（サインオフ）が必須です。
