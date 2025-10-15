package config

type Config struct {
	HTTPPort  string `mapstructure:"HTTP_PORT"`
	LogLevel  string `mapstructure:"LOG_LEVEL"`
	JWTSecret string `mapstructure:"JWT_SECRET"`

	DB DBConfig `mapstructure:",squash"`
}

type DBConfig struct {
	DSN          string `mapstructure:"DB_DSN"`
	MaxOpenConns int    `mapstructure:"DB_MAX_OPEN_CONNS"`
	MaxIdleConns int    `mapstructure:"DB_MAX_IDLE_CONNS"`
}
