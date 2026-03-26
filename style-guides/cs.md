# Czech (cs) — Translation Style Guide

## Language Profile
- Locale: `cs`
- Script: Latin, LTR
- CLDR plural forms: one, few, many, other
- Text expansion vs English: 15-20%

## Tone & Formality
Pro uživatelské rozhraní (UI) a softwarovou dokumentaci používejte vždy formální oslovení (vykání), pokud není výslovně požadován neformální tón. Tón by měl být profesionální, jasný, přátelský a nápomocný. Vyhněte se příliš strohému, robotickému nebo naopak příliš familiárnímu vyjadřování. Preferujte aktivní rod před trpným (např. „Systém uložil soubor“ zní lépe než „Soubor byl uložen systémem“). Text musí působit přirozeně a srozumitelně pro běžného uživatele.

## Grammar
*   **Tlačítka a akce (Infinitiv):** Pro popis akcí na tlačítkách, v menu a odkazech používejte vždy infinitiv (např. „Uložit“, „Smazat“, nikoliv „Uložte“ nebo „Smažte“).
*   **Pokyny pro uživatele (Rozkazovací způsob):** Při přímých pokynech v textu nebo chybových hláškách používejte rozkazovací způsob ve 2. osobě množného čísla (vykání) (např. „Zadejte heslo“, „Vyberte soubor pro nahrání“).
*   **Velká písmena (Capitalization):** Čeština nepoužívá anglický „Title Case“. V nadpisech, názvech oken, menu a popiscích pište velké pouze první písmeno prvního slova a vlastní jména (např. „Nastavení uživatelského účtu“, nikoliv „Nastavení Uživatelského Účtu“).
*   **Přivlastňovací zájmena:** Preferujte zvratné přivlastňovací zájmeno „svůj“ před „váš“, pokud se vztahuje k podmětu věty (např. „Zadejte své heslo“, nikoliv „Zadejte vaše heslo“). Anglické „your“ lze v češtině často zcela vynechat (např. „Open your profile“ -> „Otevřít profil“).
*   **Shoda přísudku s podmětem:** U zástupných proměnných (např. „{User} přidal...“) volte neutrální formulace nebo se vyhněte minulému času, pokud není znám rod uživatele, abyste předešli gramatickým chybám (např. „Uživatel {User} přidal...“).

## Pluralization
Čeština používá čtyři kategorie plurálu podle standardu CLDR. Je nutné dodržovat správné skloňování podstatných jmen:
*   **one (jeden):** Pro číslo 1 (např. „1 soubor“, „1 položka“).
*   **few (několik):** Pro čísla 2, 3 a 4 (např. „2 soubory“, „4 položky“).
*   **many (mnoho):** Pro desetinná čísla (např. „1,5 souboru“). V běžném UI se používá zřídka.
*   **other (ostatní):** Pro číslo 0 a všechna celá čísla od 5 výše (např. „0 souborů“, „5 souborů“, „100 položek“).

## Punctuation & Typography
*   **Uvozovky:** Používejte výhradně české typografické uvozovky: dolní na začátku a horní na konci („text“). Nepoužívejte anglické ("text").
*   **Čísla:** Jako oddělovač desetinných míst se používá čárka (3,14). Jako oddělovač tisíců se používá pevná mezera (1 000 000).
*   **Datum:** Standardní formát je DD. MM. YYYY (např. 31. 12. 2023). Za tečkou vždy následuje mezera.
*   **Čas:** Používá se 24hodinový formát. Oddělovačem hodin a minut je v UI obvykle dvojtečka (14:30).
*   **Zkratky a symboly:** Mezi číslem a symbolem procenta nebo měny se píše pevná mezera (např. „100 %“, „50 Kč“). Výjimkou jsou přídavná jména (např. „100%“ bez mezery znamená „stoprocentní“).

## Terminology

| English | Czech | Notes |
|---------|---------|-------|
| Save | Uložit | Infinitiv pro tlačítka a akce. |
| Cancel | Zrušit | Používá se pro přerušení akce nebo zavření dialogu bez uložení. |
| Delete | Smazat | Lze použít i „Odstranit“, ale „Smazat“ je pro běžné UI častější a kratší. |
| Settings | Nastavení | Používá se pro sekci konfigurace aplikace. |
| Search | Hledat | Infinitiv pro tlačítko nebo zástupný text (placeholder) ve vyhledávacím poli. |
| Error | Chyba | Obecný termín pro chybové stavy. |
| Loading | Načítání | Používá se jako indikátor průběhu (progress state). |
| Dashboard | Ovládací panel | Lze použít i „Nástěnka“ u méně formálních aplikací, ale „Ovládací panel“ je standard. |
| Notifications | Oznámení | Nikdy nepřekládat jako „Notifikace“, pokud to není nezbytně nutné z prostorových důvodů. |
| Sign in | Přihlásit se | Infinitiv pro tlačítko vstupu do účtu. |
| Sign out | Odhlásit se | Infinitiv pro tlačítko opuštění účtu. |
| Submit | Odeslat | Používá se pro odeslání formulářů nebo dat. |
| Profile | Profil | Osobní nastavení a údaje uživatele. |
| Help | Nápověda | Odkaz na dokumentaci nebo podporu. |
| Close | Zavřít | Infinitiv pro zavření okna, modálu nebo panelu. |
