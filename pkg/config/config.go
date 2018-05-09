package config

// Server defines the server configuration.
type Server struct {
	Health        string
	Secure        string
	Public        string
	Host          string
	Root          string
	Cert          string
	Key           string
	AutoCert      bool
	StrictCurves  bool
	StrictCiphers bool
	Templates     string
	Assets        string
	Storage       string
}

// Logs defines the logging configuration.
type Logs struct {
	Level   string
	Colored bool
	Pretty  bool
}

// Proxy defines the proxy configuration.
type Proxy struct {
	Title      string
	Endpoints  []string
	UserHeader string
}

// LDAP defines the ldap configuration.
type LDAP struct {
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

// Config defines the general configuration.
type Config struct {
	Server Server
	Logs   Logs
	Proxy  Proxy
	LDAP   LDAP
}

// New prepares a new default configuration.
func New() *Config {
	return &Config{}
}
