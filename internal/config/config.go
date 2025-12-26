package config

type Config struct {
	ServerAddress  string
	ProducerApiKey string
	ConsumerApiKey string
}

func Load() Config {
	return Config{
		ServerAddress:  ":8080",
		ProducerApiKey: "producer-123",
		ConsumerApiKey: "consumer-123",
	}
}
