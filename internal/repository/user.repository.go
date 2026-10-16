package repository

import (
	"context"

	"github.com/rcontrisha/api-workers-ewallet/internal/model"
)

type UserRepo struct{}

func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

func (u *UserRepo) UpdateProfileInfo(ctx context.Context, userId string, db DBTX, payload model.Account) (model.Account, error) {
	sql := `UPDATE account
					SET fullname = COALESCE(NULLIF($1, ''),fullname),
							phone = COALESCE(NULLIF($2,''),phone),
							avatar = COALESCE(NULLIF($3, ''),avatar),
							updated_at = now()
					WHERE id = $4
					RETURNING
					fullname, phone, avatar, id`

	var account model.Account

	if err := db.QueryRow(ctx, sql,
		payload.FullName,
		payload.Phone,
		payload.Avatar,
		userId).Scan(&account.FullName, &account.Phone, &account.Avatar, &account.Id); err != nil {
		return model.Account{}, err
	}

	return account, nil

}

func (u *UserRepo) UpdateUserEmail(ctx context.Context, email, userId string, db DBTX) (string, error) {
	sql := `UPDATE users
				SET email = COALESCE(NULLIF($1,''),email)
				WHERE id = $2
				RETURNING email`

	var userEmail string
	if err := db.QueryRow(ctx, sql, email, userId).Scan(&userEmail); err != nil {
		return "", err
	}
	return userEmail, nil
}

func (u *UserRepo) GetProfile(ctx context.Context, db DBTX, userID string) (model.UserResponse, error) {
	query := `
		SELECT u.id, a.fullname, a.phone, a.avatar, u.email, a.created_at, a.updated_at
		FROM users u
		JOIN account a ON a.id = u.id
		WHERE u.id = $1`
	var user model.UserResponse
	err := db.QueryRow(ctx, query, userID).Scan(
		&user.ID,
		&user.Fullname,
		&user.Phone,
		&user.Avatar,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return model.UserResponse{}, err
	}
	return user, nil
}
