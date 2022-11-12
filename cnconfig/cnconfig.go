package cnconfig

import (
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type AppConfig struct {
	MongoURI               string `json:"mongoURI" validate:"required"`
	LineChannelSecert      string `validate:"required"`
	LineChannelAccessToken string `validate:"required"`
}

// TODO
// priority: cli -> env -> yaml or json file
func LoadConfig() AppConfig {
	viper.SetConfigFile(".env")

	viper.SetDefault("MONGO_URI", "mongodb://localhost:27017")

	viper.AutomaticEnv()

	var err error
	err = viper.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
	}

	appConfig := AppConfig{
		MongoURI:               viper.GetString("MONGO_URI"),
		LineChannelSecert:      viper.GetString("LINE_CHANNEL_SECERT"),
		LineChannelAccessToken: viper.GetString("LINE_CHANNEL_ACCESS_TOKEN"),
	}

	err = validator.New().Struct(&appConfig)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println(err.Kind())
			fmt.Println(err.Type())
			fmt.Println(err.Value())
		}
		log.Fatalln("app config validate failed")
	}

	return appConfig
}
