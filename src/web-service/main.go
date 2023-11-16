package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/health"
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type ApplicationConfig struct {
	Port uint16 `env:"PORT" envDefault:"8080"`
}

func main() {
	godotenv.Load()
	config := ApplicationConfig{}
	if err := env.Parse(&config); err != nil {
		log.Fatalf("Couldn't parse environment %s", err.Error())
	}

	healthController := health.NewDefaultController()

	router := http.NewServeMux()

	router.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./dist/assets/"))))
	router.HandleFunc("/vite.svg", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./dist/vite.svg")
	})
	router.HandleFunc("/health", healthController.ProvideHealth)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./dist/index.html")
	})

	log.Println("Server Started!")
	addr := fmt.Sprintf("0.0.0.0:%d", config.Port)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("error while listen and serve: %s", err.Error())
	}
}
