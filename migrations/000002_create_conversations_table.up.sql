CREATE TABLE conversations (
                               id BIGSERIAL PRIMARY KEY,
                               created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                               updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);