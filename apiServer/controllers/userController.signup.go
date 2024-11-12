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
	"github.com/Alfazal007/gather-town/types"
	"github.com/Alfazal007/gather-town/validators"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func (apiCfg *ApiConf) Signup(w http.ResponseWriter, r *http.Request) {
	validate := validator.New()
	var signUpParams validators.SignupValidators

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&signUpParams)
	if err != nil {
		helpers.RespondWithError(w, 400, fmt.Sprintf("Invalid json object provided %v", err.Error()), []string{})
		return
	}

	signUpParams.Username = strings.TrimSpace(signUpParams.Username)
	signUpParams.Email = strings.TrimSpace(signUpParams.Email)
	signUpParams.Password = strings.TrimSpace(signUpParams.Password)

	err = validate.Struct(signUpParams)
	if err != nil {
		var errMessageArray []string
		for _, fieldError := range err.(validator.ValidationErrors) {
			errMessageArray = append(errMessageArray, fmt.Sprintf("Validation failed on %v at the tag %v", fieldError.Field(), fieldError.Tag()))
		}
		helpers.RespondWithError(w, 400, "Validation errors", errMessageArray)
		return
	}
	// validation is done now interact with the database
	existingUser, err := apiCfg.DB.FindUsernameOrEmail(r.Context(), database.FindUsernameOrEmailParams{
		Username: signUpParams.Username,
		Email:    signUpParams.Email,
	})

	if err != nil && err != sql.ErrNoRows {
		helpers.RespondWithError(w, 400, "Issue talking to the database", []string{})
		return
	}
	if existingUser.Username == signUpParams.Username {
		helpers.RespondWithError(w, 400, "User with this username exists", []string{})
		return
	}
	if existingUser.Email == signUpParams.Email {
		helpers.RespondWithError(w, 400, "User with this email exists", []string{})
		return
	}
	createdUser, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		Username:  signUpParams.Username,
		Email:     signUpParams.Email,
		Password:  signUpParams.Password, // TODO:: need to hash this password
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		helpers.RespondWithError(w, 400, "Issue creating the new user", []string{})
		return
	}
	helpers.RespondWithJSON(w, 200, types.ReturnCreatedUser(createdUser))
}
