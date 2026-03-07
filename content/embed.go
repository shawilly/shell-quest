// Package content provides embedded game content files.
package content

import "embed"

// WorldsFS holds all embedded world JSON files.
//go:embed worlds
var WorldsFS embed.FS
