package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

const (
	ImgDir = "images"
)

type Items struct {
	Items []Item `json:"items"`
}

type Item struct {
	Name     string `json:"name"`
	Category string `json:"category"`
	Image    string `json:"image_name"`
}

type Response struct {
	Message string `json:"message"`
}

func root(c echo.Context) error {
	res := Response{Message: "Hello, world!"}
	return c.JSON(http.StatusOK, res)
}

func addItem(c echo.Context) error {
	var items Items
	// Get form data
	name := c.FormValue("name")
	category := c.FormValue("category")
	image, err := c.FormFile("image")
	if err != nil {
		return err
	}

	src, err := image.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Create a new SHA256 hash
	hash := sha256.New()

	// Copy the file content to the hash
	if _, err := io.Copy(hash, src); err != nil {
		return err
	}

	hashInBytes := hash.Sum(nil)

	// Convert hash bytes to hex string
	hashString := hex.EncodeToString(hashInBytes)

	image_jpg := hashString + ".jpg"

	item := Item{Name: name, Category: category, Image: image_jpg}

	c.Logger().Infof("Receive item: %s, %s", item.Name, item.Category, item.Image)
	message := fmt.Sprintf("item received: %s, %s, %s", item.Name, item.Category, item.Image)

	res := Response{Message: message}

	items.Items = append(items.Items, item)

	f, err := os.OpenFile("items.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	output, err := json.Marshal(&items)
	if err != nil {
		return err
	}

	_, err = f.Write(output)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

func getItems(c echo.Context) error {
	var items Items
	jsonBytes, err := os.ReadFile("items.json")
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonBytes, &items)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, items)
}

func getItemById(c echo.Context) error {
	var items Items
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}
	jsonBytes, err := os.ReadFile("items.json")
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonBytes, &items)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, items.Items[id-1])
}

func getImg(c echo.Context) error {
	// Create image path
	imgPath := path.Join(ImgDir, c.Param("imageFilename"))

	if !strings.HasSuffix(imgPath, ".jpg") {
		res := Response{Message: "Image path does not end with .jpg"}
		return c.JSON(http.StatusBadRequest, res)
	}
	if _, err := os.Stat(imgPath); err != nil {
		c.Logger().Debugf("Image not found: %s", imgPath) //log.DEBUGじゃないと表示されない
		imgPath = path.Join(ImgDir, "default.jpg")
	}
	//os.Stat
	return c.File(imgPath)
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	//e.Logger.SetLevel(log.INFO)
	e.Logger.SetLevel(log.DEBUG) //log levelをDEBUGに設定

	frontURL := os.Getenv("FRONT_URL")
	if frontURL == "" {
		frontURL = "http://localhost:3000"
	}
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{frontURL},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	// Routes
	e.GET("/", root)
	e.POST("/items", addItem)
	e.GET("/items", getItems)
	e.GET("/image/:imageFilename", getImg)
	e.GET("/items/:id", getItemById)

	// Start server
	e.Logger.Fatal(e.Start(":9000"))
}
