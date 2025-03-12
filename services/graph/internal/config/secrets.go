package config

type Secrets struct {
	Neo4j    *Neo4jSecret    `yaml:"neo4j"`
	Redis    *RedisSecret    `yaml:"redis"`
	Postgres *PostgresSecret `yaml:"postgres"`
}

type Neo4jSecret struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type RedisSecret struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type PostgresSecret struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}
