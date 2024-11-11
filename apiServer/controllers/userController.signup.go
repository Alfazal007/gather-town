package controllers

import (
	"net/http"

	"github.com/Alfazal007/gather-town/helpers"
)

func (apiCfg *ApiConf) Signup(w http.ResponseWriter, r *http.Request) {
	helpers.RespondWithJSON(w, 200, "message")
}
