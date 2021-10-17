package constants

const (
	SchemaTXS = `
	CREATE TABLE IF NOT EXISTS last_txs
	(
		blockid         int PRIMARY KEY        NOT NULL,
		data  JSONB
	);
	`
)
