# Panduan Kontribusi — Quran API Go

Jazakallahu khairan atas minatmu untuk berkontribusi! Proyek ini dikelola oleh kontributor sukarela di bawah **Yayasan Ilmunara Teknologi Indonesia**. Mohon baca panduan ini dengan seksama sebelum mulai mengerjakan apapun — ini membantu kita menjaga keberlanjutan proyek dan menghargai waktu sesama relawan.

---

## Daftar Isi

- [Kode Etik](#kode-etik)
- [Sebelum Mulai](#sebelum-mulai)
- [Persiapan Lokal](#persiapan-lokal)
- [Alur Kerja Pengembangan](#alur-kerja-pengembangan)
- [Konvensi Pesan Commit](#konvensi-pesan-commit)
- [Panduan Pull Request](#panduan-pull-request)
- [Kebijakan Governance & Scope](#kebijakan-governance--scope)
- [Quality Gates](#quality-gates)
- [Butuh Bantuan?](#butuh-bantuan)

---

## Kode Etik

Proyek ini adalah bagian dari inisiatif teknologi Islam nirlaba. Kami mengharapkan seluruh kontributor untuk:

- Berkomunikasi dengan sopan dan konstruktif.
- Memberi dan menerima masukan dengan niat baik (*husnuzhon*).
- Menyampaikan ketidaksetujuan melalui jalur yang tepat — bukan di kolom komentar PR.
- Menghargai waktu dan usaha sesama relawan.

Pelanggaran dapat mengakibatkan pencabutan akses sebagai kontributor.

---

## Sebelum Mulai

> **Penting:** Kami beroperasi di bawah kebijakan governance MVP yang wajib dipatuhi. Sebelum mengerjakan fitur atau perbaikan apapun, pastikan:

1. Pekerjaan sudah tercakup dalam User Story yang telah disetujui di PRD.
2. Kamu **tidak** menambahkan dependensi baru tanpa diskusi terlebih dahulu (terutama layanan berbayar, Redis, atau DI framework).
3. Jika ragu, tanyakan di **grup WA terlebih dahulu** sebelum mulai menulis kode.

Pekerjaan yang di luar scope atau belum disetujui tidak akan di-merge, terlepas dari kualitas kodenya.

---

## Persiapan Lokal

### Prasyarat

- Go 1.21+ — [install](https://go.dev/dl/)
- `make` — biasanya sudah terinstall (Linux/macOS). Windows: lewat WSL.
- Docker & Docker Compose (opsional) — [install](https://docs.docker.com/get-docker/) — hanya jika ingin jalankan dalam container

### Setup

```bash
# 1. Fork lalu clone repo
git clone https://github.com/Yayasan-Digital-Islami-Indonesia/quran-api-go
cd quran-api-go

# 2. Salin environment variables
cp .env.example .env

# 3. Buat directory data (untuk SQLite)
mkdir -p data

# 4. Jalankan migrasi (membuat quran.db)
make migrate

# 5. Seed data
make seed

# 6. Jalankan dev server
make run
```

---

## Alur Kerja Pengembangan

1. **Ambil task** — Hanya kerjakan issue yang sudah di-assign kepadamu atau sudah didiskusikan di grup WA.
2. **Buat branch dari `main`** menggunakan konvensi penamaan di bawah.
3. **Tulis kodenya** mengikuti struktur proyek yang ada di PRD.
4. **Pastikan semua quality gates lolos** sebelum membuka PR.
5. **Buka PR** ke `main` dan isi PR template dengan lengkap.

### Penamaan Branch

| Tipe | Pola | Contoh |
|------|------|--------|
| Fitur | `feat/<nomor-us>-deskripsi-singkat` | `feat/us-006-surah-ayat-endpoint` |
| Bug fix | `fix/<deskripsi-singkat>` | `fix/search-ilike-case` |
| Chore / infra | `chore/<deskripsi-singkat>` | `chore/update-makefile` |

---

## Konvensi Pesan Commit

Kami mengikuti [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <ringkasan singkat>

[isi opsional]
```

**Tipe yang digunakan:** `feat`, `fix`, `chore`, `docs`, `test`, `refactor`

**Contoh:**
```
feat(surah): tambah endpoint GET /surah/:id/ayat
fix(search): tangani query parameter kosong
chore(docker): tambah healthcheck ke service postgres
```

Jaga agar setiap commit bersifat atomik — satu perubahan logis per commit.

---

## Panduan Pull Request

- **Satu PR per User Story** — jangan menggabungkan beberapa US dalam satu PR.
- Isi PR template secara lengkap.
- Cantumkan nomor US terkait di judul PR, contoh: `[US-006] Endpoint ayat per surat`.
- PR membutuhkan **persetujuan dari Tech Lead** sebelum bisa di-merge.
- Jangan merge PR milikmu sendiri.
- Jaga agar PR tetap kecil dan mudah di-review. PR yang terlalu besar akan diminta untuk dipecah.

### Template PR

Saat membuka PR, sertakan:

```markdown
## Ringkasan
<!-- Apa yang dilakukan PR ini? -->

## User Story Terkait
<!-- contoh: US-006 -->

## Perubahan yang Dilakukan
<!-- Daftar perubahan utama -->

## Quality Gates
- [ ] `go test ./...` lolos
- [ ] `go vet ./...` lolos
- [ ] `gofmt -d .` tidak menampilkan diff

## Catatan untuk Reviewer
<!-- Hal yang perlu diperhatikan reviewer -->
```

---

## Kebijakan Governance & Scope

Proyek ini mengikuti **Kebijakan Strategis & Tata Kelola MVP** Yayasan Ilmunara Teknologi Indonesia. Sebagai kontributor, kamu wajib memahami aturan berikut:

| Aturan | Detail |
|--------|--------|
| ⛔ Dilarang microservices | Hanya arsitektur monolith modular |
| ⛔ Dilarang SaaS berbayar | Semua layanan harus free-tier atau self-hosted |
| ⛔ Dilarang scope creep | Hanya implementasikan yang ada di PRD yang sudah disetujui |
| ✅ Free-tier BaaS | Supabase / Firebase / Vercel direkomendasikan |
| ✅ Utamakan kesederhanaan | Jika solusi lebih sederhana sudah cukup, gunakan itu |

Fitur yang dibangun di luar batasan ini **tanpa persetujuan tertulis** dari Kepala Produk & Tech Lead tidak akan di-merge dan akan di-rollback.

---

## Quality Gates

Setiap PR **wajib** melewati pemeriksaan berikut sebelum masuk tahap review:

```bash
go test ./...   # Semua test harus lolos
go vet ./...    # Tidak ada error static analysis
gofmt -d .      # Tidak ada masalah formatting
```

PR yang gagal pada pemeriksaan ini tidak akan di-review sampai diperbaiki.

---

## Butuh Bantuan?

- **Grup WA** — Untuk pertanyaan cepat, klarifikasi task, atau diskusi sebelum mulai mengerjakan. Ini adalah saluran komunikasi utama.
- **GitHub Issues** — Untuk laporan bug atau pertanyaan teknis yang terkait dengan US tertentu.

Jika kamu tidak yakin apakah sesuatu termasuk dalam scope, **tanyakan di grup WA terlebih dahulu**. Ini menghemat waktu semua orang.

---

*Barakallahu fiikum. Setiap baris kode yang kamu kontribusikan adalah bentuk sedekah jariyah — semoga Allah menerimanya dari kita semua.*
