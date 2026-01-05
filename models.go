package main

import "github.com/google/uuid"

type Rank struct {
	Option   string
	Position uint
}

type Vote struct {
	Username string
	Ranking  []Rank
}

type PollSummary struct {
}

type Room struct {
	Id      uuid.UUID
	Options []string
	Votes   map[string]Vote
	Result  PollSummary
}
