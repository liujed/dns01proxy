package gen

type Build struct {
	// Identifies the GitHub project for the build's Caddy DNS package.
	ProjectURL string `json:"project_url"`

	// The Go module path for the build's Caddy DNS package.
	GoModPath string `json:"go_mod_path"`

	// The Caddy DNS package's version.
	GoModVersion string `json:"go_mod_version"`

	// Points to the Caddy documentation page for the build's DNS module.
	CaddyDocURL string `json:"caddy_doc_url"`
}
