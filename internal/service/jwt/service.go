package jwt

import "github.com/jackc/pgx/v5/pgxpool"

type JWTService struct {
	Salt       string
	Pool       *pgxpool.Pool
	WebhookURL string
}
