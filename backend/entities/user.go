package entities

import (
	"time"

	"github.com/lib/pq"
	"github.com/shekhar8352/PostEaze/constants"
	"github.com/shekhar8352/PostEaze/utils/database"
)

const (
	CreateUserWithFirebase = iota
	InsertRefreshToken
	GetUserByEmail
	GetUserByToken
	GetUserByID
	GetUserByFirebaseID
	UpdateUserPlatforms
	UpdateUser
	RevokeTokens
)

type User struct {
	ID           string         `json:"id"`
	FirebaseID   string         `json:"firebase_id"`
	Name         string         `json:"name"`
	Email        string         `json:"email"`
	Platforms    pq.StringArray `json:"platforms" db:"platforms"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	RefreshToken string         `json:"refresh_token"`
	ExpiresAt    time.Time      `json:"expire_at"`
}

func (o *User) GetQuery(code int) string {
	switch code {
	case CreateUserWithFirebase:
		return `INSERT INTO users (firebase_id, name, email, platforms) 
		        VALUES ($1, $2, $3, $4) 
		        RETURNING id, created_at, updated_at;`
	case InsertRefreshToken:
		return `INSERT INTO refresh_tokens (user_id, token, expires_at, revoked) 
		        VALUES ($1, $2, $3, $4);`
	case GetUserByEmail:
		return `SELECT id, firebase_id, name, email, platforms, created_at, updated_at 
		        FROM users WHERE email = $1;`
	case GetUserByToken:
		return `SELECT user_id FROM refresh_tokens 
		        WHERE token = $1 AND revoked = false AND expires_at > NOW();`
	case GetUserByID:
		return `SELECT id, firebase_id, name, email, platforms, created_at, updated_at 
		        FROM users WHERE id = $1;`
	case GetUserByFirebaseID:
		return `SELECT id, firebase_id, name, email, platforms, created_at, updated_at 
		        FROM users WHERE firebase_id = $1;`
	case UpdateUserPlatforms:
		return `UPDATE users SET platforms = $2, updated_at = NOW() 
		        WHERE id = $1 RETURNING updated_at;`
	case UpdateUser:
		return `UPDATE users SET email = $2, updated_at = NOW() 
		        WHERE id = $1 RETURNING updated_at;`
	case RevokeTokens:
		return `UPDATE refresh_tokens SET revoked = TRUE, updated_at = NOW() 
		        WHERE user_id = $1;`
	}
	return constants.Empty
}

func (o *User) GetQueryValues(code int) []any {
	switch code {
	case CreateUserWithFirebase:
		return []any{o.FirebaseID, o.Name, o.Email, pq.Array(o.Platforms)}
	case InsertRefreshToken:
		return []any{o.ID, o.RefreshToken, o.ExpiresAt, false}
	case GetUserByEmail:
		return []any{o.Email}
	case GetUserByToken:
		return []any{o.RefreshToken}
	case GetUserByID:
		return []any{o.ID}
	case GetUserByFirebaseID:
		return []any{o.FirebaseID}
	case UpdateUserPlatforms:
		return []any{o.ID, pq.Array(o.Platforms)}
	case UpdateUser:
		return []any{o.ID, o.Email}
	case RevokeTokens:
		return []any{o.ID}
	}
	return nil
}

func (o *User) GetMultiQuery(code int) string {
	switch code {
	}
	return constants.Empty
}

func (o *User) GetMultiQueryValues(code int) []any {
	switch code {
	}
	return nil
}

func (o *User) GetNextRaw() database.RawEntity {
	return new(User)
}

func (o *User) BindRawRow(code int, row database.Scanner) error {
	switch code {
	case CreateUserWithFirebase:
		return row.Scan(&o.ID, &o.CreatedAt, &o.UpdatedAt)

	case GetUserByEmail:
		return row.Scan(&o.ID, &o.FirebaseID, &o.Name, &o.Email, &o.Platforms, &o.CreatedAt, &o.UpdatedAt)

	case GetUserByToken:
		return row.Scan(&o.ID)

	case GetUserByID:
		return row.Scan(&o.ID, &o.FirebaseID, &o.Name, &o.Email, &o.Platforms, &o.CreatedAt, &o.UpdatedAt)

	case GetUserByFirebaseID:
		return row.Scan(&o.ID, &o.FirebaseID, &o.Name, &o.Email, &o.Platforms, &o.CreatedAt, &o.UpdatedAt)

	case UpdateUserPlatforms:
		return row.Scan(&o.UpdatedAt)

	case UpdateUser:
		return row.Scan(&o.UpdatedAt)

	}
	return nil
}

func (o *User) GetExec(code int) string {
	switch code {
	default:
		return constants.Empty
	}
}

func (o *User) GetExecValues(code int, _ string) []any {
	return nil
}
