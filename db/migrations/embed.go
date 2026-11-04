package migrations

import "embed"

// Files embeds the versioned SQL migration files used by the API image.
//
//go:embed *.sql
var Files embed.FS
