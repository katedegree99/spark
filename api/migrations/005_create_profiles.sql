CREATE TABLE IF NOT EXISTS profiles (
    user_id       BIGINT UNSIGNED NOT NULL,
    name          VARCHAR(100)    NOT NULL,
    icon_image_id BIGINT UNSIGNED NULL,
    bio           TEXT            NULL,
    created_at    DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at    DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (user_id),
    CONSTRAINT fk_profiles_user  FOREIGN KEY (user_id)       REFERENCES users(id)  ON DELETE CASCADE,
    CONSTRAINT fk_profiles_image FOREIGN KEY (icon_image_id) REFERENCES images(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
