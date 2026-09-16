-- comment_history is an audit trail and must survive the deletion of the
-- comment it describes (e.g. ResetCaseStatus hard-deletes comments). A FK
-- here blocks exactly the case it exists to record.
ALTER TABLE comment_history
DROP CONSTRAINT IF EXISTS comment_history_comment_id_fkey;

CREATE INDEX IF NOT EXISTS idx_comment_history_comment_id ON comment_history (comment_id);
