package router

import (
	"net/http"

	"github.com/Alfazal007/gather-town/controllers"
	"github.com/go-chi/chi/v5"
)

func RoomRouter(apiCfg *controllers.ApiConf) *chi.Mux {
	r := chi.NewRouter()
	r.Post("/create-room", controllers.VerifyJWT(apiCfg, http.HandlerFunc(apiCfg.CreateRoom)).ServeHTTP)
	r.Delete("/delete-room", controllers.VerifyJWT(apiCfg, http.HandlerFunc(apiCfg.DeleteRoom)).ServeHTTP)
	r.Post("/add-member", controllers.VerifyJWT(apiCfg, http.HandlerFunc(apiCfg.AddMembersToRoom)).ServeHTTP)
	r.Put("/remove-member", controllers.VerifyJWT(apiCfg, http.HandlerFunc(apiCfg.RemoveMemberFromTheRoomByAdmin)).ServeHTTP)
	r.Put("/leave-member", controllers.VerifyJWT(apiCfg, http.HandlerFunc(apiCfg.LeaveTheRoom)).ServeHTTP)
	r.Get("/roomId/{roomId}", controllers.VerifyJWT(apiCfg, http.HandlerFunc(apiCfg.GetRoomInfo)).ServeHTTP)
	r.Get("/join-room/{roomId}", controllers.VerifyJWT(apiCfg, http.HandlerFunc(apiCfg.UserCanJoinRoom)).ServeHTTP)
	r.Get("/join-room/{roomId}/token/{token}/username/{username}", apiCfg.UserCanJoinRoomApiCall)
	return r
}
