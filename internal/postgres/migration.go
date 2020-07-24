package postgres

import (
	"database/sql"

	"github.com/pkg/errors"
)

type dbMigration struct {
	number  int
	command string
}

func dbMigrations() []dbMigration {
	return []dbMigration{
		{
			number: 1,
			command: `
			CREATE TABLE IF NOT EXISTS users(
				id				SERIAL NOT NULL PRIMARY KEY
			,	login	 		VARCHAR(50) NOT NULL UNIQUE
			,	email			VARCHAR(50) NOT NULL UNIQUE
			,	password		VARCHAR(50) NOT NULL 
			,	phone_number	VARCHAR(12) NOT NULL UNIQUE
			);`,
		},

		{
			number: 2,
			command: `
			CREATE UNIQUE INDEX ON users(login) USING hash;
			`,
		},

		{
			number: 3,
			command: `
			CREATE TABLE IF NOT EXISTS auth(
				id				SERIAL NOT NULL PRIMARY KEY
			,	token	 		VARCHAR(164) NOT NULL UNIQUE
			,	expires_at		TIMESTAMP NOT NULL UNIQUE
			);`,
		},

		{
			number: 4,
			command: `
			CREATE UNIQUE INDEX ON auth(token) USING hash;
			`,
		},
	}
}

func (s *DB) applied(num int) (bool, error) {
	var ex = 0
	var q = `SELECT 1 FROM migrations WHERE num = $1;`
	var err = s.QueryRow(q, num).Scan(&ex)
	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *DB) apply(num int) error {
	var q = `INSERT INTO migrations(num) VALUES ($1);`
	var _, err = s.Exec(q, num)
	if err != nil {
		return err
	}

	return nil
}

// Migrate prepares db to work.
func (s *DB) Migrate() error {
	var q = `CREATE TABLE IF NOT EXISTS migrations(
				num 		INT NOT NULL PRIMARY KEY,
				created_at	TIMESTAMP NOT NULL DEFAULT LOCALTIMESTAMP
			);`
	var _, err = s.Exec(q)
	if err != nil {
		return errors.Wrap(err, "can't create migrations table: ")
	}

	var migs = dbMigrations()
	for _, mig := range migs {
		var migrated, err = s.applied(mig.number)
		if err != nil {
			return errors.Wrapf(err, "can't check %d migration: ", mig.number)
		}

		if migrated {
			continue
		}

		if _, err := s.Exec(mig.command); err != nil {
			errors.Wrapf(err, "can't apply %d migration: ", mig.number)
		}

		if err := s.apply(mig.number); err != nil {
			errors.Wrapf(err, "can't set %d migration as applied: ", mig.number)
		}
	}
	return nil
}
