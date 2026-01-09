package middleware

import (
	"context"
	"net/http"

	"github.com/kaviraj-j/dispatch/internal/auth"
)

type Middleware struct {
	auth *auth.Auth
}

func NewMiddleware(auth *auth.Auth) *Middleware {
	return &Middleware{
		auth: auth,
	}
}

type contextKey string

const (
	consumerIDKey contextKey = "consumerID"
	producerIDKey contextKey = "producerID"
)

func (m *Middleware) IsProducerAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" {
			http.Error(w, "missing api key", http.StatusUnauthorized)
			return
		}

		err := m.auth.IsAuthenticated(apiKey, auth.ClientTypeProducer)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), producerIDKey, apiKey)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) IsConsumerAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" {
			http.Error(w, "missing api key", http.StatusUnauthorized)
			return
		}

		err := m.auth.IsAuthenticated(apiKey, auth.ClientTypeConsumer)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), consumerIDKey, apiKey)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func ConsumerIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(consumerIDKey).(string)
	return id, ok
}

func ProducerIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(producerIDKey).(string)
	return id, ok
}
