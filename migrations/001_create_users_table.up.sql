CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR (30) UNIQUE NOT NULL,
    email VARCHAR (60) UNIQUE NOT NULL,
    balance INTEGER NOT NULL,
    referrer_id INTEGER,
    FOREIGN KEY (referrer_id) REFERENCES users (id)
);

CREATE TABLE tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    description TEXT NOT NULL,
    reward INTEGER NOT NULL
);

CREATE TABLE user_tasks (
    user_id INTEGER NOT NULL,
    task_id INTEGER NOT NULL,
    completed BOOLEAN NOT NULL,
    PRIMARY KEY (user_id, task_id),
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (task_id) REFERENCES tasks (id)
);