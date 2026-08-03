-- ================================================
-- SEED DATA FOR DEMO / EVALUATION
-- ================================================

-- ================================================
-- STEP 1: CLEAN UP EXISTING DATA
-- (Delete in reverse dependency order)
-- ================================================

-- Group posts
DELETE FROM group_posts WHERE id = 'seed_group_post_1';
DELETE FROM group_comments WHERE post_id IN (SELECT id FROM group_posts WHERE group_id = 'seed_group_1');
DELETE FROM group_posts WHERE group_id = 'seed_group_1';

-- Group members
DELETE FROM group_members WHERE group_id = 'seed_group_1';

-- Groups
DELETE FROM groups WHERE id = 'seed_group_1';

-- Comments
DELETE FROM comments WHERE id IN ('seed_comment_1', 'seed_comment_2');

-- Posts
DELETE FROM posts WHERE id IN ('seed_post_1', 'seed_post_2', 'seed_post_3');

-- Followers
DELETE FROM followers 
WHERE follower_id IN (SELECT id FROM users WHERE username IN ('alice', 'bob', 'charlie'))
   OR following_id IN (SELECT id FROM users WHERE username IN ('alice', 'bob', 'charlie'));

-- Users
DELETE FROM users WHERE username IN ('alice', 'bob', 'charlie');

-- Reset autoincrement
DELETE FROM sqlite_sequence WHERE name = 'users';

-- ================================================
-- STEP 2: SEED USERS
-- Password is "password123" for all users
-- ================================================

INSERT INTO users (email, username, password_hash, first_name, last_name, bio, nickname, date_of_birth, is_private)
VALUES
    ('alice@example.com', 'alice', '$2a$10$8Cg9w1iY.Y8J34nC6b37.exhTqR7t4YxK858.m/u87C6x6XFmBvpe', 'Alice', 'Smith', 'Software developer, traveler, and coffee enthusiast.', 'Ali', '1995-05-15', 0),
    ('bob@example.com', 'bob', '$2a$10$8Cg9w1iY.Y8J34nC6b37.exhTqR7t4YxK858.m/u87C6x6XFmBvpe', 'Bob', 'Jones', 'Writer, sports fan, and weekend hiker.', 'Bobby', '1990-10-20', 0),
    ('charlie@example.com', 'charlie', '$2a$10$8Cg9w1iY.Y8J34nC6b37.exhTqR7t4YxK858.m/u87C6x6XFmBvpe', 'Charlie', 'Brown', 'Private profile for close friends.', 'Chuck', '1998-12-25', 1);

-- ================================================
-- STEP 3: SEED FOLLOWERS
-- ================================================

INSERT INTO followers (follower_id, following_id, created_at)
VALUES 
    ((SELECT id FROM users WHERE username = 'alice'), (SELECT id FROM users WHERE username = 'bob'), 1680000000),
    ((SELECT id FROM users WHERE username = 'bob'), (SELECT id FROM users WHERE username = 'alice'), 1680000000);

-- ================================================
-- STEP 4: SEED POSTS
-- ================================================

INSERT INTO posts (id, author_id, content, privacy, created_at)
VALUES 
    ('seed_post_1', (SELECT id FROM users WHERE username = 'alice'), 'Welcome to my new social network profile! Excited to share updates here.', 'public', 1680000100),
    ('seed_post_2', (SELECT id FROM users WHERE username = 'bob'), 'Just finished a beautiful hike in the mountains today. Highly recommend it!', 'public', 1680000200),
    ('seed_post_3', (SELECT id FROM users WHERE username = 'alice'), 'Looking for recommendations on Go design patterns. Anyone got good books/articles?', 'followers', 1680000300);

-- ================================================
-- STEP 5: SEED COMMENTS
-- ================================================

INSERT INTO comments (id, post_id, author_id, content, created_at)
VALUES 
    ('seed_comment_1', 'seed_post_1', (SELECT id FROM users WHERE username = 'bob'), 'Welcome, Alice! Glad to connect here.', 1680000400),
    ('seed_comment_2', 'seed_post_2', (SELECT id FROM users WHERE username = 'alice'), 'Awesome picture! Where was this hike?', 1680000500);

-- ================================================
-- STEP 6: SEED GROUPS
-- ================================================

INSERT INTO groups (id, title, description, creator_id, created_at)
VALUES 
    ('seed_group_1', 'Go Developers', 'A space to share knowledge, libraries, and best practices about the Go programming language.', (SELECT id FROM users WHERE username = 'alice'), 1680000600);

-- ================================================
-- STEP 7: SEED GROUP MEMBERS
-- ================================================

INSERT INTO group_members (group_id, user_id, joined_at)
VALUES 
    ('seed_group_1', (SELECT id FROM users WHERE username = 'alice'), 1680000600),
    ('seed_group_1', (SELECT id FROM users WHERE username = 'bob'), 1680000700);

-- ================================================
-- STEP 8: SEED GROUP POSTS
-- ================================================

INSERT INTO group_posts (id, group_id, author_id, content, created_at)
VALUES 
    ('seed_group_post_1', 'seed_group_1', (SELECT id FROM users WHERE username = 'bob'), 'Has anyone tried the new standard library features in Go 1.22?', 1680000800);
