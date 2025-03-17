package service

import (
	"context"

	"github.com/marioscordia/egov/internal/models"
)

type Service interface {
	CheckIIN(ctx context.Context, iin string) (*models.IINInfo, error)
	GetByIIN(ctx context.Context, iin string) (*models.User, error)
	GetBySearch(ctx context.Context, search string) ([]models.User, error)
	CreatePerson(ctx context.Context, person *models.User) error
}
