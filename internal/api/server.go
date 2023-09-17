package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"lab1/internal/models"
)

func StartServer() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	recipients := []models.Recipients{
		{
			Id: 0,
			Name: models.FIO{
				First_name:  "Олег",
				Second_name: "Орлов",
				Third_name:  "Никитович",
			},
			ImageSrc: "https://i.pinimg.com/originals/10/ad/ab/10adabc386ba646f7df5f4e4d3156272.jpg",
			Email:    "OlegO@mail.ru",
			Age:      27,
			Adress:   "Москва, ул. Измайловская, д.13, кв.54",
		},
		{
			Id: 1,
			Name: models.FIO{
				First_name:  "Василий",
				Second_name: "Гречко",
				Third_name:  "Валентинович",
			},
			ImageSrc: "https://catherineasquithgallery.com/uploads/posts/2021-02/1614511031_164-p-na-belom-fone-chelovek-185.jpg",
			Email:    "Grechko_101@mail.ru",
			Age:      31,
			Adress:   "Москва, ул. Тверская, д.25, кв.145",
		},
	}

	// r.GET("/recipients", func(c *gin.Context) {
	// 	// Отображение списка услуг (получателей уведомлений) в виде карточек
	// 	c.HTML(http.StatusOK, "index.tmpl", gin.H{"Recipients": recipients})
	// })

	// r.GET("/recipients/:id", func(c *gin.Context) {
	// 	id := c.Param("Id")
	// 	// Отображение подробной информации об услуге по её ID
	// 	c.HTML(http.StatusOK, "item.tmpl", gin.H{"RecipientsID": id})
	// })
	r.GET("/recipients", func(c *gin.Context) {
		filter := c.Query("filter")
		field := c.Query("field")
		filteredRecipients := filterRecipients(recipients, filter, field)

		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			// "Recipients": recipients,
			"Recipients": filteredRecipients,
			"filter":     filter,
			"field":      field,
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

	r.Static("/image", "/resources/index.css")

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	log.Println("Server down")

	r.Run()
	log.Println("Server down")
}

func filterRecipients(recipients []models.Recipients, filter string, field string) []models.Recipients {
	if filter == "" {
		return recipients
	}
	var filtered []models.Recipients
	for _, recipient := range recipients {
		if field == "First_name" && contains(recipient.Name.First_name, filter) {
			filtered = append(filtered, recipient)
		} else if field == "Email" && contains(recipient.Email, filter) {
			filtered = append(filtered, recipient)
		}
	}

	return filtered
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}
