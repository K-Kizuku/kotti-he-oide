-- Web Push Notification Database Schema
-- Based on TODO.md specifications

-- Users table (existing - for reference)
-- CREATE TABLE users (
--   id BIGSERIAL PRIMARY KEY,
--   email TEXT UNIQUE,
--   created_at TIMESTAMPTZ NOT NULL DEFAULT now()
-- );

-- Push subscriptions: Browser subscription data
CREATE TABLE push_subscriptions (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
  endpoint TEXT NOT NULL UNIQUE,            -- Push service endpoint (keep confidential)
  p256dh   TEXT NOT NULL,                  -- Browser public key for encryption
  auth     TEXT NOT NULL,                  -- Browser auth secret for encryption
  ua       TEXT,                           -- User-Agent string
  expiration_time TIMESTAMPTZ,             -- MDN expiration info if available
  is_valid BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Notification templates (optional)
CREATE TABLE notification_templates (
  id BIGSERIAL PRIMARY KEY,
  key TEXT UNIQUE NOT NULL,
  title TEXT NOT NULL,
  body  TEXT NOT NULL,
  url   TEXT,
  icon  TEXT,
  data  JSONB,                             -- Additional payload data
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- User notification preferences
CREATE TABLE notification_prefs (
  user_id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  enabled BOOLEAN NOT NULL DEFAULT TRUE,
  topics JSONB NOT NULL DEFAULT '{}'::jsonb, -- e.g., {"news": true, "promo": false}
  quiet_hours JSONB                           -- Silent periods configuration
);

-- Push job statuses
CREATE TYPE job_status AS ENUM ('pending','sending','succeeded','failed','cancelled');

-- Push jobs (async sending queue)
CREATE TABLE push_jobs (
  id BIGSERIAL PRIMARY KEY,
  idempotency_key TEXT UNIQUE,             -- Prevent duplicate sends
  template_key TEXT REFERENCES notification_templates(key),
  user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
  topic TEXT,                              -- Web Push Topic header (for collapsing)
  urgency TEXT CHECK (urgency IN ('very-low','low','normal','high')), -- RFC 8030 urgency levels
  ttl_seconds INT CHECK (ttl_seconds >= 0),
  payload JSONB NOT NULL,                  -- JSON data parsed by Service Worker
  schedule_at TIMESTAMPTZ,                 -- Scheduled delivery time
  status job_status NOT NULL DEFAULT 'pending',
  retry_count INT NOT NULL DEFAULT 0,
  last_error TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Push delivery logs
CREATE TABLE push_logs (
  id BIGSERIAL PRIMARY KEY,
  job_id BIGINT REFERENCES push_jobs(id) ON DELETE SET NULL,
  subscription_id BIGINT REFERENCES push_subscriptions(id) ON DELETE SET NULL,
  response_status INT,                     -- HTTP response status from push service
  response_headers JSONB,                  -- Response headers from push service
  error TEXT,                              -- Error message if delivery failed
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Indexes for performance
CREATE INDEX idx_push_subscriptions_user_id ON push_subscriptions(user_id);
CREATE INDEX idx_push_subscriptions_endpoint ON push_subscriptions(endpoint);
CREATE INDEX idx_push_subscriptions_valid ON push_subscriptions(is_valid) WHERE is_valid = true;
CREATE INDEX idx_push_subscriptions_expiration ON push_subscriptions(expiration_time) WHERE expiration_time IS NOT NULL;

CREATE INDEX idx_push_jobs_status ON push_jobs(status);
CREATE INDEX idx_push_jobs_schedule_at ON push_jobs(schedule_at) WHERE schedule_at IS NOT NULL;
CREATE INDEX idx_push_jobs_user_id ON push_jobs(user_id);
CREATE INDEX idx_push_jobs_idempotency ON push_jobs(idempotency_key) WHERE idempotency_key IS NOT NULL;

CREATE INDEX idx_push_logs_job_id ON push_logs(job_id);
CREATE INDEX idx_push_logs_subscription_id ON push_logs(subscription_id);
CREATE INDEX idx_push_logs_created_at ON push_logs(created_at);

-- Sample data for testing
INSERT INTO notification_templates (key, title, body, url, icon) VALUES 
('welcome', 'Welcome!', 'Welcome to our application!', '/', '/icons/welcome.png'),
('news', 'News Update', 'Check out the latest news', '/news', '/icons/news.png'),
('reminder', 'Reminder', 'You have a pending task', '/tasks', '/icons/reminder.png');