package models

type (
	StatusToCreate struct {
		Name string `json:"name" binding:"required" db:"id"`
	}

	//ImportanceStatus struct {
	//	ID   int64  `json:"id" binding:"required" db:"id"`
	//	Name string `json:"name" binding:"required" db:"id"`
	//}
	//
	//ProgressStatus struct {
	//	ID   int64  `json:"id" binding:"required" db:"id"`
	//	Name string `json:"name" binding:"required" db:"id"`
	//}
)
