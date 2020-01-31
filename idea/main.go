package main

import (
	"idea/rout"
	"idea/session"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("view/*.html")
	router.Static("/assets", "./assets")

	store := session.NewDummyStore()
	router.Use(session.StartDefaultSession(store))

	user := router.Group("/user")
	{
		user.POST("/signup", rout.UserSignUp)
		user.POST("/login", rout.UserLogIn)
		user.POST("/logout", rout.UserLogOut)
	}

	router.GET("/", rout.Home)
	router.GET("/login", rout.LogIn)
	router.GET("/signup", rout.SignUp)
	router.NoRoute(rout.NoRoute)

	router.Run(":8080")
}
