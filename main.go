package main

import (
	"context"
	"github.com/RobotApocalypse/configuration"
	"github.com/RobotApocalypse/constants"
	"github.com/RobotApocalypse/router"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	//init loggers
	log := logrus.NewEntry(logrus.New())

	//config init
	cfg := configuration.NewViperConfig(log)
	Config, err := cfg.Init("configuration/config.json", "configuration/schema.json")
	if err != nil {
		log.Fatalln(err)
	}

	ginRouter := router.SetupRouter(log, cfg)

	//init srv
	srv := &http.Server{
		Addr:    ":" + Config.GetString(constants.Port),
		Handler: ginRouter,
	}

	// Graceful shut down of server
	graceful := make(chan os.Signal)
	signal.Notify(graceful, syscall.SIGINT)
	signal.Notify(graceful, syscall.SIGTERM)
	go func() {
		<-graceful
		log.Infoln("Shutting down server...")
		ctx, cancelFunc := context.WithTimeout(context.Background(), time.Duration(Config.GetInt(constants.SrvTimeOut)) * time.Second)
		defer cancelFunc()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("Could not do graceful shutdown: %v\n", err)
		}
	}()

	log.Println("Listening server on ", Config.GetString(constants.Port))
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Not able to start server on port: %v error: %s", Config.GetString(constants.Port), err)
	}

	log.Println("Server gracefully stopped...")
}
