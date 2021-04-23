CREATE OR REPLACE FUNCTION trigger_set_timestamp()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE users
(
    id                 BIGSERIAL PRIMARY KEY,
    email              VARCHAR(320) NOT NULL UNIQUE,
    first_name         VARCHAR(255) NOT NULL,
    last_name          VARCHAR(255) NOT NULL,
    password           VARCHAR(255) NOT NULL,
    is_email_confirmed BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at         TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
CREATE TRIGGER set_timestamp
    BEFORE UPDATE
    ON users
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TABLE importance_status
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
CREATE TRIGGER set_timestamp
    BEFORE UPDATE
    ON importance_status
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TABLE progress_status
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
CREATE TRIGGER set_timestamp
    BEFORE UPDATE
    ON progress_status
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TABLE projects
(
    id                   BIGSERIAL PRIMARY KEY,
    title                VARCHAR(255)                 NOT NULL,
    description          TEXT                         NOT NULL,
    assignee_id          BIGINT REFERENCES users (id) NOT NULL,
    importance_status_id INT REFERENCES users (id)    NOT NULL,
    progress_status_id   INT REFERENCES users (id)    NOT NULL,
    created_at           TIMESTAMPTZ                  NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ                  NOT NULL DEFAULT NOW()
);
CREATE TRIGGER set_timestamp
    BEFORE UPDATE
    ON projects
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TABLE tasks
(
    id                   BIGSERIAL PRIMARY KEY,
    title                VARCHAR(255)                 NOT NULL,
    description          TEXT                         NOT NULL,
    assignee_id          BIGINT REFERENCES users (id) NOT NULL,
    importance_status_id INT REFERENCES users (id)    NOT NULL,
    progress_status_id   INT REFERENCES users (id)    NOT NULL,
    created_at           TIMESTAMPTZ                  NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ                  NOT NULL DEFAULT NOW()
);
CREATE TRIGGER set_timestamp
    BEFORE UPDATE
    ON tasks
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TABLE projects_tasks
(
    project_id BIGINT REFERENCES projects (id) ON DELETE CASCADE NOT NULL,
    task_id    BIGINT REFERENCES tasks (id) ON DELETE CASCADE    NOT NULL
);
