package main

import (
	"log"
	handler "sanjay/api/Handler"
	router "sanjay/api/Router"
	"sanjay/api/query"
	"sanjay/api/service"
	"sanjay/api/util"
	"sanjay/config"
	"sanjay/config/connection"
	"sanjay/config/migration"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	util.InitValidator()
	//redish connection
	conn, err := connection.NewRedisConnection()
	redishHandler := service.NewRadishImplement(conn)
	if err != nil {
		log.Fatalf("failed to create connection %v", err)
	}

	//twillow connection
	TwillCon, err := connection.NewTwillowConn()
	if err != nil {
		log.Fatalf("failed to create connection %v", err)
	}
	twillowHandler := service.NewTwillowClient(TwillCon)

	// Db connection
	Db, er := connection.NewDbConnection()
	if er != nil {
		log.Fatalf("failed to connect db %v", er)
	}
	_, err = Db.Exec(migration.CreateTable())
	if er != nil {
		log.Fatalf("failed to connect schima %v", er)
	}
	queryHandler := query.GetQueryHandler(Db)

	//attach Handler services
	handlerFunc := handler.NewHandlerAttach(twillowHandler, redishHandler, queryHandler)

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Authorization", "token"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	app := router.NewEngineRout(r, handlerFunc)
	app.Routes()

	if err := r.Run(config.LoadEnv().GetAppPort()); err != nil {
		log.Fatalf("failed to create server")
	}
}
