# Finnish (fi) — Translation Style Guide

## Language Profile
- Locale: `fi`
- Script: Latin, LTR
- CLDR plural forms: one, other
- Text expansion vs English: 20–30 % (suomen kielen sanat ovat usein pidempiä ja teksti tiiviimpää päätteiden vuoksi).

## Tone & Formality
Ohjelmistojen käyttöliittymissä sävyn tulee olla selkeä, asiantunteva ja luonnollinen. Puhuttele käyttäjää yksikön toisessa persoonassa (sinä), mutta vältä pronominin "sinä" toistamista (esim. *Tilisi* on parempi kuin *Sinun tilisi*). Älä käytä teitittelyä (Te) normaaleissa ohjelmistoissa, ellei kyseessä ole erittäin muodollinen palvelu (esim. perinteinen pankki). Käytä aktiivisia ja suoria ilmauksia. Vältä liiallista tuttavallisuutta tai slangia; pyri neutraaliin ja ystävälliseen yleiskieleen.

## Grammar
- **Painikkeet ja toiminnot (Imperatiivi):** Käytä painikkeissa ja valikoissa aina yksikön 2. persoonan imperatiivia (esim. *Tallenna*, *Poista*, *Kopioi*), älä perusmuotoa (*Tallentaa*) tai passiivia.
- **Otsikot ja kentät (Substantiivit):** Käytä käyttöliittymän otsikoissa ja kenttien nimissä substantiiveja (esim. *Asetukset*, *Käyttäjänimi*).
- **Yhdyssanat:** Englannin erilliset sanat kirjoitetaan suomeksi usein yhteen. Esimerkiksi "User settings" on *Käyttäjäasetukset*, ei "Käyttäjä asetukset". Tarkista yhdyssanat huolellisesti.
- **Omistusliitteet:** Vältä englannin "your"-sanan suoraa kääntämistä. Käytä omistusliitettä (esim. *Profiilisi*, ei "Sinun profiilisi") tai jätä se kokonaan pois, jos asiayhteys on selvä (esim. *Omat tiedot*).
- **Isot alkukirjaimet (Sentence case):** Suomessa vain virkkeen tai otsikon ensimmäinen sana ja erisnimet kirjoitetaan isolla alkukirjaimella. Englannin "Account Settings" on suomeksi *Tilin asetukset*, ei "Tilin Asetukset".

## Pluralization
Suomen kielessä on kaksi CLDR-monikkomuotoa: `one` (yksi) ja `other` (muut).
- **one:** Käytetään, kun luku on tasan 1. Esimerkki: "1 tiedosto" (1 file).
- **other:** Käytetään kaikkien muiden lukujen kohdalla, mukaan lukien nolla (0, 2, 3, 4...). Luvun jälkeinen substantiivi on tällöin yksikön partitiivissa. Esimerkki: "0 tiedostoa", "2 tiedostoa", "5 tiedostoa".
Huom: Jos luku puuttuu kokonaan ja puhutaan yleisesti monikosta, käytetään monikon nominatiivia (esim. "Tiedostot").

## Punctuation & Typography
- **Lainausmerkit:** Käytä suomalaisia lainausmerkkejä (”teksti”). Älä käytä englannin kaarevia merkkejä (“ ”).
- **Desimaali- ja tuhaterottimet:** Desimaalierotin on pilkku (esim. 3,14). Tuhaterotin on sitova välilyönti (esim. 1 000 tai 10 000).
- **Päivämäärät:** Käytä muotoa PP.KK.VVVV (esim. 24.12.2023).
- **Kellonajat:** Käytä 24 tunnin kelloa ja erota tunnit ja minuutit pisteellä (esim. 14.30, ei 14:30).

## Terminology

| English | Finnish | Notes |
|---------|---------|-------|
| Save | Tallenna | Käytä painikkeissa aina imperatiivia. |
| Cancel | Peruuta | Toiminnon peruminen (imperatiivi). |
| Delete | Poista | Tiedon tai kohteen poistaminen (imperatiivi). |
| Settings | Asetukset | Monikkomuoto, yleinen käyttöliittymäterminä. |
| Search | Hae / Haku | "Hae" painikkeessa (verbi), "Haku" otsikossa tai kentässä (substantiivi). |
| Error | Virhe | Yleinen virheilmoituksen otsikko. |
| Loading | Ladataan | Passiivin preesens, kun prosessi on käynnissä. |
| Dashboard | Ohjausnäkymä | Myös "Kojelauta" tai "Etusivu" kontekstista riippuen. |
| Notifications | Ilmoitukset | Monikon nominatiivi. |
| Sign in | Kirjaudu sisään | Myös pelkkä "Kirjaudu" on usein riittävä. |
| Sign out | Kirjaudu ulos | Istunnon päättäminen. |
| Submit | Lähetä | Lomakkeen tietojen lähettäminen. |
| Profile | Profiili | Käyttäjän omat tiedot. |
| Help | Ohje | Yksikkömuoto, ei "Ohjeet" ellei viitata useisiin dokumentteihin. |
| Close | Sulje | Ikkunan tai valikon sulkeminen (imperatiivi). |
