CREATE UNIQUE INDEX posts_id_idx ON posts (id);

CREATE INDEX posts_feed_id_posted_at_idx ON posts (feed_id, posted_at);

CREATE INDEX sessions_user_id ON sessions (user_id);

CREATE UNIQUE INDEX sessions_token_idx ON sessions (token, invalidated_at, expires_at);

CREATE UNIQUE INDEX user_email_idx ON users (lower(email));

CREATE INDEX user_folder_ids_idx ON users USING GIN (folder_ids);

CREATE INDEX folder_feed_ids_idx ON folders USING GIN (feed_ids);
