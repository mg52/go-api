package main

import (
	"github.com/joho/godotenv"
	"github.com/mg52/go-api/handler"
	"github.com/mg52/go-api/helper"
	"github.com/mg52/go-api/middleware"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
)

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

	godotenv.Load(".env")

	env, ok := os.LookupEnv("ENV")
	if !ok {
		env = "dev"
	}
	logrusEntry := helper.NewLogger(env)

	mux := http.NewServeMux()
	theUser := handler.NewUserHandler(logrusEntry)
	mux.Handle("/user", middleware.ChainingMiddleware(theUser, commonMiddlewares...))

	log.Print("Listening on :3000...")
	err := http.ListenAndServe(":3000", mux)
	return err
}
