package config

type Listener struct {
	Hostname string
	Port     string
}

type SslSettings struct {
	Enabled bool
}

type TlsManager struct {
	Email         string
	HostWhitelist string // comma separated list of hostnames
	CacheDir      string
}

type HttpConfig struct {
	Listener   Listener
	TlsManager TlsManager
	Ssl        SslSettings
}

type StorableFileConfig struct {
	Path string
}
