package service

import (
	"context"

	"github.com/IceMAN2377/kaspitest/internal/models"
)

type Service interface {
	CheckIIN(ctx context.Context, iin string) (*models.IINInfo, error)
	GetByIIN(ctx context.Context, iin string) (*models.User, error)
	GetBySearch(ctx context.Context, search string) ([]models.User, error)
	CreatePerson(ctx context.Context, person *models.User) error
}
