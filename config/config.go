package config

type server struct {
	Addr      string
	Cert      string
	Key       string
	Templates string
	Assets    string
	Pprof     bool
}

var (
	// Debug represents the flag to enable or disable debug logging.
	Debug bool

	// Server represents the informations about the server bindings.
	Server = &server{}
)
