CREATE TABLE service_definition (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    uid UUID NOT NULL DEFAULT uuid_generate_v4(),
    price NUMERIC(10,2) NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);