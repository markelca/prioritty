package migrations

import _ "embed"

//go:embed sql/schema.sql
var SchemaSQL string

//go:embed sql/seed.sql
var SeedSQL string
