package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/websocket"
)

// html to websocket upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// reouter url handler
type spaHandler struct {
	staticPath string
	indexPath  string
}

func corsConf(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

//spa must implement ServeHTTP
func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		Debugf("[ERROR] absolute request path %s -> %s\n", r.Method, path)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	path = filepath.Join(h.staticPath, path)
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		Debugf("[OK FALL BACK NOT EXIST] absolute request path: %s -> %s\n", r.Method, path)
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		Debugf("[ERROR] absolute request path: %s -> %s\n", r.Method, path)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Debugf("[OK] joined path %s -> %s\n", r.Method, path)
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		Debugf("readercall: {\n\t%s\n}", string(p))
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

func serveWebsocket(w http.ResponseWriter, r *http.Request) {
	Debugf("serveWebsocket call\n")
	Debugf("%s\n", r.URL.Host)
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("An error inside the serveWebsocket occurred: \n\t", err)
	}
	reader(ws)
}
