package http

import (
	"credit-layer/internal/application"

	"github.com/gin-gonic/gin"
)

func NewRouter(appUseCase *application.AppUseCase, creditLedgerUseCase *application.CreditLedgerUseCase) *gin.Engine {
	r := gin.Default()

	appHandler := NewAppHandler(appUseCase)
	creditLedgerHandler := NewCreditLedgerHandler(creditLedgerUseCase)

	v1 := r.Group("/v1")

	apps := v1.Group("/apps")
	apps.POST("", appHandler.CreateApp)

	credit := v1.Group("/credit")
	credit.Use(AuthMiddleware(appUseCase))
	credit.POST("/add", creditLedgerHandler.AddCredits)
	credit.POST("/deduct", creditLedgerHandler.DeductCredits)
	credit.GET("/balance", creditLedgerHandler.GetBalance)

	return r
}
