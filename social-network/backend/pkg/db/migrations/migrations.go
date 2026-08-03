package migrations

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// Migrate runs all `*.up.sql` files in the provided migrationsPath in lexical order.
func Migrate(db *sql.DB, migrationsPath string) error {
	pattern := filepath.Join(migrationsPath, "*.up.sql")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("failed to glob migrations: %w", err)
	}

	sort.Strings(files)

	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", f, err)
		}

		if _, err := db.Exec(string(data)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", f, err)
		}
	}

	return nil
}
