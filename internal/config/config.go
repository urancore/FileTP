package config

type Config struct {
	RootPath string `json:"root_path"`
}

func NewConfig(root_path string) *Config {
	return &Config{
		RootPath: root_path,
	}
}
