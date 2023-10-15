package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	_ "github.com/mg52/go-api/docs"
	"github.com/mg52/go-api/handler"
	"github.com/mg52/go-api/helper"
	"github.com/mg52/go-api/middleware"
	"github.com/mg52/go-api/repository"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title          Swagger Example API
// @version        1.0
// @description    This is a sample server.
// @termsOfService http://swagger.io/terms/

// @contact.name  API Support
// @contact.url   http://www.swagger.io/support
// @contact.email support@swagger.io

// @host     localhost:3000
// @BasePath /
func main() {
	if err := run(os.Args); err != nil {
		slog.Error("main error", slog.String("err", err.Error()))
		os.Exit(1)
	}
}

func run(_ []string) error {
	commonMiddlewares := []middleware.Middleware{
		middleware.MiddlewareOne,
		middleware.MiddlewareTwo,
		middleware.CORSMiddleware,
	}

	commonMiddlewaresWithAuth := []middleware.Middleware{
		middleware.MiddlewareOne,
		middleware.MiddlewareTwo,
		middleware.CORSMiddleware,
		middleware.RequireAuthentication,
	}

	godotenv.Load(".env")
	dbPort, _ := strconv.Atoi(os.Getenv("DBPORT"))
	servicePort, _ := strconv.Atoi(os.Getenv("PORT"))
	shutdownTimeOut, _ := strconv.Atoi(os.Getenv("SHUTDOWNTIMEOUT"))

	mux := http.NewServeMux()

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DBHOST"), dbPort, os.Getenv("DBUSER"), os.Getenv("DBPASSWORD"), os.Getenv("DBNAME"))
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	err = helper.CreateDatabaseObjects(db)
	if err != nil {
		return err
	}

	userRepository := repository.NewUserEntity(db)
	todoRepository := repository.NewTodoEntity(db)

	todoHandler := handler.NewTodoHandler(todoRepository)
	authHandler := handler.NewAuthHandler(userRepository)

	mux.Handle("/auth", middleware.ChainingMiddleware(authHandler, commonMiddlewares...))
	mux.Handle("/todo", middleware.ChainingMiddleware(todoHandler, commonMiddlewaresWithAuth...))
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", servicePort),
		Handler: mux,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() error {
		slog.Info(fmt.Sprintf("%s listening on 0.0.0.0:%d with %v timeout", os.Getenv("SERVICE"), servicePort, time.Duration(shutdownTimeOut)*time.Second))
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				slog.Error(err.Error())
			}
			return err
		}
		return nil
	}()

	<-stop

	slog.Info(fmt.Sprintf("%s shutting down...", os.Getenv("SERVICE")))

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(shutdownTimeOut)*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error(err.Error())
		return err
	}

	slog.Info(fmt.Sprintf("%s shut down completed.", os.Getenv("SERVICE")))

	return nil
}
