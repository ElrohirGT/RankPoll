package main

import (
	"time"

	"github.com/google/uuid"
)

type Rank struct {
	Option   string
	Position uint
}

type Vote struct {
	Username string
	Ranking  []Rank
}

type PollSummary struct {
	Winner          string
	WinnerVoteCount uint
	TotalVoteCount  uint
	Rounds          []map[string]uint
}

type Room struct {
	Id         uuid.UUID
	Options    []string
	Votes      map[string]Vote
	Summary    *PollSummary
	ValidUntil time.Time
}
