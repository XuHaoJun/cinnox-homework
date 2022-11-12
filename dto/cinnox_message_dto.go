package dto

import (
	"time"
)

type SourceType = string

const (
	SourceTypeLine = "line"
)

type CinnoxMessage[T any] struct {
	Id         string    `json:"id" bson:"id"`
	SourceType string    `json:"sourceType" bson:"sourceType"`
	Content    T         `json:"content" bson:"content"`
	Timestamp  time.Time `json:"timestamp" bson:"timestamp"`

	// line only
	UserId      string `json:"userId" bson:"userId"`
	EventId     string `json:"eventId" bson:"eventId"`
	ReplayToken string `json:"replayToken" bson:"replayToken"`
}
