package main

import (
	"net/http"
	"golang.org/x/net/websocket"
	"log"
	"html/template"
	"os"
	//"github.com/gin-gonic/gin"
)

const (
	listenAddr = ":9876" // server address
)

var (
	pwd, _ = os.Getwd()
	RootTemp = template.Must(template.ParseFiles(pwd + "/chat.html"))
	JSON = websocket.JSON // codec for JSON
	ActiveClients = make(map[ClientConn]int) // map containing clients
)

type ClientConn struct {
	websocket *websocket.Conn
	clientIP  string
}

type ClientMessageType struct {
	Msg  string `json:"message"`
	Name string `json:"username"`
}

func init() {
	log.Println("Init server")
	http.HandleFunc("/", RootHandler)
	http.Handle("/ws", websocket.Handler(HandlerServer))
}

func HandlerServer(ws *websocket.Conn) {
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
}

func RootHandler(w http.ResponseWriter, req *http.Request) {
	err := RootTemp.Execute(w, listenAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	err := http.ListenAndServe(listenAddr, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}

	/*values := []byte;
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "templates/chat.tmpl", values)
	})

	router.Run(listenAddr);*/
}
