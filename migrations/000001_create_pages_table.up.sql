CREATE TABLE IF NOT EXISTS pages
(
    id         UUID PRIMARY KEY,
    -- app_id mimics the scenario where pages table will be linked to a foreign key of the app or similar
    app_id     UUID                        NOT NULL,
    name       VARCHAR(255)                NOT NULL,
    route      VARCHAR(255)                NOT NULL,
    is_home    BOOLEAN                     NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE (app_id, route)
);

CREATE UNIQUE INDEX idx_pages_is_home_per_app ON pages (app_id) WHERE is_home = TRUE;
