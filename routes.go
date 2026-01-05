package main

import "net/http"

func MountHandlers(router *http.ServeMux) {
	router.HandleFunc("POST /api/user", CreateOrLoginUser)
	router.HandleFunc("POST /api/poll", CreatePoll)
	router.HandleFunc("GET /api/poll/{pollId}", GetPollInfo)
	router.HandleFunc("POST /api/vote", VoteInPoll)
}
