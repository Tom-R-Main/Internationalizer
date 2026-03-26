# Swedish (sv) — Translation Style Guide

## Language Profile
- Locale: `sv`
- Script: Latin, LTR
- CLDR plural forms: one, other
- Text expansion vs English: 10-15%

## Tone & Formality
Tonen i svensk programvara ska vara naturlig, vänlig och professionell. Använd alltid det personliga pronomenet "du" för att tilltala användaren. Använd aldrig "Ni", då detta uppfattas som föråldrat, onaturligt och distanserat i modern svenska. 

Undvik en alltför formell eller byråkratisk ton. Var direkt och tydlig. Använd aktiv form ("Du har sparat filen") snarare än passiv form ("Filen har sparats av dig") när det är möjligt. I korta systemmeddelanden kan dock passiv form vara acceptabelt för att spara utrymme ("Filen sparades").

## Grammar
- **Sammansatta ord (Särskrivning):** Detta är det absolut vanligaste felet vid översättning till svenska. Svenska sammansatta ord ska *alltid* skrivas ihop. Skriv "användarkonto" (inte "användar konto") och "e-postadress" (inte "e-post adress").
- **Versalisering (Sentence case):** Svenskan använder gemener (små bokstäver) för alla ord utom det första i rubriker, knappar och menyer. Skriv "Skapa nytt konto" (inte "Skapa Nytt Konto"). Egennamn och varumärken behåller dock sin versalisering.
- **Imperativ för knappar:** Använd alltid imperativ (uppmaningsform) för knappar och åtgärder. Exempel: "Spara", "Radera", "Skicka".
- **Bestämd form i gränssnitt:** I svenskt gränssnitt låter det ofta mer naturligt att använda bestämd form där engelskan saknar artikel. Engelskans "Delete account" översätts bäst till "Radera kontot" (snarare än "Radera konto").

## Pluralization
Svenska använder två pluralformer enligt CLDR:
- **One (1):** Används enbart för siffran 1. Exempel: "1 fil raderades", "1 timme kvar", "1 fel uppstod".
- **Other (0, 2-999...):** Används för noll och alla siffror större än 1. Exempel: "0 filer raderades", "2 filer raderades", "5 timmar kvar", "14 fel uppstod".

## Punctuation & Typography
- **Citattecken:** Använd svenska typografiska citattecken (”...”). Undvik raka ("...") eller engelska (“...”).
- **Siffror:** Använd kommatecken som decimalkiljare (ex. 3,14) och fast mellanslag (non-breaking space) som tusentalsavgränsare (ex. 1 000 000).
- **Datum:** Standardformatet är ÅÅÅÅ-MM-DD (ex. 2023-10-24).
- **Tid:** Använd 24-timmarsformat med kolon (ex. 14:30). Använd inte AM/PM.

## Terminology

| English | Swedish | Notes |
|---------|---------|-------|
| Save | Spara | Används som imperativ på knappar. |
| Cancel | Avbryt | Standard för att avbryta en åtgärd i dialogrutor. |
| Delete | Radera | "Ta bort" fungerar också, men "Radera" är tydligare för permanent borttagning. |
| Settings | Inställningar | Används för konfigurationsmenyer. |
| Search | Sök | Imperativform för sökfält och knappar. |
| Error | Fel | Används vid felmeddelanden (ex. "Ett fel uppstod"). |
| Loading | Laddar | Används vid laddningsskärmar eller snurror. |
| Dashboard | Översikt | "Instrumentpanel" är ofta för långt och tekniskt; "Översikt" är modern standard. |
| Notifications | Aviseringar | "Notiser" är också acceptabelt, men "Aviseringar" är vanligast i större system. |
| Sign in | Logga in | Skrivs som två ord. |
| Sign out | Logga ut | Skrivs som två ord. |
| Submit | Skicka | Betyder "Skicka in" (ex. ett formulär). Undvik direktöversättningar som "Underkasta". |
| Profile | Profil | Används för användarens personliga sida/inställningar. |
| Help | Hjälp | Standard för support och dokumentation. |
| Close | Stäng | Används för att stänga fönster, flikar eller dialogrutor. |
