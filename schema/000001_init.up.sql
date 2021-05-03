-- naming relied on this article: http://citforum.ru/database/articles/naming_rule/

CREATE OR REPLACE FUNCTION trigger_set_timestamp()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- users
CREATE TABLE r_user
(
    id                 BIGSERIAL PRIMARY KEY,
    email              VARCHAR(320) NOT NULL UNIQUE,
    firstname          VARCHAR(255) NOT NULL,
    lastname           VARCHAR(255) NOT NULL,
    password           VARCHAR(255) NOT NULL,
    is_email_confirmed BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at         TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
CREATE TRIGGER set_timestamp
    BEFORE UPDATE
    ON r_user
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

-- all importance_statuses for tasks
CREATE TABLE s_importance_status
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
CREATE TRIGGER set_timestamp
    BEFORE UPDATE
    ON s_importance_status
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

-- all progress_statuses for tasks
CREATE TABLE s_progress_status
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
CREATE TRIGGER set_timestamp
    BEFORE UPDATE
    ON s_progress_status
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

-- projects to work on
CREATE TABLE r_project
(
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    description TEXT         NOT NULL,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    closed_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
CREATE TRIGGER set_timestamp
    BEFORE UPDATE
    ON r_project
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

-- users working on projects
CREATE TABLE nn_project_user
(
    project_id BIGINT REFERENCES r_project (id) ON DELETE CASCADE NOT NULL,
    user_id    BIGINT REFERENCES r_user (id) ON DELETE CASCADE    NOT NULL,
    is_owner   BOOLEAN                                            NOT NULL DEFAULT FALSE,
    UNIQUE (project_id, user_id)
);

-- importance_statuses for project tasks
CREATE TABLE s_project_importance_status
(
    id                   SERIAL PRIMARY KEY,
    project_id           BIGINT REFERENCES r_project (id) ON DELETE CASCADE        NOT NULL,
    importance_status_id INT REFERENCES s_importance_status (id) ON DELETE CASCADE NOT NULL,
    UNIQUE (project_id, importance_status_id)
);

-- progress_statuses for project tasks
CREATE TABLE s_project_progress_status
(
    id                 SERIAL PRIMARY KEY,
    project_id         BIGINT REFERENCES r_project (id) ON DELETE CASCADE      NOT NULL,
    progress_status_id INT REFERENCES s_progress_status (id) ON DELETE CASCADE NOT NULL,
    UNIQUE (project_id, progress_status_id)
);

-- tasks to project
CREATE TABLE r_task
(
    id                   BIGSERIAL PRIMARY KEY,
    project_id           BIGINT REFERENCES r_project (id) ON DELETE CASCADE NOT NULL,
    title                VARCHAR(255)                                       NOT NULL,
    description          TEXT                                               NOT NULL,
    assignee_id          BIGINT REFERENCES r_user (id)                      NOT NULL,
    importance_status_id INT REFERENCES s_project_importance_status (id)    NOT NULL,
    progress_status_id   INT REFERENCES s_project_progress_status (id)      NOT NULL,
    created_at           TIMESTAMPTZ                                        NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ                                        NOT NULL DEFAULT NOW()
);
CREATE TRIGGER set_timestamp
    BEFORE UPDATE
    ON r_task
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();
CREATE INDEX idx_r_task_project_id ON r_task (project_id);

-- sprints have tasks that should be done for sprint time
CREATE TABLE r_sprint
(
    id         BIGSERIAL PRIMARY KEY,
    project_id BIGINT REFERENCES r_project (id) ON DELETE CASCADE NOT NULL,
    created_at TIMESTAMPTZ                                        NOT NULL DEFAULT NOW(),
    closed_at  TIMESTAMPTZ
);

CREATE TABLE nn_sprint_task
(
    sprint_id BIGINT REFERENCES r_sprint (id) ON DELETE CASCADE NOT NULL,
    task_id   BIGINT REFERENCES r_task (id) ON DELETE CASCADE   NOT NULL,
    PRIMARY KEY (sprint_id, task_id)
);

-- insert default values
INSERT INTO s_importance_status (name)
VALUES ('LOW'),
       ('MEDIUM'),
       ('HIGH');

INSERT INTO s_progress_status (name)
VALUES ('TO DO'),
       ('IN HOLD'),
       ('IN PROGRESS'),
       ('IN REVIEW'),
       ('IN TESTING'),
       ('DONE / CLOSED');
