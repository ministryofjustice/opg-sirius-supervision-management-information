package api

import (
	"github.com/opg-sirius-supervision-management-information/management-information-api/internal/auth"
	"github.com/opg-sirius-supervision-management-information/shared"
	"net/http"
	"strconv"
	"strings"
)

func (s *Server) authenticateAPI(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := auth.NewContext(r)
		logger := s.Logger(ctx)

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			logger.Error("Unable to authorise user token: ", "err", "missing bearer token")
			http.Error(w, "missing bearer token", http.StatusUnauthorized)
			return
		}

		requestToken := strings.Split(authHeader, "Bearer ")[1]
		token, err := s.JWT.Verify(requestToken)

		if err != nil {
			logger.Error("Unable to authorise user token: ", "err", err.Error())
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		claims := token.Claims.(*auth.Claims)
		userID, _ := strconv.ParseInt(claims.ID, 10, 32)

		ctx.User = &shared.User{
			ID:    int32(userID),
			Roles: claims.Roles,
		}

		h.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (s *Server) authorise(role string) func(http.Handler) http.HandlerFunc {
	return func(h http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context().(auth.Context)

			user, err := s.GetCurrentUserDetails(ctx)
			if err != nil {
				s.Logger(ctx).Error("Unable to authorise user: ", "err", err.Error())
				http.Error(w, "Error", http.StatusInternalServerError)
				return
			}

			ctx.User = &user

			if !user.HasRole(role) {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			h.ServeHTTP(w, r)
		}
	}
}
