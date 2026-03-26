# Turkish (tr) — Translation Style Guide

## Language Profile
- Locale: `tr`
- Script: Latin, LTR
- CLDR plural forms: one, other
- Text expansion vs English: ~10-20%

## Tone & Formality
Yazılım arayüzlerinde (UI) profesyonel, anlaşılır ve saygılı bir dil kullanılmalıdır. Kullanıcıya hitap ederken genellikle ikinci çoğul şahıs ("siz") veya edilgen yapı tercih edilmelidir. Ancak, çok resmi veya soğuk bir tondan kaçınılmalı, modern ve kullanıcı dostu bir üslup benimsenmelidir. Metinlerde geniş zaman veya mastar ekleri (-mek/-mak) kullanılabilir. Butonlar ve kısa menü komutlarında ise kısa ve net olması amacıyla ikinci tekil şahıs emir kipleri (örn. "Kaydet", "Sil") standarttır.

## Grammar
- **Butonlar ve Kısa Eylemler:** UI butonlarında ve menü öğelerinde her zaman ikinci tekil şahıs emir kipi kullanın (örn. *Save* -> *Kaydet*, *Kaydedin* DEĞİL).
- **Hata Mesajları ve Bildirimler:** Edilgen yapı veya geniş zaman kullanın. Kullanıcıyı suçlayıcı bir dilden kaçının (örn. *You entered the wrong password* yerine *Hatalı şifre girildi* veya *Şifre hatalı*).
- **İsim Tamlamaları:** İngilizcedeki zincirleme isim tamlamalarını Türkçeye çevirirken iyelik eklerine dikkat edin ve anlam karmaşası yaratmamaya özen gösterin (örn. *User Account Settings* -> *Kullanıcı Hesabı Ayarları*).
- **Büyük/Küçük Harf Kullanımı (Capitalization):** İngilizcede başlıkların her kelimesi büyük harfle başlasa da (Title Case), Türkçede genellikle sadece ilk kelimenin ilk harfi büyük yazılır (Sentence case). Arayüz metinlerinde doğal bir görünüm için bu kurala uyulmalıdır.

## Pluralization
Türkçede CLDR çoğul kuralları `one` ve `other` olmak üzere ikiye ayrılır.
- **one:** Miktar belirtmeyen tekil durumlar için kullanılır.
- **other:** Sıfır ve birden büyük tüm sayılar için kullanılır. 
**ÖNEMLİ KURAL:** Türkçede sayı belirtildiğinde isim her zaman tekil kalır. İngilizcedeki gibi isme çoğul eki (-ler/-lar) eklenmez. 
*Doğru:* 0 dosya, 1 dosya, 3 dosya, 10 dosya.
*Yanlış:* 3 dosyalar, 10 dosyalar.
Çoğul eki sadece sayı belirtilmediğinde kullanılır (örn. *Dosyalar silindi*).

## Punctuation & Typography
- **Tırnak İşaretleri:** Standart çift tırnak ("...") kullanılır.
- **Sayılar:** Binlik ayracı olarak nokta (.), ondalık ayracı olarak virgül (,) kullanılır (örn. *1.234.567,89*).
- **Tarih Formatı:** GG.AA.YYYY veya GG/AA/YYYY formatı standarttır (örn. *31.12.2023*).
- **Saat Formatı:** 24 saatlik dilim kullanılır. Saat ile dakika arasına iki nokta (:) veya nokta (.) konulabilir (örn. *14:30*).
- **Kısaltmalar:** İngilizce kısaltmaların sonuna gelen ekler, kısaltmanın Türkçe okunuşuna göre kesme işaretiyle ayrılır (örn. *API'ye*, *URL'yi*).

## Terminology

| English | Turkish | Notes |
|---------|---------|-------|
| Save | Kaydet | Butonlarda ikinci tekil şahıs emir kipi. |
| Cancel | İptal | Genellikle sadece "İptal" yeterlidir. |
| Delete | Sil | Butonlarda kısa ve net emir kipi. |
| Settings | Ayarlar | Menü ve başlıklarda standart kullanım. |
| Search | Ara | Arama kutusu veya buton için. |
| Error | Hata | Genel hata bildirimleri için. |
| Loading | Yükleniyor | Devam eden eylem durumu. |
| Dashboard | Kontrol Paneli | "Pano" da kullanılabilir ancak "Kontrol Paneli" daha yaygındır. |
| Notifications | Bildirimler | - |
| Sign in | Giriş yap | "Oturum aç" veya sadece "Giriş" de bağlama göre uygundur. |
| Sign out | Çıkış yap | "Oturumu kapat" veya sadece "Çıkış" da bağlama göre uygundur. |
| Submit | Gönder | Form onaylama butonları için. |
| Profile | Profil | - |
| Help | Yardım | - |
| Close | Kapat | Pencere veya diyalog kapatma eylemi. |
