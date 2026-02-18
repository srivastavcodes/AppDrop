ALTER TABLE widgets
    ADD CONSTRAINT unique_widgets_pos_per_page UNIQUE (page_id, position) DEFERRABLE INITIALLY DEFERRED;
