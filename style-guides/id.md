# Indonesian (id) — Translation Style Guide

## Language Profile
- Locale: `id`
- Script: Latin, LTR
- CLDR plural forms: other
- Text expansion vs English: 15-25%

## Tone & Formality
Gunakan nada yang profesional, sopan, namun tetap ramah dan mudah dipahami. Hindari bahasa gaul, jargon yang tidak perlu, atau gaya bahasa yang terlalu kaku dan birokratis. 

Gunakan sapaan "Anda" untuk merujuk pada pengguna (jangan pernah menggunakan "kamu", "kalian", "saudara", atau "anda" dengan huruf kecil). Untuk antarmuka pengguna (UI), gunakan kalimat langsung, ringkas, dan berorientasi pada tindakan. Utamakan bentuk aktif karena lebih mudah dibaca dan dipahami dengan cepat dibandingkan bentuk pasif.

## Grammar
1. **Kata Kerja Perintah (Imperatif) untuk Tombol:** Gunakan kata kerja dasar tanpa imbuhan awalan untuk tombol, tautan, atau tindakan. Contoh: Gunakan "Simpan" (bukan "Menyimpan"), "Hapus" (bukan "Menghapus"), "Bagikan" (bukan "Membagikan").
2. **Hukum DM (Diterangkan-Menerangkan):** Bahasa Indonesia menempatkan kata benda utama di depan kata sifat atau penjelasnya. Jangan menerjemahkan frasa bahasa Inggris secara harfiah. Contoh: "Network Settings" menjadi "Pengaturan Jaringan" (bukan "Jaringan Pengaturan").
3. **Frasa Kata Benda untuk Label:** Label UI, judul menu, atau tab harus berupa frasa kata benda. Gunakan imbuhan *pe-an* atau *per-an* untuk menunjukkan proses atau kumpulan. Contoh: "Settings" -> "Pengaturan", "Search" (sebagai fitur) -> "Pencarian".
4. **Pesan Kesalahan (Error Messages):** Gunakan kalimat yang jelas, objektif, dan berikan solusi. Hindari menyalahkan pengguna. Jangan menerjemahkan kata seru seperti "Oops!" atau "Sorry" secara harfiah kecuali konteks aplikasinya sangat kasual. Contoh: "Invalid password" diterjemahkan menjadi "Kata sandi tidak valid" (bukan "Maaf, kata sandi Anda salah").
5. **Hindari Redundansi Jamak:** Jangan mengulang kata benda untuk menunjukkan bentuk jamak jika sudah didahului oleh angka atau kata penunjuk kuantitas (seperti "beberapa", "semua"). Contoh: "5 files" diterjemahkan menjadi "5 berkas" (BUKAN "5 berkas-berkas").

## Pluralization
Bahasa Indonesia tidak memiliki perubahan bentuk kata secara tata bahasa untuk bentuk jamak berdasarkan angka (CLDR hanya menggunakan kategori `other`). 

Bentuk jamak umumnya ditunjukkan dengan pengulangan kata (contoh: "berkas-berkas") atau dengan menambahkan kata bantu bilangan. Namun, dalam konteks UI perangkat lunak, jika sudah ada angka yang mendahului, kata benda **TIDAK BOLEH** diulang.
- 0 items: "0 berkas"
- 1 item: "1 berkas"
- 2 items: "2 berkas" (BUKAN "2 berkas-berkas")
- 100 items: "100 berkas"

## Punctuation & Typography
- **Tanda Kutip:** Gunakan tanda kutip ganda ("...") untuk kutipan utama dan tanda kutip tunggal ('...') untuk kutipan di dalam kutipan.
- **Pemisah Angka:** Gunakan koma (,) sebagai pemisah desimal dan titik (.) sebagai pemisah ribuan. Contoh: 1.234.567,89 (Satu juta dua ratus tiga puluh empat ribu lima ratus enam puluh tujuh koma delapan sembilan).
- **Format Tanggal:** Gunakan format DD/MM/YYYY atau DD MMMM YYYY. Contoh: 31/12/2023 atau 31 Desember 2023.
- **Format Waktu:** Gunakan format 24 jam dengan pemisah titik (.), bukan titik dua (:). Contoh: 14.30 WIB.
- **Kapitalisasi:** Gunakan *Sentence case* (huruf kapital hanya di awal kalimat) untuk deskripsi, pesan kesalahan, dan teks panjang. Gunakan *Title Case* untuk judul halaman, tab, dan tombol utama, dengan catatan kata depan/sambung (dan, di, ke, dari, untuk) tetap huruf kecil. Contoh: "Syarat dan Ketentuan".

## Terminology

| English | Indonesian | Notes |
|---------|------------|-------|
| Save | Simpan | Gunakan kata dasar untuk tombol aksi. |
| Cancel | Batal | Singkat dan umum digunakan di UI. |
| Delete | Hapus | Gunakan kata dasar untuk tombol aksi. |
| Settings | Pengaturan | Gunakan bentuk kata benda untuk menu/label. |
| Search | Cari | Untuk tombol/placeholder. Gunakan "Pencarian" untuk judul halaman/fitur. |
| Error | Galat | Bisa juga "Kesalahan", namun "Galat" lebih baku untuk sistem IT. |
| Loading | Memuat | Bentuk aktif yang menunjukkan proses sedang berjalan. |
| Dashboard | Dasbor | Istilah serapan baku sesuai KBBI. |
| Notifications | Notifikasi | Istilah serapan yang lebih umum dan modern dibanding "Pemberitahuan". |
| Sign in | Masuk | Standar industri untuk autentikasi (login). |
| Sign out | Keluar | Standar industri untuk autentikasi (logout). |
| Submit | Kirim | Lebih natural dan umum digunakan daripada "Serahkan". |
| Profile | Profil | Istilah serapan baku. |
| Help | Bantuan | Gunakan bentuk kata benda untuk menu atau tautan. |
| Close | Tutup | Gunakan kata dasar untuk tombol aksi/menutup jendela. |
