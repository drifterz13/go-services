CREATE TABLE IF NOT EXISTS tasks (
    id UUID DEFAULT GEN_RANDOM_UUID() PRIMARY KEY,
    title text NOT NULL,
    status INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS users (
    id UUID DEFAULT GEN_RANDOM_UUID() PRIMARY KEY,
    email VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS task_members (
    user_id UUID NOT NULL,
    task_id UUID NOT NULL,
    role INTEGER NOT NULL CHECK (
        role = 0
        OR role = 1
    ),
    UNIQUE (user_id, task_id, role),
    PRIMARY KEY (user_id, task_id),
    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id),
    CONSTRAINT fk_task FOREIGN KEY(task_id) REFERENCES tasks(id)
);