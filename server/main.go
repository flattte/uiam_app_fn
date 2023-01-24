package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type sessionResolutionPoints struct {
	store *sessions.CookieStore
	users map[string]string
	times TimesStore
}

func init_router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})
	router.HandleFunc("/logout", srp.logoutHandlerCA).Methods("GET")
	router.HandleFunc("/healthcheck", srp.healtcheckCA).Methods("GET")

	router.HandleFunc("/login_req", srp.loginHandlerCA).Methods("POST", "OPTIONS")
	//login_router := router.Path("/login_req").Subrouter()
	//login_router.Methods("POST").HandlerFunc(srp.logoutHandlerCA)
	//login_router.Methods("OPTIONS").HandlerFunc(corsConf)

	router.HandleFunc("/ws", serveWebsocket)
	router.HandleFunc("/timer_collect", srp.timerCollectCA).Methods("POST", "OPTIONS")
	router.HandleFunc("/registration_req", srp.registerUserCA).Methods("POST", "OPTIONS")

	spa := spaHandler{
		staticPath: "../Rubik-Cube/build",
		indexPath:  "/index.html",
	}
	router.PathPrefix("/").Handler(spa)
	return router
}

var srp = &sessionResolutionPoints{
	store: sessions.NewCookieStore([]byte("my_secret_key")),
	users: map[string]string{"user1@lmao.com": "password", "user2@lmao.com": "password2"},
}

func main() {
	// origins := handlers.AllowedOrigins([]string{"*"})
	// headers := handlers.AllowedHeaders([]string{"*"})
	// creds := handlers.AllowCredentials()

	router := init_router()
	router.Use(mux.CORSMethodMiddleware(router))

	listen_on := os.Args[1] + ":3000"

	server := http.Server{
		Handler:      router,
		Addr:         listen_on,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	Debugf("Server running on %s\n", server.Addr)
	log.Fatal(server.ListenAndServeTLS("./ssl/cert.pem", "./ssl/key.pem"))
}
