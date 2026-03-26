# Romanian (ro) — Translation Style Guide

## Language Profile
- Locale: `ro`
- Script: Latin, LTR
- CLDR plural forms: one, few, other
- Text expansion vs English: 15-25%

## Tone & Formality
Folosiți un ton profesional, clar, politicos și accesibil. Pentru interfețele software standard, adresați-vă utilizatorului la persoana a II-a plural (pronumele de politețe „dumneavoastră” este subînțeles, dar evitați folosirea lui explicită și excesivă pentru a nu încărca textul). Folosiți diateza activă și formulări directe. Evitați jargonul tehnic excesiv și traducerile literale (mot-à-mot). Limbajul trebuie să sune natural în limba română, fiind în același timp respectuos și de încredere.

## Grammar
- **Diacritice (Critic):** Folosiți exclusiv diacriticele corecte cu virgulă dedesubt (ș, ț, Ș, Ț). Nu folosiți sub nicio formă diacriticele vechi cu sedilă (ş, ţ). Aceasta este o regulă strictă pentru orice text modern în limba română.
- **Majuscule (Capitalization):** Spre deosebire de limba engleză care folosește „Title Case”, în limba română se folosește scrierea cu literă mare doar pentru primul cuvânt din titluri, meniuri, butoane sau etichete (Sentence case), cu excepția numelor proprii. (Corect: *Setări de confidențialitate*, Greșit: *Setări De Confidențialitate*).
- **Butoane și acțiuni:** Pentru butoane (Call-to-Action), folosiți verbe la modul imperativ, persoana a II-a plural (ex. *Salvați*, *Ștergeți*) sau substantive de acțiune (ex. *Salvare*, *Ștergere*). Păstrați consecvența în întreaga interfață.
- **Diateza pasivă vs. reflexivă:** Evitați folosirea excesivă a diatezei pasive din engleză („File was saved”). Preferați diateza reflexiv-pasivă („Fișierul s-a salvat”) sau formulările impersonale/active („Ați salvat fișierul” sau „Fișier salvat”).

## Pluralization
Limba română folosește trei forme de plural conform standardului CLDR. Este esențială folosirea prepoziției „de” pentru categoria `other`.
- **one (un/o):** Se aplică pentru numărul 1 (ex. *1 fișier*, *1 utilizator*).
- **few (câteva):** Se aplică pentru 0 și numerele care se termină în 01-19 (ex. *0 fișiere*, *2 fișiere*, *14 fișiere*, *115 fișiere*). Substantivul se leagă direct de număr.
- **other (restul/de):** Se aplică pentru numerele care se termină în 20-99 sau 00, cu excepția lui 0 (ex. *20 de fișiere*, *100 de fișiere*, *1.000 de fișiere*). **Regulă critică:** Necesită adăugarea prepoziției „de” între număr și substantiv.

## Punctuation & Typography
- **Ghilimele:** Folosiți ghilimelele românești: „jos-sus” pentru citate/evidențieri principale și «franțuzești» pentru citate în interiorul altor citate (ex. *„Apasă pe «Salvează» acum”*).
- **Numere:** Folosiți virgula (,) pentru separatorul zecimal și punctul (.) pentru gruparea miilor (ex. *1.234.567,89*).
- **Data și ora:** Formatul datei este ZZ.LL.AAAA sau ZZ/LL/AAAA. Formatul orei este de 24 de ore (ex. *14:30*, nu *2:30 PM*).
- **Spațiere:** Nu lăsați spațiu înainte de semnele de punctuație duble (: ; ! ?).

## Terminology

| English | Romanian | Notes |
|---------|----------|-------|
| Save | Salvați | Verb la imperativ, persoana a II-a plural. |
| Cancel | Anulați | Verb la imperativ, persoana a II-a plural. |
| Delete | Ștergeți | Verb la imperativ, persoana a II-a plural. |
| Settings | Setări | Substantiv, plural. |
| Search | Căutare | Substantiv. Folosit pentru câmpuri de căutare și etichete. |
| Error | Eroare | Substantiv, singular. |
| Loading | Se încarcă | Verb reflexiv. Evitați „Încărcare” dacă indică o acțiune în curs de desfășurare. |
| Dashboard | Panou de control | Traducerea standard și clară pentru interfețe. |
| Notifications | Notificări | Substantiv, plural. |
| Sign in | Conectare | Substantiv de acțiune. Preferat în loc de „Autentificare” sau „Logare”. |
| Sign out | Deconectare | Substantiv de acțiune. |
| Submit | Trimiteți | Verb la imperativ. Se folosește pentru formulare. |
| Profile | Profil | Substantiv, singular. |
| Help | Ajutor | Substantiv, singular. |
| Close | Închideți | Verb la imperativ, persoana a II-a plural. |
