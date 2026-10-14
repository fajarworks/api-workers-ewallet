package service

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/rcontrisha/api-workers-ewallet/internal/dto"
	"github.com/rcontrisha/api-workers-ewallet/internal/repository"
	"github.com/rcontrisha/api-workers-ewallet/pkg"
)

type AuthService struct {
	repo *repository.AuthRepo
}

func NewAuthService(repo *repository.AuthRepo) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (a *AuthService) LoginService(ctx context.Context, payload dto.LoginRequest) (dto.LoginResponse, string, error) {
	result, err := a.repo.FindAccount(ctx, payload.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.LoginResponse{}, "", errors.New("Incorrect Email or Password.	")
		}
		return dto.LoginResponse{}, "", err
	}

	user := dto.LoginResponse{
		Id:    result.Id,
		Email: result.Email,
	}

	if err := pkg.Compare(payload.Password, result.Password); err != nil {
		log.Println(err)
		return dto.LoginResponse{}, "", err
	}

	claims := pkg.NewJWTClaims(result.Id, result.Email)
	token, err := claims.GenToken()

	return user, token, nil
}
