package application

import (
	"context"
	"slices"
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type caseActionService struct {
	caseRepository domain.CaseRepository
	commentService CommentService
	reportService  ReportService
}

type CaseActionService interface {
	ChangeOwner(ctx context.Context, caseID string, newOwner domain.ChangeOwner) error
	ChangeStatus(ctx context.Context, caseID string, newStatus domain.ChangeStatus) error
	ChangePartner(ctx context.Context, caseID string, newPartner domain.ChangePartner) error
	GenerateReport(ctx context.Context, caseID string) ([]byte, string, error)
}

func NewCaseActionService(
	caseRepository domain.CaseRepository,
	commentService CommentService,
	reportService ReportService,
) CaseActionService {
	return &caseActionService{
		caseRepository: caseRepository,
		commentService: commentService,
		reportService:  reportService,
	}
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
	crmCase.MergeUpdate(caseUpdate)

	return c.caseRepository.Update(ctx, *crmCase)
}

func (c *caseActionService) ChangeStatus(ctx context.Context, caseID string, newStatus domain.ChangeStatus) error {
	crmCase, err := c.caseRepository.GetByID(ctx, caseID)
	if err != nil {
		return err
	}

	caseUpdate := domain.CaseUpdate{
		Status:    &newStatus.Status,
		UpdatedBy: newStatus.UpdatedBy,
	}

	if newStatus.Status == domain.CLOSED {
		now := time.Now().UTC()
		caseUpdate.ClosedAt = &now
	}

	crmCase.MergeUpdate(caseUpdate)

	if newStatus.Content != nil {
		err = c.createChangeStatusComment(ctx, caseID, newStatus)
		if err != nil {
			return err
		}
	}

	return c.caseRepository.Update(ctx, *crmCase)
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
	crmCase.MergeUpdate(caseUpdate)

	return c.caseRepository.Update(ctx, *crmCase)
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
