package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/oxodao/caisseoverflow/routes"
	"github.com/oxodao/caisseoverflow/services"
	"github.com/rs/cors"

	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("CaisseOverflow")

	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found! Using currently exported vars")
	}

	dsn, exists := os.LookupEnv("POSTGRES_DSN")
	if !exists {
		panic("NO DATABASE")
	}

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		fmt.Println(err)
		panic("cant connect db")
	}

	provider := services.New(db)

	r := mux.NewRouter()

	r.HandleFunc("/items", routes.ListItemsRoute(provider)).Methods("GET")
	r.HandleFunc("/order", routes.AddOrderRoute(provider)).Methods("POST")
	r.HandleFunc("/order", routes.ListOrderRoute(provider)).Methods("GET")

	spa := spaHandler{staticPath: "frontend", indexPath: "index.html"}
	r.PathPrefix("/").Handler(spa)

	port, exists := os.LookupEnv("WEB_PORT")
	if !exists {
		panic("No listening port found! (WEB_PORT)")
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
	})

	srv := &http.Server{
		Handler:      c.Handler(r),
		Addr:         "127.0.0.1:" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

/*****

	https://github.com/gorilla/mux#serving-single-page-applications

*****/

// spaHandler implements the http.Handler interface, so we can use it
// to respond to HTTP requests. The path to the static directory and
// path to the index file within that static directory are used to
// serve the SPA in the given static directory.
type spaHandler struct {
	staticPath string
	indexPath  string
}

// ServeHTTP inspects the URL path to locate a file within the static dir
// on the SPA handler. If a file is found, it will be served. If not, the
// file located at the index path on the SPA handler will be served. This
// is suitable behavior for serving an SPA (single page application).
func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// prepend the path with the path to the static directory
	path = filepath.Join(h.staticPath, path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}
