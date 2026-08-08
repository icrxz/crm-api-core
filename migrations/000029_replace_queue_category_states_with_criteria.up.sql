ALTER TABLE queues ADD COLUMN IF NOT EXISTS criteria JSONB NOT NULL DEFAULT '{}';

UPDATE queues
SET criteria = jsonb_build_object('category', category) ||
    CASE WHEN array_length(states, 1) > 0 THEN jsonb_build_object('state', to_jsonb(states)) ELSE '{}'::jsonb END
WHERE criteria = '{}';

DROP INDEX IF EXISTS idx_queues_states;
DROP INDEX IF EXISTS idx_queues_category;

ALTER TABLE queues DROP COLUMN IF EXISTS category;
ALTER TABLE queues DROP COLUMN IF EXISTS states;

CREATE INDEX IF NOT EXISTS idx_queues_criteria ON queues USING GIN (criteria jsonb_path_ops);
