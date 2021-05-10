package models

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/pkg/errors"
)

type (
	ProjectBoardTask struct {
		TaskID            uint64 `json:"taskId" binding:"required"`
		TaskTitle         string `json:"taskTitle"`
		TaskOrderNum      int    `json:"taskOrderNum"`
		AssigneeID        uint64 `json:"assigneeId"`
		AssigneeFirstname string `json:"assigneeFirstname"`
		AssigneeLastname  string `json:"assigneeLastname"`
		AssigneeAvatarURL string `json:"assigneeAvatarURL"`
	}
	ProjectBoardProgressStatus struct {
		ProgressStatusId       int64  `json:"progressStatusId" binding:"required"`
		ProgressStatusName     string `json:"progressStatusName"`
		ProgressStatusOrderNum int    `json:"progressStatusOrderNum"`
	}
	ProjectBoardProgressStatusWithTasks struct {
		ProjectBoardProgressStatus
		Tasks []ProjectBoardTask `json:"tasks"`
	}
	ProjectBoardProgressStatusTasks []ProjectBoardTask
	ProjectBoardProgressStatuses    []ProjectBoardProgressStatus
	ProjectBoard                    []ProjectBoardProgressStatusWithTasks
)

func (tasks ProjectBoardProgressStatusTasks) Value() (driver.Value, error) {
	return json.Marshal(tasks)
}

func (tasks *ProjectBoardProgressStatusTasks) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan into ProjectBoardTasks: type assertion to []byte failed")
	}

	return json.Unmarshal(b, &tasks)
}

func (statuses ProjectBoardProgressStatuses) Value() (driver.Value, error) {
	return json.Marshal(statuses)
}

func (statuses *ProjectBoardProgressStatuses) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan into ProjectBoardProgressStatuses: type assertion to []byte failed")
	}

	return json.Unmarshal(b, &statuses)
}

func (pb ProjectBoard) Value() (driver.Value, error) {
	return json.Marshal(pb)
}

func (pb *ProjectBoard) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan into ProjectBoard: type assertion to []byte failed")
	}

	return json.Unmarshal(b, &pb)
}
