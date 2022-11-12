package dto

import (
	"fmt"
	"log"
	"time"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LineEvent struct {
	ReplyToken        string
	Type              linebot.EventType
	Mode              linebot.EventMode
	Timestamp         time.Time
	Source            *linebot.EventSource
	Message           linebot.Message
	Joined            *linebot.Members
	Left              *linebot.Members
	Postback          *linebot.Postback
	Beacon            *linebot.Beacon
	AccountLink       *linebot.AccountLink
	Things            *linebot.Things
	Members           []*linebot.EventSource
	Unsend            *linebot.Unsend
	VideoPlayComplete *linebot.VideoPlayComplete
	WebhookEventID    string
	DeliveryContext   linebot.DeliveryContext
	//
	UniqueId string `json:"uniqueId" bson:"uniqueId"`
}

func (e *LineEvent) ToCinnoxMessage() *CinnoxMessage[linebot.Message] {
	if e.Type != "message" {
		log.Fatalln("try convert non-message line event to cinnox message")
	}

	msg := CinnoxMessage[linebot.Message]{
		SourceType:  SourceTypeLine,
		Id:          e.GetId(),
		Content:     e.Message,
		Timestamp:   e.Timestamp,
		UserId:      e.Source.UserID,
		EventId:     e.UniqueId,
		ReplayToken: e.ReplyToken,
	}
	return &msg
}

func (e *LineEvent) GetId() string {
	switch v := e.Message.(type) {
	case *linebot.TextMessage:
		return v.ID
	case *linebot.ImageMessage:
		return v.ID
	case *linebot.VideoMessage:
		return v.ID
	case *linebot.AudioMessage:
		return v.ID
	case *linebot.FileMessage:
		return v.ID
	case *linebot.LocationMessage:
		return v.ID
	case *linebot.StickerMessage:
		return v.ID
	default:
		// TODO
		// why template, imagemap, flex no id?
		return primitive.NewObjectID().Hex()
	}
}

func NewLineEvent(raw *linebot.Event, index int) *LineEvent {
	return &LineEvent{
		ReplyToken:        raw.ReplyToken,
		Type:              raw.Type,
		Mode:              raw.Mode,
		Timestamp:         raw.Timestamp,
		Source:            raw.Source,
		Message:           raw.Message,
		Joined:            raw.Joined,
		Left:              raw.Left,
		Postback:          raw.Postback,
		Beacon:            raw.Beacon,
		AccountLink:       raw.AccountLink,
		Things:            raw.Things,
		Members:           raw.Members,
		Unsend:            raw.Unsend,
		VideoPlayComplete: raw.VideoPlayComplete,
		WebhookEventID:    raw.WebhookEventID,
		DeliveryContext:   raw.DeliveryContext,
		//
		UniqueId: fmt.Sprintf("%s-%d", raw.WebhookEventID, index),
	}
}
