CREATE TABLE IF NOT EXISTS user_tasks (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    task_id INTEGER NOT NULL,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_task FOREIGN KEY (task_id) REFERENCES tasks(id),
    UNIQUE (user_id, task_id)
);
CREATE INDEX idx_user_tasks_user_id ON user_tasks (user_id);
CREATE INDEX idx_user_tasks_task_id ON user_tasks (task_id);