package source

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"sample/source/tool"
)

var Config *viper.Viper

func init() {
	Config = viper.New()
	Config.AddConfigPath("./config")
	env := tool.GetTool().GetEnv()
	switch env {
	case "product":
		Config.SetConfigName("config")
		gin.SetMode(gin.ReleaseMode)
	case "test":
		Config.SetConfigName("configTest")
	case "dev":
		Config.SetConfigName("configDev")
	}
	err := Config.ReadInConfig()
	if err != nil {
		logger.Fatal(err)
	}
}

func GetConfig() *viper.Viper {
	return Config
}
