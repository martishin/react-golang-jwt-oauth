package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/martishin/react-golang-oauth/internal/models"
	"google.golang.org/api/oauth2/v2"
)

type UserRepository interface {
	FindOrCreateUser(ctx context.Context, userInfo *oauth2.Userinfo) (string, error)
	GetUserByID(ctx context.Context, userID string) (*models.User, error)
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (repo *userRepository) FindOrCreateUser(ctx context.Context, userInfo *oauth2.Userinfo) (string, error) {
	var userID string
	query := `SELECT id FROM users WHERE email = $1`

	err := repo.db.QueryRow(ctx, query, userInfo.Email).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			query = `INSERT INTO users (name, email, picture) VALUES ($1, $2, $3) RETURNING id`
			err = repo.db.QueryRow(ctx, query, userInfo.Name, userInfo.Email, userInfo.Picture).Scan(&userID)
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}

	return userID, nil
}

func (repo *userRepository) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	var user models.User
	query := `SELECT id, name, email, picture FROM users WHERE id = $1`

	err := repo.db.QueryRow(ctx, query, userID).Scan(&user.ID, &user.Name, &user.Email, &user.Picture)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
