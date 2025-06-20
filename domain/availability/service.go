package availability

import (
	"clean-architecture/domain/models"
	"clean-architecture/pkg/errorz"
	"clean-architecture/pkg/framework"
	"context"
	"fmt"
	"time"
)

type Service struct {
	repository *Repository
	logger     framework.Logger
}

func NewService(
	repository *Repository,
	logger framework.Logger,
) *Service {
	return &Service{
		repository: repository,
		logger:     logger,
	}
}

func (s *Service) CreateAvailability(
	ctx context.Context,
	userID uint,
	req *CreateAvailabilityRequest,
) (*AvailabilityResponse, error) {
	// Validate time format
	_, err := time.Parse("15:04", req.StartTime)
	if err != nil {
		return nil, errorz.NewBadRequestError("Invalid start time format. Use HH:MM")
	}
	_, err = time.Parse("15:04", req.EndTime)
	if err != nil {
		return nil, errorz.NewBadRequestError("Invalid end time format. Use HH:MM")
	}

	// Create availability
	availability := &models.Availability{
		UserID:    userID,
		Weekday:   req.Weekday,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}

	if err := s.repository.Create(ctx, availability); err != nil {
		return nil, fmt.Errorf("failed to create availability: %w", err)
	}

	return ToAvailabilityResponse(availability), nil
}

func (s *Service) GetUserAvailability(
	ctx context.Context,
	userID uint,
) ([]*AvailabilityResponse, error) {
	availabilities, err := s.repository.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user availability: %w", err)
	}

	return ToAvailabilityResponseList(availabilities), nil
}
