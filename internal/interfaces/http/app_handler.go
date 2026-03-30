package http

import (
	"credit-layer/internal/application"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AppHandler struct {
	appUseCase *application.AppUseCase
}

func NewAppHandler(appUseCase *application.AppUseCase) *AppHandler {
	return &AppHandler{
		appUseCase: appUseCase,
	}
}

type createAppRequest struct {
	Name string `json:"name" binding:"required"`
}

type createAppResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	APIKey string `json:"api_key"`
}

func (h *AppHandler) CreateApp(c *gin.Context) {
	var req createAppRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	app, err := h.appUseCase.CreateApp(c.Request.Context(), req.Name)
	if err != nil {
		switch err {
		case application.ErrAppNameTaken:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, createAppResponse{
		ID:     app.ID.String(),
		Name:   app.Name,
		APIKey: app.APIKey,
	})
}
