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

	if len(req.PollOptions) == 0 {
		_ = response.NewResponseBuilder(http.StatusBadRequest).
			SetError("Invalid option count!", errors.New("can't have a poll with 0 options")).
			SendAsJSON(w)
		return
	}

	if len(req.PollOptions) == 1 {
		_ = response.NewResponseBuilder(http.StatusBadRequest).
			SetError("Invalid option count!", errors.New("can't have a poll with only 1 option")).
			SendAsJSON(w)
		return
	}

	id := uuid.New()
	key := id.String()
	room := Room[time.Time]{
		Id:         id,
		Title:      req.Title,
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
	now := time.Now()

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
	pollInfo, found := GlobalState.Rooms[pollId.String()]
	GlobalState.Lock.RUnlock()
	log.Printf("Room already exists? %t", found)

	if !found {
		_ = response.NewResponseBuilder(http.StatusNotFound).
			SetError("Poll not found!", errors.New("the poll was not found")).
			SendAsJSON(w)
		return
	}

	shouldComputeSummary := now.After(pollInfo.ValidUntil) && pollInfo.Summary == nil
	if shouldComputeSummary {
		computeSummary(&pollInfo)
	}

	GlobalState.Lock.Lock()
	GlobalState.Rooms[pollId.String()] = pollInfo
	GlobalState.Lock.Unlock()

	posixInfo := toPosixTime(pollInfo)

	_ = response.NewResponseBuilder(http.StatusOK).
		SetBody(posixInfo).
		SendAsJSON(w)
}

func toPosixTime(r Room[time.Time]) Room[int64] {
	posixTime := r.ValidUntil.UnixMilli()
	return Room[int64]{
		Id:         r.Id,
		Title:      r.Title,
		Options:    r.Options,
		Votes:      r.Votes,
		Summary:    r.Summary,
		ValidUntil: posixTime,
	}
}

func computeSummary(room *Room[time.Time]) {
	summary := &PollSummary{
		Rounds: make([]map[string]uint, 0),
	}

	for roundIdx := range len(room.Options) {
		round := uint(roundIdx + 1)
		roundTally := make(map[string]uint)

		for _, vote := range room.Votes {
			for _, r := range vote.Ranking {
				if r.Position > round {
					continue
				}

				roundTally[r.Option] += 1
			}
		}

		summary.Rounds = append(summary.Rounds, roundTally)

		maxOpt := "INVALID"
		var maxCount uint = 0
		var totalCount uint = 0
		isUnique := true
		for opt, voteCount := range roundTally {
			totalCount += voteCount

			if voteCount > maxCount {
				maxOpt = opt
				maxCount = voteCount
				isUnique = true
			} else if voteCount == maxCount {
				isUnique = false
			}
		}

		moreThanFraction := maxCount > (totalCount / uint(len(room.Options)))
		log.Printf("Round: %d - IsUnique: %t\n", round, isUnique)
		log.Printf("Computing summary: %d > (%d / %d)\n", maxCount, totalCount, len(room.Options))
		if moreThanFraction && isUnique || round == uint(len(room.Options)) {
			log.Printf("Winner: %s", maxOpt)
			summary.Winner = maxOpt
			summary.WinnerVoteCount = maxCount
			summary.TotalVoteCount = totalCount
			break
		}
	}

	room.Summary = summary
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
