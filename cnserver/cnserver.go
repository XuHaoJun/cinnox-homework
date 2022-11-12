package cnserver

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/xuhaojun/cinnox-homework/cnconfig"
	"github.com/xuhaojun/cinnox-homework/controller"
	"github.com/xuhaojun/cinnox-homework/repository"
	"github.com/xuhaojun/cinnox-homework/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AppServer struct {
	AppConfig  *cnconfig.AppConfig
	MgoClient  *mongo.Client
	MgoMainDb  *mongo.Database
	LineClient *linebot.Client
}

func NewServer() *AppServer {
	appConfig := cnconfig.LoadConfig()
	mgoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(appConfig.MongoURI))
	if err != nil {
		log.Fatalln(err)
	}
	mainDb := mgoClient.Database(repository.DbNameMain)
	lineClient, err := linebot.New(appConfig.LineChannelSecert, appConfig.LineChannelAccessToken)
	if err != nil {
		log.Fatalln(err)
	}
	return &AppServer{
		AppConfig:  &appConfig,
		MgoClient:  mgoClient,
		MgoMainDb:  mainDb,
		LineClient: lineClient,
	}
}

func (s *AppServer) Run() {
	eventRepo := repository.NewLineEventRepository(s.MgoMainDb)
	err := eventRepo.CreateIndexes()
	if err != nil {
		log.Fatalln("create LineEventRepository indexes failed")
	}

	cnMsgRepo := repository.NewCinnoxMessageRepository(s.MgoMainDb)
	err = cnMsgRepo.CreateIndexes()
	if err != nil {
		log.Fatalln("create CinnoxMessageRepository indexes failed")
	}

	luserRepo := repository.NewLineUserRepository(s.MgoMainDb)
	err = luserRepo.CreateIndexes()
	if err != nil {
		log.Fatalln("create LineUser indexes failed")
	}

	lineSvc := service.NewLineService(s.AppConfig, &service.LineServiceRepos{
		Event:    eventRepo,
		CnMsg:    cnMsgRepo,
		LineUser: luserRepo,
	}, s.LineClient)
	lineController := controller.NewLineController(s.AppConfig, lineSvc)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("lineMessageTypeEnum", ValidateEnum)
	}
	r := gin.Default()
	r.StaticFile("/", "./demo.html")
	v1 := r.Group("/v1")
	lineController.AddRoutes(v1)

	r.Run()
}

type Enum interface {
	IsValid() bool
}

func ValidateEnum(fl validator.FieldLevel) bool {
	value := fl.Field().Interface().(Enum)
	return value.IsValid()
}
