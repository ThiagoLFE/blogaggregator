-- +goose Up
CREATE TABLE feeds (
    id UUID PRIMARY KEY DEFAULT uuidv7(),

    name TEXT NOT NULL,
    url TEXT NOT NULL,
    user_id UUID NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

ALTER TABLE feeds ADD CONSTRAINT fk_feeds_user
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE;

ALTER TABLE feeds ADD CONSTRAINT uq_feeds_user_name
    UNIQUE (user_id, name);

-- +goose Down
ALTER TABLE feeds DROP CONSTRAINT fk_feeds_user;
DROP TABLE feeds;
