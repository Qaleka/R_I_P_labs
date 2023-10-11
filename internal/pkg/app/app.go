package app

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"R_I_P_labs/internal/app/config"
	"R_I_P_labs/internal/app/ds"
	"R_I_P_labs/internal/app/dsn"
	"R_I_P_labs/internal/app/repository"
)

type Application struct {
	repo   *repository.Repository
	config *config.Config
	// dsn string
}

type GetRecipientsBack struct {
	Recipients []ds.Recipient
	Name       string
}

func (a *Application) Run() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/html/*")

	r.GET("/recipients", func(c *gin.Context) {
		FIO := c.Query("FIO")
		recipients, err := a.repo.GetRecipientByName(FIO)

		if err != nil {
			log.Println("cant get recipients", err)
			c.Error(err)
			return
		}
		c.HTML(http.StatusOK, "index.tmpl", GetRecipientsBack{
			Name:       FIO,
			Recipients: recipients,
		})

	})

	r.GET("/recipients/:id", func(c *gin.Context) {
		id := c.Param("id")
		recipient, err := a.repo.GetRecipientByID(id)
		if err != nil { // если не получилось
			log.Printf("cant get product by id %v", err)
			c.Error(err)
			return
		}

		c.HTML(http.StatusOK, "item.tmpl", *recipient)
	})

	r.POST("/recipients", func(c *gin.Context) {
		id := c.PostForm("delete")

		a.repo.DeleteRecipient(id)

		recipients, err := a.repo.GetRecipientByName("")
		if err != nil {
			log.Println("cant get recipients", err)
			c.Error(err)
			return
		}

		c.HTML(http.StatusOK, "index.tmpl", GetRecipientsBack{
			Name:       "",
			Recipients: recipients,
		})
	})

	r.Static("/image", "./resources")
	r.Static("/css", "./static/css")
	r.Run("localhost:9000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	log.Println("Server down")
}

func New() (*Application, error) {
	var err error
	app := Application{}
	app.config, err = config.NewConfig()
	if err != nil {
		return nil, err
	}

	app.repo, err = repository.New(dsn.FromEnv())
	if err != nil {
		return nil, err
	}

	return &app, nil
}
