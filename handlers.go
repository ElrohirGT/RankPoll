package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ElrohirGT/RankPoll/response"
)

type CreateOrLoginUserRequest struct {
	Username string
	Password string
}

type CreateOrLoginUserResponse struct {
	Msg string
}

func CreateOrLoginUser(w http.ResponseWriter, r *http.Request) {
	var req CreateOrLoginUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		_ = response.
			NewResponseBuilder(http.StatusBadRequest).
			SetError("Invalid object received!", err).
			SendAsJSON(w)
		return
	}

	key := req.Username + "-_-" + req.Password
	_, loaded := GlobalState.Users.LoadOrStore(key, 0)
	var msg string
	if loaded {
		msg = fmt.Sprintf("User %s logging in!", req.Username)
	} else {
		msg = fmt.Sprintf("Registered %s user!", req.Username)
	}
	log.Println(msg)

	_ = response.
		NewResponseBuilder(http.StatusOK).
		SetBody(CreateOrLoginUserResponse{Msg: msg}).
		SendAsJSON(w)
}

func CreatePoll(w http.ResponseWriter, r *http.Request) {

}

func GetPollInfo(w http.ResponseWriter, r *http.Request) {

}

func VoteInPoll(w http.ResponseWriter, r *http.Request) {

}
