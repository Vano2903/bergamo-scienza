package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

var (
	enabled  bool = false
	accepted bool = false
	start    int64
	delay    int64
)

func mainPageHandler(w http.ResponseWriter, r *http.Request) {
	home, _ := os.ReadFile("pages/home.html")
	w.Write(home)
}

func adminPageHandler(w http.ResponseWriter, r *http.Request) {
	admin, _ := os.ReadFile("pages/admin.html")
	w.Write(admin)
}

func startHandler(w http.ResponseWriter, r *http.Request) {
	start = time.Now().UnixNano() / int64(time.Millisecond)
	enabled = true
}

func richiestaHandler(w http.ResponseWriter, r *http.Request) {
	if enabled {
		if !accepted {
			accepted = true
			end := time.Now().UnixNano() / int64(time.Millisecond)
			delay = end - start
			w.Write([]byte(fmt.Sprintf(`{"ok": true, "message":"hai premuto dopo %d millisecondi"}`, delay)))
			return
		}
		w.Write([]byte(fmt.Sprintf(`{"ok": false, "message":"il bottone é stato premuto dopo %d millisecondi"}`, delay)))
		return
	}
	w.Write([]byte(`{"msg": "troppo presto"}`))
}

func resetHandler(w http.ResponseWriter, r *http.Request) {
	accepted = false
	enabled = false
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := mux.NewRouter()
	r.HandleFunc("/", mainPageHandler).Methods("get")
	r.HandleFunc("/admin", adminPageHandler).Methods("get")
	r.HandleFunc("/start", startHandler).Methods("get")
	r.HandleFunc("/reset", resetHandler).Methods("get")
	r.HandleFunc("/richiesta", richiestaHandler).Methods("get")

	log.Println("starting on:" + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
