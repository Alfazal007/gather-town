package router

import (
	"github.com/Alfazal007/gather-town/controllers"
	"github.com/go-chi/chi/v5"
)

func UserRouter(apiCfg *controllers.ApiConf) *chi.Mux {
	r := chi.NewRouter()
	r.Post("/sign-up", apiCfg.Signup)
	return r
}
