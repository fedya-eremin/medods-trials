package user

import (
	"context"

	"example.com/m/internal/dto"
	"github.com/google/uuid"
)

func (a *AuthService) CreateUser(ctx context.Context, id uuid.UUID) (*dto.User, error) {
	conn, err := a.Pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	userRow := conn.QueryRow(
		ctx,
		`insert into auth.user (id)
		values ($1)
		returning id`,
		id,
	)
	var user dto.User
	err = userRow.Scan(
		&user.ID,
	)

	return &user, err
}

func (a *AuthService) GetUser(ctx context.Context, id uuid.UUID) (*dto.User, error) {
	conn, err := a.Pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	userRow := conn.QueryRow(
		ctx,
		`select id from auth.user where id = $1 limit 1`,
		id,
	)
	var user dto.User
	err = userRow.Scan(
		&user.ID,
	)

	return &user, err
}
