package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
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
	_ "github.com/mattn/go-sqlite3"
)

const (
	ImgDir  = "images"
	DB_PATH = "../db/mercari.sqlite3"
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
	var categoryID int
	const getCategoryFromNameQuery = "SELECT id FROM categories WHERE name = $1"
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

	hashInBytes := hash.Sum(nil)

	// Convert hash bytes to hex string
	hashString := hex.EncodeToString(hashInBytes)

	image_jpg := hashString + ".jpg"

	new_image, err := os.Create("images/" + image_jpg)
	if err != nil {
		return err
	}

	// Copy the file content to the hash
	if _, err := io.Copy(new_image, src); err != nil {
		return err
	}

	item := Item{Name: name, Category: category, Image: image_jpg}

	c.Logger().Infof("Receive item: %s, %s", item.Name, item.Category, item.Image)
	message := fmt.Sprintf("item received: %s, %s, %s", item.Name, item.Category, item.Image)

	res := Response{Message: message}

	items.Items = append(items.Items, item)

	db, err := sql.Open("sqlite3", DB_PATH)
	if err != nil {
		return err
	}
	defer db.Close()

	row := db.QueryRow(getCategoryFromNameQuery, item.Category)
	err = row.Scan(&categoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err = db.Exec("INSERT INTO categories (name) VALUES ($1)", item.Category)
			if err != nil {
				return err
			}
			row := db.QueryRow(getCategoryFromNameQuery, item.Category)
			err = row.Scan(&categoryID)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	cmd2 := "INSERT INTO items (name, category_id, image_name) VALUES ($1, $2, $3)"
	_, err = db.Exec(cmd2, item.Name, categoryID, item.Image)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

func getItems(c echo.Context) error {
	var items Items

	db, err := sql.Open("sqlite3", DB_PATH)
	if err != nil {
		return err
	}
	defer db.Close()

	rows, err := db.Query("SELECT items.name, categories.name, items.image_name FROM items JOIN categories ON items.category_id = categories.id")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var name, category, image string
		if err := rows.Scan(&name, &category, &image); err != nil {
			return err
		}
		item := Item{Name: name, Category: category, Image: image}
		items.Items = append(items.Items, item)
	}

	return c.JSON(http.StatusOK, items)
}

func getItemById(c echo.Context) error {
	var items Items
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", DB_PATH)
	if err != nil {
		return err
	}
	defer db.Close()

	cmd := "SELECT items.name, categories.name, items.image_name FROM items JOIN categories ON items.category_id = categories.id WHERE items.id LIKE ?"
	rows, err := db.Query(cmd, id)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var name, category, image string
		if err := rows.Scan(&name, &category, &image); err != nil {
			return err
		}
		item := Item{Name: name, Category: category, Image: image}
		items.Items = append(items.Items, item)
	}

	return c.JSON(http.StatusOK, items.Items[id-1])
}

func searchItem(c echo.Context) error {
	var items Items
	keyword := c.FormValue("keyword")
	db, err := sql.Open("sqlite3", DB_PATH)
	if err != nil {
		return err
	}
	defer db.Close()

	cmd := "SELECT items.name, categories.name, items.image_name FROM items JOIN categories ON items.category_id = categories.id WHERE items.name LIKE ?"
	rows, err := db.Query(cmd, "%"+keyword+"%")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var name, category, image string
		if err := rows.Scan(&name, &category, &image); err != nil {
			return err
		}
		item := Item{Name: name, Category: category, Image: image}
		items.Items = append(items.Items, item)
	}
	return c.JSON(http.StatusOK, items)
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
	e.GET("/search", searchItem)

	// Start server
	e.Logger.Fatal(e.Start(":9000"))
}
