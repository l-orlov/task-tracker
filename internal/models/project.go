package models

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/pkg/errors"
)

type (
	ProjectToCreate struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}
	ProjectToUpdate struct {
		ID uint64 `json:"id" binding:"required"`
		ProjectToCreate
	}
	Project struct {
		ID          uint64 `json:"id" db:"id"`
		Name        string `json:"name" db:"name"`
		Description string `json:"description" db:"description"`
	}
	ProjectParams struct {
		ID          *uint64 `json:"id"`
		Name        *string `json:"name"`
		Description *string `json:"description"`
	}
	ProjectUser struct {
		ID        uint64 `json:"id" db:"id"`
		Email     string `json:"email" db:"email"`
		FirstName string `json:"firstName" db:"firstname"`
		LastName  string `json:"lastName" db:"lastname"`
		IsOwner   bool   `json:"isOwner" db:"is_owner"`
	}
	ProgressStatusTask struct {
		TaskID            uint64 `json:"taskId"`
		TaskTitle         string `json:"taskTitle"`
		AssigneeID        uint64 `json:"assigneeId"`
		AssigneeFirstname string `json:"assigneeFirstname"`
		AssigneeLastname  string `json:"assigneeLastname"`
		AssigneeAvatarURL string `json:"assigneeAvatarURL"`
	}
	ProjectBoard struct {
		ProgressStatusTasks []ProgressStatusTask `json:"get_project_board"`
	}
)

func (pst ProgressStatusTask) Value() (driver.Value, error) {
	return json.Marshal(pst)
}

func (pst *ProgressStatusTask) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan into ProgressStatusTask: type assertion to []byte failed")
	}

	return json.Unmarshal(b, &pst)
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
