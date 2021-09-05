package config

type TomlConfig struct {
	Title    string
	Settings Settings
	DB       Database
	Log      LogConfig `toml:"Logger"`
}

type Settings struct {
	Host      string
	HTTPSPort string
}

type Database struct {
	Name     string
	Host     string
	User     string
	Password string
	Port     string
}

type LogConfig struct {
	Output string
	Level  string
}
