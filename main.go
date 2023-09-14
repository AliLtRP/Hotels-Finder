package main

import (
	controller "main/Controller"
	"main/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	route := gin.Default()

	// cors middleware
	route.Use(controller.CORS())

	// 80MB
	route.MaxMultipartMemory = 80 << 20

	// return all the hotels
	route.GET("/", routes.GetHotels)

	// get a specific hotel
	route.GET("/app/get/hotel", routes.GetSpecHotel)

	// add new hotel
	route.POST("/post", routes.AddHotel)

	//check if the user if auth or not
	route.GET("/user/auth", routes.Auth)

	// assign new user
	route.POST("/user/auth/newuser", routes.AddNewUser)

	// no route
	route.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"error": "Not found !",
		})

		return
	})

	route.Run()
}
