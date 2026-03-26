# German (de) — Translation Style Guide

## Language Profile
- Locale: `de`
- Script: Latin, LTR
- CLDR plural forms: one, other
- Text expansion vs English: 20-30%

## Tone & Formality
Verwende standardmäßig die informelle Anrede „du“ (Kleinschreibung, außer am Satzanfang), da dies für moderne Software, SaaS und Apps der etablierte Standard ist. Der Tonfall sollte freundlich, direkt und professionell sein. Vermeide stark umgangssprachliche Ausdrücke, aber klinge natürlich und menschlich. Bei Fehlermeldungen ist Sachlichkeit wichtig: Vermeide Schuldzuweisungen an den Nutzer (schreibe z. B. „Ein Fehler ist aufgetreten“ anstelle von „Du hast einen Fehler gemacht“). Wenn das Produkt explizit für den traditionellen B2B- oder Bankensektor gedacht ist, wechsle zur formellen Anrede „Sie“ (immer großgeschrieben).

## Grammar
- **Schaltflächen (Buttons) und Menüs:** Verwende immer den Infinitiv (z. B. „Speichern“, „Abbrechen“) und niemals den Imperativ („Speichere“).
- **Komposita (Zusammengesetzte Nomen):** Schreibe zusammengesetzte Wörter im Deutschen immer zusammen (z. B. „Benutzerkonto“, nicht „Benutzer Konto“). Vermeide Leerzeichen in Komposita strikt. Wenn das Wort unübersichtlich lang wird oder englische Begriffe enthält, verwende einen Bindestrich (z. B. „E-Mail-Adresse“, „Desktop-App“).
- **Groß- und Kleinschreibung (Capitalization):** Im Gegensatz zum englischen „Title Case“ wird im Deutschen in UI-Elementen (Überschriften, Buttons) nur das erste Wort sowie alle Substantive großgeschrieben („Sentence case“). Beispiel: „Neues Projekt erstellen“ (nicht „Neues Projekt Erstellen“).
- **Aktiv statt Passiv:** Formuliere Sätze aktiv und direkt, um die UI dynamischer zu machen. Statt „Die Datei wurde erfolgreich hochgeladen“ schreibe besser „Datei erfolgreich hochgeladen“ oder „Du hast die Datei hochgeladen“.

## Pluralization
Das Deutsche verwendet zwei Pluralformen nach dem CLDR-Standard:
- **one (Eins):** Wird ausschließlich für die exakte Zahl 1 verwendet. Beispiel: „1 Datei gelöscht“ oder „1 Benachrichtigung“.
- **other (Andere):** Wird für alle anderen Zahlen verwendet, einschließlich 0, Brüche und negative Zahlen. Beispiele: „0 Dateien gelöscht“, „2 Dateien gelöscht“, „5 Benachrichtigungen“.

## Punctuation & Typography
- **Anführungszeichen:** Verwende die typografisch korrekten deutschen Anführungszeichen („...“) für Zitate oder Hervorhebungen in der UI. Verwende keine geraden ("...") oder englischen (“...”) Anführungszeichen.
- **Zahlenformate:** Verwende das Komma als Dezimaltrennzeichen (z. B. 3,14) und den Punkt als Tausendertrennzeichen (z. B. 1.000.000).
- **Datum:** Das Standardformat ist TT.MM.JJJJ (z. B. 31.12.2023).
- **Uhrzeit:** Verwende das 24-Stunden-Format mit einem Doppelpunkt als Trennzeichen (z. B. 14:30 Uhr). Die Zusätze AM/PM werden im Deutschen nicht verwendet.

## Terminology

| English | German | Notes |
|---------|---------|-------|
| Save | Speichern | Infinitiv für Schaltflächen verwenden. |
| Cancel | Abbrechen | Standardübersetzung für Dialoge und Formulare. |
| Delete | Löschen | Infinitiv für Schaltflächen verwenden. |
| Settings | Einstellungen | Immer im Plural verwenden. |
| Search | Suchen | Als Button (Infinitiv) oder Platzhalter im Suchfeld. |
| Error | Fehler | Sachlich und neutral bleiben. |
| Loading | Wird geladen... | Alternativ auch kurz „Lädt...“ bei Platzmangel. |
| Dashboard | Dashboard | Wird im Deutschen als Lehnwort beibehalten, nicht übersetzen. |
| Notifications | Benachrichtigungen | Standardbegriff für Alerts und System-Updates. |
| Sign in | Anmelden | Bevorzugt gegenüber dem Anglizismus „Einloggen“. |
| Sign out | Abmelden | Bevorzugt gegenüber dem Anglizismus „Ausloggen“. |
| Submit | Senden | Alternativ „Absenden“, je nach Kontext des Formulars. |
| Profile | Profil | Standardbegriff für Benutzerprofile. |
| Help | Hilfe | Standardbegriff für Support-Bereiche. |
| Close | Schließen | Infinitiv für das Schließen von Fenstern oder Dialogen. |
