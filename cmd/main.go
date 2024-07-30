package main

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	_ "github.com/swaggo/http-swagger"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	_ "message-service/docs"
	"message-service/internal/config"
	"message-service/internal/handlers"
	"message-service/internal/kafka"
	"message-service/internal/repository"
	"message-service/internal/service"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title Message Service API
// @version 1.0
// @description This is a server for a message service.
// @host localhost:8080
// @BasePath /
func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}

func run() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	cfg := config.Read()

	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
		}
	}(db)

	//DI
	repos := repository.NewRepository(db)
	producer := kafka.NewProducer(cfg.KafkaAddr)
	consumer := kafka.NewConsumer(cfg.KafkaAddr)
	services := service.NewService(repos, producer)
	handler := handlers.NewHandler(services)

	router := mux.NewRouter()

	router.HandleFunc("/api/message", handler.CreateMessage).Methods("POST")
	router.HandleFunc("/api/messages/stats", handler.GetStats).Methods("GET")

	//doc
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	srv := http.Server{
		Addr:           cfg.HTTPAddr,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		timer := time.NewTimer(5 * time.Second)
		defer timer.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Println("send to Kafka:", ctx.Err())
				return
			case <-timer.C:
				err := services.SendMsgToKafka(ctx)
				if err != nil {
					log.Println("send to Kafka:", err)
				}
				timer.Reset(5 * time.Second)
			}
		}
	}()
	go consumer.ConsumeMessage(ctx, services.ProcessMsg)

	// listen to OS signals and gracefully shutdown HTTP server
	stopped := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("HTTP Server Shutdown Error: %v", err)
		}
		close(stopped)
	}()

	log.Printf("Starting HTTP server on %s", cfg.HTTPAddr)

	// start HTTP server
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("HTTP server ListenAndServe Error: %v", err)
	}

	<-stopped

	log.Printf("Have a nice day!")

	return nil
}
