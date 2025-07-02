package jwt

import (
	"context"
	"encoding/base64"

	"github.com/google/uuid"
)

func (s *JWTService) DeleteRefreshTokenByID(ctx context.Context, id uuid.UUID) error {
	conn, err := s.Pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()
	var jti string
	err = conn.QueryRow(
		ctx,
		`delete from auth.refresh_token
		where id = $1
		returning id`,
		id,
	).Scan(&jti)
	return err
}

func (s *JWTService) DeleteRefreshToken(ctx context.Context, refreshToken string) error {
	realRefreshToken, err := base64.RawStdEncoding.DecodeString(refreshToken)
	conn, err := s.Pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()
	_, err = conn.Exec(
		ctx,
		`delete from auth.refresh_token
		where token_hash = $1`,
		realRefreshToken,
	)
	return err
}
