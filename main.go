package main

import (
	"github.com/joho/godotenv"
	_ "github.com/mg52/go-api/docs"
	"github.com/mg52/go-api/handler"
	"github.com/mg52/go-api/helper"
	"github.com/mg52/go-api/middleware"
	"github.com/mg52/go-api/repository"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"os"
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
		logrus.Printf("could not start application, %v", err)
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

	log.Print("Listening on :3000...")
	err := http.ListenAndServe(":3000", mux)
	return err
}
