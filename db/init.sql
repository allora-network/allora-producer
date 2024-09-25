
CREATE TABLE processed_blocks (
	id SERIAL PRIMARY KEY,
	height BIGINT NOT NULL,
	processed_at TIMESTAMP NOT NULL DEFAULT now(),
	status TEXT NOT NULL
);
CREATE INDEX ON processed_blocks (height);
CREATE INDEX ON processed_blocks (processed_at);
CREATE INDEX ON processed_blocks (status);

CREATE TABLE processed_block_events (
	id SERIAL PRIMARY KEY,
	height BIGINT NOT NULL,
	processed_at TIMESTAMP NOT NULL DEFAULT now(),
	status TEXT NOT NULL
);
CREATE INDEX ON processed_block_events (height);
CREATE INDEX ON processed_block_events (processed_at);
CREATE INDEX ON processed_block_events (status);
