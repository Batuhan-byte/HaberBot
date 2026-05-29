-- Migration Down: Make topic_id column NOT NULL in articles table
-- First we ensure there are no NULL values by assigning a default or handling it (though in practice rollback is manual)
-- Since we don't have a specific default, we leave it to database admin or just try to enforce it.
ALTER TABLE articles ALTER COLUMN topic_id SET NOT NULL;
