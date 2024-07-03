package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/icrxz/crm-api-core/internal/application"
	"github.com/icrxz/crm-api-core/internal/domain"
)

type TransactionController struct {
	transactionService application.TransactionService
}

func NewTransactionController(transactionService application.TransactionService) TransactionController {
	return TransactionController{
		transactionService: transactionService,
	}
}

func (c *TransactionController) CreateTransaction(ctx *gin.Context) {
	caseID := ctx.Param("caseID")
	if caseID == "" {
		ctx.Error(domain.NewValidationError("param caseID cannot be empty", nil))
		return
	}

	var transactionDTO *CreateTransactionDTO
	err := ctx.BindJSON(&transactionDTO)
	if err != nil {
		ctx.Error(err)
		return
	}

	transaction, err := mapCreateTransactionDTOToTransaction(*transactionDTO, caseID)
	if err != nil {
		ctx.Error(err)
		return
	}

	transactionID, err := c.transactionService.CreateTransaction(ctx.Request.Context(), transaction)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(201, gin.H{"transaction_id": transactionID})
}

func (c *TransactionController) GetTransaction(ctx *gin.Context) {
	transactionID := ctx.Param("transactionID")
	if transactionID == "" {
		ctx.Error(domain.NewValidationError("param transactionID cannot be empty", nil))
		return
	}

	transaction, err := c.transactionService.GetTransaction(ctx.Request.Context(), transactionID)
	if err != nil {
		ctx.Error(err)
		return
	}

	transactionDTO := mapTransactionToTransactionDTO(transaction)

	ctx.JSON(200, transactionDTO)
}

func (c *TransactionController) UpdateTransaction(ctx *gin.Context) {
	transactionID := ctx.Param("transactionID")
	if transactionID == "" {
		ctx.Error(domain.NewValidationError("param transactionID cannot be empty", nil))
		return
	}

	var transactionUpdateDTO *TransactionUpdateDTO
	err := ctx.BindJSON(&transactionUpdateDTO)
	if err != nil {
		ctx.Error(err)
		return
	}

	transactionUpdate := mapTransactionUpdateDTOToTransactionUpdate(*transactionUpdateDTO)

	err = c.transactionService.UpdateTransaction(ctx.Request.Context(), transactionID, transactionUpdate)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.Status(204)
}

func (c *TransactionController) SearchTransactions(ctx *gin.Context) {
	filters := c.parseQueryToFilters(ctx)

	transactions, err := c.transactionService.SearchTransactions(ctx.Request.Context(), filters)
	if err != nil {
		ctx.Error(err)
		return
	}

	transactionDTOs := mapTransactionsToTransactionsDTO(transactions)

	ctx.JSON(200, transactionDTOs)
}

func (c *TransactionController) parseQueryToFilters(ctx *gin.Context) domain.TransactionFilters {
	filters := domain.TransactionFilters{}

	if caseIDs := ctx.QueryArray("case_id"); len(caseIDs) > 0 {
		filters.CaseIDs = caseIDs
	}

	if status := ctx.QueryArray("status"); len(status) > 0 {
		filters.Status = status
	}

	if types := ctx.QueryArray("type"); len(types) > 0 {
		filters.Types = types
	}

	return filters
}
