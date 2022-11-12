package dto

import (
	"time"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type CinnoxMessageContents interface {
	linebot.Message
}

type SourceType = string

const (
	SourceTypeLine = "line"
)

type CinnoxMessage[T CinnoxMessageContents] struct {
	Id         string    `json:"id" bson:"id"`
	SourceType string    `json:"sourceType" bson:"sourceType"`
	Content    T         `json:"content" bson:"content"`
	Timestamp  time.Time `json:"timestamp" bson:"timestamp"`

	// line only
	UserId      string `json:"userId" bosn:"userId"`
	EventId     string `json:"eventId" bson:"eventId"`
	ReplayToken string `json:"replayToken" bson:"replayToken"`
}
