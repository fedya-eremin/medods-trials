package dto

type WebhookRequest struct {
	UserID       string `json:"user_id"`
	OldUserAgent string `json:"old_user_agent"`
	NewUserAgent string `json:"new_user_agent"`
}
