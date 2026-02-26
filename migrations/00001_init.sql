-- +goose Up
CREATE TABLE IF NOT EXISTS surahs (
	id INTEGER PRIMARY KEY,
	number INTEGER NOT NULL,
	name_arabic TEXT NOT NULL,
	name_latin TEXT NOT NULL,
	name_transliteration TEXT NOT NULL,
	number_of_ayahs INTEGER NOT NULL,
	revelation_type TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS ayahs (
	id INTEGER PRIMARY KEY,
	surah_id INTEGER NOT NULL,
	number_in_surah INTEGER NOT NULL,
	text_uthmani TEXT NOT NULL,
	translation_indo TEXT NOT NULL,
	translation_en TEXT NOT NULL,
	juz_number INTEGER NOT NULL,
	sajda_type TEXT,
	revelation_type TEXT NOT NULL,
	FOREIGN KEY (surah_id) REFERENCES surahs(id)
);

CREATE TABLE IF NOT EXISTS juzs (
	id INTEGER PRIMARY KEY,
	juz_number INTEGER NOT NULL,
	first_ayah_id INTEGER NOT NULL,
	last_ayah_id INTEGER NOT NULL,
	FOREIGN KEY (first_ayah_id) REFERENCES ayahs(id),
	FOREIGN KEY (last_ayah_id) REFERENCES ayahs(id)
);

CREATE VIRTUAL TABLE IF NOT EXISTS ayahs_fts USING fts5(
	ayah_id UNINDEXED,
	text_uthmani,
	translation_indo,
	translation_en,
	content='ayahs',
	content_rowid='id'
);

CREATE INDEX IF NOT EXISTS idx_ayahs_surah_id ON ayahs (surah_id);
CREATE INDEX IF NOT EXISTS idx_ayahs_juz_number ON ayahs (juz_number);

-- +goose Down
DROP INDEX IF EXISTS idx_ayahs_juz_number;
DROP INDEX IF EXISTS idx_ayahs_surah_id;

DROP TABLE IF EXISTS ayahs_fts;
DROP TABLE IF EXISTS juzs;
DROP TABLE IF EXISTS ayahs;
DROP TABLE IF EXISTS surahs;
