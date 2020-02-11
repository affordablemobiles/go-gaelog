package main

import (
	"log"
	"net/http"
	"os"

	glog "github.com/a1comms/go-gaelog"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// Just give a 403 to nosey people trying to hit the index.
	r.NotFoundHandler = r.NewRoute().BuildOnly().HandlerFunc(defaultHandler).GetHandler()

	// Handle all HTTP requests with our router.
	http.Handle("/", r)

	// [START setting_port]
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
	// [END setting_port]
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	glog.Printf(r, nil, "I'm logging, %s", "wuhoo!")

	glog.Errorf(r, map[string]interface{}{
		"code":    403,
		"message": "Permission Denied",
	}, "HTTP ERROR")

	http.Error(w, "Permission Denied", 403)
}
