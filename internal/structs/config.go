package structs

type Config struct {
	Apps    AppsConfig    `yaml:"apps"`
	Storage StorageConfig `yaml:"storage"`
	OSM     OSMConfig     `yaml:"osm"`
}

type AppsConfig struct {
	Api ApiConfig `yaml:"api"`
}

type ApiConfig struct {
	Name     string `yaml:"name"`
	Version  string `yaml:"version"`
	Listener struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"listener"`
	Ssl struct {
		Enabled bool   `yaml:"enabled"`
		Cert    string `yaml:"cert"`
		Key     string `yaml:"key"`
		Manager struct {
			Email         string   `yaml:"email"`
			HostWhitelist []string `yaml:"host_whitelist"`
			CacheDir      string   `yaml:"cache_dir"`
		} `yaml:"manager"`
	} `yaml:"ssl"`
}

type StorageConfig struct {
	Root string `yaml:"root"`
	PBFs string `yaml:"pbfs"`
}

type OSMConfig struct {
	Sources struct {
		PBF struct {
			Url      string `yaml:"url"`
			FileName string `yaml:"file"`
		} `yaml:"pbf"`
	} `yaml:"sources"`
}
