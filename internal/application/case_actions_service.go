package application

import (
	"context"
	"errors"
	"slices"
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type caseActionService struct {
	caseRepository        domain.CaseRepository
	caseHistoryRepository domain.CaseHistoryRepository
	transactionManager    domain.TransactionManager
	commentService        CommentService
	attachmentService     AttachmentService
	transactionService    TransactionService
	reportService         ReportService
}

//go:generate mockgen -source=case_actions_service.go -destination=mock_application/mock_case_actions_service.go -package=mock_application
type CaseActionService interface {
	ChangeOwner(ctx context.Context, caseID string, newOwner domain.ChangeOwner) error
	ChangeStatus(ctx context.Context, caseID string, newStatus domain.ChangeStatus) error
	ChangePartner(ctx context.Context, caseID string, newPartner domain.ChangePartner) error
	GenerateReport(ctx context.Context, caseID string) ([]byte, string, error)
	ResetCaseStatus(ctx context.Context, caseID, author string) error
}

func NewCaseActionService(
	caseRepository domain.CaseRepository,
	caseHistoryRepository domain.CaseHistoryRepository,
	transactionManager domain.TransactionManager,
	commentService CommentService,
	reportService ReportService,
	attachmentService AttachmentService,
	transactionService TransactionService,
) CaseActionService {
	return &caseActionService{
		caseRepository:        caseRepository,
		caseHistoryRepository: caseHistoryRepository,
		transactionManager:    transactionManager,
		commentService:        commentService,
		reportService:         reportService,
		attachmentService:     attachmentService,
		transactionService:    transactionService,
	}
}

func (c *caseActionService) recordHistory(ctx context.Context, caseID, eventName, author string, oldValues, newValues map[string]any) error {
	if len(oldValues) == 0 {
		return nil
	}

	history, err := domain.NewCaseHistory(caseID, eventName, author, oldValues, newValues)
	if err != nil {
		return err
	}

	return c.caseHistoryRepository.Create(ctx, history)
}

func (c *caseActionService) ChangeOwner(ctx context.Context, caseID string, newOwner domain.ChangeOwner) error {
	crmCase, err := c.caseRepository.GetByID(ctx, caseID)
	if err != nil {
		return err
	}

	caseUpdate := domain.CaseUpdate{
		OwnerID:   &newOwner.OwnerID,
		Status:    &newOwner.Status,
		UpdatedBy: newOwner.UpdatedBy,
	}

	eventName, oldValues, newValues := crmCase.DetectChanges(caseUpdate)

	crmCase.MergeUpdate(caseUpdate)

	return c.transactionManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := c.caseRepository.Update(txCtx, *crmCase); err != nil {
			return err
		}

		return c.recordHistory(txCtx, caseID, eventName, newOwner.UpdatedBy, oldValues, newValues)
	})
}

func (c *caseActionService) ChangeStatus(ctx context.Context, caseID string, newStatus domain.ChangeStatus) error {
	crmCase, err := c.caseRepository.GetByID(ctx, caseID)
	if err != nil {
		return err
	}

	caseUpdate := domain.CaseUpdate{
		Status:    &newStatus.Status,
		Type:      newStatus.Type,
		UpdatedBy: newStatus.UpdatedBy,
	}

	if newStatus.Status == domain.CLOSED {
		now := time.Now().UTC()
		caseUpdate.ClosedAt = &now
	}

	eventName, oldValues, newValues := crmCase.DetectChanges(caseUpdate)

	crmCase.MergeUpdate(caseUpdate)

	if newStatus.Content != nil {
		err = c.createChangeStatusComment(ctx, caseID, newStatus)
		if err != nil {
			return err
		}
	}

	return c.transactionManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := c.caseRepository.Update(txCtx, *crmCase); err != nil {
			return err
		}

		return c.recordHistory(txCtx, caseID, eventName, newStatus.UpdatedBy, oldValues, newValues)
	})
}

func (c *caseActionService) ChangePartner(ctx context.Context, caseID string, newPartner domain.ChangePartner) error {
	crmCase, err := c.caseRepository.GetByID(ctx, caseID)
	if err != nil {
		return err
	}

	caseUpdate := domain.CaseUpdate{
		PartnerID:  &newPartner.PartnerID,
		Status:     &newPartner.Status,
		TargetDate: &newPartner.TargetDate,
		UpdatedBy:  newPartner.UpdatedBy,
	}

	eventName, oldValues, newValues := crmCase.DetectChanges(caseUpdate)

	crmCase.MergeUpdate(caseUpdate)

	return c.transactionManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := c.caseRepository.Update(txCtx, *crmCase); err != nil {
			return err
		}

		return c.recordHistory(txCtx, caseID, eventName, newPartner.UpdatedBy, oldValues, newValues)
	})
}

func (c *caseActionService) GenerateReport(ctx context.Context, caseID string) ([]byte, string, error) {
	if caseID == "" {
		return nil, "", domain.NewValidationError("case_id is required", nil)
	}

	crmCase, err := c.caseRepository.GetByID(ctx, caseID)
	if err != nil {
		return nil, "", err
	}

	if !slices.Contains([]domain.CaseStatus{domain.REPORT, domain.PAYMENT, domain.RECEIPT, domain.CLOSED}, crmCase.Status) {
		return nil, "", domain.NewValidationError("case is not in status REPORT", map[string]any{"status": crmCase.Status})
	}

	return c.reportService.GenerateReport(ctx, *crmCase)
}

func (c *caseActionService) createChangeStatusComment(ctx context.Context, caseID string, newStatus domain.ChangeStatus) error {
	var commentType domain.CommentType
	switch newStatus.Status {
	case domain.WAITING_PARTNER:
		commentType = domain.COMMENT_CONTENT
	case domain.REPORT:
		commentType = domain.COMMENT_RESOLUTION
	case domain.PAYMENT:
		commentType = domain.COMMENT_REPORT
	case domain.REJECTED:
		commentType = domain.COMMENT_REJECTION
	}

	newComment, err := domain.NewComment(caseID, *newStatus.Content, newStatus.UpdatedBy, commentType, newStatus.Attachments)
	if err != nil {
		return err
	}

	_, err = c.commentService.Create(ctx, newComment)
	return err
}

func (c *caseActionService) ResetCaseStatus(ctx context.Context, caseID, author string) error {
	crmCase, err := c.caseRepository.GetByID(ctx, caseID)
	if err != nil {
		var customErr *domain.CustomError
		if errors.As(err, &customErr) && customErr.IsNotFound() {
			return nil
		}
		return err
	}

	emptyString := ""

	newStatus := domain.NEW
	cleanCase := domain.CaseUpdate{
		Status:    &newStatus,
		PartnerID: &emptyString,
		OwnerID:   &emptyString,
		// TargetDate: nil,
		// ClosedAt:   nil,
		UpdatedBy: author,
	}

	_, oldValues, newValues := crmCase.DetectChanges(cleanCase)

	crmCase.MergeUpdate(cleanCase)

	caseComments, err := c.commentService.GetByCaseID(ctx, caseID)
	if err != nil {
		return err
	}

	commentIDs := make([]string, 0, len(caseComments))
	for _, comment := range caseComments {
		commentIDs = append(commentIDs, comment.CommentID)
	}

	return c.transactionManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := c.attachmentService.DeleteByComments(txCtx, commentIDs); err != nil {
			return err
		}

		if err := c.transactionService.DeleteByCaseID(txCtx, caseID); err != nil {
			return err
		}

		if err := c.commentService.DeleteByCaseID(txCtx, caseID); err != nil {
			return err
		}

		if err := c.caseRepository.Update(txCtx, *crmCase); err != nil {
			return err
		}

		return c.recordHistory(txCtx, caseID, domain.CaseResetEvent, author, oldValues, newValues)
	})
}
