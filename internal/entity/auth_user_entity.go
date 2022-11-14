package entity

import "time"

// UserCacheInfo User Cache Information
type UserCacheInfo struct {
	UserID      string    `json:"user_id"`
	UserStatus  int       `json:"user_status"`
	EmailStatus int       `json:"email_status"`
	ExpiredAt   time.Time `json:"expired_at"`
}
