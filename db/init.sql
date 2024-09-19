
CREATE TABLE processed_block_transactions (
	id SERIAL PRIMARY KEY,
	height BIGINT NOT NULL,
	processed_at TIMESTAMP NOT NULL DEFAULT now(),
	status TEXT NOT NULL
);
CREATE INDEX ON processed_block_transactions (height);
CREATE INDEX ON processed_block_transactions (processed_at);
CREATE INDEX ON processed_block_transactions (status);

CREATE TABLE processed_block_events (
	id SERIAL PRIMARY KEY,
	height BIGINT NOT NULL,
	processed_at TIMESTAMP NOT NULL DEFAULT now(),
	status TEXT NOT NULL
);
CREATE INDEX ON processed_block_events (height);
CREATE INDEX ON processed_block_events (processed_at);
CREATE INDEX ON processed_block_events (status);
