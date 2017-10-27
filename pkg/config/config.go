package config

type server struct {
	Host          string
	Addr          string
	Cert          string
	Key           string
	Root          string
	Storage       string
	Templates     string
	Assets        string
	Endpoint      string
	Title         string
	LetsEncrypt   bool
	StrictCurves  bool
	StrictCiphers bool
	Prometheus    bool
	Pprof         bool
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
	// LogLevel defines the log level used by our logging package.
	LogLevel string

	// Server represents the information about the server bindings.
	Server = &server{}

	// LDAP represents the information about the ldap server bindings.
	LDAP = &ldap{}
)
