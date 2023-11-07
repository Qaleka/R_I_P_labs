package app

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"R_I_P_labs/internal/app/config"
	"R_I_P_labs/internal/app/dsn"
	"R_I_P_labs/internal/app/repository"
)

type Application struct {
	repo        *repository.Repository
	minioClient *minio.Client
	config      *config.Config
	// dsn string
}

func (app *Application) Run() {
	log.Println("Server start up")

	r := gin.Default()

	r.Use(ErrorHandler())

	// Услуги - получатели
	r.GET("/api/recipients", app.GetAllRecipients)                                     // Список с поиском
	r.GET("/api/recipients/:recipient_id", app.GetRecipient)                           // Одна услуга
	r.DELETE("/api/recipients/:recipient_id", app.DeleteRecipient)              // Удаление
	r.PUT("/api/recipients/:recipient_id", app.ChangeRecipient)                 // Изменение
	r.POST("/api/recipients", app.AddRecipient)                                    // Добавление
	r.POST("/api/recipients/:recipient_id/add_to_notification", app.AddToNotification) // Добавление в заявку

	// Заявки - уведомления
	r.GET("/api/notifications", app.GetAllNotifications)                                                       // Список (отфильтровать по дате формирования и статусу)
	r.GET("/api/notifications/:notification_id", app.GetNotification)                                          // Одна заявка
	r.PUT("/api/notifications/:notification_id/update", app.UpdateNotification)                                // Изменение (добавление транспорта)
	r.DELETE("/api/notifications/:notification_id", app.DeleteNotification)                             //Удаление
	r.DELETE("/api/notifications/:notification_id/delete_recipient/:recipient_id", app.DeleteFromNotification) // Изменеие (удаление услуг)
	r.PUT("/api/notifications/:notification_id/user_confirm", app.UserConfirm)                                 // Сформировать создателем
	r.PUT("/api/notifications/:notification_id/moderator_confirm", app.ModeratorConfirm)                        // Завершить отклонить модератором

	r.Static("/image", "./resources")
	r.Static("/css", "./static/css")
	r.Run("localhost:7000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	log.Println("Server down")
}

func New() (*Application, error) {
	var err error
	loc, _ := time.LoadLocation("UTC")
	time.Local = loc
	app := Application{}
	app.config, err = config.NewConfig()
	if err != nil {
		return nil, err
	}

	app.repo, err = repository.New(dsn.FromEnv())
	if err != nil {
		return nil, err
	}

	app.minioClient, err = minio.New(app.config.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4("", "", ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	return &app, nil
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		for _, err := range c.Errors {
			log.Println(err.Err)
		}
		lastError := c.Errors.Last()
		if lastError != nil {
			switch c.Writer.Status() {
			case http.StatusBadRequest:
				c.JSON(-1, gin.H{"error": "wrong request"})
			case http.StatusNotFound:
				c.JSON(-1, gin.H{"error": lastError.Error()})
			case http.StatusMethodNotAllowed:
				c.JSON(-1, gin.H{"error": lastError.Error()})
			default:
				c.Status(-1)
			}
		}
	}
}
