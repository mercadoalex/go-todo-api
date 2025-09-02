-- Create the tasks table
CREATE TABLE IF NOT EXISTS tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    completed BOOLEAN NOT NULL DEFAULT 0
);

-- Insert a new task
INSERT INTO tasks (title, completed) VALUES ('Learn Go', 0);

-- Select all tasks
SELECT * FROM tasks;

-- Select a single task by ID
SELECT * FROM tasks WHERE id = ?;

-- Update a task (change title and/or completed status)
UPDATE tasks SET title = ?, completed = ? WHERE id = ?;

-- Mark a task as completed
UPDATE tasks SET completed = 1 WHERE id = ?;

-- Delete a task by ID
DELETE FROM tasks WHERE id = ?;