-- Up Migration: Add parent_id column to todo_items table
ALTER TABLE todo_items
    ADD COLUMN parent_id BIGINT;
