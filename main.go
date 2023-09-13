package main

import (
	"main/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	route := gin.Default()

	route.MaxMultipartMemory = 80 << 20

	// return all the hotels
	route.GET("/", routes.GetHotels)

	// get a specific hotel
	route.GET("/app/get/hotel", routes.GetSpecHotel)

	// add new hotel
	route.POST("/post", routes.AddHotel)

	//check if the user if auth or not
	route.GET("/user/auth", routes.Auth)

	route.POST("/user/auth/newuser", routes.AddNewUser)

	route.Run("localhost:3001")
}
