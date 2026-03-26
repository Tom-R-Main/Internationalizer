> [English (original)](../../README.md)

<p align="center">
  <img src="../../assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

Pipeline d'internationalisation natif IA pour les projets logiciels. Traduisez, validez et gérez les fichiers i18n à l'aide de LLM.

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](LICENSE)

## Pourquoi Internationalizer ?

La plupart des outils i18n sont soit des bibliothèques d'exécution (i18next, react-intl), soit des plateformes SaaS de gestion de clés (Crowdin, Lokalise). Aucun d'entre eux ne résout correctement le véritable problème de la traduction :

- **La traduction manuelle** ne passe pas à l'échelle au-delà de quelques langues
- **Les API de traduction automatique** (Google Translate, DeepL) ignorent votre terminologie, votre ton et vos conventions d'interface utilisateur
- **La traduction LLM générique** fonctionne mieux, mais sans glossaires ni guides de style, vous obtenez des résultats incohérents

Internationalizer est différent. C'est un **pipeline CLI** qui combine la traduction LLM avec :

- **Des glossaires par langue** — appliquent une terminologie cohérente dans toute votre application
- **Des guides de style par langue** — contrôlent le ton, le niveau de formalité, la pluralisation et la typographie
- **Une mémoire de traduction** — ignore les chaînes inchangées, économise de l'argent sur les appels d'API
- **La validation des clés** — détecte les traductions manquantes et les incohérences d'interpolation avant le déploiement

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

Ou compiler depuis les sources :

```bash
git clone https://github.com/Tom-R-Main/Internationalizer.git
cd Internationalizer
go build -o internationalizer ./cmd/internationalizer
```

## Paquets npm

- Les tags Git et les versions des paquets npm doivent correspondre, par exemple `v0.1.0` et `0.1.0`
- Le paquet racine `internationalizer` dépend de paquets de plateforme tels que `internationalizer-darwin-arm64`
- Cibles npm prises en charge : macOS arm64/x64, Linux arm64/x64, Windows x64
- La publication CI nécessite un secret GitHub nommé `NPM_TOKEN`

## Démarrage rapide

1. Créez un fichier de configuration à la racine de votre projet :

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

## Commandes

### `translate`

Trouve les clés manquantes et les traduit via un LLM.

```bash
internationalizer translate                    # traduire toutes les locales
internationalizer translate -l fr              # traduire uniquement le français
internationalizer translate --dry-run          # prévisualiser sans appels d'API
internationalizer translate --batch-size 20    # lots plus petits
internationalizer translate --concurrency 2    # moins d'appels parallèles
```

### `validate`

Vérifie tous les fichiers de locale pour trouver les clés manquantes, les clés en trop et les incohérences d'interpolation.

```bash
internationalizer validate                     # sortie lisible par un humain
internationalizer validate --json              # JSON lisible par une machine
internationalizer validate -q                  # code de sortie uniquement
```

### `detect`

Détecte automatiquement le framework i18n et suggère une configuration.

```bash
internationalizer detect
```

Prend en charge : react-i18next, next-intl, vue-i18n, JSON natif, documents Markdown.

### `glossary`

Gère les termes du glossaire par langue qui sont appliqués lors de la traduction.

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

Gère la mémoire de traduction (cache JSONL des chaînes précédemment traduites).

```bash
internationalizer tm stats                     # afficher le nombre d'enregistrements
internationalizer tm export                    # exporter en JSON
internationalizer tm clear --force             # supprimer tous les enregistrements
```

## Référence de configuration

```yaml
# .internationalizer.yml

# Langue source (par défaut : en)
source_locale: en

# Langues vers lesquelles traduire (requis)
target_locales: [fr, de, es, ja, zh-CN, ar]

# Chemin vers le fichier de locale source (requis)
source_path: locales/en.json

# Paramètres du fournisseur LLM
llm:
  # Fournisseur : "anthropic", "openai", "gemini" ou "openrouter" (par défaut : gemini)
  provider: gemini

  # Noms de modèles par défaut selon le fournisseur :
  #   anthropic:  claude-sonnet-4-6
  #   openai:     gpt-5.4
  #   gemini:     gemini-3.1-pro-preview
  #   openrouter: google/gemini-3-flash-preview
  model: gemini-3.1-pro-preview

  # Variable d'environnement contenant la clé d'API
  api_key_env: GOOGLE_AI_STUDIO_API_KEY

  # URL de base pour les points de terminaison compatibles OpenAI (optionnel)
  # base_url: https://api.openai.com

# Clés par appel LLM (par défaut : 40)
batch_size: 40

# Appels LLM parallèles (par défaut : 4)
concurrency: 4

# Répertoire contenant les fichiers Markdown des guides de style par locale (par défaut : style-guides)
style_guides_dir: style-guides

# Répertoire contenant les fichiers JSON du glossaire par locale (par défaut : glossary)
glossary_dir: glossary

# Chemin vers le fichier de mémoire de traduction (par défaut : .internationalizer/tm.jsonl)
tm_path: .internationalizer/tm.jsonl
```

## Guides de style

Les guides de style sont des fichiers Markdown qui sont injectés dans le prompt de traduction du LLM. Ils contrôlent le ton, le niveau de formalité, la typographie et d'autres conventions spécifiques à la langue.

```
style-guides/
  _conventions.md    # règles partagées pour toutes les langues
  fr.md              # règles spécifiques au français
  ja.md              # règles spécifiques au japonais
  ar.md              # règles spécifiques à l'arabe
```

### Conventions partagées (`_conventions.md`)

Définissez les règles qui s'appliquent à toutes les langues : syntaxe d'interpolation, préservation du HTML, conventions de type de chaîne (boutons vs étiquettes vs erreurs), etc.

### Guides par langue (`{locale}.md`)

Définissez les règles spécifiques à la langue : registre de formalité (tu vs vous), ponctuation (guillemets, points d'interrogation inversés), formes plurielles, formatage des dates/nombres et un glossaire terminologique.

Consultez [`examples/react-app/style-guides/`](examples/react-app/style-guides/) pour un exemple fonctionnel.

## Format du glossaire

Les fichiers de glossaire sont des tableaux JSON stockés dans `{glossary_dir}/{locale}.json` :

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

Les termes sont injectés dans le prompt du LLM sous forme de table terminologique, garantissant une traduction cohérente des termes clés dans toute votre application.

## Mémoire de traduction

La mémoire de traduction est stockée sous forme de fichier JSONL (un enregistrement JSON par ligne). Chaque enregistrement contient :

- La clé et la valeur sources
- La valeur traduite
- Un hachage SHA-256 de la valeur source
- Un horodatage

Lors des exécutions ultérieures, les chaînes inchangées sont servies depuis le cache de la MT sans appeler le LLM, ce qui permet d'économiser du temps et des coûts d'API. Le fichier de MT est adapté à Git et peut être commité avec vos fichiers de locale.

## Formats pris en charge

| Format | Extensions | Mode |
|--------|-----------|------|
| JSON | `.json` | Clé-valeur (imbriqué, aplati par notation pointée) |
| YAML | `.yml`, `.yaml` | Clé-valeur (préserve les commentaires et l'ordre) |
| Markdown | `.md`, `.mdx` | Traduction de document entier |

## Détection du type de projet

`internationalizer detect` identifie votre configuration i18n en vérifiant :

- Les dépendances `package.json` pour react-i18next, next-intl ou vue-i18n
- Les structures de répertoires correspondant aux modèles de locales courants
- Les extensions de fichiers et les conventions de nommage

## Architecture

```
cmd/internationalizer/     Point d'entrée CLI et définitions des commandes
internal/
  config/                  Chargement de la configuration YAML avec valeurs par défaut
  detect/                  Détection automatique du type de projet
  formats/                 Analyseurs de format (JSON, YAML, Markdown)
  glossary/                Gestion du glossaire par locale
  llm/                     Interface du fournisseur LLM + implémentations
    anthropic.go           Backend Anthropic Claude
    openai.go              Backend OpenAI / compatible
    gemini.go              Backend Google Gemini via AI Studio
                           OpenRouter utilise openai.go avec une base_url personnalisée
  styleguide/              Chargeur de guide de style
  tm/                      Mémoire de traduction JSONL
  translate/               Orchestrateur de traduction
  validate/                Validation et comparaison des locales
```

## Comparaison avec les alternatives

| Fonctionnalité | Internationalizer | i18next | Crowdin | LLM générique |
|---------|------------------|---------|---------|-------------|
| Traduction par LLM | Oui | Non | Partiel | Oui |
| Guides de style par langue | Oui | Non | Non | Non |
| Application du glossaire | Oui | Non | Oui | Non |
| Mémoire de traduction | Oui | Non | Oui | Non |
| Exécution CLI / locale | Oui | N/A | Non | Manuel |
| Fichiers adaptés à Git | Oui | Oui | Partiel | Manuel |
| Aucune dépendance SaaS | Oui | Oui | Non | Variable |
| Open source (AGPL-3.0) | Oui | Oui | Non | Variable |

## Licence

[AGPL-3.0](LICENSE)

## Contribution

Consultez [CONTRIBUTING.md](CONTRIBUTING.md) pour la configuration de développement et les directives. Toutes les contributions nécessitent une approbation DCO.

