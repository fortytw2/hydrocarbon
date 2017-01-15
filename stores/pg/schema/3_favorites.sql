CREATE TABLE user_favorites (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

  post_id UUID REFERENCES posts NOT NULL,
  user_id UUID REFERENCES users NOT NULL,

  CHECK(EXTRACT(TIMEZONE FROM created_at) = '0'),
  CHECK(EXTRACT(TIMEZONE FROM updated_at) = '0')
);

CREATE TRIGGER update_user_favorites_updated_at BEFORE UPDATE ON user_favorites FOR EACH ROW EXECUTE PROCEDURE set_updated_at();
