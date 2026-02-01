package www

import "embed"

// DistFS contains the embedded frontend assets built to dist/.
// The embed directive requires the dist/ directory to exist at build time.
// The "all:" prefix includes files in subdirectories recursively.
//
//go:embed all:dist
var DistFS embed.FS
