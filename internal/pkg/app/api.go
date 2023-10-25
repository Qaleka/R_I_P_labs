package app

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"R_I_P_labs/internal/app/ds"
	"R_I_P_labs/internal/app/schemes"


	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"mime/multipart"
)

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