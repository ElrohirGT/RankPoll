package main

import "net/http"

func MountHandlers(router *http.ServeMux) {
	router.HandleFunc("/api/user", CreateOrLoginUser)
	router.HandleFunc("/api/poll", CreatePoll)
	router.HandleFunc("GET /api/poll/{pollId}", GetPollInfo)
	router.HandleFunc("/api/vote", VoteInPoll)
}
