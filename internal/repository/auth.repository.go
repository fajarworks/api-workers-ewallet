package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rcontrisha/api-workers-ewallet/internal/model"
)

type DBTX interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

type AuthRepo struct {}

func NewAuthRepo() *AuthRepo {
	return &AuthRepo{}
}

func (a *AuthRepo) FindAccount(ctx context.Context, db DBTX, email string) (model.Auth, error) {
	query := `
		SELECT u.id, u.email, u.password, (a.pin IS NOT NULL) AS has_pin
		FROM users u
		LEFT JOIN account a ON a.id = u.id
		WHERE u.email = $1`

	var acc model.Auth
	err := db.QueryRow(ctx, query, email).Scan(&acc.Id, &acc.Email, &acc.Password, &acc.HasPin)
	if err != nil {
		return model.Auth{}, err
	}
	return acc, nil
}

func (a *AuthRepo) CreateUser(ctx context.Context, db DBTX, payload model.User) (string, error) {
	query := "INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id"
	args := []any{payload.Email, payload.Password}

	var newUserId string
	err := db.QueryRow(ctx, query, args...).Scan(
		&newUserId,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", err
		}

		return "", err
	}

	return newUserId, nil
}

func (a *AuthRepo) CreateAccount(ctx context.Context, db DBTX, userId string) error {
	query := "INSERT INTO account (id) VALUES ($1)"

	cmd, err := db.Exec(ctx, query, userId)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("No Row(s) Affected")
	}

	return nil
}

func (a *AuthRepo) CreateWallet(ctx context.Context, db DBTX, userId string) error {
	query := "INSERT INTO wallet (id) VALUES ($1)"

	cmd, err := db.Exec(ctx, query, userId)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("No Row(s) Affected")
	}

	return nil
}

