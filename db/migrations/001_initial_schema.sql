-- QFlow Initial Schema
-- Migration: 001_initial_schema.sql

CREATE TABLE IF NOT EXISTS users (
    id         SERIAL PRIMARY KEY,
    phone      VARCHAR(20) UNIQUE NOT NULL,
    name       VARCHAR(100),
    role       VARCHAR(20) NOT NULL DEFAULT 'user',  -- user, provider, admin
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS otps (
    id         SERIAL PRIMARY KEY,
    phone      VARCHAR(20) NOT NULL,
    code       VARCHAR(6) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    used       BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS categories (
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS providers (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS zones (
    id          SERIAL PRIMARY KEY,
    provider_id INTEGER NOT NULL REFERENCES providers(id) ON DELETE CASCADE,
    name        VARCHAR(100) NOT NULL,
    is_open     BOOLEAN NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS queues (
    id           SERIAL PRIMARY KEY,
    queue_number INTEGER NOT NULL,
    zone_id      INTEGER NOT NULL REFERENCES zones(id) ON DELETE CASCADE,
    user_id      INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status       VARCHAR(20) NOT NULL DEFAULT 'waiting',  -- waiting, called, completed, skipped, cancelled
    created_at   TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (zone_id, queue_number)
);

CREATE TABLE IF NOT EXISTS notifications (
    id         SERIAL PRIMARY KEY,
    user_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    message    TEXT NOT NULL,
    is_read    BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indexes for OTP lookups
CREATE INDEX IF NOT EXISTS idx_otps_phone      ON otps (phone);
CREATE INDEX IF NOT EXISTS idx_otps_expires_at ON otps (expires_at);

-- Indexes for FK columns (providers, zones)
CREATE INDEX IF NOT EXISTS idx_providers_category_id ON providers (category_id);
CREATE INDEX IF NOT EXISTS idx_zones_provider_id     ON zones (provider_id);

-- Indexes for FK columns and status filter (queues)
CREATE INDEX IF NOT EXISTS idx_queues_zone_id ON queues (zone_id);
CREATE INDEX IF NOT EXISTS idx_queues_user_id ON queues (user_id);
CREATE INDEX IF NOT EXISTS idx_queues_status  ON queues (status);

-- Index for notification lookup by user
CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications (user_id);
