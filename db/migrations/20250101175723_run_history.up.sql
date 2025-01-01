CREATE TABLE run_types (
    id TINYINT PRIMARY KEY,
    name TEXT NOT NULL,

    -- System columns
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE run_histories (
    id UUID PRIMARY KEY,
    run_type TINYINT NOT NULL,
    started_at TIMESTAMP NOT NULL,
    finished_at TIMESTAMP,
    max_seen_status_id TEXT NOT NULL,

    -- System columns
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (run_type) REFERENCES run_types (id),
    INDEX (started_at)
);

INSERT INTO run_types (id, name) VALUES
    (1, 'Recurring'),
    (2, 'One-off');
