CREATE TABLE user_interests (
    from_user_id BIGINT UNSIGNED NOT NULL,
    to_user_id   BIGINT UNSIGNED NOT NULL,
    created_at   DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (from_user_id, to_user_id),
    INDEX idx_user_interests_to_user_id (to_user_id)
);
