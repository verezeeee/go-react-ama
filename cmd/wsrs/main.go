package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/rocketseat-education/semana-tech-go-react-server/internal/api"
	"github.com/rocketseat-education/semana-tech-go-react-server/internal/store/pgstore"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {

	ctx := context.Background()
	port := os.Getenv("WSRS_PORT")
	if port == "" {
    		port = "4000"
	}
	pool, err := pgxpool.New(ctx, fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s",
		os.Getenv("WSRS_DATABASE_USER"),
		os.Getenv("WSRS_DATABASE_PASSWORD"),
		os.Getenv("WSRS_DATABASE_HOST"),
		os.Getenv("WSRS_DATABASE_PORT"),
		os.Getenv("WSRS_DATABASE_NAME"),
	))
	if err != nil {
		panic(err)
	}

	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		panic(err)
	}

	handler := api.NewHandler(pgstore.New(pool))

	go func() {
		if err := http.ListenAndServe(":"+port, handler); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				panic(err)
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
