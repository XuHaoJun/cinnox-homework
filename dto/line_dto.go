package dto

import (
	"errors"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type MessageType linebot.MessageType

const (
	MessageTypeText MessageType = "text"
)

func (s MessageType) IsValid() bool {
	switch s {
	case MessageTypeText:
		return true
	}

	return false
}

type LineMessage struct {
	Type MessageType `json:"type" binding:"required,lineMessageTypeEnum"`
	Text string      `json:"text" binding:"required_if=Type text"`
}

func (msg *LineMessage) ToLineSendingMessage() (linebot.SendingMessage, error) {
	switch msg.Type {
	case MessageTypeText:
		return linebot.NewTextMessage(msg.Text), nil
	default:
		return nil, errors.New("no support line message type")
	}
}

type LineBroadcastBody struct {
	Messages []LineMessage `json:"messages" binding:"required,dive,min=1,max=10"`
}
