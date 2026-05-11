-- ============================================================================
-- Migration 006: Sync notify_email from users.email
--
-- 确保所有用户的 settings.notify_email 都有值。
-- 对于 notify_email 为空的记录，从 users.email 同步过来。
-- 这样后续通知和密码找回都能正常工作。
-- ============================================================================

UPDATE settings
SET notify_email = (
    SELECT u.email FROM users u WHERE u.id = settings.id
)
WHERE notify_email IS NULL OR notify_email = '';
