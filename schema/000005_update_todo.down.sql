-- Down Migration: Remove action_made column from todo_items table
ALTER TABLE todo_items
    DROP COLUMN action_made;
