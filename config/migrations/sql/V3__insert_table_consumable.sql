CREATE TABLE consumable (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    uid UUID NOT NULL DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    code TEXT NOT NULL,
    description TEXT,
    price NUMERIC(10,2) NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);