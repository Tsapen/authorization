package postgres

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

// Config contains settings for db
type Config struct {
	Connect string
}

// DB contains db connection.
type DB struct {
	*sql.DB
}

// CreateDBConnection create new connection with postgres.
func CreateDBConnection(c *Config) (*DB, error) {
	db, err := sql.Open("postgres", c.Connect)
	if err != nil {
		return nil, errors.Wrap(err, "can't open connection")
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

// RunExpiredTokensCleaning starts removing expired tokens from db.
func (db *DB) RunExpiredTokensCleaning() {
	var query = `DELETE FROM auth WHERE expires_at < localtimestamp;`

	for {
		time.Sleep(1 * time.Minute)
		if _, err := db.Exec(query); err != nil {
			log.Printf("can't clean expired tokens: %s\n", err)
		}
	}
}
