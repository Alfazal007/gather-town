package router

import (
	"net/http"

	"github.com/Alfazal007/gather-town/controllers"
	"github.com/go-chi/chi/v5"
)

func UserRouter(apiCfg *controllers.ApiConf) *chi.Mux {
	r := chi.NewRouter()
	r.Post("/sign-up", apiCfg.Signup)
	r.Post("/sign-in", apiCfg.SignIn)
	r.Get("/current-user", controllers.VerifyJWT(apiCfg, http.HandlerFunc(apiCfg.GetCurrentUser)).ServeHTTP)
	r.Delete("/delete-user", controllers.VerifyJWT(apiCfg, http.HandlerFunc(apiCfg.DeleteUser)).ServeHTTP)
	r.Get("/username/{username}", controllers.VerifyJWT(apiCfg, http.HandlerFunc(apiCfg.FindCurrentUser)).ServeHTTP)
	r.Get("/get-rooms", controllers.VerifyJWT(apiCfg, http.HandlerFunc(apiCfg.GetRoomsUserIsPartOf)).ServeHTTP)
	r.Get("/get-admin/roomId/{roomId}", controllers.VerifyJWT(apiCfg, http.HandlerFunc(apiCfg.GetAdminInformation)).ServeHTTP)
	r.Get("/token/{token}/username/{username}", apiCfg.IsValidUser)
	return r
}
