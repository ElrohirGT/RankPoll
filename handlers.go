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
		_ = response.NewResponseBuilder(http.StatusBadRequest).
			SetError("Invalid object received!", err).
			SendAsJSON(w)
		return
	}

	GlobalState.Lock.Lock()
	defer GlobalState.Lock.Unlock()

	var msg string
	if password, found := GlobalState.Users[req.Username]; !found {
		GlobalState.Users[req.Username] = req.Password
		msg = fmt.Sprintf("Registered %s user!", req.Username)
	} else {
		if password != req.Password {
			_ = response.NewResponseBuilder(http.StatusBadRequest).
				SetError("Invalid credentials", errors.New("password/username don't match")).
				SendAsJSON(w)
			return
		}
		msg = fmt.Sprintf("User %s logging in!", req.Username)
	}
	log.Println(msg)

	_ = response.NewResponseBuilder(http.StatusOK).
		SetBody(CreateOrLoginUserResponse{Msg: msg}).
		SendAsJSON(w)
}

type CreatePollRequest struct {
	Title           string
	PollOptions     []string
	PollingDuration time.Duration
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
		Id:         id,
		Options:    req.PollOptions,
		Votes:      make(map[string]Vote),
		ValidUntil: time.Now().Add(req.PollingDuration),
	}
	log.Printf("Storing new poll with id: %s\n", id)

	GlobalState.Lock.Lock()
	defer GlobalState.Lock.Unlock()
	GlobalState.Rooms[key] = room

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

	GlobalState.Lock.RLock()
	defer GlobalState.Lock.RUnlock()

	pollInfo, found := GlobalState.Rooms[pollId.String()]
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

type VoteInPollRequest struct {
	Username string
	PollId   uuid.UUID
	Options  map[string]uint
}

func VoteInPoll(w http.ResponseWriter, r *http.Request) {
	now := time.Now()

	var req VoteInPollRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		_ = response.NewResponseBuilder(http.StatusBadRequest).
			SetError("Invalid body!", err).
			SendAsJSON(w)
		return
	}

	GlobalState.Lock.RLock()
	roomInfo, found := GlobalState.Rooms[req.PollId.String()]
	GlobalState.Lock.RUnlock()
	if !found {
		_ = response.NewResponseBuilder(http.StatusNotFound).
			SetError("The room was not found!", err).
			SendAsJSON(w)
		return
	}

	if _, found := roomInfo.Votes[req.Username]; found {
		_ = response.NewResponseBuilder(http.StatusBadRequest).
			SetError("The user has already voted!", errors.New("a user can't vote twice")).
			SendAsJSON(w)
		return
	}

	if now.After(roomInfo.ValidUntil) {
		_ = response.NewResponseBuilder(http.StatusBadRequest).
			SetError("The poll already ended!", errors.New("the poll has ended")).
			SendAsJSON(w)
		return
	}

	ranking := make([]Rank, 0, len(roomInfo.Options))
	for _, opt := range roomInfo.Options {
		position, found := req.Options[opt]
		if !found {
			_ = response.NewResponseBuilder(http.StatusBadRequest).
				SetError("Incomplete voting options!", fmt.Errorf("no option %s found", opt)).
				SendAsJSON(w)
			return
		}

		if position == 0 {
			_ = response.NewResponseBuilder(http.StatusBadRequest).
				SetError("The 0 rank is not existent!", fmt.Errorf("option %s has 0 rank", opt)).
				SendAsJSON(w)
			return
		}

		if position > uint(len(roomInfo.Options)) {
			_ = response.NewResponseBuilder(http.StatusBadRequest).
				SetError("An option has a rank greater than voting options!", fmt.Errorf("option %s has a big rank", opt)).
				SendAsJSON(w)
			return
		}

		ranking = append(ranking, Rank{
			Option:   opt,
			Position: req.Options[opt],
		})
	}

	roomInfo.Votes[req.Username] = Vote{
		Username: req.Username,
		Ranking:  ranking,
	}

	GlobalState.Lock.Lock()
	GlobalState.Rooms[req.PollId.String()] = roomInfo
	GlobalState.Lock.Unlock()

	_ = response.NewResponseBuilder(http.StatusOK).
		SendAsJSON(w)
}
