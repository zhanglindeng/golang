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

// JSONMessageCode json message code
type JSONMessageCode int

const (
	// OnlineMessageCode 上线消息code
	OnlineMessageCode JSONMessageCode = iota
	// OfflineMessageCode 下线消息code
	OfflineMessageCode
)

// JSONMessage 传输的数据json格式
type JSONMessage struct {
	Code    JSONMessageCode `json:"code"`
	Content interface{}     `json:"content"`
}

// App app
type App struct {
	id      string
	userIds map[string]bool
}

func (a *App) hasUser(userID string) bool {
	// todo 从数据库读取
	if _, ok := a.userIds[userID]; ok {
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

// User user
type User struct {
	id           string
	messageToken string
}

// GetChannelIds get user channel ids
func (u *User) GetChannelIds() map[string]bool {
	chIds := make(map[string]bool)
	chIds["ch1"] = true
	chIds["ch2"] = true
	chIds["ch3"] = true
	return chIds
}

func findMsgTokenByUserID(userID string) string {
	// todo 从数据库读取
	token, ok := users[userID]
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
	token := findMsgTokenByUserID(id)
	if token != msgToken {
		return nil, errors.New("auth user failed")
	}
	user := &User{
		id:           id,
		messageToken: msgToken,
	}
	return user, nil
}

// Client client
type Client struct {
	app        *App
	user       *User
	hub        *Hub
	conn       *websocket.Conn
	send       chan []byte
	channelIds map[string]bool
}

// Channel channel
type Channel struct {
	id      string
	clients map[*Client]bool
}

// Message chan 之间传递的消息
type Message struct {
	client  *Client
	msgType int
	message *ClientMessage
}

// MessageType 消息类型
type MessageType int

const (
	// TextMessageType 文本消息
	TextMessageType MessageType = iota
	// OnlineMessageType 上线消息
	OnlineMessageType
	// OfflineMessageType 下线消息
	OfflineMessageType
	// ImageMessageType 图片消息
	ImageMessageType
	// 暂时不支持文件消息
)

const (
	// TextMessageMaxSize 文本消息的最大长度
	TextMessageMaxSize int = 1024
	// OnlineMessageContent 上线消息内容
	OnlineMessageContent string = "online"
	// OfflineMessageContent 下线消息内容
	OfflineMessageContent string = "offline"
)

// ClientMessage client 消息体
type ClientMessage struct {
	// 哪个 app 的
	AppID string `json:"app_id"`
	// 发送到哪个 channel
	ChannelID string `json:"channel_id"`
	// 发送者的 user id
	UserID string `json:"user_id"`
	// 消息类型，固定是 TextMessageType
	MessageType MessageType `json:"message_type"`
	// 消息内容，不能超过 TextMessageMaxSize 个字节
	Content string `json:"content"`
}

// ImageMessage 图片消息
// 通过 http 上传后得到图片 URL
// type ImageMessage struct {
// 	ClientMessage
// }

// 新建图片消息
func newImageMessage(appID, chID, userID, url string) *ClientMessage {
	return &ClientMessage{
		AppID:       appID,
		ChannelID:   chID,
		UserID:      userID,
		MessageType: ImageMessageType,
		Content:     url,
	}
}

// OfflineMessage 下线消息
// type OfflineMessage struct {
// 	ClientMessage
// }

// 新建下线消息
func newOfflineMessage(appID, chID, userID string) *ClientMessage {
	return &ClientMessage{
		AppID:       appID,
		ChannelID:   chID,
		UserID:      userID,
		MessageType: OfflineMessageType,
		Content:     OfflineMessageContent,
	}
}

// OnlineMessage 上线消息
// type OnlineMessage struct {
// 	ClientMessage
// }

// 新建上线消息
func newOnlineMessage(appID, chID, userID string) *ClientMessage {
	return &ClientMessage{
		AppID:       appID,
		ChannelID:   chID,
		UserID:      userID,
		MessageType: OnlineMessageType,
		Content:     OnlineMessageContent,
	}
}

// TextMessage 文本消息
// type TextMessage struct {
// 	ClientMessage
// }

// 新建文本消息
func newTextMessage(appID, chID, userID string, content []byte) *ClientMessage {
	return &ClientMessage{
		AppID:       appID,
		ChannelID:   chID,
		UserID:      userID,
		MessageType: TextMessageType,
		Content:     string(content),
	}
}

// Hub hub
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

// 组装上线通知消息json包
// chID 给哪个channel发的
// userID 哪个user上线了
// func onlineMessage(chID, userID string) ([]byte, error) {
// 	cm := &ClientMessage2{
// 		ChannelID:   chID,
// 		ContentType: 1,
// 		Content:     "online",
// 		UserID:      userID,
// 	}
// 	return json.Marshal(cm)
// }

// 组装下线通知消息json包
// chID 给哪个channel发的
// userID 哪个user下线了
// func offlineMessage(chID, userID string) ([]byte, error) {
// 	cm := &ClientMessage2{
// 		ChannelID:   chID,
// 		ContentType: 0,
// 		Content:     "offline",
// 		UserID:      userID,
// 	}
// 	return json.Marshal(cm)
// }

// 给指定channel中的所有client发送online通知
// client是online状态才发送通知
// func notifyOnline(chID string, user *User, channel *Channel) {
// 	msg, err := onlineMessage(chID, user.id)
// 	if err != nil {
// 		log.Println("notifyOnline", chID, user.id, err)
// 		return
// 	}
// 	for c, online := range channel.clients {
// 		if online {
// 			c.send <- msg
// 		}
// 	}
// }

// 有client上线
// 遍历client拥有的channel
// 判断channel是否已经创建，把client添加到channel中
// 没有创建新的channel，并添加到hub中
// 如果channel已经创建，给channel中状态是online(true)的client发送online通知
func online(h *Hub, client *Client) {
	// chID => channel
	tempMap := make(map[string]*Channel)

	for chID := range client.channelIds {
		channel, ok := h.channels[chID]
		if ok {
			tempMap[chID] = channel
			// 有bug，多个json包会同时发送，client解析json错误
			// 组装成一个json包发送，包含一个channel数组
			// notifyOnline(chID, client.user, channel)
			// 给channel 添加 client
			channel.clients[client] = true
			// 客户端登陆成功后，获取未读消息，channel
		} else {
			// 创建新的 channel
			newChannel := &Channel{
				id:      chID,
				clients: make(map[*Client]bool),
			}
			newChannel.clients[client] = true
			// 添加到 hub
			h.channels[chID] = newChannel
		}
	}

	// 要通知的 channel 数组
	onlineMsgArr := make([]*ClientMessage, len(tempMap))
	index := 0
	for id := range tempMap {
		onlineMsgArr[index] = newOnlineMessage(client.app.id, id, client.user.id)
		index++
	}

	jsonMsg := &JSONMessage{
		Code:    OnlineMessageCode,
		Content: onlineMsgArr,
	}
	send, err := json.Marshal(jsonMsg)
	if err != nil {
		log.Println("online 2", err)
		return
	}

	log.Println(string(send))

	// todo 暂时不发送
	// 遍历通知 channel 中在线的 client
	// for _, ch := range tempMap {
	// 	for client, online := range ch.clients {
	// 		if online {
	// 			log.Println(ch.id)
	// 			client.send <- send
	// 		}
	// 	}
	// }
}

// 发送下线通知
// 给指定channel中online的client发送有client下线的通知
// func notifyOffline(chID string, user *User, channel *Channel) {
// 	msg, err := offlineMessage(chID, user.id)
// 	if err != nil {
// 		log.Println("notifyOffline", chID, user.id)
// 		return
// 	}
// 	for c, online := range channel.clients {
// 		if online {
// 			c.send <- msg
// 		}
// 	}
// }

// 有client下线
// 遍历client所拥有的channel
// 判断channel是否在hub中（正常来说都有）
// 从channel中删除client
func offline(h *Hub, client *Client) {
	// chID => channel
	tempMap := make(map[string]*Channel)

	for chID := range client.channelIds {
		channel, ok := h.channels[chID]
		if ok {
			tempMap[chID] = channel
			// 先把在线状态设置为false，避免给自个发送通知
			channel.clients[client] = false
			// 有bug，多个json包会同时发送，client解析json错误
			// notifyOffline(chID, client.user, channel)
			// 从 channel 删除 client
			delete(channel.clients, client)
			// 这里应该是第二次关闭
			// close(client.send)
			// 如果 channel 没有人在线了，从 hub 删除掉
			if len(channel.clients) == 0 {
				delete(h.channels, channel.id)
			}
		} else {
			// 这里一般不会运行到
			fmt.Printf("unregister error: client %v, channel id %s", client, chID)
		}
	}

	// 要通知的 channel 数组
	offlineMsgArr := make([]*ClientMessage, len(tempMap))
	index := 0
	for id := range tempMap {
		offlineMsgArr[index] = newOfflineMessage(client.app.id, id, client.user.id)
		index++
	}

	jsonMsg := &JSONMessage{
		Code:    OfflineMessageCode,
		Content: offlineMsgArr,
	}
	send, err := json.Marshal(jsonMsg)
	if err != nil {
		log.Println("offline 2", err)
		return
	}

	log.Println(string(send))

	// todo 暂时不发送
	// 遍历通知 channel 中在线的 client
	// for _, ch := range tempMap {
	// 	for client, online := range ch.clients {
	// 		if online {
	// 			log.Println(ch.id)
	// 			client.send <- send
	// 		}
	// 	}
	// }

}

// 有新的消息到来时
// 组装要转发出去的消息json包
// 从hub中查找是给哪个channel发送的消息（正常来说都可以找到）
// 遍历channel中的client
// 给online的client转发消息
func onMessage(h *Hub, msg *Message) {
	// cm := &ClientMessage{
	// 	ChannelID:   message.channelID,
	// 	ContentType: message.contentType,
	// 	Content:     string(message.content),
	// 	UserID:      message.client.user.id,
	// }
	send, err := json.Marshal(msg.message)
	if err != nil {
		log.Println("onMessage", err)
		return
	}
	channel, ok := h.channels[msg.message.ChannelID]
	if ok {
		for client, online := range channel.clients {
			if online {
				client.send <- send
			}
		}
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			online(h, client)
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			offline(h, client)
		case message := <-h.message:
			onMessage(h, message)
		}
	}
}

// type ClientMessage2 struct {
// 	ChannelID   string `json:"channel_id"`
// 	ContentType int    `json:"content_type"`
// 	Content     string `json:"content"`
// 	UserID      string `json:"user_id"`
// }

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
	appID := r.URL.Query().Get("app_id")
	app, err := verifyApp(appID)
	if err != nil {
		log.Println("verifyApp", err)
		return
	}

	// user id
	userID := r.URL.Query().Get("user_id")

	// user message token
	messageToken := r.URL.Query().Get("message_token")

	user, err := authUser(app, userID, messageToken)
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
		// 解析client的json消息
		msg, err := parseClientMessage(message)
		if err != nil {
			log.Println(err)
			break
		}
		msg.client = c
		msg.msgType = messageType
		c.hub.message <- msg
	}
}

// type ClientMessage struct {
// 	ChannelID   string `json:"channel_id"`
// 	ContentType int    `json:"content_type"`
// 	Content     string `json:"content"`
// }

func parseClientMessage(b []byte) (*Message, error) {
	cm := &ClientMessage{}
	err := json.Unmarshal(b, cm)
	if err != nil {
		return nil, err
	}
	msg := &Message{
		message: cm,
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
