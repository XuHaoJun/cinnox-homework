package service

import (
	"time"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/samber/lo"
	"github.com/xuhaojun/cinnox-homework/cnconfig"
	"github.com/xuhaojun/cinnox-homework/dto"
	"github.com/xuhaojun/cinnox-homework/repository"
)

type LineService struct {
	appConfig  *cnconfig.AppConfig
	repos      *LineServiceRepos
	lineClient *linebot.Client
}

type LineServiceRepos struct {
	Event    *repository.LineEventRepository
	CnMsg    *repository.CinnoxMessageRepository
	LineUser *repository.LineUserRepository
}

func NewLineService(appConfig *cnconfig.AppConfig, repos *LineServiceRepos, lineClient *linebot.Client) *LineService {
	return &LineService{
		appConfig:  appConfig,
		repos:      repos,
		lineClient: lineClient,
	}
}

func (svc *LineService) SaveEvents(events []*dto.LineEvent) error {
	if len(events) == 0 {
		return nil
	}
	_, err := svc.repos.Event.InsertMany(events)
	if err != nil {
		return err
	}

	msgEvents := lo.Filter(events, func(item *dto.LineEvent, index int) bool {
		return item.Type == "message"
	})
	err = svc.saveMessages(msgEvents)
	if err != nil {
		return err
	}

	eventsHasUser := lo.Filter(events, func(item *dto.LineEvent, index int) bool {
		return item.Source.Type == "user" && item.Source.UserID != ""
	})
	now := time.Now()
	lusers := lo.Map(eventsHasUser, func(item *dto.LineEvent, index int) *dto.LineUser {
		return &dto.LineUser{
			Id:        item.Source.UserID,
			UpdatedAt: now,
		}
	})
	err = svc.saveUsers(lusers)
	if err != nil {
		return err
	}

	return nil
}

func (svc *LineService) saveMessages(msgEvents []*dto.LineEvent) error {
	if len(msgEvents) == 0 {
		return nil
	}
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

func (svc *LineService) saveUsers(lusers []*dto.LineUser) error {
	if len(lusers) == 0 {
		return nil
	}
	_, err := svc.repos.LineUser.UpsertMany(lusers)
	if err != nil {
		return err
	}
	return nil
}

func (svc *LineService) GetUsers() ([]*dto.LineUser, error) {
	return svc.repos.LineUser.FindAll()
}

func (svc *LineService) GetUserMessages(userId string) ([]*dto.CinnoxMessage[map[string]interface{}], error) {
	return svc.repos.CnMsg.FindByUserId(userId)
}

func (svc *LineService) BroadcastMessage(messages ...linebot.SendingMessage) error {
	_, err := svc.lineClient.BroadcastMessage(messages...).Do()
	if err != nil {
		return err
	}
	return nil
}
