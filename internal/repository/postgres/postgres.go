package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	errs "github.com/marioscordia/egov/internal/errors"
	"github.com/marioscordia/egov/internal/models"
	"github.com/marioscordia/egov/internal/repository"
)

type postgres struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) repository.Repository {
	return &postgres{
		db: db,
	}
}

func (p *postgres) GetByIIN(ctx context.Context, iin string) (*models.User, error) {
	stmt, err := p.db.PreparexContext(ctx, `select name, iin, phone from users where iin=$1`)
	if err != nil {
		return nil, err
	}

	var person models.User
	if err := stmt.GetContext(ctx, &person, iin); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}

	return &person, nil
}

func (p *postgres) GetBySearch(ctx context.Context, search string) ([]models.User, error) {
	stmt, err := p.db.PreparexContext(ctx, `select name, iin, phone from users where name like concat('%%',$1::text,'%%')`)
	if err != nil {
		return nil, err
	}

	var users []models.User
	if err := stmt.SelectContext(ctx, &users, search); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return users, errs.ErrNotFound
		}
		return nil, err
	}

	return users, nil
}

func (p *postgres) CreatePerson(ctx context.Context, user *models.User) error {
	stmt, err := p.db.PreparexContext(ctx, `select exists(
    											select 1
												from users
												where iin = $1 or phone=$2
											);`)
	if err != nil {
		return err
	}

	var exists bool
	if err := stmt.GetContext(ctx, &exists, user.IIN, user.Phone); err != nil {
		return err
	}

	if exists {
		return errs.ErrAlreadyExists
	}

	stmt, err = p.db.PreparexContext(ctx, `insert into users (name, iin, phone) values ($1, $2, $3)`)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, user.Name, user.IIN, user.Phone)
	return err
}
