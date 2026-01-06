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
	Lock  *sync.RWMutex
	Users map[string]string
	Rooms map[string]Room
}

var GlobalState = State{
	Lock:  &sync.RWMutex{},
	Users: make(map[string]string),
	Rooms: make(map[string]Room),
}

func CleanGlobalState() {
	log.Println("Cleaning global state...")
	clear(GlobalState.Users)
	clear(GlobalState.Rooms)
	log.Println("DONE!")
}

func main() {
	params := ParseParams()

	router := http.NewServeMux()
	MountHandlers(router)
	withMiddlewares := ApplyMiddlewares(router,
		RequestLoggerMiddleware,
		CORSMiddleware,
	)

	srv := &http.Server{
		Addr:    params.Addr,
		Handler: withMiddlewares,
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

	if err := srv.Shutdown(ctx); err != nil {
		log.Panicf("Failed to shutdown server: %s", err)
	}

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
