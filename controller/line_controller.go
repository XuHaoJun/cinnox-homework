package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xuhaojun/cinnox-homework/cnconfig"
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
}

func (ct *LineController) webhook(c *gin.Context) {
	events, err := util.LineParseRequest(ct.appConfig.LineChannelSecert, c.Request)
	if err != nil {
		c.JSON(http.StatusForbidden, "incorrect body or X-Line-Signature")
		return
	}

	err = ct.lineSvc.SaveEvents(events)
	fmt.Println(err)
	if err != nil {
		c.JSON(http.StatusBadRequest, "save error")
		return
	}

	c.JSON(http.StatusOK, "ok")
}
