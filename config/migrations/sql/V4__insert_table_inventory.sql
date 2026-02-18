CREATE TABLE inventory (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    uid UUID NOT NULL DEFAULT uuid_generate_v4(),
    quantity INTEGER NOT NULL DEFAULT 0 CHECK (quantity >= 0),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    consumable_id BIGINT NOT NULL,
    CONSTRAINT fk_consumable
        FOREIGN KEY (consumable_id)
        REFERENCES consumable(id)
        ON DELETE RESTRICT,
    CONSTRAINT uq_consumable_inventory UNIQUE (consumable_id)
);