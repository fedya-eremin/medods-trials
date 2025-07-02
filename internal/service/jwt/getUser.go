package jwt

import (
	"context"
	"encoding/base64"

	"example.com/m/internal/dto"
)

func (s *JWTService) GetUserByRefreshToken(ctx context.Context, refreshToken string) (*dto.User, error) {
	realRefreshToken, err := base64.RawStdEncoding.DecodeString(refreshToken)
	if err != nil {
		return nil, err
	}
	conn, err := s.Pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	userRow := conn.QueryRow(
		ctx,
		`select u.id from auth.user u
		join auth.refresh_token r on u.id = r.user_id
		where token_hash = $1
		limit 1`,
		realRefreshToken,
	)
	var user dto.User
	err = userRow.Scan(&user.ID)
	return &user, err
}
