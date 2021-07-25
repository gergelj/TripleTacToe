package router

import (
	"net/http"

	httpHandler "gregvader/triple-tac-toe/delivery/http"
	"gregvader/triple-tac-toe/delivery/http/middleware"
	"gregvader/triple-tac-toe/domain"
	"gregvader/triple-tac-toe/websocket"

	"github.com/gorilla/mux"
)

type Router struct {
	Router *mux.Router
}

func (r *Router) Initialize( /*dbConn *database.DBConn*/ ) {
	r.Router.Use(middleware.CORSMiddleware)
	r.Router.Use(middleware.JSONMiddleware)
	//r.Router.Use(middleware.DBTransactionMiddleware(dbConn.DB))
	r.Router.HandleFunc("/api/test", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("Success")); w.WriteHeader(http.StatusOK) }).Methods("GET", "OPTIONS")

	//userRepo := postgres.NewPostgresUserRepository(dbConn)
	//authUsecase := usecase.NewAuthUsecase(userRepo)
	//authHandler := httpHandler.NewAuthHandler(authUsecase)

	/*authSubRouter := r.Router.PathPrefix("/api/auth").Subrouter()
	authSubRouter.Path("/login").Handler(http.HandlerFunc(authHandler.Login)).Methods("POST", "OPTIONS")
	authSubRouter.Path("/register").Handler(http.HandlerFunc(authHandler.Register)).Methods("POST", "OPTIONS")
	*/
	setupWebsocket(r)
}

func setupWebsocket(r *Router) {
	pool := websocket.NewPool()
	go pool.Start()
	websocketHandler := httpHandler.NewWebsocketHandler()
	r.Router.HandleFunc("/api/ws", func(w http.ResponseWriter, r *http.Request) {
		username := r.URL.Query().Get("username")
		if username == "" {
			http.Error(w, "please fill in your username", http.StatusBadRequest)
			return
		}

		websocketHandler.ServeWs(pool, domain.User{Username: username}, w, r)
	}).Methods("GET", "OPTIONS")
}
