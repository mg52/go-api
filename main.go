package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/mg52/go-api/docs"
	"github.com/mg52/go-api/handler"
	"github.com/mg52/go-api/helper"
	"github.com/mg52/go-api/middleware"
	"github.com/mg52/go-api/repository"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	service         = "go-api"
	port            = 3000
	shutdownTimeout = 10 * time.Second
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
		logrus.Printf("main error, %v", err)
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
	env, ok := os.LookupEnv("ENV")
	if !ok {
		env = "dev"
	}
	logrusEntry := helper.NewLogger(env)

	mux := http.NewServeMux()

	userRepository := repository.NewUserEntity()
	userHandler := handler.NewUserHandler(logrusEntry, userRepository)
	loginHandler := handler.NewLoginHandler(logrusEntry, userRepository)

	mux.Handle("/login", middleware.ChainingMiddleware(loginHandler, commonMiddlewares...))
	mux.Handle("/user", middleware.ChainingMiddleware(userHandler, commonMiddlewaresWithAuth...))
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() error {
		logrusEntry.Infof("%s listening on 0.0.0.0:%d with %v timeout", service, port, shutdownTimeout)
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				logrusEntry.Fatal(err)
			}
			return err
		}
		return nil
	}()

	<-stop

	logrusEntry.Infof("%s shutting down ...", service)

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logrusEntry.Fatal(err)
		return err
	}

	logrusEntry.Infof("%s down", service)

	return nil
}
