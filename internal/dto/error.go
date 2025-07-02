package dto

import "fmt"

type TokenUsedError struct {
	Token string
}

func (e *TokenUsedError) Error() string {
	return fmt.Sprintf("Token already used: %s", e.Token)
}

type TokenExpiredError struct {
	Token string
}

func (e *TokenExpiredError) Error() string {
	return fmt.Sprintf("Token already expired: %s", e.Token)
}

type WrongUserAgent struct {
	UserAgent string
}

func (e *WrongUserAgent) Error() string {
	return fmt.Sprintf("Wrong user agent detected: %s. Aborting", e.UserAgent)
}
