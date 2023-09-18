package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// Post model
type Post struct {
	UserID int    `json:"userId"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

// ...

// GetPostsController akan menampilkan data yang diperoleh dari API
func GetPostsController(c echo.Context) error {
	// Membuat request GET ke API
	resp, err := http.Get("https://jsonplaceholder.typicode.com/posts")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Membaca respons dari API
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Mengubah data JSON ke dalam bentuk slice Post
	var posts []Post
	if err := json.Unmarshal(body, &posts); err != nil {
		return err
	}

	// Merender respons dalam bentuk JSON
	return c.JSON(http.StatusOK, posts)
}

// GetPostController akan menampilkan data dengan ID 3 dari API
func GetPostController(c echo.Context) error {
	// Membuat request GET ke API dengan ID
	id := c.Param("id")
	resp, err := http.Get("https://jsonplaceholder.typicode.com/posts/"+id)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Membaca respons dari API
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Mengubah data JSON ke dalam bentuk struct Post
	var post Post
	if err := json.Unmarshal(body, &post); err != nil {
		return err
	}

	// Merender respons dalam bentuk JSON
	return c.JSON(http.StatusOK, post)
}

// CreatePostController akan menyimpan data postingan ke server melalui API
func CreatePostController(c echo.Context) error {
	// Mendapatkan data postingan dari permintaan
	post := new(Post)
	if err := c.Bind(post); err != nil {
		return err
	}

	// Mengirim data postingan ke API dengan metode POST
	postBytes, err := json.Marshal(post)
	if err != nil {
		return err
	}

	resp, err := http.Post("https://jsonplaceholder.typicode.com/posts", "application/json", bytes.NewBuffer(postBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Membaca respons dari API
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Mengubah data JSON ke dalam bentuk struct Post
	var createdPost Post
	if err := json.Unmarshal(body, &createdPost); err != nil {
		return err
	}

	// Merender respons dalam bentuk JSON
	return c.JSON(http.StatusCreated, createdPost)
}

// DeletePostController akan menghapus data melalui API
func DeletePostController(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	// Mengirim permintaan DELETE ke API
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("https://jsonplaceholder.typicode.com/posts/%d", id), nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Merender respons dalam bentuk JSON
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Post deleted successfully",
	})
}

func main() {
	e := echo.New()

	// Routing
	e.GET("/posts", GetPostsController)                    // Menampilkan semua data postingan
	e.GET("/posts/:id", GetPostController)                  // Menampilkan data postingan dengan ID 3
	e.POST("/posts", CreatePostController)                  // Menyimpan data postingan
	e.DELETE("/posts/:id", DeletePostController)            // Menghapus data postingan dengan ID tertentu

	// Start server
	e.Logger.Fatal(e.Start(":8000"))
}
