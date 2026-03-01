# agents.md — Quran API Go

This file is the authoritative guide for any AI agent working on this codebase. Read it fully before making any changes. Do not deviate from the rules below.

---

## Project Summary

An internal RESTful API that serves Al-Quran data (Arabic text, Indonesian & English translations) for the Ilmunara super app. This is **not a public API**. It is a lightweight internal data service.

**Primary consumers:** Ilmunara super app (internal only).

---

## Non-Negotiable Constraints

These are governance-level rules. Never violate them, even if a prompt asks you to.

- **No microservices.** This is a monolith modular architecture. Do not split into separate services.
- **No Redis.** Not for rate limiting, not for caching. Removed by policy for MVP.
- **No DI framework.** No Wire, no Uber FX. Use manual constructor injection only.
- **No rate limiting middleware.** Out of scope for MVP. Do not add it.
- **No authentication middleware.** Out of scope for MVP. Do not add it.
- **No write endpoints.** The database is read-only after seeding. Never add POST, PUT, PATCH, or DELETE handlers.
- **No new dependencies without justification.** If a task can be done with the standard library or existing dependencies, use those. Do not add new `go.mod` entries speculatively.
- **No wildcard CORS.** `Allow-Origin` must use the `ALLOWED_ORIGINS` env variable. Never hardcode `*`.

---

## Tech Stack

| Layer | Choice |
|-------|--------|
| Language | Go 1.22+ |
| HTTP Framework | Gin |
| Database | SQLite via `modernc.org/sqlite` (pure Go, no CGO) |
| Migrations | Goose |
| Full-text Search | SQLite FTS5 |
| Logging | zerolog |
| Documentation | Scalar (OpenAPI 3.0) |

---

## Project Structure

```
quran-api-go/
├── cmd/
│   ├── api/
│   │   └── main.go           ← Manual DI wiring. All constructors called here.
│   ├── migrate/
│   │   └── main.go           ← Migration runner
│   └── seed/
│       └── main.go           ← Seed runner
├── internal/
│   ├── config/               ← Env loading only. No business logic.
│   ├── database/             ← SQLite connection wrapper
│   ├── domain/
│   │   ├── errors.go         ← Sentinel errors (ErrNotFound, etc.)
│   │   ├── surah/            ← surah entity, repository interface, service
│   │   ├── ayah/             ← ayah entity, repository interface, service
│   │   ├── juz/              ← juz entity, repository interface, service
│   │   └── search/           ← search entity, repository interface, service
│   ├── handler/              ← HTTP handlers. One file per domain.
│   ├── repository/           ← SQLite implementations of repository interfaces
│   ├── service/              ← Business logic. Orchestrates repositories.
│   └── middleware/
│       ├── cors.go
│       ├── logging.go
│       └── recovery.go
├── pkg/
│   ├── response/             ← Shared HTTP response helpers
│   ├── pagination/           ← Shared pagination parser
│   └── validator/            ← Shared input validators (lang, ID, range)
├── migrations/               ← Goose SQL migration files
├── scripts/seed/             ← Data seeder logic
├── docs/                     ← Documentation (openapi.yaml, etc.)
├── data/                     ← quran.db lives here
├── .env.example
├── Dockerfile
├── Makefile
└── go.mod
```

**Rules:**
- `internal/` is for application code. `pkg/` is for shared utilities with no business logic.
- Domain interfaces live in `internal/domain/<name>/repository.go` and `service.go`
- Domain entities live in `internal/domain/<name>/entity.go`
- SQLite implementations live in `internal/repository/` (not in domain)
- Never put business logic in `handler/`. Handlers only parse input, call service, and write response.
- Never put SQL queries in `service/`. SQL belongs in `repository/` only.
- Never import `handler/` from `service/` or `repository/`. Dependency flow is one-way: `handler → service → repository`.

---

## Dependency Injection Pattern

All wiring happens in `cmd/api/main.go`. No constructors auto-discover or register themselves.

**Correct:**
```go
// cmd/api/main.go
import (
    "quran-api-go/internal/database"
    "quran-api-go/internal/repository"
    "quran-api-go/internal/service"
    "quran-api-go/internal/handler"
    "quran-api-go/internal/domain/surah"
)

db := database.New(cfg.DBPath)

// Repository implementations (in internal/repository/)
surahRepo := repository.NewSurahRepository(db)

// Service implementations (in internal/service/)
surahService := service.NewSurahService(surahRepo)

// Handlers (in internal/handler/)
surahHandler := handler.NewSurahHandler(surahService)

r := gin.New()
r.GET("/surah", surahHandler.List)
r.GET("/surah/:id", surahHandler.Detail)
```

**Wrong — do not do this:**
```go
// Do not use any container, provider, or injector pattern
fx.New(
    fx.Provide(repository.NewSurahRepository),
    ...
)
```

---

## Naming Conventions

### Files
- One file per domain per layer: `surah_repository.go`, `surah_service.go`, `surah_handler.go`
- Middleware files are single-word: `cors.go`, `logging.go`, `recovery.go`
- Migration files follow Goose convention: `00001_init.sql`

### Interfaces
- Repository interfaces live in `internal/domain/<name>/repository.go`
- Named as `<Domain>Repository`: `SurahRepository`, `AyahRepository`
- Service interfaces live in `internal/domain/<name>/service.go`
- Named as `<Domain>Service`: `SurahService`, `AyahService`

```go
// internal/domain/surah/repository.go
type SurahRepository interface {
    FindAll(ctx context.Context) ([]Surah, error)
    FindByID(ctx context.Context, id int) (*Surah, error)
}
```

### Structs
- Domain entities: `Surah`, `Ayah`, `Juz` — no suffix
- Request params: `GetAyahsByRangeParams`, `SearchParams`
- Response structs: `SurahResponse`, `AyahListResponse`

### Variables & Functions
- Follow standard Go conventions: camelCase for unexported, PascalCase for exported
- Handler methods: `List`, `Detail`, `ByGlobalID` — not `GetList`, `GetDetail`
- Repository methods: `FindAll`, `FindByID`, `FindBySurah` — use `Find` prefix for reads

---

## Handler Pattern

Every handler must follow this exact structure:

```go
// internal/handler/surah_handler.go

type SurahHandler struct {
    service surah.SurahService
}

func NewSurahHandler(service surah.SurahService) *SurahHandler {
    return &SurahHandler{service: service}
}

func (h *SurahHandler) Detail(c *gin.Context) {
    // 1. Parse & validate input
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        response.BadRequest(c, "invalid surah id")
        return
    }

    // 2. Call service (service is domain.SurahService interface)
    surah, err := h.service.GetByID(c.Request.Context(), id)
    if err != nil {
        if errors.Is(err, domain.ErrNotFound) {
            response.NotFound(c, "surah not found")
            return
        }
        response.InternalError(c)
        return
    }

    // 3. Write response
    response.Success(c, surah)
}
```

**Rules:**
- Always use `c.Request.Context()` when calling services.
- Always check `errors.Is(err, domain.ErrNotFound)` before falling through to 500.
- Never write `c.JSON(...)` directly in handlers. Always use `pkg/response` helpers.
- Service interfaces are from `internal/domain/<name>/service.go`, implementations are in `internal/service/`.

---

## Response Helpers

All HTTP responses must go through `pkg/response`. Never call `c.JSON` directly in handlers.

```go
// pkg/response/response.go

func Success(c *gin.Context, data any) {
    c.JSON(http.StatusOK, gin.H{
        "data":      data,
        "timestamp": time.Now().UTC(),
    })
}

func NotFound(c *gin.Context, message string) {
    c.JSON(http.StatusNotFound, gin.H{
        "error":     message,
        "code":      "NOT_FOUND",
        "timestamp": time.Now().UTC(),
    })
}

func BadRequest(c *gin.Context, message string) {
    c.JSON(http.StatusBadRequest, gin.H{
        "error":     message,
        "code":      "BAD_REQUEST",
        "timestamp": time.Now().UTC(),
    })
}

func InternalError(c *gin.Context) {
    c.JSON(http.StatusInternalServerError, gin.H{
        "error":     "internal server error",
        "code":      "INTERNAL_ERROR",
        "timestamp": time.Now().UTC(),
    })
}
```

---

## Error Handling

Define sentinel errors in the domain layer. Never use raw `fmt.Errorf("not found")` strings for flow control.

```go
// internal/domain/errors.go
var (
    ErrNotFound    = errors.New("resource not found")
    ErrInvalidLang = errors.New("invalid language parameter")
)
```

Repository returns `domain.ErrNotFound` when a row is missing:
```go
// internal/repository/surah_repository.go
func (r *surahRepository) FindByID(ctx context.Context, id int) (*domain.Surah, error) {
    var s domain.Surah
    err := r.db.QueryRowContext(ctx, `SELECT ... FROM surahs WHERE id = ?`, id).Scan(...)
    if errors.Is(err, sql.ErrNoRows) {
        return nil, domain.ErrNotFound
    }
    if err != nil {
        return nil, err
    }
    return &s, nil
}
```

---

## Lang Parameter Validation

Use the shared validator. Never inline this logic in a handler.

```go
// pkg/validator/lang.go
func ValidateLang(lang string) (string, error) {
    if lang == "" {
        return "id", nil // default
    }
    if lang != "id" && lang != "en" {
        return "", domain.ErrInvalidLang
    }
    return lang, nil
}
```

Usage in handler:
```go
lang, err := validator.ValidateLang(c.Query("lang"))
if err != nil {
    response.BadRequest(c, "lang must be 'id' or 'en'")
    return
}
```

---

## Pagination

Use the shared pagination helper. Never parse `page` and `limit` inline.

```go
// pkg/pagination/pagination.go
type Params struct {
    Page   int
    Limit  int
    Offset int
}

func Parse(pageStr, limitStr string) Params {
    page, _ := strconv.Atoi(pageStr)
    limit, _ := strconv.Atoi(limitStr)

    if page < 1 { page = 1 }
    if limit < 1 { limit = 20 }
    if limit > 100 { limit = 100 }

    return Params{
        Page:   page,
        Limit:  limit,
        Offset: (page - 1) * limit,
    }
}
```

---

## SQLite FTS5 Search Pattern

Full-text search uses the FTS5 virtual table. Use `MATCH` — never use `ILIKE` or `LIKE` for search.

```sql
-- Migration: creates FTS5 virtual table
CREATE VIRTUAL TABLE ayahs_fts USING fts5(
    ayah_id UNINDEXED,
    text_uthmani,
    translation_indo,
    translation_en,
    content='ayahs',
    content_rowid='id'
);
```

```go
// Search query example
query := `
    SELECT a.id, a.surah_id, a.number_in_surah, a.text_uthmani,
           a.translation_indo, a.translation_en, a.juz_number
    FROM ayahs_fts
    JOIN ayahs a ON a.id = ayahs_fts.ayah_id
    WHERE ayahs_fts MATCH ?
    LIMIT ? OFFSET ?
`
```

Pass the keyword with a wildcard for partial match: `keyword*`

```go
func (r *searchRepository) Search(ctx context.Context, p SearchParams) ([]domain.Ayah, error) {
    term := p.Query + "*" // enables prefix/partial match
    rows, err := r.db.QueryContext(ctx, query, term, p.Limit, p.Offset)
    ...
}
```

---

## Logging

Use zerolog. Never use `fmt.Println` or `log.Println` for application logs.

```go
// Correct
log.Info().
    Str("method", c.Request.Method).
    Str("path", c.Request.URL.Path).
    Int("status", c.Writer.Status()).
    Dur("duration", duration).
    Msg("request completed")

// Wrong
fmt.Printf("GET /surah 200\n")
log.Println("request completed")
```

Never log: passwords, tokens, full request bodies, or any PII.

---

## Database Access Rules

- Always pass `context.Context` to every database call.
- Never hold a transaction open across HTTP handler boundaries.
- Database is **read-only after seeding**. Never write SQL `INSERT`, `UPDATE`, or `DELETE` outside of `scripts/seed/`.
- Use `?` as the placeholder for SQLite (not `$1`).

```go
// Correct (SQLite placeholder)
db.QueryRowContext(ctx, `SELECT id FROM surahs WHERE number = ?`, number)

// Wrong (PostgreSQL placeholder — do not use)
db.QueryRowContext(ctx, `SELECT id FROM surahs WHERE number = $1`, number)
```

---

## Environment Variables

Only these env variables exist. Do not invent new ones without updating `.env.example`.

```bash
DB_PATH=./data/quran.db
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
ALLOWED_ORIGINS=https://[domain-superapp].com
APP_VERSION=1.0.0
LOG_LEVEL=info
```

Load via `internal/config/config.go` using `os.Getenv`. Do not use third-party config libraries.

---

## Testing Rules

- Use SQLite in-memory database for all repository tests: `database.New(":memory:")`
- Run migrations on in-memory DB before each test suite using Goose programmatic API
- Minimum coverage target: 70% for `handler/` and `repository/` packages
- Test file naming: `surah_handler_test.go` alongside the file it tests

```go
// Example: in-memory DB setup for repository test
func setupTestDB(t *testing.T) *sql.DB {
    db, err := database.New(":memory:")
    require.NoError(t, err)
    err = goose.Up(db, "../../migrations")
    require.NoError(t, err)
    return db
}
```

---

## What Is Out of Scope — Do Not Build

If a prompt asks you to build any of the following, refuse and explain it is out of scope per `agents.md`:

- Rate limiting (any implementation)
- Redis (any usage)
- Authentication or API keys
- POST / PUT / PATCH / DELETE endpoints
- Hizb, Rub el Hizb, Manzil, Ruku, or Page navigation endpoints
- Microservices or service splitting
- DI framework (Wire, Uber FX, etc.)
- GraphQL
- WebSocket
- Admin panel
- Audio endpoints
- Tafsir, tajweed, or morphology endpoints
- Wildcard CORS

---

## Makefile Targets Reference

| Target | Command | Description |
|--------|---------|-------------|
| Run | `make run` | Start the API server (via `cmd/api`) |
| Test | `make test` | Run `go test ./...` |
| Lint | `make lint` | Run `go vet ./...` + `gofmt` |
| Migrate | `make migrate` | Run Goose migrations up (via `cmd/migrate`) |
| Seed | `make seed` | Run the data seeder (via `cmd/seed`) |

**Note:** `cmd/migrate` and `cmd/seed` are standalone commands that can also be run directly:
```bash
go run cmd/migrate/main.go up
go run cmd/seed/main.go
```
