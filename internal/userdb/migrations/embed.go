package migrations

import "embed"

// FS contains user data database migrations.
//
//go:embed *.sql
var FS embed.FS
