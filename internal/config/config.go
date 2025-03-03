package config

type Config struct {
	App *AppConfig `yaml:"app"`
	DB  *DBConfig  `yaml:"db"`
}

type AppConfig struct {
	Port int `yaml:"port"`
}

type DBConfig struct {
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}
