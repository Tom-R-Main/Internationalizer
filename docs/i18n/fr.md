> [English (original)](../../README.md)

<p align="center">
  <img src="../../assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

Pipeline d'internationalisation natif pour l'IA destiné aux projets logiciels. Traduisez, validez et gérez les fichiers i18n à l'aide de LLM.

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](../../LICENSE)

<p align="center">
<a href="ar.md">العربية</a> · <a href="bn.md">বাংলা</a> · <a href="cs.md">Čeština</a> · <a href="da.md">Dansk</a> · <a href="de.md">Deutsch</a> · <a href="el.md">Ελληνικά</a> · <a href="es.md">Español</a> · <a href="fi.md">Suomi</a> · <a href="fr.md">Français</a> · <a href="he.md">עברית</a> · <a href="hi.md">हिन्दी</a> · <a href="id.md">Indonesia</a> · <a href="it.md">Italiano</a> · <a href="ja.md">日本語</a> · <a href="ko.md">한국어</a> · <a href="ms.md">Bahasa Melayu</a><br><a href="nl.md">Nederlands</a> · <a href="pa.md">ਪੰਜਾਬੀ</a> · <a href="pl.md">Polski</a> · <a href="pt-BR.md">Português</a> · <a href="ro.md">Română</a> · <a href="ru.md">Русский</a> · <a href="sv.md">Svenska</a> · <a href="te.md">తెలుగు</a> · <a href="th.md">ไทย</a> · <a href="tr.md">Türkçe</a> · <a href="uk.md">Українська</a> · <a href="vi.md">Tiếng Việt</a> · <a href="yue.md">粵語</a> · <a href="zh-CN.md">简体中文</a> · <a href="zh-TW.md">繁體中文</a>
</p>

---
<!-- internationalizer:unit markdown:why-internationalizer -->
## Pourquoi Internationalizer ?

La plupart des outils i18n sont soit des bibliothèques d'exécution (i18next, react-intl), soit des plateformes SaaS de gestion de clés (Crowdin, Lokalise). Aucun d'entre eux ne résout correctement le véritable problème de traduction :

- **La traduction manuelle** ne passe pas à l'échelle au-delà de quelques langues
- **Les API de traduction automatique** (Google Translate, DeepL) ignorent votre terminologie, votre ton et vos conventions d'interface utilisateur
- **La traduction par LLM générique** fonctionne mieux, mais sans glossaires ni guides de style, vous obtenez des résultats incohérents

Internationalizer est différent. Il s'agit d'un **pipeline CLI** qui combine la traduction par LLM avec :

- **Des glossaires par langue** — appliquent une terminologie cohérente dans l'ensemble de votre application
- **Des guides de style par langue** — contrôlent le ton, le niveau de formalité, la pluralisation et la typographie
- **Une mémoire de traduction** — ignore les chaînes inchangées, réduisant les coûts d'appels d'API
- **Une validation déterministe** — intercepte les clés manquantes ou superflues, les dérives de structures protégées, les non-respects du glossaire, ainsi que les erreurs de pluriel ou ICU avant la mise en production
<!-- internationalizer:unit markdown:installation -->
## Installation

Installer depuis npm :

```bash
npm install -g internationalizer
```

Ou exécuter sans installation globale :

```bash
npx internationalizer --help
```

Le paquet npm installe le binaire précompilé correspondant depuis npm via des dépendances optionnelles spécifiques à la plateforme.

Installer avec Go :

```bash
go install github.com/Tom-R-Main/Internationalizer/cmd/internationalizer@latest
```

Ou compiler à partir des sources :

```bash
git clone https://github.com/Tom-R-Main/Internationalizer.git
cd Internationalizer
go build -o internationalizer ./cmd/internationalizer
```
<!-- internationalizer:unit markdown:npm-packages -->
## Paquets npm

- Les tags Git et les versions de paquets npm doivent correspondre, par exemple `v0.1.0` et `0.1.0`
- Le paquet racine `internationalizer` dépend de paquets de plateforme tels que `internationalizer-darwin-arm64`
- Cibles npm prises en charge : macOS arm64/x64, Linux arm64/x64, Windows x64
- La publication CI nécessite un secret GitHub nommé `NPM_TOKEN`
<!-- internationalizer:unit markdown:quick-start -->
## Démarrage rapide

1. Créez un fichier de configuration à la racine de votre projet :

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

2. Définissez votre clé d'API :

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

3. Prévisualisez ce qui sera traduit :

```bash
internationalizer translate --dry-run
```

4. Lancez la traduction :

```bash
internationalizer translate
```

5. Validez toutes les locales :

```bash
internationalizer validate
```
<!-- internationalizer:unit markdown:commands -->
## Commandes

### `translate`

Trouvez les clés manquantes ou obsolètes et traduisez-les via un LLM.

```bash
internationalizer translate                    # translate all locales
internationalizer translate -l fr              # translate French only
internationalizer translate --dry-run          # preview without API calls
internationalizer translate --adopt-existing   # baseline existing translations without API calls
internationalizer translate --refresh-policy   # refresh prompt/style/model-stale entries
internationalizer translate --batch-size 20    # smaller batches
internationalizer translate --concurrency 2    # fewer parallel calls
```

L'état de la traduction signale de manière indépendante les conditions manquante, source-stale, policy-stale, actuelle et modifiée manuellement, de sorte qu'une modification manuelle ne puisse dissimuler un changement de source ou de stratégie. Les valeurs policy-stale sont signalées mais ne sont retraduites qu'avec `--refresh-policy`. Les valeurs modifiées manuellement ne sont jamais écrasées automatiquement. Utilisez `--adopt-existing` lors de l'intégration du manifeste à des traductions révisées ou lors de l'acceptation explicite d'une modification manuelle révisée en tant que nouvelle référence.

### `validate`

Vérifiez tous les fichiers de locale par rapport à leurs bundles sources. La validation par défaut contrôle la couverture structurelle (le pourcentage de clés cibles requises présentes), signale les clés superflues comme des avertissements et échoue en cas de clés manquantes, de discordances d'interpolation ou de structure ICU MessageFormat non valide.

```bash
internationalizer validate                     # human-readable output
internationalizer validate --json              # machine-readable JSON
internationalizer validate -q                  # exit code only
internationalizer validate --strict             # enforce translation quality rules
internationalizer validate --require-state      # require current manifest provenance
```

`--strict` signale également la couverture traduite. Une valeur linguistique identique à sa source est considérée comme non traduite, à moins que le glossaire ne contienne explicitement une entrée exacte avec source identique et cible identique pour la valeur complète ; `ignore_case` est respecté, mais un terme de glossaire intégré dans une valeur plus longue ne constitue pas une exemption. Le mode strict échoue en cas de clés superflues, de valeurs identiques à la source, de modifications de structure d'interpolation/HTML/code/liens Markdown, d'infractions au glossaire et de formes de pluriel configurées manquantes.

`--require-state` vérifie chaque cible par rapport à `.internationalizer.lock`. Il échoue lorsqu'une clé n'est pas suivie, ou lorsque sa source enregistrée, sa stratégie de traduction ou son hachage cible est obsolète. Il peut être combiné avec `--strict`.

Les rapports lisibles par un humain et JSON utilisent des codes de constatation stables :

| Code | Signification |
| --- | --- |
| `missing_key` / `extra_key` | Les ensembles de clés source et cible diffèrent |
| `blank_translation` | Une source non vide a une cible vide en mode strict |
| `source_identical` | Une valeur linguistique en mode strict reste non traduite |
| `protected_structure_mismatch` | La structure d'interpolation, de HTML, de code ou de lien a changé |
| `glossary_violation` | Aucun terme cible ou variante approuvé n'a été trouvé |
| `plural_form_missing` | Une forme de pluriel configurée pour la locale est absente |
| `icu_message_syntax` | Un message ICU source ou cible est malformé |
| `icu_argument_mismatch` | Les noms, les types ou les styles de formateur des arguments ICU diffèrent |
| `icu_selector_mismatch` | Les sélecteurs diffèrent ou une catégorie de pluriel n'est pas valide pour la locale cible |
| `untracked` | Aucun enregistrement de manifeste n'existe pour la cible |
| `source_stale` | Le contenu source a changé après la traduction enregistrée |
| `policy_stale` | Le prompt généré ou les paramètres du modèle ont changé |
| `target_modified` | Le contenu cible diffère de l'enregistrement dans le manifeste |

### `detect`

Détectez automatiquement le framework i18n et suggérez une configuration.

```bash
internationalizer detect
```

Prend en charge : react-i18next, next-intl, vue-i18n, JSON standard, documents Markdown.

### `glossary`

Gérez les termes de glossaire par langue qui sont appliqués lors de la traduction.

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

Gérez la mémoire de traduction (cache JSONL des chaînes précédemment traduites).

```bash
internationalizer tm stats                     # show record counts
internationalizer tm export                    # dump as JSON
internationalizer tm clear --force             # delete all records
```
<!-- internationalizer:unit markdown:configuration-reference -->
## Référence de configuration

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

Les identifiants de locale doivent être des balises BCP 47 bien formées telles que `fr`, `pt-BR` ou `sr-Latn-RS`. Les locales cibles équivalentes au sens canonique sont rejetées en tant que doublons, et les surcharges de fournisseur spécifiques à une locale s'alignent sur la graphie canonique équivalente. Dans l'exemple ci-dessus, les locales sans surcharge — y compris le japonais — héritent de la configuration globale Gemini.

Les valeurs ICU MessageFormat sont analysées structurellement. Les arguments simples, `select`, `plural`, `selectordinal`, `number`, `date` et `time` sont pris en charge, y compris les messages imbriqués, les décalages de pluriel, les sélecteurs de nombres exacts et `#`. La validation vérifie la syntaxe, les types d'arguments et les styles de formateur, les décalages de pluriel, l'identité des branches select, ainsi que les catégories de pluriel CLDR de la locale cible. Les résultats du fournisseur qui enfreignent ces invariants sont rejetés avant l'écriture d'un fichier de locale ou d'un enregistrement de mémoire de traduction.

Avec `i18next-v4`, les familles de pluriels sources reconnues sont développées lors de la traduction selon les catégories CLDR de la locale cible. Une catégorie propre à la cible utilise la valeur `_other` de la famille source comme modèle de traduction. La validation stricte exige ces catégories cibles ; les catégories propres à la source sont facultatives pour les locales cibles qui ne les utilisent pas.
<!-- internationalizer:unit markdown:style-guides -->
## Guides de style

Les guides de style sont des fichiers Markdown injectés dans le prompt de traduction du LLM. Ils contrôlent le ton, le niveau de formalité, la typographie et d'autres conventions propres à chaque langue.

```
style-guides/
  _conventions.md    # shared rules for all languages
  fr.md              # French-specific rules
  ja.md              # Japanese-specific rules
  ar.md              # Arabic-specific rules
```

### Conventions partagées (`_conventions.md`)

Définissez les règles qui s'appliquent à toutes les langues : syntaxe d'interpolation, préservation du code HTML, conventions sur les types de chaînes (boutons vs étiquettes vs erreurs), etc.

### Guides par langue (`{locale}.md`)

Définissez les règles propres à chaque langue : registre de formalité (tu vs vous), ponctuation (guillemets, points d'interrogation inversés), formes de pluriel, formatage des dates et nombres, ainsi qu'un glossaire terminologique.

Les guides de style constituent des entrées de stratégie durables et non des sorties générées. Internationalizer les lit mais ne les réécrit jamais. Leur contenu est haché séparément du glossaire et du contrat de prompt, de sorte qu'une modification du code de l'application ne rende pas une traduction obsolète. L'édition d'un guide marque délibérément cette locale pour révision de stratégie ; la modification de la formulation interne du prompt ne le fait pas, sauf si la version du contrat de prompt change également.

Consultez [`examples/react-app/style-guides/`](../../examples/react-app/style-guides/) pour un exemple fonctionnel.
<!-- internationalizer:unit markdown:glossary-format -->
## Format du glossaire

Les fichiers de glossaire sont des tableaux JSON stockés dans `{glossary_dir}/{locale}.json` :

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

`variants` répertorie les autres formes cibles approuvées. `enforcement` peut prendre la valeur `error`, `warning`, ou être omis pour appliquer le comportement d'erreur par défaut. Les termes sont injectés dans le prompt du LLM sous forme de tableau terminologique, garantissant ainsi une traduction cohérente dans toute votre application. Une entrée exacte telle que `{"source":"API","target":"API"}` exempte également cette valeur complète identique à la source des constats de valeurs non traduites en mode strict ; elle n'exempte pas une valeur plus longue contenant simplement `API`.
<!-- internationalizer:unit markdown:translation-memory -->
## Mémoire de traduction

La mémoire de traduction est stockée sous forme de fichier JSONL (un enregistrement JSON par ligne). Chaque enregistrement contient :

- Le bundle, la clé, la valeur source, la valeur traduite et la locale cible canonique
- Les hachages de la source, du guide de style, du glossaire, du contrat de prompt et de la stratégie combinée
- Le fournisseur et le modèle qui ont produit la traduction
- Un horodatage

Lors des exécutions ultérieures, les chaînes présentant les mêmes hachages de source et de stratégie sont servies depuis le cache sans faire appel au LLM. Le chemin par défaut se situe sous le répertoire ignoré `.internationalizer/`, de sorte qu'il reste un cache local. Définissez `tm_path` vers un emplacement suivi si votre projet partage intentionnellement la mémoire de traduction. Le manifeste révisable `.internationalizer.lock` est versionné séparément.
<!-- internationalizer:unit markdown:supported-formats -->
## Formats pris en charge

| Format | Extensions | Mode |
| --- | --- | --- |
| JSON | `.json` | Clé-valeur (imbriqué, aplati en notation avec points) |
| YAML | `.yml`, `.yaml` | Clé-valeur (préserve les commentaires et l'ordre) |
| Markdown | `.md`, `.mdx` | Préambule et sections de niveau H2 |

Les cibles Markdown contiennent des commentaires invisibles `internationalizer:unit` avant les sections H2. Ces repères stables permettent à Internationalizer d'ajouter, déplacer ou modifier une section source sans retraduire les sections non concernées. Les documents existants non marqués reçoivent ces repères lors de leur prochaine mise à jour réussie.
<!-- internationalizer:unit markdown:project-type-detection -->
## Détection du type de projet

`internationalizer detect` identifie votre configuration i18n en vérifiant :

- Les dépendances dans `package.json` pour react-i18next, next-intl ou vue-i18n
- Les structures de répertoires correspondant aux modèles de locales courants
- Les extensions de fichiers et conventions de nommage
<!-- internationalizer:unit markdown:architecture -->
## Architecture

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
## Comparaison avec d'autres solutions

| Fonctionnalité | Internationalizer | i18next | Crowdin | LLM générique |
| --- | --- | --- | --- | --- |
| Traduction basée sur l'IA / LLM | Oui | Non | Partiel | Oui |
| Guides de style par langue | Oui | Non | Non | Non |
| Application des glossaires | Oui | Non | Oui | Non |
| Mémoire de traduction | Oui | Non | Oui | Non |
| Exécution en CLI / locale | Oui | N/A | Non | Manuel |
| Fichiers compatibles avec Git | Oui | Oui | Partiel | Manuel |
| Aucune dépendance SaaS | Oui | Oui | Non | Variable |
| Open source (AGPL-3.0) | Oui | Oui | Non | Variable |
<!-- internationalizer:unit markdown:license -->
## Licence

[AGPL-3.0](../../LICENSE)

Consultez [THIRD_PARTY_NOTICES.md](../../THIRD_PARTY_NOTICES.md) pour les mentions relatives aux dépendances tierces.
<!-- internationalizer:unit markdown:contributing -->
## Contribution

Consultez [CONTRIBUTING.md](../../CONTRIBUTING.md) pour la configuration du développement et les directives. Toutes les contributions requièrent la signature du DCO.
