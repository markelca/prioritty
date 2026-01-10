package sqlite

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

const testSchemaSQL = `
CREATE TABLE status (
   id INTEGER PRIMARY KEY,
   name TEXT NOT NULL UNIQUE
);

CREATE TABLE tag (
   id INTEGER PRIMARY KEY,
   name TEXT NOT NULL UNIQUE
);

CREATE TABLE task (
   id INTEGER PRIMARY KEY,
   title TEXT NOT NULL,
   body TEXT,
   status_id INTEGER NOT NULL,
   tag_id INTEGER,
   created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
   FOREIGN KEY (status_id) REFERENCES status(id)
   FOREIGN KEY (tag_id) REFERENCES tag(id)
);

CREATE TABLE note (
   id INTEGER PRIMARY KEY,
   title TEXT NOT NULL,
   body TEXT,
   tag_id INTEGER,
   created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
   FOREIGN KEY (tag_id) REFERENCES tag(id)
);

INSERT INTO status (id, name) VALUES
   (0, 'Pending'),
   (1, 'In Progress'),
   (2, 'Completed'),
   (3, 'Cancelled');
`

func setupTestDB(t *testing.T) *SQLiteRepository {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	_, err = db.Exec(testSchemaSQL)
	require.NoError(t, err)

	repo := NewSQLiteRepository(db, ":memory:")

	t.Cleanup(func() {
		db.Close()
	})

	return repo
}
