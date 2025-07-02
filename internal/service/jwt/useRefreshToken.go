package jwt

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"example.com/m/internal/contextkeys"
	"example.com/m/internal/dto"
)

func (s *JWTService) UseRefreshToken(ctx context.Context, refreshToken string, userAgent string, ip string) error {
	realRefreshToken, err := base64.RawStdEncoding.DecodeString(refreshToken)
	if err != nil {
		return err
	}
	conn, err := s.Pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	var alreadyUsed bool
	var expiresAt time.Time
	var newUserAgent string
	var newIP string
	var userID string

	err = conn.QueryRow(
		ctx,
		`update auth.refresh_token
        set used = true
        where token_hash = $1
        returning 
            (select used from auth.refresh_token where token_hash = $1) as old_used,
            expires_at,
			user_agent,
			ip_addr,
			user_id`,
		realRefreshToken,
	).Scan(&alreadyUsed, &expiresAt, &newUserAgent, &newIP, &userID)

	if err != nil {
		return err
	}

	if alreadyUsed {
		return &dto.TokenUsedError{Token: refreshToken}
	}

	if time.Now().After(expiresAt) {
		return &dto.TokenExpiredError{Token: refreshToken}
	}

	if userAgent != newUserAgent {
		err = s.DeleteRefreshToken(ctx, refreshToken)
		if err != nil {
			return err
		}
		return &dto.WrongUserAgent{UserAgent: newUserAgent}
	}

	if ip != newIP {
		webhookRequestBody, err := json.Marshal(dto.WebhookRequest{
			UserID:       userID,
			OldUserAgent: userAgent,
			NewUserAgent: newUserAgent,
		})
		if err != nil {
			return nil
		}
		_, err = http.Post(s.WebhookURL, "application/json", bytes.NewBuffer(webhookRequestBody))
		if err != nil {
			logger := contextkeys.GetLogger(ctx)
			logger.Error("Error sending request to webhook", "err", err)
		}
	}

	return nil
}
