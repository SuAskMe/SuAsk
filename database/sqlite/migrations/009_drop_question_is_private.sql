-- 当前运行环境的 SQLite 版本已支持 DROP COLUMN，直接删除该列即可，
-- 避免重建父表时对子表外键引用造成额外扰动。
ALTER TABLE questions DROP COLUMN is_private;
