# PRD: Quran API Go

## Overview
Membangun RESTful API publik untuk Al-Quran menggunakan Golang dan Gin framework. API menyediakan endpoints untuk mengakses data surat, ayat dengan terjemahan Indonesia & English, serta fitur navigasi lengkap (juz, halaman, hizb, manzil, ruku) dan pencarian teks dengan filter. Data sudah tersedia dan siap di-import.

## Goals
- Menyediakan API publik yang cepat dan ringan untuk mengakses data Al-Quran
- Mendukung 2 bahasa terjemahan: Indonesia & English
- Navigasi lengkap: surat, ayat, juz, halaman, hizb, manzil, ruku
- Pencarian teks ayat dengan filter yang fleksibel
- Open API tanpa autentikasi dengan rate limiting by IP via Redis
- API documentation yang user-friendly dengan Scalar
- Siap untuk dikonsumsi publik

## Quality Gates

These commands must pass for every user story:
- `go test ./...` - Unit tests
- `go vet ./...` - Static analysis
- `gofmt -d .` - Code formatting check

## User Stories

### US-001: Project setup dan infrastructure
**Description:** As a developer, I want to set up the project with proper infrastructure so that development can begin.

**Acceptance Criteria:**
- [ ] Initialize Go module dengan `go mod init`
- [ ] Setup Gin framework sebagai HTTP router
- [ ] Setup PostgreSQL connection dengan connection pooling
- [ ] Setup Redis untuk rate limiting
- [ ] Setup Wire atau Uber FX untuk dependency injection
- [ ] Setup structured logging (misal: zap atau zerolog)
- [ ] Setup basic metrics middleware
- [ ] Create Makefile dengan target: `run`, `test`, `lint`, `migrate`, `seed`
- [ ] Setup Docker Compose untuk PostgreSQL + Redis dev environment

### US-002: Database schema dan migrations
**Description:** As a developer, I want to create database schema and migration system so that data Quran tersimpan terstruktur.

**Acceptance Criteria:**
- [ ] Setup Dbmate atau Goose untuk migrations
- [ ] Create migration untuk tabel `surahs`: id, number, name_arabic, name_latin, name_transliteration, number_of_ayahs, revelation_type, pages
- [ ] Create migration untuk tabel `ayahs`: id, surah_id, number_in_surah, text_uthmani, translation_indo, translation_en, juz_number, page_number, manzil_number, ruku_number, sajda_type, revelation_type
- [ ] Create migration untuk tabel `juzs`: id, juz_number, first_ayah_id, last_ayah_id
- [ ] Create migration untuk tabel `pages`: id, page_number, first_ayah_id, last_ayah_id
- [ ] Create migration untuk tabel `hizbs`: id, hizb_number, first_ayah_id, last_ayah_id
- [ ] Create migration untuk tabel `rub_el_hizbs`: id, hizb_number, quarter_number, first_ayah_id, last_ayah_id
- [ ] Create migration untuk tabel `manzils`: id, manzil_number, first_ayah_id, last_ayah_id
- [ ] Create migration untuk tabel `rukus`: id, ruku_number, surah_id, first_ayah_id, last_ayah_id
- [ ] Add indexes untuk pencarian: surah_id, juz_number, page_number, manzil_number, ruku_number, text_uthmani (GIN/trigram index untuk full-text search)
- [ ] Migration harus reversible (up/down)

### US-003: Seed data dari dataset yang sudah disiapkan
**Description:** As a developer, I want to seed Quran data dari dataset yang sudah ada so that API bisa return data yang valid.

**Acceptance Criteria:**
- [ ] Create seeder script untuk import data dari dataset yang sudah disiapkan
- [ ] Import semua 114 surat
- [ ] Import semua 6236 ayat dengan lengkap (text, translations, navigasi data)
- [ ] Import data 30 juz
- [ ] Import data 604 halaman (sistem penomoran Madinah)
- [ ] Import data 240 rub el hizb (8 per juz)
- [ ] Import data 7 manzil
- [ ] Import data 556 ruku (kurang lebih)
- [ ] Seeder harus idempotent (bisa di-run berkali-kali)
- [ ] Add logging progress saat seeding
- [ ] Validasi data kelengkapan setelah seeding

### US-004: Endpoint GET /surah - List semua surat
**Description:** As an API consumer, I want to get list of all surahs so that I can browse Al-Quran content.

**Acceptance Criteria:**
- [ ] GET /surah return array semua surat
- [ ] Support query parameter: `?page=1&limit=10` untuk pagination
- [ ] Response structure: `[{ id, number, name_arabic, name_latin, name_transliteration, number_of_ayahs, revelation_type }]`
- [ ] Return HTTP 200 dengan proper headers (Content-Type: application/json)
- [ ] Return HTTP 500 jika database error

### US-005: Endpoint GET /surah/:id - Detail surat
**Description:** As an API consumer, I want to get detail surat information so that I can see surat metadata.

**Acceptance Criteria:**
- [ ] GET /surah/:id return detail surat tanpa ayat
- [ ] Response structure: `{ id, number, name_arabic, name_latin, name_transliteration, number_of_ayahs, revelation_type, pages: [start_page, end_page] }`
- [ ] Return HTTP 404 jika surat tidak ditemukan

### US-006: Endpoint GET /surah/:id/ayat - Ayat-ayat dalam surat
**Description:** As an API consumer, I want to get all ayahs in a surah so that I can read complete surah content.

**Acceptance Criteria:**
- [ ] GET /surah/:id/ayat return array ayat dalam surat
- [ ] Support query parameter: `?lang=id` atau `?lang=en` (default: id)
- [ ] Support query parameter: `?from=1&to=10` untuk range ayat (opsional)
- [ ] Support query parameter: `?page=1&limit=20` untuk pagination
- [ ] Response structure: `{ surah: { id, number, name_latin }, ayahs: [{ number, number_in_surah, text_uthmani, translation, juz, page, manzil, ruku, sajda }] }`
- [ ] Return HTTP 400 jika `lang` tidak valid
- [ ] Return HTTP 404 jika surat tidak ditemukan

### US-007: Endpoint GET /surah/:id/ayat/:number - Ayat spesifik dalam surat
**Description:** As an API consumer, I want to get specific ayah by surah and ayah number so that I can reference individual ayah.

**Acceptance Criteria:**
- [ ] GET /surah/:id/ayat/:number return detail ayat
- [ ] Support query parameter: `?lang=id` atau `?lang=en`
- [ ] Response structure: `{ id, surah_id, number, number_in_surah, text_uthmani, translation, surah_info: { id, name_latin }, juz, page, manzil, ruku, sajda, revelation_type }`
- [ ] Return HTTP 400 jika `lang` tidak valid
- [ ] Return HTTP 404 jika ayat tidak ditemukan

### US-008: Endpoint GET /ayah/:id - Detail ayat by global ID
**Description:** As an API consumer, I want to get ayah by global ID so that I can reference ayah across surahs.

**Acceptance Criteria:**
- [ ] GET /ayah/:id return detail ayat by global ID (1-6236)
- [ ] Support query parameter: `?lang=id` atau `?lang=en`
- [ ] Response structure sama seperti US-007
- [ ] Return HTTP 400 jika `lang` tidak valid
- [ ] Return HTTP 404 jika ayat tidak ditemukan

### US-009: Endpoint GET /juz - List semua juz
**Description:** As an API consumer, I want to get list of all juzs so that I can navigate Al-Quran by juz.

**Acceptance Criteria:**
- [ ] GET /juz return array semua juz (1-30)
- [ ] Response structure: `[{ juz_number, first_ayah_id, last_ayah_id, total_ayahs }]`
- [ ] Return HTTP 200

### US-010: Endpoint GET /juz/:number - Detail juz dengan ayat
**Description:** As an API consumer, I want to get all ayahs in a juz so that I can read complete juz content.

**Acceptance Criteria:**
- [ ] GET /juz/:number return ayat-ayat dalam juz
- [ ] Support query parameter: `?lang=id` atau `?lang=en`
- [ ] Support query parameter: `?page=1&limit=20` untuk pagination
- [ ] Response structure: `{ juz: { juz_number }, ayahs: [{ number, surah_info, number_in_surah, text_uthmani, translation, page }] }`
- [ ] Return HTTP 404 jika juz tidak valid (bukan 1-30)

### US-011: Endpoint GET /page - List semua halaman
**Description:** As an API consumer, I want to get list of all pages so that I can navigate Al-Quran by page.

**Acceptance Criteria:**
- [ ] GET /page return array semua halaman (1-604)
- [ ] Response structure: `[{ page_number, first_ayah_id, last_ayah_id, surahs: [{ id, name_latin }] }]`
- [ ] Return HTTP 200

### US-012: Endpoint GET /page/:number - Detail halaman dengan ayat
**Description:** As an API consumer, I want to get all ayahs in a page so that I can read complete page content.

**Acceptance Criteria:**
- [ ] GET /page/:number return ayat-ayat dalam halaman
- [ ] Support query parameter: `?lang=id` atau `?lang=en`
- [ ] Response structure: `{ page: { page_number }, ayahs: [{ number, surah_info, number_in_surah, text_uthmani, translation }] }`
- [ ] Return HTTP 404 jika halaman tidak valid (bukan 1-604)

### US-013: Endpoint GET /hizb - List semua hizb
**Description:** As an API consumer, I want to get list of all hizbs so that I can navigate Al-Quran by hizb.

**Acceptance Criteria:**
- [ ] GET /hizb return array semua hizb (1-60, atau 1-8 per juz)
- [ ] Response structure: `[{ hizb_number, juz_number, first_ayah_id, last_ayah_id }]`
- [ ] Return HTTP 200

### US-014: Endpoint GET /hizb/:number - Detail hizb dengan ayat
**Description:** As an API consumer, I want to get all ayahs in a hizb so that I can read complete hizb content.

**Acceptance Criteria:**
- [ ] GET /hizb/:number return ayat-ayat dalam hizb
- [ ] Support query parameter: `?lang=id` atau `?lang=en`
- [ ] Response structure: `{ hizb: { hizb_number, juz_number }, ayahs: [{ number, surah_info, number_in_surah, text_uthmani, translation }] }`
- [ ] Return HTTP 404 jika hizb tidak valid

### US-015: Endpoint GET /rub-el-hizb - List semua rub el hizb
**Description:** As an API consumer, I want to get list of all rub el hizbs (quarter hizb) so that I can navigate finely.

**Acceptance Criteria:**
- [ ] GET /rub-el-hizb return array semua rub el hizb (1-240)
- [ ] Response structure: `[{ id, hizb_number, quarter_number, juz_number, first_ayah_id, last_ayah_id }]`
- [ ] Return HTTP 200

### US-016: Endpoint GET /rub-el-hizb/:number - Detail rub el hizb dengan ayat
**Description:** As an API consumer, I want to get all ayahs in a rub el hizb so that I can read quarter hizb content.

**Acceptance Criteria:**
- [ ] GET /rub-el-hizb/:number return ayat-ayat dalam rub el hizb
- [ ] Support query parameter: `?lang=id` atau `?lang=en`
- [ ] Response structure: `{ rub_el_hizb: { id, hizb_number, quarter_number }, ayahs: [{ number, surah_info, number_in_surah, text_uthmani, translation }] }`
- [ ] Return HTTP 404 jika rub el hizb tidak valid

### US-017: Endpoint GET /manzil - List semua manzil
**Description:** As an API consumer, I want to get list of all manzils so that I can navigate Al-Quran by manzil.

**Acceptance Criteria:**
- [ ] GET /manzil return array semua manzil (1-7)
- [ ] Response structure: `[{ manzil_number, first_ayah_id, last_ayah_id, total_ayahs }]`
- [ ] Return HTTP 200

### US-018: Endpoint GET /manzil/:number - Detail manzil dengan ayat
**Description:** As an API consumer, I want to get all ayahs in a manzil so that I can read complete manzil content.

**Acceptance Criteria:**
- [ ] GET /manzil/:number return ayat-ayat dalam manzil
- [ ] Support query parameter: `?lang=id` atau `?lang=en`
- [ ] Response structure: `{ manzil: { manzil_number }, ayahs: [{ number, surah_info, number_in_surah, text_uthmani, translation }] }`
- [ ] Return HTTP 404 jika manzil tidak valid (bukan 1-7)

### US-019: Endpoint GET /ruku - List semua ruku
**Description:** As an API consumer, I want to get list of all rukus so that I can navigate Al-Quran by ruku.

**Acceptance Criteria:**
- [ ] GET /ruku return array semua ruku (±556)
- [ ] Support query parameter: `?surah_id=1` untuk filter per surat
- [ ] Response structure: `[{ ruku_number, surah_id, surah_name, first_ayah_id, last_ayah_id, total_ayahs }]`
- [ ] Return HTTP 200

### US-020: Endpoint GET /ruku/:id - Detail ruku dengan ayat
**Description:** As an API consumer, I want to get all ayahs in a ruku so that I can read complete ruku content.

**Acceptance Criteria:**
- [ ] GET /ruku/:id return ayat-ayat dalam ruku
- [ ] Support query parameter: `?lang=id` atau `?lang=en`
- [ ] Response structure: `{ ruku: { ruku_number, surah_info }, ayahs: [{ number, number_in_surah, text_uthmani, translation }] }`
- [ ] Return HTTP 404 jika ruku tidak ditemukan

### US-021: Endpoint GET /search - Pencarian ayat dengan filter
**Description:** As an API consumer, I want to search ayahs by keyword with filters so that I can find specific content.

**Acceptance Criteria:**
- [ ] GET /search?q=keyword untuk pencarian full-text
- [ ] Support filter: `?surah_id=1` untuk limit ke surat tertentu
- [ ] Support filter: `?juz=1` untuk limit ke juz tertentu
- [ ] Support filter: `?manzil=1` untuk limit ke manzil tertentu
- [ ] Support filter: `?lang=id` atau `?lang=en` untuk bahasa terjemahan (default: id)
- [ ] Query parameter: `?page=1&limit=20` untuk pagination result
- [ ] Pencarian menggunakan PostgreSQL ILIKE atau tsvector untuk case-insensitive search
- [ ] Response structure: `{ query, total, page, limit, results: [{ id, surah_info, number, number_in_surah, text_uthmani, translation, juz, page }] }`
- [ ] Return HTTP 400 jika parameter `q` kosong

### US-022: Endpoint GET /random - Ayat random
**Description:** As an API consumer, I want to get a random ayah so that I can discover Quran content.

**Acceptance Criteria:**
- [ ] GET /random return satu ayat random
- [ ] Support query parameter: `?lang=id` atau `?lang=en`
- [ ] Support query parameter: `?surah_id=1` untuk random dalam surat tertentu
- [ ] Response structure: `{ ayah: { number, surah_info, number_in_surah, text_uthmani, translation, juz, page } }`
- [ ] Return HTTP 200

### US-023: Rate limiting middleware dengan Redis
**Description:** As an API owner, I want to implement rate limiting by IP using Redis so that API tidak disalahgunakan.

**Acceptance Criteria:**
- [ ] Implement rate limiting middleware dengan Redis backend
- [ ] Limit: 100 request per menit per IP
- [ ] Return HTTP 429 Too Many Requests dengan headers: `X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`, `Retry-After`
- [ ] Redis connection menggunakan connection pooling
- [ ] Configurable limit lewat environment variable `RATE_LIMIT_PER_MINUTE`
- [ ] Handle Redis connection failure gracefully (fallback to allow all)

### US-024: API Documentation dengan Scalar
**Description:** As an API consumer, I want to access interactive API documentation so that I can explore endpoints easily.

**Acceptance Criteria:**
- [ ] Setup Scalar untuk API documentation
- [ ] Serve Scalar UI di endpoint /docs
- [ ] Generate OpenAPI 3.0 spec dari annotations
- [ ] Document semua endpoint dengan method, path, parameters, response examples
- [ ] Include contoh response untuk setiap endpoint
- [ ] Scalar UI harus standalone

### US-025: Observability - Metrics endpoint
**Description:** As an operator, I want to expose metrics so that I can monitor API health and performance.

**Acceptance Criteria:**
- [ ] Expose Prometheus metrics at /metrics
- [ ] Track: request count, request duration, error rate by endpoint
- [ ] Track database connection pool stats
- [ ] Track Redis connection stats
- [ ] Track rate limiting: blocked requests count

### US-026: Health check endpoint
**Description:** As an operator/deployer, I want to check API health so that I can monitor service availability.

**Acceptance Criteria:**
- [ ] GET /health return status API
- [ ] Response structure: `{ status: "ok", timestamp: "...", version: "..." }`
- [ ] GET /health/ready untuk readiness probe (check DB & Redis connection)
- [ ] GET /health/live untuk liveness probe
- [ ] Return HTTP 503 jika dependency tidak ready

### US-027: CORS dan Security headers
**Description:** As an API owner, I want proper CORS dan security headers so that API bisa dikonsumsi dari frontend manapun.

**Acceptance Criteria:**
- [ ] Setup CORS middleware untuk allow all origins (public API)
- [ ] Add security headers: `X-Content-Type-Options: nosniff`, `X-Frame-Options: DENY`, `X-XSS-Protection`
- [ ] Add `Access-Control-Allow-Origin: *`
- [ ] Add `Access-Control-Allow-Headers: Content-Type, Authorization`
- [ ] Handle OPTIONS preflight request properly

## Functional Requirements
- FR-1: API harus mendukung CORS untuk semua origin (public API)
- FR-2: Semua response harus JSON dengan proper Content-Type header
- FR-3: Error response harus konsisten: `{ error: string, code: string, details: any, timestamp: string }`
- FR-4: Database connection harus menggunakan connection pooling
- FR-5: Redis connection harus menggunakan connection pooling
- FR-6: Logging harus mencatat: request method, path, status, duration, IP address
- FR-7: Rate limiting harus berbasis IP address dari request dengan storage di Redis
- FR-8: Pencarian harus case-insensitive dan support partial match
- FR-9: Parameter `lang` harus valid: hanya `id` atau `en`
- FR-10: Scalar docs harus accessible di /docs

## API Endpoints Summary

### Surah
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/surah` | List semua surat |
| GET | `/surah/:id` | Detail surat |
| GET | `/surah/:id/ayat` | Ayat-ayat dalam surat |
| GET | `/surah/:id/ayat/:number` | Ayat spesifik dalam surat |

### Ayah
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/ayah/:id` | Detail ayat by global ID |

### Juz
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/juz` | List semua juz |
| GET | `/juz/:number` | Detail juz dengan ayat |

### Page
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/page` | List semua halaman |
| GET | `/page/:number` | Detail halaman dengan ayat |

### Hizb
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/hizb` | List semua hizb |
| GET | `/hizb/:number` | Detail hizb dengan ayat |

### Rub el Hizb
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/rub-el-hizb` | List semua rub el hizb |
| GET | `/rub-el-hizb/:number` | Detail rub el hizb dengan ayat |

### Manzil
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/manzil` | List semua manzil |
| GET | `/manzil/:number` | Detail manzil dengan ayat |

### Ruku
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/ruku` | List semua ruku |
| GET | `/ruku/:id` | Detail ruku dengan ayat |

### Utility
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/search` | Pencarian ayat |
| GET | `/random` | Ayat random |
| GET | `/health` | Health check |
| GET | `/metrics` | Prometheus metrics |
| GET | `/docs` | API documentation (Scalar) |

**Total: ~25 endpoints**

## Non-Goals (Out of Scope)
- Autentikasi API Key (Phase 2)
- Multi-language translations selain Indo/English (Phase 2: Melayu, dll)
- Audio recitation endpoints (Phase 2)
- Tafsir endpoints (Phase 2)
- Tajweed rules endpoints (Phase 2)
- Word morphology endpoints (Phase 2)
- Redis caching layer untuk response (Phase 2)
- GraphQL support (Phase 2)
- Vector search / semantic search (Phase 2)
- WebSocket untuk real-time updates (not needed)
- Admin panel untuk manage data (Phase 2)

## Technical Considerations
- **PostgreSQL Full-Text Search**: Menggunakan `tsvector` column atau `ILIKE` untuk pencarian sederhana. Tambahkan trigram extension untuk better partial matching.
- **Redis untuk Rate Limiting**: Using sliding window algorithm dengan Redis INCR dan EXPIRE.
- **Scalar untuk Docs**: Scalar provides modern API documentation UI yang bisa di-serve langsung dari Go binary.
- **DI Framework**: Wire (compile-time) atau Uber FX (runtime) untuk dependency injection.
- **Migration**: Dbmate (simpel, SQL-based) atau Goose (Go-based).
- **Data Structure**: Ayat text dalam Uthmani script untuk Arabic, translations dalam kolom terpisah.

## Success Metrics
- API response time P95 < 200ms untuk endpoint surat
- API response time P95 < 500ms untuk endpoint search
- Zero data loss dalam seeding (114 surat, 6236 ayat, 30 juz, 604 halaman, 240 rub el hizb, 7 manzil)
- All endpoints documented di Scalar /docs
- Rate limiting aktif dan terbukti membatasi request berlebih
- Health check endpoints respond < 50ms

## Open Questions
- Rate limiting: jika Redis down, fallback ke allow all atau reject all? (diset ke allow all untuk now)
- Apakah perlu endpoint /languages untuk list bahasa yang tersedia?

## File Structure (Suggested)
```
quran-api-go/
├── cmd/api/main.go
├── internal/
│   ├── config/
│   ├── domain/
│   │   ├── surah/
│   │   ├── ayah/
│   │   ├── juz/
│   │   ├── page/
│   │   ├── hizb/
│   │   ├── manzil/
│   │   └── ruku/
│   ├── handler/
│   ├── repository/
│   ├── service/
│   ├── middleware/
│   │   ├── cors.go
│   │   ├── ratelimit.go
│   │   ├── logging.go
│   │   └── metrics.go
│   └── pkg/
│       ├── db/
│       ├── redis/
│       └── logger/
├── migrations/
├── scripts/
│   └── seed/
├── docs/
│   └── openapi.yaml
├── .env.example
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── go.mod
```

## Environment Variables
```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=quran
DB_PASSWORD=quran
DB_NAME=quran_db

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# Rate Limiting
RATE_LIMIT_PER_MINUTE=100

# Server
SERVER_PORT=8080
SERVER_HOST=0.0.0.0

# App
APP_VERSION=1.0.0
LOG_LEVEL=info
```
