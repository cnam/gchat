package main

import (
	"net/http"
	//"golang.org/x/net/websocket"
	"encoding/json"
	"log"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	listenAddr = "localhost:9876" // server address
	pingMsgMeta = "system.ping"
	pingMsgData = "ping"
)

var (
	ActiveClients = make(map[ClientConn]int) // map containing clients
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
	Meta EventStruct  `json:"meta"`
}

type MessageStruct struct {
	Data string    `json:"data"`
	Meta EventStruct  `json:"meta"`
}

type EventStruct struct {
	Event string `json:"event"`
}

/*func init() {
	log.Println("Init server")
	//http.Handle("/ws", websocket.Handler(WsHandler))
}*/

/*func WsHandler(ws *websocket.Conn) {
	var err error
	var clientMessage ClientMessageType

	defer func() {
		if err = ws.Close(); err != nil {
			log.Println("Websocket could not be closed", err.Error())
		}
	}()

	client := ws.Request().RemoteAddr
	log.Println("Client connected:", client)

	sockCli := ClientConn{ws, client}
	ActiveClients[sockCli] = 0
	log.Println("Number of clients connected ...", len(ActiveClients))

	for {
		if err = JSON.Receive(ws, &clientMessage); err != nil {
			log.Println("Websocket Disconnected waiting", err.Error())
			delete(ActiveClients, sockCli)
			log.Println("Number of clients still connected ...", len(ActiveClients))
			return
		}

		log.Printf("%+v", clientMessage)
		for cs, _ := range ActiveClients {
			if err = JSON.Send(cs.websocket, clientMessage); err != nil {
				// we could not send the message to a peer
				log.Println("Could not send message to ", cs.clientIP, err.Error())
			}
		}
	}
}*/

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

func RootHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "chat.tmpl", gin.H{"host": listenAddr})
}

func main() {
	router := gin.Default()

	router.LoadHTMLGlob("templates/*")
	router.GET("/", RootHandler)
	router.GET("/ws", func(c *gin.Context){
		wshandler(c.Writer, c.Request)
		/*log.Println("Connect to ws")

		s := websocket.Server{Handler: WsHandler, Handshake:
			func (config *websocket.Config, req *http.Request) (err error) {
				config.Origin, err = websocket.Origin(config, req)
				if err == nil && config.Origin == nil {
					log.Println("null origin")
				}
			return err
		}}

		s.ServeHTTP(c.Writer, c.Request)*/
	})

	err := router.Run(listenAddr);

	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
