package service

import (
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/samber/lo"
	"github.com/xuhaojun/cinnox-homework/cnconfig"
	"github.com/xuhaojun/cinnox-homework/dto"
	"github.com/xuhaojun/cinnox-homework/repository"
)

type LineService struct {
	appConfig *cnconfig.AppConfig
	repos     *LineServiceRepos
}

type LineServiceRepos struct {
	Event *repository.LineEventRepository
	CnMsg *repository.CinnoxMessageRepository
}

func NewLineService(appConfig *cnconfig.AppConfig, repos *LineServiceRepos) *LineService {
	return &LineService{
		appConfig: appConfig,
		repos:     repos,
	}
}

func (svc *LineService) SaveEvents(events []*dto.LineEvent) error {
	if len(events) > 0 {
		_, err := svc.repos.Event.InsertMany(events)
		if err != nil {
			return err
		}
	}

	msgEvents := lo.Filter(events, func(item *dto.LineEvent, index int) bool {
		return item.Type == "message"
	})
	if len(msgEvents) > 0 {
		err := svc.saveMessages(msgEvents)
		if err != nil {
			return err
		}
	}

	return nil
}

func (svc *LineService) saveMessages(msgEvents []*dto.LineEvent) error {
	msgs := lo.Map(msgEvents,
		func(item *dto.LineEvent, index int) *dto.CinnoxMessage[linebot.Message] {
			return item.ToCinnoxMessage()
		})
	_, err := svc.repos.CnMsg.SaveLineMessages(msgs)
	if err != nil {
		return err
	}
	return nil
}
