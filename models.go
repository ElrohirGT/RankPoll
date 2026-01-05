package main

import "github.com/google/uuid"

type Room struct {
	Id      uuid.UUID
	Options []string
}
