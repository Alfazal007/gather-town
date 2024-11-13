package router

import (
	"net/http"

	"github.com/Alfazal007/gather-town/controllers"
	"github.com/go-chi/chi/v5"
)

func RoomRouter(apiCfg *controllers.ApiConf) *chi.Mux {
	r := chi.NewRouter()
	r.Post("/create-room", controllers.VerifyJWT(apiCfg, http.HandlerFunc(apiCfg.CreateRoom)).ServeHTTP)
	r.Post("/add-member", controllers.VerifyJWT(apiCfg, http.HandlerFunc(apiCfg.AddMembersToRoom)).ServeHTTP)
	return r
}
