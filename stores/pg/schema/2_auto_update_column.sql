CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE PROCEDURE set_updated_at();

CREATE TRIGGER update_feeds_updated_at BEFORE UPDATE ON feeds FOR EACH ROW EXECUTE PROCEDURE set_updated_at();

CREATE TRIGGER update_read_statuses_updated_at BEFORE UPDATE ON read_statuses FOR EACH ROW EXECUTE PROCEDURE set_updated_at();

CREATE TRIGGER update_posts_updated_at BEFORE UPDATE ON posts FOR EACH ROW EXECUTE PROCEDURE set_updated_at();

CREATE TRIGGER update_sessions_updated_at BEFORE UPDATE ON sessions FOR EACH ROW EXECUTE PROCEDURE set_updated_at();
