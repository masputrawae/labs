package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"reflect"
	"strings"
	"syscall"
	"time"
	"todo/internal/controller/handler"
	"todo/internal/controller/middleware"
	"todo/internal/infra/database"
	"todo/internal/repo"
	"todo/internal/usecase"

	"github.com/go-playground/validator/v10"
)

func main() {
	appCtx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name, _, _ := strings.Cut(field.Tag.Get("json"), ",")
		if name == "-" {
			return ""
		}
		if name == "" {
			return field.Name
		}
		return name
	})

	initCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db, err := database.InitSQLite3(initCtx, "./app.db")
	if err != nil {
		log.Fatal("error connect db:", err)
	}

	repoUser := repo.NewUser(db)
	repoSession := repo.NewSession(db)

	usecaseUser := usecase.NewUser(repoUser)
	usecaseSession := usecase.NewSession(repoSession)

	handlerUser := handler.NewUser(usecaseUser, usecaseSession, validate)
	mw := middleware.New(usecaseSession)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/register", handlerUser.Register)
	mux.HandleFunc("POST /api/login", handlerUser.Login)
	mux.HandleFunc("GET /api/me", mw.Auth(handlerUser.Get))
	mux.HandleFunc("PATCH /api/me/update", mw.Auth(handlerUser.Update))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mw.Logger(mux),
	}

	go func() {
		log.Println("server running on :8080")

		if err := srv.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			log.Printf("server error: %v", err)
			stop()
		}
	}()

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				log.Println("last clean sessions at", time.Now().UTC())
				cleanupCtx, cancel := context.WithTimeout(
					appCtx,
					5*time.Second,
				)

				err := usecaseSession.DeleteHasExpired(cleanupCtx)
				cancel()

				if err != nil {
					log.Println("error auto delete session:", err)
				}

			case <-appCtx.Done():
				log.Println("cleanup worker stopped")
				return
			}
		}
	}()

	<-appCtx.Done()
	log.Println("signal received, shutting down")

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Println("server shutdown error:", err)
	}

	log.Println("server successfully shutdown")
}
