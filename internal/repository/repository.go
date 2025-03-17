package repository

import (
	"context"

	"github.com/IceMAN2377/kaspitest/internal/models"
)

type Repository interface {
	GetByIIN(ctx context.Context, iin string) (*models.User, error)
	GetBySearch(ctx context.Context, search string) ([]models.User, error)
	CreatePerson(ctx context.Context, person *models.User) error
}
