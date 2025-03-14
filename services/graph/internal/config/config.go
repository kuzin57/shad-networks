package config

type Config struct {
	App      *AppConfig      `yaml:"app"`
	Neo4j    *Neo4jConfig    `yaml:"neo4j"`
	Redis    *RedisConfig    `yaml:"redis"`
	Postgres *PostgresConfig `yaml:"postgres"`
}

type AppConfig struct {
	Port int `yaml:"port"`
}

type Neo4jConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type RedisConfig struct {
	Port int `yaml:"port"`
}

type PostgresConfig struct {
	Port int `yaml:"port"`
}
