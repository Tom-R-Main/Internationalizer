# Malay (ms) — Translation Style Guide

## Language Profile
- Locale: `ms`
- Script: Latin, LTR
- CLDR plural forms: other
- Text expansion vs English: 15-25%

## Tone & Formality
Gunakan nada yang profesional, mesra, dan jelas. Untuk antara muka perisian (UI), bahasa haruslah ringkas, tepat, dan terus kepada intipati. Gunakan kata ganti nama "anda" untuk merujuk kepada pengguna secara sopan. Elakkan bahasa slanga, dialek daerah, atau singkatan tidak rasmi (seperti "yg", "dgn", "tak"). Walaupun profesional, elakkan laras bahasa istana atau sastera yang terlalu berbunga dan sukar difahami. Untuk arahan, gunakan ayat aktif yang jelas agar pengguna tahu tindakan yang perlu diambil.

## Grammar
- **Hukum D-M (Diterangkan - Menerangkan):** Kata nama utama mesti mendahului penerang (kata sifat atau kata nama lain). Contoh: "User settings" diterjemahkan sebagai "Tetapan pengguna", bukan "Pengguna tetapan".
- **Butang dan Arahan (Imperatif):** Gunakan kata kerja dasar tanpa imbuhan awalan "me-" untuk butang atau tindakan. Contoh: Gunakan "Simpan" (bukan "Menyimpan"), "Cari" (bukan "Mencari"), dan "Hantar" (bukan "Menghantar").
- **Mesej Ralat:** Gunakan ayat pasif atau struktur yang objektif tanpa menyalahkan pengguna. Contoh: "Fail tidak dapat dijumpai" adalah lebih baik daripada "Anda gagal menjumpai fail".
- **Penggunaan "Sila":** Gunakan perkataan "Sila" untuk arahan yang memerlukan kesopanan (contoh: "Sila log masuk untuk meneruskan"), tetapi abaikan perkataan ini pada label butang UI yang sempit untuk menjimatkan ruang.

## Pluralization
Bahasa Melayu hanya mempunyai satu kategori plural CLDR iaitu `other`. Peraturan paling penting: kata nama tidak perlu digandakan jika nombor atau kuantiti telah dinyatakan secara spesifik.
- **Contoh salah:** "3 fail-fail", "5 mesej-mesej"
- **Contoh betul:** "1 fail", "3 fail", "5 mesej"
Penggandaan kata nama (contoh: "fail-fail", "pengguna-pengguna") hanya digunakan untuk menunjukkan jamak apabila tiada angka dinyatakan. Dalam UI, format pembolehubah seperti "{count} fail" adalah standard dan tepat untuk semua nilai angka.

## Punctuation & Typography
- **Format Nombor:** Bahasa Melayu di Malaysia menggunakan titik (.) untuk perpuluhan dan koma (,) untuk pemisah ribuan. Contoh: 1,234,567.89 (Sama seperti format Inggeris AS/UK).
- **Format Tarikh:** Gunakan format HH/BB/TTTT (contoh: 31/12/2023) untuk format ringkas, atau format panjang seperti "31 Disember 2023".
- **Format Masa:** Sistem 12 jam (contoh: 2:30 petang, 9:00 pagi) atau sistem 24 jam (contoh: 14:30) boleh diterima bergantung pada konteks aplikasi, tetapi pastikan ia konsisten.
- **Tanda Petik:** Gunakan tanda petik berganda ("...") untuk petikan utama dan tanda petik tunggal ('...') untuk petikan di dalam petikan.
- **Huruf Besar (Capitalization):** Gunakan "Sentence case" (Huruf besar pada awal ayat sahaja) untuk mesej ralat dan penerangan panjang. Untuk tajuk menu dan butang, gunakan "Title Case" (Huruf besar pada setiap perkataan utama).

## Terminology

| English | Malay | Notes |
|---------|---------|-------|
| Save | Simpan | Gunakan kata dasar untuk butang tindakan. |
| Cancel | Batal | Digunakan untuk membatalkan tindakan atau proses. |
| Delete | Padam | Boleh juga menggunakan "Hapus", tetapi "Padam" lebih lazim untuk UI. |
| Settings | Tetapan | Kata nama untuk menu konfigurasi sistem atau aplikasi. |
| Search | Cari | Gunakan "Cari" untuk butang/ruangan teks, bukan "Carian" (kata nama). |
| Error | Ralat | Digunakan untuk mesej kesilapan sistem atau pepijat. |
| Loading | Memuatkan | Menunjukkan proses sedang berjalan. |
| Dashboard | Papan Pemuka | Istilah standard untuk skrin utama analitik atau kawalan. |
| Notifications | Pemberitahuan | Boleh juga menggunakan "Notifikasi" jika ruang UI sangat terhad. |
| Sign in | Log masuk | Dua perkataan. Huruf kecil 'm' jika berada di tengah ayat. |
| Sign out | Log keluar | Dua perkataan. Huruf kecil 'k' jika berada di tengah ayat. |
| Submit | Hantar | Digunakan untuk tindakan menghantar borang atau data. |
| Profile | Profil | Ejaan bahasa Melayu tanpa huruf 'e' di hujung. |
| Help | Bantuan | Kata nama untuk menu sokongan atau panduan pengguna. |
| Close | Tutup | Digunakan untuk tindakan menutup tetingkap, tab, atau dialog. |
