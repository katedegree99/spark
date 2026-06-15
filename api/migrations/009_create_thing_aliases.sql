CREATE TABLE IF NOT EXISTS thing_aliases (
    thing_id BIGINT UNSIGNED NOT NULL,
    alias    VARCHAR(100)    NOT NULL,
    PRIMARY KEY (thing_id, alias),
    CONSTRAINT fk_thing_aliases_thing FOREIGN KEY (thing_id) REFERENCES things(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
