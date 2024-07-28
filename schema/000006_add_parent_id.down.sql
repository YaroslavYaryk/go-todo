-- Down Migration: Remove parent_id column from todo_items table
ALTER TABLE todo_items
    DROP COLUMN parent_id;