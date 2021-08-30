package main

import (
	"context"

	"github.com/drifterz13/go-services/internal/common/db"
	"github.com/drifterz13/go-services/internal/common/models"
)

type UserRepository interface {
	Find(ctx context.Context) ([]*models.User, error)
	Create(ctx context.Context, email string) error
}

type userRepository struct {
	db *db.PostgresDBRepository
}

func NewUserRepository(postgresDB *db.PostgresDBRepository) UserRepository {
	return &userRepository{db: postgresDB}
}

func (r *userRepository) Find(ctx context.Context) ([]*models.User, error) {
	rows, err := r.db.Query(ctx, "SELECT * FROM users")
	if err != nil {
		return nil, err
	}

	var users []*models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

func (r *userRepository) Create(ctx context.Context, email string) error {
	err := r.db.Exec(ctx, "INSERT INTO users (email) VALUES $1", email)

	return err
}
