CREATE TABLE messages (
    id             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    room_id        BIGINT UNSIGNED NOT NULL,
    sender_user_id BIGINT UNSIGNED NOT NULL,
    content        TEXT            NOT NULL,
    created_at     DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    INDEX idx_messages_room_id (room_id),
    CONSTRAINT fk_messages_room   FOREIGN KEY (room_id)        REFERENCES rooms (id) ON DELETE CASCADE,
    CONSTRAINT fk_messages_sender FOREIGN KEY (sender_user_id) REFERENCES users (id) ON DELETE CASCADE
);
