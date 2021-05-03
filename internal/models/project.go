package models

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
	ProjectWithTasks struct {
		Project
		Tasks []Task `json:"tasks"`
	}
	ProjectImportanceStatusToAdd struct {
		ProjectID          uint64 `json:"projectId" binding:"required"`
		ImportanceStatusID int64  `json:"importanceStatusId" binding:"required"`
	}
	ProjectImportanceStatus struct {
		ID                 int64  `json:"id" db:"id"`
		ProjectID          uint64 `json:"projectId" db:"project_id"`
		ImportanceStatusID int64  `json:"importanceStatusId" db:"importance_status_id"`
	}
	ProjectProgressStatusToAdd struct {
		ProjectID        uint64 `json:"projectId" binding:"required"`
		ProgressStatusID int64  `json:"progressStatusId" binding:"required"`
	}
	ProjectProgressStatus struct {
		ID               int64  `json:"id" db:"id"`
		ProjectID        uint64 `json:"projectId" db:"project_id"`
		ProgressStatusID int64  `json:"progressStatusId" db:"progress_status_id"`
	}
	ProjectUser struct {
		ID        uint64 `json:"id" db:"id"`
		Email     string `json:"email" db:"email"`
		FirstName string `json:"firstName" db:"firstname"`
		LastName  string `json:"lastName" db:"lastname"`
		IsOwner   bool   `json:"isOwner" db:"is_owner"`
	}
)

func (project Project) ToProjectWithTasks() ProjectWithTasks {
	return ProjectWithTasks{
		Project: project,
	}
}
