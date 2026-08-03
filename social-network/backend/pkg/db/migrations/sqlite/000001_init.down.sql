DROP TABLE IF EXISTS group_messages;
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS conversation_participants;
DROP TABLE IF EXISTS conversations;
DROP TABLE IF EXISTS notifications;
DROP TABLE IF EXISTS event_responses;
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS group_comments;
DROP TABLE IF EXISTS group_posts;
DROP TABLE IF EXISTS group_join_requests;
DROP TABLE IF EXISTS group_invites;
DROP TABLE IF EXISTS group_members;
DROP TABLE IF EXISTS groups;
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS post_permissions;
DROP TABLE IF EXISTS post_images;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS followers;
DROP TABLE IF EXISTS follow_requests;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS users;
-- Drop indexes
DROP INDEX IF EXISTS idx_sessions_expires_at;
DROP INDEX IF EXISTS idx_sessions_user_id;
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_email;

-- Drop tables
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS users;
