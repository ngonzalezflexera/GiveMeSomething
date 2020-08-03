package middleware

import (
	"auth/model"
	"context"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
)

func ValidateToken(tokeString string) (*model.Token, error) {
	token := &model.Token{}
	_, err := jwt.ParseWithClaims(tokeString, token, func(t *jwt.Token) (interface{}, error) {
		return []byte("secretword"), nil
	})
	return token, err
}
func JwsVerification(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("x-access-token")
		header = strings.TrimSpace(header)
		if header == "" {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"message": "Missing Auth token"})
			return
		}

		token, err := ValidateToken(header)

		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
			return
		}
		ctx := context.WithValue(r.Context(), "user", token)
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}
