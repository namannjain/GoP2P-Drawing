package main

import (
	"fmt"
	"log"
	"net/http"

	"goP2Pbackend/config"
	"goP2Pbackend/internal/delivery/http/handler"
	"goP2Pbackend/internal/repository/postgres"
	"goP2Pbackend/internal/repository/s3"
	"goP2Pbackend/internal/usecase"
	"goP2Pbackend/pkg/auth"
	websocket "goP2Pbackend/pkg/ws"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading the ENV file")
		panic(err)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := postgres.NewDB(cfg.GetDatabaseURL())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	s3Client, err := s3.NewClient(cfg.AWS.Region, cfg.AWS.AccessKeyID, cfg.AWS.SecretAccessKey)
	if err != nil {
		log.Fatalf("Failed to create S3 client: %v", err)
	}

	userRepo := postgres.NewUserRepository(db)
	artboardRepo := postgres.NewArtboardRepository(db)
	artboardStorage := s3.NewArtboardStorage(s3Client, cfg.AWS.BucketName)

	userUsecase := usecase.NewUserUsecase(userRepo)
	artboardUsecase := usecase.NewArtboardUsecase(artboardRepo, artboardStorage)

	oauthConfig := auth.NewOAuthConfig(cfg.OAuth.GoogleClientID, cfg.OAuth.GoogleClientSecret, cfg.OAuth.OAuthRedirectURL)

	userHandler := handler.NewUserHandler(userUsecase, oauthConfig)
	artboardHandler := handler.NewArtboardHandler(artboardUsecase)

	hub := websocket.NewHub()
	go hub.Run()

	r := mux.NewRouter()

	// User routes
	r.HandleFunc("/auth/google", userHandler.GoogleLogin)
	r.HandleFunc("/auth/google/callback", userHandler.GoogleCallback)

	// Artboard routes
	r.HandleFunc("/artboards", artboardHandler.Create).Methods("POST")
	r.HandleFunc("/artboards", artboardHandler.List).Methods("GET")
	r.HandleFunc("/artboards/{id}", artboardHandler.Get).Methods("GET")
	r.HandleFunc("/artboards/{id}", artboardHandler.Update).Methods("PUT")
	r.HandleFunc("/artboards/{id}", artboardHandler.Delete).Methods("DELETE")
	r.HandleFunc("/artboards/{id}/share", artboardHandler.GenerateShareableLink).Methods("POST")

	// WebSocket route
	r.HandleFunc("/ws/{artboardID}", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(hub, w, r)
	})

	log.Printf("Server starting on port %s", cfg.Server.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Server.Port, r))
}
