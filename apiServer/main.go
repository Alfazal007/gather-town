package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/Alfazal007/gather-town/controllers"
	"github.com/Alfazal007/gather-town/internal/database"
	router "github.com/Alfazal007/gather-town/routes"
	"github.com/Alfazal007/gather-town/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"
)

func main() {
	// LOAD ALL THE ENV VARIABLES
	envVariables := utils.LoadEnvVariables()
	// SETUP A DATABASE CONNECTION TO THE QUERY OBJECT
	conn, err := sql.Open("postgres", envVariables.DatabaseUrl)
	if err != nil {
		log.Fatal("Issue connecting to the database", err)
	}

	apiCfg := controllers.ApiConf{DB: database.New(conn)}

	// SETUP A ROUTER TO HANDLE REQUESTS
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Use(middleware.Logger)

	r.Mount("/api/v1/user", router.UserRouter(&apiCfg))
	r.Mount("/api/v1/room", router.RoomRouter(&apiCfg))
	log.Println("Starting the server at port", envVariables.Port)
	//	err = http.ListenAndServe(fmt.Sprintf(":%v", envVariables.Port), r)
	err = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%v", envVariables.Port), r)
	if err != nil {
		log.Fatal("There was an error starting the server", err)
	}
}
