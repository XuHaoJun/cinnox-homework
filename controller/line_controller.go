package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/samber/lo"
	"github.com/xuhaojun/cinnox-homework/cnconfig"
	"github.com/xuhaojun/cinnox-homework/dto"
	"github.com/xuhaojun/cinnox-homework/service"
	"github.com/xuhaojun/cinnox-homework/util"
)

type LineController struct {
	appConfig *cnconfig.AppConfig
	lineSvc   *service.LineService
}

func NewLineController(appConfig *cnconfig.AppConfig, lineSvc *service.LineService) *LineController {
	return &LineController{
		appConfig: appConfig,
		lineSvc:   lineSvc,
	}
}

func (ct *LineController) AddRoutes(v1 *gin.RouterGroup) {
	v1g := v1.Group("/line")
	v1g.POST("/webhook", ct.webhook)
	v1g.POST("/broadcast", ct.broadcast)
}

func (ct *LineController) webhook(c *gin.Context) {
	events, err := util.LineParseRequest(ct.appConfig.LineChannelSecert, c.Request)
	if err != nil {
		c.JSON(http.StatusForbidden, "incorrect body or X-Line-Signature")
		return
	}

	err = ct.lineSvc.SaveEvents(events)
	if err != nil {
		c.JSON(http.StatusBadRequest, "save error")
		return
	}

	c.JSON(http.StatusOK, "ok")
}

func (ct *LineController) broadcast(c *gin.Context) {
	body := dto.LineBroadcastBody{}
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	msgs := lo.Map(body.Messages, func(item dto.LineMessage, index int) linebot.SendingMessage {
		msg, _ := item.ToLineSendingMessage()
		return msg
	})
	err = ct.lineSvc.BroadcastMessage(msgs...)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, "ok")
}
