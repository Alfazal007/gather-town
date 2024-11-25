package controllers

import (
	"fmt"
	"net/http"

	"github.com/Alfazal007/gather-town/helpers"
	"github.com/Alfazal007/gather-town/utils"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
)

func (apiCfg *ApiConf) IsValidUser(w http.ResponseWriter, r *http.Request) {
	usernameRequested := chi.URLParam(r, "username")
	jwtToken := chi.URLParam(r, "token")

	jwtSecret := utils.LoadEnvVariables().AccessTokenSecret
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		helpers.RespondWithError(w, 400, "Invalid token here", []string{})
		return
	}

	if !token.Valid {
		helpers.RespondWithError(w, 400, "Invalid token", []string{})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		helpers.RespondWithError(w, 400, "Invalid claims login again", []string{})
		return
	}

	username := claims["username"].(string)

	if username != usernameRequested {
		helpers.RespondWithError(w, 400, "Invalid username trying to be accessed", []string{})
		return
	}

	_, err = apiCfg.DB.GetUserByName(r.Context(), username)
	if err != nil {
		helpers.RespondWithError(w, 400, "Some manpulation done with the token", []string{})
		return
	}

	helpers.RespondWithJSON(w, 200, CanJoin{CanJoin: true})
}
