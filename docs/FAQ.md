# FAQ — Quran API Go

Pertanyaan yang sering ditanya developer.

---

### 1. Kenapa gak pake PostgreSQL? SQLite kan untuk local dev

Itu misconception umum. SQLite itu database production-ready, bukan cuma untuk development. Bukannya:
- SQLite dipakai di banyak app production (macOS, iOS, Android, banyak browser)
- Untuk use case read-only seperti ini, SQLite lebih cepat karena tidak ada network overhead
- 6.236 ayat itu cuma ~2-3 MB data — sangat kecil

PostgreSQL akan overkill: perlu setup server, connection pool, maintenance — padahal datanya statis. Kalau suatu saat butuh fitur write yang heavy (bookmark, user data, dsb), baru pikir migrate.

---

### 2. Repository Pattern buat CRUD sederhana? serious?

Iya, sounds like overkill. Tapi coba bayangkan kalau tidak pakai layering:

- Handler langsung query SQL — logic business dan database bercampur
- Mau bikin unit test handler? Harus setup database dulu
- Mau ganti cara query? Harus ubah di semua handler

Dengan Repository:
- Test bisa mock repository, tidak perlu database betulan
- Kalau suatu saat mau ganti SQLite ke sesuatu yang lain, tinggal ganti implementation repository-nya aja
- SQL query logic terkonsentrasi di satu tempat

Untuk project yang bakal dikerjain banyak orang, struktur seperti ini mempermudah onboarding.

---

### 3. Search Arabic pakai FTS5 — works really?

Masalahnya: Arabic text punya complexity yang FTS5 default tangkap nggak sempurna.

Contoh masalah:
- Search "rahman" mungkin tidak match "ar-rahman" atau "الرحمن"
- Diacritics (tasykil) bisa bikin hasil nggak konsisten
- Prefix "al-" (ال) bisa bikin mismatch

Untuk MVP masih usable, tapi jangan expect perfect. Kalau mau search yang bener-bener bagus untuk Arabic, perlu:
1. Normalisasi text saat seeding (hilangkan tasykil, handle prefix)
2. Atau pakai search engine yang proper (Elasticsearch/Typesense)

---

### 4. API ini bisa handle berapa request per detik?

Jawaban jujur: belum ada benchmark resmi.

Tapi secara teori dan pengalaman:
- SQLite di WAL mode handle concurrent read tanpa masalah
- Bottleneck biasanya di CPU (JSON marshaling), bukan di I/O database
- Untuk dataset sekecil ini, 1000+ RPS seharusnya achievable di hardware standar

Yang pasti: sebelum deploy ke production, **WAJIB** load test dulu. Jangan asal tebak. Pakai vegeta, k6, atau Apache Bench — dapat angka real di environment production yang sebenarnya.

---

### 5. Kalau database corrupt gimana?

Worst case: API down sampai database diganti.

Corruption biasanya terjadi karena:
- Power failure pas lagi write (tapi di sini kan write-nya cuma sekali saat seed)
- Filesystem corrupt (rare di modern OS)
- Naruh database di network storage (NFS) — inirecipe bencana

Mitigasi yang realistis:
1. Backup database secara rutin (cron job tiap malam)
2. Health check yang cek `PRAGMA integrity_check`
3. Pertimbangkan bundling database dalam Docker image — jadi file-nya immutable

---

### 6. Manual DI gak scalable? Kenapa gak Wire/FX?

Scalable dalam arti bisa handle banyak code? Iya, scalable.

Pertanyaan sebenarnya: kapan perlu Wire/FX? Jawabnya: ketika constructor chain udah terlalu panjang dan rumit (5+ level, banyak optional dependencies).

Untuk scale project ini (~10 handlers), manual DI:
- Lebih gampang dimengerti kontributor baru
- Stack trace jelas, no magic
- Build time lebih cepat

Kalau nanti projectnya beneran grow besar dan constructor chain udah bikin pusing, baru pertimbangkan Wire. Untuk sekarang, keep it simple.

---

### 7. Offset pagination lambat kalau data udah gede

Iya, offset-based memang lambat buat dataset besar. `LIMIT 20 OFFSET 10000` artinya scan 10,020 rows dulu.

Tapi untuk Quran API ini:
- Data static (read-only), jadi nggak ada masalah inconsistent result
- Max offset kecil — 6,236 rows / 20 per page = 313 pages
- User jarang banget scroll ke page 200-an

Kalau nanti data-nya jutaan rows, baru consider cursor-based pagination. Untuk sekarang, offset-based simpler dan cukup fast.

---

### 8. Gimana cara test repository? Mock atau database betulan?

Pendekatan yang dipakai: in-memory SQLite, bukan mock.

```go
db, _ := sql.Open("sqlite", ":memory:")
goose.Up(db, "../../migrations")  // Run schema
```

Kenapa?
- Test actual SQL query, bukan mock behavior
- Test migration schema
- Test itu nyata, bukan pretend

Trade-off: test lebih lambat daripada pure mock. Tapi untuk repository layer, reliability lebih penting daripada speed test. Kalau mau super fast, unit test bisa mock service layer, bukan repository.

---

### 9. CORS — bahaya banget kan kalau salah config?

Iya, bisa jadi masalah.

Kalau `ALLOWED_ORIGINS=*`:
- Semua website bisa call API kamu dari browser user
- Bisa dipakai untuk scraping atau DDoS

Kalau specific origin:
```bash
ALLOWED_ORIGINS=https://ilmunara.com
```
- Hanya domain yang terdaftar yang bisa akses via browser

Catatan penting: CORS **hanya** apply di browser. curl, Postman, atau server-to-server request tetap bisa hit API langsung tanpa ada restriction. Kalau butuh security beneran, perlu auth token/API key — tapi itu out of scope MVP.

---

### 10. Deployment zero-downtime — gimana cara nya kalau SQLite pakai file lock?

SQLite pakai file lock, jadi multiple instance nggak bisa baca file yang sama kalau ada write operation.

Strategi realistis:

1. **Blue-green** — Run version baru di port lain, test, then switch traffic di load balancer. Ada brief downtime (few seconds) saat switch.

2. **Immutable infra** — Bundle database dalam Docker image. Setiap deploy = image baru. No file copying, no locking issue.

3. **Accept downtime** — Untuk internal API, 5-10 seconds downtime mungkin acceptable.

Yang **nggak** bisa: rolling update yang gradual (karena file lock issue). Kalau mau zero-downtime beneran, perlu database yang support multi-reader proper (PostgreSQL, MySQL, dsb) — dan itu berarti migration dari SQLite.

---

*Punya pertanyaan lain? Open issue di GitHub.*
