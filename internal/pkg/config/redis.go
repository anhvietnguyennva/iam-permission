package config

// RedisConfig contains all configuration of redis
type RedisConfig struct {
	Addresses  string
	Password   string
	MasterName string
}
