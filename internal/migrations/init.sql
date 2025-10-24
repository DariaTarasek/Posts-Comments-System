CREATE TABLE IF NOT EXISTS posts (
                                     id SERIAL PRIMARY KEY,
                                     title TEXT NOT NULL,
                                     content TEXT NOT NULL,
                                     author TEXT NOT NULL,
                                     are_comments_allowed BOOLEAN DEFAULT TRUE,
                                     created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS comments (
                                        id SERIAL PRIMARY KEY,
                                        post_id SERIAL REFERENCES posts(id) ON DELETE CASCADE,
                                        author TEXT NOT NULL,
                                        content TEXT NOT NULL CHECK (length(content) <= 2000),
                                        parent_comment_id INT REFERENCES comments(id) ON DELETE CASCADE,
                                        path TEXT NOT NULL,
                                        created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_post_comment_path ON comments(post_id, path);
