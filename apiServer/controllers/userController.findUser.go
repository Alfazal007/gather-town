package controllers

import (
	"net/http"

	"github.com/Alfazal007/gather-town/helpers"
	"github.com/Alfazal007/gather-town/internal/database"
	"github.com/Alfazal007/gather-town/types"
	"github.com/go-chi/chi/v5"
)

func (apiCfg *ApiConf) FindCurrentUser(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value("user").(database.User)
	if !ok {
		helpers.RespondWithError(w, 400, "Issue with finding the user from the database", []string{})
		return
	}
	username := chi.URLParam(r, "username")
	userToBeFound, err := apiCfg.DB.GetUserByName(r.Context(), username)
	if err != nil {
		helpers.RespondWithError(w, 400, "Issue finding the user", []string{})
		return
	}
	helpers.RespondWithJSON(w, 200, types.ReturnCreatedUser(userToBeFound))
}
