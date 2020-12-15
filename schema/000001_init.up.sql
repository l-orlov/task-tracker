CREATE TABLE users
(
    id            bigserial    primary key,
    first_name    varchar(255) not null,
    last_name     varchar(255) not null,
    email         varchar(255) not null,
    password_hash varchar(255) not null
);

CREATE TABLE projects
(
    id            bigserial                primary key,
    title         varchar(255)             not null,
    description   text                     not null,
    creation_date timestamp with time zone not null
);

CREATE TABLE tasks
(
    id            bigserial                primary key,
    title         varchar(255)             not null,
    description   text                     not null,
    creation_date timestamp with time zone not null
);

CREATE TABLE subtasks
(
    id            bigserial                primary key,
    title         varchar(255)             not null,
    description   text                     not null,
    creation_date timestamp with time zone not null
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