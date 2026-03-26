# French (fr) — Translation Style Guide

## Language Profile
- Locale: `fr`
- Script: Latin, LTR
- CLDR plural forms: one, many, other
- Text expansion vs English: 15-20%

## Tone & Formality
Le ton doit être professionnel, clair, concis et courtois. Dans les interfaces logicielles, utilisez toujours le vouvoiement (« vous », « votre ») pour vous adresser à l'utilisateur. Le tutoiement (« tu ») est à proscrire, sauf consigne contraire spécifique à une marque très décontractée. 

Privilégiez une voix active et directe pour guider l'utilisateur de manière fluide. Évitez le jargon technique inutile et les formulations trop familières ou excessivement rigides. L'objectif est d'instaurer un climat de confiance et de clarté.

## Grammar
- **Boutons et menus (Infinitif vs Impératif) :** Utilisez toujours l'infinitif pour les actions déclenchées par l'utilisateur (ex. : *Enregistrer*, *Modifier*, *Copier*). N'utilisez pas l'impératif (*Enregistrez*), car c'est l'utilisateur qui donne l'ordre au système.
- **Capitalisation (Casse de phrase) :** Le français n'utilise pas la casse titre (Title Case) de l'anglais. Seul le premier mot d'un titre, d'un bouton ou d'une étiquette prend une majuscule, ainsi que les noms propres (ex. : *Paramètres du compte*, et non *Paramètres Du Compte*).
- **Étiquettes et champs (Phrases nominales) :** Pour les titres de sections ou les étiquettes de champs, privilégiez les groupes nominaux courts (ex. : *Nom d'utilisateur* au lieu de *Entrez votre nom*).
- **Voix passive :** Limitez l'usage de la voix passive, souvent lourde en français. Traduisez *The file was deleted by the system* par *Le système a supprimé le fichier* ou simplement *Fichier supprimé*.
- **Inclusivité et neutralité :** Privilégiez des tournures neutres pour éviter la surcharge typographique (points médians). Utilisez *Droits d'accès* plutôt que *Droits de l'utilisateur / l'utilisatrice*.

## Pluralization
- **one :** S'applique à 0, 1, et aux nombres décimaux inférieurs à 2. Contrairement à l'anglais, le zéro prend la marque du singulier en français (ex. : *0 fichier*, *1 erreur*, *1,5 heure*).
- **many :** S'applique aux grands nombres ronds (comme les millions) suivis de la préposition « de » (ex. : *1 million de fichiers*). Dans la plupart des interfaces UI standards, cette catégorie est rarement distincte de « other ».
- **other :** S'applique à tous les autres nombres à partir de 2 (ex. : *2 fichiers*, *100 éléments*).

## Punctuation & Typography
- **Espaces avec la ponctuation :** En français, il faut obligatoirement une espace insécable avant les signes de ponctuation doubles (`:`, `;`, `?`, `!`).
- **Guillemets :** Utilisez les guillemets français en chevrons avec des espaces insécables à l'intérieur (« texte »). Évitez les guillemets anglais (" ").
- **Nombres :** Utilisez la virgule (`,`) comme séparateur décimal et l'espace insécable comme séparateur de milliers (ex. : *1 234,56*).
- **Dates :** Le format standard est JJ/MM/AAAA (ex. : *31/12/2023*).
- **Heures :** Utilisez le format 24 heures avec la lettre « h » comme séparateur (ex. : *14 h 30*). Le format *14:30* est toléré si l'espace est très limité.

## Terminology

| English | French | Notes |
|---------|---------|-------|
| Save | Enregistrer | Préféré à « Sauvegarder » (réservé aux backups de données). |
| Cancel | Annuler | Standard pour interrompre une action ou fermer une modale. |
| Delete | Supprimer | Préféré à « Effacer » pour la suppression d'éléments ou de fichiers. |
| Settings | Paramètres | Standard. « Réglages » est parfois utilisé sur macOS/iOS, mais « Paramètres » est universel. |
| Search | Rechercher | À l'infinitif pour un bouton. Utilisez le nom « Recherche » pour un titre ou un champ. |
| Error | Erreur | Standard. |
| Loading | Chargement | Utilisez « Chargement... » ou « Chargement en cours » pour les indicateurs de progression. |
| Dashboard | Tableau de bord | Ne pas laisser en anglais. |
| Notifications | Notifications | Standard. |
| Sign in | Se connecter | Préféré à « S'identifier ». |
| Sign out | Se déconnecter | Standard. |
| Submit | Envoyer | Ou « Valider ». Évitez « Soumettre » qui est un faux-ami partiel et sonne peu naturel. |
| Profile | Profil | Attention à l'orthographe (un seul « l » en français). |
| Help | Aide | Standard. |
| Close | Fermer | Standard pour quitter une fenêtre ou une boîte de dialogue. |
