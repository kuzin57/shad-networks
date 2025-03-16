package config

type Config struct {
	App      *AppConfig      `yaml:"app"`
	Neo4j    *Neo4jConfig    `yaml:"neo4j"`
	Redis    *RedisConfig    `yaml:"redis"`
	Postgres *PostgresConfig `yaml:"postgres"`
	Kafka    *KafkaConfig    `yaml:"kafka"`
}

type AppConfig struct {
	Port int `yaml:"port"`
}

type Neo4jConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type RedisConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type PostgresConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type KafkaConfig struct {
	Host   string             `yaml:"host"`
	Port   int                `yaml:"port"`
	Topics []KafkaTopicConfig `yaml:"topics"`
}

type KafkaTopicConfig struct {
	Topic      string `yaml:"topic"`
	Partitions int32  `yaml:"partitions"`
}
