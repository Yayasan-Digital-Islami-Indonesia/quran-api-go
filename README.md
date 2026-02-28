<div align="center">

# Quran API Go

### Lightweight RESTful API untuk Data Al-Quran

<p align="center">
  <a href="https://deepwiki.com/Yayasan-Digital-Islami-Indonesia/quran-api-go"><img src="https://deepwiki.com/badge.svg"></a>
  <a href="https://github.com/moeru-ai/airi/blob/main/LICENSE"><img src="https://img.shields.io/github/license/moeru-ai/airi.svg?style=flat&colorA=080f12&colorB=1fa669"></a>
  <a href="https://discord.gg/hJtr47KXaK"><img src="https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fdiscord.com%2Fapi%2Finvites%2FhJtr47KXaK%3Fwith_counts%3Dtrue&query=%24.approximate_member_count&suffix=%20members&logo=discord&logoColor=white&label=%20&color=7389D8&labelColor=6A7EC2"></a>
  <a href="https://github.com/Yayasan-Digital-Islami-Indonesia/quran-api-go/network/members"><img src="https://img.shields.io/github/forks/Yayasan-Digital-Islami-Indonesia/quran-api-go?style=flat&logo=github&logoColor=white&label=Fork" alt="Forks"></a>
  <a href="https://github.com/Yayasan-Digital-Islami-Indonesia/quran-api-go/stargazers"><img src="https://img.shields.io/github/stars/Yayasan-Digital-Islami-Indonesia/quran-api-go?style=flat&logo=github&logoColor=white&label=Star" alt="Stars"></a>
  <a href="https://github.com/Yayasan-Digital-Islami-Indonesia/quran-api-go/issues"><img src="https://img.shields.io/github/issues/Yayasan-Digital-Islami-Indonesia/quran-api-go?style=flat&logo=github&logoColor=white&label=Issues" alt="Issues"></a>
  <a href="https://go.dev/"><img src="https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go" alt="Go"></a>
  <a href="https://www.sqlite.org/"><img src="https://img.shields.io/badge/SQLite-FTS5-07405E?style=flat&logo=sqlite&logoColor=white" alt="SQLite"></a>
</p>

</div>

---

**Quran API Go** adalah API internal untuk super app Ilmunara. Menyediakan 114 surah, 6.236 ayat, 30 juz dengan terjemahan ID/EN.

- Cepat — P95 < 200ms
- Ringan — Single binary, SQLite embedded
- Simple — JSON response

---

## Quick Start

```bash
git clone https://github.com/quran-api-go/quran-api-go.git
cd quran-api-go
go mod download
make migrate && make seed && make run
```

Server jalan di `http://localhost:8080`

**Docker:**
```bash
docker build -t quran-api-go .
docker run -p 8080:8080 -e ALLOWED_ORIGINS=https://yourapp.com quran-api-go
```

---

## Endpoint

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| GET | `/surah` | Daftar 114 surah |
| GET | `/surah/:id` | Detail surah |
| GET | `/surah/:id/ayat` | Ayat dalam surah |
| GET | `/surah/:id/ayat/:number` | Ayat spesifik |
| GET | `/ayah/:id` | Ayat by global ID (1-6236) |
| GET | `/juz` | Daftar 30 juz |
| GET | `/juz/:number` | Ayat dalam juz |
| GET | `/search` | Cari ayat by keyword |
| GET | `/random` | Ayat acak |
| GET | `/health` | Health check |
| GET | `/docs` | Dokumentasi API |

---

## Contoh

**Daftar Surah:**
```bash
curl http://localhost:8080/surah
```

**Baca Surah:**
```bash
curl "http://localhost:8080/surah/1/ayat?lang=id"
```

**Cari:**
```bash
curl "http://localhost:8080/search?q=rahman&page=1&limit=10"
```

---

## Query Parameters

| Param | Value |
|-------|-------|
| `lang` | `id` atau `en` (default: `id`) |
| `from` / `to` | Range ayat |
| `page` / `limit` | Pagination (default: `1`, `20`; max limit: `100`) |

---

## Konfigurasi

| Env Variable | Default |
|--------------|---------|
| `DB_PATH` | `./data/quran.db` |
| `SERVER_PORT` | `8080` |
| `ALLOWED_ORIGINS` | - |
| `LOG_LEVEL` | `info` |

---

## Tech Stack

```
Go 1.22+ • Gin • SQLite FTS5 • Goose • Zerolog
```

---

## Development

```bash
make test    # run tests
make lint    # static analysis
```

## Kontribusi via Fork

```bash
# 1. Fork repo, lalu clone fork kamu
git clone https://github.com/YOUR_USERNAME/quran-api-go.git
cd quran-api-go

# 2. Tambah upstream
git remote add upstream https://github.com/Yayasan-Digital-Islami-Indonesia/quran-api-go.git

# 3. Buat branch, coding, test
git checkout -b feature/fitur-kamu
# ... edit code ...
make test && make lint

# 4. Push ke fork, buat PR
git push origin feature/fitur-kamu
# Buka PR di GitHub → "Compare across forks"
```

Lihat [CONTRIBUTING.md](CONTRIBUTING.md) untuk detail.

---

## License

MIT
