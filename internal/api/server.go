package api

import (
	"log"
	"models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func StartServer() {
	log.Println("Server start up")

	recipients := []models.Container{
		{
			Id:       0,
			Name:     "Андрей Петрович",
			ImageSrc: "http:localhost:8080/man1.png",
		},
	}

	r := gin.Default()

	r.GET("/recipients", func(c *gin.Context) {
		c.JSON(http.StatusOK, recipients)
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.LoadHTMLGlob("templates/*")

	r.GET("/home", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title":  "Йо",
			"image1": "C:/Users/ANDREY/Desktop/РИП/lab1/resources/man1.png",
		})
	})

	r.Static("/image", "./resources")

	r.Run()
	log.Println("Server down")
}
