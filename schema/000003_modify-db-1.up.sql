DO $$ BEGIN
    CREATE TYPE theme_enum AS ENUM ('light', 'dark');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

ALTER TABLE todo_lists
    ADD COLUMN created_at timestamp DEFAULT CURRENT_TIMESTAMP;

ALTER TABLE todo_lists
    ADD COLUMN IF NOT EXISTS date date;

CREATE TABLE category (
    id BIGINT PRIMARY KEY,
    name CHAR(255),
    icon_name CHAR(255)
);

CREATE TABLE rate (
    id BIGINT PRIMARY KEY,
    task_completed BIGINT,
    name CHAR(255)
);

ALTER TABLE todo_items
    ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ADD COLUMN note CHAR(255),
    ADD COLUMN notification_time TIMESTAMP,
    ADD COLUMN predicted_time_to_spend TIME,
    ADD COLUMN priority BIGINT,
    ADD COLUMN is_deleted BOOLEAN DEFAULT FALSE,
    ADD COLUMN category BIGINT;


ALTER TABLE todo_items
    ADD CONSTRAINT fk_category
    FOREIGN KEY (category) REFERENCES category(id);



ALTER TABLE users
    ADD COLUMN theme theme_enum DEFAULT 'light',
    ADD COLUMN timezone CHAR(50),
    ADD COLUMN is_paid_member BOOLEAN DEFAULT FALSE,
    ADD COLUMN rate BIGINT;




ALTER TABLE users
    ADD CONSTRAINT fk_rate
    FOREIGN KEY (rate) REFERENCES rate(id);