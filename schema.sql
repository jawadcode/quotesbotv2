CREATE TABLE IF NOT EXISTS quotes (
	id SERIAL,
	guild_id BIGINT NOT NULL,
	channel_id BIGINT NOT NULL,
	message_id BIGINT NOT NULL,
	content text NOT NULL,
	tsv TSVECTOR,
	author BIGINT NOT NULL,
	added_by BIGINT NOT NULL,
	added_at BIGINT NOT NULL
);

DROP TRIGGER IF EXISTS tsvectorupdate ON quotes;

CREATE TRIGGER tsvectorupdate BEFORE INSERT OR UPDATE
ON quotes FOR EACH ROW EXECUTE PROCEDURE
tsvector_update_trigger(tsv, 'pg_catalog.english', content);
