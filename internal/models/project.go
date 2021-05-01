package models

type (
	ProjectToCreate struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	ProjectToUpdate struct {
		ID uint64 `json:"id" binding:"required"`
		ProjectToCreate
		Owner uint64 `json:"owner" binding:"required"`
	}

	Project struct {
		ID          uint64 `json:"id" db:"id"`
		Name        string `json:"name" db:"name"`
		Description string `json:"description" db:"description"`
		Owner       uint64 `json:"owner" db:"owner"`
	}

	ProjectParams struct {
		ID          *uint64 `json:"id"`
		Name        *string `json:"name"`
		Description *string `json:"description"`
		Owner       *uint64 `json:"owner"`
	}

	ProjectWithTasks struct {
		Project
		Tasks []Task `json:"tasks"`
	}
)

func (project Project) ToProjectWithTasks() ProjectWithTasks {
	return ProjectWithTasks{
		Project: project,
	}
}
