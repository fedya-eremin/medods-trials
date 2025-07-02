package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/fedya-eremin/medods-trials/internal/dto"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("User not created with id %s", id)
		}
		return nil, err
	}
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

func (a *AuthService) GetUserByIDAndJTI(ctx context.Context, id uuid.UUID, jti uuid.UUID) (*dto.User, error) {
	conn, err := a.Pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	userRow := conn.QueryRow(
		ctx,
		`select u.id from auth.user u
		join auth.refresh_token r on u.id = r.user_id
		where u.id = $1 and r.id = $2 and r.used is false
		limit 1`,
		id,
		jti,
	)
	var user dto.User
	err = userRow.Scan(
		&user.ID,
	)

	return &user, err
}
