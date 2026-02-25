# PRD: Quran API Go (MVP)

---

## Latar Belakang & Rumusan Masalah

Super app Ilmunara membutuhkan akses terstruktur ke data Al-Quran untuk berbagai fitur inti. Saat ini belum ada internal data service — setiap fitur super app harus mengelola data Quran secara mandiri, yang berpotensi menimbulkan duplikasi dan inkonsistensi.

Proyek ini membangun sebuah layanan API internal yang ringan dan terpusat sebagai sumber data Al-Quran untuk super app.

**Ini adalah layanan internal, bukan utilitas publik.** Akses publik dapat dipertimbangkan di fase berikutnya.

---

## Tujuan

- Menyediakan API internal yang cepat dan ringan untuk data Al-Quran
- Mendukung terjemahan Bahasa Indonesia & Inggris
- Navigasi inti: surah, ayat, juz
- Pencarian teks dengan filter dasar
- Zero dependensi infrastruktur eksternal di MVP

## Non-Goals (Fase Berikutnya)

- Akses publik / open API
- Autentikasi & API key
- Rate limiting (tidak ada traffic publik di MVP)
- Endpoint navigasi Hizb, Rub el Hizb, Manzil, Ruku, Halaman
- Audio murattal, Tafsir, Tajweed, Morfologi kata
- Redis caching layer, GraphQL
- Bahasa selain ID/EN
- Admin panel

---

## Keputusan Teknis

| Keputusan | Pilihan | Alasan |
|-----------|---------|--------|
| Bahasa & Framework | Go + Gin | Ringan, cepat, cocok untuk API |
| Database | SQLite (FTS5) | Data statis read-only; tidak perlu infrastruktur server database |
| Dependency Injection | Manual constructor injection | Ramah relawan; tidak perlu code generation |
| Rate Limiting | Tidak ada (MVP) | Traffic internal saja; tambahkan di fase berikutnya |
| Autentikasi | Tidak ada (MVP) | Layanan internal; tambahkan di fase berikutnya |
| Dokumentasi | Scalar di /docs | Referensi developer internal |
| Deployment | [VPS / Railway / Render — TBD] | [Pemilik: TBD] |

---

## Quality Gates

Wajib lolos sebelum setiap issue di-close:

```
go test ./...     # Unit tests
go vet ./...      # Static analysis
gofmt -d .        # Code formatting
```

---

## Sprint 0 — Fondasi & Data
> Tujuan: Project dapat berjalan secara lokal dengan data Quran yang lengkap dan tervalidasi.
>
> ⚠️ Seluruh issue di Sprint 0 harus selesai sebelum Sprint 1 dimulai.

---

### #1 — Setup project & infrastruktur dasar

```
Labels: sprint-0, setup
Depends on: -
```

**Acceptance Criteria:**
- [ ] Inisialisasi Go module dengan `go mod init`
- [ ] Setup Gin sebagai HTTP router
- [ ] Setup koneksi SQLite menggunakan `modernc.org/sqlite` (pure Go, tanpa CGO)
- [ ] Manual constructor injection di `cmd/api/main.go` — tanpa DI framework
- [ ] Buat Makefile dengan target: `run`, `test`, `lint`, `migrate`, `seed`
- [ ] Buat `.env.example` (lihat bagian Environment Variables)

---

### #2 — Skema database & migrasi

```
Labels: sprint-0, database
Depends on: #1
```

**Acceptance Criteria:**
- [ ] Setup Goose untuk migrasi (kompatibel SQLite)
- [ ] Migrasi tabel `surahs`: id, number, name_arabic, name_latin, name_transliteration, number_of_ayahs, revelation_type
- [ ] Migrasi tabel `ayahs`: id, surah_id, number_in_surah, text_uthmani, translation_indo, translation_en, juz_number, sajda_type, revelation_type
- [ ] Migrasi tabel `juzs`: id, juz_number, first_ayah_id, last_ayah_id
- [ ] Migrasi FTS5 virtual table untuk pencarian pada `text_uthmani`, `translation_indo`, `translation_en`
- [ ] Indexes: surah_id, juz_number
- [ ] Semua migrasi reversible (up/down)

---

### #3 — Seed data Quran

```
Labels: sprint-0, database
Depends on: #2
```

**Acceptance Criteria:**
- [ ] Import 114 surah
- [ ] Import 6.236 ayat (teks, terjemahan, referensi juz)
- [ ] Import 30 juz
- [ ] Seeder bersifat idempotent (aman dijalankan berulang kali)
- [ ] Seeder menampilkan log progres
- [ ] Validasi pasca-seed: assert jumlah data sesuai (114 / 6.236 / 30)

---

### #4 — Setup logging middleware

```
Labels: sprint-0, middleware
Depends on: #1
```

**Acceptance Criteria:**
- [ ] Setup zerolog sebagai logger
- [ ] Logging middleware mencatat: method, path, status, durasi, IP
- [ ] Tidak mencatat data sensitif
- [ ] Log menggunakan format structured (JSON)

---

## Sprint 1 — Endpoint Surah & Ayat
> Tujuan: Fitur utama super app yang menampilkan konten Quran per surah dan ayat dapat terpenuhi.
>
> Issue #5 harus selesai lebih dulu. Setelah itu, #6, #7, dan #8 dapat dikerjakan **secara paralel**. Issue #9–#12 dapat dimulai segera setelah issue dependency-nya selesai.

---

### #5 — Setup shared DTOs & error response structs

```
Labels: sprint-1, infrastructure
Depends on: #1
```

**Deskripsi:**
Buat shared struct untuk response dan error agar seluruh endpoint menggunakan format yang konsisten.

**Acceptance Criteria:**
- [ ] Buat `pkg/response` dengan helper: `Success()`, `Error()`, `NotFound()`, `BadRequest()`
- [ ] Format error konsisten: `{ error: string, code: string, timestamp: string }`
- [ ] Format success konsisten: `{ data: any, timestamp: string }`
- [ ] Unit test untuk setiap helper function

---

### #6 — Surah repository layer

```
Labels: sprint-1, repository
Depends on: #3
```

**Deskripsi:**
Buat layer repository untuk mengakses data surah dari SQLite.

**Acceptance Criteria:**
- [ ] Interface `SurahRepository` dengan method: `FindAll()`, `FindByID(id)`
- [ ] Implementasi SQLite untuk interface tersebut
- [ ] Unit test dengan SQLite in-memory database

---

### #7 — Ayah repository layer

```
Labels: sprint-1, repository
Depends on: #3
```

**Deskripsi:**
Buat layer repository untuk mengakses data ayat dari SQLite.

**Acceptance Criteria:**
- [ ] Interface `AyahRepository` dengan method: `FindByID(id)`, `FindBySurah(surahID, from, to)`, `FindBySurahAndNumber(surahID, number)`
- [ ] Implementasi SQLite untuk interface tersebut
- [ ] Unit test dengan SQLite in-memory database

---

### #8 — GET /surah & GET /surah/:id — List dan detail surah

```
Labels: sprint-1, endpoint
Depends on: #5, #6
```

**Deskripsi:**
Implementasi dua endpoint surah dasar untuk menampilkan daftar dan detail surah.

**Acceptance Criteria:**
- [ ] `GET /surah` mengembalikan array seluruh surah
- [ ] Response: `[{ id, number, name_arabic, name_latin, name_transliteration, number_of_ayahs, revelation_type }]`
- [ ] `GET /surah/:id` mengembalikan metadata surah tanpa ayat
- [ ] HTTP 404 jika surah tidak ditemukan
- [ ] HTTP 500 jika database error

---

### #9 — GET /surah/:id/ayat — Ayat-ayat dalam surah

```
Labels: sprint-1, endpoint
Depends on: #5, #7
```

**Acceptance Criteria:**
- [ ] Mengembalikan semua ayat dalam surah
- [ ] Query param: `?lang=id` atau `?lang=en` (default: id)
- [ ] Query param: `?from=1&to=10` untuk range (opsional)
- [ ] Response: `{ surah: { id, number, name_latin }, ayahs: [{ number, number_in_surah, text_uthmani, translation, juz, sajda }] }`
- [ ] HTTP 400 jika `lang` tidak valid
- [ ] HTTP 404 jika surah tidak ditemukan

---

### #10 — GET /surah/:id/ayat/:number — Ayat spesifik dalam surah

```
Labels: sprint-1, endpoint
Depends on: #5, #7
```

**Acceptance Criteria:**
- [ ] Mengembalikan satu ayat berdasarkan surah + nomor ayat
- [ ] Query param: `?lang=id` atau `?lang=en`
- [ ] Response: `{ id, surah_id, number, number_in_surah, text_uthmani, translation, surah_info: { id, name_latin }, juz, sajda, revelation_type }`
- [ ] HTTP 400 jika `lang` tidak valid
- [ ] HTTP 404 jika tidak ditemukan

---

### #11 — GET /ayah/:id — Detail ayat by global ID

```
Labels: sprint-1, endpoint
Depends on: #5, #7
```

**Acceptance Criteria:**
- [ ] Mengembalikan ayat berdasarkan global ID (1–6.236)
- [ ] Query param: `?lang=id` atau `?lang=en`
- [ ] Struktur response sama seperti #10
- [ ] HTTP 400 jika `lang` tidak valid
- [ ] HTTP 404 jika tidak ditemukan

---

### #12 — Input validation middleware

```
Labels: sprint-1, middleware
Depends on: #1
```

**Deskripsi:**
Buat middleware/helper terpusat untuk validasi query parameter yang digunakan lintas endpoint.

**Acceptance Criteria:**
- [ ] Helper `ValidateLang(lang string)` — hanya menerima `id` atau `en`
- [ ] Helper `ValidateIDParam(id string)` — harus angka positif
- [ ] Helper `ValidateRangeParam(from, to string)` — from ≤ to, keduanya positif
- [ ] Unit test untuk setiap validator
- [ ] Digunakan oleh endpoint #9, #10, #11

---

## Sprint 2 — Juz, Pencarian & Utilitas
> Tujuan: Navigasi juz, pencarian teks, dan widget ayat harian tersedia untuk super app.
>
> Issue #13 dan #14 dapat dikerjakan **secara paralel** sejak awal sprint. Issue #15 dan #16 dapat dikerjakan paralel setelah #13 selesai. Issue #17 dapat dikerjakan paralel setelah #14 selesai. Issue #18 dan #19 dapat dikerjakan kapan saja dalam sprint ini.

---

### #13 — Juz repository layer

```
Labels: sprint-2, repository
Depends on: #3
```

**Acceptance Criteria:**
- [ ] Interface `JuzRepository` dengan method: `FindAll()`, `FindByNumber(number)`, `FindAyahsByJuz(juzNumber)`
- [ ] Implementasi SQLite untuk interface tersebut
- [ ] Unit test dengan SQLite in-memory database

---

### #14 — Search repository layer (FTS5)

```
Labels: sprint-2, repository
Depends on: #3
```

**Deskripsi:**
Buat layer repository khusus untuk pencarian full-text menggunakan SQLite FTS5.

**Acceptance Criteria:**
- [ ] Interface `SearchRepository` dengan method: `Search(query, filters, pagination)`
- [ ] Implementasi menggunakan FTS5 virtual table
- [ ] Mendukung filter: `surah_id`, `juz`
- [ ] Mendukung pencarian pada `text_uthmani`, `translation_indo`, `translation_en`
- [ ] Case-insensitive dan partial match
- [ ] Unit test dengan SQLite in-memory database

---

### #15 — GET /juz & GET /juz/:number — List dan detail juz

```
Labels: sprint-2, endpoint
Depends on: #5, #13
```

**Acceptance Criteria:**
- [ ] `GET /juz` mengembalikan semua 30 juz
- [ ] Response list: `[{ juz_number, first_ayah_id, last_ayah_id, total_ayahs }]`
- [ ] `GET /juz/:number` mengembalikan semua ayat dalam juz
- [ ] Query param: `?lang=id` atau `?lang=en`
- [ ] Query param: `?page=1&limit=20` untuk pagination
- [ ] Response detail: `{ juz: { juz_number }, ayahs: [{ number, surah_info, number_in_surah, text_uthmani, translation }] }`
- [ ] HTTP 404 jika juz tidak dalam rentang 1–30

---

### #16 — GET /search — Pencarian full-text

```
Labels: sprint-2, endpoint
Depends on: #5, #14
```

**Acceptance Criteria:**
- [ ] `?q=keyword` — wajib diisi, HTTP 400 jika kosong
- [ ] Filter: `?surah_id=1`, `?juz=1`
- [ ] Filter: `?lang=id` atau `?lang=en` (default: id)
- [ ] Pagination: `?page=1&limit=20`
- [ ] Response: `{ query, total, page, limit, results: [{ id, surah_info, number, number_in_surah, text_uthmani, translation, juz }] }`

---

### #17 — GET /random — Ayat acak

```
Labels: sprint-2, endpoint
Depends on: #5, #13
```

**Acceptance Criteria:**
- [ ] Mengembalikan satu ayat secara acak
- [ ] Query param: `?lang=id` atau `?lang=en`
- [ ] Query param: `?surah_id=1` untuk membatasi ke surah tertentu
- [ ] HTTP 200

---

### #18 — Pagination helper

```
Labels: sprint-2, infrastructure
Depends on: #1
```

**Deskripsi:**
Buat helper terpusat untuk pagination agar seluruh endpoint yang mendukung `?page` dan `?limit` menggunakan logika yang sama.

**Acceptance Criteria:**
- [ ] Helper `ParsePagination(page, limit string)` menghasilkan offset + limit
- [ ] Default: `page=1`, `limit=20`
- [ ] Batas maksimal limit: 100
- [ ] Unit test mencakup edge case (page 0, limit negatif, limit > 100)
- [ ] Digunakan oleh #15 dan #16

---

### #19 — Unit test suite Sprint 1

```
Labels: sprint-2, testing
Depends on: #8, #9, #10, #11
```

**Deskripsi:**
Tulis unit test lengkap untuk seluruh handler dan repository yang dibangun di Sprint 1.

**Acceptance Criteria:**
- [ ] Coverage ≥ 70% untuk package `handler` dan `repository`
- [ ] Test mencakup happy path dan error path (404, 400, 500)
- [ ] Menggunakan SQLite in-memory untuk test database
- [ ] `go test ./...` lolos tanpa error

---

## Sprint 3 — Kesiapan Rilis
> Tujuan: API siap dikonsumsi secara internal, terdokumentasi, aman, dan lolos release checklist governance.
>
> Issue #20, #21, #22, #23, dan #24 dapat dikerjakan **secara paralel** sejak awal sprint. Issue #25 menunggu #22 selesai. Issue #26 menunggu #25 dan semua endpoint selesai. Issue #27 dikerjakan **terakhir** setelah semua issue lain di sprint ini selesai.

---

### #20 — GET /health — Health check

```
Labels: sprint-3, ops
Depends on: #1
```

**Acceptance Criteria:**
- [ ] `GET /health` mengembalikan `{ status: "ok", timestamp: "...", version: "..." }`
- [ ] `GET /health/ready` memeriksa apakah file SQLite dapat diakses
- [ ] HTTP 503 jika tidak siap

---

### #21 — CORS & security headers

```
Labels: sprint-3, security
Depends on: #1
```

**Acceptance Criteria:**
- [ ] CORS `Allow-Origin` dikonfigurasi melalui env variable `ALLOWED_ORIGINS`
- [ ] Default: domain super app spesifik (bukan wildcard)
- [ ] Security headers: `X-Content-Type-Options: nosniff`, `X-Frame-Options: DENY`
- [ ] OPTIONS preflight ditangani dengan benar

---

### #22 — Unit test suite Sprint 2

```
Labels: sprint-3, testing
Depends on: #15, #16, #17
```

**Acceptance Criteria:**
- [ ] Coverage ≥ 70% untuk handler dan repository Sprint 2
- [ ] Test mencakup happy path dan error path
- [ ] `go test ./...` lolos tanpa error

---

### #23 — Dockerfile & konfigurasi deployment

```
Labels: sprint-3, ops
Depends on: #1
```

**Deskripsi:**
Siapkan Dockerfile dan konfigurasi deployment agar API dapat di-deploy ke target environment yang dipilih.

**Acceptance Criteria:**
- [ ] Dockerfile multi-stage build (builder + runner)
- [ ] Image final sekecil mungkin (gunakan `scratch` atau `alpine`)
- [ ] File SQLite di-copy ke dalam image atau di-mount sebagai volume
- [ ] Dokumentasi singkat cara deploy di README
- [ ] Build berhasil dengan `docker build`

---

### #24 — README & panduan kontributor

```
Labels: sprint-3, docs
Depends on: #1
```

**Deskripsi:**
Tulis dokumentasi onboarding agar kontributor baru dapat menjalankan project dalam waktu singkat.

**Acceptance Criteria:**
- [ ] Langkah setup lokal: clone → install → migrate → seed → run
- [ ] Penjelasan singkat struktur folder
- [ ] Cara menjalankan test
- [ ] Cara menjalankan migrasi dan seeder
- [ ] Kontributor baru dapat menjalankan project dalam < 10 menit

---

### #25 — Setup analytics & kanal feedback

```
Labels: sprint-3, ops
Depends on: #22
```

**Acceptance Criteria:**
- [ ] PostHog atau setara dipasang untuk melacak penggunaan endpoint
- [ ] Kanal feedback (link WA/email) tersedia di dalam super app
- [ ] Diperlukan untuk mengukur D7 retention dan feedback kualitatif

---

### #26 — Dokumentasi API (Scalar)

```
Labels: sprint-3, docs
Depends on: #8, #9, #10, #11, #15, #16, #17, #20
```

**Acceptance Criteria:**
- [ ] Scalar UI tersedia di `/docs`
- [ ] OpenAPI 3.0 spec mencakup seluruh endpoint MVP
- [ ] Setiap endpoint memiliki minimal satu contoh response
- [ ] Standalone (tanpa ketergantungan CDN eksternal)

---

### #27 — Release checklist & verifikasi governance

```
Labels: sprint-3, governance
Depends on: #20, #21, #25, #26
```

**Deskripsi:**
Verifikasi bahwa seluruh persyaratan governance terpenuhi sebelum API dirilis ke lingkungan internal.

**Acceptance Criteria:**
- [ ] Core loop berjalan tanpa critical error
- [ ] Seluruh endpoint P1 sudah diuji dan passing
- [ ] Tidak ada data sensitif yang terekspos di sisi client
- [ ] Konten sudah ditinjau (aman dari isu SARA/Syariah)
- [ ] Privacy Policy & Terms of Service sudah tersedia dan dapat diakses
- [ ] Lisensi data sudah dikonfirmasi untuk semua sumber
- [ ] Analytics sudah terpasang
- [ ] Kanal feedback tersedia
- [ ] Disetujui oleh Dewan Pendiri
- [ ] Disetujui oleh Dewan Pembina

---

## Ringkasan Issue

| Issue | Judul | Sprint | Dapat Paralel Dengan |
|-------|-------|--------|----------------------|
| #1 | Setup project | 0 | - |
| #2 | Skema & migrasi | 0 | #4 (setelah #1) |
| #3 | Seed data | 0 | - |
| #4 | Logging middleware | 0 | #2 |
| #5 | Shared DTOs & error structs | 1 | - |
| #6 | Surah repository | 1 | #7, #12 |
| #7 | Ayah repository | 1 | #6, #12 |
| #8 | GET /surah & /surah/:id | 1 | #9, #10, #11 |
| #9 | GET /surah/:id/ayat | 1 | #8, #10, #11 |
| #10 | GET /surah/:id/ayat/:number | 1 | #8, #9, #11 |
| #11 | GET /ayah/:id | 1 | #8, #9, #10 |
| #12 | Input validation middleware | 1 | #6, #7 |
| #13 | Juz repository | 2 | #14 |
| #14 | Search repository (FTS5) | 2 | #13 |
| #15 | GET /juz & /juz/:number | 2 | #16, #17 |
| #16 | GET /search | 2 | #15, #17 |
| #17 | GET /random | 2 | #15, #16 |
| #18 | Pagination helper | 2 | semua issue Sprint 2 |
| #19 | Unit test Sprint 1 | 2 | semua issue Sprint 2 |
| #20 | GET /health | 3 | #21, #22, #23, #24 |
| #21 | CORS & security headers | 3 | #20, #22, #23, #24 |
| #22 | Unit test Sprint 2 | 3 | #20, #21, #23, #24 |
| #23 | Dockerfile & deployment | 3 | #20, #21, #22, #24 |
| #24 | README & panduan kontributor | 3 | #20, #21, #22, #23 |
| #25 | Analytics & feedback channel | 3 | setelah #22 |
| #26 | Dokumentasi Scalar | 3 | setelah semua endpoint |
| #27 | Release checklist & governance | 3 | terakhir |

---

## Ringkasan Endpoint

| Method | Endpoint | Fitur Super App | Issue | Sprint |
|--------|----------|-----------------|-------|--------|
| GET | `/surah` | Halaman daftar surah | #8 | 1 |
| GET | `/surah/:id` | Halaman detail surah | #8 | 1 |
| GET | `/surah/:id/ayat` | Halaman baca surah | #9 | 1 |
| GET | `/surah/:id/ayat/:number` | Referensi ayat | #10 | 1 |
| GET | `/ayah/:id` | Ayat by global ID | #11 | 1 |
| GET | `/juz` | Navigasi juz | #15 | 2 |
| GET | `/juz/:number` | Halaman baca juz | #15 | 2 |
| GET | `/search` | Fitur pencarian | #16 | 2 |
| GET | `/random` | Widget ayat harian | #17 | 2 |
| GET | `/health` | Monitoring operasional | #20 | 3 |
| GET | `/docs` | Referensi developer | #26 | 3 |

**Total: 11 endpoint**

---

## Persyaratan Fungsional

- FR-1: CORS origin harus dapat dikonfigurasi, tidak boleh hardcode wildcard
- FR-2: Semua response harus JSON dengan Content-Type yang sesuai
- FR-3: Format error response harus konsisten: `{ error: string, code: string, timestamp: string }`
- FR-4: Parameter `lang` hanya menerima `id` atau `en`
- FR-5: Pencarian harus case-insensitive dan mendukung partial match via FTS5
- FR-6: Semua log menggunakan format structured (zerolog)
- FR-7: Database bersifat read-only setelah seeding — tidak ada write endpoint di MVP

---

## Metrik Keberhasilan

### Teknis
- P95 response time < 200ms untuk endpoint surah
- P95 response time < 500ms untuk pencarian
- Zero data loss setelah seeding (114 surah, 6.236 ayat, 30 juz)
- Seluruh endpoint terdokumentasi di Scalar

### Produk (Wajib per Kebijakan MVP)
- D7 Retention: ≥ 20% fitur super app yang menggunakan API ini kembali aktif dalam 7 hari
- Feedback kualitatif: ≥ 50 respons dalam 30 hari setelah rilis (via WA/form)
- Core loop: developer dapat mengambil surah lengkap beserta terjemahan dalam < 3 menit setelah membaca /docs

---

## Environment Variables

```bash
# Database
DB_PATH=./data/quran.db

# Server
SERVER_PORT=8080
SERVER_HOST=0.0.0.0

# CORS
ALLOWED_ORIGINS=https://[domain-superapp].com

# App
APP_VERSION=1.0.0
LOG_LEVEL=info
```

---

## Struktur File

```
quran-api-go/
├── cmd/api/
│   └── main.go               # Manual DI wiring di sini
├── internal/
│   ├── config/
│   ├── domain/
│   │   ├── surah/
│   │   ├── ayah/
│   │   └── juz/
│   ├── handler/
│   ├── repository/
│   ├── service/
│   └── middleware/
│       ├── cors.go
│       ├── logging.go
│       └── recovery.go
├── pkg/
│   ├── response/
│   └── pagination/
├── migrations/
├── scripts/seed/
├── docs/openapi.yaml
├── data/                     # File SQLite .db disimpan di sini
├── .env.example
├── Dockerfile
├── Makefile
└── go.mod
```

---
