package main

import (
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

// Normally T will be time.Time but sometimes it needs to be something else.
// For example an int64 for representing Unix time.
type Room[T any] struct {
	Id         uuid.UUID
	Title      string
	Options    []string
	Votes      map[string]Vote
	Summary    *PollSummary
	ValidUntil T
}
