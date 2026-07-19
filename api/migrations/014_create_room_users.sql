CREATE TABLE room_users (
    room_id    BIGINT UNSIGNED NOT NULL,
    user_id    BIGINT UNSIGNED NOT NULL,
    created_at DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (room_id, user_id),
    INDEX idx_room_users_user_id (user_id),
    CONSTRAINT fk_room_users_room FOREIGN KEY (room_id) REFERENCES rooms (id) ON DELETE CASCADE,
    CONSTRAINT fk_room_users_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
