package seed

import "embed"

// Files embeds fictional, development-only seed data.
//
//go:embed *.sql
var Files embed.FS
