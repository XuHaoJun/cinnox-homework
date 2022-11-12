package cnserver

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/xuhaojun/cinnox-homework/cnconfig"
	"github.com/xuhaojun/cinnox-homework/controller"
	"github.com/xuhaojun/cinnox-homework/repository"
	"github.com/xuhaojun/cinnox-homework/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AppServer struct {
	AppConfig *cnconfig.AppConfig
	MgoClient *mongo.Client
	MgoMainDb *mongo.Database
}

func NewServer() *AppServer {
	appConfig := cnconfig.LoadConfig()
	mgoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(appConfig.MongoURI))
	if err != nil {
		log.Fatalln(err)
	}
	mainDb := mgoClient.Database(repository.DbNameMain)
	return &AppServer{
		AppConfig: &appConfig,
		MgoClient: mgoClient,
		MgoMainDb: mainDb,
	}
}

func (s *AppServer) Run() {
	eventRepo := repository.NewLineEventRepository(s.MgoMainDb)
	err := eventRepo.CreateIndexes()
	if err != nil {
		log.Fatalln("create indexes failed")
	}

	cnMsgRepo := repository.NewCinnoxMessageRepository(s.MgoMainDb)
	err = cnMsgRepo.CreateIndexes()
	if err != nil {
		log.Fatalln("create indexes failed")
	}

	lineSvc := service.NewLineService(s.AppConfig, &service.LineServiceRepos{
		Event: eventRepo,
		CnMsg: cnMsgRepo,
	})
	lineController := controller.NewLineController(s.AppConfig, lineSvc)

	r := gin.Default()
	v1 := r.Group("/v1")
	lineController.AddRoutes(v1)

	r.Run()
}
