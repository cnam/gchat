package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"../gchat/lib/auth"
	"../gchat/lib/websocket"
)

const (
	listenAddr = "localhost:9876"
	privateToken = "VerySecretToken"
)

func RootHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "chat.tmpl", gin.H{"host": listenAddr})
}

func main() {
	router := gin.Default()

	auth.Register(router, privateToken)
	websocket.Register(router)

	router.LoadHTMLGlob("templates/*")
	router.GET("/", RootHandler)

	err := router.Run(listenAddr);

	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
