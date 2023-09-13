package routes

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	database "main/DataBase"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/olahol/go-imageupload"
	"go.mongodb.org/mongo-driver/bson"
)

var client, ctx = database.Connection()
var collectionHotel = database.Collection(client, ctx, "Hotel")
var collectionHAuth = database.Collection(client, ctx, "users")

var currentImage *imageupload.Image

func GetHotels(c *gin.Context) {
	dbCtx, dbCancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer dbCancel()

	// get all elements as an cursor object from db
	cursor, err := collectionHotel.Find(dbCtx, bson.D{{}})

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

	// get the params from url
	name := c.Query("name")
	rate, _ := strconv.ParseFloat(c.Query("rate")[0:3], 64)
	price, err := strconv.Atoi(c.Query("price"))
	location := c.Query("location")

	// get the comments array
	comments := c.Request.URL.Query()["comments"]

	fmt.Println(err)

	fmt.Println(comments)
	// get the images
	form, _ := c.MultipartForm()
	files := form.File["images"]

	// get the directory
	dist, err := os.Getwd()

	// save each image in ./images dir
	for _, file := range files {
		c.SaveUploadedFile(file, path.Join(dist, "images", file.Filename))
	}

	dbCtx, dbCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer dbCancel() // Make sure to cancel the context when the handler exits

	// files in ./images dir
	filesList, err := os.ReadDir(path.Join(dist, "images"))

	// 2nd array to store array of image bytes
	var imageByte [][]byte

	// loop into each image and read it bytes
	for _, file := range filesList {
		Bytes, err := ioutil.ReadFile(path.Join(dist, "images", file.Name()))

		if err != nil {
			log.Println(err)
		}

		imageByte = append(imageByte, Bytes)
	}

	// schema or the info i need to save it in database
	filter := bson.D{
		{"name", name},
		{"rate", rate},
		{"price", price},
		{"location", location},
		{"images", imageByte},
		{"comments", comments},
	}

	// insert data
	res := database.InsertOne(dbCtx, collectionHotel, filter)

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

	var hotel bson.M

	if err := collectionHotel.FindOne(dbCtx, filter).Decode(&hotel); err != nil {
		c.JSON(400, gin.H{"message": "field"})
	}

	c.JSON(200, hotel)
}

// auth function

func Auth(c *gin.Context) {

	user := c.Query("user")
	password := c.Query("password")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	userInfo := bson.D{
		{"user", user},
		{"passowrd", password},
	}

	if res := collectionHAuth.FindOne(ctx, userInfo); res != nil {
		c.String(http.StatusFound, "%v is auth", user)
		return
	}

	c.String(http.StatusNonAuthoritativeInfo, "%v is not auth", user)
}

// new user sign up
func AddNewUser(c *gin.Context) {

	userName := c.Query("username")
	passowrd := c.Request.URL.Query()["password"]
	phone := c.Query("phone")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	fmt.Println(passowrd[0], passowrd[1])

	if passowrd[0] != passowrd[1] {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "bad request, can't assign new user",
		})

		return
	}

	newUser := bson.D{
		{"name", userName},
		{"password", passowrd[0]},
		{"phone", phone},
	}

	if res := collectionHAuth.FindOne(ctx, bson.D{{"user", userName}}); res != nil {
		userResponse := database.InsertOne(ctx, collectionHAuth, newUser)

		c.JSON(http.StatusCreated, gin.H{
			"userID": userResponse.InsertedID,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "bad request, can't assign new user",
		})
	}

}
