-- naming relied ON this article: http://citforum.ru/database/articles/naming_rule/

CREATE OR REPLACE FUNCTION trigger_set_timestamp()
    RETURNS TRIGGER
    LANGUAGE plpgsql
AS
$$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;

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

-- projects to work ON
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
    name          VARCHAR(255) NOT NULL DEFAULT '',
    order_num     INT          NOT NULL DEFAULT 0,
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
    RETURNS TRIGGER
    LANGUAGE plpgsql
AS
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

    RETURN NEW;
END;
$$;

CREATE TRIGGER insert_default_project_statuses
    AFTER INSERT
    ON r_project
    FOR EACH ROW
EXECUTE PROCEDURE trigger_insert_default_project_statuses();

-- users working ON projects
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
    id                           BIGSERIAL PRIMARY KEY,
    project_id                   BIGINT REFERENCES r_project (id) ON DELETE CASCADE NOT NULL,
    title                        VARCHAR(255)                                       NOT NULL DEFAULT '',
    description                  TEXT                                               NOT NULL DEFAULT '',
    assignee_id                  BIGINT REFERENCES r_user (id)                      NOT NULL,
    importance_status_id         INT REFERENCES s_project_importance_status (id)    NOT NULL,
    progress_status_id           INT REFERENCES s_project_progress_status (id)      NOT NULL,
    order_num_in_progress_status INT                                                NOT NULL DEFAULT 0,
    created_at                   TIMESTAMPTZ                                        NOT NULL DEFAULT NOW(),
    updated_at                   TIMESTAMPTZ                                        NOT NULL DEFAULT NOW()
);
CREATE TRIGGER set_timestamp
    BEFORE UPDATE
    ON r_task
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();
CREATE INDEX idx_r_task_project_id ON r_task (project_id);

CREATE OR REPLACE FUNCTION trigger_set_r_task_order_num_in_progress_status()
    RETURNS TRIGGER
    LANGUAGE plpgsql
AS
$$
DECLARE
    rec RECORD;
BEGIN
    FOR rec IN SELECT *
               FROM r_task
               WHERE order_num_in_progress_status >= NEW.order_num_in_progress_status
                 AND progress_status_id = NEW.progress_status_id
               ORDER BY order_num_in_progress_status
        LOOP
            UPDATE r_task
            SET order_num_in_progress_status = order_num_in_progress_status + 1
            WHERE id = rec.id;
        END LOOP;

    RETURN NEW;
END;
$$;

CREATE TRIGGER set_r_task_order_num_in_progress_status
    BEFORE INSERT
    ON r_task
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_r_task_order_num_in_progress_status();

CREATE OR REPLACE FUNCTION get_project_board(_project_id BIGINT)
    RETURNS JSONB
    LANGUAGE plpgsql
AS
$$
BEGIN
    RETURN (SELECT COALESCE(jsonb_agg(
                                    jsonb_build_object(
                                            'progressStatusId', spps.id,
                                            'progressStatusName', spps.name,
                                            'progressStatusOrderNum', spps.order_num,
                                            'tasks', COALESCE(t.tasks, '[]'::JSONB)
                                        )
                                    ORDER BY (spps.order_num)
                                ), '[]'::JSONB) board
            FROM s_project_progress_status spps
                     LEFT JOIN LATERAL (
                SELECT rt.progress_status_id,
                       jsonb_agg(
                               jsonb_build_object(
                                       'taskId', rt.id,
                                       'taskTitle', rt.title,
                                       'taskOrderNum', rt.order_num_in_progress_status,
                                       'assigneeId', rt.assignee_id,
                                       'assigneeFirstname', ru.firstname,
                                       'assigneeLastname', ru.lastname,
                                       'assigneeAvatarURL', ru.avatar_url
                                   )
                               ORDER BY (rt.order_num_in_progress_status)
                           ) tasks
                FROM r_task rt
                         INNER JOIN r_user ru ON ru.id = rt.assignee_id
                GROUP BY rt.progress_status_id
                ) t ON spps.id = t.progress_status_id
            WHERE spps.project_id = _project_id);
END;
$$;

CREATE OR REPLACE FUNCTION update_project_board_parts(_board JSONB)
    RETURNS VOID
    LANGUAGE plpgsql
AS
$$
DECLARE
    rec RECORD;
BEGIN
    FOR rec IN SELECT board."progressStatusId", tasks."taskId", tasks."taskOrderNum"
               FROM jsonb_to_recordset(_board) AS board("progressStatusId" INT, "tasks" JSONB)
                        LEFT JOIN jsonb_to_recordset(board."tasks") AS tasks("taskId" BIGINT, "taskOrderNum" INT)
                                  ON TRUE
        LOOP
            IF rec."progressStatusId" IS NOT NULL AND
               rec."taskId" IS NOT NULL AND rec."taskOrderNum" IS NOT NULL THEN
                UPDATE r_task
                SET progress_status_id           = rec."progressStatusId",
                    order_num_in_progress_status = rec."taskOrderNum"
                WHERE id = rec."taskId";
            END IF;
        END LOOP;
END;
$$;

CREATE OR REPLACE FUNCTION update_project_board_progress_statuses(_progress_statuses JSONB)
    RETURNS VOID
    LANGUAGE plpgsql
AS
$$
DECLARE
    rec RECORD;
BEGIN
    FOR rec IN SELECT board."progressStatusId", board."progressStatusOrderNum"
               FROM jsonb_to_recordset(_progress_statuses) AS board("progressStatusId" INT, "progressStatusOrderNum" INT)
        LOOP
            IF rec."progressStatusId" IS NOT NULL AND rec."progressStatusOrderNum" IS NOT NULL THEN
                UPDATE s_project_progress_status
                SET order_num = rec."progressStatusOrderNum"
                WHERE id = rec."progressStatusId";
            END IF;
        END LOOP;
END;
$$;

CREATE OR REPLACE FUNCTION update_project_board_progress_status_tasks(_tasks JSONB)
    RETURNS VOID
    LANGUAGE plpgsql
AS
$$
DECLARE
    rec RECORD;
BEGIN
    FOR rec IN SELECT tasks."taskId", tasks."taskOrderNum"
               FROM jsonb_to_recordset(_tasks) AS tasks("taskId" BIGINT, "taskOrderNum" INT)
        LOOP
            IF rec."taskId" IS NOT NULL AND rec."taskOrderNum" IS NOT NULL THEN
                UPDATE r_task
                SET order_num_in_progress_status = rec."taskOrderNum"
                WHERE id = rec."taskId";
            END IF;
        END LOOP;
END;
$$;
