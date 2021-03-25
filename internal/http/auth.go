package http

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"home24-technical-test/pkg/appcontext"
	"home24-technical-test/pkg/http/response"

	userAdapter "home24-technical-test/internal/user/adapter"
)

func (hs *Server) authorizedOnly(getUserAdapter userAdapter.GetUserAdapter, getLoginSessionAdapter userAdapter.GetLoginSessionAdapter) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			var userID int

			ctx := r.Context()
			session := getSessionToken(r)
			if session == "" {
				response.Error(w, "Forbidden", http.StatusForbidden, fmt.Errorf("Access denied"))
				return
			} else {
				userSession, err := getLoginSessionAdapter.Execute(ctx, session)
				if err != nil {
					response.Error(w, "Internal Server Error", http.StatusInternalServerError, err)
					return
				}
				if userSession == nil {
					response.Error(w, "Unauthorized", http.StatusUnauthorized, fmt.Errorf("Unauthorized"))
					return
				}

				userData, err := getUserAdapter.Execute(ctx, int(userSession.Info["UserID"].(float64)))
				if err != nil {
					response.Error(w, "Internal Server Error", http.StatusInternalServerError, err)
					return
				}
				userID = userData.ID
			}
			ctx = context.WithValue(ctx, appcontext.KeySessionID, session)

			if userID != 0 {
				ctx = context.WithValue(ctx, appcontext.KeyUserID, userID)
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}

func getSessionToken(r *http.Request) string {
	token := r.Header.Get("Authorization")
	splitToken := strings.Split(token, "session")

	if len(splitToken) < 2 {
		return ""
	}

	token = strings.Trim(splitToken[1], " ")
	return token
}
