package handler

import (
	"github.com/LevOrlov5404/task-tracker/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) CreateImportanceStatus(c *gin.Context) {
	var status models.StatusToCreate
	if err := c.BindJSON(&status); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.ImportanceStatus.Create(c, status)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) GetAllImportanceStatuses(c *gin.Context) {

}

func (h *Handler) GetImportanceStatusByID(c *gin.Context) {

}

func (h *Handler) UpdateImportanceStatus(c *gin.Context) {

}

func (h *Handler) DeleteImportanceStatus(c *gin.Context) {

}

func (h *Handler) CreateProgressStatus(c *gin.Context) {
	var status models.StatusToCreate
	if err := c.BindJSON(&status); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.ProgressStatus.Create(c, status)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) GetAllProgressStatuses(c *gin.Context) {

}

func (h *Handler) GetProgressStatusByID(c *gin.Context) {

}

func (h *Handler) UpdateProgressStatus(c *gin.Context) {

}

func (h *Handler) DeleteProgressStatus(c *gin.Context) {

}
