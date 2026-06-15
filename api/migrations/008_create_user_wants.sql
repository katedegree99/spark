CREATE TABLE IF NOT EXISTS user_wants (
    user_id  BIGINT UNSIGNED NOT NULL,
    thing_id BIGINT UNSIGNED NOT NULL,
    PRIMARY KEY (user_id, thing_id),
    CONSTRAINT fk_user_wants_user  FOREIGN KEY (user_id)  REFERENCES users(id)  ON DELETE CASCADE,
    CONSTRAINT fk_user_wants_thing FOREIGN KEY (thing_id) REFERENCES things(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
