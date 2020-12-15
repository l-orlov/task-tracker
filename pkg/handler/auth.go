package handler

import (
	"net/http"

	"github.com/LevOrlov5404/task-tracker/models"
	"github.com/gin-gonic/gin"
)

func (h * Handler) signUp(c *gin.Context) {
	var user models.User

	if err := c.BindJSON(&user); err != nil {
		newErrorRespnse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.CreateUser(user)
	if err != nil {
		newErrorRespnse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h * Handler) signIn(c *gin.Context) {

}
