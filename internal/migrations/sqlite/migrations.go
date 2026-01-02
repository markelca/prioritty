package sqlite

import (
	"database/sql"
	_ "embed"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"github.com/markelca/prioritty/pkg/items/repository/sqlite"
	"github.com/spf13/viper"
)

//go:embed sql/schema.sql
var SchemaSQL string

//go:embed sql/seed.sql
var SeedSQL string

func NewSQLiteRepository(dbPath string) (*sqlite.SQLiteRepository, error) {
	dbExists := false
	if _, err := os.Stat(dbPath); err == nil {
		dbExists = true
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	repo := sqlite.NewSQLiteRepository(db, dbPath)

	if !dbExists {
		if _, err := db.Exec(SchemaSQL); err != nil {
			log.Printf("Error executing schema: %v", err)
			db.Close()
			return nil, err
		}
		if viper.GetBool("demo") {
			if _, err := db.Exec(SeedSQL); err != nil {
				db.Close()
				log.Printf("Error executing seed data: %v", err)
				return nil, err
			}
		}
	}
	return repo, nil
}
