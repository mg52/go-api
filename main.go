package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
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
	"strconv"
	"syscall"
	"time"
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
	dbPort, _ := strconv.Atoi(os.Getenv("DBPORT"))
	servicePort, _ := strconv.Atoi(os.Getenv("PORT"))
	shutdownTimeOut, _ := strconv.Atoi(os.Getenv("SHUTDOWNTIMEOUT"))

	env, ok := os.LookupEnv("ENV")
	if !ok {
		env = "dev"
	}
	logrusEntry := helper.NewLogger(env)

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
	//	sqlStatement := `
	//CREATE TABLE IF NOT EXISTS users (
	//  id SERIAL PRIMARY KEY,
	//  username TEXT UNIQUE NOT NULL,
	//  password TEXT NOT NULL
	//);`
	//	_, err = db.Exec(sqlStatement)
	//	if err != nil {
	//		return err
	//	}
	//
	//	//	sqlStatement2 := `
	//	//INSERT INTO users (username, password)
	//	//VALUES ($1, $2)
	//	//RETURNING id`
	//	//	id2 := 0
	//	//	err = db.QueryRow(sqlStatement2, "mg52", "password").Scan(&id2)
	//	//	if err != nil {
	//	//		panic(err)
	//	//	}
	//	//	fmt.Println("New record ID2 is:", id2)
	//	//
	//	//	sqlStatement3 := `
	//	//INSERT INTO users (username, password)
	//	//VALUES ($1, $2)
	//	//RETURNING id`
	//	//	id3 := 0
	//	//	err = db.QueryRow(sqlStatement3, "pinar", "pass2").Scan(&id3)
	//	//	if err != nil {
	//	//		panic(err)
	//	//	}
	//	//	fmt.Println("New record ID3 is:", id3)
	//
	//	sqlStatementSelect := `SELECT * FROM users WHERE id=$1;`
	//	var user domain.User
	//	row := db.QueryRow(sqlStatementSelect, 3)
	//	errSelect := row.Scan(&user.ID, &user.Username, &user.Password)
	//	if errSelect != nil && errSelect == sql.ErrNoRows {
	//		fmt.Println("No rows were returned!")
	//	} else {
	//		fmt.Println(user)
	//	}

	userRepository := repository.NewUserEntity(db)
	// TODO: User handler will be removed and todo handler will be implemented and injected.
	userHandler := handler.NewUserHandler(logrusEntry, userRepository)
	authHandler := handler.NewAuthHandler(logrusEntry, userRepository)

	mux.Handle("/auth", middleware.ChainingMiddleware(authHandler, commonMiddlewares...))
	mux.Handle("/user", middleware.ChainingMiddleware(userHandler, commonMiddlewaresWithAuth...))
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", servicePort),
		Handler: mux,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() error {
		logrusEntry.Infof("%s listening on 0.0.0.0:%d with %v timeout", os.Getenv("SERVICE"), servicePort, time.Duration(shutdownTimeOut)*time.Second)
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				logrusEntry.Fatal(err)
			}
			return err
		}
		return nil
	}()

	<-stop

	logrusEntry.Infof("%s shutting down ...", os.Getenv("SERVICE"))

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(shutdownTimeOut)*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logrusEntry.Fatal(err)
		return err
	}

	logrusEntry.Infof("%s down", os.Getenv("SERVICE"))

	return nil
}
