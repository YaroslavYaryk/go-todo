-- Remove foreign key constraints
ALTER TABLE users
    DROP CONSTRAINT IF EXISTS fk_rate;

ALTER TABLE todo_items
    DROP CONSTRAINT IF EXISTS fk_category;

-- Remove added columns from todo_items table
ALTER TABLE todo_items
    DROP COLUMN IF EXISTS created_at,
    DROP COLUMN IF EXISTS updated_at,
    DROP COLUMN IF EXISTS note,
    DROP COLUMN IF EXISTS notification_time,
    DROP COLUMN IF EXISTS predicted_time_to_spend,
    DROP COLUMN IF EXISTS priority,
    DROP COLUMN IF EXISTS is_deleted,
    DROP COLUMN IF EXISTS category;

-- Drop the category table
DROP TABLE IF EXISTS category;

-- Remove added columns from todolist table
ALTER TABLE todolist
    DROP COLUMN IF EXISTS created_at,
    DROP COLUMN IF EXISTS date;

-- Remove added columns from users table
ALTER TABLE users
    DROP COLUMN IF EXISTS theme,
    DROP COLUMN IF EXISTS timezone,
    DROP COLUMN IF EXISTS is_paid_member,
    DROP COLUMN IF EXISTS rate;

-- Drop the rate table
DROP TABLE IF EXISTS rate;