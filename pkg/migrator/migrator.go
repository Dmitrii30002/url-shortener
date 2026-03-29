package migrator

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func Up(db *sql.DB, dir string) error {
	if dir == "" {
		dir = "."
	}

	if err := goose.Up(db, dir); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
