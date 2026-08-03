DELETE FROM group_posts WHERE id = 'seed_group_post_1';
DELETE FROM group_members WHERE group_id = 'seed_group_1';
DELETE FROM groups WHERE id = 'seed_group_1';
DELETE FROM comments WHERE id IN ('seed_comment_1', 'seed_comment_2');
DELETE FROM posts WHERE id IN ('seed_post_1', 'seed_post_2', 'seed_post_3');
DELETE FROM followers WHERE (follower_id = '1' AND following_id = '2') OR (follower_id = '2' AND following_id = '1');
DELETE FROM users WHERE id IN (1, 2, 3);
