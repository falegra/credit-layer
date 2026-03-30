package http

import (
	"credit-layer/internal/application"
	"credit-layer/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreditLedgerHandler struct {
	creditLedgerUseCase *application.CreditLedgerUseCase
}

func NewCreditLedgerHandler(creditLedgerUseCase *application.CreditLedgerUseCase) *CreditLedgerHandler {
	return &CreditLedgerHandler{
		creditLedgerUseCase: creditLedgerUseCase,
	}
}

type creditRequest struct {
	UserID         string `json:"user_id" binding:"required"`
	Amount         int64  `json:"amount" binding:"required,gt=0"`
	Description    string `json:"description" binding:"required"`
	IdempotencyKey string `json:"idempotency_key" binding:"required"`
}

type creditResponse struct {
	ID             string `json:"id"`
	AppID          string `json:"app_id"`
	UserID         string `json:"user_id"`
	Amount         int64  `json:"amount"`
	Description    string `json:"description,omitempty"`
	IdempotencyKey string `json:"idempotency_key,omitempty"`
}

func toCreditResponse(cl *domain.CreditLedger) creditResponse {
	resp := creditResponse{
		ID:     cl.ID.String(),
		AppID:  cl.AppID.String(),
		UserID: cl.UserID,
		Amount: cl.Amount,
	}
	if cl.Description != nil {
		resp.Description = *cl.Description
	}
	if cl.IdempotencyKey != nil {
		resp.IdempotencyKey = *cl.IdempotencyKey
	}
	return resp
}

func (h *CreditLedgerHandler) AddCredits(c *gin.Context) {
	var req creditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	app := c.MustGet(appContextKey).(*domain.App)

	result, err := h.creditLedgerUseCase.AddCredits(
		c.Request.Context(),
		app.ID.String(),
		req.UserID,
		req.Amount,
		&req.Description,
		&req.IdempotencyKey,
	)
	if err != nil {
		switch err {
		case application.ErrInvalidAmount:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, toCreditResponse(result))
}

func (h *CreditLedgerHandler) DeductCredits(c *gin.Context) {
	var req creditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	app := c.MustGet(appContextKey).(*domain.App)

	result, err := h.creditLedgerUseCase.DeductCredits(
		c.Request.Context(),
		app.ID.String(),
		req.UserID,
		req.Amount,
		&req.Description,
		&req.IdempotencyKey,
	)
	if err != nil {
		switch err {
		case application.ErrInvalidAmount:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case application.ErrInsufficientCredits:
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, toCreditResponse(result))
}

func (h *CreditLedgerHandler) GetBalance(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	app := c.MustGet(appContextKey).(*domain.App)

	balance, err := h.creditLedgerUseCase.GetBalance(c.Request.Context(), app.ID.String(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": balance})
}
