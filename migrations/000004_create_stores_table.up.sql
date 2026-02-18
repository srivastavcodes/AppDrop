CREATE TABLE IF NOT EXISTS stores
(
    id        UUID PRIMARY KEY,
    name      VARCHAR(255)                NOT NULL,
    slug      VARCHAR(100)                NOT NULL UNIQUE,
    createdAt TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updatedAt TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);

ALTER TABLE pages
    ADD CONSTRAINT fk_pages_store FOREIGN KEY (store_id) REFERENCES stores (id) ON DELETE CASCADE;
