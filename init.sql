CREATE TABLE mappings (
	id	BIGSERIAL,
	short	VARCHAR(32),
	full	VARCHAR(2048),

	CONSTRAINT mappings_short_key UNIQUE(short)
);
