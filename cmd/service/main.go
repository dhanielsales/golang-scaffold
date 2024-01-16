package main

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"

	"github.com/dhanielsales/golang-scaffold/config"
	_ "github.com/dhanielsales/golang-scaffold/docs"
	"github.com/dhanielsales/golang-scaffold/init/shutdown"
	"github.com/dhanielsales/golang-scaffold/internal/http"
	"github.com/dhanielsales/golang-scaffold/internal/postgres"
	"github.com/dhanielsales/golang-scaffold/internal/redis"

	"github.com/dhanielsales/golang-scaffold/modules/store"
)

type service struct {
	http     *http.HttpServer
	postgres *postgres.Storage
	redis    *redis.Storage
	env      *config.EnvVars
	validate *validator.Validate
}

func new(env *config.EnvVars) (*service, error) {
	// init the Postgres storage
	postgres, err := postgres.Bootstrap(env.POSTGRES_URL)
	if err != nil {
		return nil, err
	}

	// init the Redis storage
	redis, err := redis.Bootstrap(env.REDIS_URL)
	if err != nil {
		return nil, err
	}

	// init http server
	httpServer := http.Bootstrap(env.PORT)

	// init validator
	validate := validator.New(validator.WithRequiredStructEnabled())
	validator := http.NewValidator(validate)

	// Start store module
	store.Bootstrap(postgres, redis, httpServer, validator)

	return &service{
		http:     httpServer,
		postgres: postgres,
		redis:    redis,
		env:      env,
	}, nil
}

func (s *service) run() {
	s.http.Start()
}

func (s *service) cleanup() {
	s.http.Cleanup()
	s.postgres.Cleanup()
}

// @title Golang scaffold
// @version 1.0
// @description A simple Golang backend scaffold
// @contact.name Dhaniel Sales
// @license.name MIT
// @BasePath /
func main() {
	// setup exit code for graceful shutdown
	var exitCode int
	defer func() {
		os.Exit(exitCode)
	}()

	// load config
	env, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("error: %v", err)
		exitCode = 1
		return
	}

	// Create new service
	srv, err := new(env)
	if err != nil {
		fmt.Printf("error: %v", err)
		exitCode = 1
		return
	}

	// run the cleanup after the server is terminated
	defer srv.cleanup()

	// run the server
	go srv.run()

	// ensure the server is shutdown gracefully & app runs
	shutdown.Gracefully()
}
