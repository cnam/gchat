package auth

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"
	"time"
	"log"
)

var (
	privateToken []byte
)

type LoginFormType struct  {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Register(router *gin.Engine, PrivateToken string) {
	privateToken = []byte(PrivateToken)
	router.POST("/auth/login", LoginHandler)
	router.GET("/auth/login", LoginForm)
}

func LoginForm(c *gin.Context) {
	c.HTML(http.StatusOK, "auth.tmpl", gin.H{})
}

func LoginHandler(c *gin.Context) {
	var form LoginFormType

	c.Bind(&form)

	log.Printf("%+v", form)

	token := jwt.New(jwt.SigningMethodHS256);
	token.Claims["username"] = form.Username
	token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	tokenString, err := token.SignedString(privateToken)

	if (err != nil) {
		log.Println("Fatal encode", err.Error())
	}

	c.JSON(http.StatusOK, gin.H{"token":tokenString})
}