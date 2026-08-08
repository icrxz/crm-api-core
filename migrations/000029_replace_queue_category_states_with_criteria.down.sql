ALTER TABLE queues ADD COLUMN IF NOT EXISTS category TEXT NOT NULL DEFAULT 'mobile';
ALTER TABLE queues ADD COLUMN IF NOT EXISTS states TEXT[] NOT NULL DEFAULT '{}';

UPDATE queues
SET category = COALESCE(criteria->>'category', 'mobile'),
    states = CASE WHEN jsonb_typeof(criteria->'state') = 'array'
        THEN ARRAY(SELECT jsonb_array_elements_text(criteria->'state'))
        ELSE '{}'::text[]
    END;

DROP INDEX IF EXISTS idx_queues_criteria;
ALTER TABLE queues DROP COLUMN IF EXISTS criteria;

CREATE INDEX IF NOT EXISTS idx_queues_states ON queues USING GIN (states);
CREATE INDEX IF NOT EXISTS idx_queues_category ON queues (category);
