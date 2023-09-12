package routes

import (
	"context"
	"fmt"
	database "main/DataBase"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/olahol/go-imageupload"
	"go.mongodb.org/mongo-driver/bson"
)

var client, ctx = database.Connection()
var collection = database.Collection(client, ctx)

var currentImage *imageupload.Image

func GetHotels(c *gin.Context) {
	dbCtx, dbCancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer dbCancel()

	// get all elements as an cursor object from db
	cursor, err := collection.Find(dbCtx, bson.D{{}})

	// array for each element
	var results []bson.M

	// loop each element and extract the data
	for cursor.Next(dbCtx) {
		var result bson.M

		if err := cursor.Decode(&result); err != nil {
			return
		}
		results = append(results, result)
	}

	if err != nil {
		fmt.Println(err)
	}

	c.JSON(200, results)
}

func AddHotel(c *gin.Context) {

	name := c.Query("name")
	rate, _ := strconv.ParseFloat(c.Query("rate")[0:3], 64)
	price, err := strconv.Atoi(c.Query("price"))
	location := c.Query("location")
	images, err := c.FormFile("images")
	comments := c.Query("comments")

	fmt.Println(err)

	// should use query
	response := c.Request.URL.Query()

	for key, value := range response {
		fmt.Printf("Key: %s, Value: %d\n", key, value)
	}
	// fmt.Println(images)

	// for i := 0; i < len(images); i++ {
	// 	c.SaveUploadedFile(images[i], "temp%d.jpg")
	// }

	// imageByte, err := ioutil.ReadFile("temp.jpg")

	dbCtx, dbCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer dbCancel() // Make sure to cancel the context when the handler exits

	filter := bson.D{
		{"name", name},
		{"rate", rate},
		{"price", price},
		{"location", location},
		{"images", images},
		{"comments", comments},
	}

	res := database.InsertOne(dbCtx, collection, filter)

	c.String(200, "Done Uploading %v", res)
}

func GetSpecHotel(c *gin.Context) {
	name := c.Query("name")
	rate, _ := strconv.ParseFloat(c.Query("rate")[0:3], 64)
	price, err := strconv.Atoi(c.Query("price"))

	dbCtx, dbCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer dbCancel()

	if err != nil {
		fmt.Println(err)
	}

	filter := bson.D{
		{"name", name},
		{"rate", rate},
		{"price", price},
	}

	// res := collection.FindOne(dbCtx, bson.D{{Key: "name", Value: name}, {Key: "rate", Value: rate}, {Key: "price", Value: price}})

	var hotel bson.M

	if err := collection.FindOne(dbCtx, filter).Decode(&hotel); err != nil {
		c.JSON(400, gin.H{"message": "field"})
	}

	c.JSON(200, hotel)
}
