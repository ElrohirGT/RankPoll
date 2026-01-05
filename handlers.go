package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ElrohirGT/RankPoll/response"
	"github.com/google/uuid"
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

type CreatePollRequest struct {
	Title       string
	PollOptions []string
	PollUntil   time.Time
}
type CreatePollResponse struct {
	PollId uuid.UUID
	Msg    string
}

func CreatePoll(w http.ResponseWriter, r *http.Request) {
	var req CreatePollRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		_ = response.NewResponseBuilder(http.StatusBadRequest).
			SetError("Invalid object received!", err).
			SendAsJSON(w)
		return
	}

	id := uuid.New()
	key := id.String()
	room := Room{
		Id:      id,
		Options: req.PollOptions,
	}
	log.Printf("Storing new poll with id: %s\n", id)
	GlobalState.Rooms.Store(key, room)

	_ = response.NewResponseBuilder(http.StatusOK).
		SetBody(CreatePollResponse{Msg: "Success!", PollId: id}).
		SendAsJSON(w)
}

func GetPollInfo(w http.ResponseWriter, r *http.Request) {
	pollStrId := r.PathValue("pollId")
	if pollStrId == "" {
		_ = response.NewResponseBuilder(http.StatusNotFound).
			SetError("Poll not found!", errors.New("no pollId supplied")).
			SendAsJSON(w)
		return
	}
	log.Printf("The poll id from the path is: %s", pollStrId)

	pollId, err := uuid.Parse(pollStrId)
	if err != nil {
		_ = response.NewResponseBuilder(http.StatusNotFound).
			SetError("Poll not found!", err).
			SendAsJSON(w)
		return
	}

	pollInfo, found := GlobalState.Rooms.Load(pollId.String())
	log.Printf("Room already exists? %t", found)
	if !found {
		_ = response.NewResponseBuilder(http.StatusNotFound).
			SetError("Poll not found!", err).
			SendAsJSON(w)
		return
	}

	_ = response.NewResponseBuilder(http.StatusOK).
		SetBody(pollInfo).
		SendAsJSON(w)
}

func VoteInPoll(w http.ResponseWriter, r *http.Request) {

}
