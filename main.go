package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
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
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hey :)")
	})

	r.HandleFunc("/items", routes.ListItemsRoute(provider)).Methods("GET")
	r.HandleFunc("/order", routes.AddOrderRoute(provider)).Methods("POST")
	r.HandleFunc("/order", routes.ListOrderRoute(provider)).Methods("GET")

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
