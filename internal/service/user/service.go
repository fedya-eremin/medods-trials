package user

import "github.com/jackc/pgx/v5/pgxpool"

type AuthService struct {
	Pool *pgxpool.Pool
}
