CREATE TABLE IF NOT EXISTS widgets
(
    id         UUID PRIMARY KEY,
    page_id    UUID                        NOT NULL REFERENCES pages (id) ON DELETE CASCADE,
    type       VARCHAR(50)                 NOT NULL,
    position   INTEGER                     NOT NULL,
    config     JSONB,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT valid_widget_type CHECK ( type IN ('banner', 'product_grid', 'text', 'image', 'spacer'))
);

CREATE INDEX idx_widgets_page_position ON widgets (page_id, position);

CREATE INDEX idx_widgets_page_id ON widgets (page_id);
