-- Rollback the indexes if needed
DROP INDEX IF EXISTS idx_followers_following;
DROP INDEX IF EXISTS idx_follow_requests_target;
DROP INDEX IF EXISTS idx_follow_requests_requester;
