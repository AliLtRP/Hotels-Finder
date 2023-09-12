package main

import (
	"main/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	route := gin.Default()

	// return all the hotels
	route.GET("/", routes.GetHotels)

	// get a specific hotel
	route.GET("/app/get/hotel", routes.GetSpecHotel)

	// add new hotel
	route.POST("/post", routes.AddHotel)

	route.Run("localhost:3000")
}
