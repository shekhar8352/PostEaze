package entities

import (
	"time"

	"github.com/shekhar8352/PostEaze/constants"
	"github.com/shekhar8352/PostEaze/utils/database"
)

const (
	CreateUser = iota
	InsertRefreshToken
	GetUserByEmail
	GetUserByToken
	GetUserByID
	RevokeTokens
)

type User struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Password     string    `json:"-"`
	UserType     string    `json:"user_type"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expire_at"`
}

func (o *User) GetQuery(code int) string {
	switch code {
	case CreateUser:
		return `INSERT INTO users (name, email, password, user_type) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id, created_at, updated_at;`
	case InsertRefreshToken:
		return `INSERT INTO refresh_tokens (user_id, token, expires_at, revoked ) 
		VALUES ($1, $2, $3, $4 ) ;`
	case GetUserByEmail:
		return `SELECT id , name , password, user_type, created_at , updated_at FROM users WHERE email = $1 ;`
	case GetUserByToken:
		return `SELECT user_id FROM refresh_tokens WHERE token = $1 AND revoked = false AND expires_at > NOW() ;`
	case GetUserByID:
		return `SELECT id, name, email , user_type WHERE id = $1 ;`
	case RevokeTokens:
		return `UPDATE refresh_tokens SET revoked = TRUE, updated_at = NOW() WHERE user_id = $1;`

	}
	return constants.Empty
}

func (o *User) GetQueryValues(code int) []any {
	switch code {
	case CreateUser:
		return []interface{}{o.Name, o.Email, o.Password, o.UserType}
	case InsertRefreshToken:
		return []interface{}{o.ID, o.RefreshToken, o.ExpiresAt, false}
	case GetUserByEmail:
		return []interface{}{o.Email}
	case GetUserByToken:
		return []interface{}{o.RefreshToken}
	case GetUserByID:
		return []interface{}{o.ID}
	case RevokeTokens:
		return []interface{}{o.ID}
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
	case CreateUser:
		row.Scan(&o.ID, &o.CreatedAt, &o.UpdatedAt)
	case GetUserByEmail:
		row.Scan(&o.ID, &o.Name, &o.Password, &o.UserType, &o.CreatedAt, &o.UpdatedAt)
	case GetUserByToken:
		row.Scan(&o.ID)
	case GetUserByID:
		row.Scan(&o.ID, &o.Name, &o.Email, &o.UserType)
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
