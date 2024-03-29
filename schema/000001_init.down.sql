DROP TABLE IF EXISTS
    r_task,
    s_project_importance_status,
    s_project_progress_status,
    nn_project_user,
    r_project,
    r_user
    CASCADE;

DROP FUNCTION IF EXISTS
    trigger_insert_default_project_statuses(),
    trigger_set_r_task_order_num_in_progress_status(),
    get_project_board(BIGINT),
    update_project_board_parts(JSONB),
    update_project_board_progress_statuses(JSONB),
    update_project_board_progress_status_tasks(JSONB),
    trigger_set_timestamp(),
    CASCADE;
