DROP INDEX IF EXISTS idx_comment_history_comment_id;

ALTER TABLE comment_history ADD CONSTRAINT comment_history_comment_id_fkey FOREIGN KEY (comment_id) REFERENCES comments (comment_id);
