package models

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/pkg/errors"
)

type (
	ProgressStatusTask struct {
		TaskID            uint64 `json:"taskId"`
		TaskTitle         string `json:"taskTitle"`
		AssigneeID        uint64 `json:"assigneeId"`
		AssigneeName      string `json:"assigneeName"`
		AssigneeAvatarURL string `json:"assigneeAvatarURL"`
	}
	ProgressStatusTasks []ProgressStatusTask

	ProgressStatusToCreate struct {
		ProjectID uint64 `json:"projectId" binding:"required"`
		Name      string `json:"name" binding:"required"`
		OrderNum  int    `json:"orderNum" binding:"required"`
	}
	ProgressStatus struct {
		ID        int64               `json:"id" db:"id"`
		ProjectID uint64              `json:"projectId" db:"project_id"`
		Name      string              `json:"name" db:"name"`
		OrderNum  int                 `json:"orderNum" db:"order_num"`
		Tasks     ProgressStatusTasks `json:"tasks" db:"ordered_tasks"`
	}
)

func (task ProgressStatusTask) Value() (driver.Value, error) {
	return json.Marshal(task)
}

func (task *ProgressStatusTask) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan into ProgressStatusTask: type assertion to []byte failed")
	}

	return json.Unmarshal(b, &task)
}

func (tasks ProgressStatusTasks) Value() (driver.Value, error) {
	return json.Marshal(tasks)
}

func (tasks *ProgressStatusTasks) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan into ProgressStatusTasks: type assertion to []byte failed")
	}

	return json.Unmarshal(b, &tasks)
}
