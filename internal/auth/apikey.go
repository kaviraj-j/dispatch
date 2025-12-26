package auth

import "errors"

type Auth struct {
	producerApiKey string
	consumerApiKey string
}

type ClientType string

const (
	ClientTypeProducer ClientType = "PRODUCER"
	ClientTypeConsumer ClientType = "CONSUMER"
)

// Errors
var (
	ErrInvalidApiKey error = errors.New("invalid api key")
)

func NewAuth(producerApiKey string, consumerApiKey string) *Auth {
	return &Auth{
		producerApiKey: producerApiKey,
		consumerApiKey: consumerApiKey,
	}
}

func (auth *Auth) IsAuthenticated(apikey string, clientType ClientType) error {
	valid := false
	switch clientType {
	case ClientTypeProducer:
		if auth.producerApiKey == apikey {
			valid = true
		}
	case ClientTypeConsumer:
		if auth.consumerApiKey == apikey {
			valid = true
		}
	}
	if !valid {
		return ErrInvalidApiKey
	}
	return nil
}
