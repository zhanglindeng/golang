package main

import (
	"log"
	"net/http"
	"github.com/satori/go.uuid"
	"github.com/gorilla/websocket"
	"gopkg.in/gin-contrib/cors.v1"
	"gopkg.in/gin-gonic/gin.v1"
	//"github.com/go-redis/redis"
	"sync"
	//"github.com/go-xorm/xorm"
	//"github.com/go-redis/redis"
)

// websocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// redis
//var redisClient *redis.Client

// uid
var uid int

var mutex sync.Mutex

// mysql
//var mysqlEngine *xorm.Engine

//func init() {
//	var err error
//	mysqlEngine, err = xorm.NewEngine("mysql", "websocket:websocket@(tcp://localhost:3306)/websocket?charset=utf8")
//	if err != nil {
//		log.Panic(err)
//	}
//	err = mysqlEngine.Ping()
//	if err != nil {
//		log.Panic(err)
//	}
//
//	redisClient = redis.NewClient(&redis.Options{
//		Addr:     "localhost:6379",
//		Password: "", // no password set
//		DB:       5,  // use default DB
//	})
//
//	_, err = redisClient.Ping().Result()
//	if err != nil {
//		log.Panic(err)
//	}
//}

// clients
var clients map[string]*websocket.Conn

// uuid uid
var uuids map[string]int

type clientData struct {
	SendUuid string `json:"send_uuid"`
	RecvUuid string `json:"recv_uuid"`
	Message  string `json:"message"`
}

type temp struct {
	Uuid string `json:"uuid"`
	Id   int `json:"id"`
}

func main() {
	router := gin.New()

	// cors
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:8989"}
	router.Use(cors.New(config))

	clients = make(map[string]*websocket.Conn)
	uuids = make(map[string]int)

	router.GET("/", func(c *gin.Context) {
		mutex.Lock()
		uid++
		mutex.Unlock()
		u := uuid.NewV4()

		uuids[u.String()] = uid

		c.String(http.StatusOK, "UUID %s => %d", u, uid)
	})

	router.GET("/users", func(c *gin.Context) {

		var length = len(uuids)
		var users = make([]*temp, length)
		i := 0
		for k, v := range uuids {
			users[i] = &temp{Uuid: k, Id: v}
			i++
			//users = append(users, &temp{Uuid: k, Id: v})
		}

		c.JSON(http.StatusOK, gin.H{"users": users[1:]})
	})

	router.GET("/ws/:uuid", func(c *gin.Context) {
		userUuid := c.Param("uuid")
		if userUuid == "" {
			// 不会运行到，直接404了
			c.AbortWithStatus(400)
			return
		}
		// 验证
		log.Println("user uuid", userUuid)

		if _, ok := uuids[userUuid]; !ok {
			c.AbortWithStatus(401)
			return
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		// token := string(com.RandomCreateBytes(32))
		clients[userUuid] = conn;
		defer func() {
			conn.Close()
			delete(clients, userUuid)
		}()

		d := &clientData{}
		for {
			err := conn.ReadJSON(d)
			if err != nil {
				log.Println("read:", err)
				break
			}
			// log.Printf("recv: %s", d)

			sendUserId, ok1 := uuids[d.SendUuid];
			recvUserId, ok2 := uuids[d.RecvUuid];
			_, ok3 := clients[d.SendUuid];
			recvConn, ok4 := clients[d.RecvUuid];
			if ok3 && ok1 && ok2 {
				if ok4 { // 在线
					recvConn.WriteJSON(d)
				} else { // 离线
					log.Println(sendUserId, recvUserId, d.Message)
				}
			} else {
				return
			}
		}
	})

	err := router.Run(":8080")
	if err != nil {
		log.Println(err)
	}
}
