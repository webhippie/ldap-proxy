package config

type server struct {
	Addr       string
	Cert       string
	Key        string
	Templates  string
	Assets     string
	Endpoint   string
	Title      string
	Pprof      bool
	Prometheus bool
}

type ldap struct {
	Addr         string
	BindUsername string
	BindPassword string
	BaseDN       string
	FilterDN     string
	UserAttr     string
	UserHeader   string
	MailAttr     string
	MailHeader   string
}

var (
	// Debug represents the flag to enable or disable debug logging.
	Debug bool

	// Server represents the information about the server bindings.
	Server = &server{}

	// LDAP represents the information about the ldap server bindings.
	LDAP = &ldap{}
)
