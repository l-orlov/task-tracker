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
    email              VARCHAR(320) NOT NULL DEFAULT '' UNIQUE,
    firstname          VARCHAR(255) NOT NULL DEFAULT '',
    lastname           VARCHAR(255) NOT NULL DEFAULT '',
    password           VARCHAR(255) NOT NULL DEFAULT '',
    is_email_confirmed BOOLEAN      NOT NULL DEFAULT FALSE,
    avatar_url         VARCHAR(500) NOT NULL DEFAULT '',
    created_at         TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
CREATE TRIGGER set_timestamp
    BEFORE UPDATE
    ON r_user
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

-- projects to work on
CREATE TABLE r_project
(
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL DEFAULT '',
    description TEXT         NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    closed_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
CREATE TRIGGER set_timestamp
    BEFORE UPDATE
    ON r_project
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

-- importance_statuses for project tasks
CREATE TABLE s_project_importance_status
(
    id         SERIAL PRIMARY KEY,
    project_id BIGINT REFERENCES r_project (id) ON DELETE CASCADE,
    name       VARCHAR(255) NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    UNIQUE (project_id, name)
);
CREATE TRIGGER set_timestamp
    BEFORE UPDATE
    ON s_project_importance_status
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

-- progress_statuses for project tasks
CREATE TABLE s_project_progress_status
(
    id            SERIAL PRIMARY KEY,
    project_id    BIGINT REFERENCES r_project (id) ON DELETE CASCADE,
    name          VARCHAR(255) NOT NULL DEFAULT '' UNIQUE,
    order_num     INT          NOT NULL DEFAULT 0,
    ordered_tasks JSONB        NOT NULL DEFAULT '[]'::JSONB,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    UNIQUE (project_id, name)
);
CREATE TRIGGER set_timestamp
    BEFORE UPDATE
    ON s_project_progress_status
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

-- insert default statuses for new project
CREATE OR REPLACE FUNCTION trigger_insert_default_project_statuses()
    RETURNS TRIGGER AS
$$
BEGIN
    INSERT INTO s_project_importance_status (project_id, name)
    VALUES (NEW.id, 'LOW'),
           (NEW.id, 'MEDIUM'),
           (NEW.id, 'HIGH');

    INSERT INTO s_project_progress_status (project_id, name, order_num)
    VALUES (NEW.id, 'TO DO', 0),
           (NEW.id, 'IN PROGRESS', 1),
           (NEW.id, 'DONE', 2);

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER insert_default_project_statuses
    AFTER INSERT
    ON r_project
    FOR EACH ROW
EXECUTE PROCEDURE trigger_insert_default_project_statuses();

-- users working on projects
CREATE TABLE nn_project_user
(
    project_id BIGINT REFERENCES r_project (id) ON DELETE CASCADE NOT NULL,
    user_id    BIGINT REFERENCES r_user (id) ON DELETE CASCADE    NOT NULL,
    is_owner   BOOLEAN                                            NOT NULL DEFAULT FALSE,
    UNIQUE (project_id, user_id)
);

-- tasks to project
CREATE TABLE r_task
(
    id                   BIGSERIAL PRIMARY KEY,
    project_id           BIGINT REFERENCES r_project (id) ON DELETE CASCADE NOT NULL,
    title                VARCHAR(255)                                       NOT NULL DEFAULT '',
    description          TEXT                                               NOT NULL DEFAULT '',
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
