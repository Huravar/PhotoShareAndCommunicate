package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"photo_service/gadgets"
	"photo_service/utils"
	"sync"
	"time"
)

var upgrader = websocket.Upgrader{ //允许跨域访问
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ClintWebsocket struct {
	UserId        uint
	UserName      string
	WebsocketConn *websocket.Conn
}

type ClintPool struct {
	clients map[*ClintWebsocket]bool
	mu      sync.Mutex
}

var clintPool = ClintPool{
	clients: make(map[*ClintWebsocket]bool),
}

type Message struct {
	Channel string
	Data    string
}

func SubscribeToRedis(PChannel string) {
	IRedisSubClient := utils.Red.Subscribe(context.Background(), PChannel)
	defer func(IRedisSubClient *redis.PubSub) {
		err := IRedisSubClient.Close()
		if err != nil {
			log.Println(PChannel, "RedisSubClient关闭失败！", err)
		}
	}(IRedisSubClient)
	for {
		msg, err := IRedisSubClient.ReceiveMessage(context.Background())
		fmt.Println(msg)
		if err != nil {
			log.Println(PChannel, "Redis 订阅失败:", err)
			return
		}
		for client := range clintPool.clients {
			if client.UserId == gadgets.StringToUint(PChannel) && clintPool.clients[client] {
				err := client.WebsocketConn.WriteMessage(websocket.TextMessage, []byte(msg.Payload))
				if err != nil {
					log.Println(PChannel, "发送信息失败", err)
					break
				}
			}
			if client.UserId == gadgets.StringToUint(PChannel) && !clintPool.clients[client] {
				log.Println(PChannel, "用户未在线！")
				break
			}

		}
	}
}

// UserBasicCommunicate
// @Summary      建立用户WebSocket长连接
// @Description  通过Token验证用户身份，升级为WebSocket连接，用于实时通讯和消息推送
// @Tags         WebSocket communicate
// @Param id     query string true "用户id"
// @Param Authorization header string true "Bearer 用户令牌"
// @Param Connection header string true "websocket protocol"
// @Success      101  {object} map[string]interface{}
// @Failure      400  {object} map[string]interface{}
// @Failure      401  {object} map[string]interface{}
// @Router       /ws [get]
func UserBasicCommunicate(c *gin.Context) {
	IUserTokenBasicInfo, err := VerifyToken(c)
	if err != nil {
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket 升级失败:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "创建websocket失败",
		})
		return
	}
	defer func(conn *websocket.Conn) {
		err = conn.Close()
		if err != nil {
			log.Println("WebSocket 关闭失败", err)
		}
	}(conn)
	IClintWebsocket := ClintWebsocket{
		UserId:        gadgets.StringToUint(IUserTokenBasicInfo.UserId),
		UserName:      IUserTokenBasicInfo.UserName,
		WebsocketConn: conn,
	}
	clintPool.mu.Lock()
	clintPool.clients[&IClintWebsocket] = true
	clintPool.mu.Unlock()
	go SubscribeToRedis(IUserTokenBasicInfo.UserId)
	go WebsocketHeartbeatCheck(conn, &IClintWebsocket)
	for {
		_, Imessage, err := conn.ReadMessage()
		if err != nil {
			log.Println("读取消息失败:", err)
			break
		}
		var msg Message
		if err = json.Unmarshal(Imessage, &msg); err != nil {
			log.Println("消息解析失败:", err)
			continue
		}
		fmt.Println(msg)
		if err = utils.Red.Publish(context.Background(), msg.Channel, msg.Data).Err(); err != nil {
			log.Println("Redis 发布失败:", err)
		}
	}
	clintPool.mu.Lock()
	delete(clintPool.clients, &IClintWebsocket)
	clintPool.mu.Unlock()
}

func WebsocketHeartbeatCheck(conn *websocket.Conn, PClintWebsocket *ClintWebsocket) {
	ITicker := time.NewTicker(30 * time.Second)
	defer ITicker.Stop()
	for {
		select {
		case <-ITicker.C:
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println("心跳检测失败，关闭连接", err)
				clintPool.mu.Lock()
				delete(clintPool.clients, PClintWebsocket)
				if err = conn.Close(); err != nil {
					log.Println("WebSocket 关闭失败", err)
					return
				}
				return
			}
		}
	}
}
