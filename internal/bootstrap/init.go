package bootstrap

import (
	AppEvent "gorm-gen-skeleton/internal/event"
	"log"
	"time"

	"gorm-gen-skeleton/app/amqp"
	"gorm-gen-skeleton/app/event"
	"gorm-gen-skeleton/app/task"
	"gorm-gen-skeleton/internal/config"
	"gorm-gen-skeleton/internal/config/driver"
	"gorm-gen-skeleton/internal/crontab"
	"gorm-gen-skeleton/internal/logger"
	"gorm-gen-skeleton/internal/mq"
	"gorm-gen-skeleton/internal/redis"
	"gorm-gen-skeleton/internal/variable"
	"gorm-gen-skeleton/internal/variable/consts"
)

func init() {
	var err error
	if variable.Config, err = config.New(driver.New(), config.Options{
		BasePath: variable.BasePath,
	}); err != nil {
		log.Fatal(consts.ErrorInitConfig)
	}
	if variable.Log, err = logger.New(
		logger.WithDebug(true),
		logger.WithEncode("json"),
		logger.WithFilename(variable.BasePath+"/storage/logs/system.log"),
	); err != nil {
		log.Fatal(consts.ErrorInitLogger)
	}
	if variable.DB, err = InitMysql(); err != nil {
		log.Fatal(consts.ErrorInitDb)
	}
	if err = InitMongo(); err != nil {
		log.Fatal(consts.ErrorInitMongoDb)
	}

	// Elastic
	if variable.Elastic, err = InitElastic(); err != nil {
		log.Fatal(consts.ErrorInitElastic)
	}

	// Redis
	redisConfig := variable.Config.Get("Redis").(map[string]any)
	if redisConfig != nil && !redisConfig["disabled"].(bool) {
		variable.Redis = redis.New(
			redis.WithAddr(redisConfig["addr"].(string)),
			redis.WithPwd(redisConfig["pwd"].(string)),
			redis.WithDb(redisConfig["db"].(int)),
			redis.WithPoolSize(redisConfig["poolsize"].(int)),
			redis.WithMaxIdleConn(redisConfig["maxidleconn"].(int)),
			redis.WithMinIdleConn(redisConfig["minidleconn"].(int)),
			redis.WithMaxLifetime(time.Duration(redisConfig["maxlifetime"].(int))),
			redis.WithMaxIdleTime(time.Duration(redisConfig["maxidletime"].(int))),
		)
	}

	// Crontab
	if variable.Config.GetBool("Crontab.Enable") {
		variable.Crontab = crontab.New()
		variable.Crontab.AddFunc(task.New().Tasks()...)
		variable.Crontab.Start()
	}

	// RocketMQ
	if variable.Config.GetBool("MQ.Enable") {
		if variable.RocketMQ, err = mq.New(
			mq.WithNameServers(variable.Config.GetStringSlice("MQ.Servers")),
			mq.WithConsumerGroupName(variable.Config.GetString("MQ.ConsumerGroupName")),
			mq.WithProducerGroupName(variable.Config.GetString("MQ.ProducerGroupName")),
			mq.WithRetries(variable.Config.GetInt("MQ.Retries")),
		); err != nil {
			log.Fatal(consts.ErrorInitMQ)
		}
	}

	// Amqp
	if variable.Config.GetBool("Amqp.Enable") {
		variable.Amqp = mq.NewRabbitMq(variable.Config.GetString("Amqp.Addr"))
		consumers := (&amqp.Amqp{}).InitConsumers()
		if len(consumers) > 0 {
			variable.Amqp.Consumers(consumers...)
		}
	}

	// Event
	variable.Event = AppEvent.New()
	if err = (&event.Event{}).Init(); err != nil {
		log.Fatal(consts.ErrorInitEvent)
	}
}
