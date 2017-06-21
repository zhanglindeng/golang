package main

import (
	"bytes"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Hub struct {
	// 在线的 channel
	channels map[string]*Channel
	// 有新的客户端上线
	register chan *Client
	// 有客户端下线
	unregister chan *Client
	// 新消息
	message chan *Message
}

type Message struct {
	// 哪个 channel 的消息
	channelId string
	// 消息类型
	msgType int
	// 消息正文
	content []byte
}

func newHub() *Hub {
	return &Hub{
		channels:   make(map[string]*Channel),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		message:    make(chan *Message),
	}
}

func (h *Hub) run() {
	for {
		select {
		case register := <-h.register:
			//  如果已经存在 channel
			if channel, ok := h.channels[register.channelId]; ok {
				channel.clients[register] = true
			} else {
				// 不存在，获取channel信息，创建一个新的channel
				newChannel := &Channel{
					id:      register.channelId,
					clients: make(map[*Client]bool),
					// online:    make(chan *Client),
					// offline:   make(chan *Client),
					// broadcast: make(chan []byte),
				}
				h.channels[register.channelId] = newChannel
			}
			// 广播通知有人上线了
		case unregister := <-h.unregister:
			if channel, ok := h.channels[unregister.channelId]; ok {
				if _, ok2 := channel.clients[unregister]; ok2 {
					// 广播通知有人下线了
					delete(channel.clients, unregister)
					close(unregister.send)
				}
				// channel 是否还有人在线
				if len(channel.clients) == 0 {
					delete(h.channels, unregister.channelId)
					// close(channel.online)
					// close(channel.offline)
					// close(channel.broadcast)
				}
			}
		case message := <-h.message:
			if channel, ok := h.channels[message.channelId]; ok {
				for client := range channel.clients {
					select {
					case client.send <- message.content:
					default:
						close(client.send)
						delete(channel.clients, client)
					}
				}
			}
		}
	}
}

type Client struct {
	hub *Hub
	// 客户端连接
	conn *websocket.Conn
	// 发送的消息
	send chan []byte
	// channel id
	channelId string
}

type Channel struct {
	// channel id
	id string
	// 有哪些客户端
	clients map[*Client]bool
	// 有客户端上线
	// online chan *Client
	// 有客户端下线
	// offline chan *Client
	// 新消息广播
	// broadcast chan []byte
}

var addr = flag.String("addr", ":8080", "http service address")

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
	// 模仿 github.com/gorilla/websocket 的 chat 例子
	// https://github.com/gorilla/websocket/blob/master/examples/chat/home.html
	http.ServeFile(w, r, "home.html")
}

func main() {

	flag.Parse()
	hub := newHub()
	go hub.run()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func writeMessageType(n chan int) {
	time.Sleep(5 * time.Second)
	n <- 1
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// 这里的 channel id 写死了，相当于只有一个聊天频道
	// todo 动态设置 channel id
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), channelId: "abc"}
	client.hub.register <- client

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
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
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
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

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
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		msg := &Message{
			// 这里的 channel id 写死了，相当于只有一个聊天频道
			// todo 动态设置 channel id
			channelId: "abc",
			msgType:   messageType,
			content:   message,
		}
		c.hub.message <- msg
	}
}
