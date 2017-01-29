CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE feeds (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

  plugin TEXT NOT NULL,
  initial_url TEXT NOT NULL,

  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  last_refreshed_at TIMESTAMPTZ,
  last_enqueued_at TIMESTAMPTZ,
  name TEXT NOT NULL UNIQUE,
  description TEXT NOT NULL,

  hex_color TEXT NOT NULL,
  icon_url TEXT NOT NULL,

  CHECK(EXTRACT(TIMEZONE FROM created_at) = '0'),
  CHECK(EXTRACT(TIMEZONE FROM updated_at) = '0'),
  CHECK(EXTRACT(TIMEZONE FROM last_refreshed_at) = '0'),
  CHECK(EXTRACT(TIMEZONE FROM last_enqueued_at) = '0')
);

CREATE TABLE posts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  feed_id UUID REFERENCES feeds NOT NULL,

  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  refreshed_at TIMESTAMPTZ,
  posted_at TIMESTAMPTZ,

  title TEXT NOT NULL,
  url TEXT NOT NULL UNIQUE,
  content TEXT NOT NULL,

  CHECK(EXTRACT(TIMEZONE FROM created_at) = '0'),
  CHECK(EXTRACT(TIMEZONE FROM updated_at) = '0'),
  CHECK(EXTRACT(TIMEZONE FROM posted_at) = '0'),
  CHECK(EXTRACT(TIMEZONE FROM refreshed_at) = '0')
);

CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

  analytics BOOLEAN NOT NULL DEFAULT false,
  email TEXT NOT NULL UNIQUE,
  encrypted_password TEXT NOT NULL,

  active BOOLEAN NOT NULL DEFAULT 'false',
  confirmed BOOLEAN NOT NULL DEFAULT 'false',
  confirmation_token TEXT NOT NULL,
  token_created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

  stripe_customer_id TEXT NOT NULL DEFAULT '',
  -- set the default to work with stripe trial
  paid_until TIMESTAMPTZ NOT NULL DEFAULT now() + interval '28 days',

  folder_ids text[] NOT NULL,

  CHECK(EXTRACT(TIMEZONE FROM token_created_at) = '0'),
  CHECK(EXTRACT(TIMEZONE FROM paid_until) = '0'),
  CHECK(EXTRACT(TIMEZONE FROM created_at) = '0'),
  CHECK(EXTRACT(TIMEZONE FROM updated_at) = '0')
);

CREATE TABLE folders (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

  name text NOT NULL,
  feed_ids text[] NOT NULL,

  CHECK(EXTRACT(TIMEZONE FROM created_at) = '0'),
  CHECK(EXTRACT(TIMEZONE FROM updated_at) = '0')
);

CREATE TABLE read_statuses (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users NOT NULL,
  post_id  UUID REFERENCES posts NOT NULL,

  read_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  device_id TEXT NOT NULL,
  location TEXT NOT NULL,

  CHECK(EXTRACT(TIMEZONE FROM read_at) = '0')
);

CREATE TABLE sessions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users NOT NULL,

  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  invalidated_at TIMESTAMPTZ,

  expires_at TIMESTAMPTZ NOT NULL DEFAULT now() + interval '28 days',
  token TEXT NOT NULL,

  CHECK(EXTRACT(TIMEZONE FROM created_at) = '0'),
  CHECK(EXTRACT(TIMEZONE FROM invalidated_at) = '0')
);
