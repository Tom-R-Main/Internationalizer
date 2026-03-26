# Italian (it) — Translation Style Guide

## Language Profile
- Locale: `it`
- Script: Latin, LTR
- CLDR plural forms: one, other
- Text expansion vs English: 15-25%

## Tone & Formality
Usa un tono professionale, chiaro e amichevole. Per le interfacce utente (UI) moderne, rivolgiti all'utente in modo diretto usando la seconda persona singolare informale ("tu"). Evita il "Lei" formale a meno che non sia strettamente richiesto dal contesto (es. software bancario, medico o strettamente legale). Mantieni le frasi brevi, concise e orientate all'azione, evitando un linguaggio eccessivamente burocratico o colloquiale.

## Grammar
- **Pulsanti e Call to Action (CTA):** Usa l'imperativo di seconda persona singolare per i comandi (es. "Salva", "Copia", "Invia"). Sii coerente in tutta l'interfaccia.
- **Maiuscole (Capitalization):** A differenza dell'inglese che usa il "Title Case", in italiano si usa il "Sentence case" per titoli, menu, etichette e pulsanti. Solo la prima lettera della prima parola va in maiuscolo (es. usa "Impostazioni di sistema", non "Impostazioni Di Sistema").
- **Forma passiva e impersonale:** Evita la forma passiva quando possibile, poiché appesantisce la lettura. Usa la forma attiva o impersonale con il "si" (es. preferisci "Impossibile trovare il file" o "File non trovato" rispetto a "Il file non è stato trovato").
- **Posizione degli aggettivi:** Inserisci l'aggettivo dopo il sostantivo, come da regola generale italiana, specialmente nelle etichette UI (es. "Impostazioni avanzate", non "Avanzate impostazioni").

## Pluralization
L'italiano utilizza due categorie plurali CLDR:
- **One:** usato esclusivamente per il numero 1 (es. "1 file", "1 nuovo messaggio").
- **Other:** usato per lo zero e per tutti gli altri numeri maggiori di 1 (es. "0 file", "2 messaggi", "100 utenti").
Assicurati sempre che articoli, aggettivi e participi concordino in genere e numero con il sostantivo a cui si riferiscono (es. "1 riga eliminata" vs "3 righe eliminate").

## Punctuation & Typography
- **Virgolette:** Usa le virgolette alte (" ") per l'interfaccia utente standard. Riserva le virgolette caporali (« ») solo per citazioni lunghe o documentazione formale.
- **Numeri:** Usa la virgola (,) come separatore decimale e il punto (.) come separatore delle migliaia (es. 1.234,56).
- **Data:** Il formato standard è GG/MM/AAAA (es. 31/12/2023).
- **Ora:** Usa il formato a 24 ore separato dai due punti (es. 14:30, non 2:30 PM).
- **Punteggiatura:** Non inserire mai lo spazio prima dei segni di punteggiatura doppi (: ; ! ?).

## Terminology

| English | Italian | Notes |
|---------|---------|-------|
| Save | Salva | Imperativo (2ª persona singolare). |
| Cancel | Annulla | Usato per interrompere un'azione o chiudere un modale senza salvare. |
| Delete | Elimina | Preferito a "Cancella" per la rimozione definitiva di file o dati. |
| Settings | Impostazioni | Sostantivo plurale. |
| Search | Cerca | Imperativo. |
| Error | Errore | Sostantivo singolare. |
| Loading | Caricamento | Sostantivo (spesso usato come "Caricamento in corso..."). Evitare il gerundio "Caricando". |
| Dashboard | Dashboard | Lasciare in inglese. È un prestito ormai standard (femminile: "la dashboard"). |
| Notifications | Notifiche | Sostantivo plurale. |
| Sign in | Accedi | Imperativo. Evitare "Log in" o "Entra". |
| Sign out | Esci | Imperativo. Evitare "Log out" o "Disconnettiti" se lo spazio è limitato. |
| Submit | Invia | Imperativo. Usato per confermare moduli e form. |
| Profile | Profilo | Sostantivo singolare. |
| Help | Guida | Preferito a "Aiuto" nei menu di navigazione del software. |
| Close | Chiudi | Imperativo. Usato per chiudere finestre o pannelli. |
