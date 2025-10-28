package router

import (
	handler "sanjay/api/Handler"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Router  *gin.Engine
	Handler *handler.HandlerMeth
}

func NewEngineRout(rout *gin.Engine, handler *handler.HandlerMeth) *Config {
	return &Config{
		Router:  rout,
		Handler: handler,
	}
}
func (app *Config) Routes() {
	app.Router.POST("/otp", app.Handler.SendSMS())
	app.Router.POST("/verifyOTP", app.Handler.VerifyOTP())

	profile := app.Router.Group("/profileSection")
	profile.Use(handler.AuthMiddleware())
	profile.GET("", app.Handler.GetProfile())
	profile.PUT("", app.Handler.UpdateProfile())

	chatBox := app.Router.Group("/ws")
	chatBox.Use(handler.QueryAuthMiddleware())
	chatBox.GET("/chat", app.Handler.WsConnection())
	app.Router.GET("/connection", app.Handler.LiveConnection())
}
