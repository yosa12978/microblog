CREATE TABLE IF NOT EXISTS posts (
    id BIGSERIAL,
    content VARCHAR NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),

    PRIMARY KEY(id)
);
