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
		ProgressStatusId       int64              `json:"progressStatusId" binding:"required"`
		ProgressStatusName     string             `json:"progressStatusName"`
		ProgressStatusOrderNum int                `json:"progressStatusOrderNum"`
		Tasks                  []ProjectBoardTask `json:"tasks"`
	}
	ProjectBoard []ProjectBoardProgressStatus
)

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
