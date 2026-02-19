# PRD: Quran API Go

## Overview
Membangun RESTful API publik untuk Al-Quran menggunakan Golang dan Gin framework. API menyediakan endpoints untuk mengakses data surat, ayat dengan terjemahan Indonesia & English, serta fitur pencarian teks dengan filter surat/juz. Data sudah tersedia dan siap di-import.

## Goals
- Menyediakan API publik yang cepat dan ringan untuk mengakses data Al-Quran
- Mendukung 2 bahasa terjemahan: Indonesia & English
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
- [ ] Create migration untuk tabel `ayahs`: id, surah_id, number_in_surah, text_uthmani, translation_indo, translation_en, juz_number, page_number, sajda_type, revelation_type
- [ ] Create migration untuk tabel `juzs`: id, juz_number, first_ayah_surah_id, last_ayah_surah_id
- [ ] Add indexes untuk pencarian: surah_id, juz_number, text_uthmani (GIN/trigram index untuk full-text search)
- [ ] Migration harus reversible (up/down)

### US-003: Seed data dari dataset yang sudah disiapkan
**Description:** As a developer, I want to seed Quran data dari dataset yang sudah ada so that API bisa return data yang valid.

**Acceptance Criteria:**
- [ ] Create seeder script untuk import data dari dataset yang sudah disiapkan
- [ ] Import semua 114 surat
- [ ] Import semua 6236 ayat dengan text_uthmani, translation_indo, translation_en
- [ ] Import data 30 juz
- [ ] Seeder harus idempotent (bisa di-run berkali-kali)
- [ ] Add logging progress saat seeding
- [ ] Validasi data kelengkapan setelah seeding (count surahs = 114, count ayahs = 6236)

### US-004: Endpoint GET /surah - List semua surat
**Description:** As an API consumer, I want to get list of all surahs so that I can browse Al-Quran content.

**Acceptance Criteria:**
- [ ] GET /surah return array semua surat
- [ ] Support query parameter: `?page=1&limit=10` untuk pagination
- [ ] Response structure: `[{ id, number, name_arabic, name_latin, name_transliteration, number_of_ayahs, revelation_type }]`
- [ ] Return HTTP 200 dengan proper headers (Content-Type: application/json)
- [ ] Return HTTP 500 jika database error

### US-005: Endpoint GET /surah/:id - Detail surat dengan ayat
**Description:** As an API consumer, I want to get detail surat with all ayahs so that I can read complete surah content.

**Acceptance Criteria:**
- [ ] GET /surah/:id return detail surat dengan array ayat
- [ ] Support query parameter: `?lang=id` atau `?lang=en` untuk pilih bahasa terjemahan (default: id)
- [ ] Response structure: `{ surah: { id, number, name_arabic, name_latin, number_of_ayahs, revelation_type }, ayahs: [{ number, text_uthmani, translation, juz, page, sajda, revelation_type }] }`
- [ ] Return HTTP 400 jika `lang` tidak valid (hanya id/en)
- [ ] Return HTTP 404 jika surat tidak ditemukan
- [ ] Response time < 100ms untuk surat pendek (< 50 ayat)

### US-006: Endpoint GET /ayah/:id - Detail ayat spesifik
**Description:** As an API consumer, I want to get specific ayah by ID so that I can reference individual ayah.

**Acceptance Criteria:**
- [ ] GET /ayah/:id return detail ayat
- [ ] Support query parameter: `?lang=id` atau `?lang=en`
- [ ] Response structure: `{ id, surah_id, number, text_uthmani, translation, surah_info: { id, name_latin }, juz, page, sajda, revelation_type }`
- [ ] Return HTTP 400 jika `lang` tidak valid
- [ ] Return HTTP 404 jika ayat tidak ditemukan

### US-007: Endpoint GET /juz - List dan detail juz
**Description:** As an API consumer, I want to get juz information so that I can navigate Al-Quran by juz.

**Acceptance Criteria:**
- [ ] GET /juz return array semua juz (1-30)
- [ ] Response structure: `[{ juz_number, first_ayah, last_ayah, total_ayahs }]`
- [ ] GET /juz/:number return detail ayat-ayat dalam juz tersebut
- [ ] Support `?lang=id` atau `?lang=en` untuk endpoint detail juz
- [ ] Return HTTP 404 jika juz tidak valid (bukan 1-30)

### US-008: Endpoint GET /search - Pencarian ayat dengan filter
**Description:** As an API consumer, I want to search ayahs by keyword with filters so that I can find specific content.

**Acceptance Criteria:**
- [ ] GET /search?q=keyword untuk pencarian full-text
- [ ] Support filter: `?surah_id=1` untuk limit ke surat tertentu
- [ ] Support filter: `?juz=1` untuk limit ke juz tertentu
- [ ] Support filter: `?lang=id` atau `?lang=en` untuk bahasa terjemahan (default: id)
- [ ] Query parameter: `?page=1&limit=20` untuk pagination result
- [ ] Pencarian menggunakan PostgreSQL ILIKE atau tsvector untuk case-insensitive search
- [ ] Response structure: `{ query, total, page, limit, results: [{ id, surah_info, number, text_uthmani, translation, juz, page }] }`
- [ ] Return HTTP 400 jika parameter `q` kosong

### US-009: Rate limiting middleware dengan Redis
**Description:** As an API owner, I want to implement rate limiting by IP using Redis so that API tidak disalahgunakan.

**Acceptance Criteria:**
- [ ] Implement rate limiting middleware dengan Redis backend
- [ ] Limit: 100 request per menit per IP
- [ ] Return HTTP 429 Too Many Requests dengan headers: `X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`, `Retry-After`
- [ ] Redis connection menggunakan connection pooling
- [ ] Configurable limit lewat environment variable `RATE_LIMIT_PER_MINUTE`
- [ ] Handle Redis connection failure gracefully (fallback to allow all atau reject all)

### US-010: API Documentation dengan Scalar
**Description:** As an API consumer, I want to access interactive API documentation so that I can explore endpoints easily.

**Acceptance Criteria:**
- [ ] Setup Scalar untuk API documentation
- [ ] Serve Scalar UI di endpoint /docs
- [ ] Generate OpenAPI 3.0 spec dari annotations
- [ ] Document semua endpoint dengan method, path, parameters, response examples
- [ ] Include contoh response untuk setiap endpoint
- [ ] Scalar UI harus standalone (bisa berjalan tanpa build steps tambahan)

### US-011: Observability - Metrics endpoint
**Description:** As an operator, I want to expose metrics so that I can monitor API health and performance.

**Acceptance Criteria:**
- [ ] Expose Prometheus metrics at /metrics
- [ ] Track: request count, request duration, error rate by endpoint
- [ ] Track database connection pool stats
- [ ] Track Redis connection stats
- [ ] Track rate limiting: blocked requests count

### US-012: Health check endpoint
**Description:** As an operator/deployer, I want to check API health so that I can monitor service availability.

**Acceptance Criteria:**
- [ ] GET /health return status API
- [ ] Response structure: `{ status: "ok", timestamp: "...", version: "..." }`
- [ ] GET /health/ready untuk readiness probe (check DB & Redis connection)
- [ ] GET /health/live untuk liveness probe
- [ ] Return HTTP 503 jika dependency tidak ready

### US-013: CORS dan Security headers
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

## Non-Goals (Out of Scope)
- Autentikasi API Key (future)
- Multi-language translations selain Indo/English (Phase 2: Melayu, dll)
- Audio recitation endpoints (Phase 2)
- Tafsir endpoints (Phase 2)
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
- Zero data loss dalam seeding (114 surat, 6236 ayat)
- All endpoints documented di Scalar /docs
- Rate limiting aktif dan terbukti membatasi request berlebih
- Health check endpoints respond < 50ms

## Open Questions
- Rate limiting: jika Redis down, fallback ke allow all atau reject all?
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
│   │   └── juz/
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
