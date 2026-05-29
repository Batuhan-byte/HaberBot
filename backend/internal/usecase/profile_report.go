package usecase

import (
	"context"
	"errors"
	"fmt"

	"haberbot/internal/domain/entity"
	"haberbot/internal/domain/port"
)

// ReportProfileRequest defines the report input.
type ReportProfileRequest struct {
	ReporterUserID string
	ReportedUserID string
	Reason         string // 'spam', 'offensive_content', 'harassment', 'other'
	Comment        string
}

// ReportProfileUseCase coordinates report submission and validation.
type ReportProfileUseCase struct {
	profileRepo port.ProfileRepository
	userRepo    port.UserRepository
}

// NewReportProfileUseCase creates a new ReportProfileUseCase.
func NewReportProfileUseCase(profileRepo port.ProfileRepository, userRepo port.UserRepository) *ReportProfileUseCase {
	return &ReportProfileUseCase{
		profileRepo: profileRepo,
		userRepo:    userRepo,
	}
}

// Execute performs validation and files the report.
func (uc *ReportProfileUseCase) Execute(ctx context.Context, req ReportProfileRequest) error {
	if req.ReporterUserID == "" || req.ReportedUserID == "" {
		return errors.New("kullanıcı kimlikleri gereklidir")
	}

	if req.ReporterUserID == req.ReportedUserID {
		return errors.New("kendi profilinizi şikayet edemezsiniz")
	}

	// Validate reason
	validReasons := map[string]bool{
		"spam":              true,
		"offensive_content": true,
		"harassment":        true,
		"other":             true,
	}
	if !validReasons[req.Reason] {
		return errors.New("geçersiz şikayet nedeni")
	}

	// Validate target user exists
	targetUser, err := uc.userRepo.FindByID(ctx, req.ReportedUserID)
	if err != nil {
		return fmt.Errorf("checking target user: %w", err)
	}
	if targetUser == nil {
		return errors.New("şikayet edilecek kullanıcı bulunamadı")
	}

	// Rate limiting check (1 report per profile per user per 24 hours)
	hasRecent, err := uc.profileRepo.HasRecentReport(ctx, req.ReporterUserID, req.ReportedUserID)
	if err != nil {
		return fmt.Errorf("checking active reports: %w", err)
	}
	if hasRecent {
		return errors.New("bu profili son 24 saat içinde zaten şikayet ettiniz")
	}

	report := &entity.ProfileReport{
		ReporterUserID: req.ReporterUserID,
		ReportedUserID: req.ReportedUserID,
		Reason:         req.Reason,
		Comment:        req.Comment,
		Status:         "open",
	}

	return uc.profileRepo.SubmitReport(ctx, report)
}
