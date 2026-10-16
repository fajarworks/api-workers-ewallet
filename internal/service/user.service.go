package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rcontrisha/api-workers-ewallet/internal/dto"
	"github.com/rcontrisha/api-workers-ewallet/internal/model"
	"github.com/rcontrisha/api-workers-ewallet/internal/repository"
)

type UserService struct {
	repo *repository.UserRepo
	db   *pgxpool.Pool
}

func NewUserService(repo *repository.UserRepo, db *pgxpool.Pool) *UserService {
	return &UserService{
		repo: repo,
		db:   db,
	}

}

func (u *UserService) UpdateProfileInfo(ctx context.Context, userId, email string, payload dto.UserUpdateResponse) (dto.UserUpdateResponse, error) {
	tx, err := u.db.Begin(ctx)
	if err != nil {
		return dto.UserUpdateResponse{}, err
	}

	defer tx.Rollback(ctx)

	data := model.Account{
		Id:       payload.Id,
		FullName: payload.FullName,
		Phone:    payload.Phone,
		Avatar:   payload.AvatarUrl,
	}

	account, err := u.repo.UpdateProfileInfo(ctx, userId, tx, data)
	if err != nil {
		return dto.UserUpdateResponse{}, err
	}

	userEmail, err := u.repo.UpdateUserEmail(ctx, email, userId, tx)
	if err != nil {
		return dto.UserUpdateResponse{}, err
	}

	profile := dto.UserUpdateResponse{
		Id:        account.Id,
		FullName:  account.FullName,
		Phone:     account.Phone,
		AvatarUrl: account.Avatar,
		Email:     userEmail,
	}

	if err := tx.Commit(ctx); err != nil {
		return dto.UserUpdateResponse{}, err
	}
	return profile, nil

}
