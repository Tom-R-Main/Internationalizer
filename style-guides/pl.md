# Polish (pl) — Translation Style Guide

## Language Profile
- Locale: `pl`
- Script: Latin, LTR
- CLDR plural forms: one, few, many, other
- Text expansion vs English: 20-30%

## Tone & Formality
Nowoczesne interfejsy oprogramowania (UI) w języku polskim powinny brzmieć naturalnie, przyjaźnie i profesjonalnie. Unikaj zbytniej formalności i form grzecznościowych typu "Pan/Pani", chyba że tłumaczysz system bankowy, medyczny lub prawny. 

Zwracaj się do użytkownika bezpośrednio, używając drugiej osoby liczby pojedynczej (np. "Zaloguj się", "Twój profil"). Komunikaty o błędach powinny być jasne, pomocne i nieobwiniające użytkownika (zamiast "Wprowadziłeś złe hasło" użyj "Podane hasło jest nieprawidłowe").

## Grammar
*   **Przyciski i akcje (Call to Action):** Używaj trybu rozkazującego w 2. osobie liczby pojedynczej (np. *Zapisz*, *Usuń*, *Wyślij*). Unikaj bezokoliczników (nie pisz *Zapisać*, *Usunąć*).
*   **Etykiety i nagłówki:** Stosuj rzeczowniki odczasownikowe dla procesów (np. *Pobieranie*, *Wyszukiwanie*) oraz formy rzeczownikowe dla sekcji (np. *Ustawienia*, *Konto*).
*   **Wielkie litery (Capitalization):** W języku polskim w tytułach, nagłówkach i na przyciskach wielką literą piszemy tylko pierwsze słowo (tzw. sentence case). Piszemy *Ustawienia konta*, a nie *Ustawienia Konta*.
*   **Neutralność płciowa:** Unikaj form wskazujących na płeć użytkownika (np. *Zrobiłeś to*, *Byłaś zalogowana*). Zamiast *Jesteś pewny?* użyj *Czy na pewno chcesz kontynuować?*. Zamiast *Zostałeś wylogowany* użyj formy bezosobowej: *Wylogowano*.

## Pluralization
Język polski posiada cztery formy liczby mnogiej (według standardu CLDR), które zależą od końcówki liczby:
*   **one (jeden):** Używana wyłącznie dla liczby 1 (np. *1 plik*, *1 błąd*).
*   **few (kilka):** Używana dla liczb kończących się na 2, 3, 4, z wyjątkiem tych kończących się na 12, 13, 14 (np. *2 pliki*, *4 błędy*, *23 pliki*, *104 błędy*).
*   **many (wiele):** Używana dla liczb całkowitych, które nie wpadają w kategorie "one" i "few", w tym dla zera (np. *0 plików*, *5 plików*, *11 błędów*, *26 plików*).
*   **other (inne):** Używana dla ułamków (np. *1,5 pliku*, *0,5 błędu*).

## Punctuation & Typography
*   **Cudzysłowy:** Używaj polskich cudzysłowów typograficznych: „tekst”. W zagnieżdżeniach stosuj cudzysłowy francuskie: «tekst». Nie używaj prostych cudzysłowów ("tekst").
*   **Liczby:** Separatorem dziesiętnym jest przecinek (np. *3,14*). Separatorem tysięcy jest spacja nierozdzielająca (np. *1 000 000*).
*   **Data i czas:** Standardowy format daty to DD.MM.YYYY (np. *31.12.2023*) lub międzynarodowy YYYY-MM-DD. Format czasu to system 24-godzinny (np. *14:30*, unikaj formatu AM/PM).
*   **Interpunkcja:** Nie stawiaj spacji przed znakami zapytania, wykrzyknikami, dwukropkami czy średnikami.

## Terminology

| English | Polish | Notes |
|---------|---------|-------|
| Save | Zapisz | Tryb rozkazujący, standard dla przycisków zapisu. |
| Cancel | Anuluj | Tryb rozkazujący, przerywanie akcji. |
| Delete | Usuń | Tryb rozkazujący (nie "Skasuj"). |
| Settings | Ustawienia | Rzeczownik, zawsze w liczbie mnogiej. |
| Search | Szukaj | Tryb rozkazujący dla przycisku. Jako etykieta pola: "Wyszukiwanie". |
| Error | Błąd | |
| Loading | Wczytywanie | Rzeczownik odczasownikowy. Dopuszczalne również "Ładowanie". |
| Dashboard | Panel główny | Unikaj angielskiego słowa "Dashboard". Można też użyć "Kokpit". |
| Notifications | Powiadomienia | |
| Sign in | Zaloguj się | Tryb rozkazujący. |
| Sign out | Wyloguj się | Tryb rozkazujący. |
| Submit | Prześlij | Zależnie od kontekstu formularza, dopuszczalne też "Zatwierdź". |
| Profile | Profil | |
| Help | Pomoc | |
| Close | Zamknij | Tryb rozkazujący, zamykanie okna/modalu. |
