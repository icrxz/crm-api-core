package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/icrxz/crm-api-core/internal/application"
	"github.com/icrxz/crm-api-core/internal/domain"
)

type QueueController struct {
	queueService application.QueueService
}

func NewQueueController(queueService application.QueueService) QueueController {
	return QueueController{
		queueService: queueService,
	}
}

func (c *QueueController) CreateQueue(ctx *gin.Context) {
	var queueDTO *CreateQueueDTO
	if err := ctx.BindJSON(&queueDTO); err != nil {
		_ = ctx.Error(err)
		return
	}

	queue, err := mapCreateQueueDTOToQueue(*queueDTO)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	queueID, err := c.queueService.Create(ctx.Request.Context(), queue)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"queue_id": queueID})
}

func (c *QueueController) UpdateQueue(ctx *gin.Context) {
	queueID := ctx.Param("queueID")
	if queueID == "" {
		_ = ctx.Error(domain.NewValidationError("param queueID cannot be empty", nil))
		return
	}

	var queueDTO *UpdateQueueDTO
	if err := ctx.BindJSON(&queueDTO); err != nil {
		_ = ctx.Error(err)
		return
	}

	update := mapUpdateQueueDTOToUpdateQueue(*queueDTO)

	if err := c.queueService.Update(ctx.Request.Context(), queueID, update); err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (c *QueueController) GetQueue(ctx *gin.Context) {
	queueID := ctx.Param("queueID")
	if queueID == "" {
		_ = ctx.Error(domain.NewValidationError("param queueID cannot be empty", nil))
		return
	}

	queue, err := c.queueService.GetByID(ctx.Request.Context(), queueID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, mapQueueToQueueDTO(*queue))
}

func (c *QueueController) SearchQueues(ctx *gin.Context) {
	filters := c.parseQueryToFilters(ctx)

	queues, err := c.queueService.Search(ctx.Request.Context(), filters)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, mapSearchResultToSearchResultDTO(queues, mapQueuesToQueueDTOs))
}

func (c *QueueController) DeleteQueue(ctx *gin.Context) {
	queueID := ctx.Param("queueID")
	if queueID == "" {
		_ = ctx.Error(domain.NewValidationError("param queueID cannot be empty", nil))
		return
	}

	if err := c.queueService.Delete(ctx.Request.Context(), queueID); err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (c *QueueController) AddMember(ctx *gin.Context) {
	queueID := ctx.Param("queueID")
	if queueID == "" {
		_ = ctx.Error(domain.NewValidationError("param queueID cannot be empty", nil))
		return
	}

	var memberDTO *AddQueueMemberDTO
	if err := ctx.BindJSON(&memberDTO); err != nil {
		_ = ctx.Error(err)
		return
	}

	if err := c.queueService.AddMember(ctx.Request.Context(), queueID, memberDTO.UserID); err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (c *QueueController) RemoveMember(ctx *gin.Context) {
	queueID := ctx.Param("queueID")
	userID := ctx.Param("userID")
	if queueID == "" || userID == "" {
		_ = ctx.Error(domain.NewValidationError("params queueID and userID cannot be empty", nil))
		return
	}

	if err := c.queueService.RemoveMember(ctx.Request.Context(), queueID, userID); err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (c *QueueController) GetMembers(ctx *gin.Context) {
	queueID := ctx.Param("queueID")
	if queueID == "" {
		_ = ctx.Error(domain.NewValidationError("param queueID cannot be empty", nil))
		return
	}

	members, err := c.queueService.GetMembers(ctx.Request.Context(), queueID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, mapUsersToUserDTOs(members))
}

func (c *QueueController) parseQueryToFilters(ctx *gin.Context) domain.QueueFilters {
	filters := domain.QueueFilters{
		PagingFilter: domain.PagingFilter{
			Limit:  10,
			Offset: 0,
		},
	}

	if queueIDs := ctx.QueryArray("queue_id"); len(queueIDs) > 0 {
		filters.QueueID = queueIDs
	}

	if categories := ctx.QueryArray("category"); len(categories) > 0 {
		filters.Category = categories
	}

	if states := ctx.QueryArray("state"); len(states) > 0 {
		filters.State = states
	}

	if active := ctx.Query("active"); active != "" {
		activeBool, err := strconv.ParseBool(active)
		if err == nil {
			filters.Active = &activeBool
		}
	}

	if limitParam := ctx.Query("limit"); limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil {
			filters.Limit = parsedLimit
		}
	}

	if offsetParam := ctx.Query("offset"); offsetParam != "" {
		if parsedOffset, err := strconv.Atoi(offsetParam); err == nil {
			filters.Offset = parsedOffset
		}
	}

	return filters
}
