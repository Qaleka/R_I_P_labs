package api

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"R_I_P_labs/internal/models"
)

func StartServer() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/html/*")

	recipients := models.GetCardsInfo()

	r.GET("/recipients", func(c *gin.Context) {
		Name := c.Query("Name")
		filteredRecipients := filterRecipients(recipients, Name)

		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			// "Recipients": recipients,
			"Recipients": filteredRecipients,
			"Name":       Name,
		})
	})

	r.GET("/recipients/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id < 0 || id >= len(recipients) {
			c.String(http.StatusNotFound, "Страница не найдена")
			return
		}
		recipient := recipients[id]

		c.HTML(http.StatusOK, "item.tmpl", gin.H{
			"Recipients": recipient,
		})
	})

	r.Static("/image", "./resources")
	r.Static("/css", "./static/css")
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	log.Println("Server down")

	r.Run()
	log.Println("Server down")
}

func filterRecipients(recipients []models.Recipients, filter string) []models.Recipients {
	if filter == "" {
		return recipients
	}
	var filtered []models.Recipients
	for _, recipient := range recipients {
		nameParts := strings.Fields(filter)
		matches := false
		for _, part := range nameParts {
			if contains(recipient.Name.First_name, part) || contains(recipient.Name.Second_name, part) {
				matches = true
				break
			}
		}
		if matches {
			filtered = append(filtered, recipient)
		}
	}

	return filtered
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}
