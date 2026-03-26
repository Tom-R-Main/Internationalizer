# Danish (da) — Translation Style Guide

## Language Profile
- Locale: `da`
- Script: Latin, LTR
- CLDR plural forms: one, other
- Text expansion vs English: 15-20%

## Tone & Formality
Tonen i dansk software og brugergrænseflader (UI) skal være uformel, men professionel. Brug altid tiltaleformen "du" og "dig" (med småt). Undgå den forældede og stive høflige form "De/Dem". Sproget skal være direkte, venligt og hjælpsomt uden at være overdrevet entusiastisk. Undgå unødvendige fyldord, og hold sætningerne korte, klare og præcise, så brugeren nemt kan afkode budskabet.

## Grammar
- **Sammensatte navneord (Compound nouns):** På dansk skrives sammensatte ord altid i ét ord (f.eks. "brugerkonto", "kontrolpanel", "fejlmeddelelse"). Undgå særskrivning (at dele ordene op), da dette er en meget hyppig fejl forårsaget af engelsk indflydelse.
- **Bydeform (Imperativ) på knapper:** Brug altid bydeform til handlinger, menuer og knapper (f.eks. "Gem", "Slet", "Kopier"). Undgå at bruge navnemåde (infinitiv) som "At gemme" eller "Gemme".
- **Store og små bogstaver (Sentence case):** Brug kun stort begyndelsesbogstav i det første ord i overskrifter, knapper og menuer, medmindre der indgår egennavne (f.eks. "Opret ny bruger", ikke "Opret Ny Bruger"). Kopier ikke den engelske "Title Case"-stil.
- **Aktivt sprog:** Brug aktiv form frem for passiv, hvor det er muligt, for at gøre teksten mere direkte (f.eks. "Du har slettet filen" frem for "Filen blev slettet af dig"). Passiv (s-form) kan dog bruges i korte systembeskeder (f.eks. "Ændringerne gemmes").

## Pluralization
Dansk bruger to pluraliskategorier ifølge CLDR-standarden:
- **One (ental):** Bruges udelukkende om tallet 1 (f.eks. "1 fil", "1 bruger", "1 fejl").
- **Other (flertal):** Bruges om alle andre tal, inklusive 0 og decimaltal (f.eks. "0 filer", "2 brugere", "10,5 sekunder").

## Punctuation & Typography
- **Anførselstegn:** Brug dobbelte anførselstegn ("...") eller vinkelanførselstegn (»...«).
- **Talformat:** Brug komma som decimalseparator og punktum som tusindtalsseparator (f.eks. 1.234.567,89).
- **Datoformat:** Standardformatet er DD.MM.ÅÅÅÅ (f.eks. 31.12.2023) eller DD/MM-ÅÅÅÅ.
- **Tidsformat:** Brug 24-timers format med kolon (f.eks. 14:30). Undgå de engelske AM/PM-betegnelser.

## Terminology

| English | Danish | Notes |
|---------|---------|-------|
| Save | Gem | Bydeform (imperativ) til knapper og handlinger. |
| Cancel | Annuller | Standard for at afbryde en handling. |
| Delete | Slet | Bydeform til knapper. |
| Settings | Indstillinger | Altid i flertal. |
| Search | Søg | Bydeform til knapper. Brug "Søgning", hvis det er et navneord. |
| Error | Fejl | Samme ord i både ental og flertal. |
| Loading | Indlæser | Angiver en igangværende handling (nutid). |
| Dashboard | Oversigt | "Dashboard" kan bruges, men "Oversigt" eller "Kontrolpanel" foretrækkes ofte på dansk. |
| Notifications | Notifikationer | "Meddelelser" kan også bruges afhængigt af konteksten. |
| Sign in | Log ind | Skrives i to ord som handling (udsagnsord). |
| Sign out | Log ud | Skrives i to ord som handling (udsagnsord). |
| Submit | Indsend | "Send" eller "Godkend" kan også bruges afhængigt af formularen. |
| Profile | Profil | Standardoversættelse. |
| Help | Hjælp | Standardoversættelse. |
| Close | Luk | Bydeform til knapper og dialogbokse. |
