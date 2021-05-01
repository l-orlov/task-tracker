CREATE OR REPLACE FUNCTION trigger_set_timestamp()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- users of system
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

-- all importance_statuses for tasks
CREATE TABLE importance_statuses
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
CREATE TRIGGER set_timestamp
    BEFORE UPDATE
    ON importance_statuses
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

-- all progress_statuses for tasks
CREATE TABLE progress_statuses
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
CREATE TRIGGER set_timestamp
    BEFORE UPDATE
    ON progress_statuses
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

-- projects to work on
CREATE TABLE projects
(
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(255)                 NOT NULL,
    description TEXT                         NOT NULL,
    owner       BIGINT REFERENCES users (id) NOT NULL,
    created_at  TIMESTAMPTZ                  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ                  NOT NULL DEFAULT NOW(),
    closed_at   TIMESTAMPTZ                  NOT NULL DEFAULT NOW()
);
CREATE TRIGGER set_timestamp
    BEFORE UPDATE
    ON projects
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

-- users working on projects
CREATE TABLE projects_users
(
    project_id BIGINT REFERENCES projects (id) ON DELETE CASCADE NOT NULL,
    user_id    BIGINT REFERENCES users (id) ON DELETE CASCADE    NOT NULL
);
CREATE INDEX idx_projects_users ON projects_users (project_id, user_id);

-- importance_statuses for project tasks
CREATE TABLE projects_importance_statuses
(
    id                   SERIAL PRIMARY KEY,
    project_id           BIGINT REFERENCES projects (id) ON DELETE CASCADE         NOT NULL,
    importance_status_id INT REFERENCES importance_statuses (id) ON DELETE CASCADE NOT NULL
);
CREATE INDEX idx_projects_importance_statuses ON projects_importance_statuses (project_id, importance_status_id);

-- progress_statuses for project tasks
CREATE TABLE projects_progress_statuses
(
    id                 SERIAL PRIMARY KEY,
    project_id         BIGINT REFERENCES projects (id) ON DELETE CASCADE       NOT NULL,
    progress_status_id INT REFERENCES progress_statuses (id) ON DELETE CASCADE NOT NULL
);
CREATE INDEX idx_projects_progress_statuses ON projects_progress_statuses (project_id, progress_status_id);

-- tasks to project
CREATE TABLE tasks
(
    id                   BIGSERIAL PRIMARY KEY,
    project_id           BIGINT REFERENCES projects (id) ON DELETE CASCADE NOT NULL,
    title                VARCHAR(255)                                      NOT NULL,
    description          TEXT                                              NOT NULL,
    assignee_id          BIGINT REFERENCES users (id)                      NOT NULL,
    importance_status_id INT REFERENCES projects_importance_statuses (id)  NOT NULL,
    progress_status_id   INT REFERENCES projects_progress_statuses (id)    NOT NULL,
    created_at           TIMESTAMPTZ                                       NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ                                       NOT NULL DEFAULT NOW()
);
CREATE TRIGGER set_timestamp
    BEFORE UPDATE
    ON tasks
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();
CREATE INDEX idx_tasks_project_id ON tasks (project_id);

-- sprints have tasks that should be done for sprint time
CREATE TABLE sprints
(
    id         BIGSERIAL PRIMARY KEY,
    project_id BIGINT REFERENCES projects (id) ON DELETE CASCADE NOT NULL,
    created_at TIMESTAMPTZ                                       NOT NULL DEFAULT NOW(),
    closed_at  TIMESTAMPTZ                                       NOT NULL DEFAULT NOW()
);

CREATE TABLE projects_sprints
(
    project_id BIGINT REFERENCES projects (id) ON DELETE CASCADE NOT NULL,
    sprint_id  BIGINT REFERENCES sprints (id) ON DELETE CASCADE  NOT NULL
);
CREATE INDEX idx_projects_sprints ON projects_sprints (project_id, sprint_id);

CREATE TABLE sprints_tasks
(
    sprint_id BIGINT REFERENCES sprints (id) ON DELETE CASCADE NOT NULL,
    task_id   BIGINT REFERENCES tasks (id) ON DELETE CASCADE   NOT NULL
);
CREATE INDEX idx_sprints_tasks ON sprints_tasks (sprint_id, task_id);

-- insert default values
INSERT INTO importance_statuses (name)
VALUES ('LOW'),
       ('MEDIUM'),
       ('HIGH');

INSERT INTO progress_statuses (name)
VALUES ('TO DO'),
       ('IN HOLD'),
       ('IN PROGRESS'),
       ('IN REVIEW'),
       ('IN TESTING'),
       ('DONE / CLOSED');
