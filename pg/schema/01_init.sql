CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE users (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v1mc(),

	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

	admin BOOLEAN NOT NULL DEFAULT false,

	stripe_customer_id CITEXT,
	stripe_subscription_id CITEXT,
	last_payment_date TIMESTAMPTZ,
	email CITEXT NOT NULL,

	UNIQUE (email)
);

-- login tokens are one-time tokens used to login if oauth 2.0 is not used
CREATE TABLE login_tokens (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v1mc(),
	user_id UUID REFERENCES users NOT NULL,

	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	expires_at TIMESTAMPTZ NOT NULL DEFAULT now() + INTERVAL '24 HOURS',

	user_agent TEXT NOT NULL,
	ip CIDR NOT NULL,

	token TEXT DEFAULT encode(gen_random_bytes(16), 'hex'),
	used BOOLEAN NOT NULL DEFAULT false
);

CREATE UNIQUE INDEX login_tokens_token_idx ON login_tokens (token);

-- session tokens are used to authenticate sessions
CREATE TABLE sessions (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v1mc(),
	user_id UUID REFERENCES users NOT NULL,

	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	user_agent TEXT NOT NULL,
	ip CIDR NOT NULL,

	key CITEXT DEFAULT encode(gen_random_bytes(16), 'hex'),
	active BOOLEAN NOT NULL DEFAULT true
);

CREATE UNIQUE INDEX sessions_key_idx ON sessions (key);

-- folders are used to maintain collections of feeds
CREATE TABLE folders (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v1mc(),
	user_id UUID REFERENCES USERS NOT NULL,
	
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

	name TEXT NOT NULL DEFAULT 'default',

	UNIQUE (user_id, name)
);

-- feeds are individual feeds
CREATE TABLE feeds (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v1mc(),

	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

	plugin TEXT NOT NULL,
	url TEXT NOT NULL,
	title TEXT NOT NULL,

	public BOOLEAN NOT NULL DEFAULT true
);

CREATE UNIQUE INDEX feeds_plugin_url_public_uniq_idx ON feeds (plugin, url) WHERE public;

-- feed_folders is a join table between feeds and folders
CREATE TABLE feed_folders (
	user_id UUID REFERENCES users NOT NULL,
	folder_id UUID REFERENCES folders NOT NULL,
	feed_id UUID REFERENCES feeds NOT NULL,

	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

	priority INT NOT NULL DEFAULT 0,

	-- display this feed in collapsed or full mode
	display_mode TEXT NOT NULL DEFAULT 'full',

	PRIMARY KEY (user_id, folder_id, feed_id)
);

CREATE TABLE posts (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v1mc(),
	feed_id UUID NOT NULL REFERENCES feeds (id),

	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	posted_at TIMESTAMPTZ NOT NULL DEFAULT '01-01-1970 00:00:00'::timestamptz,

	content_hash CITEXT NOT NULL,
	title TEXT NOT NULL,
	author TEXT NOT NULL DEFAULT '',
	body TEXT NOT NULL,
	url TEXT NOT NULL,

	extra JSONB,

	UNIQUE (content_hash)
);

CREATE INDEX posts_feed_idx ON posts (feed_id);

-- read statuses tracking
CREATE TABLE read_statuses (
	post_id UUID NOT NULL REFERENCES posts,
	user_id UUID NOT NULL REFERENCES users,

	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

	PRIMARY KEY (post_id, user_id)
);

CREATE TYPE scrape_state AS ENUM (
	'STOPPED',
	'RUNNING',
	'ERRORED',
	'SUCCESS'
);

CREATE TABLE scrapes (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v1mc(),
	feed_id UUID NOT NULL REFERENCES feeds (id),

	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

	scheduled_start_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	started_at TIMESTAMPTZ,
	ended_at TIMESTAMPTZ,

	state scrape_state NOT NULL DEFAULT 'STOPPED',
	errors TEXT[] NOT NULL DEFAULT '{}',

	config JSONB DEFAULT '{}',

	total_datums INT DEFAULT 0,
	total_retries INT DEFAULT 0,
	total_tasks INT DEFAULT 0
);