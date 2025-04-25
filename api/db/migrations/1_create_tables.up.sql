--CREATE TABLE note (
--	id TEXT PRIMARY KEY,
--	text TEXT,
--	cover_url TEXT
--);

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

--------------------------
-- Users Table
--------------------------
CREATE TABLE
    users (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
        username VARCHAR(50) NOT NULL UNIQUE,
        password_hash VARCHAR(128) NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

--------------------------
-- Posts Table
--------------------------
CREATE TABLE
    posts (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
        user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
        title VARCHAR(255) NOT NULL,
        content TEXT NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

--------------------------
-- Comments Table
--------------------------
CREATE TABLE
    comments (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
        post_id UUID NOT NULL REFERENCES posts (id) ON DELETE CASCADE,
        user_id UUID REFERENCES users (id) ON DELETE SET NULL,
        content TEXT NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

--------------------------
-- Other things
--------------------------
CREATE
OR REPLACE FUNCTION update_updated_at_column () RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_users_updated_at BEFORE
UPDATE ON users FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column ();

CREATE TRIGGER trg_posts_updated_at BEFORE
UPDATE ON posts FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column ();

CREATE TRIGGER trg_comments_updated_at BEFORE
UPDATE ON comments FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column ();

CREATE INDEX idx_posts_user_id ON posts (user_id);

CREATE INDEX idx_comments_post_id ON comments (post_id);

CREATE INDEX idx_comments_user_id ON comments (user_id);