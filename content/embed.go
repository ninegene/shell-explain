// Package content provides the original Markdown sources embedded in the CLI.
package content

import "embed"

// Files is the single authored explanation collection: command pages in en,
// beginner guides in en/guides, and the glossary in en/glossary. Subdirectories
// must be listed here explicitly; the embed directive does not descend into them.
//
//go:embed en/*.md en/guides/*.md en/glossary/*.md
var Files embed.FS
