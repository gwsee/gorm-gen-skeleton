节省时间与精力，更高效地打造稳定可靠的Web项目：基于Go语言和Gin框架的完善Web项目骨架。无需从零开始，直接利用这个骨架，快速搭建一个功能齐全、性能优异的Web应用。充分发挥Go语言和Gin框架的优势，轻松处理高并发、大流量的请求。构建可扩展性强、易于维护的代码架构，保证项目的长期稳定运行。同时，通过集成常用功能模块和最佳实践，减少繁琐的开发工作，使您专注于业务逻辑的实现。

该骨架每个组件之间可单独使用，组件之间松耦合，高内聚，组件的实现基于其他三方依赖包的封装。
目前该骨架实现了大多数的组件，比如事件,中间件,日志,配置,参数验证,命令行,定时任务等功能，目前可以满足大多数开发需求，后续会持续维护更新功能。

#### 单元测试

```shell
go test -v -run=TestConfig ./test
```

#### 设置环境变量并下载项目依赖
```shell
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
go mod download
```

#### 运行项目
```shell
go run ./cmd/main.go
```

#### 项目编译打包运行
```shell
go build ./cmd/main.go

// 编译
make build

// 运行
make run

// 编译与运行
make

// 运行项目
./main
```

#### 项目目录结构说明
```text
├─app
│  ├─command ---> 命令行
│  ├─controller
│  │    └─base.go ---> BaseController，主要定义了request参数验证器validator
│  ├─event
│  │  ├─entity ---> 事件实体目录
│  │  ├─listen ---> 事件监听执行脚本目录
│  │  └─event.go ---> 事件注册代码
│  │       
│  ├─middleware ---> 中间件代码目录
│  ├─request ---> 请求参数校验代码目录
│  │   └─request.go ---> 参数验证器
│  └─task ---> 定时任务代码目录
│     └─task.go ---> 注册定时任务脚本
├─cmd ---> 项目入口目录
│  └─cli ---> 项目命令行模式入口目录
├─config
│  └─config.yaml ---> 配置文件
├─internal ---> 包含第三方包的封装
├─router ---> 路由目录
│  └─router.go
├─storage ---> 日志、资源存储目录
│  └─logs
└─test ---> 单元测试目录
```

### 骨架全局变量

  该骨架中全局变量如下，可直接查看`internal/variable.go`文件。
  ```go
  var (
  // 项目更目录
	BasePath string
  
  // Log日志
	Log      *zap.Logger
  
  // 配置，Viper封装
	Config   *config.Config
  
  // 数据库Gorm
	DB       *gorm.DB
	MongoDB  *mongo.MongoDB
	Redis    *redis.Client
	Crontab  *crontab.Crontab

  // RabbitMQ
  Amqp     mq.RabbitMQInterface

  // rocketmq， 目前官方RocketMQ Golang SDK一些功能尚未完善，暂时不可用
	MQ       mq.Interface
  
  // 事件
	Event    *event.Event
  )
  ```

### 基础功能

---

#### 路由

该骨架的web框架是gin，所以路由定义可直接阅读Gin框架的文档。

在该骨架中定义注册路由需要在`router`文件夹下面的`router.go`文件中的`func (*AppRouter) Add(server *gin.Engine)`方法定义注册：

```go
server.GET("/foo", func(ctx *gin.Context) {
    ctx.String(http.StatusOK, "hello word!")
})
```

也可以通过自己定义路由的定义注册，只需要实现`gorm-gen-skeleton/internal/server/router`下面的`Interface`接口。如下示例：
在router目录下定义了一个`CustomRouter`结构体，该结构体实现了`Interface`接口

```go
package router

import (
    "net/http"
    
    "gorm-gen-skeleton/internal/server"
    "github.com/gin-gonic/gin"
)

type CustomRouter struct {
    server server.HttpServer
}

func NewCustom(srv server.HttpServer) *CustomRouter {
    return &CustomRouter{
        srv,
    }
}

func (*CustomRouter) Add(srv *gin.Engine) {
    srv.GET("/custom", func(ctx *gin.Context) {
        ctx.String(http.StatusOK, "custom router")
    })
}
```

> 需要注意的是，如果是自定义路由注册，需要修改项目`cmd`文件夹下面的`main.go`入口文件，通过`http.SetRouters(router.NewCustom(http))`注册给`gin`

#### 中间件

定义中间件与`gin`框架一样，该骨架默认实现了panic异常的中间件，可以查看`internal/server/middleware`文件夹中的`exception.go`文件。

如果需要定义其他的中间件并加载注册，可以将定义好的中间件通过`server.HttpServer`接口的`SetMiddleware(middlewares ...middleware.Interface)`方法注册加载，
比如我们实现如下自定义全局中间件`middleware/custom.go`：

```go
type Custom struct{}

func (c *Custom) Handle() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        fmt.Println("Custom middleware exec...")
    }
}
```
然后在定义路由的地方使用`server.SetMiddleware(&middleware.Custom{})`注册中间件。
定义全局路由中间件可以参考`router/router.go`中的`New`方法。

> 如果是局部中间件，可以直接在具体的路由上注册，参考gin路由中间件的用法

#### 日志

在该骨架中的日志是直接对`go.uber.org/zap`的封装，使用时，直接通过全局变量`variable.Log`访问写入日志，可直接使用zap支持的所有方法。

```go
package demo
import "gorm-gen-skeleton/internal/variable"
func Demo() {
    variable.Log.Info("info message")
}
```

日志文件默认是以`json`格式写入到`storage/logs/system.log`里面

#### 配置

配置项的定义直接在`config/config.yaml`文件中定义，并且配置的读取写入是通过封装`github.com/spf13/viper`实现，在该骨架中，只提供了如下一些获取配置的方法：

```go
type ConfigInterface interface {
	Get(key string) any
	GetString(key string) string
	GetBool(key string) bool
	GetInt(key string) int
	GetInt32(key string) int32
	GetInt64(key string) int64
	GetFloat64(key string) float64
	GetDuration(key string) time.Duration
	GetStringSlice(key string) []string
}
```

需要注意的是，骨架中对配置项的获取做了缓存的处理，第一次加载是在文件中获取，后面每次回去都是在`cache`中获取，目前`cache`默认只支持`memory`，骨架中也支持自定义`cache`的方法，只需要实现`config.CacheInterface`接口就可以，比如需要使用`redis`作为配置缓存，可以通过下面的方式处理:

```go
type ConfigRedisCache struct {}

var _ config.CacheInterface = (*ConfigRedisCache)(nil)

func (c *ConfigRedisCache) Get(key string) any {
    return nil
}

func (c *ConfigRedisCache) Set(key string, value any) bool {
    return true
}

func (c *ConfigRedisCache) Has(key string) bool {
    return true
}

func (c *ConfigRedisCache) FuzzyDelete(key string) {
    
}
```

然后将`ConfigRedisCache`结构体配置到`config.Options`中，如下所示，修改`internal/bootstrap/init.go`初始化配置的方法：

```go
variable.Config, err := config.New(driver.New(), config.Options{
	BasePath: './',
    Cache: &ConfigRedisCache{}
})
```

`config.yaml`基础配置如下：

```yaml
# http配置
HttpServer:
  Port: ":8888"
  
  # 服务模式，和gin的gin.SetMode的值是一样的
  Mode: "debug"
# socket配置
Websocket:
  WriteReadBufferSize: 2048
  HeartbeatFailMaxTimes: 4
  PingPeriod: 20
  ReadDeadline: 100
  WriteDeadline: 35
  PingMsg: "ping"
  
# 数据库配置
Database:
  # 可以查看GORM相关的配置选项
  Mysql:
    SlowThreshold: 5
    LogLevel: 4
    ConnMaxLifetime: 1
    MaxIdleConn: 2
    MaxOpenConn: 2
    ConnMaxIdleTime: 12
    Reade:
      - "root:root@tcp(192.168.1.4:3306)/test?charset=utf8mb4&loc=Local&parseTime=True"
    Write: "root:root@tcp(192.168.1.4:3306)/test?charset=utf8mb4&loc=Local&parseTime=True"
  # mongo数据库的基础配置
  Mongo:
    Enable: false
    Uri:
    MinPoolSize: 10
    MaxPoolSize: 20


Redis:
  Disabled: false
  Addr: "192.168.1.4:6379"
  Pwd: ""
  Db: 0
  PoolSize: 20
  MaxIdleConn: 30
  MinIdleConn: 10
  # 单位（秒）
  MaxLifeTime: 60
  # 单位（分）
  MaxIdleTime: 30

# 定时任务
Crontab:
  Enable: true

# 消息队列，使用rocketmq
MQ:
  Enable: true
  Servers:
    - "127.0.0.1:9876"
  Retries: 1
  ProducerGroupName: "ProducerGroup"
  ConsumerGroupName: "ConsumerGroup"

# RabbitMQ
Amqp:
  Enable: true
  Addr: "amqp://guest:guest@127.0.0.1:5672/"  
```

#### 事件机制

- 定义事件实体

  在`app/event/entity`目录下定义一个事件实体，该实体实现了`event.EventInterface`接口：

  ```go
  package entity
  
  type DemoEvent struct {}
  
  func (d *DemoEvent) GetData() any {
      return "demo param"
  }
  ```

  

- 定义事件监听

  在`app/event/listen`目录中定义一个`DemoEventListen`事件监听，并且该`DemoEventListen`结构体必须要实现`event.Interface`接口：

  ```go
  package listen
  
  import (
  	"fmt"
  	event2 "gorm-gen-skeleton/app/event/entity"
  	"gorm-gen-skeleton/internal/event"
  )
  
  type DemoEventListen struct {
  }
  
  // 可同时监听多个事件
  func (*DemoEventListen) Listen() []event.EventInterface {
  	return []event.EventInterface{&event2.DemoEvent{}}
  }
  
  func (*DemoEventListen) Process(data any) {
  	fmt.Printf("%v --> %s", data, "exec DemoEventListen.Process")
  }
  ```

- 最后需要将事件进行注册，在`app/event/event.go`文件中的`Init`方法内执行：

  ```go
  variable.Event.Register(&listen.DemoEventListen{})
  ```

- 调用事件执行

  ```go
  // 同步
  variable.Event.Dispatch(&entity.DemoEvent{})

  // 异步
  variable.Event.DispatchAsync(&entity.DemoEvent{})
  ```

#### 验证器

gin框架本身内置了`validator`校验，骨架里面只是对其参数的校验做了统一的校验入口。

通过如下方式获取进行参数的校验，并设置中文错误提示：

```go
type Param struct {
    Name  int    `binding:"required" form:"name" query:"name" json:"name"`
}
appRequest, err := AppRequest.New("zh")
if err != nil {
    return
}
var data Param
errMap := appRequest.Validator(ctx, &data)
fmt.Println(errMap)
```

骨架里面已经实现了默认的参数校验，可以在`app/request/request.go`文件中查看。并且在`controller`目录中`base.go`有一个`Validate(ctx *gin.Context, param any)`方法，在其他controller中要进行参数校验的时候，只需要继承`base`结构体，然后调用`Validate`方法。

```go
package controller

import "github.com/gin-gonic/gin"

type DemoController struct {
    base
}

type DemoRequest struct {
    Id int `binding:"required" form:"id" query:"id" json:"id"`
}

func (d *DemoController) Index(ctx *gin.Context) {
    var param DemoRequest
    if err := d.base.Validate(ctx, &param); err == nil {
		ctx.JSON(http.StatusOK, gin.H{"data": param})
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
	}
}
```

> 验证规格参考`github.com/go-playground/validator`官方文档

#### 命令行

基于`github.com/spf13/cobra`封装，骨架实现了基于gorm的gen代码生成器命令行。

- 生成数据库model与dao方法命令

  ```
  go run cmd/cli/main.go gen:model
  ```

  > 关于代码生成命令行的配置，可参考`app/command/command.go`文件下的`newGenCommand() AppCommand.Interface`方法

- 自定义定义命令

  在`app/command`目录中定义自己的命令，比如自定义一个输出`success ok`的命令

  ```go
  package command
  
  import (
      "fmt"
      "github.com/spf13/cobra"
  )
  
  type FooCommand struct {}
  
  func (f *FooCommand) Command() *cobra.Command {
      return &cobra.Command{
  		Use:   "foo",
  		Short: "命令使用简介.",
  		Long: `命令介绍.`,
  		Run: func(cmd *cobra.Command, args []string) {
  			str, _ := cmd.Flags().GetString("name")
               fmt.Printf("success, %s", str)
  		},
  	}
  }
  
  func (f *FooCommand) Flags(root *cobra.Command) {
  	root.Flags().StringP("name", "n", "world!", "命令参数")
  }
  ```

- 注册命令

  需要在`app/command/command.go`中的`RegisterCmds() []AppCommand.Interface`方法内注册自定义命令,具体使用方法可查看原文件`app/command/command.go`

- 全局标志  

  需要在`app/command/command.go`中的`GlobalFlags()`方法内命令,具体使用方法可查看原文件`app/command/command.go`

- 执行命令 

  ```go
  go run cmd/cli/main.go foo --name ok
  ```

- 查看命令信息

  ```go
  go run cmd/cli/main.go help
  
  // 或者
  go run cmd/cli/main.go foo --help
  ```
  
#### 定时任务

定时是通过封装`github.com/robfig/cron/v3`实现

- 定义定时任务方法

  在`app/task`目录下定义执行方法，比如每一分钟打印`success`字符

  ```go
  package task
  
  import "fmt"
  
  type SuccessTask struct {}
  
  // 时间规则
  func (s *SuccessTask) Rule() string {
  	return "* * * * *"
  }
  
  func (s *SuccessTask) Execute() func() {
  	return func() {
  		fmt.Println("success")
  	}
  }
  ```

- 加载定时任务

  需要在`app/task/task.go`文件中的`Tasks`方法内，加载自定义的任务，参考task目录下的`task.go`文件

#### Websocket

- 消息处理与链接关闭监听

  该骨架中的`websocket`是对`github.com/gorilla/websocket`依赖库的封装，在编写通信功能时，只需要实现`websocket.MessageHandler`接口：

  ```go
  import 	(
    AppSocket "gorm-gen-skeleton/internal/server/websocket"
  )

  type socketHandler struct {}
  
  // 消息处理
  func (s *socketHandler) OnMessage(message AppSocket.Message) {
    fmt.Println(fmt.Sprintf("mt: %v，data: %s, uuids: %v", message.MessageType, message.Data, message.Subkeys))
  }
  
  func (s *socketHandler) OnError(key string, err error) {
    fmt.Printf("socket err: %s, client: %s", err, key)
  }
  
  func (s *socketHandler) OnClose(key string) {
    fmt.Printf("socket closed. client:%s", key)
  }
  ```

- 创建链接Websocket

  ```go
  import 	(
    "github.com/google/uuid"
    "github.com/gin-gonic/gin"
    AppSocket "gorm-gen-skeleton/internal/server/websocket"
  )

  var client AppSocket.SocketClientInterface

  func init() {
    client, _ = AppSocket.NewSocket(AppSocket.WithHandler(&socketHandler{}))
  }

  func Connect(ctx *gin.Context) {
    subkey := uuid.New().String()
	  client.Connect(ctx, subkey)
  }
  ```

- 发送消息

  这里可以发送全部信息给全部用户，或者发送给用户，`AppSocket.Message`结构体中`Subkeys []string`表示需要发送给哪些用户：

  ```go
  client.WriteMessage(AppSocket.Message{
		MessageType: websocket.TextMessage,
		Data:        []byte("服务端收到消息并回复ok"),
    Subkeys:     []string{"9bae9c4f-87a8-46b1-b1b9-4f37b63a7d19"}
	})
  ```

  > 如果`Subkeys`是空切片数组，会将消息推送给全部在线用户

- 心跳消息

  websocket标准协议实现隐式心跳，Server端向Client端发送ping格式数据包,浏览器收到ping标准格式，自动将消息原路返回给服务器

- 其他方法

  - `GetAllKeys() []string`:获取所有websocket连接uuid
  - `GetClientState(key string) ClientState`:获取指定客户端在线状态

### 消息中间件

#### RabbitMQ

RabbitMQ消息中间件的使用可参考`test/rmq_test.go`单元测试，同时骨架中也实现了RabbitMQ的简单模式示例：`app/amqp`目录下可查看

全局变量`variable.Amqp`返回`mq.RabbitMQInterface`接口，接口中的方法可查看源文件查看

- 定义消费者`consumer`

  定义`consumer`，需要实现`mq.ConsumerHandler`接口，该接口可查看`internal/mq/rabbitmq.go`文件中的定义，比如定义一个简单模式的消费者：

  ```go
  import (
    "fmt"
    "gorm-gen-skeleton/internal/mq"

    "github.com/streadway/amqp"
  )

  type FooConsumer struct{}

  func (*FooConsumer) Option() mq.ConsumerOption {
    return mq.ConsumerOption{
      CommonOption: mq.CommonOption{
        Mode:       mq.SimpleMode,
        QueueName:  "foo",
        Durable:    false,
        AutoDelete: false,
        Exclusive:  false,
        NoWait:     false,
        Args:       nil,
      },
      AutoAck: true,
      NoLocal: false,
    }
  }

  func (*FooConsumer) Exec(msg <-chan amqp.Delivery) {
    for v := range msg {
      fmt.Printf("consumer message:%v\n", string(v.Body))
    }
  }
  ```

  当我们定义好消费者后，需要通过`mq.InitConsumer`中的`InitConsumers() []ConsumerHandler`进行注册，骨架中注册文件已写死，项目启动时，会执行`InitConsumers`方法，查看`app/amqp/amqp.go`文件的代码：

  ```go
  package amqp

  import (
    "gorm-gen-skeleton/app/amqp/consumer"
    "gorm-gen-skeleton/internal/mq"
  )

  type Amqp struct{}

  func (*Amqp) InitConsumers() []mq.ConsumerHandler {
    return []mq.ConsumerHandler{
      &consumer.FooConsumer{},
    }
  }
  ```

- 发送消息`producer`

  通过全局变量`variable.Amqp`的方法发送消息：

  ```go
  opts := mq.ProducerOption{
		CommonOption: mq.CommonOption{
			Mode:       mq.SimpleMode,
			QueueName:  "foo",
			Durable:    false,
			AutoDelete: false,
			Exclusive:  false,
			NoWait:     false,
			Args:       nil,
		},
		Message: amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("test message"),
		},
		Mandatory: false,
		Immediate: false,
	}
	variable.Amqp.Publish(opts)
  ```

#### RocketMQ

> 注意：目前官方RocketMQ Golang SDK一些功能尚未完善，暂时不可用

消息中间件的使用可参考`test/mq_test.go`单元测试

骨架对外提供了`mq.Interface`接口中的方法，可查看`internal/mq/rocketmq.go`文件下的`Interface`接口定义

定义consumer时，只需要实现`mq.ConsumerInterface`接口即可，详见`mq_test`单元测试

  
