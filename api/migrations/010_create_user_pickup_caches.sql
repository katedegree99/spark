CREATE TABLE user_pickup_caches (
    user_id     BIGINT UNSIGNED NOT NULL,
    cache_date  DATE            NOT NULL,
    picked_user_ids JSON        NOT NULL,
    created_at  DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (user_id, cache_date)
);
