CREATE TABLE system_user
(
    system_user_id BIGSERIAL PRIMARY KEY,
    firstname      VARCHAR(255) NOT NULL DEFAULT '',
    lastname       VARCHAR(255) NOT NULL DEFAULT '',
    email          VARCHAR(320) NOT NULL DEFAULT '' UNIQUE,
    password       VARCHAR(255) NOT NULL DEFAULT '',
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE project
(
    project_id  BIGSERIAL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL DEFAULT '',
    description TEXT         NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    closed_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE project_user
(
    project_user_id BIGSERIAL PRIMARY KEY,
    project_id      BIGINT REFERENCES project (project_id) ON DELETE CASCADE         NOT NULL,
    user_id         BIGINT REFERENCES system_user (system_user_id) ON DELETE CASCADE NOT NULL,
    is_owner        BOOLEAN                                                          NOT NULL DEFAULT FALSE,
    UNIQUE (project_id, user_id)
);

CREATE TABLE task_importance_status
(
    task_importance_status_id SERIAL PRIMARY KEY,
    name                      VARCHAR(255) NOT NULL DEFAULT '' UNIQUE
);

CREATE TABLE task_progress_status
(
    task_progress_status_id SERIAL PRIMARY KEY,
    name                    VARCHAR(255) NOT NULL DEFAULT '' UNIQUE
);

CREATE TABLE project_task_progress_status
(
    project_task_progress_status_id SERIAL PRIMARY KEY,
    task_progress_status_id         BIGINT REFERENCES task_progress_status (task_progress_status_id) ON DELETE CASCADE NOT NULL,
    project_id                      BIGINT REFERENCES project (project_id) ON DELETE CASCADE                           NOT NULL
);

CREATE TABLE project_task
(
    project_task_id                 BIGSERIAL PRIMARY KEY,
    title                           VARCHAR(255)                                                                                       NOT NULL DEFAULT '',
    description                     TEXT                                                                                               NOT NULL DEFAULT '',
    project_user_id                 BIGINT REFERENCES project_user (project_user_id) ON DELETE CASCADE                                 NOT NULL,
    task_importance_status_id       BIGINT REFERENCES task_importance_status (task_importance_status_id) ON DELETE CASCADE             NOT NULL,
    project_task_progress_status_id BIGINT REFERENCES project_task_progress_status (project_task_progress_status_id) ON DELETE CASCADE NOT NULL,
    created_at                      TIMESTAMPTZ                                                                                        NOT NULL DEFAULT NOW(),
    updated_at                      TIMESTAMPTZ                                                                                        NOT NULL DEFAULT NOW()
);

CREATE TABLE project_sprint
(
    project_sprint_id BIGSERIAL PRIMARY KEY,
    project_id        BIGINT REFERENCES project (project_id) ON DELETE CASCADE NOT NULL,
    created_at        TIMESTAMPTZ                                              NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ                                              NOT NULL DEFAULT NOW(),
    closed_at         TIMESTAMPTZ                                              NOT NULL DEFAULT NOW()
);

CREATE TABLE project_sprint_task
(
    task_id           BIGINT REFERENCES project_task (project_task_id) ON DELETE CASCADE     NOT NULL,
    project_sprint_id BIGINT REFERENCES project_sprint (project_sprint_id) ON DELETE CASCADE NOT NULL,
    added_at          TIMESTAMPTZ                                                            NOT NULL DEFAULT NOW()
);

CREATE TABLE task_change_log
(
    task_change_log_id           BIGSERIAL PRIMARY KEY,
    project_task_id              BIGINT REFERENCES project_task (project_task_id) ON DELETE CASCADE NOT NULL,
    project_task_progress_status INT                                                                NOT NULL,
    project_user                 BIGINT                                                             NOT NULL,
    updated_at                   TIMESTAMPTZ                                                        NOT NULL DEFAULT NOW()
);
