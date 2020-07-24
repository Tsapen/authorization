package postgres

import (
	"database/sql"

	"github.com/Tsapen/authorization/internal/auth"

	"github.com/lib/pq"
)

const (
	dataException            = "22"
	uniqueConstrainViolation = "23"
)

func translateError(err error) error {
	if err == sql.ErrNoRows {
		return auth.ErrNotFound
	}

	if errType, ok := err.(*pq.Error); ok {
		var class = errType.Code.Class()

		if class == dataException ||
			class == uniqueConstrainViolation {
			return auth.ErrBadParameters
		}
	}

	return err
}

