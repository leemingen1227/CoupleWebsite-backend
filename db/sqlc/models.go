// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.0

package db

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Blog struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	PairID     int64     `json:"pair_id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Picture    string    `json:"picture"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

type Invitation struct {
	ID              int64     `json:"id"`
	InviterID       uuid.UUID `json:"inviter_id"`
	InviteeEmail    string    `json:"invitee_email"`
	InvitationToken string    `json:"invitation_token"`
	IsAccepted      bool      `json:"is_accepted"`
	CreateTime      time.Time `json:"create_time"`
}

type Pair struct {
	ID         int64        `json:"id"`
	CreateTime time.Time    `json:"create_time"`
	UpdateTime time.Time    `json:"update_time"`
	StartDate  sql.NullTime `json:"start_date"`
}

type Session struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type User struct {
	ID              uuid.UUID     `json:"id"`
	Email           string        `json:"email"`
	PasswordDigest  string        `json:"password_digest"`
	Name            string        `json:"name"`
	IsEmailVerified bool          `json:"is_email_verified"`
	PairID          sql.NullInt64 `json:"pair_id"`
	CreateTime      time.Time     `json:"create_time"`
	UpdateTime      time.Time     `json:"update_time"`
}

type UserPair struct {
	PairID int64     `json:"pair_id"`
	UserID uuid.UUID `json:"user_id"`
}

type VerifyEmail struct {
	ID          int64     `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Email       string    `json:"email"`
	SecretCode  string    `json:"secret_code"`
	IsUsed      bool      `json:"is_used"`
	CreateTime  time.Time `json:"create_time"`
	ExpiredTime time.Time `json:"expired_time"`
}
