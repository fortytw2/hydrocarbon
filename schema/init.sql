CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

	stripe_customer_id TEXT,
	email TEXT NOT NULL
);

CREATE UNIQUE INDEX users_email_uniq_idx ON users (lower(email));

-- login tokens are one-time tokens used to login if oauth 2.0 is not used
CREATE TABLE login_tokens (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id UUID REFERENCES users NOT NULL,

	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	expires_at TIMESTAMPTZ NOT NULL DEFAULT now() + INTERVAL '24 HOURS',

	token TEXT DEFAULT encode(gen_random_bytes(16), 'hex'),
	used BOOLEAN NOT NULL DEFAULT false
);

CREATE UNIQUE INDEX login_tokens_token_idx ON login_tokens (token);

-- session tokens are used to authenticate sessions
CREATE TABLE sessions (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id UUID REFERENCES users NOT NULL,

	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	user_agent TEXT NOT NULL,
	ip TEXT NOT NULL,

	key TEXT DEFAULT encode(gen_random_bytes(16), 'hex'),
	active BOOLEAN NOT NULL DEFAULT true
);

CREATE UNIQUE INDEX sessions_token_idx ON sessions (token);

-- payments are where we record all payments made via stripe
CREATE TABLE payments (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id UUID REFERENCES users NOT NULL,

	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- feed_folders are folders of feeds
CREATE TABLE feed_folders (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- feeds are individual feeds
CREATE TABLE feeds (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	feed_folder_id UUID REFERENCES feed_folders NOT NULL,

	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

	plugin TEXT NOT NULL,
	url TEXT NOT NULL,
	title TEXT NOT NULL
);

CREATE UNIQUE INDEX feed_plugin_url_idx ON feeds(plugin, url);

CREATE TABLE posts (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	hash TEXT NOT NULL,

	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

	title TEXT NOT NULL,
	content TEXT NOT NULL
);
