package main

import (
	"log"
	"net/http"
)

type Middleware = func(http.Handler) http.Handler

func ApplyMiddlewares(mux http.Handler, mds ...Middleware) http.Handler {
	finalMux := mux
	for _, md := range mds {
		finalMux = md(finalMux)
	}
	return finalMux
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Vary", "Origin")
		w.Header().Add("Vary", "Access-Control-Request-Method")
		w.Header().Add("Vary", "Access-Control-Request-Headers")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, token")
		w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS")

		if r.Method == http.MethodOptions {
			log.Println("`Options` CORS request detected! Responding ok!")
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
