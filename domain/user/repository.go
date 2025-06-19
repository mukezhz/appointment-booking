package user

import (
	"context"

	"clean-architecture/domain/common"
	"clean-architecture/domain/models"
	"clean-architecture/pkg/errorz"
	"clean-architecture/pkg/framework"
	"clean-architecture/pkg/infrastructure"
	"clean-architecture/pkg/types"
)

type Repository struct {
	db     infrastructure.Database
	logger framework.Logger
}

func NewRepository(
	db infrastructure.Database,
	logger framework.Logger,
) *Repository {
	return &Repository{
		db:     db,
		logger: logger}
}

func (r *Repository) Create(ctx context.Context, user *models.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		r.logger.Error("failed to create user: ", err)
		return common.HandleDBError(errorz.NewInternalError("failed to create user"), ErrUserMap)
	}
	return nil
}

func (r *Repository) FindByID(ctx context.Context, id types.BinaryUUID) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, "uuid = ?", id).Error; err != nil {
		r.logger.Error("failed to find user by ID: ", err)
		return nil, common.HandleDBError(err, ErrUserMap)
	}
	return &user, nil
}

func (r *Repository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.Model(&models.User{}).WithContext(ctx).First(&user, "email = ?", email).Error; err != nil {
		r.logger.Error("failed to find user by email: ", err)
		return nil, common.HandleDBError(err, ErrUserMap)
	}
	r.logger.Debug("found user by email: ", user.Email)
	return &user, nil
}

func (r *Repository) Update(ctx context.Context, user *models.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		r.logger.Error("failed to update user: ", err)
		return common.HandleDBError(err, ErrUserMap)
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, id types.BinaryUUID) error {
	if err := r.db.WithContext(ctx).Delete(&models.User{}, "uuid = ?", id).Error; err != nil {
		r.logger.Error("failed to delete user: ", err)
		return common.HandleDBError(err, ErrUserMap)
	}
	return nil
}
