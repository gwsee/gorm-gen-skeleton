package variable

import (
	"gorm-gen-skeleton/internal/elasticsearch"
	"gorm-gen-skeleton/internal/event"
	"gorm-gen-skeleton/internal/mongo"
	"log"
	"os"
	"strings"

	"gorm-gen-skeleton/internal/config"
	"gorm-gen-skeleton/internal/crontab"
	"gorm-gen-skeleton/internal/mq"
	"gorm-gen-skeleton/internal/variable/consts"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	BasePath string
	Log      *zap.Logger
	Config   *config.Config
	DB       *gorm.DB
	MongoDB  *mongo.MongoDB
	Redis    *redis.Client
	Crontab  *crontab.Crontab
	Amqp     mq.RabbitMQInterface
	Event    *event.Event
	Elastic  *elasticsearch.Elasticsearch

	// RocketMQ 目前官方RocketMQ Golang SDK一些功能尚未完善，暂时不可用
	RocketMQ mq.Interface
)

func init() {
	if curPath, err := os.Getwd(); err == nil {
		if len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "-test") {
			BasePath = strings.Replace(strings.Replace(curPath, `\test`, "", 1), `/test`, "", 1)
		} else {
			BasePath = curPath
		}
	} else {
		log.Fatal(consts.ErrorsBasePath)
	}
}
