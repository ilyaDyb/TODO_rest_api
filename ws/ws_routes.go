package ws

import "github.com/gin-gonic/gin"

func RegisterWsRoutes(router *gin.Engine) {
	router.GET("/ws/:chatID/:username", WsHandler)
}