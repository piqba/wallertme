package constants

const (
	SchemaMTTS = `
	CREATE TABLE IF NOT EXISTS last_txs
	(
		blockid         int PRIMARY KEY        NOT NULL,
		data  JSONB
	);
	`
)
