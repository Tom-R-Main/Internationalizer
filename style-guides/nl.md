# Dutch (nl) — Translation Style Guide

## Language Profile
- Locale: `nl`
- Script: Latin, LTR
- CLDR plural forms: one, other
- Text expansion vs English: 15-25%

## Tone & Formality
Gebruik een actieve, vriendelijke en professionele toon. Voor moderne software en apps is het informeel aanspreken met 'je' en 'jouw' de absolute standaard. Vermijd het formele 'u' en 'uw', tenzij de software specifiek bedoeld is voor zeer formele medische, juridische of financiële contexten. Wees direct en duidelijk, maar klink niet te amicaal of populair. Vermijd passieve zinsconstructies (met 'worden' of 'zijn') zoveel mogelijk; kies voor een actieve zinsbouw om de interface vlot en begrijpelijk te houden.

## Grammar
- **Knoppen en menu's (Infinitief):** Gebruik het hele werkwoord (infinitief) voor losse acties op knoppen, tabbladen en in menu's. Vertaal 'Save' als 'Opslaan' (niet 'Sla op') en 'Delete' als 'Verwijderen'.
- **Instructies (Gebiedende wijs):** Gebruik de gebiedende wijs (de stam van het werkwoord) voor instructies in lopende tekst en dialoogvensters. Bijvoorbeeld: 'Vul je wachtwoord in' (niet 'Vul uw wachtwoord in' of 'Invullen').
- **Samenstellingen (Engelse ziekte):** Schrijf samengestelde zelfstandige naamwoorden in het Nederlands altijd aan elkaar vast. LLM's nemen vaak onterecht de Engelse spaties over. Vertaal 'user account' als 'gebruikersaccount' (niet 'gebruikers account') en 'email address' als 'e-mailadres'.
- **Hoofdlettergebruik:** Gebruik 'sentence case' voor UI-elementen. Alleen de eerste letter van een zin, label of knop krijgt een hoofdletter. Neem het Engelse 'Title Case' niet over. 'Account Settings' wordt 'Accountinstellingen' (niet 'Account Instellingen').

## Pluralization
Het Nederlands kent twee meervoudsvormen volgens de CLDR-regels:
- **one (enkelvoud):** Wordt uitsluitend gebruikt voor het getal 1. Bijvoorbeeld: "1 bestand geselecteerd", "1 map verwijderd".
- **other (meervoud):** Wordt gebruikt voor alle andere getallen, inclusief 0. Bijvoorbeeld: "0 bestanden geselecteerd", "2 mappen verwijderd", "145 gebruikers gevonden".

## Punctuation & Typography
- **Aanhalingstekens:** Gebruik enkele aanhalingstekens ('...') om de nadruk te leggen op UI-elementen in de tekst. Gebruik dubbele aanhalingstekens ("...") alleen voor letterlijke citaten.
- **Getallen:** Gebruik een komma als decimaalscheidingsteken en een punt als duizendtallenscheidingsteken (bijv. 1.234.567,89).
- **Datum:** De standaard numerieke notatie is DD-MM-YYYY (bijv. 31-12-2023). Bij uitgeschreven datums gebruiken we de notatie '31 december 2023' (maanden met een kleine letter).
- **Tijd:** Gebruik altijd de 24-uursnotatie met een dubbele punt (bijv. 14:30). Gebruik geen AM/PM.

## Terminology

| English | Dutch | Notes |
|---------|---------|-------|
| Save | Opslaan | Gebruik de infinitief voor knoppen. |
| Cancel | Annuleren | Standaard vertaling voor het afbreken van een actie. |
| Delete | Verwijderen | Gebruik 'Verwijderen', vermijd 'Wissen' tenzij het om cache/geschiedenis gaat. |
| Settings | Instellingen | Altijd met een hoofdletter aan het begin, meervoud. |
| Search | Zoeken | Infinitief voor de zoekknop of het zoekveld. |
| Error | Fout | Gebruik 'Fout' of 'Foutmelding' afhankelijk van de context. |
| Loading | Laden... | Vaak gebruikt met beletselteken (...) om een lopend proces aan te geven. |
| Dashboard | Dashboard | Dit Engelse leenwoord is volledig ingeburgerd in Nederlandse software. |
| Notifications | Meldingen | Gebruik 'Meldingen', dit is natuurlijker dan het letterlijke 'Notificaties'. |
| Sign in | Inloggen | 'Aanmelden' is ook acceptabel, maar 'Inloggen' is de moderne standaard. |
| Sign out | Uitloggen | Tegenhanger van Inloggen. (Gebruik 'Afmelden' als 'Aanmelden' is gebruikt). |
| Submit | Verzenden | Gebruik 'Verzenden' voor formulieren. 'Indienen' kan voor formele aanvragen. |
| Profile | Profiel | Standaard vertaling voor gebruikersprofiel. |
| Help | Help | Blijft 'Help' in UI-context (niet 'Helpen'). |
| Close | Sluiten | Infinitief voor het sluiten van vensters of modale dialogen. |
