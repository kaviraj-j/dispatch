package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerAddress  string
	ProducerApiKey string
	ConsumerApiKey string
}

func Load() Config {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	return Config{
		ServerAddress:  os.Getenv("SERVER_ADDRESS"),
		ProducerApiKey: os.Getenv("PRODUCER_API_KEY"),
		ConsumerApiKey: os.Getenv("CONSUMER_API_KEY"),
	}
}
