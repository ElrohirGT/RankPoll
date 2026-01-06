package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
)

type Middleware = func(http.Handler) http.Handler

func ApplyMiddlewares(mux http.Handler, mds ...Middleware) http.Handler {
	finalMux := mux
	for i := len(mds) - 1; i >= 0; i-- {
		md := mds[i]
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

func RequestLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Failed to read request body: %s\n", err)
		} else {
			log.Printf("[REQ-%s %s] %s\n", r.Method, r.URL.String(), body)
		}

		newBodyStream := io.NopCloser(bytes.NewReader(body))
		r.Body = newBodyStream

		rec := httptest.NewRecorder()
		next.ServeHTTP(rec, r)

		resp := rec.Result()
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Failed to read response body: %s\n", err)
		} else {
			log.Printf("[RESP-%s] %s\n", resp.Status, body)
		}

		for key, val := range resp.Header {
			for _, v := range val {
				w.Header().Add(key, v)
			}
		}

		_, err = io.Copy(w, bytes.NewBuffer(body))
		if err != nil {
			log.Printf("Failed to copy body bytes into response: %s\n", err)
		}
	})
}
