package middleware

import (
	"net/http"
	"strings"

	"github.com/nemolucky/patisserie-go/internal/ctxutil"
	"github.com/nemolucky/patisserie-go/internal/domain"
	tokenprovider "github.com/nemolucky/patisserie-go/internal/token_provider"
)

func AuthMiddleware(tokenProvider *tokenprovider.TokenProvider) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")

			if header == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(header, " ", 2)

			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			claims, err := tokenProvider.Parse(
				parts[1],
				domain.AccessTokenType,
			)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := ctxutil.SetUserID(r.Context(), claims.UserID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}