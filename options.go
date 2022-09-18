package unisrv

// Options describes options for unisrv handler.
type Options struct {
	// Base specifies the base path for the Unity application.
	Base string
	// NoCache specifies whether to set `Cache-Control: no-cache` header.
	NoCache bool
}
