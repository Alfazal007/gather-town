package controllers

import (
	"net/http"

	"github.com/Alfazal007/gather-town/helpers"
	"github.com/Alfazal007/gather-town/internal/database"
)

func (apiCfg *ApiConf) DeleteUser(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(database.User)
	if !ok {
		helpers.RespondWithError(w, 400, "Issue with finding the user from the database", []string{})
		return
	}
	_, err := apiCfg.DB.DeleteUserViaId(r.Context(), user.ID)
	if err != nil {
		helpers.RespondWithError(w, 400, "Failed to delete the user", []string{})
		return
	}
	helpers.RespondWithJSON(w, 200, map[string]interface{}{})
}
