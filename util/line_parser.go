package util

import (
	"net/http"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/samber/lo"
	"github.com/xuhaojun/cinnox-homework/dto"
)

func LineParseRequest(channelSecret string, r *http.Request) ([]*dto.LineEvent, error) {
	raws, err := linebot.ParseRequest(channelSecret, r)
	if err != nil {
		return nil, err
	}
	return lo.Map(raws, func(item *linebot.Event, index int) *dto.LineEvent {
		return dto.NewLineEvent(item, index)
	}), nil
}
