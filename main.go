package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Params struct {
	Addr string
}

type State struct {
	Users sync.Map
	Rooms sync.Map
}

var GlobalState = State{}

func CleanGlobalState() {
	log.Println("Cleaning global state...")
	GlobalState.Users.Clear()
	GlobalState.Rooms.Clear()
	log.Println("DONE!")
}

func main() {
	params := ParseParams()

	router := http.NewServeMux()
	MountHandlers(router)
	srv := &http.Server{
		Addr:    params.Addr,
		Handler: router,
	}

	go func() {
		log.Println("Listening on:", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	go func() {
		if err := srv.Shutdown(ctx); err != nil {
			log.Panicf("Failed to shutdown server: %s", err)
		}
	}()
	<-ctx.Done()

	log.Println("Shut down... Goodbye!")
}

func ParseParams() Params {
	params := Params{}
	params.Addr, _ = LookUpParam("SRV_ADDRESS")
	return params
}

func LookUpParam(key string) (string, bool) {
	const ParamNotFoundFormat = "WARNING: %s ENV VARIABLE NOT SET!\n"

	param, found := os.LookupEnv(key)
	if !found {
		log.Printf(ParamNotFoundFormat, key)
	}

	return param, found
}
