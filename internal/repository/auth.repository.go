package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rcontrisha/api-workers-ewallet/internal/model"
)

type AuthRepo struct {
	db *pgxpool.Pool
}

func NewAuthRepo(db *pgxpool.Pool) *AuthRepo {
	return &AuthRepo{
		db: db,
	}
}

func (a *AuthRepo) FindAccount(ctx context.Context, email string) (model.User, error) {
	query := "SELECT id, email, password FROM users WHERE email = $1"

	var user model.User

	err := a.db.QueryRow(ctx, query, email).Scan(
		&user.Id,
		&user.Email,
		&user.Password,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, err
		}

		return model.User{}, err
	}

	return user, nil
}
