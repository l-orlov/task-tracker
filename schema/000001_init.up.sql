CREATE TABLE users
(
    id            bigserial    primary key,
    first_name    varchar(255) not null,
    last_name     varchar(255) not null,
    email         varchar(255) not null,
    password_hash varchar(255) not null,
    is_active     boolean      not null
);

CREATE TABLE importance_status
(
    id   serial       primary key,
    name varchar(255) not null
);

CREATE TABLE progress_status
(
    id   serial       primary key,
    name varchar(255) not null
);

CREATE TABLE projects
(
    id                   bigserial                   primary key,
    title                varchar(255)                not null,
    description          text                        not null,
    creation_date        timestamp with time zone    not null,
    assignee_id          bigint references users(id) not null,
    importance_status_id int references users(id)    not null,
    progress_status_id   int references users(id)    not null
);

CREATE TABLE tasks
(
    id                   bigserial                   primary key,
    title                varchar(255)                not null,
    description          text                        not null,
    creation_date        timestamp with time zone    not null,
    assignee_id          bigint references users(id) not null,
    importance_status_id int references users(id)    not null,
    progress_status_id   int references users(id)    not null
);

CREATE TABLE subtasks
(
    id                   bigserial                   primary key,
    title                varchar(255)                not null,
    description          text                        not null,
    creation_date        timestamp with time zone    not null,
    assignee_id          bigint references users(id) not null,
    importance_status_id int references users(id)    not null,
    progress_status_id   int references users(id)    not null
);

CREATE TABLE projects_tasks
(
    project_id bigint references projects(id) on delete cascade not null,
    task_id    bigint references tasks(id)    on delete cascade not null
);

CREATE TABLE tasks_subtasks
(
    task_id    bigint references tasks(id)    on delete cascade not null,
    subtask_id bigint references subtasks(id) on delete cascade not null
);
