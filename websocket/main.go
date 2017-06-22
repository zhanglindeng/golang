package main

import (
	"errors"
	"flag"
	"html/template"
	"log"
	"net/http"
	"time"

	"encoding/json"

	"fmt"

	"github.com/gorilla/websocket"
)

type App struct {
	id      string
	userIds map[string]bool
}

func (a *App) hasUser(userId string) bool {
	// todo 从数据库读取
	if _, ok := a.userIds[userId]; ok {
		return true
	}
	return false
}

func verifyApp(id string) (*App, error) {
	// todo 读取数据验证 app id 是否可用
	userIds, ok := apps[id]
	if !ok {
		return nil, errors.New("app id not found")
	}
	app := &App{
		id:      id,
		userIds: userIds,
	}
	return app, nil
}

type User struct {
	id           string
	messageToken string
}

func (u *User) GetChannelIds() map[string]bool {
	chIds := make(map[string]bool)
	chIds["ch1"] = true
	chIds["ch2"] = true
	chIds["ch3"] = true
	return chIds
}

func findMsgTokenByUserId(userId string) string {
	// todo 从数据库读取
	token, ok := users[userId]
	if ok {
		return token
	}
	return ""
}

func authUser(app *App, id, msgToken string) (*User, error) {
	if msgToken == "" {
		return nil, errors.New("message token required")
	}
	if ok := app.hasUser(id); !ok {
		return nil, errors.New("app user not found")
	}
	token := findMsgTokenByUserId(id)
	if token != msgToken {
		return nil, errors.New("auth user failed")
	}
	user := &User{
		id:           id,
		messageToken: msgToken,
	}
	return user, nil
}

type Client struct {
	app        *App
	user       *User
	hub        *Hub
	conn       *websocket.Conn
	send       chan []byte
	channelIds map[string]bool
}

type Channel struct {
	id      string
	clients map[*Client]bool
}

type Message struct {
	channelID   string
	client      *Client
	msgType     int
	content     []byte
	contentType int
}

type Hub struct {
	clients    map[*Client]bool
	channels   map[string]*Channel
	register   chan *Client
	unregister chan *Client
	message    chan *Message
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		channels:   make(map[string]*Channel),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		message:    make(chan *Message),
	}
}

func onlineMessage(chId, userId string) ([]byte, error) {
	cm := &ClientMessage2{
		ChannelID:   chId,
		ContentType: 1,
		Content:     "online",
		UserID:      userId,
	}
	return json.Marshal(cm)
}

func offlineMessage(chId, userId string) ([]byte, error) {
	cm := &ClientMessage2{
		ChannelID:   chId,
		ContentType: 0,
		Content:     "offline",
		UserID:      userId,
	}
	return json.Marshal(cm)
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			for channelId, _ := range client.channelIds {
				channel, ok := h.channels[channelId]
				if ok {
					// 通知上线
					for c, online := range channel.clients {
						if online {
							// 上线通知消息
							m, err := onlineMessage(channelId, client.user.id)
							log.Println("m :", string(m))
							if err == nil {
								c.send <- m
								m = []byte{}
							} else {
								log.Println(channelId, client.user.id, err)
							}
						}
					}
					// 给channel 添加 client
					channel.clients[client] = true
					// 客户端登陆成功后，获取未读消息，channel
				} else {
					// 创建新的 channel
					newChannel := &Channel{
						id:      channelId,
						clients: make(map[*Client]bool),
					}
					newChannel.clients[client] = true
					// 添加到 hub
					h.channels[channelId] = newChannel
				}
			}
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			for channelId, _ := range client.channelIds {
				channel, ok := h.channels[channelId]
				if ok {
					// 通知下线
					for c, online := range channel.clients {
						if online {
							// 通知下线消息
							m, err := offlineMessage(channelId, client.user.id)
							if err == nil {
								if c != client {
									c.send <- m
								}
							} else {
								log.Println(channelId, client.user.id, err)
							}
						}
					}
					// 从 channel 删除 client
					channel.clients[client] = false
					// 如果 channel 没有人在线了，是否删除掉？
				} else {
					// 这里一般不会运行到
					fmt.Printf("unregister error: client %v", client)
				}
			}
		case message := <-h.message:
			// log.Println(message.channelID)
			cm := &ClientMessage2{
				ChannelID:   message.channelID,
				ContentType: message.contentType,
				Content:     string(message.content),
				UserID:      message.client.user.id,
			}
			send, err := json.Marshal(cm)
			if err == nil {
				channel, ok := h.channels[message.channelID]
				if ok {
					for client := range channel.clients {
						select {
						case client.send <- send:
							if client.app.id == message.client.app.id {
								client.send <- send
							}
							// default:
							// 	close(client.send)
							// 	delete(h.clients, client)
						}
					}
				}
			}
		}
	}
}

type ClientMessage2 struct {
	ChannelID   string `json:"channel_id"`
	ContentType int    `json:"content_type"`
	Content     string `json:"content"`
	UserID      string `json:"user_id"`
}

var addr = flag.String("addr", ":8080", "http service address")

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	t, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, "Server Error", 500)
		return
	}

	var data = struct {
		App          string
		User         string
		MessageToken string
	}{
		App:          r.URL.Query().Get("app_id"),
		User:         r.URL.Query().Get("user_id"),
		MessageToken: r.URL.Query().Get("message_token"),
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, "Server Error", 500)
		return
	}

	// http.ServeFile(w, r, "index.html")
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {

	// application id
	appId := r.URL.Query().Get("app_id")
	app, err := verifyApp(appId)
	if err != nil {
		log.Println("verifyApp", err)
		return
	}

	// user id
	userId := r.URL.Query().Get("user_id")

	// user message token
	messageToken := r.URL.Query().Get("message_token")

	user, err := authUser(app, userId, messageToken)
	if err != nil {
		log.Println("authUser", err)
		return
	}

	// get user all channel ids

	// websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// client
	client := &Client{
		app:        app,
		user:       user,
		hub:        hub,
		conn:       conn,
		send:       make(chan []byte, 256),
		channelIds: make(map[string]bool),
	}
	client.hub.register <- client
	client.channelIds = user.GetChannelIds()

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				log.Println("error closed")
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Println(err, "NextWriter")
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				log.Println(err, "Close")
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println(err, "ticker WriteMessage")
				return
			}
		}
	}
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		messageType, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		log.Println(messageType, string(message))
		//content := bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		msg, err := parseClientMessage(message)
		if err != nil {
			log.Println(err)
			break
		}
		msg.client = c
		msg.msgType = messageType
		log.Println(msg.client, msg.channelID, msg.content, msg.contentType, msg.msgType)
		c.hub.message <- msg
	}
}

type ClientMessage struct {
	ChannelID   string `json:"channel_id"`
	ContentType int    `json:"content_type"`
	Content     string `json:"content"`
}

func parseClientMessage(b []byte) (*Message, error) {
	cm := &ClientMessage{}
	err := json.Unmarshal(b, cm)
	if err != nil {
		return nil, err
	}
	msg := &Message{
		content:     []byte(cm.Content),
		channelID:   cm.ChannelID,
		contentType: cm.ContentType,
	}
	return msg, nil
}

var (
	// app id => user ids => true
	apps map[string]map[string]bool
	// user id => message token
	users map[string]string
)

func initData() {
	users = make(map[string]string)
	users["u1"] = "u1_mt1"
	users["u2"] = "u2_mt2"
	users["u3"] = "u2_mt2"
	users["u4"] = "u2_mt2"
	users["u5"] = "u2_mt2"

	userIds1 := make(map[string]bool)
	userIds1["u1"] = true
	userIds1["u2"] = true
	userIds1["u3"] = true
	userIds1["u4"] = true

	userIds2 := make(map[string]bool)
	userIds2["u2"] = true
	userIds2["u3"] = true
	userIds2["u5"] = true

	apps = make(map[string]map[string]bool)
	apps["app1"] = userIds1
	apps["app2"] = userIds2
}

func main() {
	flag.Parse()
	hub := newHub()
	go hub.run()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	// init apps, users
	initData()
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
