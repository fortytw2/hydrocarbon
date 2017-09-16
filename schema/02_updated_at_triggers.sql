CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    IF row(NEW.*) IS DISTINCT FROM row(OLD.*) THEN
        NEW.updated_at = now(); 
        RETURN NEW;
    ELSE
        RETURN OLD;
    END IF;
END;
$$ language 'plpgsql';

CREATE TRIGGER users_updated_at
    BEFORE UPDATE ON users 
    FOR EACH ROW EXECUTE PROCEDURE set_updated_at();

CREATE TRIGGER login_tokens_updated_at
    BEFORE UPDATE ON login_tokens 
    FOR EACH ROW EXECUTE PROCEDURE set_updated_at();

CREATE TRIGGER sessions_updated_at
    BEFORE UPDATE ON sessions 
    FOR EACH ROW EXECUTE PROCEDURE set_updated_at();

CREATE TRIGGER folders_updated_at
    BEFORE UPDATE ON folders 
    FOR EACH ROW EXECUTE PROCEDURE set_updated_at();

CREATE TRIGGER feeds_updated_at
    BEFORE UPDATE ON feeds 
    FOR EACH ROW EXECUTE PROCEDURE set_updated_at();

CREATE TRIGGER feed_folders_updated_at
    BEFORE UPDATE ON feed_folders 
    FOR EACH ROW EXECUTE PROCEDURE set_updated_at();

CREATE TRIGGER posts_updated_at
    BEFORE UPDATE ON posts 
    FOR EACH ROW EXECUTE PROCEDURE set_updated_at();