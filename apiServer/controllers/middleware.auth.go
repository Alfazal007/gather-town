package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Alfazal007/gather-town/helpers"
	"github.com/Alfazal007/gather-town/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func VerifyJWT(apiCfg *ApiConf, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var jwtToken string
		authorization := r.Header.Get("Authorization")
		if authorization == "" || !strings.HasPrefix(authorization, "Bearer ") {
			helpers.RespondWithError(w, 400, "No headers provided", []string{})
			return
		}

		jwtToken = strings.TrimPrefix(authorization, "Bearer ")
		// Verify the JWT token
		jwtSecret := utils.LoadEnvVariables().AccessTokenSecret

		if jwtToken == "" {
			helpers.RespondWithError(w, 400, "Provide access token", []string{})
			return
		}
		token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})
		if err != nil {
			helpers.RespondWithError(w, 401, "Invalid token here", []string{})
			return
		}
		if !token.Valid {
			helpers.RespondWithError(w, 401, "Invalid token", []string{})
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			helpers.RespondWithError(w, 400, "Invalid claims login again", []string{})
			return
		}

		username := claims["username"].(string)
		id := claims["user_id"].(string)

		user, err := apiCfg.DB.GetUserByName(r.Context(), username)
		if err != nil {
			helpers.RespondWithError(w, 400, "Some manpulation done with the token", []string{})
			return
		}
		idUUID, err := uuid.Parse(id)
		if err != nil {
			helpers.RespondWithError(w, 400, "Some manpulation done with the token", []string{})
			return
		}
		if idUUID != user.ID {
			helpers.RespondWithError(w, 400, "Some manipulations done with the token try again", []string{})
			return
		}
		ctx := context.WithValue(r.Context(), "user", user)
		r = r.WithContext(ctx)
		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
