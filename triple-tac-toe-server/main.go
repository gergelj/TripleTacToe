package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"gregvader/triple-tac-toe/globals"
	routing "gregvader/triple-tac-toe/router"

	"github.com/gorilla/mux"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	//dbConn := database.DBConn{}
	//dbConn.ConnectToDatabase()

	/*if globals.IsDevEnvironment() {
		migrations.Migrate(dbConn)
		seed.Run(&dbConn)
	}*/

	server := &http.Server{
		Addr: ":" + globals.Port,
		//Handler: getPreparedRouter(&dbConn).Router,
		Handler: getPreparedRouter().Router,
	}
	server.ListenAndServe()
	log.Println("[TRIPLE-TAC-TOE-SERVER] Listening on " + globals.Port)
}

/*
func getPreparedRouter(conn *database.DBConn) routing.Router {
	router := routing.Router{Router: &mux.Router{}}
	router.Initialize(conn)
	return router
}*/

func getPreparedRouter() routing.Router {
	router := routing.Router{Router: &mux.Router{}}
	router.Initialize()
	return router
}
