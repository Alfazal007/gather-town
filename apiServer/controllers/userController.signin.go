package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Alfazal007/gather-town/helpers"
	"github.com/Alfazal007/gather-town/internal/database"
	"github.com/Alfazal007/gather-town/utils"
	"github.com/Alfazal007/gather-town/validators"
	"github.com/go-playground/validator/v10"
)

func (apiCfg *ApiConf) SignIn(w http.ResponseWriter, r *http.Request) {
	validate := validator.New()
	var signInParams validators.SigninValidators

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&signInParams)
	if err != nil {
		helpers.RespondWithError(w, 400, fmt.Sprintf("Invalid json object provided %v", err.Error()), []string{})
		return
	}

	signInParams.Username = strings.TrimSpace(signInParams.Username)
	signInParams.Password = strings.TrimSpace(signInParams.Password)

	err = validate.Struct(signInParams)
	if err != nil {
		var errMessageArray []string
		for _, fieldError := range err.(validator.ValidationErrors) {
			errMessageArray = append(errMessageArray, fmt.Sprintf("Validation failed on %v at the tag %v", fieldError.Field(), fieldError.Tag()))
		}
		helpers.RespondWithError(w, 400, "Validation errors", errMessageArray)
		return
	}
	existingUser, err := apiCfg.DB.FindUsernameOrEmailForLogin(r.Context(), signInParams.Username)
	if err != nil {
		helpers.RespondWithError(w, 404, "Error logging in, recheck the credentials or try again later", []string{})
		return
	}
	isPasswordCorrect := utils.ValidatePassword(signInParams.Password, existingUser.Password)
	if !isPasswordCorrect {
		helpers.RespondWithError(w, 400, "Incorrect password", []string{})
		return
	}
	// generate access and refresh tokens
	accessToken, refreshToken, err := utils.GenerateTokens(existingUser)
	if err != nil {
		helpers.RespondWithError(w, 400, "Issue generating the tokens", []string{})
		return
	}
	_, err = apiCfg.DB.UpdateRefreshToken(r.Context(), database.UpdateRefreshTokenParams{
		RefreshToken: sql.NullString{
			String: refreshToken,
			Valid:  true,
		},
		Username: existingUser.Username,
	})
	if err != nil {
		helpers.RespondWithError(w, 400, "Issue updating the database", []string{})
		return
	}
	// add refresh token to the database
	cookie1 := http.Cookie{
		Name:     "accessToken",
		Value:    accessToken,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		Expires:  time.Now().Add(time.Hour * 24),
	}
	http.SetCookie(w, &cookie1)
	cookie2 := http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
		Expires:  time.Now().Add(time.Hour * 240),
	}
	http.SetCookie(w, &cookie2)
	cookie3 := http.Cookie{
		Name:     "username",
		Value:    existingUser.Username,
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
		Expires:  time.Now().Add(time.Hour * 240),
	}
	http.SetCookie(w, &cookie3)
	cookie4 := http.Cookie{
		Name:     "id",
		Value:    existingUser.ID.String(),
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
		Expires:  time.Now().Add(time.Hour * 240),
	}
	http.SetCookie(w, &cookie4)
	type Tokens struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
		Username     string `json:"username"`
		Id           string `json:"userid"`
	}
	helpers.RespondWithJSON(w, 200, Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Username:     existingUser.Username,
		Id:           existingUser.ID.String(),
	})
}
