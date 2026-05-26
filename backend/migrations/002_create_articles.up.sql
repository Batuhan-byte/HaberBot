CREATE TABLE IF NOT EXISTS articles (
    id VARCHAR(36) PRIMARY KEY,
    title TEXT NOT NULL,
    turkish_title TEXT NOT NULL DEFAULT '',
    original_url TEXT NOT NULL UNIQUE,
    source_type VARCHAR(50) NOT NULL,
    original_content TEXT NOT NULL,
    turkish_summary TEXT NOT NULL DEFAULT '',
    score INTEGER NOT NULL DEFAULT 0,
    topic_id VARCHAR(36) NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
    image_url TEXT NOT NULL DEFAULT '',
    processed_at TIMESTAMP WITH TIME ZONE,
    fetched_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_articles_topic_id ON articles(topic_id);
CREATE INDEX IF NOT EXISTS idx_articles_processed_at ON articles(processed_at) WHERE processed_at IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_articles_fetched_at_desc ON articles(fetched_at DESC);
CREATE INDEX IF NOT EXISTS idx_articles_topic_id_fetched_at_desc ON articles(topic_id, fetched_at DESC);
