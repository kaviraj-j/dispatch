package config

type Config struct {
	ServerAddress string
}

func Load() Config {
	return Config{
		ServerAddress: ":8080",
	}
}
