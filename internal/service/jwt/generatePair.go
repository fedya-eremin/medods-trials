package jwt

import (
	"context"
	"crypto/sha512"
	"encoding/base64"
	"time"

	"example.com/m/internal/dto"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *JWTService) generateAccessToken(userID uuid.UUID) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.RegisteredClaims{
		ID:        uuid.NewString(),
		Subject:   userID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(600 * time.Second)),
	})
	return accessToken.SignedString([]byte(s.Salt))
}

func (s *JWTService) generateAndStoreRefreshToken(ctx context.Context, userID uuid.UUID, accessToken string, userAgent string, ip string) (string, error) {
	claims, err := s.ParseAccessToken(accessToken)
	if err != nil {
		return "", err
	}
	shaAccessToken := sha512.Sum512([]byte(accessToken))
	hashedToken, err := bcrypt.GenerateFromPassword(shaAccessToken[:], 12)
	if err != nil {
		return "", err
	}
	conn, err := s.Pool.Acquire(ctx)
	if err != nil {
		return "", err
	}
	defer conn.Release()
	conn.Exec(
		ctx,
		`insert into auth.refresh_token (id, user_id, token_hash, used, user_agent, ip_addr, expires_at)
		values ($1, $2, $3, $4, $5, $6, $7)`,
		claims.ID,
		userID,
		hashedToken,
		false,
		userAgent,
		ip,
		claims.ExpiresAt.Time,
	)

	return base64.RawStdEncoding.EncodeToString(hashedToken), nil
}

func (s *JWTService) GenerateTokenPair(ctx context.Context, userID uuid.UUID, userAgent string, clientIP string) (*dto.TokenPair, error) {
	accessToken, err := s.generateAccessToken(userID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateAndStoreRefreshToken(
		ctx,
		userID,
		accessToken,
		userAgent,
		clientIP,
	)
	if err != nil {
		return nil, err
	}
	return &dto.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Lifetime:     600,
	}, nil
}
