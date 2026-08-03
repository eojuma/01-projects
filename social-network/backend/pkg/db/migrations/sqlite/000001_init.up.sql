-- ================================================
-- USERS & AUTH
-- ================================================
-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT UNIQUE NOT NULL,
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    first_name TEXT,
    last_name TEXT,
    bio TEXT,
    avatar_url TEXT,
    nickname TEXT,
    date_of_birth TEXT,
    is_private INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Create sessions table
CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL,
    expires_at DATETIME NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);


-- ================================================
-- SOCIAL GRAPH
-- ================================================

CREATE TABLE IF NOT EXISTS follow_requests (
    id           TEXT PRIMARY KEY,
    requester_id TEXT NOT NULL REFERENCES users(id),
    target_id    TEXT NOT NULL REFERENCES users(id),
    status       TEXT DEFAULT 'pending'
                   CHECK(status IN ('pending','accepted','declined')),
    created_at   INTEGER NOT NULL,
    UNIQUE(requester_id, target_id)
);

CREATE TABLE IF NOT EXISTS followers (
    follower_id  TEXT NOT NULL REFERENCES users(id),
    following_id TEXT NOT NULL REFERENCES users(id),
    created_at   INTEGER NOT NULL,
    PRIMARY KEY (follower_id, following_id)
);

-- ================================================
-- POSTS
-- ================================================

CREATE TABLE IF NOT EXISTS posts (
    id         TEXT PRIMARY KEY,
    author_id  TEXT NOT NULL REFERENCES users(id),
    content    TEXT NOT NULL,
    privacy    TEXT DEFAULT 'public'
                 CHECK(privacy IN ('public','followers','selected')),
    created_at INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS post_images (
    id         TEXT PRIMARY KEY,
    post_id    TEXT NOT NULL REFERENCES posts(id),
    image_path TEXT NOT NULL,
    created_at INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS post_permissions (
    post_id    TEXT NOT NULL REFERENCES posts(id),
    user_id    TEXT NOT NULL REFERENCES users(id),
    PRIMARY KEY (post_id, user_id)
);

CREATE TABLE IF NOT EXISTS comments (
    id         TEXT PRIMARY KEY,
    post_id    TEXT NOT NULL REFERENCES posts(id),
    author_id  TEXT NOT NULL REFERENCES users(id),
    content    TEXT NOT NULL,
    image_path TEXT,
    created_at INTEGER NOT NULL
);

-- ================================================
-- GROUPS
-- ================================================

CREATE TABLE IF NOT EXISTS groups (
    id          TEXT PRIMARY KEY,
    title       TEXT NOT NULL,
    description TEXT,
    creator_id  TEXT NOT NULL REFERENCES users(id),
    created_at  INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS group_members (
    group_id   TEXT NOT NULL REFERENCES groups(id),
    user_id    TEXT NOT NULL REFERENCES users(id),
    joined_at  INTEGER NOT NULL,
    PRIMARY KEY (group_id, user_id)
);

CREATE TABLE IF NOT EXISTS group_invites (
    id         TEXT PRIMARY KEY,
    group_id   TEXT NOT NULL REFERENCES groups(id),
    inviter_id TEXT NOT NULL REFERENCES users(id),
    invitee_id TEXT NOT NULL REFERENCES users(id),
    status     TEXT DEFAULT 'pending'
                 CHECK(status IN ('pending','accepted','declined')),
    created_at INTEGER NOT NULL,
    UNIQUE(group_id, invitee_id)
);

CREATE TABLE IF NOT EXISTS group_join_requests (
    id         TEXT PRIMARY KEY,
    group_id   TEXT NOT NULL REFERENCES groups(id),
    user_id    TEXT NOT NULL REFERENCES users(id),
    status     TEXT DEFAULT 'pending'
                 CHECK(status IN ('pending','approved','declined')),
    created_at INTEGER NOT NULL,
    UNIQUE(group_id, user_id)
);

CREATE TABLE IF NOT EXISTS group_posts (
    id         TEXT PRIMARY KEY,
    group_id   TEXT NOT NULL REFERENCES groups(id),
    author_id  TEXT NOT NULL REFERENCES users(id),
    content    TEXT NOT NULL,
    created_at INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS group_comments (
    id         TEXT PRIMARY KEY,
    post_id    TEXT NOT NULL REFERENCES group_posts(id),
    author_id  TEXT NOT NULL REFERENCES users(id),
    content    TEXT NOT NULL,
    created_at INTEGER NOT NULL
);

-- ================================================
-- EVENTS
-- ================================================

CREATE TABLE IF NOT EXISTS events (
    id          TEXT PRIMARY KEY,
    group_id    TEXT NOT NULL REFERENCES groups(id),
    creator_id  TEXT NOT NULL REFERENCES users(id),
    title       TEXT NOT NULL,
    description TEXT,
    event_date  INTEGER NOT NULL,
    created_at  INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS event_responses (
    id         TEXT PRIMARY KEY,
    event_id   TEXT NOT NULL REFERENCES events(id),
    user_id    TEXT NOT NULL REFERENCES users(id),
    response   TEXT NOT NULL CHECK(response IN ('going','not_going')),
    created_at INTEGER NOT NULL,
    UNIQUE(event_id, user_id)
);

-- ================================================
-- NOTIFICATIONS
-- ================================================

CREATE TABLE IF NOT EXISTS notifications (
    id         TEXT PRIMARY KEY,
    user_id    TEXT NOT NULL REFERENCES users(id),
    type       TEXT NOT NULL,
    content    TEXT NOT NULL,
    is_read    INTEGER DEFAULT 0,
    ref_id     TEXT,
    created_at INTEGER NOT NULL
);

-- ================================================
-- CHAT
-- ================================================

CREATE TABLE IF NOT EXISTS conversations (
    id         TEXT PRIMARY KEY,
    created_at INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS conversation_participants (
    conversation_id TEXT NOT NULL REFERENCES conversations(id),
    user_id         TEXT NOT NULL REFERENCES users(id),
    PRIMARY KEY (conversation_id, user_id)
);

CREATE TABLE IF NOT EXISTS messages (
    id              TEXT PRIMARY KEY,
    conversation_id TEXT NOT NULL REFERENCES conversations(id),
    sender_id       TEXT NOT NULL REFERENCES users(id),
    content         TEXT NOT NULL,
    is_read         INTEGER DEFAULT 0,
    created_at      INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS group_messages (
    id         TEXT PRIMARY KEY,
    group_id   TEXT NOT NULL REFERENCES groups(id),
    sender_id  TEXT NOT NULL REFERENCES users(id),
    content    TEXT NOT NULL,
    created_at INTEGER NOT NULL
);


-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions(expires_at);
