package websocket

import (
	"net/http"
	"encoding/json"
	"log"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	pingMsgMeta = "system.ping"
	pingMsgData = "ping"
)

var (
	ActiveClients = make(map[ClientConn]int)
)

type ClientConn struct {
	conn *websocket.Conn
	clientIP  string
}

type ClientMessageType struct {
	Msg  string `json:"message"`
	Name string `json:"username"`
}

type ClientMessageData struct  {
	Data ClientMessageType `json:"data"`
	Meta EventStruct       `json:"meta"`
}

type MessageStruct struct {
	Data string       `json:"data"`
	Meta EventStruct  `json:"meta"`
}

type EventStruct struct {
	Event string `json:"event"`
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wshandler(writer http.ResponseWriter, req *http.Request) {
	var clientMessageData ClientMessageData
	var msgStruct MessageStruct
	var pingMsg = `{"data":"pong","meta":{"event":"system.ping"}}`

	conn, err := wsupgrader.Upgrade(writer, req, nil)
	if err != nil {
		log.Println("Failed to set websocket upgrade: %+v", err)
		return
	}

	client := req.RemoteAddr
	log.Println("Client connected:", client)

	sockCli := ClientConn{conn, client}
	ActiveClients[sockCli] = 0
	log.Println("Number of clients connected ...", len(ActiveClients))

	for {
		t, msg, err := conn.ReadMessage()

		json.Unmarshal(msg, &msgStruct)

		if err != nil {
			delete(ActiveClients, sockCli)
			return
		}

		log.Printf("%+v", msgStruct)

		if msgStruct.Meta.Event == pingMsgMeta && msgStruct.Data == pingMsgData {
			log.Println("Write system ping")
			conn.WriteMessage(t, []byte(pingMsg))
			continue
		}

		json.Unmarshal(msg, &clientMessageData)

		log.Printf("%+v", clientMessageData)

		msg, err  = json.Marshal(clientMessageData)

		for cs, _ := range ActiveClients {
			if err = cs.conn.WriteMessage(t, msg); err != nil {
				// we could not send the message to a peer
				log.Println("Could not send message to ", cs.clientIP, err.Error())
			}
		}
	}
}

func Register(router *gin.Engine) {
	router.GET("/ws", func(c *gin.Context){
		wshandler(c.Writer, c.Request)
	})
}