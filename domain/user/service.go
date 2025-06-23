package user

import (
	"context"

	"github.com/google/uuid"

	"github.com/mukezhz/appointment-booking/domain/constants"
	"github.com/mukezhz/appointment-booking/domain/models"
	"github.com/mukezhz/appointment-booking/pkg/errorz"
	"github.com/mukezhz/appointment-booking/pkg/framework"
	"github.com/mukezhz/appointment-booking/pkg/types"
)

type Service struct {
	repo   *Repository
	logger framework.Logger
	auth   AuthService
}

func NewService(
	repo *Repository,
	logger framework.Logger,
) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
		auth:   NewAuthService("hello"),
	}
}

func (s *Service) Register(
	ctx context.Context,
	email, firstName, lastName, firstNameJa, lastNameJa string,
	role constants.UserRole,
) (*models.User, error) {
	// Check if user already exists
	if _, err := s.repo.FindByEmail(ctx, email); err == nil {
		return nil, ErrEmailTaken
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return nil, errorz.NewInternalError("failed to generate uuid")
	}

	user := &models.User{
		UUID:        types.BinaryUUID(id),
		Email:       email,
		Role:        role,
		FirstName:   firstName,
		LastName:    lastName,
		FirstNameJa: firstNameJa,
		LastNameJa:  lastNameJa,
		IsActive:    true,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) VerifyCredentials(ctx context.Context, email string) (*models.User, error) {
	s.logger.Info("attempting to verify user credentials for email: ", email)
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	return user, nil
}

func (s *Service) GetUserByID(ctx context.Context, id types.BinaryUUID) (*models.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) UpdateUser(ctx context.Context, user *models.User) error {
	if err := s.repo.Update(ctx, user); err != nil {
		return err
	}
	return nil
}
