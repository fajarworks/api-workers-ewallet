package service

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rcontrisha/api-workers-ewallet/internal/dto"
	"github.com/rcontrisha/api-workers-ewallet/internal/model"
	"github.com/rcontrisha/api-workers-ewallet/internal/repository"
	"github.com/rcontrisha/api-workers-ewallet/pkg"
)

type AuthService struct {
	repo *repository.AuthRepo
	db   *pgxpool.Pool
}

func NewAuthService(repo *repository.AuthRepo, db *pgxpool.Pool) *AuthService {
	return &AuthService{
		repo: repo,
		db:   db,
	}
}

func (a *AuthService) LoginService(ctx context.Context, payload dto.LoginRequest) (dto.LoginResponse, string, error) {
	result, err := a.repo.FindAccount(ctx, a.db, payload.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.LoginResponse{}, "", errors.New("Incorrect Email or Password.	")
		}
		return dto.LoginResponse{}, "", err
	}

	user := dto.LoginResponse{
		Id:     result.Id,
		Email:  result.Email,
		HasPin: result.HasPin,
	}

	if err := pkg.Compare(payload.Password, result.Password); err != nil {
		log.Println(err)
		return dto.LoginResponse{}, "", err
	}

	claims := pkg.NewJWTClaims(result.Id, result.Email)
	token, err := claims.GenToken()

	return user, token, nil
}

func (a *AuthService) RegisterService(ctx context.Context, payload dto.LoginRequest) error {
	_, err := a.repo.FindAccount(ctx, a.db, payload.Email)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}
	if err == nil {
		return errors.New("User Already Exist.")
	}

	hashedPwd := pkg.NewRecommendedHashConfig().GenHash(payload.Password)

	tx, err := a.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			log.Println(err.Error())
		}
	}()

	var user = model.User{
		Email:    payload.Email,
		Password: hashedPwd,
	}

	userId, err := a.repo.CreateUser(ctx, tx, user)
	if err != nil {
		log.Println("[Service-CreateUser]-Error: ", err)
		return err
	}
	if err := a.repo.CreateAccount(ctx, tx, userId); err != nil {
		log.Println("[Service-CreateAccount]-Error: ", err)
		return err
	}
	if err := a.repo.CreateWallet(ctx, tx, userId); err != nil {
		log.Println("[Service-CreateWallet]-Error: ", err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
