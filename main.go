package main

import (
	"games/controllers"
	"games/database"
	"games/middlewares"
	"games/models"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	err := database.InitDatabase()
	if err != nil {

		log.Fatalln("could not create database", err)
	}

	database.GlobalDB.AutoMigrate(&models.Player{})
	database.GlobalDB.AutoMigrate(&models.Account{})

	r := setupRouter()

	r.Run(":8080")
}

func setupRouter() *gin.Engine {

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(200, "Welcome To This Website")
	})

	api := r.Group("/api")
	{

		public := api.Group("/public")
		{

			public.POST("/login", controllers.Login)

			public.POST("/register", controllers.RegisterPlayer)

		}

		protected := api.Group("/protected").Use(middlewares.Authz())
		{

			protected.GET("/profile", controllers.Profile) // detail player
			protected.POST("/logout", controllers.Logout) // 
			protected.POST("/account", controllers.Account) // register account bank
			protected.PUT("/topupbalance", controllers.TopUpBalance)
			protected.GET("/getlistplayer", controllers.GetListPlayer) // return list player and filter
		}
	}

	return r
}
