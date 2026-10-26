// Package content provides the original Markdown sources embedded in the CLI.
package content

import "embed"

// Files is the single authored explanation collection.
//
//go:embed en/*.md
var Files embed.FS
