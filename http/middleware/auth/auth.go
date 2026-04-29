//go:build !solution

package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

type User struct {
	Name  string
	Email string
}

func ContextUser(ctx context.Context) (*User, bool) {
	user, ok := ctx.Value(struct{}{}).(*User)
	return user, ok
}

var ErrInvalidToken = errors.New("invalid token")

type TokenChecker interface {
	CheckToken(ctx context.Context, token string) (*User, error)
}

func CheckAuth(checker TokenChecker) func(next http.Handler) http.Handler {
	f := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if !strings.HasPrefix(header, "Bearer ") {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			token := strings.Split(header, " ")[1]

			user, err := checker.CheckToken(r.Context(), token)

			if err == ErrInvalidToken {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			ctx := r.Context()
			newctx := context.WithValue(ctx, struct{}{}, user)

			newreq := r.WithContext(newctx)

			next.ServeHTTP(w, newreq)

		})
	}
	return f
}
