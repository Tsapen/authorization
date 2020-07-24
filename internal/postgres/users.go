package postgres

import (
	"github.com/Tsapen/authorization/internal/auth"
)

// CreateUser creates user in db.
func (db *DB) CreateUser(u auth.User) error {
	var q = `INSERT INTO users(email, login, password, phone_number) VALUES ($1, $2, $3, $4);`

	if _, err := db.Exec(q, u.Email, u.Login, u.Password, u.PhoneNumber); err != nil {
		return translateError(err)
	}

	return nil
}

// CheckUser checks if the user exists.
func (db *DB) CheckUser(a auth.Auth) error {
	var q = `SELECT 1 from users WHERE login = $1 AND password = $2;`

	var res int
	if err := db.QueryRow(q, a.Login, a.Password).Scan(&res); err != nil {
		return translateError(err)
	}

	return nil
}

// CreateToken creates token in db.
func (db *DB) CreateToken(td auth.TokenDetails) error {
	var q = `INSERT INTO auth(token, expires_at) VALUES ($1, $2), ($3, $4);`

	if _, err := db.Exec(q, td.AccessToken, td.AtExpires, td.RefreshToken, td.RtExpires); err != nil {
		return translateError(err)
	}

	return nil
}
