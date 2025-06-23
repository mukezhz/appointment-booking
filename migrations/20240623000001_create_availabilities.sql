-- Create availabilities table
CREATE TABLE availabilities (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    weekday VARCHAR(10) NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT fk_user_availability FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Add index for availability lookup
CREATE INDEX idx_availabilities_user_weekday ON availabilities(user_id, weekday) WHERE deleted_at IS NULL;
