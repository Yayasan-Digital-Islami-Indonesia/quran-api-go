<div align="center">

# Quran API Go

### Lightweight RESTful API + MCP Server untuk Data Al-Quran

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

REST API Al-Quran Indonesia dengan MCP Server bawaan. Menyediakan 114 surah, 6.236 ayat, 30 juz dengan terjemahan ID/EN — bisa diakses via HTTP maupun langsung dari AI assistant.

- Cepat — P95 < 200ms
- Ringan — Single binary, SQLite embedded
- AI-ready — MCP Server untuk Claude, Cursor, dan tools lainnya

---

## Quick Start

```bash
git clone https://github.com/Yayasan-Digital-Islami-Indonesia/quran-api-go.git
cd quran-api-go
go mod download
go run ./cmd/migrate && go run ./cmd/seed --data ./data/seed && go run ./cmd/api
```

Server jalan di `http://localhost:8080` · Docs di `http://localhost:8080/docs`

**Docker:**
```bash
docker build -t quran-api-go .
docker run -p 8080:8080 -e ALLOWED_ORIGINS=https://yourapp.com quran-api-go
```

> Docker otomatis jalankan migrasi sebelum server start via `entrypoint.sh`.

---

## Endpoint

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| GET | `/surah` | Daftar 114 surah |
| GET | `/surah?type=meccan\|medinan` | Filter surah by revelation type |
| GET | `/surah/:id` | Detail surah |
| GET | `/surah/:id/ayah` | Ayat dalam surah (optional range) |
| GET | `/surah/:id/ayah/:number` | Ayat spesifik dalam surah |
| GET | `/ayah/:id` | Ayat by global ID (1-6236) |
| GET | `/sajda` | Daftar 15 ayat sajda tilawah |
| GET | `/random` | Ayat acak |
| GET | `/juz` | Daftar 30 juz |
| GET | `/juz/:number` | Detail juz |
| GET | `/juz/:number/ayah` | Ayat dalam juz (paginated) |
| GET | `/juz/:number/surah` | Surah yang ada dalam juz |
| GET | `/search` | Full-text search (Arab, ID, EN) |
| GET | `/health` | Health check |
| GET | `/health/ready` | Readiness check |
| GET | `/docs` | Dokumentasi API (Scalar) |

---

## MCP Server

API ini dilengkapi **MCP (Model Context Protocol) server** sehingga bisa digunakan langsung dari AI assistant.

**Koneksi via Streamable HTTP:**

| | |
|---|---|
| URL | `https://quran.wahyuikbal.com/mcp` |
| Transport | Streamable HTTP |

**Setup Claude Desktop** — tambahkan ke `claude_desktop_config.json`:
```json
{
  "mcpServers": {
    "quran": {
      "type": "http",
      "url": "https://quran.wahyuikbal.com/mcp"
    }
  }
}
```

**Tools yang tersedia:**

| Tool | Deskripsi |
|------|-----------|
| `list_surahs` | Daftar semua 114 surah |
| `get_surah` | Detail surah by ID |
| `get_ayahs_by_surah` | Ayat dalam surah |
| `get_ayah` | Ayat by global ID |
| `get_ayah_by_ref` | Ayat by nomor surah + ayat |
| `random_ayah` | Ayat acak |
| `list_juz` | Daftar semua 30 juz |
| `get_juz` | Detail juz |
| `get_ayahs_by_juz` | Ayat dalam juz |
| `search_quran` | Full-text search |

**Atau jalankan lokal via stdio** (untuk Claude Desktop lokal):
```bash
go run ./cmd/mcp
```

---

## Contoh

```bash
# Filter surah Makkiyah
curl "http://localhost:8080/surah?type=meccan"

# Ayat sajda tilawah
curl "http://localhost:8080/sajda?lang=id"

# Surah yang ada di Juz 30
curl "http://localhost:8080/juz/30/surah"

# Baca surah dengan terjemahan
curl "http://localhost:8080/surah/1/ayah?lang=id"

# Cari ayat
curl "http://localhost:8080/search?q=sabar&lang=id&page=1&limit=10"
```

---

## Query Parameters

| Param | Value |
|-------|-------|
| `lang` | `id` atau `en` (default: `id`) |
| `type` | `meccan` atau `medinan` (khusus `/surah`) |
| `from` / `to` | Range ayat |
| `page` / `limit` | Pagination (default: `1`, `20`; max: `100`) |

---

## Konfigurasi

| Env Variable | Default | Keterangan |
|--------------|---------|------------|
| `DB_PATH` | `./data/quran.db` | Path ke SQLite database |
| `SERVER_PORT` | `8080` | Port server |
| `SERVER_HOST` | `0.0.0.0` | Host server |
| `ALLOWED_ORIGINS` | - | Allowed CORS origins. Gunakan `*` untuk allow semua (MCP public) |
| `APP_VERSION` | `1.0.0` | Versi aplikasi |
| `LOG_LEVEL` | `info` | Level logging |

---

## Tech Stack

```
Go 1.22+ · Gin · SQLite FTS5 · Goose · Zerolog · MCP Go SDK · swaggo
```

---

## Development

```bash
go run ./cmd/api          # Jalankan API server
go run ./cmd/mcp          # Jalankan MCP server (stdio)
go test ./...             # Run tests
go vet ./...              # Lint
go run ./cmd/migrate      # Jalankan migrasi
go run ./cmd/seed --data ./data/seed  # Seed database

# Regenerate OpenAPI docs (setelah ubah handler)
swag init -g cmd/api/main.go -o docs --outputTypes go,yaml
cp docs/swagger.yaml docs/api-reference/openapi.yaml
```

---

## Kontribusi via Fork

```bash
# 1. Fork repo, lalu clone fork kamu
git clone https://github.com/YOUR_USERNAME/quran-api-go.git
cd quran-api-go

# 2. Tambah upstream
git remote add upstream https://github.com/Yayasan-Digital-Islami-Indonesia/quran-api-go.git

# 3. Buat branch, coding, test
git checkout -b feature/fitur-kamu
go test ./... && go vet ./...

# 4. Push ke fork, buat PR
git push origin feature/fitur-kamu
```

---

## License

MIT
