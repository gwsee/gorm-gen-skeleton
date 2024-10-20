package controller

import (
	"fmt"

	AppSocket "gorm-gen-skeleton/internal/server/websocket"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var client AppSocket.SocketClientInterface

func init() {
	client, _ = AppSocket.NewSocket(AppSocket.WithHandler(&socketHandler{}))
}

type Socket struct{}

func (s *Socket) Connect(ctx *gin.Context) {
	subkey := uuid.New().String()
	client.Connect(ctx, subkey)
	client.WriteMessage(AppSocket.Message{
		MessageType: websocket.TextMessage,
		Data:        []byte(fmt.Sprintf("uuid: %s", subkey)),
	})
}

type socketHandler struct{}

func (s *socketHandler) OnMessage(message AppSocket.Message) {
	fmt.Println(fmt.Sprintf("mt: %v，data: %s, uuid: %v", message.MessageType, message.Data, message.Subkeys))
	fmt.Println(client.GetAllKeys())
	client.WriteMessage(AppSocket.Message{
		MessageType: websocket.TextMessage,
		Data:        []byte("服务端收到消息并回复ok"),
	})
}

func (s *socketHandler) OnError(key string, err error) {
	fmt.Printf("socket err: %s, client: %s", err, key)
}

func (s *socketHandler) OnClose(key string) {
	fmt.Printf("socket closed. client:%s", key)
}
