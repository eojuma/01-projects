CREATE TABLE IF NOT EXISTS post_likes (
    post_id    TEXT NOT NULL,
    user_id    TEXT NOT NULL,
    created_at INTEGER NOT NULL,
    PRIMARY KEY (post_id, user_id)
);
