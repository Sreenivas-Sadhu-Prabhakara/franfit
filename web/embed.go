// Package web embeds the self-contained frontend (no CDNs, no external
// assets) so the whole product ships as one binary.
package web

import "embed"

// Files holds the static UI served at /.
//
//go:embed index.html style.css app.js
var Files embed.FS
