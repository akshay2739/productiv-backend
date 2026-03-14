-- Users table
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL DEFAULT 'Default User',
    timezone VARCHAR(100) NOT NULL DEFAULT 'Asia/Kolkata',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Insert default user
INSERT INTO users (name, timezone) VALUES ('Akshay', 'Asia/Kolkata');

-- Pillars table
CREATE TABLE IF NOT EXISTS pillars (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    icon VARCHAR(10) NOT NULL,
    color VARCHAR(20) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    display_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, type)
);

-- Insert default pillars for the default user
INSERT INTO pillars (user_id, type, name, icon, color, display_order) VALUES
    (1, 'fasting',    'Fasting',    '🍽️', '#e94560', 1),
    (1, 'gym',        'Gym',        '💪', '#4a9eff', 2),
    (1, 'meditation', 'Meditation', '🧘', '#9b59b6', 3),
    (1, 'retention',  'Retention',  '🔥', '#2ecc71', 4);

-- Fasting sessions table
CREATE TABLE IF NOT EXISTS fasting_sessions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    protocol VARCHAR(20) NOT NULL,
    target_hours INT NOT NULL,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ,
    actual_duration_hours DOUBLE PRECISION,
    target_reached BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_fasting_sessions_user_active ON fasting_sessions(user_id) WHERE end_time IS NULL;
CREATE INDEX idx_fasting_sessions_user_date ON fasting_sessions(user_id, start_time);

-- Gym sessions table
CREATE TABLE IF NOT EXISTS gym_sessions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    workout_type VARCHAR(50) NOT NULL,
    duration_min INT,
    energy_level INT CHECK (energy_level IS NULL OR (energy_level >= 1 AND energy_level <= 5)),
    logged_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_gym_sessions_user_date ON gym_sessions(user_id, logged_at);

-- Meditation sessions table
CREATE TABLE IF NOT EXISTS meditation_sessions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    target_minutes INT NOT NULL,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ,
    actual_duration_minutes DOUBLE PRECISION,
    mood_before INT CHECK (mood_before IS NULL OR (mood_before >= 1 AND mood_before <= 5)),
    mood_after INT CHECK (mood_after IS NULL OR (mood_after >= 1 AND mood_after <= 5)),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_meditation_sessions_user_active ON meditation_sessions(user_id) WHERE end_time IS NULL;
CREATE INDEX idx_meditation_sessions_user_date ON meditation_sessions(user_id, start_time);

-- Retention streaks table
CREATE TABLE IF NOT EXISTS retention_streaks (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ,
    days_count INT NOT NULL DEFAULT 0,
    reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_retention_streaks_user_active ON retention_streaks(user_id) WHERE end_date IS NULL;
