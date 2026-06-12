-- +goose Up
-- Recreate ayahs_fts without the invalid `ayah_id UNINDEXED` column.
-- The original schema listed ayah_id as a content column, but the `ayahs`
-- table has no such column (it uses `id`), causing every FTS query to fail
-- with "no such column: T.ayah_id". Removing it and rebuilding from the
-- content table fixes the issue without any data loss.
DROP TABLE IF EXISTS ayahs_fts;

CREATE VIRTUAL TABLE IF NOT EXISTS ayahs_fts USING fts5(
	text_uthmani,
	translation_indo,
	translation_en,
	content='ayahs',
	content_rowid='id'
);

-- Rebuild the FTS5 index from the ayahs content table.
INSERT INTO ayahs_fts(ayahs_fts) VALUES('rebuild');

-- +goose Down
DROP TABLE IF EXISTS ayahs_fts;

CREATE VIRTUAL TABLE IF NOT EXISTS ayahs_fts USING fts5(
	ayah_id UNINDEXED,
	text_uthmani,
	translation_indo,
	translation_en,
	content='ayahs',
	content_rowid='id'
);
